package kluster

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// fakeDocker is a docker on PATH that logs each call to $FAKE_LOG; lists
// $FAKE_NODES for ps; saves an image as a file naming it; logs what an
// import read on its standard input; and exits 4 for the subcommand
// $FAKE_DOCKER_FAIL.
const fakeDocker = `#!/bin/sh
sub="$1"
if [ "$sub" = "$FAKE_DOCKER_FAIL" ]; then echo "docker: $sub broke" >&2; exit 4; fi
case "$sub" in
  ps) echo "docker $*" >> "$FAKE_LOG"; printf '%b' "$FAKE_NODES" ;;
  save) echo "docker save $4" >> "$FAKE_LOG"; echo "tar of $4" > "$3" ;;
  exec) echo "docker exec $3 $(cat) <- $4 $5 $6 $7 $8 $9" >> "$FAKE_LOG" ;;
esac
`

// fakeStep is a project script: it logs its name, its directory, and its
// arguments to $FAKE_LOG, writes the FLAI_ variables and KUBECONFIG it got to
// $FAKE_ENV/<name>, lists $FAKE_IMAGES_OUT in FLAI_KLUSTER_IMAGES when that is
// set, and exits 3 when $FAKE_FAIL names it.
const fakeStep = `#!/bin/sh
name=$(basename "$0" .sh)
echo "$name $(pwd) $*" >> "$FAKE_LOG"
env | grep -E '^(FLAI_|KUBECONFIG=)' | sort > "$FAKE_ENV/$name"
if [ -n "$FLAI_KLUSTER_IMAGES" ]; then printf '%b' "$FAKE_IMAGES_OUT" > "$FLAI_KLUSTER_IMAGES"; fi
echo "$name ran"
if [ "$name" = "$FAKE_FAIL" ]; then echo "$name broke" >&2; exit 3; fi
`

type deployFixture struct {
	fixture
	worktree, env string
	m             manifest.Manifest
}

const (
	commitA = "0123456789abcdef0123456789abcdef01234567"
	commitB = "fedcba9876543210fedcba9876543210fedcba98"
)

// deploySetup makes a project with a cluster that is up, a story's worktree
// with build and deploy scripts, and a fake docker first on PATH, with two
// nodes of the cluster and one other container.
func deploySetup(t *testing.T) deployFixture {
	t.Helper()
	f := deployFixture{fixture: setup(t, authstarConfig), worktree: t.TempDir(), env: t.TempDir()}
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "docker"), []byte(fakeDocker), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := os.MkdirAll(filepath.Join(f.worktree, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"build", "deploy"} {
		if err := os.WriteFile(filepath.Join(f.worktree, "scripts", s+".sh"), []byte(fakeStep), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(KubeconfigPath(f.root)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(KubeconfigPath(f.root), []byte("apiVersion: v1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_ENV", f.env)
	t.Setenv("FAKE_NODES", `k3d-authstar-integration-server-0 server\nk3d-authstar-integration-agent-0 agent\nk3d-authstar-integration-serverlb loadbalancer\n`)
	t.Setenv("FAKE_IMAGES_OUT", `authstar/tower:s-0578-0123456789ab\n\nauthstar/gate:s-0578-0123456789ab\n`)
	t.Setenv("FAKE_DOCKER_FAIL", "")
	f.m = manifest.Manifest{Name: "authstar", Kluster: &manifest.Kluster{
		Config: "integration/kluster.yaml",
		Build:  []string{"scripts/build.sh", "--all"},
		Deploy: []string{"scripts/deploy.sh"},
	}}
	return f
}

func (f deployFixture) envOf(t *testing.T, step string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.env, step))
	if err != nil {
		t.Fatalf("%s did not run: %v", step, err)
	}
	m := map[string]string{}
	for _, l := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		k, v, _ := strings.Cut(l, "=")
		m[k] = v
	}
	return m
}

func TestDeployBuildsLoadsAndDeploysIntoTheStorysNamespace(t *testing.T) {
	f := deploySetup(t)
	var out bytes.Buffer
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &out)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"build " + f.worktree + " --all",
		"docker ps --filter label=k3d.cluster=authstar-integration --format {{.Names}} {{.Label \"k3d.role\"}}",
		"docker save authstar/tower:s-0578-0123456789ab",
		"docker exec k3d-authstar-integration-server-0 tar of authstar/tower:s-0578-0123456789ab <- ctr -n k8s.io images import -",
		"docker exec k3d-authstar-integration-agent-0 tar of authstar/tower:s-0578-0123456789ab <- ctr -n k8s.io images import -",
		"docker save authstar/gate:s-0578-0123456789ab",
		"docker exec k3d-authstar-integration-server-0 tar of authstar/gate:s-0578-0123456789ab <- ctr -n k8s.io images import -",
		"docker exec k3d-authstar-integration-agent-0 tar of authstar/gate:s-0578-0123456789ab <- ctr -n k8s.io images import -",
		"deploy " + f.worktree,
	}
	if got := f.calls(t); !slices.Equal(got, want) {
		t.Errorf("build, then each image into each node, then deploy, in the worktree:\n got %q\nwant %q", got, want)
	}
	env := f.envOf(t, "deploy")
	wantEnv := map[string]string{
		"FLAI_STORY":             "S-0578",
		"FLAI_KLUSTER_NAME":      "authstar-integration",
		"FLAI_KLUSTER_NAMESPACE": "s-0578",
		"FLAI_KLUSTER_PREFIX":    "s0578_",
		"FLAI_KLUSTER_COMMIT":    commitA,
		"FLAI_KLUSTER_TAG":       "s-0578-0123456789ab",
		"FLAI_KLUSTER_STATE":     StatePath(f.root, "s-0578"),
		"KUBECONFIG":             KubeconfigPath(f.root),
		"FLAI_KLUSTER_CONTEXT":   "k3d-authstar-integration",
	}
	for k, v := range wantEnv {
		if env[k] != v {
			t.Errorf("deploy gets %s=%q, got %q", k, v, env[k])
		}
	}
	if _, ok := env["FLAI_KLUSTER_IMAGES"]; ok {
		t.Error("only build gets FLAI_KLUSTER_IMAGES")
	}
	if b := f.envOf(t, "build"); b["FLAI_KLUSTER_IMAGES"] == "" || b["FLAI_KLUSTER_TAG"] != "s-0578-0123456789ab" {
		t.Errorf("build gets the images file and the variables: %v", b)
	}
	if fi, err := os.Stat(StatePath(f.root, "s-0578")); err != nil || !fi.IsDir() {
		t.Errorf("the namespace's state folder is made: %v", err)
	}

	if r.Outcome != OutcomePassed || r.Story != "S-0578" || r.Namespace != "s-0578" || r.Commit != commitA || r.Tag != "s-0578-0123456789ab" || r.Checkout != f.worktree || r.EndedAt.IsZero() || r.FailedStep != "" {
		t.Errorf("the record says what was deployed and that it passed: %+v", r)
	}
	if !slices.Equal(r.Images, []string{"authstar/tower:s-0578-0123456789ab", "authstar/gate:s-0578-0123456789ab"}) {
		t.Errorf("the record names the images: %q", r.Images)
	}
	var names []string
	for _, s := range r.Steps {
		names = append(names, s.Name+"="+s.Outcome)
		if s.Log != LogPath(f.root, "s-0578", s.Name) {
			t.Errorf("step %s logs to %s, got %s", s.Name, LogPath(f.root, "s-0578", s.Name), s.Log)
		}
	}
	if !slices.Equal(names, []string{"build=passed", "load=passed", "deploy=passed"}) {
		t.Errorf("the record has the three steps: %q", names)
	}
	if data, _ := os.ReadFile(LogPath(f.root, "s-0578", StepBuild)); !strings.Contains(string(data), "build ran") {
		t.Errorf("build's output is in its log: %q", data)
	}
	if data, _ := os.ReadFile(LogPath(f.root, "s-0578", StepLoad)); !strings.Contains(string(data), "loading authstar/gate:s-0578-0123456789ab into k3d-authstar-integration-agent-0") {
		t.Errorf("the load is logged: %q", data)
	}
	if !strings.Contains(out.String(), "s-0578: deploy, log "+LogPath(f.root, "s-0578", StepDeploy)) {
		t.Errorf("each step is named as it starts: %q", out.String())
	}

	kept, ok, err := ReadRecord(f.root, "s-0578")
	if err != nil || !ok || kept.Outcome != OutcomePassed || kept.Commit != commitA || len(kept.Steps) != 3 {
		t.Errorf("the record is kept: %+v %v %v", kept, ok, err)
	}
	data, err := os.ReadFile(RecordPath(f.root, "s-0578"))
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"story": "S-0578"`, `"namespace": "s-0578"`, `"commit":`, `"images":`, `"started_at":`, `"ended_at":`, `"outcome": "passed"`, `"steps":`, `"log":`} {
		if !strings.Contains(string(data), field) {
			t.Errorf("the record's JSON has %s:\n%s", field, data)
		}
	}
	if _, ok := Running(f.root); ok {
		t.Error("no job runs once the deploy ends")
	}
}

func TestADeployStopsAtTheStepThatFails(t *testing.T) {
	for _, tc := range []struct {
		step, fail, dockerFail string
		code                   int
		ran                    string
	}{
		{step: StepBuild, fail: "build", code: 3, ran: "build"},
		{step: StepLoad, dockerFail: "save", code: 4, ran: "build,docker"},
		{step: StepDeploy, fail: "deploy", code: 3, ran: "build,docker,docker,docker,docker,docker,docker,docker,deploy"},
	} {
		t.Run(tc.step, func(t *testing.T) {
			f := deploySetup(t)
			t.Setenv("FAKE_FAIL", tc.fail)
			t.Setenv("FAKE_DOCKER_FAIL", tc.dockerFail)
			r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "in-progress", f.worktree, commitA, &bytes.Buffer{})
			var se *StepError
			if !errors.As(err, &se) || se.Step != tc.step {
				t.Fatalf("the %s step's failure is a StepError naming it: %v", tc.step, err)
			}
			log := LogPath(f.root, "s-0578", tc.step)
			if !strings.Contains(err.Error(), "the "+tc.step+" step") || !strings.Contains(err.Error(), log) {
				t.Errorf("the error names the step and its log: %v", err)
			}
			if se.ExitCode != tc.code || !strings.Contains(err.Error(), "exit code "+strconv.Itoa(tc.code)) {
				t.Errorf("the error names the exit code %d: %v", tc.code, err)
			}
			var ran []string
			for _, c := range f.calls(t) {
				ran = append(ran, strings.Fields(c)[0])
			}
			if strings.Join(ran, ",") != tc.ran {
				t.Errorf("nothing runs after the failing step: %q", ran)
			}
			if r.Outcome != OutcomeFailed || r.FailedStep != tc.step || r.Error == "" {
				t.Errorf("the record says the deploy failed at %s: %+v", tc.step, r)
			}
			kept, _, _ := ReadRecord(f.root, "s-0578")
			if kept.Outcome != OutcomeFailed || kept.FailedStep != tc.step || kept.ExitCode != se.ExitCode {
				t.Errorf("the kept record says so too: %+v", kept)
			}
			if last := kept.Steps[len(kept.Steps)-1]; last.Name != tc.step || last.Outcome != OutcomeFailed {
				t.Errorf("the failing step is the record's last, failed: %+v", kept.Steps)
			}
		})
	}
}

func TestALoadFailsWhenNoNodeRuns(t *testing.T) {
	f := deploySetup(t)
	t.Setenv("FAKE_NODES", `k3d-other-serverlb loadbalancer\n`)
	_, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "the load step") || !strings.Contains(err.Error(), "flai kluster up") {
		t.Errorf("a load with no node says to bring the cluster up: %v", err)
	}
}

func TestADeployRefusesAStoryNotInProgressOrReview(t *testing.T) {
	f := deploySetup(t)
	for _, status := range []string{"backlog", "ready", "done", "cancelled"} {
		_, err := f.c.Deploy(context.Background(), f.m, "S-0578", status, f.worktree, commitA, &bytes.Buffer{})
		if err == nil || !strings.Contains(err.Error(), "S-0578 is "+status) || !strings.Contains(err.Error(), "in progress or in review") {
			t.Errorf("a story that is %s is refused: %v", status, err)
		}
	}
	if calls := f.calls(t); len(calls) != 0 {
		t.Errorf("nothing runs for a refused story: %q", calls)
	}
	if _, ok, _ := ReadRecord(f.root, "s-0578"); ok {
		t.Error("a refused story is not recorded")
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "T-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "not a story ID") {
		t.Errorf("a task's ID is refused: %v", err)
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, "main", &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "whole hash") {
		t.Errorf("a commit that is not a hash is refused: %v", err)
	}
	if err := os.Remove(KubeconfigPath(f.root)); err != nil {
		t.Fatal(err)
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "flai kluster up") {
		t.Errorf("a cluster with no kubeconfig says to bring it up: %v", err)
	}
}

func TestADeployRefusesWhileAnotherJobRuns(t *testing.T) {
	f := deploySetup(t)
	started := time.Date(2026, 11, 7, 2, 0, 0, 0, time.UTC)
	release, err := acquire(f.root, Job{Kind: "promote", Namespace: "authstar", StartedAt: started, PID: 4242})
	if err != nil {
		t.Fatal(err)
	}
	if j, ok := Running(f.root); !ok || j.Kind != "promote" || j.Namespace != "authstar" {
		t.Errorf("the running job is seen: %+v %v", j, ok)
	}
	_, err = f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	if !errors.Is(err, ErrBusy) || !strings.Contains(err.Error(), "a promote of the namespace authstar") || !strings.Contains(err.Error(), "2026-11-07T02:00:00Z") || !strings.Contains(err.Error(), "4242") {
		t.Errorf("a second job is refused, naming the one that runs: %v", err)
	}
	if calls := f.calls(t); len(calls) != 0 {
		t.Errorf("nothing runs while another job holds the lock: %q", calls)
	}
	release()
	if _, ok := Running(f.root); ok {
		t.Error("no job runs once it is released")
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Errorf("a deploy runs once the lock is free: %v", err)
	}
}

func TestAJobFileLeftBehindIsNotRunning(t *testing.T) {
	f := deploySetup(t)
	if err := os.MkdirAll(Dir(f.root), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(jobPath(f.root), []byte(`{"kind":"deploy","namespace":"s-0001"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok := Running(f.root); ok {
		t.Error("a job file no process holds the lock for is not a running job")
	}
}

func TestADeployWithoutABuildLoadsNothing(t *testing.T) {
	f := deploySetup(t)
	f.m.Kluster.Build = nil
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if got := f.calls(t); !slices.Equal(got, []string{"deploy " + f.worktree}) {
		t.Errorf("only deploy runs: %q", got)
	}
	if len(r.Images) != 0 || r.Steps[0].Outcome != OutcomeSkipped || r.Steps[1].Outcome != OutcomeSkipped || r.Steps[2].Outcome != OutcomePassed {
		t.Errorf("build and load are skipped, and no image is recorded: %+v", r)
	}
}

func TestARedeployAtANewCommitReplacesTheRecord(t *testing.T) {
	f := deploySetup(t)
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_FAIL", "deploy")
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitB, &bytes.Buffer{}); err == nil {
		t.Fatal("the second deploy fails at deploy")
	}
	r, ok, err := ReadRecord(f.root, "s-0578")
	if err != nil || !ok {
		t.Fatal(ok, err)
	}
	if r.Commit != commitB || r.Tag != "s-0578-fedcba987654" || r.Outcome != OutcomeFailed || len(r.Steps) != 3 {
		t.Errorf("the record is the second deploy's alone: %+v", r)
	}
	if e := f.envOf(t, "deploy"); e["FLAI_KLUSTER_TAG"] != "s-0578-fedcba987654" {
		t.Errorf("the second deploy is tagged with its commit: %v", e)
	}
}

func TestTargets(t *testing.T) {
	tg, err := StoryTarget("s-0042", "/w", commitA)
	if err != nil || tg.Story != "S-0042" || tg.Namespace != "s-0042" || tg.Prefix != "s0042_" || tg.Tag() != "s-0042-0123456789ab" {
		t.Errorf("a story's target: %+v %v", tg, err)
	}
	m := manifest.Manifest{Name: "Auth Star", Kluster: &manifest.Kluster{}}
	d, err := DefaultTarget(m, "/main", commitB)
	if err != nil || d.Story != "" || d.Namespace != "auth-star" || d.Prefix != "" || d.Tag() != "main-fedcba987654" {
		t.Errorf("the default namespace's target: %+v %v", d, err)
	}
}

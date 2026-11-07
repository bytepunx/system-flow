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

// fakeKubectl is a kubectl on PATH that logs its arguments and the
// KUBECONFIG it got to $FAKE_LOG, and exits 5 when $FAKE_KUBECTL_FAIL is set.
const fakeKubectl = `#!/bin/sh
echo "kubectl $* KUBECONFIG=$KUBECONFIG" >> "$FAKE_LOG"
if [ -n "$FAKE_KUBECTL_FAIL" ]; then echo "kubectl broke" >&2; exit 5; fi
echo "namespace deleted"
`

// fakeGit is a git on PATH that logs its directory and arguments to
// $FAKE_LOG; answers rev-parse with $FAKE_HEAD, or fails it when
// $FAKE_GIT_FAIL is set; and makes the worktree worktree add names, with a
// .git file and a copy of $FAKE_SCRIPTS as its scripts.
const fakeGit = `#!/bin/sh
echo "git $(pwd) $*" >> "$FAKE_LOG"
case "$1" in
  rev-parse) if [ -n "$FAKE_GIT_FAIL" ]; then echo "fatal: bad revision" >&2; exit 128; fi; echo "$FAKE_HEAD" ;;
  worktree) if [ "$2" = "add" ]; then mkdir -p "$5"; echo "gitdir: x" > "$5/.git"; cp -r "$FAKE_SCRIPTS" "$5/scripts"; fi ;;
esac
`

// teardownSetup is deploySetup with fake kubectl and git on PATH and a
// teardown script in the story's worktree.
func teardownSetup(t *testing.T) deployFixture {
	t.Helper()
	f := deploySetup(t)
	bin := t.TempDir()
	for name, script := range map[string]string{"kubectl": fakeKubectl, "git": fakeGit} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(f.worktree, "scripts", "teardown.sh"), []byte(fakeStep), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_KUBECTL_FAIL", "")
	t.Setenv("FAKE_GIT_FAIL", "")
	t.Setenv("FAKE_HEAD", commitA)
	t.Setenv("FAKE_SCRIPTS", filepath.Join(f.worktree, "scripts"))
	return f
}

// reset forgets the calls logged so far.
func (f deployFixture) reset(t *testing.T) {
	t.Helper()
	if err := os.Remove(f.log); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func (f deployFixture) deployed(t *testing.T, story string) {
	t.Helper()
	if _, err := f.c.Deploy(context.Background(), f.m, story, "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
}

const deleteS0578 = "kubectl --context k3d-authstar-integration delete namespace s-0578 --wait --ignore-not-found KUBECONFIG="

func TestATeardownRunsTheTeardownCommandThenDeletesTheNamespace(t *testing.T) {
	f := teardownSetup(t)
	f.deployed(t, "S-0578")
	f.reset(t)
	f.m.Kluster.Teardown = []string{"scripts/teardown.sh", "--purge"}
	r, err := f.c.Teardown(context.Background(), f.m, "S-0578", f.worktree, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"teardown " + f.worktree + " --purge", deleteS0578 + KubeconfigPath(f.root)}
	if got := f.calls(t); !slices.Equal(got, want) {
		t.Errorf("the teardown command, then the namespace deleted:\n got %q\nwant %q", got, want)
	}
	env := f.envOf(t, "teardown")
	if env["FLAI_KLUSTER_NAMESPACE"] != "s-0578" || env["FLAI_KLUSTER_PREFIX"] != "s0578_" || env["FLAI_KLUSTER_COMMIT"] != commitA || env["FLAI_STORY"] != "S-0578" {
		t.Errorf("the teardown command gets the namespace's variables at the recorded commit: %v", env)
	}
	if r.Namespace != "s-0578" || r.Commit != commitA {
		t.Errorf("the record torn down is returned: %+v", r)
	}
	for _, p := range []string{RecordPath(f.root, "s-0578"), StatePath(f.root, "s-0578"), LogDir(f.root, "s-0578")} {
		if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("%s is removed once the namespace is torn down: %v", p, err)
		}
	}
	if _, ok := Running(f.root); ok {
		t.Error("no job runs once the teardown ends")
	}
}

func TestATeardownWithoutATeardownCommandDeletesOnlyTheNamespace(t *testing.T) {
	f := teardownSetup(t)
	f.deployed(t, "S-0578")
	f.reset(t)
	if _, err := f.c.Teardown(context.Background(), f.m, "s-0578", "", &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if got := f.calls(t); !slices.Equal(got, []string{deleteS0578 + KubeconfigPath(f.root)}) {
		t.Errorf("only the namespace is deleted: %q", got)
	}
	if _, ok, _ := ReadRecord(f.root, "s-0578"); ok {
		t.Error("the record is removed")
	}
}

func TestATeardownOfAStoryNeverDeployedRunsNothing(t *testing.T) {
	f := teardownSetup(t)
	_, err := f.c.Teardown(context.Background(), f.m, "S-0578", f.worktree, &bytes.Buffer{})
	if !errors.Is(err, ErrNotDeployed) || !strings.Contains(err.Error(), "S-0578") {
		t.Errorf("a story never deployed says so: %v", err)
	}
	if calls := f.calls(t); len(calls) != 0 {
		t.Errorf("nothing runs: %q", calls)
	}
}

func TestATeardownThatFailsKeepsTheRecord(t *testing.T) {
	for _, tc := range []struct {
		name, fail, kubectl string
		code                int
		calls               int
	}{
		{name: "command", fail: "teardown", code: 3, calls: 1},
		{name: "kubectl", kubectl: "yes", code: 5, calls: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := teardownSetup(t)
			f.deployed(t, "S-0578")
			f.reset(t)
			f.m.Kluster.Teardown = []string{"scripts/teardown.sh"}
			t.Setenv("FAKE_FAIL", tc.fail)
			t.Setenv("FAKE_KUBECTL_FAIL", tc.kubectl)
			r, err := f.c.Teardown(context.Background(), f.m, "S-0578", f.worktree, &bytes.Buffer{})
			var se *StepError
			if !errors.As(err, &se) || se.Step != StepTeardown || se.ExitCode != tc.code || !strings.Contains(err.Error(), LogPath(f.root, "s-0578", StepTeardown)) {
				t.Fatalf("the failure names the teardown step, its exit code %d, and its log: %v", tc.code, err)
			}
			if n := len(f.calls(t)); n != tc.calls {
				t.Errorf("nothing runs after the failure: %q", f.calls(t))
			}
			kept, ok, _ := ReadRecord(f.root, "s-0578")
			if !ok || kept.Outcome != OutcomeFailed || kept.FailedStep != StepTeardown || kept.ExitCode != tc.code || kept.Commit != commitA {
				t.Errorf("the record is kept, failed at the teardown: %+v", kept)
			}
			if last := kept.Steps[len(kept.Steps)-1]; last.Name != StepTeardown || last.Outcome != OutcomeFailed {
				t.Errorf("the teardown is the record's last step: %+v", kept.Steps)
			}
			if r.FailedStep != StepTeardown {
				t.Errorf("the record returned says so too: %+v", r)
			}
			if _, err := os.Stat(StatePath(f.root, "s-0578")); err != nil {
				t.Errorf("the state folder stays: %v", err)
			}
		})
	}
}

func TestPromoteDeploysMainsHeadIntoTheDefaultNamespace(t *testing.T) {
	f := teardownSetup(t)
	wt := MainWorktree(f.root)
	r, err := f.c.Promote(context.Background(), f.m, "main", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	calls := f.calls(t)
	want := []string{
		"git " + f.root + " rev-parse --verify main^{commit}",
		"git " + f.root + " worktree prune",
		"git " + f.root + " worktree add --detach --force " + wt + " " + commitA,
		"build " + wt + " --all",
	}
	if len(calls) < len(want)+1 || !slices.Equal(calls[:len(want)], want) || calls[len(calls)-1] != "deploy "+wt {
		t.Errorf("the worktree of main's head is made, then built and deployed there:\n got %q\nwant %q then deploy", calls, want)
	}
	env := f.envOf(t, "deploy")
	if env["FLAI_STORY"] != "" || env["FLAI_KLUSTER_NAMESPACE"] != "authstar" || env["FLAI_KLUSTER_PREFIX"] != "" || env["FLAI_KLUSTER_TAG"] != "main-0123456789ab" || env["FLAI_KLUSTER_COMMIT"] != commitA {
		t.Errorf("deploy gets the default namespace's variables: %v", env)
	}
	if r.Outcome != OutcomePassed || r.Story != "" || r.Namespace != "authstar" || r.Checkout != wt || r.Commit != commitA {
		t.Errorf("the promotion is recorded under the default namespace: %+v", r)
	}

	f.reset(t)
	t.Setenv("FAKE_HEAD", commitB)
	r, err = f.c.Promote(context.Background(), f.m, "main", &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if got := f.calls(t); len(got) < 2 || got[1] != "git "+f.root+" -C "+wt+" checkout --detach --force "+commitB {
		t.Errorf("a second promotion moves the worktree to the new head: %q", got)
	}
	if r.Commit != commitB || r.Tag != "main-fedcba987654" {
		t.Errorf("the record is the new head's: %+v", r)
	}
}

func TestAPromotionThatFails(t *testing.T) {
	f := teardownSetup(t)
	t.Setenv("FAKE_FAIL", "deploy")
	r, err := f.c.Promote(context.Background(), f.m, "main", &bytes.Buffer{})
	var se *StepError
	if !errors.As(err, &se) || se.Step != StepDeploy || se.Namespace != "authstar" {
		t.Fatalf("a failed deploy of main names the step: %v", err)
	}
	if kept, _, _ := ReadRecord(f.root, "authstar"); kept.Outcome != OutcomeFailed || kept.FailedStep != StepDeploy || r.FailedStep != StepDeploy {
		t.Errorf("the failure is recorded under the default namespace: %+v", kept)
	}

	g := teardownSetup(t)
	t.Setenv("FAKE_FAIL", "")
	t.Setenv("FAKE_GIT_FAIL", "yes")
	if _, err := g.c.Promote(context.Background(), g.m, "trunk", &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "trunk") || !strings.Contains(err.Error(), "bad revision") {
		t.Errorf("a branch git cannot resolve is named, with git's words: %v", err)
	}
	for _, c := range g.calls(t) {
		if !strings.HasPrefix(c, "git ") {
			t.Errorf("nothing but git runs when the head cannot be resolved: %q", c)
		}
	}
	if _, ok := Running(g.root); ok {
		t.Error("the lock is released")
	}
}

func TestDeploymentsListsTheDefaultNamespaceFirst(t *testing.T) {
	f := teardownSetup(t)
	if list, err := Deployments(f.root, "authstar"); err != nil || list == nil || len(list) != 0 {
		t.Errorf("no records is an empty list: %v %v", list, err)
	}
	f.deployed(t, "S-0578")
	f.deployed(t, "S-0012")
	f.deployed(t, "S-1000")
	if _, err := f.c.Promote(context.Background(), f.m, "main", &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	list, err := Deployments(f.root, "authstar")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, d := range list {
		got = append(got, d.Namespace)
		if d.Job != nil {
			t.Errorf("no job runs: %+v", d)
		}
	}
	if !slices.Equal(got, []string{"authstar", "s-0012", "s-0578", "s-1000"}) {
		t.Errorf("the default namespace first, then the stories by number: %q", got)
	}
	if !list[0].Default || list[1].Default {
		t.Errorf("only the default namespace is marked default: %+v", list[:2])
	}

	release, err := acquire(f.root, Job{Kind: "teardown", Story: "S-0578", Namespace: "s-0578", StartedAt: time.Now().UTC(), PID: 7})
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	stale, _, _ := ReadRecord(f.root, "s-0012")
	stale.Outcome, stale.EndedAt = OutcomeRunning, time.Time{}
	if err := writeRecord(f.root, stale); err != nil {
		t.Fatal(err)
	}
	list, err = Deployments(f.root, "authstar")
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range list {
		switch d.Namespace {
		case "s-0578":
			if d.Job == nil || d.Job.Kind != "teardown" {
				t.Errorf("the namespace a job acts on carries it: %+v", d)
			}
		case "s-0012":
			if d.Outcome != OutcomeInterrupted || d.Job != nil {
				t.Errorf("a record left running with no job is interrupted: %+v", d)
			}
		default:
			if d.Job != nil {
				t.Errorf("a namespace no job acts on carries none: %+v", d)
			}
		}
	}
}

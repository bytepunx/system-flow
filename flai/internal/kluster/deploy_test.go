package kluster

import (
	"bytes"
	"context"
	"encoding/json"
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
// $FAKE_ENV/<name>, copies the record $FAKE_RECORD names, when set, to
// $FAKE_ENV/<name>.record, lists acme/<component> in FLAI_KLUSTER_IMAGES for
// a component's build and $FAKE_IMAGES_OUT for the block's build, and exits
// 3 when $FAKE_FAIL names it.
const fakeStep = `#!/bin/sh
name=$(basename "$0" .sh)
echo "$name $(pwd) $*" >> "$FAKE_LOG"
env | grep -E '^(FLAI_|KUBECONFIG=)' | sort > "$FAKE_ENV/$name"
if [ -n "$FAKE_RECORD" ]; then cp "$FAKE_RECORD" "$FAKE_ENV/$name.record"; fi
if [ -n "$FLAI_KLUSTER_COMPONENT" ]; then echo "acme/$FLAI_KLUSTER_COMPONENT:$FLAI_KLUSTER_TAG" > "$FLAI_KLUSTER_IMAGES"
elif [ -n "$FLAI_KLUSTER_IMAGES" ]; then printf '%b' "$FAKE_IMAGES_OUT" > "$FLAI_KLUSTER_IMAGES"; fi
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
			var after bool
			for _, s := range kept.Steps {
				switch {
				case s.Name == tc.step:
					after = true
					if s.Outcome != OutcomeFailed {
						t.Errorf("the failing step is recorded failed: %+v", s)
					}
				case after && s.Outcome != OutcomeSkipped:
					t.Errorf("a step after the failing one is recorded skipped: %+v", s)
				case !after && s.Outcome != OutcomePassed:
					t.Errorf("a step before the failing one is recorded passed: %+v", s)
				}
			}
			if len(kept.Steps) != 3 {
				t.Errorf("the record keeps every step of the plan: %+v", kept.Steps)
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

// componentSetup is deploySetup with three components, a, b, and c, each
// built by its own script, scripts/build-<name>.sh, in place of the block's
// build.
func componentSetup(t *testing.T) deployFixture {
	t.Helper()
	f := deploySetup(t)
	f.m.Kluster.Build = nil
	for _, name := range []string{"a", "b", "c"} {
		script := filepath.Join(f.worktree, "scripts", "build-"+name+".sh")
		if err := os.WriteFile(script, []byte(fakeStep), 0o755); err != nil {
			t.Fatal(err)
		}
		f.m.Kluster.Components = append(f.m.Kluster.Components, manifest.KlusterComponent{Name: name, Build: []string{"scripts/build-" + name + ".sh", name}})
	}
	return f
}

func stepNames(steps []Step) []string {
	var names []string
	for _, s := range steps {
		names = append(names, s.Name+"="+s.Outcome)
	}
	return names
}

// S-0595: each component builds as its own step, build-<name>, with its own
// log and images file and FLAI_KLUSTER_COMPONENT, and one load takes every
// image they listed (ADR-0263).
func TestADeployBuildsEachComponentAsItsOwnStep(t *testing.T) {
	f := componentSetup(t)
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := stepNames(r.Steps), []string{"build-a=passed", "build-b=passed", "build-c=passed", "load=passed", "deploy=passed"}; !slices.Equal(got, want) {
		t.Errorf("one build step per component, then load and deploy:\n got %q\nwant %q", got, want)
	}
	for _, name := range []string{"a", "b", "c"} {
		step := ComponentStep(name)
		data, err := os.ReadFile(LogPath(f.root, "s-0578", step))
		if err != nil || !strings.Contains(string(data), step+" ran") {
			t.Errorf("%s's output is in its own log: %q %v", step, data, err)
		}
		env := f.envOf(t, step)
		if env["FLAI_KLUSTER_COMPONENT"] != name || env["FLAI_KLUSTER_IMAGES"] != filepath.Join(LogDir(f.root, "s-0578"), "images-"+name) {
			t.Errorf("%s gets its component and its own images file: %v", step, env)
		}
	}
	if _, ok := f.envOf(t, "deploy")["FLAI_KLUSTER_COMPONENT"]; ok {
		t.Error("only a component's build gets FLAI_KLUSTER_COMPONENT")
	}
	want := []string{"acme/a:s-0578-0123456789ab", "acme/b:s-0578-0123456789ab", "acme/c:s-0578-0123456789ab"}
	if !slices.Equal(r.Images, want) {
		t.Errorf("the record names every component's images: %q", r.Images)
	}
	var saved []string
	for _, c := range f.calls(t) {
		if strings.HasPrefix(c, "docker save ") {
			saved = append(saved, strings.TrimPrefix(c, "docker save "))
		}
	}
	if !slices.Equal(saved, want) {
		t.Errorf("one load takes every component's images: %q", saved)
	}
}

// S-0595: a component whose build fails stops the deploy, and the steps after
// it are recorded skipped.
func TestAComponentThatFailsSkipsTheStepsAfterIt(t *testing.T) {
	f := componentSetup(t)
	t.Setenv("FAKE_FAIL", "build-b")
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	var se *StepError
	if !errors.As(err, &se) || se.Step != "build-b" || se.ExitCode != 3 {
		t.Fatalf("build-b's failure is a StepError naming it: %v", err)
	}
	if got, want := stepNames(r.Steps), []string{"build-a=passed", "build-b=failed", "build-c=skipped", "load=skipped", "deploy=skipped"}; !slices.Equal(got, want) {
		t.Errorf("the steps after the failing component are skipped:\n got %q\nwant %q", got, want)
	}
	kept, _, _ := ReadRecord(f.root, "s-0578")
	if kept.FailedStep != "build-b" || kept.Outcome != OutcomeFailed || !slices.Equal(stepNames(kept.Steps), stepNames(r.Steps)) {
		t.Errorf("the kept record says so too: %+v", kept)
	}
	var ran []string
	for _, c := range f.calls(t) {
		ran = append(ran, strings.Fields(c)[0])
	}
	if !slices.Equal(ran, []string{"build-a", "build-b"}) {
		t.Errorf("nothing runs after build-b: %q", ran)
	}
}

// S-0595: the record lists every step as pending before the first runs, and
// the step that runs as running, with its start and its log, while it runs
// (ADR-0263).
func TestTheRecordListsEveryStepPendingAndTheOneThatRuns(t *testing.T) {
	f := componentSetup(t)
	t.Setenv("FAKE_RECORD", RecordPath(f.root, "s-0578"))
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	read := func(step string) Record {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(f.env, step+".record"))
		if err != nil {
			t.Fatalf("%s saw no record: %v", step, err)
		}
		var r Record
		if err := json.Unmarshal(data, &r); err != nil {
			t.Fatal(err)
		}
		return r
	}
	first := read("build-a")
	if got, want := stepNames(first.Steps), []string{"build-a=running", "build-b=pending", "build-c=pending", "load=pending", "deploy=pending"}; !slices.Equal(got, want) {
		t.Errorf("while the first step runs:\n got %q\nwant %q", got, want)
	}
	if s := first.Steps[0]; s.StartedAt.IsZero() || s.Log != LogPath(f.root, "s-0578", "build-a") || !s.EndedAt.IsZero() {
		t.Errorf("the running step has its start and its log and no end: %+v", s)
	}
	if s := first.Steps[1]; !s.StartedAt.IsZero() || s.Log != "" {
		t.Errorf("a pending step has no start and no log: %+v", s)
	}
	if first.Outcome != OutcomeRunning {
		t.Errorf("the deployment runs: %+v", first)
	}
	if got, want := stepNames(read("deploy").Steps), []string{"build-a=passed", "build-b=passed", "build-c=passed", "load=passed", "deploy=running"}; !slices.Equal(got, want) {
		t.Errorf("while deploy runs:\n got %q\nwant %q", got, want)
	}
	data, err := os.ReadFile(filepath.Join(f.env, "build-a.record"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), `"started_at"`) != 2 {
		t.Errorf("only the deployment and the running step have a start in the JSON:\n%s", data)
	}
}

// S-0595: the block's single build keeps its step, build, in the plan.
func TestASingleBuildIsPlannedAsBuildLoadDeploy(t *testing.T) {
	f := deploySetup(t)
	t.Setenv("FAKE_RECORD", RecordPath(f.root, "s-0578"))
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(f.env, "build.record"))
	if err != nil {
		t.Fatal(err)
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if got, want := stepNames(r.Steps), []string{"build=running", "load=pending", "deploy=pending"}; !slices.Equal(got, want) {
		t.Errorf("while build runs:\n got %q\nwant %q", got, want)
	}
	if e := f.envOf(t, "build"); e["FLAI_KLUSTER_IMAGES"] != filepath.Join(LogDir(f.root, "s-0578"), "images") || e["FLAI_KLUSTER_COMPONENT"] != "" {
		t.Errorf("the block's build keeps its images file and names no component: %v", e)
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

func TestTheChecksOutcomeIsKeptInTheRecordUntilTheNextDeploy(t *testing.T) {
	f := deploySetup(t)
	if err := RecordChecks(f.root, "s-0578", commitA, OutcomePassed); err == nil || !strings.Contains(err.Error(), "no deployment record") {
		t.Errorf("checks kept for a namespace never deployed: %v", err)
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if err := RecordChecks(f.root, "s-0578", commitB, OutcomePassed); err == nil || !strings.Contains(err.Error(), "deployed again, at 0123456789ab") {
		t.Errorf("checks at another commit than the deploy's were kept: %v", err)
	}
	if err := RecordChecks(f.root, "s-0578", commitA, OutcomeFailed); err != nil {
		t.Fatal(err)
	}
	r, _, err := ReadRecord(f.root, "s-0578")
	if err != nil || r.ChecksCommit != commitA || r.ChecksOutcome != OutcomeFailed || r.Outcome != OutcomePassed || len(r.Steps) != 3 {
		t.Errorf("the record keeps the deploy and the checks' outcome: %+v %v", r, err)
	}
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if r, _, _ := ReadRecord(f.root, "s-0578"); r.ChecksCommit != "" || r.ChecksOutcome != "" {
		t.Errorf("a deploy clears the checks' outcome: %+v", r)
	}
}

func TestAcceptGate(t *testing.T) {
	deployed := Record{Story: "S-0578", Namespace: "s-0578", Commit: commitA, Outcome: OutcomePassed}
	with := func(edit func(*Record)) *Record {
		r := deployed
		edit(&r)
		return &r
	}
	cases := []struct {
		name     string
		record   *Record
		checksOn bool
		want     []string // the reason holds each; none when the gate passes
	}{
		{"deploy and checks passed at the head", with(func(r *Record) { r.ChecksCommit, r.ChecksOutcome = commitA, OutcomePassed }), true, nil},
		{"checks off", &deployed, false, nil},
		{"checks off, their failure aside", with(func(r *Record) { r.ChecksCommit, r.ChecksOutcome = commitA, OutcomeFailed }), false, nil},
		{"no record", nil, true, []string{"s-0578 was never deployed"}},
		{"deploy at another commit", with(func(r *Record) {
			r.Commit, r.ChecksCommit, r.ChecksOutcome = commitB, commitB, OutcomePassed
		}), true, []string{"last deploy", "at fedcba987654", "not at the story's head 0123456789ab", "it passed"}},
		{"deploy failed", with(func(r *Record) { r.Outcome, r.FailedStep = OutcomeFailed, StepDeploy }), false, []string{"the deploy of S-0578", "at 0123456789ab failed at the deploy step"}},
		{"deploy failed before its first step", with(func(r *Record) { r.Outcome = OutcomeFailed }), true, []string{"the deploy of S-0578", "failed before its first step"}},
		{"deploy interrupted", with(func(r *Record) { r.Outcome = OutcomeRunning }), true, []string{"the deploy of S-0578", "was interrupted"}},
		{"teardown failed", with(func(r *Record) { r.Outcome, r.FailedStep = OutcomeFailed, StepTeardown }), true, []string{"the teardown of S-0578", "failed at the teardown step"}},
		{"checks failed", with(func(r *Record) { r.ChecksCommit, r.ChecksOutcome = commitA, OutcomeFailed }), true, []string{"deploy of S-0578", "at 0123456789ab passed", "checks run against it there ended failed"}},
		{"checks timed out", with(func(r *Record) { r.ChecksCommit, r.ChecksOutcome = commitA, "timed-out" }), true, []string{"checks run against it there ended timed-out"}},
		{"checks not run", &deployed, true, []string{"deploy of S-0578", "at 0123456789ab passed", "no checks run against it has ended there yet"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()
			if c.record != nil {
				if err := writeRecord(root, *c.record); err != nil {
					t.Fatal(err)
				}
			}
			ok, reason, err := AcceptGate(root, "s-0578", commitA, c.checksOn)
			if err != nil {
				t.Fatal(err)
			}
			if ok != (c.want == nil) {
				t.Fatalf("ok = %v, reason %q", ok, reason)
			}
			if ok && reason != "" {
				t.Errorf("a gate that passes gives the reason %q", reason)
			}
			for _, w := range c.want {
				if !strings.Contains(reason, w) {
					t.Errorf("reason %q does not say %q", reason, w)
				}
			}
		})
	}
}

func TestAcceptGateRefusesAStoryWithNoBranchOrNoID(t *testing.T) {
	root := t.TempDir()
	if ok, reason, err := AcceptGate(root, "S-0578", "", true); ok || err != nil || !strings.Contains(reason, "no branch story/S-0578") {
		t.Errorf("no head: %v %q %v", ok, reason, err)
	}
	if _, _, err := AcceptGate(root, "main", commitA, true); err == nil {
		t.Error("a name that is no story ID was gated")
	}
}

// headsGit is an execx.Runner whose git rev-parse answers head, or fails
// when head is empty, and records the arguments it was given.
type headsGit struct {
	head string
	args []string
}

func (g *headsGit) Run(_, _ string, args ...string) (string, error) {
	g.args = args
	if g.head == "" {
		return "", errors.New("exit status 1")
	}
	return g.head + "\n", nil
}

func (g *headsGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (g *headsGit) LookPath(name string) (string, error) { return name, nil }

func TestStoryHeadResolvesTheStorysBranch(t *testing.T) {
	g := &headsGit{head: commitA}
	head, err := StoryHead(g, "/main", "s-0578")
	if err != nil || head != commitA {
		t.Fatalf("head = %q, %v", head, err)
	}
	if want := []string{"rev-parse", "--verify", "--quiet", "refs/heads/story/S-0578^{commit}"}; !slices.Equal(g.args, want) {
		t.Errorf("git %q, want %q", g.args, want)
	}
	g.head = ""
	if _, err := StoryHead(g, "/main", "S-0578"); err == nil || !strings.Contains(err.Error(), "refs/heads/story/S-0578") {
		t.Errorf("a story with no branch: %v", err)
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

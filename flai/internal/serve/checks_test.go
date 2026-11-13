//go:build !windows

package serve

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func checksLab(t *testing.T) (dir Dir, worktree string) {
	t.Helper()
	root := t.TempDir()
	dir = Dir(filepath.Join(root, "serve"))
	worktree = filepath.Join(root, "worktree")
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir, worktree
}

func nowFixed() time.Time { return time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC) }

func TestRunChecksStopsAtTheFirstFailure(t *testing.T) {
	dir, worktree := checksLab(t)
	marker := filepath.Join(worktree, "second-ran")
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{
			{Name: "one", Command: []string{"sh", "-c", "exit 1"}},
			{Name: "two", Command: []string{"sh", "-c", "touch " + marker}},
		},
		Timeout: time.Minute,
	}
	run, err := dir.RunChecks("S-0001", "/root", worktree, cfg, nowFixed)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "failed" || len(run.Steps) != 1 || run.Steps[0].Name != "one" {
		t.Fatalf("run: %+v", run)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("the second check ran after the first failed")
	}
}

func TestRunChecksSubstitutesStoryAndRootAndRunsInTheWorktree(t *testing.T) {
	dir, worktree := checksLab(t)
	out := filepath.Join(worktree, "out.txt")
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{
			{Name: "one", Command: []string{"sh", "-c", "pwd -P > " + out + "; echo {story} >> " + out + "; echo {root} >> " + out}},
		},
		Timeout: time.Minute,
	}
	run, err := dir.RunChecks("S-0002", "/root", worktree, cfg, nowFixed)
	if err != nil || run.Outcome != "passed" {
		t.Fatalf("run: %+v %v", run, err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := filepath.EvalSymlinks(worktree)
	if err != nil {
		t.Fatal(err)
	}
	want := resolved + "\nS-0002\n" + worktree + "\n"
	if string(data) != want {
		t.Errorf("substitution: got %q, want %q", data, want)
	}
}

func TestRunChecksRefusesWithoutAWorktree(t *testing.T) {
	dir, _ := checksLab(t)
	cfg := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "one", Command: []string{"true"}}}, Timeout: time.Minute}
	_, err := dir.RunChecks("S-0003", "/root", filepath.Join(t.TempDir(), "nope"), cfg, nowFixed)
	if err == nil {
		t.Fatal("expected a refusal")
	}
	if !strings.Contains(err.Error(), "no worktree") || !strings.Contains(err.Error(), "flai stream open S-0003") {
		t.Errorf("wording: %v", err)
	}
}

// One run per story at a time (S-0082's third acceptance criterion): a
// second run for the same story is refused, not queued and not replaced,
// while the first is still active.
func TestRunChecksRefusesASecondRunWhileOneIsActive(t *testing.T) {
	dir, worktree := checksLab(t)
	releaseFile := filepath.Join(worktree, "release")
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{
			{Name: "one", Command: []string{"sh", "-c", "while [ ! -f " + releaseFile + " ]; do sleep 0.05; done"}},
		},
		Timeout: time.Minute,
	}
	done := make(chan ChecksRun, 1)
	go func() {
		run, err := dir.RunChecks("S-0007", "/root", worktree, cfg, time.Now)
		if err != nil {
			t.Error(err)
		}
		done <- run
	}()
	waitForRunning(t, dir, "S-0007")

	if _, err := dir.RunChecks("S-0007", "/root", worktree, cfg, time.Now); err == nil {
		t.Error("a second run for the same story should be refused")
	} else if !strings.Contains(err.Error(), "already active") {
		t.Errorf("wording: %v", err)
	}

	if err := os.WriteFile(releaseFile, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if run := <-done; run.Outcome != "passed" {
		t.Errorf("the first run: %+v", run)
	}
}

func waitForRunning(t *testing.T, dir Dir, story string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if run, ok := dir.ChecksState(story); ok && run.Running && run.PID != 0 {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("%s never started running", story)
}

func TestRunChecksTimesOutAndLeavesNothingRunning(t *testing.T) {
	dir, worktree := checksLab(t)
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{{Name: "slow", Command: []string{"sleep", "30"}}},
		Timeout:  300 * time.Millisecond,
	}
	start := time.Now()
	run, err := dir.RunChecks("S-0004", "/root", worktree, cfg, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "timed-out" {
		t.Fatalf("outcome: %+v", run)
	}
	if elapsed := time.Since(start); elapsed > 6*time.Second {
		t.Errorf("took too long to actually stop: %s", elapsed)
	}
}

func TestCancelChecksStopsARunningProcessAndWhatItSpawned(t *testing.T) {
	dir, worktree := checksLab(t)
	pidFile := filepath.Join(worktree, "child.pid")
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{
			{Name: "parent", Command: []string{"sh", "-c", "sleep 30 & echo $! > " + pidFile + "; wait"}},
		},
		Timeout: time.Minute,
	}
	done := make(chan ChecksRun, 1)
	go func() {
		run, err := dir.RunChecks("S-0005", "/root", worktree, cfg, time.Now)
		if err != nil {
			t.Error(err)
		}
		done <- run
	}()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(pidFile); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("the child never started: %v", err)
	}
	childPID := parsePID(t, pidData)

	if err := dir.CancelChecks("S-0005", 5*time.Second, func(d time.Duration) { time.Sleep(d) }); err != nil {
		t.Fatal(err)
	}
	run := <-done
	if run.Outcome != "cancelled" {
		t.Errorf("outcome: %+v", run)
	}
	if Alive(childPID) {
		t.Error("the grandchild (sleep) it spawned is still running; group-kill did not reach it")
	}
}

func TestCancelChecksRefusesWhenNothingIsRunning(t *testing.T) {
	dir, _ := checksLab(t)
	if err := dir.CancelChecks("S-0006", time.Second, func(time.Duration) {}); err == nil {
		t.Fatal("expected a refusal")
	}
}

func TestResolveChecksPrefersTheHostOverTheManifestNeverBoth(t *testing.T) {
	host := []manifest.NamedCommand{{Name: "host", Command: []string{"true"}}}
	man := []manifest.NamedCommand{{Name: "manifest", Command: []string{"true"}}}
	if got := ResolveChecks(nil, man); len(got) != 1 || got[0].Name != "manifest" {
		t.Errorf("falls back to the manifest: %+v", got)
	}
	if got := ResolveChecks(host, man); len(got) != 1 || got[0].Name != "host" {
		t.Errorf("the host's own, not merged: %+v", got)
	}
}

func parsePID(t *testing.T, data []byte) int {
	t.Helper()
	n, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("not a pid: %q", data)
	}
	return n
}

func TestRunChecksWithGivesEachCheckTheEnvironment(t *testing.T) {
	t.Setenv("FLAI_KLUSTER_NAMESPACE", "")
	dir, worktree := checksLab(t)
	out := filepath.Join(worktree, "ns.txt")
	cfg := ChecksConfig{
		Commands: []manifest.NamedCommand{{Name: "ns", Command: []string{"sh", "-c", `printf %s "$FLAI_KLUSTER_NAMESPACE" > ` + out}}},
		Timeout:  time.Minute,
	}
	run, err := dir.RunChecksWith(t.Context(), "S-0003", "/root", worktree, cfg, []string{"FLAI_KLUSTER_NAMESPACE=s-0003"}, nowFixed)
	if err != nil || run.Outcome != "passed" {
		t.Fatalf("run: %+v %v", run, err)
	}
	if data, _ := os.ReadFile(out); string(data) != "s-0003" {
		t.Errorf("the check saw FLAI_KLUSTER_NAMESPACE %q, want s-0003", data)
	}
	if log, _ := os.ReadFile(run.Log); !strings.Contains(string(log), "against the kluster namespace s-0003") {
		t.Errorf("the run's log does not name the namespace:\n%s", log)
	}

	run, err = dir.RunChecks("S-0004", "/root", worktree, cfg, nowFixed)
	if err != nil || run.Outcome != "passed" {
		t.Fatalf("run: %+v %v", run, err)
	}
	if data, _ := os.ReadFile(out); string(data) != "" {
		t.Errorf("a run with no environment added saw FLAI_KLUSTER_NAMESPACE %q", data)
	}
}

// ADR-0292 §6: a run keeps who asked for it, and the namespace and commit
// its variables name, in the story's checks status, where the next story
// reads it; a run from a shell names no asker, and one with no namespace
// names none.
func TestRunChecksAsKeepsWhoAskedTheNamespaceAndTheCommit(t *testing.T) {
	dir, worktree := checksLab(t)
	cfg := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "ok", Command: []string{"true"}}}, Timeout: time.Minute}
	env := []string{"FLAI_KLUSTER_NAMESPACE=s-0005", "FLAI_KLUSTER_COMMIT=" + klHeadA}
	run, err := dir.RunChecksAs(t.Context(), "agent-S-0005", "S-0005", "/root", worktree, cfg, env, nil, nowFixed)
	if err != nil || run.AskedBy != "agent-S-0005" || run.Namespace != "s-0005" || run.Commit != klHeadA {
		t.Fatalf("run: %+v %v", run, err)
	}
	kept, ok := dir.ChecksState("S-0005")
	if !ok || kept.AskedBy != "agent-S-0005" || kept.Namespace != "s-0005" || kept.Commit != klHeadA || kept.Outcome != "passed" {
		t.Errorf("the checks status = %+v, want who asked, the namespace, and the commit", kept)
	}
	run, err = dir.RunChecks("S-0006", "/root", worktree, cfg, nowFixed)
	if err != nil || run.AskedBy != "" || run.Namespace != "" || run.Commit != "" {
		t.Errorf("a run from a shell with no namespace: %+v %v", run, err)
	}
}

// withChecks has the lab's pass run cfg's checks, with their state in a
// folder of the test's own, as flai serve's options name them.
func (lab *klusterLab) withChecks(cfg ChecksConfig) {
	lab.t.Helper()
	lab.withConfig()
	lab.opts.Dir = Dir(lab.t.TempDir())
	lab.opts.Checks = func(*workitem.Repo) (ChecksConfig, error) { return cfg, nil }
	lab.pass = newKlusterPass(lab.opts, Entry{Key: "t", Name: "t", Root: lab.root}, &sync.WaitGroup{})
}

// nsCheck is a check that writes the namespace it sees to a file, and exits
// with exit.
func nsCheck(t *testing.T, exit int) (ChecksConfig, string) {
	out := filepath.Join(t.TempDir(), "ns.txt")
	return ChecksConfig{
		Enabled:  true,
		Commands: []manifest.NamedCommand{{Name: "golden", Command: []string{"sh", "-c", `printf %s "$FLAI_KLUSTER_NAMESPACE" > ` + out + "; exit " + strconv.Itoa(exit)}}},
		Timeout:  time.Minute,
	}, out
}

func TestThePassRunsAStorysChecksAgainstItsNamespaceOnceItsDeployPassed(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	cfg, out := nsCheck(t, 0)
	lab.withChecks(cfg)
	story := lab.inReview("Checked", klHeadA)

	lab.round()
	if data, _ := os.ReadFile(out); string(data) != strings.ToLower(story) {
		t.Fatalf("the checks saw the namespace %q, want %q", data, strings.ToLower(story))
	}
	var details []string
	for _, e := range lab.entries() {
		details = append(details, e.Outcome+": "+e.Detail)
	}
	if got := strings.Join(details, "\n"); !strings.Contains(got, "started: checks of "+story) || !strings.Contains(got, "done: checks of "+story) {
		t.Errorf("the journal does not hold the checks run's start and end:\n%s", got)
	}
	if run, ok := lab.opts.Dir.ChecksState(story); !ok || run.Outcome != "passed" {
		t.Errorf("checks state = %+v, %v; want a passed run", run, ok)
	}
}

func TestThePassRunsNoChecksWithChecksOffOrAfterAFailedDeploy(t *testing.T) {
	t.Run("checks off", func(t *testing.T) {
		lab := newKlusterLab(t, klusterBlock)
		cfg, out := nsCheck(t, 0)
		cfg.Enabled = false
		lab.withChecks(cfg)
		lab.inReview("Unchecked", klHeadA)
		lab.round()
		if len(lab.kluster.deploys()) != 1 {
			t.Fatalf("deploys = %v, want one", lab.kluster.deploys())
		}
		if _, err := os.Stat(out); err == nil {
			t.Error("checks ran with the checks host action off")
		}
	})
	t.Run("failed deploy", func(t *testing.T) {
		lab := newKlusterLab(t, klusterBlock)
		cfg, out := nsCheck(t, 0)
		lab.withChecks(cfg)
		lab.kluster.fail = func(string, string) error {
			return &kluster.StepError{Step: "deploy", ExitCode: 1, Err: errors.New("exit status 1")}
		}
		lab.inReview("Broken", klHeadA)
		lab.round()
		if _, err := os.Stat(out); err == nil {
			t.Error("checks ran after a failed deploy")
		}
	})
}

// fakeDockerClient writes a docker client into a folder of its own that
// records its arguments, one a line, in argv beside it, and the value of
// FLAI_KLUSTER_NAMESPACE it was given in env, then runs body.
func fakeDockerClient(t *testing.T, body string) (program, argv, env string) {
	t.Helper()
	bin := t.TempDir()
	program, argv, env = filepath.Join(bin, "docker"), filepath.Join(bin, "argv"), filepath.Join(bin, "env")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + argv + "\nprintf '%s\\n' \"$FLAI_KLUSTER_NAMESPACE\" > " + env + "\n" + body + "\n"
	if err := os.WriteFile(program, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return program, argv, env
}

// reopened is lab's project read again, with the container harness its
// manifest now names.
func reopened(t *testing.T, lab *verifyLab) *workitem.Repo {
	t.Helper()
	repo, err := workitem.Open(lab.repo.Root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// stopRecorder is the docker a check's container is stopped through: it
// records each call and runs nothing.
type stopRecorder struct {
	mu    sync.Mutex
	calls [][]string
}

func (r *stopRecorder) Run(_, name string, args ...string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, append([]string{name}, args...))
	return "", nil
}

func (r *stopRecorder) RunInput(dir, name, _ string, args ...string) (string, error) {
	return r.Run(dir, name, args...)
}

func (*stopRecorder) LookPath(name string) (string, error) { return name, nil }

func (r *stopRecorder) said() [][]string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.calls)
}

// I-0295, ADR-0274: a check of a story whose agent runs in a container runs
// as docker run in a container from the story's image with the story's
// mounts, its program the entrypoint and the run's environment passed by
// name; nothing of the check runs on the host. Before ADR-0274 the check's
// program ran on the host in the worktree.
func TestAContainerStorysCheckRunsInItsContainerAndNothingOnTheHost(t *testing.T) {
	lab, _ := newContainerLab(t, true)
	in := ChecksContainer(reopened(t, lab), lab.story, AgentConfig{}, imageDocker{has: true})
	if in == nil || in.Err != nil {
		t.Fatalf("the container story's checks have no container: %+v", in)
	}
	program, argvFile, envFile := fakeDockerClient(t, "exit 0")
	in.Host.Program = program
	marker := filepath.Join(lab.worktree, "ran-on-the-host")
	cfg := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "unit", Command: []string{"sh", "-c", "touch " + marker}}}, Timeout: time.Minute}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	run, err := dir.RunChecksIn(context.Background(), lab.story, lab.repo.Root, lab.worktree, cfg, []string{"FLAI_KLUSTER_NAMESPACE=s-0001"}, in, nowFixed)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "passed" {
		t.Errorf("the run: %+v", run)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("the check's program ran on the host")
	}
	data, err := os.ReadFile(argvFile)
	if err != nil {
		t.Fatalf("docker was not run: %v", err)
	}
	root := lab.repo.Root
	reg := filepath.Join(root, ".git", "worktrees", lab.story)
	whole := " " + strings.ReplaceAll(strings.TrimSpace(string(data)), "\n", " ") + " "
	for _, want := range []string{
		" run --rm --name flai-check-" + lab.story + "-20260922T120000Z-1 ",
		" -v " + lab.worktree + ":" + lab.worktree + " -w " + lab.worktree + " ",
		" -v " + filepath.Join(root, ".git") + ":" + filepath.Join(root, ".git") + ":ro ",
		" -v " + reg + ":" + reg + " ",
		" -v " + filepath.Join(root, "wip") + ":" + filepath.Join(root, "wip") + ":ro ",
		" -v " + filepath.Join(root, "system-flow.yaml") + ":" + filepath.Join(root, "system-flow.yaml") + ":ro ",
		" -e FLAI_KLUSTER_NAMESPACE ",
		" --entrypoint sh " + harness.DefaultImages[harness.ClaudeCode] + " -c touch " + marker + " ",
	} {
		if !strings.Contains(whole, want) {
			t.Errorf("the check's docker run lacks %q:\n%s", want, whole)
		}
	}
	if strings.Contains(whole, "s-0001") {
		t.Error("a variable's value is in the docker run's arguments")
	}
	if ns, _ := os.ReadFile(envFile); strings.TrimSpace(string(ns)) != "s-0001" {
		t.Errorf("the docker client's FLAI_KLUSTER_NAMESPACE: %q", ns)
	}
}

// ADR-0058, ADR-0274: a cancelled check of a container story stops its
// container by name, beside ending the docker client's process group.
func TestACancelledContainerCheckStopsItsContainerByName(t *testing.T) {
	lab, _ := newContainerLab(t, true)
	in := ChecksContainer(reopened(t, lab), lab.story, AgentConfig{}, imageDocker{has: true})
	if in == nil || in.Err != nil {
		t.Fatalf("the container story's checks have no container: %+v", in)
	}
	program, _, _ := fakeDockerClient(t, "exec sleep 30")
	in.Host.Program = program
	stops := &stopRecorder{}
	in.Docker = stops
	cfg := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "slow", Command: []string{"sleep", "60"}}}, Timeout: time.Minute}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	done := make(chan ChecksRun, 1)
	go func() {
		run, err := dir.RunChecksIn(context.Background(), lab.story, lab.repo.Root, lab.worktree, cfg, nil, in, nowFixed)
		if err != nil {
			t.Error(err)
		}
		done <- run
	}()
	waitForRunning(t, dir, lab.story)
	if err := dir.CancelChecks(lab.story, 10*time.Second, func(d time.Duration) { time.Sleep(d) }); err != nil {
		t.Fatal(err)
	}
	if run := <-done; run.Outcome != "cancelled" {
		t.Errorf("the run: %+v", run)
	}
	want := []string{program, "stop", "flai-check-" + lab.story + "-20260922T120000Z-1"}
	if calls := stops.said(); len(calls) != 1 || !slices.Equal(calls[0], want) {
		t.Errorf("the stops: %q, want %q", calls, want)
	}
}

// ADR-0274 §5: a container story whose image is not on the host fails its
// check with the fix, and runs nothing on the host.
func TestAContainerStorysCheckWithNoImageFailsAndRunsNothingOnTheHost(t *testing.T) {
	lab, _ := newContainerLab(t, false)
	in := ChecksContainer(reopened(t, lab), lab.story, AgentConfig{}, imageDocker{has: false})
	if in == nil || in.Err == nil {
		t.Fatalf("a story with no image has a container: %+v", in)
	}
	marker := filepath.Join(lab.worktree, "ran-on-the-host")
	cfg := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "unit", Command: []string{"sh", "-c", "touch " + marker}}}, Timeout: time.Minute}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	run, err := dir.RunChecksIn(context.Background(), lab.story, lab.repo.Root, lab.worktree, cfg, nil, in, nowFixed)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "failed" {
		t.Errorf("the run: %+v", run)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("the check's program ran on the host")
	}
	if log, _ := os.ReadFile(run.Log); !strings.Contains(string(log), "flai image ensure") {
		t.Errorf("the log does not name the fix:\n%s", log)
	}
}

// A story whose agent runs on a host harness has no container, so its checks
// run on the host, in the worktree, as before.
func TestAHostHarnessStorysChecksHaveNoContainer(t *testing.T) {
	lab := newVerifyLab(t, false)
	if in := ChecksContainer(lab.repo, lab.story, AgentConfig{}, imageDocker{has: true}); in != nil {
		t.Errorf("a host-harness story's checks run in a container: %+v", in)
	}
}

// S-0751 criterion 1, ADR-0295 §5: flai serve's checks of a container story
// whose image is missing ask for the image job that makes it rather than
// refuse, wait for it before their first check, and then run in the image it
// made. Before S-0751 the checks failed, naming make agent-images.
func TestAContainerStorysCheckWaitsForTheImageJob(t *testing.T) {
	lab, _ := newContainerLab(t, false)
	made := &atomic.Bool{}
	jobs := verifyImageJobs(t, made, nil)
	program, argvFile, _ := fakeDockerClient(t, "exit 0")
	cfg := AgentConfig{Harnesses: map[string]harness.Host{harness.Container: {Program: program}}}
	in := checksContainerOf(reopened(t, lab), lab.story, cfg, madeDocker{made: made}, jobs, "t")
	if in == nil || in.Err != nil || in.wait == nil {
		t.Fatalf("the checks of a story with no image do not wait for an image job: %+v", in)
	}
	job := in.wait.job
	checks := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "unit", Command: []string{"true"}}}, Timeout: time.Minute}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	run, err := dir.RunChecksIn(context.Background(), lab.story, lab.repo.Root, lab.worktree, checks, nil, in, nowFixed)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "passed" || !made.Load() {
		t.Errorf("the run: %+v, image made %v", run, made.Load())
	}
	if _, err := os.Stat(argvFile); err != nil {
		t.Errorf("the check did not run in its container: %v", err)
	}
	log, _ := os.ReadFile(run.Log)
	for _, want := range []string{"waiting for flai serve's image job " + job.ID(), job.Log(), "the image job " + job.ID() + " passed"} {
		if !strings.Contains(string(log), want) {
			t.Errorf("the log lacks %q:\n%s", want, log)
		}
	}
}

// S-0751 criterion 4: checks whose image job fails fail with the job's cause,
// its failing step, and its log, run nothing on the host, and never name make
// agent-images.
func TestAContainerStorysCheckFailsWithTheImageJobsCause(t *testing.T) {
	lab, _ := newContainerLab(t, false)
	step := "build flai-agent-claude-code:dev"
	jobs := verifyImageJobs(t, &atomic.Bool{}, &agentimage.StepError{Step: step, Tries: 1, Said: "no space left on device", Err: errors.New("exit status 1")})
	in := checksContainerOf(reopened(t, lab), lab.story, AgentConfig{}, imageDocker{has: false}, jobs, "t")
	if in == nil || in.wait == nil {
		t.Fatalf("the checks of a story with no image do not wait for an image job: %+v", in)
	}
	job := in.wait.job
	marker := filepath.Join(lab.worktree, "ran-on-the-host")
	checks := ChecksConfig{Commands: []manifest.NamedCommand{{Name: "unit", Command: []string{"sh", "-c", "touch " + marker}}}, Timeout: time.Minute}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	run, err := dir.RunChecksIn(context.Background(), lab.story, lab.repo.Root, lab.worktree, checks, nil, in, nowFixed)
	if err != nil {
		t.Fatal(err)
	}
	if run.Outcome != "failed" {
		t.Errorf("the run: %+v", run)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Error("the check's program ran on the host")
	}
	log, _ := os.ReadFile(run.Log)
	for _, want := range []string{job.ID(), step, "no space left on device", job.Log(), "flai image ensure", "nothing runs on the host"} {
		if !strings.Contains(string(log), want) {
			t.Errorf("the log lacks %q:\n%s", want, log)
		}
	}
	if strings.Contains(string(log), "make agent-images") {
		t.Errorf("the log names make agent-images:\n%s", log)
	}
}

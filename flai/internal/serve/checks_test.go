//go:build !windows

package serve

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
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
			{Name: "one", Command: []string{"sh", "-c", "pwd > " + out + "; echo {story} >> " + out + "; echo {root} >> " + out}},
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

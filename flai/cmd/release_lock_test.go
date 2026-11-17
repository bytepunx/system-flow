package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// heldLockRunner is the real runner with another git process holding the
// main checkout's index lock throughout: before every git add -A asked of it,
// .git/index.lock is made where git makes it.
type heldLockRunner struct {
	execx.System
	lock string
}

func (h *heldLockRunner) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && len(args) > 1 && args[0] == "add" && args[1] == "-A" {
		if err := os.WriteFile(h.lock, nil, 0o644); err != nil {
			return "", err
		}
	}
	return h.System.Run(dir, name, args...)
}

// acceptedPendingProject is pendingProject with both its stories accepted, so
// a release of cli waits to be published, and a change left uncommitted in the
// main checkout, as the orchestrator's narrative is when it publishes, so the
// publish has something to commit before it tags.
func acceptedPendingProject(t *testing.T) (root, remote string) {
	t.Helper()
	root, remote = pendingProject(t)
	for _, id := range []string{"S-0001", "S-0002"} {
		if _, errOut, code := runIn(t, root, "accept", id); code != 0 {
			t.Fatalf("accept %s: %s", id, errOut)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "wip", "agents", "orchestrator.md"), []byte("# Orchestrator\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, remote
}

// I-0359: a publish whose git add and git commit meet the index lock another
// git process holds for a moment waits it out, then commits, tags, and pushes.
func TestReleasePendingWaitsOutAnIndexLockHeldForAMoment(t *testing.T) {
	t.Parallel()
	root, remote := acceptedPendingProject(t)
	lock := filepath.Join(root, ".git", "index.lock")
	r := &lockingRunner{lock: lock, runs: map[string]int{}}

	out, errOut, code := runWithApp(t, &app{configPath: configOf(t), cwd: root, runner: r}, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.1.0") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("the publish waits out the lock: %d %s %s", code, out, errOut)
	}
	if r.runs["add"] != 2 || r.runs["commit"] != 2 {
		t.Errorf("git add and git commit each meet the lock once and run again: %v", r.runs)
	}
	if subject := gitIn(t, root, "log", "-1", "--format=%s"); subject != "chore: publish cli 1.1.0" {
		t.Errorf("the publish commit: %q", subject)
	}
	if status := gitIn(t, root, "status", "--porcelain"); status != "" {
		t.Errorf("everything is committed: %q", status)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the release is tagged: %s", tags)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(remoteTags, "cli/v1.1.0") {
		t.Errorf("the tag reaches the remote: %s", remoteTags)
	}
	if local, head := gitIn(t, root, "rev-parse", "HEAD"), gitIn(t, remote, "rev-parse", "main"); local != head {
		t.Errorf("main reaches the remote too: %s vs %s", local, head)
	}
	if _, err := os.Stat(lock); !os.IsNotExist(err) {
		t.Errorf("the lock is the other process's to remove, and it did: %v", err)
	}
}

// I-0359: a publish whose git add meets the index lock throughout fails with
// git's message and says what to do, and tags and pushes nothing.
func TestReleasePendingFailsWhenTheIndexLockStaysHeld(t *testing.T) {
	t.Parallel()
	root, remote := acceptedPendingProject(t)
	lock := filepath.Join(root, ".git", "index.lock")
	remoteHead := gitIn(t, remote, "rev-parse", "main")

	_, errOut, code := runWithApp(t, &app{configPath: configOf(t), cwd: root, runner: &heldLockRunner{lock: lock}}, "release", "--pending")
	if code == 0 {
		t.Fatal("the publish whose lock stays held exits non-zero")
	}
	for _, want := range []string{"index.lock': File exists", "another git process held the index lock"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the error names %q:\n%s", want, errOut)
		}
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("nothing is tagged: %s", tags)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list", "cli/*"); strings.Contains(remoteTags, "cli/v1.1.0") {
		t.Errorf("nothing reaches the remote: %s", remoteTags)
	}
	if head := gitIn(t, remote, "rev-parse", "main"); head != remoteHead {
		t.Errorf("main is not pushed: %s vs %s", head, remoteHead)
	}
}

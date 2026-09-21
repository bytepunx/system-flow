package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// pendingProject is a git project like researchProject's, with a second
// story ready to accept, so a batch of two items against cli can build up
// before a publish.
func pendingProject(t *testing.T) (root, remote string) {
	t.Helper()
	root, remote = researchProject(t, "remediation", true)
	if _, errOut, code := runIn(t, root, "story", "new", "Another fix", "--epic", "E-0001", "--nature", "feature"); code != 0 {
		t.Fatal(errOut)
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0002-another-fix.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0002", "ready"}, {"move", "S-0002", "in-progress"}, {"stream", "open", "S-0002"},
		{"task", "new", "Do it", "--story", "S-0002"},
		{"move", "T-0002", "ready"}, {"move", "T-0002", "in-progress"}, {"move", "T-0002", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0002")
	_ = os.WriteFile(filepath.Join(wt, "cli", "another.go"), []byte("package main\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0002] the fix")
	if _, errOut, code := runIn(t, root, "move", "S-0002", "review"); code != 0 {
		t.Fatal(errOut)
	}
	return root, remote
}

// S-0087: several stories accepted against one component release together at
// the highest level among them, only when the operator publishes.
func TestReleasePendingEndToEnd(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 { // remediation: patch
		t.Fatalf("accept S-0001: %s", errOut)
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0002"); code != 0 { // feature: minor
		t.Fatalf("accept S-0002: %s", errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
		t.Fatalf("accepting alone releases nothing: %s", tags)
	}

	dry, _, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(dry, "cli") || !strings.Contains(dry, "S-0001") || !strings.Contains(dry, "S-0002") || !strings.Contains(dry, "dry run: nothing changed") {
		t.Fatalf("dry run names the component and every item: %s", dry)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.1.0") {
		t.Fatalf("dry run changes nothing: %s", tags)
	}

	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.1.0") || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("publish: %d %s %s", code, out, errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); !strings.Contains(tags, "cli/v1.1.0") {
		t.Errorf("the highest level among the batch (feature, minor) wins: %s", tags)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(remoteTags, "cli/v1.1.0") {
		t.Errorf("the tag reaches the remote: %s", remoteTags)
	}
	if local, head := gitIn(t, root, "rev-parse", "HEAD"), gitIn(t, remote, "rev-parse", "main"); local != head {
		t.Errorf("main reaches the remote too: %s vs %s", local, head)
	}

	again, _, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(again, "nothing pending") {
		t.Errorf("nothing left after publishing: %d %s", code, again)
	}
}

// A publish that fails after tagging but before (or during) the push is
// resumable: running it again does not recreate the tag, and finishes the
// push (S-0087).
func TestReleasePendingResumesAfterAFailedPush(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	// break the remote so the push half of publish fails
	broken := gitIn(t, root, "remote", "get-url", "origin")
	gitIn(t, root, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code == 0 || !strings.Contains(errOut, "the push failed") {
		t.Fatalf("the push should fail and say so: %d %s", code, errOut)
	}
	tagMsg := gitIn(t, root, "tag", "-l", "-n1", "cli/v1.0.1")
	if !strings.Contains(tagMsg, "cli/v1.0.1") {
		t.Fatalf("the tag was still created locally: %q", tagMsg)
	}

	// a second run must not try to recreate the tag or recompute the plan
	dry, _, _ := runIn(t, root, "release", "--pending", "--dry-run")
	if !strings.Contains(dry, "nothing pending") {
		t.Errorf("nothing new is pending; what remains is a push: %s", dry)
	}

	gitIn(t, root, "remote", "set-url", "origin", broken)
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "pushed to origin") {
		t.Fatalf("resumed publish should finish the push: %d %s %s", code, out, errOut)
	}
	if remoteTags := gitIn(t, remote, "tag", "--list", "cli/*"); !strings.Contains(remoteTags, "cli/v1.0.1") {
		t.Errorf("the tag made on the failed run reaches the remote now: %s", remoteTags)
	}
	// and it was not recreated: only one tag object for it
	count := strings.TrimSpace(gitIn(t, root, "tag", "--list", "cli/v1.0.1"))
	if strings.Count(count, "cli/v1.0.1") != 1 {
		t.Errorf("the tag exists exactly once: %q", count)
	}
}

// Nothing accepted at all: publishing says so and changes nothing.
func TestReleasePendingNothingPending(t *testing.T) {
	root, _ := researchProject(t, "remediation", true)
	out, _, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "nothing pending") {
		t.Errorf("nothing accepted, nothing pending: %d %s", code, out)
	}
}

package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/preview"
	"github.com/bytepunx/system-flow/flai/internal/release"
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
// push (S-0087). The push fails here because the remote moved; one that
// cannot be reached refuses before tagging (S-0174).
func TestReleasePendingResumesAfterAFailedPush(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	// the remote refuses the push half of publish; a remote that moved is
	// refused before anything is tagged (S-0195), so a hook stands in for a
	// push that fails after tagging
	hook := filepath.Join(remote, "hooks", "pre-receive")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}

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

	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
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

// S-0174: a clone whose tags lag the remote's, as after publishing from
// another clone, offers no plan built on its stale tag and refuses to
// publish, changing nothing, until the tags are fetched.
func TestReleasePendingFromACloneMissingTheRemotesTags(t *testing.T) {
	root, remote := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	gitIn(t, remote, "tag", "cli/v1.4.0", "main") // published from another clone
	head := gitIn(t, root, "rev-parse", "HEAD")

	dry, _, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(dry, "nothing can publish") || !strings.Contains(dry, "cli/v1.4.0 (here cli/v1.0.0)") || !strings.Contains(dry, "git fetch --tags origin") || strings.Contains(dry, "S-0001") {
		t.Fatalf("dry run from a lagging clone: %d %s", code, dry)
	}
	js, _, _ := runIn(t, root, "release", "--pending", "--dry-run", "--json")
	var pub preview.Publishing
	if err := json.Unmarshal([]byte(js), &pub); err != nil || len(pub.Plans) != 0 || !pub.Remote.Lagging() ||
		pub.Remote.Behind[0] != (release.Lag{Component: "cli", Local: "cli/v1.0.0", Remote: "cli/v1.4.0"}) {
		t.Errorf("dry run --json from a lagging clone: %v %s", err, js)
	}

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitPushDiverged || !strings.Contains(errOut, "refusing to publish") || !strings.Contains(errOut, "git fetch --tags origin") {
		t.Fatalf("publish from a lagging clone: %d %s", code, errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.0.1") {
		t.Errorf("a refused publish tags nothing: %s", tags)
	}
	if now := gitIn(t, root, "rev-parse", "HEAD"); now != head {
		t.Errorf("a refused publish commits nothing: %s, was %s", now, head)
	}
	if status := gitIn(t, root, "status", "--porcelain"); strings.TrimSpace(status) != "" {
		t.Errorf("a refused publish leaves the checkout as it was: %s", status)
	}

	gitIn(t, root, "fetch", "-q", "--tags", "origin")
	dry, _, _ = runIn(t, root, "release", "--pending", "--dry-run")
	if !strings.Contains(dry, "1.4.0 -> 1.4.1") || !strings.Contains(dry, "S-0001") {
		t.Errorf("with the tags fetched, the plan builds on the remote's: %s", dry)
	}
	out, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(out, "cli 1.4.1") {
		t.Errorf("publish once in step: %d %s %s", code, out, errOut)
	}
}

// S-0174: a remote that cannot be asked still shows the plan, with a
// warning, and refuses the publish before anything changes: the push would
// need that remote anyway.
func TestReleasePendingWithAnUnreachableRemote(t *testing.T) {
	root, _ := pendingProject(t)
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("accept: %s", errOut)
	}
	gitIn(t, root, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing.git"))
	head := gitIn(t, root, "rev-parse", "HEAD")

	dry, _, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(dry, "warning: could not ask origin") || !strings.Contains(dry, "S-0001") || !strings.Contains(dry, "1.0.0 -> 1.0.1") {
		t.Fatalf("dry run with the remote unreachable shows the plan and warns: %d %s", code, dry)
	}

	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != exitPushDiverged || !strings.Contains(errOut, "could not ask origin") || !strings.Contains(errOut, "refusing to publish") {
		t.Fatalf("publish with the remote unreachable: %d %s", code, errOut)
	}
	if tags := gitIn(t, root, "tag", "--list", "cli/*"); strings.Contains(tags, "cli/v1.0.1") {
		t.Errorf("a refused publish tags nothing: %s", tags)
	}
	if now := gitIn(t, root, "rev-parse", "HEAD"); now != head {
		t.Errorf("a refused publish commits nothing: %s, was %s", now, head)
	}
}

// I-0024: a story touching two components with no tag saying which it
// delivers to is named with the reason in the dry run and warned of when
// publishing, instead of vanishing from the batch.
func TestReleasePendingNamesWhatItCannotPlan(t *testing.T) {
	root, _ := pendingProject(t)
	manifest := filepath.Join(root, "system-flow.yaml")
	s, _ := os.ReadFile(manifest)
	_ = os.WriteFile(manifest, append(s, []byte("  - name: web\n    path: web\n    kind: sveltekit\n")...), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "web"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "web", ".gitkeep"), nil, 0o644)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "chore: a second component")
	if _, errOut, code := runIn(t, root, "stream", "sync", "S-0002"); code != 0 {
		t.Fatalf("sync: %s", errOut)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0002")
	_ = os.MkdirAll(filepath.Join(wt, "web"), 0o755)
	_ = os.WriteFile(filepath.Join(wt, "web", "both.js"), []byte("// both\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "fix: [S-0002] the other side")
	for _, id := range []string{"S-0001", "S-0002"} {
		if _, errOut, code := runIn(t, root, "accept", id); code != 0 {
			t.Fatalf("accept %s: %s", id, errOut)
		}
	}

	dry, _, code := runIn(t, root, "release", "--pending", "--dry-run")
	if code != 0 || !strings.Contains(dry, "left out: S-0002 Another fix: S-0002 touches cli, web but no tag says which it delivers to") || !strings.Contains(dry, "S-0001") {
		t.Fatalf("dry run names what it cannot plan and plans the rest: %d %s", code, dry)
	}
	_, errOut, code := runIn(t, root, "release", "--pending")
	if code != 0 || !strings.Contains(errOut, "S-0002 is left out of the release") {
		t.Errorf("publishing warns of what it leaves out: %d %s", code, errOut)
	}
}

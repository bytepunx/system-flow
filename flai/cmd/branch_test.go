package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gitIn runs git in dir for tests, failing the test on error.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestStoryBranchLifecycle(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("line one\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	// flai's own git calls (rebase, merge) need an identity in the repository; CI has no global one.
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")

	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Branchy", "--epic", "E-0001", "--touches", "docs/guide.md"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "task", "new", "Do it", "--story", "S-0001"); code != 0 {
		t.Fatal(errOut)
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-branchy.md")
	s, _ := os.ReadFile(storyFile)
	body := strings.Replace(string(s), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)
	_ = os.WriteFile(storyFile, []byte(body), 0o644)
	for _, st := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", "S-0001", st); code != 0 {
			t.Fatal(errOut)
		}
	}

	// open: narrative plus branch and worktree
	out, errOut, code := runIn(t, root, "stream", "open", "S-0001")
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	if code != 0 || !strings.Contains(out, "branch story/S-0001 checked out at .flai-cache/worktrees/S-0001") {
		t.Fatalf("open: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(wt, "system-flow.yaml")); err != nil {
		t.Fatal("worktree not checked out")
	}
	if gitIn(t, wt, "rev-parse", "--abbrev-ref", "HEAD") != "story/S-0001" {
		t.Fatal("worktree is not on the story branch")
	}

	// wip written from inside the worktree lands in the main checkout
	if _, errOut, code := runIn(t, wt, "task", "new", "From the worktree", "--story", "S-0001"); code != 0 {
		t.Fatal(errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/kanban/tasks/T-0002-from-the-worktree.md")); err != nil {
		t.Error("task created from the worktree should be in the main checkout's wip")
	}
	if _, err := os.Stat(filepath.Join(wt, "wip/kanban/tasks/T-0002-from-the-worktree.md")); err == nil {
		t.Error("task must not be written into the worktree's wip")
	}
	if _, errOut, code := runIn(t, wt, "touches", "T-0002", "docs/guide.md"); code != 0 {
		t.Fatal(errOut)
	}
	if out, _, _ := runIn(t, root, "board", "--all"); !strings.Contains(out, "touches docs/guide.md") {
		t.Errorf("board should show touches: %s", out)
	}

	// work on the branch, main moves on, sync rebases cleanly
	_ = os.WriteFile(filepath.Join(wt, "docs", "new.md"), []byte("new\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] new doc")
	_ = os.WriteFile(filepath.Join(root, "README.md"), []byte("readme\n"), 0o644)
	gitIn(t, root, "add", "README.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: readme")
	out, errOut, code = runIn(t, wt, "stream", "sync", "S-0001")
	if code != 0 || !strings.Contains(out, "story/S-0001 is rebased onto main") {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(wt, "README.md")); err != nil {
		t.Error("rebase should bring main's readme into the worktree")
	}

	// a conflict stops the rebase inside the worktree and names the file
	_ = os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("main's line\n"), 0o644)
	gitIn(t, root, "commit", "-q", "-am", "docs: guide on main")
	_ = os.WriteFile(filepath.Join(wt, "docs", "guide.md"), []byte("branch's line\n"), 0o644)
	gitIn(t, wt, "commit", "-q", "-am", "docs: [S-0001] guide on branch")
	_, errOut, code = runIn(t, wt, "stream", "sync", "S-0001")
	if code == 0 || !strings.Contains(errOut, "conflicts in docs/guide.md") || !strings.Contains(errOut, "rebase --continue") {
		t.Fatalf("conflict: %d %s", code, errOut)
	}
	_ = os.WriteFile(filepath.Join(wt, "docs", "guide.md"), []byte("resolved\n"), 0o644)
	gitIn(t, wt, "add", "docs/guide.md")
	c := exec.Command("git", "-c", "core.editor=true", "rebase", "--continue")
	c.Dir = wt
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("rebase --continue: %v %s", err, out)
	}

	// accept: rebase, fast-forward into main, remove worktree and branch
	for _, id := range []string{"T-0001", "T-0002"} {
		for _, st := range []string{"ready", "in-progress", "done"} {
			if _, errOut, code := runIn(t, root, "move", id, st); code != 0 {
				t.Fatal(errOut)
			}
		}
	}
	s, _ = os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "- [ ] done", "- [x] done", 1)), 0o644)
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	// a dry run previews the merge and changes nothing
	out, errOut, code = runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 || !strings.Contains(out, `"branch": "story/S-0001"`) || !strings.Contains(out, `"dry_run": true`) {
		t.Fatalf("dry run: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Fatal("dry run must not remove the worktree")
	}
	// moving the story to done IS acceptance (S-0046): same flow, same flags
	out, errOut, code = runIn(t, root, "move", "S-0001", "done", "--by", "tester")
	if code != 0 || !strings.Contains(out, "accepted S-0001") || !strings.Contains(out, "story/S-0001 merged and removed") {
		t.Fatalf("move to done: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(wt); err == nil {
		t.Error("worktree should be removed")
	}
	if br := gitIn(t, root, "branch", "--list", "story/S-0001"); br != "" {
		t.Errorf("branch should be deleted: %q", br)
	}
	log := gitIn(t, root, "log", "--oneline", "-5")
	if !strings.Contains(log, "[S-0001] new doc") || !strings.Contains(log, "accept and archive") {
		t.Errorf("main should contain the branch commits and the acceptance: %s", log)
	}
	if guide, _ := os.ReadFile(filepath.Join(root, "docs", "guide.md")); string(guide) != "resolved\n" {
		t.Errorf("guide after accept: %q", guide)
	}
	if _, err := os.Stat(filepath.Join(root, "wip/archive/kanban/stories/S-0001-branchy.md")); err != nil {
		t.Error("story should be archived")
	}
}

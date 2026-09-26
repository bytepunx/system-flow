package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S-0140: a story went to review with work left uncommitted in its worktree,
// and the operator's acceptance failed once its agent had gone. The move to
// review refuses it, naming the paths; the acceptance preview lists them and
// blocks on them before anything is merged.
func TestAStoryGoesToReviewOnlyWithItsWorktreeCommitted(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Left over", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-left-over.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wt, "docs", "done.md"), []byte("x\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] done")
	_ = os.WriteFile(filepath.Join(wt, "docs", "forgotten.md"), []byte("x\n"), 0o644)

	_, errOut, code := runIn(t, root, "move", "S-0001", "review")
	if code == 0 || !strings.Contains(errOut, ".flai-cache/worktrees/S-0001 has uncommitted changes (docs/forgotten.md)") || !strings.Contains(errOut, "story/S-0001") {
		t.Fatalf("the move must refuse the uncommitted file and say where to commit it: %d %s", code, errOut)
	}
	if out, _, _ := runIn(t, root, "show", "S-0001", "--json"); !strings.Contains(out, `"status": "in-progress"`) && !strings.Contains(out, `"status":"in-progress"`) {
		t.Errorf("a refused move changes nothing:\n%s", out)
	}

	// From review, as a story moved there by an older flai or by hand would
	// be: the preview lists the paths and blocks before anything is merged.
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] forgotten")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatalf("a committed worktree moves: %s", errOut)
	}
	_ = os.WriteFile(filepath.Join(wt, "docs", "after.md"), []byte("x\n"), 0o644)
	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("a dry run must not refuse: %d %s", code, errOut)
	}
	var res struct {
		WorktreeUncommitted []string `json:"worktree_uncommitted"`
		Blockers            []string `json:"blockers"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("preview json: %v\n%s", err, out)
	}
	if len(res.WorktreeUncommitted) != 1 || res.WorktreeUncommitted[0] != "docs/after.md" {
		t.Errorf("worktree_uncommitted = %v", res.WorktreeUncommitted)
	}
	if len(res.Blockers) != 1 || !strings.Contains(res.Blockers[0], "docs/after.md") {
		t.Errorf("blockers = %q", res.Blockers)
	}
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code == 0 || !strings.Contains(errOut, "cannot be accepted yet") {
		t.Errorf("the real run refuses before merging: %d %s", code, errOut)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Errorf("the worktree must be left in place: %v", err)
	}
}

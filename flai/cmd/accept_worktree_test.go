package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// I-0017: the dashboard saw the repository at another path than the host, so
// the worktree's link pointed nowhere and acceptance failed with git's raw
// error. It must say what to do, in the preview and in the real run, and
// change nothing.
func TestAcceptSaysWhatToDoWhenTheWorktreeCannotBeOpened(t *testing.T) {
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
	for _, args := range [][]string{
		{"epic", "new", "Epic"},
		{"story", "new", "Elsewhere", "--epic", "E-0001"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-elsewhere.md")
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
	_ = os.WriteFile(filepath.Join(wt, "docs", "x.md"), []byte("x\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] x")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	// What the container saw: the link names a path that is not there.
	_ = os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: /host/only/path/.git/worktrees/S-0001\n"), 0o644)

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 {
		t.Fatalf("the preview must succeed and carry the blocker: %d %s", code, errOut)
	}
	var res struct {
		Blockers []string `json:"blockers"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("preview json: %v\n%s", err, out)
	}
	if len(res.Blockers) != 1 || !strings.Contains(res.Blockers[0], "flai accept S-0001") || !strings.Contains(res.Blockers[0], ".flai-cache/worktrees/S-0001") {
		t.Errorf("blockers: %q", res.Blockers)
	}
	// git's own text belongs in the log, not in what the card shows
	if strings.Contains(out, "exit status 128") || strings.Contains(out, "not a git repository") {
		t.Errorf("git's raw error must not be the message:\n%s", out)
	}

	_, errOut, code = runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "cannot be accepted yet") || !strings.Contains(errOut, "flai accept S-0001") {
		t.Errorf("the real run must refuse with the same advice: %d %s", code, errOut)
	}
	after, _ := os.ReadFile(storyFile)
	if !strings.Contains(string(after), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}
}

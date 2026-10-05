package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// I-0066: a rebase conflict during a story's sync was git-added unresolved,
// and acceptance fast-forwarded main onto the markers. Acceptance must refuse
// such a branch, naming each marker's path and line and merging nothing, and
// accept it once the conflict is resolved and committed (S-0253).
func TestAcceptRefusesABranchCarryingConflictMarkers(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Conflicted", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-conflicted.md")
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
	doc := filepath.Join(wt, "design", "x.md")
	// The markers are built so that no line of this file is one. The file
	// opens with a blank line, which a reader that trims git's output would
	// lose and so misnumber every marker.
	conflicted := "\n# X\n\n" +
		strings.Repeat("<", 7) + " HEAD\nours\n" +
		strings.Repeat("=", 7) + "\ntheirs\n" +
		strings.Repeat(">", 7) + " f0f5443 (docs: [S-0201] keep both)\nend\n"
	_ = os.WriteFile(doc, []byte(conflicted), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] x")
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	mainBefore := gitIn(t, root, "rev-parse", "main")

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 {
		t.Fatal("acceptance must refuse a branch carrying conflict markers")
	}
	for _, want := range []string{"story/S-0001 carries merge conflict markers at design/x.md:4, design/x.md:8;", ".flai-cache/worktrees/S-0001", "accept again"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the refusal must say %q:\n%s", want, errOut)
		}
	}
	if strings.Contains(errOut, "design/x.md:6") {
		t.Errorf("the divider is not a marker:\n%s", errOut)
	}
	if got := gitIn(t, root, "rev-parse", "main"); got != mainBefore {
		t.Errorf("main moved from %s to %s", mainBefore, got)
	}
	if after, _ := os.ReadFile(storyFile); !strings.Contains(string(after), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") == "" {
		t.Error("a refused acceptance must keep the story branch")
	}
	if _, err := os.Stat(wt); err != nil {
		t.Errorf("a refused acceptance must keep the worktree: %v", err)
	}

	// resolved and committed, the branch is accepted
	_ = os.WriteFile(doc, []byte("\n# X\n\nours and theirs\nend\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "docs: [S-0001] resolve x")
	if _, errOut, code := runIn(t, root, "accept", "S-0001"); code != 0 {
		t.Fatalf("a resolved branch must be accepted: %d %s", code, errOut)
	}
	if got := gitIn(t, root, "show", "main:design/x.md"); !strings.Contains(got, "ours and theirs") {
		t.Errorf("main must carry the resolution:\n%s", got)
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") != "" {
		t.Error("an accepted story's branch is removed")
	}
}

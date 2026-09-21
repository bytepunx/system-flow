package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S-0051: uncommitted paths outside wip are reported by the dry run so the
// choice is made before confirming; the real run still refuses them unless
// --yes includes them.
func TestAcceptReportsUncommittedPathsInTheDryRun(t *testing.T) {
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
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Stray", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-stray.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{
		{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001", "--no-branch"},
		{"task", "new", "Do it", "--story", "S-0001"},
		{"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "S-0001", "review"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	preview := func(args ...string) []string {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"accept", "S-0001", "--dry-run", "--json"}, args...)...)
		if code != 0 {
			t.Fatalf("a dry run must not refuse (%v): %d %s", args, code, errOut)
		}
		var res struct {
			Uncommitted []string `json:"uncommitted"`
		}
		if err := json.Unmarshal([]byte(out), &res); err != nil {
			t.Fatalf("preview json: %v\n%s", err, out)
		}
		return res.Uncommitted
	}
	if got := preview(); len(got) != 0 {
		t.Errorf("clean tree, wip changes only: uncommitted = %v", got)
	}

	_ = os.WriteFile(filepath.Join(root, "docs", "stray.md"), []byte("x\n"), 0o644)
	for _, args := range [][]string{nil, {"--yes"}} {
		if got := preview(args...); len(got) != 1 || got[0] != "docs/stray.md" {
			t.Errorf("dry run %v should report the stray file, got %v", args, got)
		}
	}
	out, _, _ := runIn(t, root, "accept", "S-0001", "--dry-run")
	if !strings.Contains(out, "uncommitted outside wip: docs/stray.md") || !strings.Contains(out, "--yes") {
		t.Errorf("text preview should name the file and the way to include it:\n%s", out)
	}

	_, errOut, code := runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "uncommitted changes outside wip (docs/stray.md)") {
		t.Errorf("the real run must still refuse without --yes: %d %s", code, errOut)
	}
	if after, _ := os.ReadFile(storyFile); !strings.Contains(string(after), "status: review") {
		t.Error("a refused acceptance must leave the story in review")
	}

	if _, errOut, code := runIn(t, root, "move", "S-0001", "done", "--yes"); code != 0 {
		t.Fatalf("move to done with --yes should accept and include the file: %s", errOut)
	}
	if files := gitIn(t, root, "show", "--name-only", "--format=", "HEAD"); !strings.Contains(files, "docs/stray.md") {
		t.Errorf("the acceptance commit should hold the stray file:\n%s", files)
	}
	if st := gitIn(t, root, "status", "--porcelain"); st != "" {
		t.Errorf("tree should be clean after acceptance: %s", st)
	}
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// syncProject is a git repository with an epic, for stories whose branches
// sync checks against each other (S-0131). Integration: it runs real git.
func syncProject(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
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
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	return root
}

// openSyncStory creates story n in progress with touches and opens its
// worktree, returning the worktree's path.
func openSyncStory(t *testing.T, root string, n int, touches string) string {
	t.Helper()
	id := fmt.Sprintf("S-%04d", n)
	title := fmt.Sprintf("Story %d", n)
	if _, errOut, code := runIn(t, root, "story", "new", title, "--epic", "E-0001", "--touches", touches); code != 0 {
		t.Fatal(errOut)
	}
	file := filepath.Join(root, "wip/kanban/stories", fmt.Sprintf("%s-story-%d.md", id, n))
	s, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] done\n", 1)), 0o644)
	for _, st := range []string{"ready", "in-progress"} {
		if _, errOut, code := runIn(t, root, "move", id, st); code != 0 {
			t.Fatal(errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "stream", "open", id); code != 0 {
		t.Fatal(errOut)
	}
	return filepath.Join(root, ".flai-cache", "worktrees", id)
}

// commitIn writes content to rel in the worktree and commits it.
func commitIn(t *testing.T, wt, rel, content string) {
	t.Helper()
	p := filepath.Join(wt, rel)
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "change "+rel)
}

func TestSyncTrialMergesOtherOpenBranches(t *testing.T) {
	root := syncProject(t)
	one := openSyncStory(t, root, 1, "docs")
	two := openSyncStory(t, root, 2, "docs/guide.md")
	three := openSyncStory(t, root, 3, "design")
	commitIn(t, one, "docs/guide.md", "one's line\n")
	commitIn(t, two, "docs/guide.md", "two's line\n")
	commitIn(t, three, "design/note.md", "three\n")

	out, errOut, code := runIn(t, one, "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{
		"story/S-0001 is rebased onto main",
		"story/S-0001 conflicts with story/S-0002 (in progress) in docs/guide.md",
		"story/S-0001 merges cleanly with story/S-0003 (in progress)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sync output lacks %q:\n%s", want, out)
		}
	}
	// the trial merge writes nothing to any worktree
	for _, wt := range []string{one, two, three} {
		if st := gitIn(t, wt, "status", "--porcelain"); st != "" {
			t.Errorf("%s is not clean after sync: %s", wt, st)
		}
	}
	if got := gitIn(t, two, "show", "HEAD:docs/guide.md"); got != "two's line" {
		t.Errorf("story/S-0002 changed: %q", got)
	}

	out, errOut, code = runIn(t, one, "--json", "stream", "sync", "S-0001")
	if code != 0 {
		t.Fatalf("sync --json: %d %s %s", code, out, errOut)
	}
	var res struct {
		OK       bool          `json:"ok"`
		Branches []branchCheck `json:"branches"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if !res.OK || len(res.Branches) != 2 {
		t.Fatalf("branches: %+v", res)
	}
	if b := res.Branches[0]; b.Story != "S-0002" || b.Clean || strings.Join(b.Conflicts, ",") != "docs/guide.md" || b.Branch != "story/S-0002" || b.Status != "in-progress" {
		t.Errorf("conflicting pair: %+v", b)
	}
	if b := res.Branches[1]; b.Story != "S-0003" || !b.Clean || len(b.Conflicts) != 0 {
		t.Errorf("clean pair: %+v", b)
	}
}

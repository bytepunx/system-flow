package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// reviewProject is a git project with story S-0001 in progress on its branch.
func reviewProject(t *testing.T) (root, wt string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	root = tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "docs"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
		_ = os.WriteFile(filepath.Join(root, d, ".gitkeep"), nil, 0o644)
	}
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "keep.md"), []byte("one\ntwo\nthree\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "gone.md"), []byte("bye\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "old-name.md"), []byte(strings.Repeat("same line\n", 20)), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "t@t")
	gitIn(t, root, "config", "user.name", "t")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Reviewed", "--epic", "E-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	file := filepath.Join(root, "wip/kanban/stories/S-0001-reviewed.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, args := range [][]string{{"move", "S-0001", "ready"}, {"move", "S-0001", "in-progress"}, {"stream", "open", "S-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	return root, filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
}

func TestStreamDiff(t *testing.T) {
	root, wt := reviewProject(t)
	_ = os.WriteFile(filepath.Join(wt, "docs", "keep.md"), []byte("one\n2\nthree\nfour\n"), 0o644)
	_ = os.WriteFile(filepath.Join(wt, "docs", "new.md"), []byte("hello\n"), 0o644)
	_ = os.Remove(filepath.Join(wt, "docs", "gone.md"))
	_ = os.Rename(filepath.Join(wt, "docs", "old-name.md"), filepath.Join(wt, "docs", "new-name.md"))
	_ = os.WriteFile(filepath.Join(wt, "docs", "blob.bin"), []byte{0, 1, 2, 0, 255, 0}, 0o644)
	_ = os.WriteFile(filepath.Join(wt, "docs", "huge.md"), []byte(strings.Repeat("a long line of text to fill the patch\n", 4000)), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] changes")
	// main moves on: the diff is against the merge base, not main's tip
	_ = os.WriteFile(filepath.Join(root, "docs", "elsewhere.md"), []byte("main only\n"), 0o644)
	gitIn(t, root, "add", "docs/elsewhere.md")
	gitIn(t, root, "commit", "-q", "-m", "docs: elsewhere")

	out, errOut, code := runIn(t, root, "stream", "diff", "S-0001", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var d streamDiff
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if d.Branch != "story/S-0001" || d.Commits != 1 || d.Base == "" || !d.Truncated {
		t.Errorf("header: %+v", d)
	}
	got := map[string]diffFile{}
	for _, f := range d.Files {
		got[f.Path] = f
	}
	if _, ok := got["docs/elsewhere.md"]; ok {
		t.Error("what main did since the branch point is not the story's change")
	}
	if f := got["docs/keep.md"]; f.Status != "modified" || f.Additions != 2 || f.Deletions != 1 || !strings.HasPrefix(f.Patch, "@@") || !strings.Contains(f.Patch, "-two") || !strings.Contains(f.Patch, "+four") {
		t.Errorf("modified: %+v", f)
	}
	if f := got["docs/new.md"]; f.Status != "added" || f.Additions != 1 || !strings.Contains(f.Patch, "+hello") {
		t.Errorf("added: %+v", f)
	}
	if f := got["docs/gone.md"]; f.Status != "deleted" || f.Deletions != 1 {
		t.Errorf("deleted: %+v", f)
	}
	if f := got["docs/new-name.md"]; f.Status != "renamed" || f.OldPath != "docs/old-name.md" {
		t.Errorf("renamed: %+v", f)
	}
	if f := got["docs/blob.bin"]; !f.Binary || f.Patch != "" {
		t.Errorf("binary: %+v", f)
	}
	if f := got["docs/huge.md"]; !f.Truncated || len(f.Patch) > diffFileLimit || f.Additions != 4000 {
		t.Errorf("truncated: status %s, %d bytes, +%d, truncated %v", f.Status, len(f.Patch), f.Additions, f.Truncated)
	}
	text, _, _ := runIn(t, root, "stream", "diff", "S-0001")
	if !strings.Contains(text, "story/S-0001 against") || !strings.Contains(text, "docs/old-name.md -> docs/new-name.md") {
		t.Errorf("text output:\n%s", text)
	}
}

func TestStreamDiffWithoutABranch(t *testing.T) {
	root, _ := reviewProject(t)
	if _, errOut, code := runIn(t, root, "story", "new", "No branch", "--epic", "E-0001"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "stream", "diff", "S-0002"); code == 0 || !strings.Contains(errOut, "has no branch story/S-0002") {
		t.Errorf("no branch: %d %s", code, errOut)
	}
}

func TestAcceptLogsEachStep(t *testing.T) {
	root, wt := reviewProject(t)
	t.Setenv("LOG_FORMAT", "json")
	_ = os.WriteFile(filepath.Join(wt, "docs", "new.md"), []byte("hello\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] new")
	for _, args := range [][]string{
		{"task", "new", "Do it", "--story", "S-0001"}, {"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "S-0001", "review"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	_, errOut, code := runIn(t, root, "accept", "S-0001", "--by", "dana", "--json")
	if code != 0 {
		t.Fatal(errOut)
	}
	var steps []string
	for _, line := range strings.Split(errOut, "\n") {
		var e struct{ Msg, Step, Item string }
		if json.Unmarshal([]byte(line), &e) == nil && e.Msg == "acceptance step" && e.Item == "S-0001" {
			steps = append(steps, e.Step)
		}
	}
	if got := strings.Join(steps, " "); got != "merged done archived committed" {
		t.Errorf("steps in order, no tags without a release and no push without a remote: %q", got)
	}
}

// Found from the review page (S-0041): a story whose criteria are not all
// ticked was refused only after its branch had been merged and its worktree
// removed. The rules for done are checked first now, and the preview says so.
func TestAcceptChecksTheRulesBeforeItMergesAnything(t *testing.T) {
	root, wt := reviewProject(t)
	file := filepath.Join(root, "wip/kanban/stories/S-0001-reviewed.md")
	s, _ := os.ReadFile(file)
	_ = os.WriteFile(file, []byte(strings.Replace(string(s), "- [x] done\n", "- [ ] done\n", 1)), 0o644)
	_ = os.WriteFile(filepath.Join(wt, "docs", "new.md"), []byte("hello\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] new")
	for _, args := range [][]string{
		{"task", "new", "Do it", "--story", "S-0001"}, {"move", "T-0001", "ready"}, {"move", "T-0001", "in-progress"}, {"move", "T-0001", "done"},
		{"move", "S-0001", "review"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	head := gitIn(t, root, "rev-parse", "main")

	out, errOut, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 || !strings.Contains(out, "S-0001 has unchecked acceptance criteria") {
		t.Errorf("the preview should carry the rule as a blocker: %d %s %s", code, out, errOut)
	}
	_, errOut, code = runIn(t, root, "accept", "S-0001")
	if code == 0 || !strings.Contains(errOut, "unchecked acceptance criteria") {
		t.Fatalf("acceptance must be refused: %d %s", code, errOut)
	}
	if gitIn(t, root, "rev-parse", "main") != head {
		t.Error("the branch was merged into main before the refusal")
	}
	if _, err := os.Stat(filepath.Join(wt, "docs", "new.md")); err != nil {
		t.Error("the worktree was removed before the refusal")
	}
	if gitIn(t, root, "branch", "--list", "story/S-0001") == "" {
		t.Error("the story branch was deleted before the refusal")
	}
	if after, _ := os.ReadFile(file); !strings.Contains(string(after), "status: review") {
		t.Error("the story must stay in review")
	}
}

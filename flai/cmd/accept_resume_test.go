package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// A story that is done but was never accepted (an older flai, a hand edit,
// a board that only changed the state) is flagged by check and completed by
// accept (S-0046).
func TestAcceptCompletesAnUnacceptedStory(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	root := tempProject(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "design/adrs", "design/system", "design/tech", "design/conventions", "docs"} {
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
		{"story", "new", "Stranded", "--epic", "E-0001"},
		{"task", "new", "Do it", "--story", "S-0001"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-stranded.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, mv := range [][]string{{"S-0001", "ready"}, {"S-0001", "in-progress"}} {
		if _, errOut, code := runIn(t, root, "move", mv[0], mv[1]); code != 0 {
			t.Fatal(errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "stream", "open", "S-0001"); code != 0 {
		t.Fatal(errOut)
	}
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0001")
	_ = os.WriteFile(filepath.Join(wt, "docs", "x.md"), []byte("x\n"), 0o644)
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feat: [S-0001] x")
	for _, st := range []string{"ready", "in-progress", "done"} {
		if _, errOut, code := runIn(t, root, "move", "T-0001", st); code != 0 {
			t.Fatal(errOut)
		}
	}
	// Strand it the way an older flai did: status done by hand, nothing else.
	s, _ = os.ReadFile(storyFile)
	stranded := strings.Replace(string(s), "status: in-progress", "status: done", 1)
	stranded = strings.Replace(stranded, "\ntags:", "\n  - to: review\n    at: 2026-09-15T21:00:00Z\n    by: t\n  - to: done\n    at: 2026-09-15T21:00:01Z\n    by: t\ntags:", 1)
	_ = os.WriteFile(storyFile, []byte(stranded), 0o644)

	out, _, _ := runIn(t, root, "check")
	if !strings.Contains(out, "story.unaccepted") || !strings.Contains(out, "flai accept S-0001") {
		t.Errorf("check should flag the stranded story:\n%s", out)
	}
	out, errOut, code := runIn(t, root, "accept", "S-0001")
	if code != 0 || !strings.Contains(out, "completed the acceptance of S-0001") || !strings.Contains(out, "story/S-0001 merged and removed") {
		t.Fatalf("resume: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "x.md")); err != nil {
		t.Error("the branch commit should be on main")
	}
	if _, err := os.Stat(filepath.Join(root, "wip/archive/kanban/stories/S-0001-stranded.md")); err != nil {
		t.Error("story should be archived")
	}
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "story.unaccepted") {
		t.Errorf("nothing should be flagged after completion:\n%s", out)
	}
	if _, _, code := runIn(t, root, "accept", "S-0001"); code == 0 {
		t.Error("an archived story cannot be accepted again")
	}
}

// Without a committer identity the acceptance commit would fail after the
// merge and the archive; the preflight refuses before anything changes.
func TestAcceptRefusesBeforeChangingAnythingWithoutIdentity(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	for _, k := range []string{"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "EMAIL"} {
		t.Setenv(k, "x") // registers the restore; an empty value would still override git config
		_ = os.Unsetenv(k)
	}
	root := tempProject(t)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "add", "-A")
	c := exec.Command("git", "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "init")
	c.Dir = root
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	for _, args := range [][]string{{"epic", "new", "Epic"}, {"story", "new", "Anon", "--epic", "E-0001"}, {"task", "new", "Do", "--story", "S-0001"}} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatal(errOut)
		}
	}
	storyFile := filepath.Join(root, "wip/kanban/stories/S-0001-anon.md")
	s, _ := os.ReadFile(storyFile)
	_ = os.WriteFile(storyFile, []byte(strings.Replace(string(s), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] done\n", 1)), 0o644)
	for _, mv := range [][]string{{"S-0001", "ready"}, {"S-0001", "in-progress"}, {"T-0001", "ready"}, {"T-0001", "in-progress"}, {"T-0001", "done"}} {
		if _, errOut, code := runIn(t, root, "move", mv[0], mv[1]); code != 0 {
			t.Fatal(errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "stream", "open", "S-0001", "--no-branch"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "move", "S-0001", "review"); code != 0 {
		t.Fatal(errOut)
	}
	out, _, code := runIn(t, root, "accept", "S-0001", "--dry-run", "--json")
	if code != 0 || !strings.Contains(out, "committer identity") {
		t.Errorf("the dry run should list the blocker: %d %s", code, out)
	}
	_, errOut, code := runIn(t, root, "move", "S-0001", "done")
	if code == 0 || !strings.Contains(errOut, "cannot be accepted yet") {
		t.Fatalf("should refuse: %d %s", code, errOut)
	}
	after, _ := os.ReadFile(storyFile)
	if !strings.Contains(string(after), "status: review") {
		t.Errorf("the story must still be in review:\n%s", after)
	}
}

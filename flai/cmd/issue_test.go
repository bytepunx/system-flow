package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIssueRecordingStoryOrder(t *testing.T) {
	env := func(vars map[string]string) func(string) string {
		return func(k string) string { return vars[k] }
	}
	branch := func(b string) func() string { return func() string { return b } }
	for _, c := range []struct {
		name, flag string
		env        map[string]string
		branch     string
		want       string
	}{
		{"flag first", "S-0001", map[string]string{"FLAI_STORY": "S-0002", "FLAI_AGENT": "agent-S-0003"}, "story/S-0004", "S-0001"},
		{"then FLAI_STORY", "", map[string]string{"FLAI_STORY": "S-0002", "FLAI_AGENT": "agent-S-0003"}, "story/S-0004", "S-0002"},
		{"then FLAI_AGENT", "", map[string]string{"FLAI_AGENT": "agent-S-0003"}, "story/S-0004", "S-0003"},
		{"an agent not of the form agent-S-nnnn names none", "", map[string]string{"FLAI_AGENT": "tester"}, "story/S-0004", "S-0004"},
		{"then the branch", "", nil, "story/S-0004", "S-0004"},
		{"outside any story", "", map[string]string{"FLAI_AGENT": "tester"}, "main", ""},
		{"no git", "", nil, "", ""},
	} {
		if got := recordingStory(c.flag, env(c.env), branch(c.branch)); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

// issueProject is a project for the issue commands, with no story named by
// the environment.
func issueProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "")
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "design", "conventions"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "conventions", "README.md"), []byte("# c\n"), 0o644)
	return root
}

func TestIssueListStory(t *testing.T) {
	root := issueProject(t)
	for _, args := range [][]string{
		{"issue", "new", "Lint on the host is v1", "--class", "efficiency", "--story", "S-0007"},
		{"issue", "new", "Fixture was ignored", "--class", "defect"},
		{"issue", "bump", "I-0002", "--story", "S-0008"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	t.Setenv("FLAI_STORY", "S-0007")
	if _, errOut, code := runIn(t, root, "issue", "bump", "I-0002"); code != 0 {
		t.Fatalf("bump under FLAI_STORY: %d %s", code, errOut)
	}
	t.Setenv("FLAI_STORY", "")

	out, _, _ := runIn(t, root, "issue", "list", "--story", "S-0008")
	if !strings.Contains(out, "I-0002") || strings.Contains(out, "I-0001") {
		t.Errorf("list --story S-0008 should give only I-0002:\n%s", out)
	}
	out, _, _ = runIn(t, root, "issue", "list", "--story", "S-7")
	if !strings.Contains(out, "I-0001") || !strings.Contains(out, "I-0002") {
		t.Errorf("list --story S-7 should give both, the ID in any padding:\n%s", out)
	}

	if out, errOut, code := runIn(t, root, "issue", "story", "I-0001"); code != 0 {
		t.Fatalf("issue story: %d %s %s", code, out, errOut)
	}
	out, _, code := runIn(t, root, "issue", "list", "--json")
	var list []map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &list) != nil || len(list) != 2 {
		t.Fatalf("list --json: %s", out)
	}
	stories := func(v any) string {
		var s []string
		for _, x := range v.([]any) {
			s = append(s, x.(string))
		}
		return strings.Join(s, ",")
	}
	if list[0]["id"] != "I-0001" || list[0]["count"].(float64) != 1 || stories(list[0]["stories"]) != "S-0007" || list[0]["story"] != "S-0001" {
		t.Errorf("I-0001 in list --json: %v", list[0])
	}
	if list[1]["id"] != "I-0002" || stories(list[1]["stories"]) != "S-0008,S-0007" || list[1]["story"] != "" {
		t.Errorf("I-0002 in list --json: %v", list[1])
	}
	out, _, _ = runIn(t, root, "issue", "list", "--story", "S-0008", "--json")
	if json.Unmarshal([]byte(out), &list) != nil || len(list) != 1 || list[0]["id"] != "I-0002" {
		t.Errorf("list --story --json: %s", out)
	}
	out, _, _ = runIn(t, root, "issue", "list", "--story", "S-0099", "--json")
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("list --story with none should be an empty array: %s", out)
	}
}

func TestIssueStory(t *testing.T) {
	root := issueProject(t)
	for _, args := range [][]string{
		{"issue", "new", "Fixture was ignored", "--class", "defect"},
		{"issue", "new", "Lint on the host is v1", "--class", "efficiency"},
		{"issue", "new", "Go not on PATH", "--class", "blocker"},
		{"epic", "new", "Upkeep"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	issueFiles := map[string][]byte{}
	for _, f := range []string{"I-0001-fixture-was-ignored.md", "I-0002-lint-on-the-host-is-v1.md", "summary.md"} {
		issueFiles[f] = readIssueFile(t, root, f)
	}
	out, errOut, code := runIn(t, root, "issue", "story", "I-0001")
	if code != 0 || out != "S-0001 Fixture was ignored\n  wip/kanban/stories/S-0001-fixture-was-ignored.md\n" {
		t.Fatalf("issue story I-0001: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "issue", "story", "I-0002", "--epic", "E-0001", "--json")
	var made map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &made) != nil {
		t.Fatalf("issue story I-0002 --json: %d %s %s", code, out, errOut)
	}
	want := map[string]any{"id": "S-0002", "title": "Lint on the host is v1", "nature": "improvement", "path": "wip/kanban/stories/S-0002-lint-on-the-host-is-v1.md", "issue": "I-0002", "committed": false}
	for k, v := range want {
		if made[k] != v {
			t.Errorf("issue story --json %s = %v, want %v", k, made[k], v)
		}
	}
	for file, nature := range map[string]string{"S-0001-fixture-was-ignored.md": "remediation", "S-0002-lint-on-the-host-is-v1.md": "improvement"} {
		data, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "stories", file))
		if !strings.Contains(string(data), "\nnature: "+nature+"\n") || !strings.Contains(string(data), "\nstatus: backlog\n") {
			t.Errorf("%s should be a %s story in backlog:\n%s", file, nature, data)
		}
	}
	epic, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "epics", "E-0001-upkeep.md"))
	if !strings.Contains(string(epic), "S-0002") {
		t.Errorf("the epic should link the story made under it:\n%s", epic)
	}
	for _, f := range []string{"I-0001-fixture-was-ignored.md", "I-0002-lint-on-the-host-is-v1.md", "summary.md"} {
		if string(issueFiles[f]) != string(readIssueFile(t, root, f)) {
			t.Errorf("issue story should not change design/issues/%s", f)
		}
	}

	storyFiles := func() int {
		m, _ := filepath.Glob(filepath.Join(root, "wip", "kanban", "stories", "*.md"))
		return len(m)
	}
	_, errOut, code = runIn(t, root, "issue", "story", "I-0001")
	if code == 0 || !strings.Contains(errOut, "already linked by open story S-0001") || storyFiles() != 2 {
		t.Errorf("an issue an open story links should be refused, naming it: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "issue", "close", "I-0003", "--reason", "fixed"); code != 0 {
		t.Fatalf("close: %s", errOut)
	}
	before, _ := os.ReadFile(filepath.Join(root, "design", "issues", "I-0003-go-not-on-path.md"))
	_, errOut, code = runIn(t, root, "issue", "story", "I-0003")
	after, _ := os.ReadFile(filepath.Join(root, "design", "issues", "I-0003-go-not-on-path.md"))
	if code == 0 || !strings.Contains(errOut, "I-0003 is closed") || storyFiles() != 2 || string(before) != string(after) {
		t.Errorf("a closed issue should be refused with nothing changed: %d %s", code, errOut)
	}
}

// The dashboard makes a story from an issue as it makes any item: as the
// manifest's owner, committed on its own with its trailer.
func TestIssueStoryAutocommit(t *testing.T) {
	root := bodyProject(t)
	t.Setenv("FLAI_STORY", "")
	if _, errOut, code := runIn(t, root, "issue", "new", "Fixture was ignored", "--class", "defect"); code != 0 {
		t.Fatalf("issue new: %s", errOut)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "the issue")
	out, errOut, code := runIn(t, root, "issue", "story", "I-0001", "--epic", "E-0001", "--owner", "olive", "--autocommit", "--trailer", "Created-with: flaiover", "--json")
	var made struct {
		ID, Path, Commit string
		Committed        bool
	}
	if code != 0 || json.Unmarshal([]byte(out), &made) != nil {
		t.Fatalf("issue story --autocommit: %d %s %s", code, out, errOut)
	}
	if made.ID != "S-0001" || !made.Committed || made.Commit == "" {
		t.Errorf("the story should be made and committed: %+v", made)
	}
	if data, _ := os.ReadFile(filepath.Join(root, made.Path)); !strings.Contains(string(data), "\nowner: olive\n") {
		t.Errorf("the story should be owned by --owner:\n%s", data)
	}
	show := gitIn(t, root, "show", "--stat", "--format=%s|%b", "HEAD")
	if !strings.HasPrefix(show, "chore: [S-0001] create story: Fixture was ignored|Created-with: flaiover") {
		t.Errorf("commit: %s", show)
	}
	if !strings.Contains(show, "S-0001-fixture-was-ignored.md") || !strings.Contains(show, "E-0001-workbench.md") || !strings.Contains(show, "2 files changed") {
		t.Errorf("the commit should hold the story and its epic only: %s", show)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("nothing should be left uncommitted: %q", st)
	}
}

func readIssueFile(t *testing.T, root, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "design", "issues", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// A story's issues are committed on its branch, so until it is accepted
// they are in its worktree and not in the main checkout: --story reads them
// there, and the story made from one goes to wip/ in the main checkout.
func TestIssueStoryWorktree(t *testing.T) {
	root := issueProject(t)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0005")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runIn(t, root, "issue", "new", "On main", "--class", "efficiency", "--story", "S-0004"); code != 0 {
		t.Fatalf("issue new in the root: %s", errOut)
	}
	// the worktree holds main's issues as of its last sync, and its own
	_ = os.MkdirAll(filepath.Join(wt, "design", "issues"), 0o755)
	for _, f := range []string{"system-flow.yaml", "design/issues/I-0001-on-main.md"} {
		data, _ := os.ReadFile(filepath.Join(root, f))
		_ = os.WriteFile(filepath.Join(wt, f), data, 0o644)
	}
	if _, errOut, code := runIn(t, wt, "issue", "new", "Only on the branch", "--class", "defect", "--story", "S-0005"); code != 0 {
		t.Fatalf("issue new in the worktree: %s", errOut)
	}

	out, _, _ := runIn(t, root, "issue", "list", "--story", "S-5")
	if !strings.Contains(out, "Only on the branch") || strings.Contains(out, "On main") {
		t.Errorf("list --story should read the story's worktree:\n%s", out)
	}
	out, _, _ = runIn(t, root, "issue", "list")
	if !strings.Contains(out, "On main") || strings.Contains(out, "Only on the branch") {
		t.Errorf("list without --story should read the root:\n%s", out)
	}
	out, _, _ = runIn(t, root, "issue", "list", "--story", "S-0006")
	if strings.TrimSpace(out) != "" {
		t.Errorf("a story with no worktree reads the root, where none names it:\n%s", out)
	}

	if _, errOut, code := runIn(t, root, "issue", "story", "I-0002"); code == 0 || !strings.Contains(errOut, "I-0002 not found") {
		t.Errorf("without --story the issue is read from the root, which has no I-0002: %d %s", code, errOut)
	}
	before := readIssueFile(t, wt, "I-0002-only-on-the-branch.md")
	out, errOut, code := runIn(t, root, "issue", "story", "I-0002", "--story", "S-0005")
	if code != 0 || out != "S-0001 Only on the branch\n  wip/kanban/stories/S-0001-only-on-the-branch.md\n" {
		t.Fatalf("issue story --story: %d %q %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(wt, "wip", "kanban", "stories")); err == nil {
		t.Error("the story should be made in the main checkout's wip, not the worktree's")
	}
	if string(before) != string(readIssueFile(t, wt, "I-0002-only-on-the-branch.md")) {
		t.Error("issue story should not change the issue in the worktree")
	}
	out, _, _ = runIn(t, root, "issue", "list", "--story", "S-0005", "--json")
	var list []map[string]any
	if json.Unmarshal([]byte(out), &list) != nil || len(list) != 1 || list[0]["title"] != "Only on the branch" || list[0]["story"] != "S-0001" {
		t.Errorf("list --story --json should give the worktree's issue, linked by S-0001: %s", out)
	}
	_, errOut, code = runIn(t, root, "issue", "story", "I-0002", "--story", "S-0005")
	if code == 0 || !strings.Contains(errOut, "already linked by open story S-0001") {
		t.Errorf("the worktree's issue, now linked, should be refused: %d %s", code, errOut)
	}
}

func TestIssueCommands(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "design", "conventions"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "design", "conventions", "README.md"), []byte("# c\n"), 0o644)
	out, errOut, code := runIn(t, root, "issue", "new", "Go not on PATH", "--class", "blocker", "--cost", "3m")
	if code != 0 || !strings.HasPrefix(out, "I-0001 Go not on PATH") {
		t.Fatalf("new: %d %s %s", code, out, errOut)
	}
	_, errOut, code = runIn(t, root, "issue", "new", "no class")
	if code == 0 || !strings.Contains(errOut, "class") {
		t.Fatalf("class required: %s", errOut)
	}
	out, _, code = runIn(t, root, "issue", "bump", "I-0001", "--cost", "5m", "--note", "again")
	if code != 0 || !strings.Contains(out, "count 2, avg cost 4m") {
		t.Fatalf("bump: %s", out)
	}
	sum, _ := os.ReadFile(filepath.Join(root, "design", "issues", "summary.md"))
	if !strings.Contains(string(sum), "| [I-0001](I-0001-go-not-on-path.md) | blocker | Go not on PATH | 2 | 4m | 8m |") {
		t.Errorf("summary after bump:\n%s", sum)
	}
	out, _, _ = runIn(t, root, "prime")
	if !strings.HasSuffix(strings.TrimSpace(out), "design/issues/summary.md") {
		t.Errorf("prime should end with the summary when issues are open:\n%s", out)
	}
	out, _, _ = runIn(t, root, "prime", "--cat")
	if !strings.Contains(out, "open issues\n===") || !strings.Contains(out, "[I-0001]") {
		t.Errorf("prime --cat should include the table:\n%s", out)
	}
	out, _, code = runIn(t, root, "issue", "list", "--json")
	var list []map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &list) != nil || len(list) != 1 || list[0]["count"].(float64) != 2 {
		t.Fatalf("list json: %s", out)
	}
	out, _, code = runIn(t, root, "issue", "close", "I-0001", "--reason", "fixed")
	if code != 0 || !strings.Contains(out, "I-0001 closed") {
		t.Fatalf("close: %s", out)
	}
	out, _, _ = runIn(t, root, "issue", "summary")
	if strings.TrimSpace(out) != "no open issues" {
		t.Errorf("summary after close: %s", out)
	}
	out, _, _ = runIn(t, root, "prime")
	if strings.Contains(out, "summary.md") {
		t.Error("prime should not list the summary when nothing is open")
	}
	out, _, _ = runIn(t, root, "check")
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "issues.") {
			t.Errorf("check reports an issues finding after the commands: %s", line)
		}
	}
}

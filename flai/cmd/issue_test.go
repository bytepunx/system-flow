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
		{"issue", "new", "Fixture was ignored", "--class", "defect", "--cost", "20m"},
		{"issue", "bump", "I-0001", "--cost", "20m"},
		{"issue", "new", "Lint on the host is v1", "--class", "efficiency"},
		{"issue", "new", "Go not on PATH", "--class", "blocker"},
		{"epic", "new", "Upkeep"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	summary := readIssueFile(t, root, "summary.md")
	out, errOut, code := runIn(t, root, "issue", "story", "I-0001")
	if code != 0 || out != "S-0001 Fixture was ignored\n  wip/kanban/stories/S-0001-fixture-was-ignored.md\n" {
		t.Fatalf("issue story I-0001: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "issue", "story", "I-0002", "--epic", "E-0001", "--json")
	var made map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &made) != nil {
		t.Fatalf("issue story I-0002 --json: %d %s %s", code, out, errOut)
	}
	want := map[string]any{"id": "S-0002", "title": "Lint on the host is v1", "nature": "improvement", "path": "wip/kanban/stories/S-0002-lint-on-the-host-is-v1.md", "issue": "I-0002", "draft": true, "committed": false}
	for k, v := range want {
		if made[k] != v {
			t.Errorf("issue story --json %s = %v, want %v", k, made[k], v)
		}
	}
	if _, ok := made["cost_of_delay"]; ok {
		t.Errorf("an issue with no cost gives no cost of delay: %v", made["cost_of_delay"])
	}
	for file, nature := range map[string]string{"S-0001-fixture-was-ignored.md": "remediation", "S-0002-lint-on-the-host-is-v1.md": "improvement"} {
		data, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "stories", file))
		if !strings.Contains(string(data), "\nnature: "+nature+"\n") || !strings.Contains(string(data), "\nstatus: backlog\n") || !strings.Contains(string(data), "\ndraft: true\n") {
			t.Errorf("%s should be a draft %s story in backlog:\n%s", file, nature, data)
		}
	}
	// S-0203: the issue's cost and count give the story's time lost per
	// cycle, set by flai and explained in its Notes
	data, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "stories", "S-0001-fixture-was-ignored.md"))
	for _, want := range []string{"\n    time_lost_per_cycle: 40m\n    by: flai\n", "time_lost_per_cycle 40m: 20m per occurrence × 2 occurrences ÷ 1 cycle of 168h"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("S-0001 should carry %q:\n%s", want, data)
		}
	}
	if data, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "stories", "S-0002-lint-on-the-host-is-v1.md")); strings.Contains(string(data), "cost_of_delay") {
		t.Errorf("an issue with no cost gives no cost of delay:\n%s", data)
	}
	epic, _ := os.ReadFile(filepath.Join(root, "wip", "kanban", "epics", "E-0001-upkeep.md"))
	if !strings.Contains(string(epic), "S-0002") {
		t.Errorf("the epic should link the story made under it:\n%s", epic)
	}
	for f, id := range map[string]string{"I-0001-fixture-was-ignored.md": "S-0001", "I-0002-lint-on-the-host-is-v1.md": "S-0002"} {
		if data := readIssueFile(t, root, f); !strings.Contains(string(data), "## Remediation\n\nStory "+id+" remediates this issue, created from it at 2026-09-15T21:00:00Z.\n") {
			t.Errorf("the issue's Remediation section should name %s:\n%s", id, data)
		}
	}
	if string(summary) != string(readIssueFile(t, root, "summary.md")) {
		t.Error("issue story should not change design/issues/summary.md")
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
// manifest's owner, committed on its own with its trailer, and with it the
// issue that now names it (S-0203). An issue in another story's worktree is
// named there and left for that story to commit.
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
	if !strings.Contains(show, "S-0001-fixture-was-ignored.md") || !strings.Contains(show, "E-0001-workbench.md") || !strings.Contains(show, "I-0001-fixture-was-ignored.md") || !strings.Contains(show, "3 files changed") {
		t.Errorf("the commit should hold the story, its epic, and the issue only: %s", show)
	}
	if data := readIssueFile(t, root, "I-0001-fixture-was-ignored.md"); !strings.Contains(string(data), "Story S-0001 remediates this issue") {
		t.Errorf("the committed issue should name the story:\n%s", data)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("nothing should be left uncommitted: %q", st)
	}

	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0005")
	_ = os.MkdirAll(filepath.Join(wt, "design", "issues"), 0o755)
	for _, f := range []string{"system-flow.yaml", "design/issues/I-0001-fixture-was-ignored.md"} {
		data, _ := os.ReadFile(filepath.Join(root, f))
		_ = os.WriteFile(filepath.Join(wt, f), data, 0o644)
	}
	if _, errOut, code := runIn(t, wt, "issue", "new", "Only on the branch", "--class", "defect", "--story", "S-0005"); code != 0 {
		t.Fatalf("issue new in the worktree: %s", errOut)
	}
	out, errOut, code = runIn(t, root, "issue", "story", "I-0002", "--story", "S-0005", "--autocommit", "--json")
	if code != 0 || json.Unmarshal([]byte(out), &made) != nil || made.ID != "S-0002" || !made.Committed {
		t.Fatalf("issue story --story --autocommit: %d %s %s", code, out, errOut)
	}
	if show := gitIn(t, root, "show", "--stat", "--format=%s", "HEAD"); !strings.Contains(show, "S-0002-only-on-the-branch.md") || !strings.Contains(show, "1 file changed") {
		t.Errorf("the commit should hold the story only, not the worktree's issue: %s", show)
	}
	if data := readIssueFile(t, wt, "I-0002-only-on-the-branch.md"); !strings.Contains(string(data), "Story S-0002 remediates this issue") {
		t.Errorf("the worktree's issue should name the story:\n%s", data)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("nothing should be left uncommitted in the main checkout: %q", st)
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
// there and names the story there, and the story made from one goes to wip/
// in the main checkout.
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
	out, errOut, code := runIn(t, root, "issue", "story", "I-0002", "--story", "S-0005")
	if code != 0 || out != "S-0001 Only on the branch\n  wip/kanban/stories/S-0001-only-on-the-branch.md\n" {
		t.Fatalf("issue story --story: %d %q %s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(wt, "wip", "kanban", "stories")); err == nil {
		t.Error("the story should be made in the main checkout's wip, not the worktree's")
	}
	if data := readIssueFile(t, wt, "I-0002-only-on-the-branch.md"); !strings.Contains(string(data), "Story S-0001 remediates this issue") {
		t.Errorf("the issue in the worktree should name the story:\n%s", data)
	}
	if _, err := os.Stat(filepath.Join(root, "design", "issues", "I-0002-only-on-the-branch.md")); err == nil {
		t.Error("the worktree's issue should not be written to the main checkout")
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

// I-0065: an issue another story recorded, committed on its branch or not
// yet committed in its worktree, holds its number, so a third story's issue
// takes the next one.
func TestIssueNewNumbersPastEveryStorysIssues(t *testing.T) {
	if testing.Short() {
		t.Skip("needs git")
	}
	root := bodyProject(t)
	t.Setenv("FLAI_STORY", "")
	writeIssue := func(dir, name string) {
		t.Helper()
		path := filepath.Join(dir, "design", "issues", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("---\nid: "+name[:6]+"\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	worktree := func(story string) string {
		t.Helper()
		wt := filepath.Join(root, ".flai-cache", "worktrees", story)
		gitIn(t, root, "worktree", "add", "-q", "-b", "story/"+story, wt)
		return wt
	}
	writeIssue(root, "I-0001-on-main.md")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "on main")

	one := worktree("S-0001")
	writeIssue(one, "I-0002-committed-on-its-branch.md")
	gitIn(t, one, "add", "-A")
	gitIn(t, one, "commit", "-q", "-m", "one")
	gitIn(t, root, "worktree", "remove", one)
	writeIssue(worktree("S-0002"), "I-0003-uncommitted-in-its-worktree.md")

	out, errOut, code := runIn(t, worktree("S-0003"), "issue", "new", "Third story's issue", "--class", "defect")
	if code != 0 || !strings.HasPrefix(out, "I-0004 Third story's issue\n") {
		t.Errorf("issue new in a third story's worktree should be I-0004: %d %q %s", code, out, errOut)
	}
}

// S-0224: the analyzer files a finding with its impact and the report that
// found it; with --report an open issue of the same title is bumped, not
// duplicated, and --json says which; bump with --report links a report and
// updates the impact; anything bad is refused with nothing written.
func TestIssueNewAndBumpWithImpactAndReport(t *testing.T) {
	root := issueProject(t)
	report := "design/analysis/2026-10-06-bottlenecks.md"
	issueFiles := func() []string {
		m, _ := filepath.Glob(filepath.Join(root, "design", "issues", "I-*.md"))
		return m
	}
	filed := func(args ...string) map[string]any {
		t.Helper()
		out, errOut, code := runIn(t, root, append([]string{"issue", "new"}, args...)...)
		var got map[string]any
		if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
			t.Fatalf("issue new %v: %d %s %s", args, code, out, errOut)
		}
		return got
	}
	got := filed("Review waits a day", "--class", "efficiency", "--revenue-per-week", "1200", "--time-lost-per-cycle", "240m",
		"--evidence", "Twelve stories waited 18h on average.", "--report", report, "--story", "S-0224", "--json")
	if got["id"] != "I-0001" || got["outcome"] != "opened" || !strings.HasSuffix(got["path"].(string), "I-0001-review-waits-a-day.md") {
		t.Errorf("new --report --json: %v", got)
	}
	data := string(readIssueFile(t, root, "I-0001-review-waits-a-day.md"))
	for _, want := range []string{"Story: S-0224.\nReport: " + report + ".\nFirst occurrence.\n",
		"## Impact\nTwelve stories waited 18h on average.\n\n- revenue_per_week: 1200\n- time_lost_per_cycle: 4h\n\n## Remediation\n",
		"Found by the analysis in [2026-10-06-bottlenecks.md](../analysis/2026-10-06-bottlenecks.md).\n"} {
		if !strings.Contains(data, want) {
			t.Errorf("the issue should hold %q:\n%s", want, data)
		}
	}

	got = filed("Review waits a day", "--class", "efficiency", "--penalty-per-week", "300", "--report", "design/analysis/2026-10-13-all.md", "--json")
	if got["id"] != "I-0001" || got["outcome"] != "bumped" || got["count"].(float64) != 2 || len(issueFiles()) != 1 {
		t.Errorf("the same title with --report is bumped: %v %v", got, issueFiles())
	}
	out, errOut, code := runIn(t, root, "issue", "new", "Review waits a day", "--class", "efficiency", "--report", "design/analysis/2026-10-13-all.md")
	if code != 0 || out != "I-0001 Review waits a day\n  design/issues/I-0001-review-waits-a-day.md\n  already recorded from this report; nothing changed\n" {
		t.Errorf("the same report again: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "issue", "new", "Review waits a day", "--class", "efficiency")
	if code != 0 || !strings.HasPrefix(out, "I-0002 Review waits a day\n") || strings.Contains(out, "opened") {
		t.Errorf("without --report the same title is a new issue, said as before: %d %q %s", code, out, errOut)
	}

	out, errOut, code = runIn(t, root, "issue", "bump", "I-0002", "--report", report, "--time-lost-per-cycle", "2h", "--evidence", "Seen in the same report.")
	if code != 0 || out != "I-0002 count 2, avg cost -\n" {
		t.Errorf("bump with --report: %d %q %s", code, out, errOut)
	}
	data = string(readIssueFile(t, root, "I-0002-review-waits-a-day.md"))
	for _, want := range []string{"Report: " + report + ".\nOccurred again.\n\n## Impact\nSeen in the same report.\n\n- time_lost_per_cycle: 2h\n\n## Remediation\n",
		"](../analysis/2026-10-06-bottlenecks.md).\n"} {
		if !strings.Contains(data, want) {
			t.Errorf("bump with --report should link the report and add the impact, %q:\n%s", want, data)
		}
	}

	before := map[string]string{}
	for _, f := range append(issueFiles(), filepath.Join(root, "design", "issues", "summary.md")) {
		b, _ := os.ReadFile(f)
		before[f] = string(b)
	}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"issue", "new", "Bad", "--class", "defect", "--revenue-per-week", "lots"}, "not an amount of zero or more"},
		{[]string{"issue", "new", "Bad", "--class", "defect", "--penalty-per-week", "-1"}, "not an amount of zero or more"},
		{[]string{"issue", "new", "Bad", "--class", "defect", "--time-lost-per-cycle", "0"}, "not a duration longer than zero"},
		{[]string{"issue", "new", "Bad", "--class", "defect", "--report", "design/issues/summary.md"}, "not an analysis report"},
		{[]string{"issue", "new", "Review waits a day", "--class", "efficiency", "--report", report, "--time-lost-per-cycle", "soon"}, "not a duration"},
		{[]string{"issue", "bump", "I-0001", "--penalty-per-week", "x"}, "not an amount"},
		{[]string{"issue", "bump", "I-0001", "--time-lost-per-cycle", "-4h"}, "not a duration longer than zero"},
		{[]string{"issue", "bump", "I-0001", "--report", "design/analysis/notes.txt"}, "not an analysis report"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v should be refused with %q: %d %s", c.args, c.want, code, errOut)
		}
	}
	if len(issueFiles()) != 2 {
		t.Errorf("a refusal wrote an issue: %v", issueFiles())
	}
	for f, b := range before {
		if after, _ := os.ReadFile(f); string(after) != b {
			t.Errorf("a refusal changed %s:\n%s", f, after)
		}
	}
}

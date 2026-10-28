package cmd

import (
	"encoding/json"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scopeFixture copies the good fixture, whose S-004 is in progress, and
// gives the archived S-002 a branch that was never merged (I-0057).
func scopeFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	refs := filepath.Join(root, ".git", "refs", "heads", "story")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refs, "S-002"), []byte("0123456789abcdef0123456789abcdef01234567\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// absentMainWorktree is S-004's worktree mounted without its main checkout,
// as in an agent container (I-0219): the good fixture in a folder S-004
// whose .git file names a main checkout's git directory that is not there.
// It returns the worktree and the absent main checkout.
func absentMainWorktree(t *testing.T) (worktree, main string) {
	t.Helper()
	dir := t.TempDir()
	worktree, main = filepath.Join(dir, "S-004"), filepath.Join(dir, "main")
	if err := os.CopyFS(worktree, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		".gitignore": ".flai-cache/\n",
		".git":       "gitdir: " + filepath.Join(main, ".git", "worktrees", "S-004") + "\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(worktree, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return worktree, main
}

// I-0219: flai check in a worktree whose main checkout is absent refuses,
// naming the main checkout and the host, and reads and writes nothing,
// rather than checking the worktree's tracked wip/ for the board's.
func TestCheckRefusesAWorktreeWhoseMainCheckoutIsAbsent(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	worktree, main := absentMainWorktree(t)
	before := treeOf(t, worktree)
	for _, args := range [][]string{
		{"check", "--strict"},
		{"check", "--story", "S-004", "--record-issues"},
		{"check", worktree},
	} {
		out, errOut, code := runIn(t, worktree, args...)
		// the whole error, one line, as the command boundary logs it
		want := `"err":"flai check needs the main checkout ` + main + ", whose wip/ is the board, which is not here: " + worktree + ` is a story worktree mounted without it, as in an agent container; run flai check on the host, or call the MCP tool verify, which runs it there"}`
		if code == 0 || out != "" || !strings.Contains(errOut, want) {
			t.Errorf("%v: exit %d\nout %q\nerr %q, want %q", args, code, out, errOut, want)
		}
	}
	if after := treeOf(t, worktree); !maps.Equal(before, after) {
		t.Error("a refused check changed the worktree")
	}
}

// S-0249: --story S-nnnn passes over the findings outside the story, says
// how many in the summary, and marks them in --json.
func TestCheckStoryPassesOverFindingsOutsideIt(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := scopeFixture(t)
	out, errOut, code := runIn(t, root, "check", "--strict")
	if code != 1 || !strings.Contains(out, "warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it\n") {
		t.Fatalf("unscoped --strict should fail on S-002's branch: %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "check", "--strict", "--story", "S-4")
	if code != 0 {
		t.Fatalf("scoped to S-004, --strict should pass: %d\n%s%s", code, out, errOut)
	}
	for _, want := range []string{
		"warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it (outside S-004)\n",
		"10 items checked, 0 errors, 2 warnings; 2 outside S-004, notes the run passes over\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output lacks %q:\n%s", want, out)
		}
	}

	out, _, code = runIn(t, root, "check", "--strict", "--story", "S-004", "--json")
	var res struct {
		Warnings int `json:"warnings"`
		Advisory int `json:"advisory"`
		Outside  int `json:"outside"`
		Findings []struct {
			Rule    string `json:"rule"`
			Outside bool   `json:"outside"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("json: %d %v %s", code, err, out)
	}
	if res.Outside != 2 || res.Warnings != 2 || res.Advisory != 0 || len(res.Findings) != 2 {
		t.Fatalf("json counts: %+v", res)
	}
	for _, f := range res.Findings {
		if !f.Outside {
			t.Errorf("%s should carry outside: %s", f.Rule, out)
		}
	}
	// Unscoped, --json has no outside at all.
	out, _, _ = runIn(t, root, "check", "--json")
	if strings.Contains(out, `"outside"`) {
		t.Errorf("unscoped json should not mention outside:\n%s", out)
	}

	if _, errOut, code := runIn(t, root, "check", "--story", "T-003"); code == 0 || !strings.Contains(errOut, "scope the check to T-003: it is a task; name a story") {
		t.Errorf("a task is not a scope: %d %s", code, errOut)
	}
}

// I-0174: --story takes several stories, comma-separated or repeated, as a
// merge queue batch names them; a finding outside both is a note, the
// summary names both, and --record-issues, which records one close-out's
// findings, is refused before checking.
func TestCheckStoryTakesSeveralStories(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", map[string]string{"S-006": "docs"})
	refs := filepath.Join(root, ".git", "refs", "heads", "story")
	if err := os.MkdirAll(refs, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refs, "S-002"), []byte("0123456789abcdef0123456789abcdef01234567\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const unaccepted = "warning: story.unaccepted: S-002 is done but its branch story/S-002 was never merged; merge or delete it"
	out, errOut, code := runIn(t, root, "check", "--strict")
	if code != 1 || !strings.Contains(out, unaccepted+"\n") {
		t.Fatalf("unscoped --strict should fail on S-002's branch: %d\n%s%s", code, out, errOut)
	}
	for _, args := range [][]string{{"--story", "S-004,S-006"}, {"--story", "S-4", "--story", "S-006"}} {
		out, errOut, code = runIn(t, root, append([]string{"check", "--strict"}, args...)...)
		if code != 0 {
			t.Fatalf("%v: scoped to S-004 and S-006, --strict should pass: %d\n%s%s", args, code, out, errOut)
		}
		for _, want := range []string{unaccepted + " (outside S-004, S-006)\n", "; 2 outside S-004, S-006, notes the run passes over\n"} {
			if !strings.Contains(out, want) {
				t.Errorf("%v: output lacks %q:\n%s", args, want, out)
			}
		}
	}

	out, errOut, code = runIn(t, root, "check", "--strict", "--story", "S-004,S-006", "--record-issues")
	if code == 0 || !strings.Contains(errOut, "--record-issues records the findings outside one story's close-out, and --story names 2 (S-004, S-006)") || strings.Contains(out, "items checked") {
		t.Errorf("--record-issues with two stories should be refused before checking: %d\n%s%s", code, out, errOut)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("a refused run should write nothing under design/issues, got %v", names)
	}
}

// gitScript answers the git commands flai stream diff and git status run,
// for a story branch that changes design/system/overview.md, or, with
// porcelain, a worktree that has it uncommitted; statusErr fails git status.
type gitScript struct {
	changed, porcelain string
	statusErr          error
}

func (g gitScript) Run(_, _ string, args ...string) (string, error) {
	cmd := strings.Join(args, " ")
	switch {
	case cmd == "rev-parse --is-inside-work-tree":
		return "true", nil
	case cmd == "rev-parse --verify --quiet refs/heads/story/S-004":
		return "abc", nil
	case cmd == "rev-parse --abbrev-ref HEAD":
		return "main", nil
	case strings.HasPrefix(cmd, "merge-base "), strings.HasPrefix(cmd, "rev-parse --short "):
		return "base", nil
	case strings.HasPrefix(cmd, "rev-list --count "):
		return "1", nil
	case strings.HasPrefix(cmd, "diff --name-status "):
		if g.changed == "" {
			return "", nil
		}
		return "M\x00" + g.changed + "\x00", nil
	case strings.HasPrefix(cmd, "diff --numstat "):
		if g.changed == "" {
			return "", nil
		}
		return "1\t1\t" + g.changed + "\x00", nil
	case strings.HasPrefix(cmd, "diff --no-color "):
		return "@@ -1 +1 @@\n", nil
	case cmd == "status --porcelain":
		return g.porcelain, g.statusErr
	}
	return "", errors.New("unscripted: git " + cmd)
}

func (g gitScript) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (gitScript) LookPath(name string) (string, error) { return name, nil }

// S-0249: a finding on a path the story's branch changes, committed or
// uncommitted in its worktree, is the story's and still fails --strict.
func TestCheckStoryFailsOnAPathItChanges(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	for _, tc := range []struct {
		name string
		git  gitScript
		code int
	}{
		{"committed on its branch", gitScript{changed: "design/system/overview.md"}, 1},
		{"uncommitted in its worktree", gitScript{porcelain: " M design/system/overview.md\n"}, 1},
		{"changed elsewhere", gitScript{changed: "flai/main.go"}, 0},
		{"a worktree that does not read", gitScript{statusErr: errors.New("fatal: not a git repository")}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeFixture(t)
			overview := filepath.Join(root, "design/system/overview.md")
			data, _ := os.ReadFile(overview)
			if err := os.WriteFile(overview, []byte(strings.Replace(string(data), "title: Overview\n", "", 1)), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Join(root, ".flai-cache", "worktrees", "S-004"), 0o755); err != nil {
				t.Fatal(err)
			}
			out, errOut, code := runWithApp(t, &app{cwd: root, runner: tc.git}, "check", "--strict", "--story", "S-004")
			if code != tc.code {
				t.Fatalf("exit %d, want %d:\n%s%s", code, tc.code, out, errOut)
			}
			line := "design/system/overview.md:1: warning: doc.title: front matter needs a title"
			if tc.git.statusErr != nil {
				if !strings.Contains(errOut, "read what S-004 has uncommitted in") || strings.Contains(out, "items checked") {
					t.Errorf("an unreadable worktree should stop the run:\n%s%s", out, errOut)
				}
				return
			}
			if tc.code == 1 && !strings.Contains(out, line+"\n") {
				t.Errorf("doc.title should be the story's:\n%s", out)
			}
			if tc.code == 0 && !strings.Contains(out, line+" (outside S-004)\n") {
				t.Errorf("doc.title should be outside the story:\n%s", out)
			}
		})
	}
}

// issueFiles is the files under design/issues in root, by name.
func issueFiles(t *testing.T, root string) []string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(root, "design", "issues"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// S-0249: --record-issues records each rule's findings outside the story in
// one issue, once per story and findings, and bumps it for another story.
func TestCheckRecordIssuesOpensThenBumpsOncePerStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := scopeFixture(t)
	issue := filepath.Join(root, "design", "issues", "I-0002-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md")
	finding := "`wip/archive/kanban/stories/S-002-two.md`: S-002 is done but its branch story/S-002 was never merged; merge or delete it"
	read := func() string {
		t.Helper()
		data, err := os.ReadFile(issue)
		if err != nil {
			t.Fatalf("the story.unaccepted issue: %v (files: %v)", err, issueFiles(t, root))
		}
		return string(data)
	}

	out, errOut, code := runIn(t, root, "check", "--strict", "--story", "S-004", "--record-issues")
	if code != 0 {
		t.Fatalf("first run: exit %d\n%s%s", code, out, errOut)
	}
	if !strings.Contains(out, "recorded story.unaccepted outside S-004 in I-0002 (opened)\n") {
		t.Errorf("first run should say it opened the issue:\n%s", out)
	}
	doc := read()
	for _, want := range []string{
		"title: \"flai check finds `story.unaccepted` outside the story at close-out\"\n",
		"class: efficiency\n", "count: 1\n", "Story: S-0004.\n", finding + "\n",
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("first run: issue lacks %q:\n%s", want, doc)
		}
	}
	if summary, err := os.ReadFile(filepath.Join(root, "design", "issues", "summary.md")); err != nil || !strings.Contains(string(summary), "[I-0002]") {
		t.Errorf("summary.md should list I-0002: %v\n%s", err, summary)
	}

	out, errOut, code = runIn(t, root, "check", "--strict", "--story", "S-004", "--record-issues")
	if code != 0 || !strings.Contains(out, "recorded story.unaccepted outside S-004 in I-0002 (already recorded, count 1)\n") {
		t.Fatalf("second run: exit %d\n%s%s", code, out, errOut)
	}
	if again := read(); again != doc {
		t.Errorf("an identical run should leave the issue as it was:\n%s", again)
	}

	out, errOut, code = runIn(t, root, "check", "--strict", "--story", "S-003", "--record-issues", "--json")
	var res struct {
		Outside  int `json:"outside"`
		Recorded []struct {
			Rule, Issue, Path, Outcome string
			Count                      int
		} `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("another story: exit %d %v\n%s%s", code, err, out, errOut)
	}
	if res.Outside != 2 || len(res.Recorded) != 2 {
		t.Fatalf("json should keep outside and list both rules: %+v", res)
	}
	got := res.Recorded[1]
	if got.Rule != "story.unaccepted" || got.Issue != "I-0002" || got.Outcome != "bumped" || got.Count != 2 ||
		got.Path != "design/issues/I-0002-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md" {
		t.Errorf("another story should bump the issue to 2: %+v", got)
	}
	doc = read()
	if !strings.Contains(doc, "count: 2\n") || !strings.Contains(doc, "Story: S-0003.\n") || strings.Count(doc, finding) != 2 {
		t.Errorf("another story should add its instance:\n%s", doc)
	}
	if names := issueFiles(t, root); len(names) != 3 {
		t.Errorf("one issue per rule and the summary, got %v", names)
	}
}

// overlapFixture copies the good fixture, whose S-004 is in progress, gives
// S-004 the touches s004 when it is not empty, and adds each story of stories
// in progress under E-001, with its touches, its narrative, and its row in
// the index, the in-progress limit raised to hold them: a board whose only
// findings a run scoped to S-004 notes are the wip.overlap ones.
func overlapFixture(t *testing.T, s004 string, stories map[string]string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	edit := func(rel, old, new string) {
		t.Helper()
		p := filepath.Join(root, rel)
		data, err := os.ReadFile(p)
		if err != nil || !strings.Contains(string(data), old) {
			t.Fatalf("%s has no %q: %v", rel, old, err)
		}
		if err := os.WriteFile(p, []byte(strings.Replace(string(data), old, new, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	edit("wip/kanban/board.md", "in-progress: 2\n", "in-progress: 5\n")
	if s004 != "" {
		edit("wip/kanban/stories/S-004-four.md", "tags: []\n", "tags: []\ntouches: ["+s004+"]\n")
	}
	narrative, err := os.ReadFile(filepath.Join(root, "wip/agents/S-004.md"))
	if err != nil {
		t.Fatal(err)
	}
	four, err := os.ReadFile(filepath.Join(root, "wip/kanban/stories/S-004-four.md"))
	if err != nil {
		t.Fatal(err)
	}
	if s004 != "" {
		four = []byte(strings.Replace(string(four), "touches: ["+s004+"]\n", "", 1))
	}
	for id, touches := range stories {
		story := strings.NewReplacer("S-004", id, "- T-003 T3\n", "", "tags: []\n", "tags: []\ntouches: ["+touches+"]\n").Replace(string(four))
		if err := os.WriteFile(filepath.Join(root, "wip/kanban/stories", id+"-four.md"), []byte(story), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "wip/agents", id+".md"), []byte(strings.ReplaceAll(string(narrative), "S-004", id)), 0o644); err != nil {
			t.Fatal(err)
		}
		edit("wip/kanban/epics/E-001-epic.md", "- S-005 Five\n", "- S-005 Five\n- "+id+" Four\n")
		edit("wip/agents/index.md", "| [S-004](S-004.md) | Four | in-progress | bot | 2026-08-31T10:00:00Z |\n",
			"| [S-004](S-004.md) | Four | in-progress | bot | 2026-08-31T10:00:00Z |\n| ["+id+"]("+id+".md) | Four | in-progress | bot | 2026-08-31T10:00:00Z |\n")
	}
	return root
}

// I-0076, ADR-0115: two other stories in progress overlap; a close-out of a
// third neither notes the overlap nor records it in an issue.
func TestCheckRecordIssuesLeavesOutAnOverlapBetweenTwoOtherStories(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", map[string]string{"S-006": "flai", "S-007": "flai/cmd"})
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, _ := runWithApp(t, &app{cwd: root, runner: git}, "check")
	if !strings.Contains(out, "warning: wip.overlap: S-006 touches flai, which S-007 (in progress) also touches as flai/cmd\n") {
		t.Fatalf("unscoped, the run should report the overlap of S-006 and S-007:\n%s%s", out, errOut)
	}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int `json:"outside"`
		Findings []struct {
			Rule string `json:"rule"`
		} `json:"findings"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "wip.overlap") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the overlap should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened or bumped, got %v", names)
	}
}

// ADR-0115: an overlap that names the closing story is printed as a note
// outside it, and is not recorded in an issue.
func TestCheckRecordIssuesNotesButDoesNotRecordTheStorysOwnOverlap(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "flai", map[string]string{"S-006": "flai/cmd"})
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: gitScript{changed: "wip/kanban/epics/E-001-epic.md"}}, "check", "--strict", "--story", "S-004", "--record-issues")
	if code != 0 {
		t.Fatalf("a note should not fail --strict: exit %d\n%s%s", code, out, errOut)
	}
	if !strings.Contains(out, "warning: wip.overlap: S-004 touches flai, which S-006 (in progress) also touches as flai/cmd (outside S-004)\n") {
		t.Errorf("the overlap should be a note outside S-004:\n%s", out)
	}
	if strings.Contains(out, "recorded ") {
		t.Errorf("the overlap should not be recorded:\n%s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened, got %v", names)
	}
}

// I-0078, ADR-0122: another story, cancelled from backlog, waits to be
// archived; a close-out of S-004 neither notes its item.archive nor records
// it in an issue, while the unscoped check still warns on it.
func TestCheckRecordIssuesLeavesOutAnItemArchiveOutsideTheStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", nil)
	story := "---\nid: S-006\ntype: story\nnature: feature\ntitle: Six\nstatus: cancelled\nparent: E-001\nowner: agent\n" +
		"created: 2026-08-25T09:00:00Z\nupdated: 2026-08-26T10:00:00Z\ntransitions:\n  - to: cancelled\n    at: 2026-08-26T10:00:00Z\n    by: alex\ntags: []\n" +
		"---\n\n# S-006 Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/stories/S-006-six.md"), []byte(story), 0o644); err != nil {
		t.Fatal(err)
	}
	epic := filepath.Join(root, "wip/kanban/epics/E-001-epic.md")
	data, err := os.ReadFile(epic)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(epic, []byte(strings.Replace(string(data), "- S-005 Five\n", "- S-005 Five\n- S-006 Six\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, _ := runWithApp(t, &app{cwd: root, runner: git}, "check")
	if !strings.Contains(out, "warning: item.archive: S-006 is cancelled; run flai archive\n") {
		t.Fatalf("unscoped, the run should report S-006's item.archive:\n%s%s", out, errOut)
	}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "item.archive") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the item.archive should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened or bumped, got %v", names)
	}
}

// I-0123, ADR-0133: in progress is over its limit on the main checkout's
// board, as at S-0321's and S-0339's close-outs; a close-out of S-004 neither
// notes the board.wip-limit nor records it in an issue, while the unscoped
// check still warns on it and fails --strict.
func TestCheckRecordIssuesLeavesOutABoardWIPLimitOutsideTheStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", map[string]string{"S-006": "docs", "S-007": "template"})
	board := filepath.Join(root, "wip/kanban/board.md")
	data, err := os.ReadFile(board)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(board, []byte(strings.Replace(string(data), "in-progress: 5\n", "in-progress: 2\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict")
	if code == 0 || !strings.Contains(out, "warning: board.wip-limit: 3 stories in in-progress, limit 2\n") {
		t.Fatalf("unscoped, the run should warn on the board's in-progress limit and fail --strict: exit %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "board.wip-limit") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the board.wip-limit should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened or bumped, got %v", names)
	}
}

// I-0134, ADR-0136: an open issue that no open story links has waited past
// issues.story_after, as I-0052 had at S-0343's, S-0289's, S-0342's, and
// S-0361's close-outs; a close-out of S-004, which did not touch it, neither
// notes the issues.no-story nor records it in an issue, while the unscoped
// check still warns on it and fails --strict.
func TestCheckRecordIssuesLeavesOutAnIssuesNoStoryOutsideTheStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", nil)
	dir := filepath.Join(root, "design", "issues")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	issue := "---\nid: I-0052\ntitle: Waits for a story\nclass: defect\nstatus: open\ncount: 1\nfirst_reported: 2026-08-01T12:00:00Z\n" +
		"last_reported: 2026-08-01T12:00:00Z\nupdated: 2026-08-01T12:00:00Z\n---\n\n# I-0052 Waits for a story\n"
	if err := os.WriteFile(filepath.Join(dir, "I-0052-waits-for-a-story.md"), []byte(issue), 0o644); err != nil {
		t.Fatal(err)
	}
	summary := "# Issues\n\n- [I-0052](I-0052-waits-for-a-story.md)\n"
	if err := os.WriteFile(filepath.Join(dir, "summary.md"), []byte(summary), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict")
	if code == 0 || !strings.Contains(out, "warning: issues.no-story: I-0052 has been open since 2026-08-01 with no open story linking it") {
		t.Fatalf("unscoped, the run should warn on I-0052 with no story and fail --strict: exit %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "issues.no-story") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the issues.no-story should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); !slices.Equal(names, []string{"I-0052-waits-for-a-story.md", "summary.md"}) {
		t.Errorf("no issue should be opened, got %v", names)
	}
	data, err := os.ReadFile(filepath.Join(dir, "I-0052-waits-for-a-story.md"))
	if err != nil || string(data) != issue {
		t.Errorf("I-0052 should be left as it was: %v\n%s", err, data)
	}
	if data, err := os.ReadFile(filepath.Join(dir, "summary.md")); err != nil || string(data) != summary {
		t.Errorf("summary.md should be left as it was: %v\n%s", err, data)
	}
}

// I-0135, ADR-0137: an open thread is anchored on an archived story, as
// TH-0389 was on S-0338 at S-0343's, S-0289's, and S-0342's close-outs; a
// close-out of S-004, which did not touch it, neither notes the
// threads.archived nor records it in an issue, while the unscoped check still
// warns on it and fails --strict.
func TestCheckRecordIssuesLeavesOutAThreadsArchivedOutsideTheStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", nil)
	thread := "---\nid: TH-0389\ntitle: Waits for the operator\nanchor:\n  path: wip/kanban/stories/S-003-three.md\n  item: S-003\n" +
		"status: open\nparticipants: [orchestrator]\ncreated: 2026-08-30T09:00:00Z\nupdated: 2026-08-30T09:00:00Z\n---\n\n" +
		"# TH-0389 Waits for the operator\n\nOn S-003 (wip/kanban/stories/S-003-three.md).\n\n## Entries\n\n### 2026-08-30T09:00:00Z orchestrator\nOver to the operator.\n"
	path := filepath.Join(root, "wip", "threads", "TH-0389-waits-for-the-operator.md")
	if err := os.WriteFile(path, []byte(thread), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict")
	if code == 0 || !strings.Contains(out, "warning: threads.archived: TH-0389 is open but S-003 is archived; resolve it with flai thread resolve TH-0389") {
		t.Fatalf("unscoped, the run should warn on TH-0389 open on the archived S-003 and fail --strict: exit %d\n%s%s", code, out, errOut)
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "threads.archived") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the threads.archived should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); len(names) != 0 {
		t.Errorf("no issue should be opened, got %v", names)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != thread {
		t.Errorf("TH-0389 should be left as it was: %v\n%s", err, data)
	}
}

// I-0161, ADR-0171: a close-out on a host or in a worktree that cannot read
// a component's props, here a PATH without node, neither notes nor records
// the components.examples note for examples the story did not change, while
// the unscoped check still reports it. When the story changes an examples
// file, the note is inside it and still not recorded. (A sub-project whose
// svelte is not installed gets the same note on its folder; the check
// package's scope tests cover it, as node cannot be faked here.)
func TestCheckRecordIssuesLeavesOutAnUncheckedExamplesNoteOutsideTheStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("PATH", t.TempDir())
	root := overlapFixture(t, "", nil)
	const examples = "design/components/web/src/lib/Card.examples.json"
	manifest := filepath.Join(root, "system-flow.yaml")
	data, err := os.ReadFile(manifest)
	if err != nil {
		t.Fatal(err)
	}
	write := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(manifest, string(data)+"projects:\n  - name: web\n    path: web\n    kind: sveltekit\n")
	write(filepath.Join(root, "web", "src", "lib", "Card.svelte"), "")
	write(filepath.Join(root, filepath.FromSlash(examples)), `{"examples": [{"name": "A", "props": {}}]}`)
	const note = "warning: components.examples: the props of the examples under design/components are not checked against their components': flai reads a component's props with node, which is not on PATH"
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check")
	if code != 0 || !strings.Contains(out, note) {
		t.Fatalf("unscoped, the run should note the examples left unchecked: exit %d\n%s%s", code, out, errOut)
	}
	var res struct {
		Findings []struct {
			Rule    string `json:"rule"`
			Outside bool   `json:"outside"`
		} `json:"findings"`
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "components.examples") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, which changes no examples, the note should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); len(names) != 0 {
		t.Errorf("no issue should be opened, got %v", names)
	}

	// Here the story leaves E-001 alone, so its epic.lags-stories is outside
	// it and recorded; the note is not.
	var own struct {
		Findings []struct {
			Rule    string `json:"rule"`
			Outside bool   `json:"outside"`
		} `json:"findings"`
		Recorded []struct {
			Rule string `json:"rule"`
		} `json:"recorded"`
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: gitScript{changed: examples}}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	if err := json.Unmarshal([]byte(out), &own); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	inside := 0
	for _, f := range own.Findings {
		if f.Rule == "components.examples" && !f.Outside {
			inside++
		}
	}
	recorded := 0
	for _, r := range own.Recorded {
		if r.Rule == "components.examples" {
			recorded++
		}
	}
	if inside != 1 || recorded != 0 {
		t.Errorf("scoped to S-004, which changes %s, the note should be inside it and not recorded: %s", examples, out)
	}
}

// I-0096, ADR-0123: another story in progress leaves a code span ending in a
// space in its narrative, as S-0229 did; a close-out of S-004 neither notes
// its markdown.MD038 nor records it in an issue, while the unscoped check
// still warns on it.
func TestCheckRecordIssuesLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", map[string]string{"S-006": "docs"})
	if err := os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD041: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	narrative := filepath.Join(root, "wip/agents/S-006.md")
	data, err := os.ReadFile(narrative)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(narrative, append(data, "\nThe `rule: ` field.\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, _ := runWithApp(t, &app{cwd: root, runner: git}, "check")
	if !strings.Contains(out, "wip/agents/S-006.md:29: warning: markdown.MD038: ") {
		t.Fatalf("unscoped, the run should report the MD038 on S-006's narrative:\n%s%s", out, errOut)
	}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "markdown.MD038") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the MD038 should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened or bumped, got %v", names)
	}
}

// I-0109, ADR-0125: another story just started, its narrative's Current
// state and Next steps still the template's placeholders, as S-0265's were;
// a close-out of S-004 neither notes its narrative.state nor records it in an
// issue, while the unscoped check still warns on it.
func TestCheckRecordIssuesLeavesOutNarrativeStateOnAnotherOpenStorysNarrative(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := overlapFixture(t, "", map[string]string{"S-006": "docs"})
	narrative := filepath.Join(root, "wip/agents/S-006.md")
	data, err := os.ReadFile(narrative)
	if err != nil {
		t.Fatal(err)
	}
	unwritten := strings.Replace(string(data), "## Current state\ns\n\n## Next steps\n1. n\n", "## Current state\n\n## Next steps\n1.\n", 1)
	if unwritten == string(data) {
		t.Fatalf("S-006's narrative has no Current state and Next steps to unwrite:\n%s", data)
	}
	if err := os.WriteFile(narrative, []byte(unwritten), 0o644); err != nil {
		t.Fatal(err)
	}
	git := gitScript{changed: "wip/kanban/epics/E-001-epic.md"}
	out, errOut, _ := runWithApp(t, &app{cwd: root, runner: git}, "check")
	if !strings.Contains(out, "wip/agents/S-006.md:14: warning: narrative.state: ") {
		t.Fatalf("unscoped, the run should report the narrative.state on S-006's narrative:\n%s%s", out, errOut)
	}
	out, errOut, code := runWithApp(t, &app{cwd: root, runner: git}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if strings.Contains(out, "narrative.state") || res.Outside != 0 || len(res.Recorded) != 0 {
		t.Errorf("scoped to S-004, the narrative.state should be neither noted nor recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("no issue should be opened or bumped, got %v", names)
	}
}

// S-0227: a finding that quotes a code span is written without backticks,
// so the issue it is recorded in passes the markdown lint (MD038).
func TestFindingTextLeavesNoCodeSpanOpen(t *testing.T) {
	got := findingText("MD038/no-space-in-code Spaces inside code span elements\n  [Context: \"`rule: `\"]")
	want := `MD038/no-space-in-code Spaces inside code span elements [Context: "'rule: '"]`
	if got != want {
		t.Errorf("findingText = %q, want %q", got, want)
	}
}

// S-0249: --record-issues needs --story, and a run with nothing outside the
// story writes nothing under design/issues.
func TestCheckRecordIssuesNeedsStoryAndWritesOnlyWhatIsOutside(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := scopeFixture(t)
	out, errOut, code := runIn(t, root, "check", "--strict", "--record-issues")
	if code == 0 || !strings.Contains(errOut, "--record-issues records the findings outside a story: give --story S-nnnn") || strings.Contains(out, "items checked") {
		t.Errorf("without --story it should be refused before checking: %d\n%s%s", code, out, errOut)
	}

	// Without S-002's branch, the only finding is the epic, which the story's
	// branch changes here, so nothing is outside.
	root = t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../internal/metrics/testdata/good")); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runWithApp(t, &app{cwd: root, runner: gitScript{changed: "wip/kanban/epics/E-001-epic.md"}}, "check", "--strict", "--story", "S-004", "--record-issues", "--json")
	var res struct {
		Outside  int               `json:"outside"`
		Findings []json.RawMessage `json:"findings"`
		Recorded []json.RawMessage `json:"recorded"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || code != 0 {
		t.Fatalf("exit %d %v\n%s%s", code, err, out, errOut)
	}
	if res.Outside != 0 || len(res.Findings) == 0 || res.Recorded == nil || len(res.Recorded) != 0 {
		t.Errorf("want findings, none outside, and an empty recorded: %s", out)
	}
	if names := issueFiles(t, root); names != nil {
		t.Errorf("nothing outside should write nothing under design/issues, got %v", names)
	}
}

// Run in a story's worktree, a finding on wip/, which lives in the main
// checkout, is recorded under the path the project names it by, not with ../.
func TestProjectPathNamesWipFromTheMainCheckout(t *testing.T) {
	main := filepath.Join(t.TempDir(), "proj")
	wt := filepath.Join(main, ".flai-cache", "worktrees", "S-0004")
	repo := &workitem.Repo{Root: wt, MainRoot: main}
	for _, c := range []struct{ in, want string }{
		{filepath.Join("..", "..", "..", "wip", "kanban", "board.md"), "wip/kanban/board.md"},
		{filepath.Join("design", "issues", "summary.md"), "design/issues/summary.md"},
		{filepath.Join(main, "wip", "agents", "S-0004.md"), "wip/agents/S-0004.md"},
	} {
		if got := projectPath(repo, c.in); got != c.want {
			t.Errorf("projectPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

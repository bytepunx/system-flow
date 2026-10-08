package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// scopeProject copies the good fixture, whose S-004 is in progress with its
// narrative, its task T-003, and the thread TH-0001 anchored on it, and gives
// the archived S-002 a branch that was never merged: story.unaccepted.
func scopeProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("../metrics/testdata/good")); err != nil {
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

// edit rewrites a fixture file, failing the test when old is not in it.
func edit(t *testing.T, root, rel, old, new string) {
	t.Helper()
	p := filepath.Join(root, rel)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), old) {
		t.Fatalf("%s has no %q", rel, old)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(data), old, new, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runScoped(t *testing.T, root string, changed []string) (*Result, *Result) {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	whole, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	scoped, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := ScopeToStory(scoped, repo, "S-004", changed); err != nil {
		t.Fatal(err)
	}
	return whole, scoped
}

// I-0057: a story's close-out stopped at flai check --strict on
// story.unaccepted for another, archived story whose branch was never
// merged. Scoped to the story in progress, the run passes; unscoped, it
// fails as before.
func TestScopedCheckPassesOverAnotherStorysUnmergedBranch(t *testing.T) {
	whole, scoped := runScoped(t, scopeProject(t), nil)
	var unaccepted *Finding
	for i, f := range scoped.Findings {
		if f.Rule == "story.unaccepted" {
			unaccepted = &scoped.Findings[i]
		}
	}
	if unaccepted == nil || !strings.Contains(unaccepted.Path, "S-002") || !strings.Contains(unaccepted.Message, "story/S-002 was never merged") {
		t.Fatalf("the fixture should give story.unaccepted on S-002: %+v", scoped.Findings)
	}
	if whole.OK(true) || !whole.OK(false) {
		t.Errorf("unscoped, story.unaccepted must fail --strict only: %+v", whole)
	}
	if !unaccepted.Outside || unaccepted.Level != Warning {
		t.Errorf("story.unaccepted on S-002 should be a warning outside S-004: %+v", unaccepted)
	}
	if !scoped.OK(true) {
		t.Errorf("scoped to S-004, the run should pass --strict: %+v", scoped)
	}
	// The other finding is E-001's advisory epic.lags-stories: outside too,
	// and counted once, as outside, not again as advisory.
	if scoped.Outside != 2 || scoped.Advisory != 0 || scoped.Warnings != whole.Warnings || whole.Advisory != 1 {
		t.Errorf("outside %d advisory %d warnings %d, want 2, 0, %d", scoped.Outside, scoped.Advisory, scoped.Warnings, whole.Warnings)
	}
	for _, f := range whole.Findings {
		if f.Outside {
			t.Errorf("unscoped, nothing is outside: %+v", f)
		}
	}
}

// S-0249: a finding on what the story owns still fails --strict when the
// run is scoped to it; a wip.overlap that names it, even on its own file, and
// an error on another story's file do not (ADR-0115).
func TestScopedCheckFailsOnlyOnTheStorysOwnFindings(t *testing.T) {
	const overview = "design/system/overview.md"
	for _, tc := range []struct {
		name    string
		rule    string
		path    string // the finding's path ends with this
		changed []string
		fails   bool
		setup   func(t *testing.T, root string)
	}{
		{name: "the story's item file", rule: "item.heading", path: "S-004-four.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/stories/S-004-four.md", "\n## Notes", "\n")
		}},
		{name: "a task's file", rule: "item.heading", path: "T-003-t3.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/tasks/T-003-t3.md", "\n## Notes", "\n")
		}},
		{name: "the narrative", rule: "narrative.section", path: "S-004.md", fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, "wip/agents/S-004.md", "\n## Open questions\n", "\n")
		}},
		{name: "a thread anchored on the story", rule: "threads.entries", path: "TH-0001-is-four-really-done.md", fails: true, setup: func(t *testing.T, root string) {
			p := filepath.Join(root, "wip/threads/TH-0001-is-four-really-done.md")
			data, _ := os.ReadFile(p)
			head, _, _ := strings.Cut(string(data), "\n### ")
			if err := os.WriteFile(p, []byte(head+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "a conversation of the story", rule: "messages.entries", path: "MS-0001-who-takes-flai-cmd.md", fails: true, setup: func(t *testing.T, root string) {
			addStory(t, root, "S-006", "E-001", "")
			writeConversation(t, root, "S-004", "S-006")
		}},
		{name: "a conversation between two other stories", rule: "messages.entries", path: "MS-0001-who-takes-flai-cmd.md", setup: func(t *testing.T, root string) {
			addStory(t, root, "S-006", "E-001", "")
			addStory(t, root, "S-007", "E-001", "")
			writeConversation(t, root, "S-006", "S-007")
		}},
		{name: "a path in the story's diff", rule: "doc.title", path: overview, changed: []string{overview}, fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "a path in a folder the story adds", rule: "doc.title", path: overview, changed: []string{"design/system/"}, fails: true, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "a path outside the story's diff", rule: "doc.title", path: overview, changed: []string{"flai/main.go"}, setup: func(t *testing.T, root string) {
			edit(t, root, overview, "title: Overview\n", "")
		}},
		{name: "wip.overlap naming the story with another in progress", rule: "wip.overlap", path: "S-004-four.md", setup: func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/stories/S-004-four.md", "tags: []\n", "tags: []\ntouches: [flai]\n")
			addStory(t, root, "S-006", "E-001", "touches: [flai/cmd]\n")
		}},
		{name: "an error on another story's file", rule: "item.parent-missing", path: "S-006-six.md", setup: func(t *testing.T, root string) {
			addStory(t, root, "S-006", "E-009", "")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeProject(t)
			tc.setup(t, root)
			whole, scoped := runScoped(t, root, tc.changed)
			var got *Finding
			for i, f := range scoped.Findings {
				if f.Rule == tc.rule && strings.HasSuffix(filepath.ToSlash(f.Path), tc.path) {
					got = &scoped.Findings[i]
				}
			}
			if got == nil {
				t.Fatalf("no %s on %s: %+v", tc.rule, tc.path, scoped.Findings)
			}
			if whole.OK(true) {
				t.Errorf("unscoped, the run should fail --strict: %+v", whole)
			}
			if got.Outside == tc.fails {
				t.Errorf("outside %v, want %v: %+v", got.Outside, !tc.fails, got)
			}
			if scoped.OK(true) == tc.fails {
				t.Errorf("scoped OK(strict) %v, want %v: %+v", scoped.OK(true), !tc.fails, scoped.Findings)
			}
		})
	}
}

// I-0076, ADR-0115: scoped to a story, a wip.overlap between two other
// stories in progress is left out of the result, and its warning is taken
// back from the counts, so the run neither notes nor records it.
func TestScopedCheckLeavesOutAnOverlapBetweenTwoOtherStories(t *testing.T) {
	root := scopeProject(t)
	edit(t, root, "wip/kanban/board.md", "in-progress: 2\n", "in-progress: 5\n")
	addStory(t, root, "S-006", "E-001", "touches: [flai]\n")
	addStory(t, root, "S-007", "E-001", "touches: [flai/cmd]\n")
	whole, scoped := runScoped(t, root, nil)
	var overlaps []Finding
	for _, f := range whole.Findings {
		if f.Rule == "wip.overlap" {
			overlaps = append(overlaps, f)
		}
	}
	if len(overlaps) != 1 || strings.Join(overlaps[0].Stories, " ") != "S-006 S-007" {
		t.Fatalf("unscoped, the run should report the one overlap of S-006 and S-007: %+v", whole.Findings)
	}
	warnings, outside := 0, 0
	for _, f := range scoped.Findings {
		if f.Rule == "wip.overlap" {
			t.Errorf("scoped to S-004, the overlap of two other stories should be left out: %+v", f)
		}
		if f.Level == Warning {
			warnings++
		}
		if f.Outside {
			outside++
		}
	}
	if len(scoped.Findings) != len(whole.Findings)-1 || scoped.Warnings != whole.Warnings-1 || scoped.Errors != whole.Errors {
		t.Errorf("scoped findings %d warnings %d errors %d, want %d, %d, %d", len(scoped.Findings), scoped.Warnings, scoped.Errors, len(whole.Findings)-1, whole.Warnings-1, whole.Errors)
	}
	if scoped.Warnings != warnings || scoped.Outside != outside {
		t.Errorf("counts should match the findings kept: warnings %d outside %d, findings give %d, %d", scoped.Warnings, scoped.Outside, warnings, outside)
	}
}

// I-0078, ADR-0122: scoped to a story, an item.archive on another story,
// cancelled from backlog and not archived, is left out of the result, and
// its warning is taken back from the counts, so the run neither notes nor
// records it. Unscoped, it is reported as before.
func TestScopedCheckLeavesOutAnItemArchiveOutsideTheStory(t *testing.T) {
	root := scopeProject(t)
	addCancelledStory(t, root, "S-006")
	whole, scoped := runScoped(t, root, nil)
	var archive []Finding
	for _, f := range whole.Findings {
		if f.Rule == "item.archive" {
			archive = append(archive, f)
		}
	}
	if len(archive) != 1 || archive[0].Message != "S-006 is cancelled; run flai archive" {
		t.Fatalf("unscoped, the run should report S-006's item.archive: %+v", whole.Findings)
	}
	warnings, outside := 0, 0
	for _, f := range scoped.Findings {
		if f.Rule == "item.archive" {
			t.Errorf("scoped to S-004, the item.archive of another story should be left out: %+v", f)
		}
		if f.Level == Warning {
			warnings++
		}
		if f.Outside {
			outside++
		}
	}
	if len(scoped.Findings) != len(whole.Findings)-1 || scoped.Warnings != whole.Warnings-1 || scoped.Errors != whole.Errors {
		t.Errorf("scoped findings %d warnings %d errors %d, want %d, %d, %d", len(scoped.Findings), scoped.Warnings, scoped.Errors, len(whole.Findings)-1, whole.Warnings-1, whole.Errors)
	}
	if scoped.Warnings != warnings || scoped.Outside != outside {
		t.Errorf("counts should match the findings kept: warnings %d outside %d, findings give %d, %d", scoped.Warnings, scoped.Outside, warnings, outside)
	}
}

// I-0123, ADR-0133: scoped to a story, a board.wip-limit on the main
// checkout's board, in progress over its limit as at S-0321's and S-0339's
// close-outs or review over it, advisory, is left out of the result, and its
// counts are taken back, advisory included, so the run neither notes nor
// records it. Unscoped, it is reported as before.
func TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory(t *testing.T) {
	for _, tc := range []struct {
		name     string
		message  string
		advisory bool
		setup    func(t *testing.T, root string)
	}{
		{"in progress over its limit", "3 stories in in-progress, limit 2", false, func(t *testing.T, root string) {
			addStory(t, root, "S-006", "E-001", "")
			addStory(t, root, "S-007", "E-001", "")
		}},
		{"review over its limit", "2 stories in review, limit 1", true, func(t *testing.T, root string) {
			edit(t, root, "wip/kanban/board.md", "review: 3\n", "review: 1\n")
			for _, id := range []string{"S-006", "S-007"} {
				addStory(t, root, id, "E-001", "")
				rel := "wip/kanban/stories/" + id + "-six.md"
				edit(t, root, rel, "status: in-progress\n", "status: review\n")
				edit(t, root, rel, "    by: agent\ntags: []\n", "    by: agent\n  - to: review\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\n")
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeProject(t)
			tc.setup(t, root)
			whole, scoped := runScoped(t, root, nil)
			var limits []Finding
			for _, f := range whole.Findings {
				if f.Rule == "board.wip-limit" {
					limits = append(limits, f)
				}
			}
			if len(limits) != 1 || limits[0].Message != tc.message || limits[0].advisory != tc.advisory ||
				!strings.HasSuffix(filepath.ToSlash(limits[0].Path), "wip/kanban/board.md") {
				t.Fatalf("unscoped, the run should report %q on the board: %+v", tc.message, whole.Findings)
			}
			if !tc.advisory && whole.OK(true) {
				t.Errorf("unscoped, in progress over its limit should fail --strict: %+v", whole)
			}
			for _, f := range scoped.Findings {
				if f.Rule == "board.wip-limit" {
					t.Errorf("scoped to S-004, the board.wip-limit should be left out: %+v", f)
				}
			}
			if len(scoped.Findings) != len(whole.Findings)-1 || scoped.Warnings != whole.Warnings-1 || scoped.Errors != whole.Errors {
				t.Errorf("scoped findings %d warnings %d errors %d, want %d, %d, %d", len(scoped.Findings), scoped.Warnings, scoped.Errors, len(whole.Findings)-1, whole.Warnings-1, whole.Errors)
			}
			assertCountsMatch(t, scoped)
			advisory := 0
			for _, f := range scoped.Findings {
				if f.advisory && !f.Outside {
					advisory++
				}
			}
			if scoped.Advisory != advisory {
				t.Errorf("scoped advisory %d, the findings kept give %d", scoped.Advisory, advisory)
			}
			if !scoped.OK(true) {
				t.Errorf("scoped to S-004, the run should pass --strict: %+v", scoped.Findings)
			}
		})
	}
}

// I-0096, ADR-0123: scoped to a story, a markdown finding on the narrative
// of another story in progress, a code span ending in a space as S-0229's
// was, is left out of the result, and its warning is taken back from the
// counts, so the run neither notes nor records it: that story's own
// close-out finds it. Unscoped, it is reported as before.
func TestScopedCheckLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative(t *testing.T) {
	root := scopeProject(t)
	lintMarkdown(t, root)
	addStory(t, root, "S-006", "E-001", "")
	addCodeSpan(t, root, "wip/agents/S-006.md")
	whole, scoped := runScoped(t, root, nil)
	if got := codeSpanFindings(whole); len(got) != 1 || filepath.ToSlash(got[0].Path) != "wip/agents/S-006.md" {
		t.Fatalf("unscoped, the run should report the MD038 on S-006's narrative: %+v", whole.Findings)
	}
	if got := codeSpanFindings(scoped); len(got) != 0 {
		t.Errorf("scoped to S-004, the MD038 on S-006's narrative should be left out: %+v", got)
	}
	if len(scoped.Findings) != len(whole.Findings)-1 || scoped.Warnings != whole.Warnings-1 || scoped.Errors != whole.Errors {
		t.Errorf("scoped findings %d warnings %d errors %d, want %d, %d, %d", len(scoped.Findings), scoped.Warnings, scoped.Errors, len(whole.Findings)-1, whole.Warnings-1, whole.Errors)
	}
	assertCountsMatch(t, scoped)
}

// ADR-0123 §2: a markdown finding outside the story on anything but another
// open story's narrative stays a note: the narrative of a done or cancelled
// story, an archived narrative, a task, or a thread of another story.
func TestScopedCheckNotesEveryOtherMarkdownFindingOutsideTheStory(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root string) string
	}{
		{"done story's narrative", func(t *testing.T, root string) string {
			addStory(t, root, "S-006", "E-001", "")
			edit(t, root, "wip/kanban/stories/S-006-six.md", "status: in-progress\n", "status: done\n")
			return "wip/agents/S-006.md"
		}},
		{"cancelled story's narrative", func(t *testing.T, root string) string {
			addStory(t, root, "S-006", "E-001", "")
			edit(t, root, "wip/kanban/stories/S-006-six.md", "status: in-progress\n", "status: cancelled\n")
			return "wip/agents/S-006.md"
		}},
		{"archived narrative", func(t *testing.T, root string) string {
			if err := os.MkdirAll(filepath.Join(root, "wip/archive/agents"), 0o755); err != nil {
				t.Fatal(err)
			}
			narrative, err := os.ReadFile(filepath.Join(root, "wip/agents/S-004.md"))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, "wip/archive/agents/S-005.md"), []byte(strings.ReplaceAll(string(narrative), "S-004", "S-005")), 0o644); err != nil {
				t.Fatal(err)
			}
			return "wip/archive/agents/S-005.md"
		}},
		{"another story's task", func(*testing.T, string) string {
			return "wip/archive/kanban/tasks/T-004-t4.md"
		}},
		{"another story's thread", func(t *testing.T, root string) string {
			addStory(t, root, "S-006", "E-001", "")
			edit(t, root, "wip/threads/TH-0001-is-four-really-done.md", "  item: S-004\n", "  item: S-006\n")
			return "wip/threads/TH-0001-is-four-really-done.md"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := scopeProject(t)
			lintMarkdown(t, root)
			rel := tc.setup(t, root)
			addCodeSpan(t, root, rel)
			whole, scoped := runScoped(t, root, nil)
			if got := codeSpanFindings(whole); len(got) != 1 || filepath.ToSlash(got[0].Path) != rel {
				t.Fatalf("unscoped, the run should report the MD038 on %s: %+v", rel, whole.Findings)
			}
			if got := codeSpanFindings(scoped); len(got) != 1 || !got[0].Outside {
				t.Errorf("scoped to S-004, the MD038 on %s should be a note outside it: %+v", rel, got)
			}
			assertCountsMatch(t, scoped)
		})
	}
}

// I-0109, ADR-0125: scoped to a story, the narrative.state findings on the
// narrative of another story just started, whose Current state and Next steps
// still hold the template's placeholder as S-0265's did, are left out of the
// result and their warnings taken back from the counts, so the run neither
// notes nor records them: that story's own close-out stops on them. Unscoped,
// they are reported as before, and the story's own stay inside it.
func TestScopedCheckLeavesOutNarrativeStateOnAnotherOpenStorysNarrative(t *testing.T) {
	root := scopeProject(t)
	addStory(t, root, "S-006", "E-001", "")
	unwriteState(t, root, "wip/agents/S-006.md")
	whole, scoped := runScoped(t, root, nil)
	if got := stateFindings(whole); len(got) != 2 || filepath.ToSlash(got[0].Path) != "wip/agents/S-006.md" {
		t.Fatalf("unscoped, the run should report S-006's unwritten Current state and Next steps: %+v", whole.Findings)
	}
	if got := stateFindings(scoped); len(got) != 0 {
		t.Errorf("scoped to S-004, the narrative.state on S-006's narrative should be left out: %+v", got)
	}
	if len(scoped.Findings) != len(whole.Findings)-2 || scoped.Warnings != whole.Warnings-2 || scoped.Errors != whole.Errors {
		t.Errorf("scoped findings %d warnings %d errors %d, want %d, %d, %d", len(scoped.Findings), scoped.Warnings, scoped.Errors, len(whole.Findings)-2, whole.Warnings-2, whole.Errors)
	}
	assertCountsMatch(t, scoped)

	unwriteState(t, root, "wip/agents/S-004.md")
	_, scoped = runScoped(t, root, nil)
	if got := stateFindings(scoped); len(got) != 2 || got[0].Outside || got[1].Outside || filepath.ToSlash(got[0].Path) != "wip/agents/S-004.md" {
		t.Errorf("scoped to S-004, its own unwritten narrative should be inside it: %+v", got)
	}
}

// ADR-0125 §2: another rule's finding on another open story's narrative stays
// a note, as a narrative.section for a section its agent removed.
func TestScopedCheckNotesAnotherNarrativeRuleOnAnotherOpenStorysNarrative(t *testing.T) {
	root := scopeProject(t)
	addStory(t, root, "S-006", "E-001", "")
	edit(t, root, "wip/agents/S-006.md", "## Open questions\n", "")
	_, scoped := runScoped(t, root, nil)
	var got []Finding
	for _, f := range scoped.Findings {
		if f.Rule == "narrative.section" {
			got = append(got, f)
		}
	}
	if len(got) != 1 || !got[0].Outside || filepath.ToSlash(got[0].Path) != "wip/agents/S-006.md" {
		t.Errorf("scoped to S-004, the narrative.section on S-006's narrative should be a note outside it: %+v", got)
	}
	assertCountsMatch(t, scoped)
}

// unwriteState puts back the template's placeholders under a fixture
// narrative's Current state and Next steps, as flai stream open writes them.
func unwriteState(t *testing.T, root, rel string) {
	t.Helper()
	edit(t, root, rel, "## Current state\ns\n\n## Next steps\n1. n\n", "## Current state\n\n## Next steps\n1.\n")
}

// stateFindings is res's narrative.state findings.
func stateFindings(res *Result) []Finding {
	var out []Finding
	for _, f := range res.Findings {
		if f.Rule == "narrative.state" {
			out = append(out, f)
		}
	}
	return out
}

// lintMarkdown gives the project a markdownlint configuration, so the wip
// markdown is linted (S-0179).
func lintMarkdown(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD041: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// addCodeSpan appends to a fixture file a code span ending in a space, which
// the markdown lint reports as MD038 (I-0096).
func addCodeSpan(t *testing.T, root, rel string) {
	t.Helper()
	p := filepath.Join(root, rel)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(data, "\nThe `rule: ` field.\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// codeSpanFindings is res's markdown.MD038 findings.
func codeSpanFindings(res *Result) []Finding {
	var out []Finding
	for _, f := range res.Findings {
		if f.Rule == "markdown.MD038" {
			out = append(out, f)
		}
	}
	return out
}

// assertCountsMatch fails when res's warnings and outside counts are not
// those of the findings it kept.
func assertCountsMatch(t *testing.T, res *Result) {
	t.Helper()
	warnings, outside := 0, 0
	for _, f := range res.Findings {
		if f.Level == Warning {
			warnings++
		}
		if f.Outside {
			outside++
		}
	}
	if res.Warnings != warnings || res.Outside != outside {
		t.Errorf("counts should match the findings kept: warnings %d outside %d, findings give %d, %d", res.Warnings, res.Outside, warnings, outside)
	}
}

// addCancelledStory writes a story under E-001 that was cancelled from
// backlog and never archived, as S-0250 was (I-0078), and lists it in the
// epic.
func addCancelledStory(t *testing.T, root, id string) {
	t.Helper()
	story := "---\nid: " + id + "\ntype: story\nnature: feature\ntitle: Six\nstatus: cancelled\nparent: E-001" +
		"\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-26T10:00:00Z\ntransitions:\n  - to: cancelled\n    at: 2026-08-26T10:00:00Z\n    by: alex\ntags: []\n" +
		"---\n\n# " + id + " Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/stories", id+"-six.md"), []byte(story), 0o644); err != nil {
		t.Fatal(err)
	}
	edit(t, root, "wip/kanban/epics/E-001-epic.md", "- S-005 Five\n", "- S-005 Five\n- "+id+" Six\n")
}

// addStory writes an in-progress story under parent, with front matter
// extra, and its narrative.
func addStory(t *testing.T, root, id, parent, extra string) {
	t.Helper()
	story := "---\nid: " + id + "\ntype: story\nnature: feature\ntitle: Six\nstatus: in-progress\nparent: " + parent +
		"\nowner: agent\ncreated: 2026-08-25T09:00:00Z\nupdated: 2026-08-31T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-25T10:00:00Z\n    by: agent\n  - to: in-progress\n    at: 2026-08-31T10:00:00Z\n    by: agent\ntags: []\n" + extra +
		"---\n\n# " + id + " Six\n\n## Goal\ng\n\n## Acceptance criteria\n- [ ] works\n\n## Tasks\n\n## Notes\n"
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/stories", id+"-six.md"), []byte(story), 0o644); err != nil {
		t.Fatal(err)
	}
	narrative, err := os.ReadFile(filepath.Join(root, "wip/agents/S-004.md"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "wip/agents", id+".md"), []byte(strings.ReplaceAll(string(narrative), "S-004", id)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestScopeToStoryNamesAStory(t *testing.T) {
	repo, err := workitem.Open(scopeProject(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"S-099", "T-003", "nonsense"} {
		if err := ScopeToStory(&Result{}, repo, id, nil); err == nil || !strings.Contains(err.Error(), "scope the check to") {
			t.Errorf("%s: want an error naming the scope, got %v", id, err)
		}
	}
}

// writeConversation writes a conversation between two stories with no dated
// entries, which the messages rule warns on (ADR-0120).
func writeConversation(t *testing.T, root, from, to string) {
	t.Helper()
	doc := "---\nid: MS-0001\ntitle: Who takes flai/cmd\nfrom: " + from + "\nto: " + to + "\nstatus: open\nparticipants: [agent-" + from + "]\ncreated: 2026-10-07T10:00:00Z\nupdated: 2026-10-07T10:00:00Z\n---\n\n# MS-0001 Who takes flai/cmd\n\n## Entries\n"
	p := filepath.Join(root, "wip/messages/MS-0001-who-takes-flai-cmd.md")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
}

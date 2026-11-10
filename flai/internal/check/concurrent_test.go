package check

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// oldOrder is every rule Run calls, by name, in the order Run called them one
// after another before S-0601. It is written out here rather than read from
// sequence, so that a rule dropped from sequence, or moved in it, fails the
// test.
var oldOrder = []struct {
	name string
	run  func(*checker)
}{
	{"layout", (*checker).layout},
	{"instructions", (*checker).instructions},
	{"workItems", (*checker).workItems},
	{"narratives", (*checker).narratives},
	{"activities", (*checker).activities},
	{"overlap", (*checker).overlap},
	{"componentTag", (*checker).componentTag},
	{"after", (*checker).after},
	{"touchesFolder", (*checker).touchesFolder},
	{"afterHub", (*checker).afterHub},
	{"afterReason", (*checker).afterReason},
	{"unaccepted", (*checker).unaccepted},
	{"epicLags", (*checker).epicLags},
	{"board", (*checker).board},
	{"documentation", (*checker).documentation},
	{"adrIndex", (*checker).adrIndex},
	{"conventions", (*checker).conventions},
	{"roles", (*checker).roles},
	{"formulas", (*checker).formulas},
	{"issues", (*checker).issues},
	{"experiments", (*checker).experiments},
	{"analysis", (*checker).analysis},
	{"digest", (*checker).digest},
	{"designTokens", (*checker).designTokens},
	{"componentExamples", (*checker).componentExamples},
	{"threads", (*checker).threads},
	{"messages", (*checker).messages},
	{"markdown", (*checker).markdown},
	{"conflictMarkers", (*checker).conflictMarkers},
}

// sequential runs every rule one after another on one checker, in oldOrder,
// as Run did before S-0601, and returns its result and the number of findings
// each rule added, by name.
func sequential(t *testing.T, repo *workitem.Repo) (*Result, map[string]int) {
	t.Helper()
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	c := &checker{repo: repo, now: now.UTC(), items: items, byID: map[string]*workitem.Item{}, res: &Result{Items: len(items)}}
	added := map[string]int{}
	for _, r := range oldOrder {
		n := len(c.res.Findings)
		r.run(c)
		added[r.name] = len(c.res.Findings) - n
	}
	sortFindings(c.res.Findings)
	return c.res, added
}

// everyRuleProject is testdata/bad copied to a temporary folder, with files
// added so that every rule Run calls finds something: two stories in
// progress whose claims overlap, one touching two sub-projects untagged and
// one keeping a folder, three open stories naming S-001 in after: with no
// reason and one naming a story that does not exist, an epic in backlog with
// a story in progress, a done story never accepted, an ADR index row that
// links no ADR file, broken role, formula, experiment, analysis, digest, activity, and
// token files, a component's examples with node missing, a thread with a
// missing anchor and a duplicate heading, a conversation with a bad status,
// and a merge conflict marker.
func everyRuleProject(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("testdata/bad")); err != nil {
		t.Fatal(err)
	}
	w := func(rel, content string) {
		t.Helper()
		writeFile(t, filepath.Join(root, filepath.FromSlash(rel)), content)
	}
	manifest, err := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	w("system-flow.yaml", string(manifest)+webAndAPI)
	w("web/src/lib/Card.svelte", "")
	w("web/src/lib/Panel.svelte", "")
	w("api/main.go", "package main\n")
	story := func(id, parent, status, touches, after string) string {
		moves := ""
		for _, to := range []string{workitem.InProgress, workitem.Review, workitem.Done} {
			if boardRank[to] > boardRank[status] {
				break
			}
			moves += "  - to: " + to + "\n    at: 2026-08-02T10:00:00Z\n    by: a\n"
		}
		return "---\nid: " + id + "\ntype: story\nnature: feature\ntitle: S\nstatus: " + status + "\nparent: " + parent + "\nowner: agent\ncreated: 2026-08-01T09:00:00Z\nupdated: 2026-08-02T10:00:00Z\ntransitions:\n  - to: ready\n    at: 2026-08-01T10:00:00Z\n    by: a\n" + moves +
			"tags: []\ntouches: [" + touches + "]\nafter: [" + after + "]\n---\n\n# " + id + " S\n\n## Goal\ng\n\n## Acceptance criteria\n\n- [x] c\n\n## Tasks\n\n## Notes\n"
	}
	w("wip/kanban/epics/E-002-e.md", "---\nid: E-002\ntype: epic\nnature: feature\ntitle: E\nstatus: backlog\nowner: agent\ncreated: 2026-08-01T09:00:00Z\nupdated: 2026-08-01T09:00:00Z\ntransitions: []\ntags: []\n---\n\n# E-002 E\n\n## Outcome\no\n\n## Stories\n\n## Notes\n")
	w("wip/kanban/stories/S-005-s.md", story("S-005", "E-002", workitem.InProgress, "web/src/lib/Card.svelte, api/main.go", "S-001"))
	w("wip/kanban/stories/S-006-s.md", story("S-006", "E-001", workitem.InProgress, "web/src/lib/", "S-001"))
	w("wip/kanban/stories/S-008-s.md", story("S-008", "E-001", workitem.Ready, "web/src/lib/Panel.svelte", "S-001, S-099"))
	w("wip/kanban/stories/S-009-s.md", story("S-009", "E-001", workitem.Done, "", ""))
	w("design/adrs/README.md", "# Architecture Decision Records\n\n| ADR | Title | Status |\n|-----|-------|--------|\n| [0007](0007-gone.md) | Gone | accepted |\n")
	w("design/workflows/roles/broken.md", "---\nname: [broken\n---\n")
	w("design/workflows/formulas/broken.md", "---\nname: [broken\n---\n")
	w("design/experiments/broken.md", "no front matter\n")
	w("design/analysis/broken.md", "no front matter\n")
	w("design/digest/broken.md", "no front matter\n")
	w("wip/agents/planner.md", "---\nkind: [broken\n---\n")
	w("design/system/design-system.tokens.json", "{")
	w("design/components/"+cardFile, `{"examples":[]}`)
	w(".markdownlint.yaml", "default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\n")
	w("wip/threads/TH-0001-t.md", "---\nid: TH-0001\ntitle: T\nanchor:\n  path: docs/missing.md\nstatus: open\nparticipants: [alex]\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\n---\n\n# TH-0001 T\n\n## Entries\n\n### 2026-09-01T12:00:00Z alex\nhi\n\n### 2026-09-01T12:00:00Z alex\nagain\n")
	w("wip/messages/MS-0001-t.md", "---\nid: MS-0001\ntitle: T\nfrom: S-005\nto: S-006\nstatus: pending\nparticipants: [alex]\ncreated: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\n---\n\n# MS-0001 T\n\n## Entries\n\n### 2026-09-01T12:00:00Z alex S-005\nhi\n")
	w("docs/users/merge.md", "---\ntitle: Merge\nupdated: 2026-09-01\n---\n\n"+strings.Repeat("<", 7)+" HEAD\na\n"+strings.Repeat("=", 7)+"\nb\n"+strings.Repeat(">", 7)+" topic\n")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// S-0601: Run, which runs its rules at once, gives the findings, their
// levels, marks, and order, and the counts that running them one after
// another in the old order gives, over a project with a finding of every
// rule. Run it under -race: the rules share the checker's items, byID, and
// vocab.
func TestRunAtOnceMatchesTheRulesOneAfterAnother(t *testing.T) {
	withNode(t, &fakeNode{missing: true})
	repo := everyRuleProject(t)
	want, added := sequential(t, repo)
	var none []string
	for _, r := range oldOrder {
		if added[r.name] == 0 {
			none = append(none, r.name)
		}
	}
	if len(none) > 0 {
		t.Errorf("the fixture gives these rules no finding, so the test does not cover them: %v", none)
	}
	if len(sequence) != len(oldOrder) {
		t.Fatalf("Run runs %d rules, the old order %d", len(sequence), len(oldOrder))
	}
	for i := range sequence {
		name, run := oldOrder[i].name, oldOrder[i].run
		if name == "workItems" {
			// Run runs workItems as indexItems alone, then oneItem in chunks (S-0677)
			run = (*checker).indexItems
			if !sequence[i].first {
				t.Errorf("rule %d of Run, indexItems, is not marked first", i)
			}
		}
		if reflect.ValueOf(sequence[i].run).Pointer() != reflect.ValueOf(run).Pointer() {
			t.Errorf("rule %d of Run is not %s", i, name)
		}
	}
	// S-0677: chunks of one item each, so that the per-item findings come from
	// more than one chunk and merge in item order.
	itemsWithFindings := map[string]bool{}
	for _, f := range want.Findings {
		if strings.HasPrefix(f.Rule, "item.") || strings.HasPrefix(f.Rule, "story.") {
			itemsWithFindings[f.Path] = true
		}
	}
	if len(itemsWithFindings) < 2 {
		t.Fatalf("the fixture gives per-item findings on %d files; the chunks need two or more", len(itemsWithFindings))
	}
	for _, size := range []int{1, 2, itemChunk} {
		setItemChunk(t, size)
		for range 5 {
			got, err := Run(repo, now)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Run with chunks of %d items differs from the rules one after another:\n got %+v\nwant %+v", size, got, want)
			}
		}
	}
}

// setItemChunk sets itemChunk to n for the test.
func setItemChunk(t *testing.T, n int) {
	t.Helper()
	old := itemChunk
	itemChunk = n
	t.Cleanup(func() { itemChunk = old })
}

// S-0677: indexItems' children index gives each item's children as
// workitem.Children does, in the same order, over a project with stories,
// tasks, and a duplicate ID.
func TestChildrenIndexMatchesChildren(t *testing.T) {
	repo := everyRuleProject(t)
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	c := &checker{repo: repo, now: now.UTC(), items: items, byID: map[string]*workitem.Item{}, res: &Result{}}
	c.indexItems()
	parents := 0
	for _, it := range items {
		want := workitem.Children(items, it.ID)
		if len(want) > 0 {
			parents++
		}
		if got := c.childrenOf(it.ID); !slices.Equal(got, want) {
			t.Errorf("children of %s: got %d items, want %d, in workitem.Children's order", it.ID, len(got), len(want))
		}
	}
	if parents < 2 {
		t.Errorf("the fixture has %d items with children; the test needs two or more", parents)
	}
	if got, want := (&checker{items: items}).childrenOf("E-001"), workitem.Children(items, "E-001"); !slices.Equal(got, want) {
		t.Errorf("childrenOf with no index differs from workitem.Children")
	}
}

func TestChunksCutItemsInOrder(t *testing.T) {
	var items []*workitem.Item
	for i := range 5 {
		items = append(items, &workitem.Item{ID: string(rune('a' + i))})
	}
	for _, n := range []int{1, 2, 5, 9} {
		var joined []*workitem.Item
		cs := chunks(items, n)
		for _, ch := range cs {
			if len(ch) == 0 || len(ch) > n {
				t.Errorf("chunks of %d: a chunk of %d", n, len(ch))
			}
			joined = append(joined, ch...)
		}
		if !slices.Equal(joined, items) {
			t.Errorf("chunks of %d do not join back to the items in order", n)
		}
		if want := (len(items) + n - 1) / n; len(cs) != want {
			t.Errorf("chunks of %d: %d chunks, want %d", n, len(cs), want)
		}
	}
	if len(chunks(nil, 3)) != 0 {
		t.Error("chunks of no items should be none")
	}
}

package check

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// writeProject is scopeProject with the markdown lint on and a finding of
// each rule RunItemWrite keeps already in place, none of them the write's:
// story.unaccepted on S-002, epic.lags-stories on E-001, board.order, a
// wip.overlap between S-004 and a second story in progress, S-006, whose
// narrative's state is unwritten, and a markdown finding on S-001's file.
func writeProject(t *testing.T) string {
	t.Helper()
	root := scopeProject(t)
	lintMarkdown(t, root)
	addStory(t, root, "S-006", "E-001", "touches: [flai/x.go]\n")
	edit(t, root, "wip/kanban/epics/E-001-epic.md", "- S-005 Five\n", "- S-005 Five\n- S-006 Six\n")
	edit(t, root, "wip/kanban/stories/S-004-four.md", "tags: []\n", "tags: []\ntouches: [flai/x.go]\n")
	unwriteState(t, root, "wip/agents/S-006.md")
	edit(t, root, "wip/kanban/board.md", "order: []", "order: [S-001]")
	addCodeSpan(t, root, "wip/archive/kanban/stories/S-001-one.md")
	return root
}

// addConflict appends a merge conflict to a fixture file.
func addConflict(t *testing.T, root, rel string) {
	t.Helper()
	p := filepath.Join(root, rel)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(data, "\n<<<<<<< ours\na\n=======\nb\n>>>>>>> theirs\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
}

func findingKey(f Finding) string { return f.Rule + "\x00" + f.Path + "\x00" + f.Message }

func findingKeys(fs []Finding) []string {
	out := make([]string, 0, len(fs))
	for _, f := range fs {
		out = append(out, findingKey(f))
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// newInside is what a checked write refuses on or warns of: the findings of
// after that before did not have, inside the written item's story once after
// is scoped to it (ScopeToItem), by rule, path, and message.
func newInside(t *testing.T, repo *workitem.Repo, before, after *Result, id string, written []string) []string {
	t.Helper()
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	known := map[string]bool{}
	for _, f := range before.Findings {
		known[findingKey(f)] = true
	}
	outside, err := ScopeToItem(after, repo, it, written)
	if err != nil {
		t.Fatal(err)
	}
	var out []Finding
	for _, f := range after.Findings {
		if !known[findingKey(f)] && !outside(f) {
			out = append(out, f)
		}
	}
	return findingKeys(out)
}

// graphRules are the rules RunItemWrite runs over every item, as Run does.
var graphRules = []string{"wip.overlap", "story.component-tag", "story.after", "task.after", "touches.folder", "wip.after-hub", "wip.after-reason", "story.unaccepted", "epic.lags-stories", "board.front-matter", "board.order", "board.wip-limit"}

func byRules(res *Result, rules []string) []string {
	var out []Finding
	for _, f := range res.Findings {
		if slices.Contains(rules, f.Rule) {
			out = append(out, f)
		}
	}
	return findingKeys(out)
}

func rulesOf(ks []string) []string {
	var out []string
	for _, k := range ks {
		out = append(out, strings.SplitN(k, "\x00", 2)[0])
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// S-0598: a checked write's run, before and after the write, brings the same
// new findings inside the written item's story as two whole runs do, for an
// edit of a story, of a task, and of an epic, and for a task's creation; it
// reports what it reports as a whole run does, finds the rules over every
// item as a whole run does, and reads no other story's narrative and lints
// no file the write did not write.
func TestItemWriteRunBringsWhatAWholeRunBrings(t *testing.T) {
	const (
		story = "wip/kanban/stories/S-004-four.md"
		task  = "wip/kanban/tasks/T-003-t3.md"
		epic  = "wip/kanban/epics/E-001-epic.md"
		newT  = "wip/kanban/tasks/T-005-t5.md"
		narr  = "wip/agents/S-004.md"
	)
	cases := []struct {
		name          string
		id            string
		before, after Write
		write         func(t *testing.T, root string)
		want          []string // rules the write brings inside the story
	}{
		{
			name:   "story edit",
			id:     "S-004",
			before: Write{Items: []string{"S-004"}, Paths: []string{story, narr}},
			after:  Write{Items: []string{"S-004"}, Paths: []string{story, narr}},
			write: func(t *testing.T, root string) {
				edit(t, root, story, "touches: [flai/x.go]\n", "touches: [flai/x.go]\nafter: [S-099]\n")
				edit(t, root, story, "- T-003 T3\n", "")
				addCodeSpan(t, root, story)
				addConflict(t, root, story)
				unwriteState(t, root, narr)
			},
			want: []string{"item.parent-list", "markdown.MD038", "markdown.conflict-marker", "narrative.state", "story.after"},
		},
		{
			name:   "task edit",
			id:     "T-003",
			before: Write{Items: []string{"T-003"}, Paths: []string{task}},
			after:  Write{Items: []string{"T-003"}, Paths: []string{task}},
			write: func(t *testing.T, root string) {
				edit(t, root, task, "stream: S-004\n", "stream: S-999\nafter: [T-099]\n")
				addCodeSpan(t, root, task)
				addConflict(t, root, task)
			},
			want: []string{"item.stream", "markdown.MD038", "markdown.conflict-marker", "task.after"},
		},
		{
			name:   "epic edit",
			id:     "E-001",
			before: Write{Items: []string{"E-001"}, Paths: []string{epic}},
			after:  Write{Items: []string{"E-001"}, Paths: []string{epic}},
			write: func(t *testing.T, root string) {
				edit(t, root, epic, "- S-004 Four\n", "")
				addCodeSpan(t, root, epic)
				addConflict(t, root, epic)
			},
			want: []string{"item.parent-list", "markdown.MD038", "markdown.conflict-marker"},
		},
		{
			name:   "task creation",
			id:     "T-005",
			before: Write{Items: []string{"S-004"}, Paths: []string{newT, story}},
			after:  Write{Items: []string{"T-005", "S-004"}, Paths: []string{newT, story}},
			write: func(t *testing.T, root string) {
				data, err := os.ReadFile(filepath.Join(root, task))
				if err != nil {
					t.Fatal(err)
				}
				text := strings.NewReplacer("T-003", "T-005", "T3", "T5", "stream: S-004\n", "stream: S-004\nafter: [T-099]\n").Replace(string(data))
				if err := os.WriteFile(filepath.Join(root, newT), []byte(text), 0o644); err != nil {
					t.Fatal(err)
				}
				addCodeSpan(t, root, newT)
			},
			want: []string{"item.parent-list", "markdown.MD038", "task.after"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeProject(t)
			repo, err := workitem.Open(root)
			if err != nil {
				t.Fatal(err)
			}
			wholeBefore, err := Run(repo, now)
			if err != nil {
				t.Fatal(err)
			}
			writeBefore, err := RunItemWrite(repo, now, tc.before)
			if err != nil {
				t.Fatal(err)
			}
			tc.write(t, root)
			wholeAfter, err := Run(repo, now)
			if err != nil {
				t.Fatal(err)
			}
			writeAfter, err := RunItemWrite(repo, now, tc.after)
			if err != nil {
				t.Fatal(err)
			}
			if writeAfter.Items != wholeAfter.Items {
				t.Errorf("items: %d, want %d as a whole run lists", writeAfter.Items, wholeAfter.Items)
			}
			// what a whole run reports, each rule the write's run keeps
			whole := map[string]bool{}
			for _, k := range findingKeys(wholeAfter.Findings) {
				whole[k] = true
			}
			for _, k := range findingKeys(writeAfter.Findings) {
				if !whole[k] {
					t.Errorf("the write's run reports %q, which a whole run does not", k)
				}
			}
			if got, want := byRules(writeAfter, graphRules), byRules(wholeAfter, graphRules); !slices.Equal(got, want) {
				t.Errorf("rules over every item:\n got %q\nwant %q", got, want)
			}
			for _, rule := range []string{"story.unaccepted", "epic.lags-stories", "board.order", "wip.overlap"} {
				if !slices.Contains(rulesOf(byRules(writeAfter, graphRules)), rule) {
					t.Errorf("the fixture's %s is not reported", rule)
				}
			}
			// what the write brings inside its story, by two runs of each
			wantNew := newInside(t, repo, wholeBefore, wholeAfter, tc.id, tc.after.Paths)
			gotNew := newInside(t, repo, writeBefore, writeAfter, tc.id, tc.after.Paths)
			if !slices.Equal(gotNew, wantNew) {
				t.Errorf("new inside the story:\n got %q\nwant %q", gotNew, wantNew)
			}
			if got := rulesOf(wantNew); !slices.Equal(got, tc.want) {
				t.Errorf("the write brings %v inside the story, want %v", got, tc.want)
			}
			// the other story's narrative is not read, nor another wip file linted
			for _, f := range writeAfter.Findings {
				if f.Path == filepath.Join(root, "wip/agents/S-006.md") {
					t.Errorf("the write's run read S-006's narrative: %+v", f)
				}
				if strings.HasPrefix(f.Rule, "markdown.") && f.Path == filepath.Join(root, "wip/archive/kanban/stories/S-001-one.md") {
					t.Errorf("the write's run linted a file it did not write: %+v", f)
				}
			}
		})
	}
}

// The index is checked only when the write wrote it, as an edit does.
func TestItemWriteRunChecksTheIndexOnlyWhenWritten(t *testing.T) {
	root := writeProject(t)
	edit(t, root, "wip/agents/index.md", "| [S-004](S-004.md) | Four | in-progress | bot | 2026-08-31T10:00:00Z |\n", "")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	index := func(res *Result) int { return len(ruleFindings(res, "narrative.index")) }
	whole, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	without, err := RunItemWrite(repo, now, Write{Items: []string{"S-004"}, Paths: []string{"wip/kanban/stories/S-004-four.md"}})
	if err != nil {
		t.Fatal(err)
	}
	with, err := RunItemWrite(repo, now, Write{Items: []string{"S-004"}, Paths: []string{"wip/kanban/stories/S-004-four.md", "wip/agents/index.md"}})
	if err != nil {
		t.Fatal(err)
	}
	if index(whole) == 0 || index(with) != index(whole) || index(without) != 0 {
		t.Errorf("narrative.index: whole %d, written %d, not written %d; want the whole run's when written, none otherwise", index(whole), index(with), index(without))
	}
}

// failNode fails the test if the check starts node.
type failNode struct{ t *testing.T }

func (f failNode) LookPath(name string) (string, error) {
	f.t.Errorf("the write's run looked for %s", name)
	return "", os.ErrNotExist
}

func (f failNode) Run(context.Context, string, string, string, ...string) ([]byte, []byte, error) {
	f.t.Error("the write's run started node")
	return nil, nil, os.ErrNotExist
}

// A write's run starts no node, where a whole run of the same project does
// to read its components' props (ADR-0164).
func TestItemWriteRunStartsNoNode(t *testing.T) {
	f := &fakeNode{props: map[string]string{"src/lib/Card.svelte": cardProps}}
	withNode(t, f)
	repo := examplesProject(t, webAndAPI, map[string]string{
		cardFile: "{\n  \"examples\": [\n    {\n      \"name\": \"Open\",\n      \"props\": { \"title\": \"A\" }\n    }\n  ]\n}\n",
	})
	if _, err := Run(repo, now); err != nil {
		t.Fatal(err)
	}
	if len(f.asked) != 1 {
		t.Fatalf("a whole run asked node %d times, want once", len(f.asked))
	}
	nodeRunner = failNode{t}
	if _, err := RunItemWrite(repo, now, Write{Paths: []string{cardPath}}); err != nil {
		t.Fatal(err)
	}
}

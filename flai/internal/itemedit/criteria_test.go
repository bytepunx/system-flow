package itemedit

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const storyBody = `## Goal

Tick what was verified.

## Acceptance criteria

- [ ] flai criteria list numbers the boxes
- [x] flai criteria tick ticks them
- [ ] the dashboard ticks them too

## Tasks

- [ ] T-0001 A task
- [x] T-0002 Another

## Notes

- [ ] a box in the notes
`

// S-0282: the acceptance criteria are the boxes under their heading, nested
// ones included, numbered in document order; boxes elsewhere are not theirs.
func TestCriteriaListsTheBoxesUnderAcceptanceCriteria(t *testing.T) {
	for _, c := range []struct {
		name string
		body string
		want []Criterion
	}{
		{"a typical story", storyBody, []Criterion{
			{N: 1, Text: "flai criteria list numbers the boxes"},
			{N: 2, Text: "flai criteria tick ticks them", Ticked: true},
			{N: 3, Text: "the dashboard ticks them too"},
		}},
		{"nested boxes and a capital X", "## Acceptance criteria  \n\n- [ ] outer\n  - [X] inner\n\tnot a box\n- [x] trailing spaces  \n- [x]  two spaces is not a box\n- [ ]\n### Detail\n- [ ] under a subheading\n", []Criterion{
			{N: 1, Text: "outer"},
			{N: 2, Text: "inner", Ticked: true},
			{N: 3, Text: "trailing spaces", Ticked: true},
			{N: 4, Text: "under a subheading"},
		}},
		{"no section", "## Goal\n\n- [ ] not a criterion\n", []Criterion{}},
		{"a section without a box", "## Acceptance criteria\n\nSomething true.\n\n## Notes\n- [ ] note\n", []Criterion{}},
		{"the last section", "## Notes\n\n- [x] note\n\n## Acceptance criteria\n- [x] last", []Criterion{{N: 1, Text: "last", Ticked: true}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := Criteria(c.body)
			if got == nil || !reflect.DeepEqual(got, c.want) {
				t.Errorf("Criteria = %#v, want %#v", got, c.want)
			}
		})
	}
}

// S-0282: a tick or an untick changes the marks it names and no other byte.
func TestTickSetsTheNamedBoxesAndNothingElse(t *testing.T) {
	for _, c := range []struct {
		name         string
		body         string
		tick, untick []int
		want         string
	}{
		{"tick", storyBody, []int{1, 3}, nil,
			strings.Replace(strings.Replace(storyBody, "- [ ] flai criteria list", "- [x] flai criteria list", 1), "- [ ] the dashboard", "- [x] the dashboard", 1)},
		{"untick", storyBody, nil, []int{2},
			strings.Replace(storyBody, "- [x] flai criteria tick", "- [ ] flai criteria tick", 1)},
		{"tick and untick at once, a number twice", storyBody, []int{3, 3}, []int{2},
			strings.Replace(strings.Replace(storyBody, "- [x] flai criteria tick", "- [ ] flai criteria tick", 1), "- [ ] the dashboard", "- [x] the dashboard", 1)},
		{"ticking a ticked box leaves the body as it was", storyBody, []int{2}, nil, storyBody},
		{"a capital X ticked again stays", "## Acceptance criteria\n- [X] done\n", []int{1}, nil, "## Acceptance criteria\n- [X] done\n"},
		{"a capital X is unticked", "## Acceptance criteria\n- [X] done\n", nil, []int{1}, "## Acceptance criteria\n- [ ] done\n"},
		{"a nested box", "## Acceptance criteria\n- [ ] outer\n  - [ ] inner\n", []int{2}, nil, "## Acceptance criteria\n- [ ] outer\n  - [x] inner\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := Tick(c.body, c.tick, c.untick)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.want {
				t.Errorf("Tick =\n%s\nwant\n%s", got, c.want)
			}
		})
	}
}

// S-0282: what cannot be ticked is refused as the caller's mistake, with
// what to do instead.
func TestTickRefusesWhatCannotBeTicked(t *testing.T) {
	for _, c := range []struct {
		name         string
		body         string
		tick, untick []int
		says         string
	}{
		{"no criteria", "## Goal\n\n- [ ] not a criterion\n", []int{1}, nil, "no acceptance criteria with a checkbox"},
		{"nothing to tick or untick", storyBody, nil, nil, "flai criteria list"},
		{"zero", storyBody, []int{0}, nil, "there are 3, numbered 1 to 3 as flai criteria list"},
		{"past the last", storyBody, nil, []int{4}, "no acceptance criterion 4: there are 3"},
		{"in both lists", storyBody, []int{1}, []int{1}, "both to tick and to untick"},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := Tick(c.body, c.tick, c.untick)
			var inv *InvalidError
			if !errors.As(err, &inv) {
				t.Fatalf("want an InvalidError, got %q, %v", got, err)
			}
			if !strings.Contains(err.Error(), c.says) {
				t.Errorf("message %q does not say %q", err, c.says)
			}
		})
	}
}

// raceRepo is a project that lints its markdown, with an epic and under it
// story A in progress, whose criterion the tests tick, with its task, story B
// in progress, each with its narrative written, and story C ready to start.
func raceRepo(t *testing.T) (repo *workitem.Repo, a, task, b, c *workitem.Item) {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD041: false\n"), 0o644)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	create := func(typ, title, parent string) *workitem.Item {
		t.Helper()
		o := workitem.NewOptions{Type: typ, Title: title, Parent: parent, Nature: "feature", Owner: "alex", Now: at}
		if typ == workitem.Story {
			o.Body = "## Goal\n\nDo it.\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"
		}
		it, err := repo.Create(o)
		if err != nil {
			t.Fatal(err)
		}
		return it
	}
	epic := create(workitem.Epic, "Epic", "")
	a = create(workitem.Story, "Story A", epic.ID)
	task = create(workitem.Task, "Task of A", a.ID)
	b = create(workitem.Story, "Story B", epic.ID)
	c = create(workitem.Story, "Story C", epic.ID)
	for _, move := range []struct {
		id string
		to []string
	}{{a.ID, []string{workitem.Ready, workitem.InProgress}}, {b.ID, []string{workitem.Ready, workitem.InProgress}}, {c.ID, []string{workitem.Ready}}} {
		for _, st := range move.to {
			s, _ := repo.Get(move.id)
			if _, err := repo.Transition(s, st, "alex", "", at); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, s := range []*workitem.Item{a, b} {
		s, _ = repo.Get(s.ID)
		if _, err := repo.OpenStream(s, workitem.StreamOptions{Agent: "claude", Now: at}); err != nil {
			t.Fatal(err)
		}
		if _, err := repo.SetStreamState(s.ID, "Working.", "- Go on.", workitem.StreamOptions{Agent: "claude", Now: at}); err != nil {
			t.Fatal(err)
		}
	}
	return repo, a, task, b, c
}

// between has another writer run write between an edit's two runs of the
// check, as another agent does in the main checkout (I-0182).
func between(t *testing.T, write func(repo *workitem.Repo)) {
	t.Helper()
	calls := 0
	runCheck = func(repo *workitem.Repo, now time.Time) (*check.Result, error) {
		if calls++; calls == 2 {
			write(repo)
		}
		return check.Run(repo, now)
	}
	t.Cleanup(func() { runCheck = check.Run })
}

// tick ticks story's first criterion through Apply.
func tick(t *testing.T, repo *workitem.Repo, story string) (*Result, error) {
	t.Helper()
	v, err := Show(repo, story)
	if err != nil {
		t.Fatal(err)
	}
	body, err := Tick(v.Body, []int{1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return Apply(repo, nil, story, Change{Body: &body}, Options{Hash: v.Hash, By: "alex", NoCommit: true, Now: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)})
}

// blankLines gives the file at path three blank lines in a row, a
// markdown.MD012 finding.
func blankLines(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, []byte("\n\n\nMore.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
}

// findings is the findings of a check of repo on path with rule.
func findings(t *testing.T, repo *workitem.Repo, rule, path string) []check.Finding {
	t.Helper()
	res, err := check.Run(repo, time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	var out []check.Finding
	for _, f := range res.Findings {
		if strings.HasPrefix(f.Rule, rule) && filepath.Join(repo.Root, f.Path) == filepath.Clean(path) {
			out = append(out, f)
		}
	}
	return out
}

// I-0182, ADR-0197: a story another agent starts between a tick's two runs
// of the check brings narrative.state findings on its narrative, which a
// check scoped to the ticked story leaves out: the tick stands, with no
// warning for them.
func TestATickStandsWhenAnotherStoryStartsBetweenTheChecks(t *testing.T) {
	repo, a, _, _, c := raceRepo(t)
	between(t, func(repo *workitem.Repo) {
		s, _ := repo.Get(c.ID)
		if _, err := repo.Transition(s, workitem.InProgress, "bob", "", time.Date(2026, 10, 1, 9, 30, 0, 0, time.UTC)); err != nil {
			t.Fatal(err)
		}
		s, _ = repo.Get(c.ID)
		if _, err := repo.OpenStream(s, workitem.StreamOptions{Agent: "bob", Now: time.Date(2026, 10, 1, 9, 30, 0, 0, time.UTC)}); err != nil {
			t.Fatal(err)
		}
		// as flai story start does once the narrative is open
		items, err := repo.List(false)
		if err != nil {
			t.Fatal(err)
		}
		if err := repo.WriteIndex(items, time.Date(2026, 10, 1, 9, 30, 0, 0, time.UTC)); err != nil {
			t.Fatal(err)
		}
	})
	res, err := tick(t, repo, a.ID)
	if r, ok := docedit.IsRefused(err); ok {
		t.Fatalf("a story started meanwhile refused the tick: %+v", r.Findings)
	}
	if err != nil {
		t.Fatal(err)
	}
	narrative := repo.NarrativePath(c.ID)
	if len(findings(t, repo, "narrative.state", narrative)) == 0 {
		t.Fatal("the other story's narrative has no narrative.state finding: the test reproduces nothing")
	}
	for _, f := range res.Warnings {
		if filepath.Join(repo.Root, f.Path) == narrative {
			t.Errorf("a finding the scoped check leaves out is no warning: %+v", f)
		}
	}
	if got := Criteria(mustGet(t, repo, a.ID).Body); len(got) != 1 || !got[0].Ticked {
		t.Errorf("the criterion is not ticked: %+v", got)
	}
}

// I-0182, ADR-0197: a markdown finding on another story's task, written
// meanwhile and untracked, is outside the ticked story: the tick stands, and
// the finding is among its warnings, marked outside.
func TestATickWarnsOfAFindingOnAnotherStorysTaskWrittenBetweenTheChecks(t *testing.T) {
	repo, a, _, b, _ := raceRepo(t)
	var theirs string
	between(t, func(repo *workitem.Repo) {
		task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Task of B", Parent: b.ID, Owner: "bob", Now: time.Date(2026, 10, 1, 9, 30, 0, 0, time.UTC)})
		if err != nil {
			t.Fatal(err)
		}
		theirs = task.Path
		blankLines(t, theirs)
	})
	res, err := tick(t, repo, a.ID)
	if r, ok := docedit.IsRefused(err); ok {
		t.Fatalf("another story's task written meanwhile refused the tick: %+v", r.Findings)
	}
	if err != nil {
		t.Fatal(err)
	}
	warned := false
	for _, f := range res.Warnings {
		if f.Rule == "markdown.MD012" && filepath.Join(repo.Root, f.Path) == theirs {
			warned = true
			if !f.Outside {
				t.Errorf("the warning is not marked outside the story: %+v", f)
			}
		}
	}
	if !warned {
		t.Errorf("no markdown.MD012 warning on %s: %+v", theirs, res.Warnings)
	}
	if got := Criteria(mustGet(t, repo, a.ID).Body); len(got) != 1 || !got[0].Ticked {
		t.Errorf("the criterion is not ticked: %+v", got)
	}
}

// ADR-0197: a new finding inside the ticked story, here on its own task,
// still refuses the tick, and the story is left as it was.
func TestATickIsRefusedOnANewFindingInsideItsStory(t *testing.T) {
	repo, a, task, _, _ := raceRepo(t)
	was, _ := os.ReadFile(a.Path)
	between(t, func(*workitem.Repo) { blankLines(t, task.Path) })
	_, err := tick(t, repo, a.ID)
	r, ok := docedit.IsRefused(err)
	if !ok {
		t.Fatalf("want a refusal, got %v", err)
	}
	if len(r.Findings) == 0 {
		t.Error("the refusal names no finding")
	}
	for _, f := range r.Findings {
		if f.Rule != "markdown.MD012" || filepath.Join(repo.Root, f.Path) != task.Path || f.Outside {
			t.Errorf("finding: %+v", f)
		}
	}
	if now, _ := os.ReadFile(a.Path); string(now) != string(was) {
		t.Errorf("a refused tick changed the story:\n%s", now)
	}
}

// ADR-0197: an epic has no story, so an edit of one refuses a new finding
// only on a path it wrote; a finding on a story's task written meanwhile is
// a warning, marked outside.
func TestAnEpicEditWarnsOfANewFindingOnAPathItDidNotWrite(t *testing.T) {
	repo, _, task, _, _ := raceRepo(t)
	epic := mustGet(t, repo, "E-0001")
	between(t, func(*workitem.Repo) { blankLines(t, task.Path) })
	title := "Epic renamed"
	res, err := Apply(repo, nil, epic.ID, Change{Title: &title}, Options{By: "alex", NoCommit: true, Now: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC)})
	if err != nil {
		t.Fatalf("a finding on a path the edit did not write refused it: %v", err)
	}
	warned := false
	for _, f := range res.Warnings {
		warned = warned || (f.Rule == "markdown.MD012" && filepath.Join(repo.Root, f.Path) == task.Path && f.Outside)
	}
	if !warned {
		t.Errorf("no markdown.MD012 warning outside on %s: %+v", task.Path, res.Warnings)
	}
}

func mustGet(t *testing.T, repo *workitem.Repo, id string) *workitem.Item {
	t.Helper()
	it, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	return it
}

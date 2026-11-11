package itemedit

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Widen adds only what the touches do not cover, a folder or a component
// covering what lies below it, keeps the touches there were, and leaves the
// edit notice flai touches leaves; with nothing to add it writes nothing.
func TestWidenAddsWhatTheTouchesDoNotCover(t *testing.T) {
	at := time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	s := story("Mine", workitem.InProgress, "design", "README.md", "cli")
	later := at.Add(time.Hour)

	added, err := Widen(repo, s.ID, []string{"design/issues/I-0001.md", "README.md", "flai/cmd/x.go", "docs/a.md"}, "claude", func() time.Time { return later })
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"docs/a.md"}; !reflect.DeepEqual(added, want) {
		t.Errorf("added %v, want %v", added, want)
	}
	got, err := repo.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"design", "README.md", "cli", "docs/a.md"}; !reflect.DeepEqual(got.Touches, want) {
		t.Errorf("touches %v, want %v", got.Touches, want)
	}
	if got.Updated != later.Format(workitem.TimeFormat) {
		t.Errorf("updated %s", got.Updated)
	}
	notices := Notices(repo)
	if len(notices) != 1 || notices[0].ID != s.ID || notices[0].By != "claude" || !reflect.DeepEqual(notices[0].Changed, []string{"touches"}) {
		t.Errorf("notices %+v", notices)
	}

	added, err = Widen(repo, s.ID, []string{"docs/a.md", "design/x.md"}, "claude", func() time.Time { return later.Add(time.Hour) })
	if err != nil || added == nil || len(added) != 0 {
		t.Fatalf("nothing to add: %v %v", added, err)
	}
	if got, _ := repo.Get(s.ID); got.Updated != later.Format(workitem.TimeFormat) || len(Notices(repo)) != 1 {
		t.Errorf("an edit with nothing to add was written: updated %s, %d notices", got.Updated, len(Notices(repo)))
	}
	if _, err := Widen(repo, "S-0099", []string{"a.md"}, "claude", func() time.Time { return later }); err == nil {
		t.Error("no item to widen is not an error")
	}
}

// S-0522, I-0183: WidenStory reads the clock it stamps the task and the
// story with after it reads each, so that a move made before the read is not
// stamped over by an earlier time. The clock here takes away the file of the
// item it is read for, so a widening that read the clock before the item
// finds nothing to widen; each stamp is the reading taken for its item.
func TestWidenStoryReadsEachItemBeforeItsClock(t *testing.T) {
	at := time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: mine.ID, Owner: "alex", Touches: []string{"docs/a.md"}, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	stamped := []*workitem.Item{task, mine} // in the order WidenStory widens them
	aside := t.TempDir()
	readings := 0
	clock := func() time.Time {
		if readings < len(stamped) {
			p := stamped[readings].Path
			if err := os.Rename(p, filepath.Join(aside, filepath.Base(p))); err != nil {
				t.Fatal(err)
			}
		}
		readings++
		return at.Add(time.Duration(readings) * time.Minute)
	}

	w, err := WidenStory(WidenOptions{Repo: repo, Story: mine.ID, Task: task.ID, By: "claude", Now: clock, Paths: []string{"flai/cmd/x.go"}})
	if err != nil {
		t.Fatalf("a widening read its clock before the item it stamps: %v", err)
	}
	if !reflect.DeepEqual(w.Task, []string{"flai/cmd/x.go"}) || !reflect.DeepEqual(w.Story, []string{"flai/cmd/x.go"}) {
		t.Errorf("added %+v, want flai/cmd/x.go to the task and the story", w)
	}
	notices := Notices(repo)
	for i, it := range stamped {
		got, err := repo.Get(it.ID)
		if err != nil {
			t.Fatal(err)
		}
		want := at.Add(time.Duration(i+1) * time.Minute).Format(workitem.TimeFormat)
		if got.Updated != want || !slices.Contains(got.Touches, "flai/cmd/x.go") {
			t.Errorf("%s updated %s, touches %v; want stamped %s with flai/cmd/x.go", it.ID, got.Updated, got.Touches, want)
		}
		if len(notices) != len(stamped) || notices[i].ID != it.ID || notices[i].At != want {
			t.Errorf("notices %+v; want %s's at %s", notices, it.ID, want)
		}
	}
}

// WidenStory leaves the wip folder out, widens the task and then the story,
// each with what its own touches do not cover, and reports the story in
// progress the story's claim grew into.
func TestWidenStoryWidensTheTaskAndTheStoryOutsideWip(t *testing.T) {
	at := time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	mine := story("Mine", workitem.InProgress, "design")
	theirs := story("Theirs", workitem.InProgress, "cli")
	task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: mine.ID, Owner: "alex", Touches: []string{"docs/a.md"}, Now: at})
	if err != nil {
		t.Fatal(err)
	}

	w, err := WidenStory(WidenOptions{
		Repo: repo, Story: mine.ID, Task: task.ID, By: "claude", Now: func() time.Time { return at },
		Paths: []string{"wip/kanban/stories/x.md", "docs/a.md", "flai/cmd/x.go", "design/y.md"},
	})
	if err != nil || w.Untold != nil {
		t.Fatal(err, w.Untold)
	}
	if want := []string{"flai/cmd/x.go", "design/y.md"}; !reflect.DeepEqual(w.Task, want) {
		t.Errorf("task added %v, want %v", w.Task, want)
	}
	if want := []string{"docs/a.md", "flai/cmd/x.go"}; !reflect.DeepEqual(w.Story, want) {
		t.Errorf("story added %v, want %v", w.Story, want)
	}
	if want := []Overlapping{{Story: theirs.ID, Title: "Theirs", Paths: []string{"flai/cmd/x.go"}}}; !reflect.DeepEqual(w.Overlaps, want) {
		t.Errorf("overlaps %+v, want %+v", w.Overlaps, want)
	}
	if got, _ := repo.Get(mine.ID); !reflect.DeepEqual(got.Touches, []string{"design", "docs/a.md", "flai/cmd/x.go"}) {
		t.Errorf("story touches %v", got.Touches)
	}

	again, err := WidenStory(WidenOptions{Repo: repo, Story: mine.ID, Paths: []string{"wip/agents/x.md", "design/z.md"}, By: "claude"})
	if err != nil || len(again.Task) != 0 || len(again.Story) != 0 || len(again.Overlaps) != 0 || again.Overlaps == nil {
		t.Errorf("nothing to add: %+v %v", again, err)
	}
}

// ADR-0284: NarrowStory removes from the closing task's touches each file the
// branch does not change, keeping a changed file, a folder in the tree, a
// folder no longer in the tree whose files the branch deleted, a
// sub-project's name or tag, and a path under wip, and removes the same file
// from the story's touches only when no other open task names it; a done or
// cancelled task's touches do not keep it.
func TestNarrowStoryRemovesThePredictedFilesTheBranchDoesNotChange(t *testing.T) {
	at := time.Date(2026, 11, 11, 9, 0, 0, 0, time.UTC)
	repo, story := claimRepo(t, at)
	tree := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tree, "docs/guide"), 0o755); err != nil {
		t.Fatal(err)
	}
	mine := story("Mine", workitem.InProgress, "docs/a.md", "docs/b.md", "docs/c.md", "design", "src/x.go")
	task := func(title string, touches ...string) *workitem.Item {
		t.Helper()
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: title, Parent: mine.ID, Owner: "alex", Touches: touches, Now: at})
		if err != nil {
			t.Fatal(err)
		}
		return it
	}
	closing := task("Closing", "src/x.go", "docs/a.md", "docs/b.md", "docs/c.md", "docs/guide", "gone", "cli", "wip/kanban/x.md")
	task("Open", "docs/b.md")
	done := task("Done", "docs/c.md")
	done.Status = workitem.Done
	if err := repo.Save(done); err != nil {
		t.Fatal(err)
	}
	later := at.Add(time.Hour)

	n, err := NarrowStory(NarrowOptions{
		Repo: repo, Story: mine.ID, Task: closing.ID, Tree: tree, By: "claude", Now: func() time.Time { return later },
		Changed: []string{"src/x.go", "flai/cmd/y.go", "gone/old.md"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"docs/a.md", "docs/b.md", "docs/c.md"}; !reflect.DeepEqual(n.Task, want) {
		t.Errorf("task narrowed %v, want %v", n.Task, want)
	}
	if want := []string{"docs/a.md", "docs/c.md"}; !reflect.DeepEqual(n.Story, want) {
		t.Errorf("story narrowed %v, want %v", n.Story, want)
	}
	if got, _ := repo.Get(closing.ID); !reflect.DeepEqual(got.Touches, []string{"src/x.go", "docs/guide", "gone", "cli", "wip/kanban/x.md"}) || got.Updated != later.Format(workitem.TimeFormat) {
		t.Errorf("task touches %v, updated %s", got.Touches, got.Updated)
	}
	if got, _ := repo.Get(mine.ID); !reflect.DeepEqual(got.Touches, []string{"docs/b.md", "design", "src/x.go"}) {
		t.Errorf("story touches %v", got.Touches)
	}
	if notices := Notices(repo); len(notices) != 2 || notices[0].ID != closing.ID || notices[1].ID != mine.ID {
		t.Errorf("notices %+v", notices)
	}

	again, err := NarrowStory(NarrowOptions{Repo: repo, Story: mine.ID, Task: closing.ID, Tree: tree, By: "claude", Changed: []string{"src/x.go", "gone/old.md"}})
	if err != nil || again.Task == nil || again.Story == nil || len(again.Task)+len(again.Story) != 0 {
		t.Errorf("nothing to narrow: %+v %v", again, err)
	}
	if len(Notices(repo)) != 2 {
		t.Error("an edit with nothing to narrow was written")
	}
}

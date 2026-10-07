package itemedit

import (
	"reflect"
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

	added, err := Widen(repo, s.ID, []string{"design/issues/I-0001.md", "README.md", "flai/cmd/x.go", "docs/a.md"}, "claude", later)
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

	added, err = Widen(repo, s.ID, []string{"docs/a.md", "design/x.md"}, "claude", later.Add(time.Hour))
	if err != nil || added == nil || len(added) != 0 {
		t.Fatalf("nothing to add: %v %v", added, err)
	}
	if got, _ := repo.Get(s.ID); got.Updated != later.Format(workitem.TimeFormat) || len(Notices(repo)) != 1 {
		t.Errorf("an edit with nothing to add was written: updated %s, %d notices", got.Updated, len(Notices(repo)))
	}
	if _, err := Widen(repo, "S-0099", []string{"a.md"}, "claude", later); err == nil {
		t.Error("no item to widen is not an error")
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

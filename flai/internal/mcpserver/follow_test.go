package mcpserver

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readyUnderBacklogEpic is a new epic in backlog with a story under it that
// is ready, put there without the epic following, as an epic moved back by
// hand would leave it (S-0200).
func (f *fixture) readyUnderBacklogEpic(t *testing.T) (epic, story *workitem.Item) {
	t.Helper()
	epic, err := f.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Second epic", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	story, err = f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Leads", Parent: epic.ID, Owner: "alex", Touches: []string{"docs/leads"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(story.Path)
	_ = os.WriteFile(story.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	story, _ = f.repo.Get(story.ID)
	items, _ := f.repo.List(false)
	if _, err := f.repo.Move(story, workitem.Ready, workitem.MoveOptions{By: "alex", Now: t0, Items: items}); err != nil {
		t.Fatal(err)
	}
	if err := f.repo.Save(story); err != nil {
		t.Fatal(err)
	}
	return epic, story
}

func TestItemMoveReportsTheEpicThatFollowed(t *testing.T) {
	f := setup(t)
	epic, story := f.readyUnderBacklogEpic(t)
	*f.clock = t0.Add(2 * time.Minute)
	out, failed := f.call(t, "item_move", map[string]any{"id": story.ID, "to": "in-progress"})
	if failed != "" {
		t.Fatal(failed)
	}
	got, _ := out["followed"].(map[string]any)
	want := map[string]any{"id": epic.ID, "type": "epic", "title": "Second epic", "from": "backlog", "to": "in-progress", "story": story.ID}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("followed[%s] = %v, want %v (followed %v)", k, got[k], v, out["followed"])
		}
	}
	if e, _ := f.repo.Get(epic.ID); e.Status != workitem.InProgress {
		t.Errorf("%s is %s, want in-progress", epic.ID, e.Status)
	}
}

func TestItemMoveThatLeavesItsEpicReportsNoFollowed(t *testing.T) {
	f := setup(t)
	if e, _ := f.repo.Get(f.story.Parent); e.Status != workitem.InProgress {
		t.Fatalf("the fixture's epic is %s, want in-progress with its story", e.Status)
	}
	s := f.readyStory(t, "Second", t0)
	*f.clock = t0.Add(2 * time.Minute)
	out, failed := f.call(t, "item_move", map[string]any{"id": s.ID, "to": "in-progress"})
	if failed != "" {
		t.Fatal(failed)
	}
	if _, ok := out["followed"]; ok {
		t.Errorf("the epic was already in progress, yet followed is %v", out["followed"])
	}
}

// Another actor's story move takes its epic along: the agent's inbox says
// which story each of the epic's moves followed.
func TestInboxSaysWhichStoryAnEpicFollowed(t *testing.T) {
	f := setup(t)
	epic, story := f.readyUnderBacklogEpic(t)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}
	if _, err := f.repo.TransitionAll(story, workitem.InProgress, "alex", "", t0.Add(5*time.Minute), false); err != nil {
		t.Fatal(err)
	}
	*f.clock = t0.Add(10 * time.Minute)
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	got := changeSummaries(out, "changes")
	want := []string{
		epic.ID + " Second epic moved to ready by alex, following " + story.ID,
		epic.ID + " Second epic moved to in-progress by alex, following " + story.ID,
		story.ID + " Leads moved to in-progress by alex",
	}
	if len(got) != len(want) {
		t.Fatalf("changes: %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("change %d: %q, want %q", i, got[i], want[i])
		}
	}
	for _, c := range out["changes"].([]any) {
		m := c.(map[string]any)
		if m["id"] == epic.ID && m["follows"] != story.ID {
			t.Errorf("%v does not name %s in follows", m, story.ID)
		}
	}
}

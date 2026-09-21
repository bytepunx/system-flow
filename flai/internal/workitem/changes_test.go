package workitem

import (
	"testing"
	"time"
)

func TestChangesSinceACursor(t *testing.T) {
	at := func(s string) time.Time {
		v, err := time.Parse(TimeFormat, s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	items := []*Item{
		{ID: "S-0002", Type: Story, Title: "Two", Transitions: []Transition{
			{To: Ready, At: "2026-09-19T10:00:05Z", By: "alex"},
			{To: InProgress, At: "2026-09-19T10:00:09Z", By: "claude"},
		}},
		{ID: "S-0001", Type: Story, Title: "One", Transitions: []Transition{
			{To: Ready, At: "2026-09-19T09:00:00Z", By: "alex"}, // before the cursor
		}, Blocked: []Block{
			{From: "2026-09-19T10:00:05Z", Until: "2026-09-19T10:00:07Z", Reason: "waiting"},
		}},
	}
	kinds := func(cs []Change) (out []string) {
		for _, c := range cs {
			out = append(out, c.ID+" "+c.Kind+" "+c.To)
		}
		return out
	}
	eq := func(name string, got, want []string) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("%s: got %v, want %v", name, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("%s: got %v, want %v", name, got, want)
			}
		}
	}

	// what others did since 10:00:00, oldest first; claude's own move is left out
	got := Changes(items, at("2026-09-19T10:00:00Z"), "claude", nil)
	eq("since 10:00:00 as claude", kinds(got), []string{"S-0001 blocked ", "S-0002 moved ready", "S-0001 unblocked "})
	if got[0].Reason != "waiting" || got[1].By != "alex" || got[1].Title != "Two" {
		t.Errorf("details: %+v", got)
	}
	// someone else sees claude's move too
	eq("as alex", kinds(Changes(items, at("2026-09-19T10:00:08Z"), "alex", nil)), []string{"S-0002 moved in-progress"})

	// The boundary second: a cursor at 10:00:05 that has reported the block
	// still gets the move stamped with the same second, once.
	seen := map[string]bool{got[0].Key(): true}
	second := Changes(items, at("2026-09-19T10:00:05Z"), "claude", seen)
	eq("boundary second", kinds(second), []string{"S-0002 moved ready", "S-0001 unblocked "})
	seen[second[0].Key()] = true
	eq("boundary second again", kinds(Changes(items, at("2026-09-19T10:00:05Z"), "claude", seen)), []string{"S-0001 unblocked "})
	// nothing after the last change
	if cs := Changes(items, at("2026-09-19T10:00:09Z"), "claude", map[string]bool{}); len(cs) != 0 {
		t.Errorf("claude's own move at the cursor second must not appear: %v", kinds(cs))
	}
}

func TestBoardViewReadyInPullOrder(t *testing.T) {
	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	mk := func(id, status string) *Item {
		return &Item{ID: id, Type: Story, Title: id, Status: status, Created: "2026-09-19T11:00:00Z"}
	}
	items := []*Item{mk("S-0001", Ready), mk("S-0002", Ready), mk("S-0003", Ready), mk("S-0004", InProgress), mk("S-0005", InProgress),
		{ID: "T-0001", Type: Task, Status: Ready, Created: "2026-09-19T11:00:00Z"},
		{ID: "S-0009", Type: Story, Status: Ready, Archived: true}}
	board := &Board{WIPLimits: map[string]int{InProgress: 2, Ready: 2}, Order: []string{"S-0003", "S-0001"}}
	v := NewBoardView(items, board, now, false, nil)
	var ids []string
	for _, c := range v.ReadyInPullOrder() {
		ids = append(ids, c.ID)
	}
	if len(ids) != 3 || ids[0] != "S-0003" || ids[1] != "S-0001" || ids[2] != "S-0002" {
		t.Errorf("ready in pull order, unlisted last: %v", ids)
	}
	if v.CanPull() {
		t.Error("two in progress against a limit of two leaves no room")
	}
	if len(v.Breaches) != 1 || v.Columns[Ready][0].Age != "1h" {
		t.Errorf("breaches %v, age %q", v.Breaches, v.Columns[Ready][0].Age)
	}
	if len(NewBoardView(items, board, now, true, nil).Columns[Ready]) != 4 {
		t.Error("all adds the task and still leaves the archived story out")
	}
}

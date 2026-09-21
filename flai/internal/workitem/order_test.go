package workitem

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func story(id, status string) *Item { return &Item{ID: id, Type: Story, Status: status} }

// A board with two ready stories named, one backlog story placed, and three
// backlog stories nobody has placed.
func orderFixture() (*Board, []*Item) {
	items := []*Item{
		story("S-0001", Ready), story("S-0002", Ready),
		story("S-0003", Backlog), story("S-0004", Backlog), story("S-0005", Backlog), story("S-0006", Backlog),
		story("S-0007", InProgress),
		{ID: "E-0001", Type: Epic, Status: Backlog},
		{ID: "T-0001", Type: Task, Status: Ready},
		{ID: "S-0008", Type: Story, Status: Backlog, Archived: true},
	}
	return &Board{Order: []string{"S-0002", "S-0001", "S-0005"}}, items
}

func TestPullSequenceNamesFirstThenByID(t *testing.T) {
	b, items := orderFixture()
	if got, want := PullSequence(b.Order, items, Ready), []string{"S-0002", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ready: got %v want %v", got, want)
	}
	if got, want := PullSequence(b.Order, items, Backlog), []string{"S-0005", "S-0003", "S-0004", "S-0006"}; !reflect.DeepEqual(got, want) {
		t.Errorf("backlog: named first, the rest by ID, no epic, no archived story: got %v want %v", got, want)
	}
	if got := PullSequence([]string{"S-0009", "S-0007", "S-0001", "S-0001"}, items, Ready); !reflect.DeepEqual(got, []string{"S-0001", "S-0002"}) {
		t.Errorf("unknown, other-state, and repeated names are ignored: got %v", got)
	}
}

func TestPlace(t *testing.T) {
	cases := []struct {
		name      string
		id        string
		p         Placement
		wantOrder []string
		wantSeq   []string
	}{
		{"ready to the top", "S-0001", Placement{Top: true},
			[]string{"S-0001", "S-0002", "S-0005"}, []string{"S-0001", "S-0002"}},
		{"ready after another", "S-0002", Placement{After: "S-0001"},
			[]string{"S-0001", "S-0002", "S-0005"}, []string{"S-0001", "S-0002"}},
		{"ready to the bottom, already there", "S-0001", Placement{Bottom: true},
			[]string{"S-0002", "S-0001", "S-0005"}, []string{"S-0002", "S-0001"}},
		{"an unplaced backlog story to the top names only itself", "S-0006", Placement{Top: true},
			[]string{"S-0002", "S-0001", "S-0006", "S-0005"}, []string{"S-0006", "S-0005", "S-0003", "S-0004"}},
		{"before an unplaced story names as far down as it must", "S-0006", Placement{Before: "S-0004"},
			[]string{"S-0002", "S-0001", "S-0005", "S-0003", "S-0006"}, []string{"S-0005", "S-0003", "S-0006", "S-0004"}},
		{"a placed story to the bottom is named last, after everything", "S-0005", Placement{Bottom: true},
			[]string{"S-0002", "S-0001", "S-0003", "S-0004", "S-0006", "S-0005"}, []string{"S-0003", "S-0004", "S-0006", "S-0005"}},
		{"an unplaced story to the bottom where it already is names nothing new", "S-0006", Placement{Bottom: true},
			[]string{"S-0002", "S-0001", "S-0005"}, []string{"S-0005", "S-0003", "S-0004", "S-0006"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, items := orderFixture()
			if err := b.Place(items, c.id, c.p); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(b.Order, c.wantOrder) {
				t.Errorf("order: got %v want %v", b.Order, c.wantOrder)
			}
			var st string
			for _, it := range items {
				if it.ID == c.id {
					st = it.Status
				}
			}
			if got := PullSequence(b.Order, items, st); !reflect.DeepEqual(got, c.wantSeq) {
				t.Errorf("read back: got %v want %v", got, c.wantSeq)
			}
		})
	}
}

func TestPlaceKeepsReadyBeforeBacklogAndDropsStaleNames(t *testing.T) {
	_, items := orderFixture()
	b := &Board{Order: []string{"S-0005", "S-0007", "S-0099", "S-0002", "E-0001", "S-0001"}}
	if err := b.Place(items, "S-0004", Placement{After: "S-0005"}); err != nil {
		t.Fatal(err)
	}
	if want := []string{"S-0002", "S-0001", "S-0005", "S-0004"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("got %v want %v", b.Order, want)
	}
}

func TestPlaceRefusals(t *testing.T) {
	cases := []struct {
		name string
		id   string
		p    Placement
		want string
	}{
		{"no position", "S-0001", Placement{}, "exactly one of"},
		{"two positions", "S-0001", Placement{Top: true, Before: "S-0002"}, "exactly one of"},
		{"an epic", "E-0001", Placement{Top: true}, "only stories are in the pull order"},
		{"a task", "T-0001", Placement{Top: true}, "only stories are in the pull order"},
		{"a story in progress", "S-0007", Placement{Top: true}, "only backlog and ready stories have a pull order"},
		{"an archived story", "S-0008", Placement{Top: true}, "not an active item"},
		{"an unknown story", "S-0099", Placement{Top: true}, "not an active item"},
		{"itself", "S-0001", Placement{Before: "S-0001"}, "relative to itself"},
		{"across columns", "S-0001", Placement{Before: "S-0003"}, "within one column"},
		{"relative to a task", "S-0001", Placement{After: "T-0001"}, "only stories are in the pull order"},
		{"relative to nothing known", "S-0001", Placement{After: "S-0099"}, "not an active item"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, items := orderFixture()
			before := append([]string{}, b.Order...)
			err := b.Place(items, c.id, c.p)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("got %v, want an error containing %q", err, c.want)
			}
			if !reflect.DeepEqual(b.Order, before) {
				t.Errorf("a refusal changed the order: %v", b.Order)
			}
		})
	}
}

// A placed backlog story that becomes ready joins the end of the ready
// section, not the end of the list behind backlog names.
func TestPlaceReadyLast(t *testing.T) {
	b, items := orderFixture()
	b.PlaceReadyLast("S-0005", items)
	if want := []string{"S-0002", "S-0001", "S-0005"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("already last of the ready section once ready: got %v want %v", b.Order, want)
	}
	b = &Board{Order: []string{"S-0002", "S-0005", "S-0003"}}
	b.PlaceReadyLast("S-0003", items)
	if want := []string{"S-0002", "S-0003", "S-0005"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("before the backlog names: got %v want %v", b.Order, want)
	}
	b = &Board{Order: []string{"S-0005"}}
	b.PlaceReadyLast("S-0004", items)
	if want := []string{"S-0004", "S-0005"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("with no ready story named it leads: got %v want %v", b.Order, want)
	}
	b = &Board{Order: []string{"S-0005"}}
	b.PlaceReadyLast("S-0004", nil)
	if want := []string{"S-0005", "S-0004"}; !reflect.DeepEqual(b.Order, want) {
		t.Errorf("without the items it appends as before: got %v want %v", b.Order, want)
	}
}

func TestBoardViewColumnsReadInPullOrder(t *testing.T) {
	b, items := orderFixture()
	now := time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)
	ids := func(cards []BoardCard) (out []string) {
		for _, c := range cards {
			out = append(out, c.ID)
		}
		return
	}
	v := NewBoardView(items, b, now, false, nil)
	if got, want := ids(v.Columns[Ready]), []string{"S-0002", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ready: got %v want %v", got, want)
	}
	if got, want := ids(v.Columns[Backlog]), []string{"S-0005", "S-0003", "S-0004", "S-0006"}; !reflect.DeepEqual(got, want) {
		t.Errorf("backlog: got %v want %v", got, want)
	}
	if got, want := ids(v.ReadyInPullOrder()), []string{"S-0002", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ReadyInPullOrder: got %v want %v", got, want)
	}
	all := NewBoardView(items, b, now, true, nil)
	if got, want := ids(all.Columns[Ready]), []string{"S-0002", "S-0001", "T-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("with tasks, stories take the sequence in the stories' places: got %v want %v", got, want)
	}
	if got, want := ids(all.Columns[Backlog]), []string{"S-0005", "S-0003", "S-0004", "S-0006", "E-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("with epics: got %v want %v", got, want)
	}
}

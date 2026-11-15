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
	v := NewBoardView(items, b, now, false, nil, nil)
	if got, want := ids(v.Columns[Ready]), []string{"S-0002", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ready: got %v want %v", got, want)
	}
	if got, want := ids(v.Columns[Backlog]), []string{"S-0005", "S-0003", "S-0004", "S-0006"}; !reflect.DeepEqual(got, want) {
		t.Errorf("backlog: got %v want %v", got, want)
	}
	if got, want := ids(v.ReadyInPullOrder()), []string{"S-0002", "S-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ReadyInPullOrder: got %v want %v", got, want)
	}
	all := NewBoardView(items, b, now, true, nil, nil)
	if got, want := ids(all.Columns[Ready]), []string{"S-0002", "S-0001", "T-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("with tasks, stories take the sequence in the stories' places: got %v want %v", got, want)
	}
	if got, want := ids(all.Columns[Backlog]), []string{"S-0005", "S-0003", "S-0004", "S-0006", "E-0001"}; !reflect.DeepEqual(got, want) {
		t.Errorf("with epics: got %v want %v", got, want)
	}
}

// S-0739, ADR-0301 §2: the pull order is read from the states the workflow
// gives an order and from the pull's source, never from their names, so
// placing a story, naming one that has just entered the pull's source, the
// refusals, and the board's columns are the same under the fixture, whose
// story in doing has no order as one in-progress has none, nor one in
// testing.
func TestOrderKindUnderBothFlows(t *testing.T) {
	for _, f := range planKindFlows(t) {
		t.Run(f.name, func(t *testing.T) {
			w := f.w
			board := func() (*Board, []*Item) {
				b, items := orderFixture()
				for _, it := range items {
					it.Status = f.state(it.Status)
				}
				b.Workflow = &w
				return b, items
			}
			for _, name := range w.Names() {
				if got, want := OrderedIn(w, name), name == Backlog || name == Ready; got != want {
					t.Errorf("%s ordered: %v, want %v", name, got, want)
				}
			}
			for _, c := range []struct {
				id   string
				p    Placement
				want []string
			}{
				{"S-0001", Placement{Top: true}, []string{"S-0001", "S-0002", "S-0005"}},
				{"S-0006", Placement{Before: "S-0004"}, []string{"S-0002", "S-0001", "S-0005", "S-0003", "S-0006"}},
				{"S-0005", Placement{Bottom: true}, []string{"S-0002", "S-0001", "S-0003", "S-0004", "S-0006", "S-0005"}},
			} {
				b, items := board()
				if err := b.Place(items, c.id, c.p); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(b.Order, c.want) {
					t.Errorf("place %s %+v: got %v want %v", c.id, c.p, b.Order, c.want)
				}
			}
			b, items := board()
			unordered := []*Item{items[6] /* S-0007, active */, story("S-0010", f.state(Review)), story("S-0011", f.state(Done))}
			if f.name == "fixture" {
				unordered = append(unordered, story("S-0012", fixtureTesting))
			}
			items = append(items, unordered[1:]...)
			for _, it := range unordered {
				want := "rule: " + it.ID + " is " + it.Status + "; only backlog and ready stories have a pull order"
				if err := b.Place(items, it.ID, Placement{Top: true}); err == nil || err.Error() != want {
					t.Errorf("got %v, want %s", err, want)
				}
			}
			b.PlaceReadyLast("S-0005", items)
			if want := []string{"S-0002", "S-0001", "S-0005"}; !reflect.DeepEqual(b.Order, want) {
				t.Errorf("ready last: got %v want %v", b.Order, want)
			}

			b, items = board()
			items = append(items, &Item{ID: "T-0002", Type: Task, Parent: "S-0007", Status: f.state(InProgress)}, &Item{ID: "T-0003", Type: Task, Parent: "S-0007", Status: f.state(Ready), After: []string{"T-0002"}})
			v := NewBoardView(items, b, time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC), false, nil, nil)
			var ready, backlog []string
			for _, c := range v.Columns[Ready] {
				ready = append(ready, c.ID)
			}
			for _, c := range v.Columns[Backlog] {
				backlog = append(backlog, c.ID)
			}
			if want := []string{"S-0002", "S-0001"}; !reflect.DeepEqual(ready, want) {
				t.Errorf("ready column: got %v want %v", ready, want)
			}
			if want := []string{"S-0005", "S-0003", "S-0004", "S-0006"}; !reflect.DeepEqual(backlog, want) {
				t.Errorf("backlog column: got %v want %v", backlog, want)
			}
			if cards := v.Columns[f.state(InProgress)]; len(cards) != 1 || cards[0].Tasks == nil || *cards[0].Tasks != (TaskSummary{Waiting: 1, InProgress: 1, Layers: 2}) {
				t.Errorf("the active story's task plan: %+v", cards)
			}
		})
	}
}

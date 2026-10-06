package workitem

import (
	"reflect"
	"testing"
)

const draftBody = "# S-0001 Draft\n\n## Goal\n\nDo it.\n\n## Acceptance criteria\n- [ ] it works\n\n## Tasks\n\n## Notes\n"

var storySections = []string{"Goal", "Acceptance criteria", "Tasks", "Notes"}

// completeDraft is a draft story under E-0001 that lacks nothing unless the
// test changes it.
func completeDraft() *Item {
	return &Item{
		ID: "S-0001", Type: Story, Status: Backlog, Parent: "E-0001", Title: "Draft", Draft: true,
		Body: draftBody, Touches: []string{"flai/cmd"},
		CostOfDelay: &CostOfDelay{Value: amount(300)},
		Forecast:    &Forecast{Duration: "4h", Delivery: "2026-10-09T12:00:00Z"},
	}
}

// S-0219: a draft is complete when it has every section of the template, a
// goal, a criterion as a checkbox, a touch, a forecast duration and
// delivery, a cost of delay value, and an open epic or none.
func TestFinalizableLacks(t *testing.T) {
	epic := &Item{ID: "E-0001", Type: Epic, Status: InProgress, Title: "Epic"}
	for _, c := range []struct {
		name   string
		change func(s *Item, e *Item) *Item
		want   []string
	}{
		{"complete", func(s, e *Item) *Item { return e }, []string{}},
		{"no epic", func(s, e *Item) *Item { s.Parent = ""; return nil }, []string{}},
		{"no criterion", func(s, e *Item) *Item {
			s.Body = "## Goal\n\nDo it.\n\n## Acceptance criteria\n- [ ]\n\n## Tasks\n\n## Notes\n"
			return e
		}, []string{"no acceptance criteria with a checkbox"}},
		{"no touches", func(s, e *Item) *Item { s.Touches = nil; return e }, []string{"no touches"}},
		{"no forecast", func(s, e *Item) *Item { s.Forecast = nil; return e }, []string{"no forecast duration", "no forecast delivery"}},
		{"no delivery", func(s, e *Item) *Item { s.Forecast.Delivery = ""; return e }, []string{"no forecast delivery"}},
		{"bad duration", func(s, e *Item) *Item { s.Forecast.Duration = "soon"; return e }, []string{"no forecast duration"}},
		{"no value", func(s, e *Item) *Item { s.CostOfDelay = &CostOfDelay{}; return e }, []string{"no cost of delay value"}},
		{"no goal or section", func(s, e *Item) *Item {
			s.Body = "## Acceptance criteria\n- [ ] it works\n\n## Notes\n"
			return e
		}, []string{`no "## Goal" section`, `no "## Tasks" section`, "no goal"}},
		{"cancelled epic", func(s, e *Item) *Item { e.Status = Cancelled; return e }, []string{"its epic E-0001 is cancelled"}},
		{"done epic", func(s, e *Item) *Item { e.Status = Done; return e }, []string{"its epic E-0001 is done"}},
		{"missing epic", func(s, e *Item) *Item { return nil }, []string{"its epic E-0001 is not found"}},
	} {
		s, e := completeDraft(), *epic
		parent := c.change(s, &e)
		if got := FinalizableLacks(s, parent, storySections); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: lacks %q, want %q", c.name, got, c.want)
		}
	}

	// A project's own template decides the sections.
	if got := FinalizableLacks(completeDraft(), epic, append(storySections, "Risks")); !reflect.DeepEqual(got, []string{`no "## Risks" section`}) {
		t.Errorf("a template section the draft lacks: %q", got)
	}
}

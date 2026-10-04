package planning

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func amount(v float64) *float64 { return &v }

// withInputs gives an item cost of delay inputs.
func withInputs(it *workitem.Item, in workitem.CostInputs) *workitem.Item {
	in.By, in.At = "olive", at(history)
	it.CostOfDelay = &workitem.CostOfDelay{Inputs: &in}
	return it
}

func epic(id string) *workitem.Item {
	return &workitem.Item{ID: id, Type: workitem.Epic, Status: workitem.InProgress, Nature: "feature"}
}

// childOf puts a story under an epic.
func childOf(it *workitem.Item, parent string) *workitem.Item {
	it.Parent = parent
	return it
}

// S-0210: revenue and penalty per week add to the hours lost per cycle at
// the hour rate times the cycles in a week, in the project's currency.
func TestCostOfDelayFromInputs(t *testing.T) {
	items := []*workitem.Item{withInputs(epic("E-0001"), workitem.CostInputs{
		RevenuePerWeek: amount(1200), PenaltyPerWeek: amount(300), TimeLostPerCycle: "10h",
	})}
	plan := manifest.Planning{Currency: "EUR", HourRate: amount(150), Cycle: "84h"}
	res, err := CostOfDelay(items, plan, "e-1")
	if err != nil {
		t.Fatal(err)
	}
	want := "1200.00 EUR of revenue a week plus 300.00 EUR of penalty a week plus 10h of time lost per 84h cycle at 150 EUR an hour, 2.00 cycles a week: 4500.00 EUR a week."
	if res.ID != "E-0001" || res.Value != 4500 || res.Currency != "EUR" || res.From != "inputs" || res.Basis != want || res.Epic != nil {
		t.Fatalf("got %+v\nbasis %q", res, res.Basis)
	}
	in := res.Inputs
	if in == nil || *in.RevenuePerWeek != 1200 || *in.PenaltyPerWeek != 300 || in.TimeLostPerCycle != "10h" || *in.HourRate != 150 || in.Cycle != "84h" || in.CyclesPerWeek != 2 {
		t.Errorf("inputs: %+v", in)
	}

	// a penalty alone needs no hour rate, and the value is rounded to cents
	items = append(items, childOf(withInputs(story("S-0001", workitem.Ready, "feature", "", 1, 0, now), workitem.CostInputs{PenaltyPerWeek: amount(33.333)}), "E-0001"))
	res, err = CostOfDelay(items, manifest.Planning{}, "S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if res.Value != 33.33 || res.Currency != "USD" || res.Basis != "33.33 USD of penalty a week: 33.33 USD a week." || res.Inputs.HourRate != nil || res.Inputs.Cycle != "" {
		t.Errorf("penalty alone: %+v %+v", res, res.Inputs)
	}
}

// S-0210: a story without inputs takes the share of its epic's value that
// its duration is of its epic's open stories without inputs: its own
// forecast.duration when it has one, else one worked out from history.
func TestCostOfDelayEpicShare(t *testing.T) {
	sibling := childOf(story("S-0011", workitem.Ready, "feature", "m1", 2, 2, now), "E-0001")
	sibling.Forecast = &workitem.Forecast{Duration: "2h", By: "planner", At: at(now)}
	items := append(past(), withInputs(epic("E-0001"), workitem.CostInputs{TimeLostPerCycle: "10h"}),
		childOf(story("S-0010", workitem.Ready, "feature", "m1", 2, 2, now), "E-0001"), // size 4: 27m
		sibling,
		childOf(story("S-0012", workitem.Done, "feature", "m1", 2, 2, now), "E-0001"),
		childOf(withInputs(story("S-0013", workitem.Backlog, "feature", "m1", 2, 2, now), workitem.CostInputs{RevenuePerWeek: amount(100)}), "E-0001"),
		childOf(story("S-0014", workitem.InProgress, "feature", "m1", 4, 4, now), "E-0001"), // size 8: 54m
		story("S-0015", workitem.Ready, "feature", "m1", 2, 2, now),
	)
	res, err := CostOfDelay(items, manifest.Planning{HourRate: amount(150)}, "S-0010")
	if err != nil {
		t.Fatal(err)
	}
	want := "S-0010's share of E-0001's 1500.00 USD a week, 27m of 3h21m forecast over its 3 open stories without inputs: 201.49 USD a week."
	if res.Value != 201.49 || res.From != "epic" || res.Basis != want || res.Inputs != nil {
		t.Fatalf("got %+v\nbasis %q", res, res.Basis)
	}
	e := res.Epic
	if e == nil || e.ID != "E-0001" || e.Value != 1500 || e.From != "inputs" || e.Share != 1620.0/12060 || len(e.Stories) != 3 {
		t.Fatalf("epic: %+v", e)
	}
	wantStories := []SharedStory{
		{ID: "S-0010", Duration: "27m", DurationSeconds: 1620},
		{ID: "S-0011", Duration: "2h", DurationSeconds: 7200, FromForecast: true},
		{ID: "S-0014", Duration: "54m", DurationSeconds: 3240},
	}
	for i := range wantStories {
		if e.Stories[i] != wantStories[i] {
			t.Errorf("story %d: got %+v, want %+v", i, e.Stories[i], wantStories[i])
		}
	}

	// an epic with a value and no inputs shares its value
	items[len(past())].CostOfDelay = &workitem.CostOfDelay{Value: amount(603), By: "olive", At: at(history)}
	res, err = CostOfDelay(items, manifest.Planning{}, "S-0014")
	if err != nil {
		t.Fatal(err)
	}
	if res.Value != 162 || res.Epic.From != "value" || res.Epic.Value != 603 {
		t.Errorf("from the epic's value: %+v %+v", res, res.Epic)
	}
}

// S-0210: only an open epic or story with inputs, or a story whose epic has
// inputs or a value, has a cost of delay; the inputs are the operator's.
func TestCostOfDelayRefusals(t *testing.T) {
	task := &workitem.Item{ID: "T-0001", Type: workitem.Task, Status: workitem.Ready, Parent: "S-0001"}
	valued := epic("E-0002")
	valued.CostOfDelay = &workitem.CostOfDelay{Value: amount(500), By: "olive", At: at(history)}
	items := []*workitem.Item{
		task, epic("E-0001"), valued,
		withInputs(epic("E-0003"), workitem.CostInputs{TimeLostPerCycle: "4h"}),
		childOf(story("S-0001", workitem.Ready, "feature", "", 1, 0, now), "E-0001"),
		story("S-0002", workitem.Ready, "feature", "", 1, 0, now),
		childOf(withInputs(story("S-0003", workitem.Done, "feature", "", 1, 0, now), workitem.CostInputs{RevenuePerWeek: amount(1)}), "E-0002"),
		childOf(story("S-0004", workitem.Ready, "feature", "", 1, 0, now), "E-0003"),
	}
	for _, c := range []struct{ id, plan, want string }{
		{"T-0001", "", "cannot work out the cost of delay of T-0001: it is a task, and only an epic or a story has a cost of delay"},
		{"S-0099", "", "cannot work out the cost of delay of S-0099: there is no such item"},
		{"S-0003", "", "cannot work out the cost of delay of S-0003: it is done"},
		{"E-0002", "", "E-0002 has no cost of delay inputs; the operator gives them with flai edit --revenue-per-week, --penalty-per-week, or --time-lost-per-cycle"},
		{"S-0001", "", "S-0001 has no cost of delay inputs, and its epic E-0001 has neither inputs nor a value; the operator gives inputs on the story or the epic"},
		{"S-0002", "", "S-0002 has no cost of delay inputs and no epic to take a share of; the operator gives its inputs"},
		{"E-0003", "", "E-0003's time lost per cycle needs planning.hour_rate in system-flow.yaml, and it is unset; the operator sets it"},
		{"S-0004", "", "cannot work out the cost of delay of S-0004: E-0003's time lost per cycle needs planning.hour_rate"},
		{"E-0003", "bad", "planning.cycle \"bad\" is not a duration"},
	} {
		plan := manifest.Planning{}
		if c.plan != "" {
			plan = manifest.Planning{HourRate: amount(100), Cycle: c.plan}
		}
		_, err := CostOfDelay(items, plan, c.id)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: got %v, want %q", c.id, err, c.want)
		}
	}
}

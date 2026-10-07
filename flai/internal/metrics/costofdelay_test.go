package metrics

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// costItem is an item with the cost of delay value given, nil for none,
// created at the time given and moved as the pairs of state and time say.
func costItem(id, typ string, value *float64, created string, moves ...string) *workitem.Item {
	it := &workitem.Item{ID: id, Type: typ, Status: workitem.Backlog, Created: created}
	for i := 0; i < len(moves); i += 2 {
		it.Transitions = append(it.Transitions, workitem.Transition{To: moves[i], At: moves[i+1]})
		it.Status = moves[i]
	}
	if value != nil {
		it.CostOfDelay = &workitem.CostOfDelay{Value: value}
	}
	return it
}

// costItems are stories with and without a cost of delay and an epic the
// story report leaves out. A value of 700 a week is 100 a day of waiting.
func costItems() []*workitem.Item {
	return []*workitem.Item{
		// waits from before the window to the middle of its second day, then is done
		costItem("S-0001", workitem.Story, val(700), "2026-08-25T06:00:00Z",
			workitem.Ready, "2026-08-29T18:00:00Z", workitem.InProgress, "2026-08-30T12:00:00Z", workitem.Done, "2026-08-31T09:00:00Z"),
		// still in ready at now
		costItem("S-0002", workitem.Story, val(100.5), "2026-08-31T00:00:00Z", workitem.Ready, "2026-09-01T06:00:00Z"),
		// no value: adds nothing
		costItem("S-0003", workitem.Story, nil, "2026-08-20T00:00:00Z"),
		// moves through three columns on one day
		costItem("S-0004", workitem.Story, val(1400), "2026-08-30T00:00:00Z",
			workitem.Ready, "2026-08-30T06:00:00Z", workitem.InProgress, "2026-08-30T08:00:00Z", workitem.Review, "2026-08-31T10:00:00Z"),
		// cancelled from backlog: nothing after it closed
		costItem("S-0005", workitem.Story, val(70), "2026-08-28T00:00:00Z", workitem.Cancelled, "2026-08-29T12:00:00Z"),
		costItem("E-0001", workitem.Epic, val(7000), "2026-08-20T00:00:00Z"),
	}
}

// threeDays starts the window at noon on Saturday 29 August 2026, in the ISO
// week that started on Monday 24 August.
const threeDays = 3 * 24 * time.Hour

// S-0205: each item's value and what its time in backlog and ready cost up
// to now, absent without a value.
func TestCostOfDelayPerItemIsItsValueOverItsWeeksWaiting(t *testing.T) {
	rep := Compute(costItems(), Options{Now: now, Since: threeDays})
	want := map[string][2]*float64{ // value, incurred
		"S-0001": {val(700), val(525)},
		"S-0002": {val(100.5), val(21.54)},
		"S-0003": {nil, nil},
		"S-0004": {val(1400), val(66.67)},
		"S-0005": {val(70), val(15)},
	}
	if len(rep.Items) != len(want) {
		t.Fatalf("items = %d, want %d", len(rep.Items), len(want))
	}
	for _, m := range rep.Items {
		w := want[m.ID]
		for k, got := range [2]*float64{m.CostOfDelay, m.CostIncurred} {
			if (got == nil) != (w[k] == nil) || (got != nil && *got != *w[k]) {
				t.Errorf("%s [value, incurred][%d] = %v, want %v", m.ID, k, deref(got), deref(w[k]))
			}
		}
	}
	data, err := json.Marshal(rep.Items[0])
	if err != nil {
		t.Fatal(err)
	}
	if keys := `"cost_of_delay":700,"cost_of_delay_incurred":525`; !strings.Contains(string(data), keys) {
		t.Errorf("S-0001 json lacks %s: %s", keys, data)
	}
	data, err = json.Marshal(rep.Items[2])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "cost_of_delay") {
		t.Errorf("S-0003 has no value and carries a cost of delay: %s", data)
	}
}

// S-0205: what is outstanding per column at the end of each day of the
// window, what waiting cost during each day up to now, and during each week,
// the first one whole.
func TestCostOfDelayByTheDayAndTheWeek(t *testing.T) {
	data, err := json.Marshal(Compute(costItems(), Options{Now: now, Since: threeDays}).CostOfDelay)
	if err != nil {
		t.Fatal(err)
	}
	// S-0003, without a value, waits in backlog throughout (S-0213).
	oneInBacklog := `"without_value":{"backlog":1,"in-progress":0,"ready":0,"review":0},`
	want := `{"days":[` +
		`{"date":"2026-08-29","outstanding":{"backlog":0,"in-progress":0,"ready":700,"review":0},` + oneInBacklog + `"incurred":105},` +
		`{"date":"2026-08-30","outstanding":{"backlog":0,"in-progress":2100,"ready":0,"review":0},` + oneInBacklog + `"incurred":116.67},` +
		`{"date":"2026-08-31","outstanding":{"backlog":100.5,"in-progress":0,"ready":0,"review":1400},` + oneInBacklog + `"incurred":14.36},` +
		`{"date":"2026-09-01","outstanding":{"backlog":0,"in-progress":0,"ready":100.5,"review":1400},` + oneInBacklog + `"incurred":7.18}],` +
		`"weeks":[` +
		`{"week":"2026-W35","start":"2026-08-24","incurred":606.67},` +
		`{"week":"2026-W36","start":"2026-08-31","incurred":21.54}]}`
	if string(data) != want {
		t.Errorf("cost_of_delay =\n%s\nwant\n%s", data, want)
	}
}

// S-0205: with no items every day and week of the window is still there, at
// zero, every column present.
func TestCostOfDelayOfNoItemsIsZeroOverTheWindow(t *testing.T) {
	data, err := json.Marshal(Compute(nil, Options{Now: now, Since: threeDays}).CostOfDelay)
	if err != nil {
		t.Fatal(err)
	}
	day := func(date string) string {
		return `{"date":"` + date + `","outstanding":{"backlog":0,"in-progress":0,"ready":0,"review":0},` +
			`"without_value":{"backlog":0,"in-progress":0,"ready":0,"review":0},"incurred":0}`
	}
	want := `{"days":[` + day("2026-08-29") + `,` + day("2026-08-30") + `,` + day("2026-08-31") + `,` + day("2026-09-01") + `],` +
		`"weeks":[{"week":"2026-W35","start":"2026-08-24","incurred":0},{"week":"2026-W36","start":"2026-08-31","incurred":0}]}`
	if string(data) != want {
		t.Errorf("cost_of_delay =\n%s\nwant\n%s", data, want)
	}
}

// unvaluedItems are stories, an epic, and a task, most without a cost of
// delay. The window's days are 29 August to 1 September 2026.
func unvaluedItems() []*workitem.Item {
	// gains its value on the third day: its value today stands for all of its
	// history, so it counts on no day
	gains := costItem("S-0012", workitem.Story, val(700), "2026-08-28T00:00:00Z")
	gains.CostOfDelay.At = "2026-08-31T10:00:00Z"
	// has inputs and no value: it counts
	inputs := costItem("S-0013", workitem.Story, nil, "2026-08-29T20:00:00Z",
		workitem.InProgress, "2026-08-30T09:00:00Z", workitem.Done, "2026-08-31T15:00:00Z")
	inputs.CostOfDelay = &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: val(500)}}
	return []*workitem.Item{
		// moves a column a day, in a different one at the end of each
		costItem("S-0011", workitem.Story, nil, "2026-08-28T00:00:00Z",
			workitem.Ready, "2026-08-30T10:00:00Z", workitem.InProgress, "2026-08-31T09:00:00Z", workitem.Review, "2026-09-01T08:00:00Z"),
		gains,
		inputs,
		// created on the third day
		costItem("S-0014", workitem.Story, nil, "2026-08-31T13:00:00Z"),
		costItem("E-0002", workitem.Epic, nil, "2026-08-20T00:00:00Z"),
		costItem("T-0001", workitem.Task, nil, "2026-08-20T00:00:00Z", workitem.InProgress, "2026-08-30T09:00:00Z"),
	}
}

// S-0213: per column, the number of the report's items in it at the end of
// each day that have no value, every column present; a task carries no value
// and is never counted.
func TestCostOfDelayCountsTheItemsWithoutAValuePerColumn(t *testing.T) {
	for _, tc := range []struct {
		typ  string
		want string
	}{
		{workitem.Story, `[` +
			`{"backlog":2,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":0,"in-progress":1,"ready":1,"review":0},` +
			`{"backlog":1,"in-progress":1,"ready":0,"review":0},` +
			`{"backlog":1,"in-progress":0,"ready":0,"review":1}]`},
		{workitem.Epic, `[` +
			`{"backlog":1,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":1,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":1,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":1,"in-progress":0,"ready":0,"review":0}]`},
		{workitem.Task, `[` +
			`{"backlog":0,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":0,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":0,"in-progress":0,"ready":0,"review":0},` +
			`{"backlog":0,"in-progress":0,"ready":0,"review":0}]`},
	} {
		var counts []map[string]int
		for _, d := range Compute(unvaluedItems(), Options{Now: now, Since: threeDays, Type: tc.typ}).CostOfDelay.Days {
			counts = append(counts, d.WithoutValue)
		}
		data, err := json.Marshal(counts)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != tc.want {
			t.Errorf("%s without_value by day =\n%s\nwant\n%s", tc.typ, data, tc.want)
		}
	}
}

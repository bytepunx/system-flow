package metrics

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// orderAt is the projection's start, T0; the report's now is a moment past
// it, which the projection drops.
var orderAt = time.Date(2026, 10, 5, 12, 0, 0, 0, time.UTC)

// orderStory is an open story in a state with the cost of delay value and
// forecast duration given, nil and "" for none; one in progress or in
// review started an hour before T0.
func orderStory(id, status string, value *float64, duration string) *workitem.Item {
	it := &workitem.Item{ID: id, Type: workitem.Story, Status: status, Created: "2026-10-01T09:00:00Z"}
	if value != nil {
		it.CostOfDelay = &workitem.CostOfDelay{Value: value}
	}
	if duration != "" {
		it.Forecast = &workitem.Forecast{Duration: duration}
	}
	if status == workitem.InProgress || status == workitem.Review {
		started := orderAt.Add(-time.Hour)
		it.Transitions = append(it.Transitions, workitem.Transition{To: workitem.InProgress, At: started.Format(workitem.TimeFormat)})
		if status == workitem.Review {
			it.Transitions = append(it.Transitions, workitem.Transition{To: workitem.Review, At: started.Add(time.Minute).Format(workitem.TimeFormat)})
		}
	}
	return it
}

// orderHistory is three archived done stories whose agents worked a time,
// in progress 1.5 times as long and done twice as long after they started:
// every story's busy factor is 1.5 and its cycle factor 2.
func orderHistory() []*workitem.Item {
	var out []*workitem.Item
	start := time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)
	for i, seconds := range []int64{1200, 1600, 2000} {
		it := &workitem.Item{ID: "S-000" + string(rune('1'+i)), Type: workitem.Story, Status: workitem.Done, Nature: "feature", Archived: true,
			Created: start.Format(workitem.TimeFormat), Usage: &usage.Usage{Source: usage.SourceLog, Seconds: seconds}}
		it.Transitions = []workitem.Transition{
			{To: workitem.InProgress, At: start.Format(workitem.TimeFormat)},
			{To: workitem.Review, At: start.Add(time.Duration(seconds) * 1500 * time.Millisecond).Format(workitem.TimeFormat)},
			{To: workitem.Done, At: start.Add(time.Duration(seconds) * 2 * time.Second).Format(workitem.TimeFormat)},
		}
		out = append(out, it)
	}
	return out
}

// orderItems is a ready column whose pull order is the dearest of the three,
// a story in progress that holds one of the two lanes until T0+2h, one in
// review that holds none, and a backlog story that is not played out.
//
// S-0010 holds its lane from T0-1h for 3h. Each placed story holds its lane
// for 1.5 times its duration: S-0020 (cod 100, wsjf 25/h) 6h, S-0021 (cod
// 700, wsjf 2098/h) 30m1.5s, S-0022 (cod 1400, wsjf 200/h) 10h30m. S-0023
// has no value and S-0024 no forecast; an epic in ready is not a story.
func orderItems() []*workitem.Item {
	return append(orderHistory(),
		orderStory("S-0010", workitem.InProgress, nil, "2h"),
		orderStory("S-0011", workitem.Review, val(5000), "8h"),
		orderStory("S-0020", workitem.Ready, val(100), "4h"),
		orderStory("S-0021", workitem.Ready, val(700), "20m1s"),
		orderStory("S-0022", workitem.Ready, val(1400), "7h"),
		orderStory("S-0023", workitem.Ready, nil, "1h"),
		orderStory("S-0024", workitem.Ready, val(300), ""),
		orderStory("S-0030", workitem.Backlog, val(9000), "1h"),
		&workitem.Item{ID: "E-0001", Type: workitem.Epic, Status: workitem.Ready, CostOfDelay: &workitem.CostOfDelay{Value: val(50)}},
	)
}

// orderOptions are a report at T0 and a moment past, on a board whose
// in-progress limit is limit and whose pull order puts the stories left out
// among the others.
func orderOptions(limit int) Options {
	return Options{Now: orderAt.Add(400 * time.Millisecond), WIPLimit: limit,
		Order: []string{"S-0024", "S-0020", "S-0021", "S-0023", "S-0022", "S-0030", "E-0001"}}
}

func orderJSON(t *testing.T, o CostOrder) string {
	t.Helper()
	data, err := json.Marshal(o)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// S-0213: on two lanes, one free at T0 and one at T0+2h, the pull order
// pulls S-0020 at T0, S-0021 at T0+2h, and S-0022 when S-0021's lane frees
// at 14:30:01.5, rounded up to 14:30:02. Cost of delay pulls S-0022 first,
// and WSJF S-0021, whose lane frees at 12:30:01.5 for S-0022. Each cost is
// the value times the seconds from T0 to the rounded pull over a week:
// current 700×7200 + 1400×9002 = 29.17, cod 700×7200 + 100×9002 = 9.82,
// wsjf 1400×1802 + 100×7200 = 5.36, over 604800, so WSJF saves 23.81.
func TestCostOrderPricesThePullOrderAgainstTheTwoPolicies(t *testing.T) {
	got := orderJSON(t, Compute(orderItems(), orderOptions(2)).CostOfDelay.Order)
	start := `{"at":"2026-10-05T12:00:00Z","incurred":0}`
	want := `{"at":"2026-10-05T12:00:00Z","horizon":"2026-10-05T14:30:02Z","series":[` +
		`{"by":"current","total":29.17,"points":[` + start +
		`,{"at":"2026-10-05T12:00:00Z","id":"S-0020","incurred":0}` +
		`,{"at":"2026-10-05T14:00:00Z","id":"S-0021","incurred":8.33}` +
		`,{"at":"2026-10-05T14:30:02Z","id":"S-0022","incurred":29.17}]},` +
		`{"by":"cod","total":9.82,"points":[` + start +
		`,{"at":"2026-10-05T12:00:00Z","id":"S-0022","incurred":0}` +
		`,{"at":"2026-10-05T14:00:00Z","id":"S-0021","incurred":8.33}` +
		`,{"at":"2026-10-05T14:30:02Z","id":"S-0020","incurred":9.82}]},` +
		`{"by":"wsjf","total":5.36,"points":[` + start +
		`,{"at":"2026-10-05T12:00:00Z","id":"S-0021","incurred":0}` +
		`,{"at":"2026-10-05T12:30:02Z","id":"S-0022","incurred":4.17}` +
		`,{"at":"2026-10-05T14:00:00Z","id":"S-0020","incurred":5.36}]}],` +
		`"saving":23.81,"cheaper":"wsjf","left_out":["S-0024","S-0023"]}`
	if got != want {
		t.Errorf("order =\n%s\nwant\n%s", got, want)
	}
	opt := orderOptions(2)
	opt.Type = workitem.Epic
	if epics := orderJSON(t, Compute(orderItems(), opt).CostOfDelay.Order); epics != want {
		t.Errorf("an epic report projects other than the ready stories:\n%s", epics)
	}
}

// S-0213: a story waits for the delivery of one in its after that is in
// progress: S-0010, started at T0-1h with 2h and no history, is delivered at
// T0+1h, so S-0040 waits an hour on a free lane, 700×3600/604800 = 4.17.
func TestCostOrderWaitsForTheStoriesAfter(t *testing.T) {
	waiting := orderStory("S-0040", workitem.Ready, val(700), "1h")
	waiting.After = []string{"S-0010"}
	items := []*workitem.Item{orderStory("S-0010", workitem.InProgress, nil, "2h"), waiting}
	got := Compute(items, Options{Now: orderAt, WIPLimit: 3, Order: []string{"S-0040"}}).CostOfDelay.Order
	if s := got.Series[0]; len(s.Points) != 2 || s.Points[1].At != "2026-10-05T13:00:00Z" || s.Total != 4.17 {
		t.Errorf("current = %+v, want S-0040 pulled at 13:00 for 4.17", s)
	}
}

// S-0213: with no ready story placed, each series is only its first point,
// every total and the saving are 0, cod is cheaper, and the horizon is at.
func TestCostOrderOfAnEmptyReadyColumn(t *testing.T) {
	start := `{"at":"2026-10-05T12:00:00Z","incurred":0}`
	series := func(by string) string { return `{"by":"` + by + `","total":0,"points":[` + start + `]}` }
	want := `{"at":"2026-10-05T12:00:00Z","horizon":"2026-10-05T12:00:00Z","series":[` +
		series("current") + `,` + series("cod") + `,` + series("wsjf") + `],"saving":0,"cheaper":"cod","left_out":[]}`
	if got := orderJSON(t, Compute(orderHistory(), orderOptions(2)).CostOfDelay.Order); got != want {
		t.Errorf("order =\n%s\nwant\n%s", got, want)
	}
	if got := Compute(nil, Options{Now: orderAt}).CostOfDelay.Order; orderJSON(t, got) != want {
		t.Errorf("order of no items =\n%s\nwant\n%s", orderJSON(t, got), want)
	}
}

// S-0213: with no in-progress limit every story is pulled at once, at T0,
// as flai forecast pulls them, and every order costs nothing.
func TestCostOrderWithoutALimitPullsEveryStoryAtOnce(t *testing.T) {
	got := Compute(orderItems(), orderOptions(0)).CostOfDelay.Order
	if got.Horizon != "2026-10-05T12:00:00Z" || got.Saving != 0 || got.Cheaper != OrderCOD || len(got.Series) != 3 {
		t.Errorf("horizon %s, saving %v, cheaper %s, %d series", got.Horizon, got.Saving, got.Cheaper, len(got.Series))
	}
	for _, s := range got.Series {
		if len(s.Points) != 4 || s.Total != 0 {
			t.Errorf("%s: %d points, total %v; want 4 at 0", s.By, len(s.Points), s.Total)
		}
		for _, p := range s.Points {
			if p.At != "2026-10-05T12:00:00Z" || p.Incurred != 0 {
				t.Errorf("%s: point %+v, want at T0 for 0", s.By, p)
			}
		}
	}
}

package metrics

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0293: the stories' turns are summed per class over every UTC day of the
// window, from its start's day to now's, every day present; per story over
// those days, by canonical ID, leaving out a story with none on them; and
// over the window. Every story counts whatever its status and the report's
// type; an epic's summed turns and a task's do not, nor days outside the
// window.
func TestTurnsAreSummedPerDayAndStoryOverTheWindow(t *testing.T) {
	turnsOn := func(days ...usage.TurnDay) *usage.Usage {
		return &usage.Usage{Source: usage.SourceLog, Seconds: 60, Turns: days, Models: []usage.Model{{Model: "m", Input: 1}}}
	}
	items := []*workitem.Item{
		{ID: "S-0010", Type: workitem.Story, Title: "Later", Status: workitem.Cancelled, Usage: turnsOn(
			usage.TurnDay{Day: "2026-08-28", Work: 50}, // before the window's start's day
			usage.TurnDay{Day: "2026-08-29", Ceremony: 2, Work: 3},
			usage.TurnDay{Day: "2026-09-01", TestRuns: 1, HandEdits: 1, EmptyWakes: 2},
		)},
		{ID: "S-0002", Type: workitem.Story, Title: "Earlier", Status: workitem.InProgress, Usage: turnsOn(
			usage.TurnDay{Day: "2026-08-29", Work: 4},
		)},
		{ID: "S-0003", Type: workitem.Story, Title: "Outside", Status: workitem.Done, Usage: turnsOn(
			usage.TurnDay{Day: "2026-08-01", Work: 9},
			usage.TurnDay{Day: "2026-09-02", Work: 9}, // after now's day
		)},
		{ID: "S-0004", Type: workitem.Story, Title: "Unmeasured", Status: workitem.Done},
		{ID: "E-0001", Type: workitem.Epic, Title: "Epic", Status: workitem.InProgress, Usage: turnsOn(usage.TurnDay{Day: "2026-08-30", Work: 7})},
		{ID: "T-0001", Type: workitem.Task, Title: "Task", Status: workitem.Done, Usage: turnsOn(usage.TurnDay{Day: "2026-08-30", Work: 7})},
	}
	// the window starts at noon on 29 August and ends at noon on 1 September
	rep := Compute(items, Options{Now: now, Since: 3 * 24 * time.Hour, Type: workitem.Task})
	data, err := json.Marshal(rep.Turns)
	if err != nil {
		t.Fatal(err)
	}
	zero := `"turns":0,"ceremony":0,"test_runs":0,"empty_wakes":0,"hand_edits":0,"work":0`
	want := `{"classes":["ceremony","test_runs","empty_wakes","hand_edits","work"],` +
		`"total":{"turns":13,"ceremony":2,"test_runs":1,"empty_wakes":2,"hand_edits":1,"work":7},` +
		`"days":[{"day":"2026-08-29","turns":9,"ceremony":2,"test_runs":0,"empty_wakes":0,"hand_edits":0,"work":7},` +
		`{"day":"2026-08-30",` + zero + `},{"day":"2026-08-31",` + zero + `},` +
		`{"day":"2026-09-01","turns":4,"ceremony":0,"test_runs":1,"empty_wakes":2,"hand_edits":1,"work":0}],` +
		`"stories":[{"id":"S-0002","title":"Earlier","status":"in-progress","turns":4,"ceremony":0,"test_runs":0,"empty_wakes":0,"hand_edits":0,"work":4},` +
		`{"id":"S-0010","title":"Later","status":"cancelled","turns":9,"ceremony":2,"test_runs":1,"empty_wakes":2,"hand_edits":1,"work":3}]}`
	if string(data) != want {
		t.Errorf("turns =\n%s\nwant\n%s", data, want)
	}
	if got := rep.Turns.Total.Count(usage.TurnEmptyWakes); got != 2 {
		t.Errorf("total.Count(empty_wakes) = %d, want 2", got)
	}
}

// S-0293: with no story carrying turns the report still carries every day of
// the window, empty lists rather than null.
func TestTurnsAreAlwaysReported(t *testing.T) {
	data, err := json.Marshal(Compute(nil, Options{Now: now, Since: 24 * time.Hour}).Turns)
	if err != nil {
		t.Fatal(err)
	}
	zero := `"turns":0,"ceremony":0,"test_runs":0,"empty_wakes":0,"hand_edits":0,"work":0`
	want := `{"classes":["ceremony","test_runs","empty_wakes","hand_edits","work"],"total":{` + zero + `},` +
		`"days":[{"day":"2026-08-31",` + zero + `},{"day":"2026-09-01",` + zero + `}],"stories":[]}`
	if string(data) != want {
		t.Errorf("turns =\n%s\nwant\n%s", data, want)
	}
}

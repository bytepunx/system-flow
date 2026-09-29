package metrics

import (
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func doneWith(id, at string, u *usage.Usage) *workitem.Item {
	return &workitem.Item{ID: id, Type: workitem.Story, Status: workitem.Done, Created: "2026-08-20T00:00:00Z",
		Transitions: []workitem.Transition{{To: workitem.InProgress, At: "2026-08-20T00:00:00Z"}, {To: workitem.Done, At: at}}, Usage: u}
}

func spend(seconds int64, models ...usage.Model) *usage.Usage {
	return &usage.Usage{Source: usage.SourceLog, Seconds: seconds, Models: models}
}

func TestUsageIsTotalledByModelAndLaidOutAgainstTimeAndCost(t *testing.T) {
	opus := func(read int64, cost float64) usage.Model {
		return usage.Model{Model: "claude-opus-5-5", CacheRead: read, Cost: cost}
	}
	haiku := usage.Model{Model: "claude-haiku-4-5", Output: 600, Cost: 0.1}
	items := []*workitem.Item{
		doneWith("S-0002", "2026-08-30T00:00:00Z", spend(1800, opus(3000, 2), haiku)),
		doneWith("S-0001", "2026-08-25T00:00:00Z", spend(3600, opus(1000, 1))),
		doneWith("S-0003", "2026-08-26T00:00:00Z", nil),                        // no usage: left out
		doneWith("S-0004", "2026-07-01T00:00:00Z", spend(60, opus(99999, 99))), // before the window
	}
	rep := Compute(items, Options{Now: now})
	u := rep.Usage
	if u.Items != 2 || u.Tokens != 4600 || !near(u.Cost, 3.1) || u.Seconds != 5400 {
		t.Fatalf("totals = %+v", u)
	}
	if len(u.Models) != 2 || u.Models[0].Model != "claude-haiku-4-5" || u.Models[1].Items != 2 || !near(*u.Models[1].TokensPerHour, 4000/1.5) {
		t.Errorf("models = %+v", u.Models)
	}
	if len(u.Done) != 2 || u.Done[0].ID != "S-0001" || u.Done[1].Done != 2 || !near(u.Done[1].Cost, 3.1) || u.Done[1].Tokens != 4600 {
		t.Errorf("done = %+v", u.Done)
	}
	if h := u.ByModel["claude-haiku-4-5"]; len(h) != 1 || h[0].ID != "S-0002" || h[0].Done != 1 || h[0].Tokens != 600 {
		t.Errorf("by model = %+v", u.ByModel)
	}
	for _, m := range rep.Items {
		switch m.ID {
		case "S-0002":
			if m.Usage == nil || m.Usage.Tokens != 3600 || !near(*m.Usage.TokensPerHour, 7200) || len(m.Usage.Models) != 2 || !near(*m.Usage.Models[1].TokensPerHour, 1200) {
				t.Errorf("S-0002 usage = %+v", m.Usage)
			}
		case "S-0003":
			if m.Usage != nil {
				t.Errorf("an item with no usage has %+v", m.Usage)
			}
		}
	}
	empty := Compute(nil, Options{Now: now, Since: time.Hour})
	if empty.Usage.Models == nil || empty.Usage.Done == nil || empty.Usage.ByModel == nil {
		t.Error("an empty report's usage lists are null")
	}
}

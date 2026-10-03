package mcpserver

import (
	"context"
	"strings"
	"testing"
)

// S-0206, ADR-0079: a strategic agent reports an activity and gets back the
// entry flai measured and the document's totals.
func TestAStrategicAgentLogsAnActivityThroughTheServer(t *testing.T) {
	type call struct {
		root, kind, summary, items, by string
	}
	var got []call
	f := setupWith(t, func(o *Options) {
		o.Activities = func(_ context.Context, root, kind, summary string, items []string, by string) (ActivityLogged, error) {
			got = append(got, call{root, kind, summary, strings.Join(items, ","), by})
			return ActivityLogged{
				Entry:    ActivityEntry{At: "2026-10-03T10:01:30Z", Summary: summary, Items: items, Seconds: 90, Cost: 0.6173, Estimated: true},
				Activity: ActivityTotals{Kind: kind, AccruedCost: 1.6173, AccruedSeconds: 390, TasksCompleted: 2, LastRun: "2026-10-03T10:01:30Z"},
			}, nil
		}
	})
	out, failed := f.call(t, "activity_log", map[string]any{"kind": "orchestrator", "summary": "  Ordered\n the board ", "items": []string{"s-1", "S-0002"}})
	if failed != "" {
		t.Fatal(failed)
	}
	entry, _ := out["entry"].(map[string]any)
	totals, _ := out["activity"].(map[string]any)
	if entry["seconds"] != float64(90) || entry["cost"] != 0.6173 || entry["estimated"] != true || entry["summary"] != "Ordered the board" {
		t.Errorf("entry: %v", entry)
	}
	if totals["kind"] != "orchestrator" || totals["accrued_cost"] != 1.6173 || totals["accrued_seconds"] != float64(390) || totals["tasks_completed"] != float64(2) || totals["last_run"] != "2026-10-03T10:01:30Z" {
		t.Errorf("totals: %v", totals)
	}
	want := call{projectRoot(f.repo), "orchestrator", "Ordered the board", "S-0001,S-0002", "claude"}
	if len(got) != 1 || got[0] != want {
		t.Errorf("asked: %+v, want %+v", got, want)
	}
}

func TestAnActivityWithoutAKindOrASummaryIsRefused(t *testing.T) {
	called := false
	f := setupWith(t, func(o *Options) {
		o.Activities = func(context.Context, string, string, string, []string, string) (ActivityLogged, error) {
			called = true
			return ActivityLogged{}, nil
		}
	})
	if _, failed := f.call(t, "activity_log", map[string]any{"kind": "story", "summary": "did it"}); !strings.Contains(failed, "planner, orchestrator, or analyzer") {
		t.Errorf("unknown kind: %q", failed)
	}
	if _, failed := f.call(t, "activity_log", map[string]any{"kind": "planner", "summary": " \n "}); !strings.Contains(failed, "needs a summary") {
		t.Errorf("empty summary: %q", failed)
	}
	if called {
		t.Error("a refused activity reached the host")
	}
}

func TestAServerThatCannotLogActivitiesSaysSo(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "activity_log", map[string]any{"kind": "planner", "summary": "Planned S-0001"}); !strings.Contains(failed, "cannot log activities") {
		t.Errorf("without a logger: %q", failed)
	}
}

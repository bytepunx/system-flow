package metrics

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// timersNow is the moment the timer fixtures are computed at.
var timersNow = time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)

// timersFlow is the baseline with three timers of its own: the wait from the
// first review to acceptance, the sum of in-progress and review counting
// cancelled items too, and one only tasks take.
func timersFlow() workflow.Workflow {
	flow := workflow.Baseline()
	flow.Timers = append(flow.Timers,
		workflow.Timer{Name: "acceptance-wait", Title: "Acceptance wait",
			From: &workflow.Point{Edge: workflow.EdgeFirst, State: workitem.Review}, To: &workflow.Point{At: workflow.PointCompleted},
			Sum: []string{}, Types: slices.Clone(workflow.Types), Counts: workflow.CountsAccepted},
		workflow.Timer{Name: "in-work", Title: "In work", Sum: []string{workitem.InProgress, workitem.Review},
			Types: slices.Clone(workflow.Types), Counts: workflow.CountsClosed},
		workflow.Timer{Name: "task-only", Title: "Task only", Sum: []string{workitem.InProgress},
			Types: []string{workflow.TypeTask}, Counts: workflow.CountsAccepted},
	)
	return flow
}

// timersItems are four stories with known transitions: one sent back from
// review once, one done months ago, one cancelled from review, and one open
// in review.
func timersItems() []*workitem.Item {
	tr := func(to, at string) workitem.Transition { return workitem.Transition{To: to, At: at} }
	story := func(id, status, created string, steps ...workitem.Transition) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Nature: "feature", Status: status, Created: created, Transitions: steps}
	}
	return []*workitem.Item{
		story("S-0001", workitem.Done, "2026-08-20T00:00:00Z",
			tr(workitem.Ready, "2026-08-20T01:00:00Z"),
			tr(workitem.InProgress, "2026-08-20T02:00:00Z"),
			tr(workitem.Review, "2026-08-20T06:00:00Z"),
			tr(workitem.InProgress, "2026-08-20T08:00:00Z"),
			tr(workitem.Review, "2026-08-20T10:00:00Z"),
			tr(workitem.Done, "2026-08-20T13:00:00Z")),
		story("S-0002", workitem.Done, "2026-01-01T00:00:00Z",
			tr(workitem.Ready, "2026-01-01T01:00:00Z"),
			tr(workitem.InProgress, "2026-01-02T00:00:00Z"),
			tr(workitem.Review, "2026-01-03T00:00:00Z"),
			tr(workitem.Done, "2026-01-05T00:00:00Z")),
		story("S-0003", workitem.Cancelled, "2026-08-25T00:00:00Z",
			tr(workitem.Ready, "2026-08-25T01:00:00Z"),
			tr(workitem.InProgress, "2026-08-25T02:00:00Z"),
			tr(workitem.Review, "2026-08-25T05:00:00Z"),
			tr(workitem.Cancelled, "2026-08-25T06:00:00Z")),
		story("S-0004", workitem.Review, "2026-09-01T00:00:00Z",
			tr(workitem.Ready, "2026-09-01T01:00:00Z"),
			tr(workitem.InProgress, "2026-09-01T02:00:00Z"),
			tr(workitem.Review, "2026-09-01T06:00:00Z")),
	}
}

// itemsByID are the report's items by their IDs.
func itemsByID(rep *Report) map[string]ItemMetrics {
	out := map[string]ItemMetrics{}
	for _, m := range rep.Items {
		out[m.ID] = m
	}
	return out
}

// S-0752 criteria 1, 3, and 6, ADR-0317 §2 and §3: each item holds the
// workflow's own timers that take its type and have a value for it, read
// from the transitions it already records, a story done months ago among
// them.
func TestComputeReportsEachItemsOwnTimers(t *testing.T) {
	flow := timersFlow()
	rep := Compute(timersItems(), Options{Now: timersNow, Since: 365 * 24 * time.Hour, Workflow: &flow})
	names := make([]string, len(rep.Timers))
	for i, d := range rep.Timers {
		names[i] = d.Name
	}
	if !slices.Equal(names, []string{"acceptance-wait", "in-work", "task-only"}) {
		t.Errorf("timers %v, want the three own timers in file order and none of the baseline's four", names)
	}
	if d := rep.Timers[1]; d.Title != "In work" || d.Counts != workflow.CountsClosed || !slices.Equal(d.Types, workflow.Types) {
		t.Errorf("in-work definition %+v", d)
	}
	items := itemsByID(rep)
	for id, want := range map[string]map[string]float64{
		// Done 13:00 less the first review at 06:00; in-progress 4h and 2h,
		// review 2h and 3h.
		"S-0001": {"acceptance-wait": 7 * 3600, "in-work": 11 * 3600},
		// Done two days after the first review; a day in progress, two in
		// review.
		"S-0002": {"acceptance-wait": 2 * 86400, "in-work": 3 * 86400},
		// Cancelled an hour after review: the point completed is the close.
		"S-0003": {"acceptance-wait": 3600, "in-work": 4 * 3600},
		// Open in review: no completion, and the open review to now counts.
		"S-0004": {"in-work": 4*3600 + 6*3600},
	} {
		if got := items[id].Timers; !mapsEqual(got, want) {
			t.Errorf("%s timers %v, want %v", id, got, want)
		}
	}
}

// S-0752 criteria 1 and 6, ADR-0317 §4: the summary gives each own timer
// that takes stories its distribution over the items completed in the
// window that its counts takes, the cancelled one only for counts closed,
// and the open one in none.
func TestComputeSummarisesOwnTimersOverTheWindow(t *testing.T) {
	flow := timersFlow()
	rep := Compute(timersItems(), Options{Now: timersNow, Since: 365 * 24 * time.Hour, Workflow: &flow, By: "nature"})
	want := map[string]Distribution{
		"acceptance-wait": {Count: 2, P50: 7 * 3600, P85: 2 * 86400, Max: 2 * 86400, Mean: (7*3600 + 2*86400) / 2},
		"in-work":         {Count: 3, P50: 11 * 3600, P85: 3 * 86400, Max: 3 * 86400, Mean: (11*3600 + 3*86400 + 4*3600) / 3},
	}
	if got := rep.Summary.Timers; !distributionsEqual(got, want) {
		t.Errorf("summary timers %+v, want %+v", got, want)
	}
	if len(rep.Groups) != 1 || !distributionsEqual(rep.Groups[0].Timers, want) {
		t.Errorf("groups %+v, want one feature group with the summary's timers", rep.Groups)
	}
	if rep.Summary.Completed != 2 || rep.Summary.Cancelled != 1 || rep.Summary.CycleTime.Count != 2 {
		t.Errorf("summary %+v: the baseline's figures changed", rep.Summary)
	}
}

// ADR-0317 §4: a timer that takes the type with no item completed in the
// window gives a distribution of zeros.
func TestComputeGivesZerosForAnOwnTimerWithNothingToCount(t *testing.T) {
	flow := timersFlow()
	rep := Compute(timersItems(), Options{Now: timersNow, Since: time.Hour, Workflow: &flow})
	want := map[string]Distribution{"acceptance-wait": {}, "in-work": {}}
	if got := rep.Summary.Timers; !distributionsEqual(got, want) {
		t.Errorf("summary timers %+v, want zeros for both", got)
	}
}

// S-0752 criterion 5, ADR-0317 §1: on the baseline workflow the report,
// its items, and its summary hold no timers, and its JSON no timers key.
func TestComputeOnTheBaselineHasNoTimers(t *testing.T) {
	rep := Compute(timersItems(), Options{Now: timersNow, Since: 365 * 24 * time.Hour, By: "nature"})
	if rep.Timers != nil || rep.Summary.Timers != nil {
		t.Errorf("report timers %v, summary timers %v, want none", rep.Timers, rep.Summary.Timers)
	}
	for _, m := range rep.Items {
		if m.Timers != nil {
			t.Errorf("%s timers %v, want none", m.ID, m.Timers)
		}
	}
	for _, g := range rep.Groups {
		if g.Timers != nil {
			t.Errorf("group %s timers %v, want none", g.Group, g.Timers)
		}
	}
	b, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"timers"`) {
		t.Errorf("the baseline's report has a timers key")
	}
}

func mapsEqual(a, b map[string]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range b {
		if got, ok := a[k]; !ok || !near(got, v) {
			return false
		}
	}
	return true
}

func distributionsEqual(a, b map[string]Distribution) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range b {
		got, ok := a[k]
		if !ok || got.Count != v.Count || !near(got.P50, v.P50) || !near(got.P85, v.P85) || !near(got.Max, v.Max) || !near(got.Mean, v.Mean) {
			return false
		}
	}
	return true
}

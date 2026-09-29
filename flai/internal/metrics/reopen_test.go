package metrics

import (
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0167, ADR-0055: an item moved back out of cancelled is open again, so it
// has no completed, lead time, or cycle time until it closes again, and then
// it completes when it last closed.
func TestReopenedItemIsNotCompleted(t *testing.T) {
	t0 := time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)
	at := func(h int) string { return t0.Add(time.Duration(h) * time.Hour).Format(workitem.TimeFormat) }
	it := &workitem.Item{
		ID: "S-0001", Type: workitem.Story, Nature: "feature", Status: workitem.Backlog, Created: t0.Format(workitem.TimeFormat),
		Transitions: []workitem.Transition{
			{To: workitem.Ready, At: at(1)},
			{To: workitem.InProgress, At: at(2)},
			{To: workitem.Cancelled, At: at(3)},
			{To: workitem.Backlog, At: at(4)},
		},
	}
	now := t0.Add(10 * time.Hour)
	m := Derive(it, now)
	if m.Completed != "" || m.LeadTime != nil || m.CycleTime != nil {
		t.Errorf("reopened: completed %q lead %v cycle %v", m.Completed, m.LeadTime, m.CycleTime)
	}
	if m.Age == nil || *m.Age != (8*time.Hour).Seconds() {
		t.Errorf("reopened item's age counts from its first start: %v", m.Age)
	}
	rep := Compute([]*workitem.Item{it}, Options{Now: now})
	if rep.Summary.Cancelled != 0 || rep.Summary.Completed != 0 {
		t.Errorf("summary: %+v", rep.Summary)
	}

	it.Transitions = append(it.Transitions, workitem.Transition{To: workitem.Cancelled, At: at(6)})
	it.Status = workitem.Cancelled
	m = Derive(it, now)
	if m.Completed != at(6) || m.LeadTime == nil || *m.LeadTime != (6*time.Hour).Seconds() {
		t.Errorf("closed again: completed %q lead %v", m.Completed, m.LeadTime)
	}
}

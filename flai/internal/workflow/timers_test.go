package workflow

import (
	"maps"
	"strings"
	"testing"
	"time"
)

// when is the RFC 3339 time s.
func when(t *testing.T, s string) time.Time {
	t.Helper()
	at, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("time %q: %v", s, err)
	}
	return at
}

// history is a story created at created whose steps are "<state> <at>".
func history(t *testing.T, created string, steps ...string) History {
	t.Helper()
	h := History{Type: TypeStory, Created: when(t, created)}
	for _, s := range steps {
		to, at, ok := strings.Cut(s, " ")
		if !ok {
			t.Fatalf("step %q is not \"<state> <at>\"", s)
		}
		h.Steps = append(h.Steps, Step{To: to, At: when(t, at)})
	}
	return h
}

// measures are the values of the named timers over h, a timer without a
// value left out.
func measures(t *testing.T, w Workflow, h History, now time.Time, names ...string) map[string]float64 {
	t.Helper()
	out := map[string]float64{}
	for _, n := range names {
		tm, ok := w.Timer(n)
		if !ok {
			t.Fatalf("the workflow has no timer %s", n)
		}
		if v, ok := w.Measure(tm, h, now); ok {
			out[n] = v
		}
	}
	return out
}

// S-0750: the baseline's four timers give today's lead, cycle, queue, and
// review time, and InState today's time in state, on the fixture items of
// internal/metrics/testdata/good and on the edge cases ADR-0305 names.
func TestBaselineTimersReadHistories(t *testing.T) {
	w := Baseline()
	for _, c := range []struct {
		name    string
		h       History
		now     string
		timers  map[string]float64
		inState map[string]float64
	}{
		{
			"S-001 done",
			history(t, "2026-08-01T09:00:00Z", "ready 2026-08-01T10:00:00Z", "in-progress 2026-08-02T10:00:00Z", "review 2026-08-03T10:00:00Z", "done 2026-08-03T12:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 183600, "cycle": 93600, "queue": 86400, "review": 7200},
			map[string]float64{"backlog": 3600, "ready": 86400, "in-progress": 86400, "review": 7200},
		},
		{
			"S-002 done",
			history(t, "2026-08-10T09:00:00Z", "ready 2026-08-10T09:30:00Z", "in-progress 2026-08-11T09:30:00Z", "review 2026-08-11T19:30:00Z", "done 2026-08-12T09:30:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 174600, "cycle": 86400, "queue": 86400, "review": 50400},
			map[string]float64{"backlog": 1800, "ready": 86400, "in-progress": 36000, "review": 50400},
		},
		{
			"T-001 done and never in review",
			history(t, "2026-08-01T09:30:00Z", "ready 2026-08-01T10:00:00Z", "in-progress 2026-08-02T10:00:00Z", "done 2026-08-03T09:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 171000, "cycle": 82800, "queue": 86400},
			map[string]float64{"backlog": 1800, "ready": 86400, "in-progress": 82800},
		},
		{
			"S-003 cancelled from backlog",
			history(t, "2026-08-20T09:00:00Z", "cancelled 2026-08-21T09:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 86400},
			map[string]float64{"backlog": 86400},
		},
		{
			"S-004 open in progress",
			history(t, "2026-08-25T09:00:00Z", "ready 2026-08-25T10:00:00Z", "in-progress 2026-08-31T10:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"queue": 518400},
			map[string]float64{"backlog": 3600, "ready": 518400, "in-progress": 86400},
		},
		{
			"open in review counts review to now",
			history(t, "2026-08-01T09:00:00Z", "ready 2026-08-01T10:00:00Z", "in-progress 2026-08-02T10:00:00Z", "review 2026-08-03T10:00:00Z"),
			"2026-08-03T15:00:00Z",
			map[string]float64{"queue": 86400, "review": 18000},
			map[string]float64{"backlog": 3600, "ready": 86400, "in-progress": 86400, "review": 18000},
		},
		{
			"reopened from cancelled is open",
			history(t, "2026-08-20T09:00:00Z", "cancelled 2026-08-21T09:00:00Z", "backlog 2026-08-22T09:00:00Z"),
			"2026-08-23T09:00:00Z",
			map[string]float64{},
			map[string]float64{"backlog": 172800, "cancelled": 86400},
		},
		{
			"sent back to backlog and to ready: queue is the ready time before in-progress",
			history(t, "2026-08-01T00:00:00Z", "ready 2026-08-01T01:00:00Z", "backlog 2026-08-01T03:00:00Z", "ready 2026-08-01T07:00:00Z", "in-progress 2026-08-01T08:00:00Z", "ready 2026-08-01T09:00:00Z", "in-progress 2026-08-01T11:00:00Z", "review 2026-08-01T12:00:00Z", "done 2026-08-01T13:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 46800, "cycle": 18000, "queue": 10800, "review": 3600},
			map[string]float64{"backlog": 18000, "ready": 18000, "in-progress": 7200, "review": 3600},
		},
		{
			"in progress before ever ready has no queue",
			history(t, "2026-08-01T00:00:00Z", "in-progress 2026-08-01T01:00:00Z", "ready 2026-08-01T02:00:00Z", "in-progress 2026-08-01T03:00:00Z", "review 2026-08-01T04:00:00Z", "done 2026-08-01T05:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{"lead": 18000, "cycle": 14400, "review": 3600},
			map[string]float64{"backlog": 3600, "in-progress": 7200, "ready": 3600, "review": 3600},
		},
		{
			"never moved",
			history(t, "2026-08-31T10:00:00Z"),
			"2026-09-01T10:00:00Z",
			map[string]float64{},
			map[string]float64{"backlog": 86400},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			now := when(t, c.now)
			if got := measures(t, w, c.h, now, "lead", "cycle", "queue", "review"); !maps.Equal(got, c.timers) {
				t.Errorf("timers = %v, want %v", got, c.timers)
			}
			if got := w.InState(c.h, now); !maps.Equal(got, c.inState) {
				t.Errorf("InState = %v, want %v", got, c.inState)
			}
		})
	}
}

// S-0750: the last interval runs to now and is open while the item is open,
// and a closed item's is left out (ADR-0305 §1); completed is the last
// step's at only while the item is closed.
func TestIntervalsAndCompleted(t *testing.T) {
	w := Baseline()
	now := when(t, "2026-09-01T10:00:00Z")
	open := history(t, "2026-08-25T09:00:00Z", "ready 2026-08-25T10:00:00Z")
	want := []Interval{
		{State: "backlog", From: when(t, "2026-08-25T09:00:00Z"), To: when(t, "2026-08-25T10:00:00Z")},
		{State: "ready", From: when(t, "2026-08-25T10:00:00Z"), To: now, Open: true},
	}
	if got := w.Intervals(open, now); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("open Intervals = %v, want %v", got, want)
	}
	cancelled := history(t, "2026-08-20T09:00:00Z", "cancelled 2026-08-21T09:00:00Z")
	if got := w.Intervals(cancelled, now); len(got) != 1 || got[0].State != "backlog" || got[0].Open {
		t.Errorf("cancelled Intervals = %v, want the backlog interval alone", got)
	}
	reopened := history(t, "2026-08-20T09:00:00Z", "cancelled 2026-08-21T09:00:00Z", "backlog 2026-08-22T09:00:00Z")
	completed := Point{At: PointCompleted}
	for _, c := range []struct {
		name   string
		h      History
		closed bool
		at     string
	}{
		{"open", open, false, ""},
		{"cancelled", cancelled, true, "2026-08-21T09:00:00Z"},
		{"reopened", reopened, false, ""},
		{"never moved", history(t, "2026-08-31T10:00:00Z"), false, ""},
	} {
		if got := w.Closed(w.Current(c.h)); got != c.closed {
			t.Errorf("%s: Closed = %v, want %v", c.name, got, c.closed)
		}
		at, ok := w.PointAt(completed, c.h)
		if ok != (c.at != "") || (ok && !at.Equal(when(t, c.at))) {
			t.Errorf("%s: completed = %v, %v, want %q", c.name, at, ok, c.at)
		}
	}
	if _, ok := w.Timer("no-such-timer"); ok {
		t.Error("Timer(no-such-timer) found a timer the baseline lacks")
	}
}

// timedFlow is the baseline with in-progress renamed doing, keeping its old
// name as an alias, and testing, of the active kind, added between it and
// review, with timers that name doing and one that names the active kind.
const timedFlow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: doing, kind: active, limit: true, aliases: [in-progress]}
  - {name: testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: doing, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - {from: doing, to: testing}
  - {from: testing, to: review}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: review, to: doing, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: [doing, testing], to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, doing, testing], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: doing}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: doing}}
  - {name: review, title: Review time, sum: [review]}
  - {name: active, title: Active time, from: {first: {kind: active}}, to: completed}
---
`

// S-0750 criterion 5: on a workflow with in-progress renamed and testing
// added between it and review, cycle time spans the time in testing, time
// in state has a testing key, queue time runs until the first doing, and an
// item recorded under the old name reads under the new one.
func TestTimersOnARenamedAndAddedState(t *testing.T) {
	w := mustParse(t, timedFlow)
	now := when(t, "2026-09-01T10:00:00Z")
	timers := map[string]float64{"lead": 36000, "cycle": 28800, "queue": 3600, "review": 3600, "active": 28800}
	inState := map[string]float64{"backlog": 3600, "ready": 3600, "doing": 10800, "testing": 14400, "review": 3600}
	for _, c := range []struct {
		name  string
		doing string
	}{
		{"current name", "doing"},
		{"old name, an alias", "in-progress"},
	} {
		t.Run(c.name, func(t *testing.T) {
			h := history(t, "2026-08-01T00:00:00Z", "ready 2026-08-01T01:00:00Z", c.doing+" 2026-08-01T02:00:00Z", "testing 2026-08-01T05:00:00Z", "review 2026-08-01T09:00:00Z", "done 2026-08-01T10:00:00Z")
			got := measures(t, w, h, now, "lead", "cycle", "queue", "review", "active")
			if !maps.Equal(got, timers) {
				t.Errorf("timers = %v, want %v", got, timers)
			}
			in := w.InState(h, now)
			if !maps.Equal(in, inState) {
				t.Errorf("InState = %v, want %v", in, inState)
			}
			if spans := in["doing"] + in["testing"]; got["cycle"] < spans {
				t.Errorf("cycle %v does not span doing and testing, %v", got["cycle"], spans)
			}
			last, ok := w.PointAt(Point{Edge: EdgeLast, Kind: KindActive}, h)
			if !ok || !last.Equal(when(t, "2026-08-01T05:00:00Z")) {
				t.Errorf("{last: {kind: active}} = %v, %v, want the move to testing", last, ok)
			}
		})
	}
}

package workflow

import (
	"slices"
	"time"
)

// History is an item's created and its transitions in order, as a timer
// reads it (ADR-0305 §1).
type History struct {
	// Type is the item's type, epic, story, or task.
	Type    string
	Created time.Time
	Steps   []Step
}

// Step is one transition of a History: the state it entered and when.
type Step struct {
	To string
	At time.Time
}

// Interval is a stretch of an item's history spent in one state; Open marks
// the last one, which runs to now.
type Interval struct {
	State    string
	From, To time.Time
	Open     bool
}

// Timer is the timer named name.
func (w Workflow) Timer(name string) (Timer, bool) {
	for _, t := range w.Timers {
		if t.Name == name {
			return t, true
		}
	}
	return Timer{}, false
}

// Current is the state the history leaves the item in, by its current name:
// the last step's, or the initial state without a step. The item is closed
// when Closed reports it of this state.
func (w Workflow) Current(h History) string {
	if n := len(h.Steps); n > 0 {
		return w.canonical(h.Steps[n-1].To)
	}
	initial, _ := w.Initial()
	return initial.Name
}

// closed reports whether the history leaves the item closed.
func (w Workflow) closed(h History) bool { return w.Closed(w.Current(h)) }

// Intervals are the stretches the item spent in each state, in history
// order, under each state's current name (ADR-0305 §1, ADR-0301 §6): the
// initial state from created to the first step, each step's state to the
// next, and the last step's state to now, which a closed item leaves out.
func (w Workflow) Intervals(h History, now time.Time) []Interval {
	initial, _ := w.Initial()
	out := make([]Interval, 0, len(h.Steps)+1)
	state, from := initial.Name, h.Created
	for _, s := range h.Steps {
		out = append(out, Interval{State: state, From: from, To: s.At})
		state, from = w.canonical(s.To), s.At
	}
	if !w.closed(h) {
		out = append(out, Interval{State: state, From: from, To: now, Open: true})
	}
	return out
}

// PointAt is the moment of the history the point names, and whether the
// history has it (ADR-0305 §2).
func (w Workflow) PointAt(p Point, h History) (time.Time, bool) {
	switch p.At {
	case PointCreated:
		return h.Created, true
	case PointCompleted:
		if len(h.Steps) == 0 || !w.closed(h) {
			return time.Time{}, false
		}
		return h.Steps[len(h.Steps)-1].At, true
	}
	into := func(to string) bool {
		if p.Kind != "" {
			return w.Kind(to) == p.Kind
		}
		return w.canonical(to) == w.canonical(p.State)
	}
	var at time.Time
	found := false
	for _, s := range h.Steps {
		if !into(s.To) {
			continue
		}
		switch p.Edge {
		case EdgeFirst:
			return s.At, true
		case EdgeLast:
			at, found = s.At, true
		}
	}
	return at, found
}

// Measure is the timer's value over the history in seconds, and whether it
// has one (ADR-0305 §3). It reads the timer's form only: which items its
// types and counts take is the caller's to decide.
func (w Workflow) Measure(t Timer, h History, now time.Time) (float64, bool) {
	switch {
	case t.From != nil && t.To != nil:
		from, ok := w.PointAt(*t.From, h)
		if !ok {
			return 0, false
		}
		to, ok := w.PointAt(*t.To, h)
		if !ok {
			return 0, false
		}
		return to.Sub(from).Seconds(), true
	case len(t.Sum) == 0:
		return 0, false
	}
	entered, ok := w.entered(t.Sum, h)
	if !ok {
		return 0, false
	}
	in := w.among(t.Sum)
	if t.Until == nil {
		var sum float64
		for _, iv := range w.Intervals(h, now) {
			if in(iv.State) {
				sum += iv.To.Sub(iv.From).Seconds()
			}
		}
		return sum, true
	}
	until, ok := w.PointAt(*t.Until, h)
	if !ok || entered.After(until) {
		return 0, false
	}
	var sum time.Duration
	for _, iv := range w.Intervals(h, now) {
		if in(iv.State) && !iv.Open && !iv.To.After(until) {
			sum += iv.To.Sub(iv.From)
		}
	}
	return sum.Seconds(), true
}

// InState is the seconds the item spent in each state it entered, by the
// state's current name, the open interval included.
func (w Workflow) InState(h History, now time.Time) map[string]float64 {
	out := map[string]float64{}
	for _, iv := range w.Intervals(h, now) {
		out[iv.State] += iv.To.Sub(iv.From).Seconds()
	}
	return out
}

// entered is when the history first enters one of the states: created for
// an item created in one, else the first step into one.
func (w Workflow) entered(states []string, h History) (time.Time, bool) {
	in := w.among(states)
	if initial, ok := w.Initial(); ok && in(initial.Name) {
		return h.Created, true
	}
	for _, s := range h.Steps {
		if in(s.To) {
			return s.At, true
		}
	}
	return time.Time{}, false
}

// among reports whether a state, or a name an alias gives it, is one of the
// states.
func (w Workflow) among(states []string) func(string) bool {
	return func(name string) bool {
		name = w.canonical(name)
		return slices.ContainsFunc(states, func(s string) bool { return w.canonical(s) == name })
	}
}

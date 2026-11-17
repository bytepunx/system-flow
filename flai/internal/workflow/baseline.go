package workflow

import (
	_ "embed"
	"slices"
	"sync"
)

// baselineFile is the baseline workflow file flai carries, today's board
// (ADR-0301 §7). A project with no workflow file of its own runs it.
//
//go:embed baseline/board.md
var baselineFile string

// BaselinePath is the path the baseline is parsed at, and its Problems and
// Path name: the embedded file's, not a project's.
const BaselinePath = "flai/internal/workflow/baseline/board.md"

// baseline is the embedded file parsed once.
var baseline = sync.OnceValue(func() Workflow {
	w, p := Parse(BaselinePath, baselineFile)
	if p != nil {
		panic("the embedded baseline workflow does not parse: " + p.String())
	}
	return w
})

// Baseline is the baseline workflow flai carries, parsed once. Each call
// gives a copy of its own, so a caller may change it without changing what
// the next caller reads.
func Baseline() Workflow {
	return baseline().Clone()
}

// Clone is a copy of w that shares no slice or point with it.
func (w Workflow) Clone() Workflow {
	out := w
	out.States = make([]State, len(w.States))
	for i, s := range w.States {
		s.Types, s.Aliases = slices.Clone(s.Types), slices.Clone(s.Aliases)
		if s.Runs != nil {
			r := *s.Runs
			r.Suites, r.Kluster = cloneSuites(r.Suites), slices.Clone(r.Kluster)
			s.Runs = &r
		}
		out.States[i] = s
	}
	out.Transitions = make([]Transition, len(w.Transitions))
	for i, t := range w.Transitions {
		t.From, t.Types, t.By = slices.Clone(t.From), slices.Clone(t.Types), slices.Clone(t.By)
		t.Rules, t.Effects = slices.Clone(t.Rules), slices.Clone(t.Effects)
		if t.Verify != nil {
			v := *t.Verify
			v.Suites.Names = slices.Clone(v.Suites.Names)
			t.Verify = &v
		}
		out.Transitions[i] = t
	}
	out.Timers = make([]Timer, len(w.Timers))
	for i, t := range w.Timers {
		t.From, t.To, t.Until = clonePoint(t.From), clonePoint(t.To), clonePoint(t.Until)
		t.Sum, t.Types = slices.Clone(t.Sum), slices.Clone(t.Types)
		out.Timers[i] = t
	}
	v := w.VerifyAt
	out.VerifyAt = VerifyAt{Batch: cloneSuites(v.Batch), Publish: cloneSuites(v.Publish), CI: cloneSuites(v.CI)}
	return out
}

func clonePoint(p *Point) *Point {
	if p == nil {
		return nil
	}
	c := *p
	return &c
}

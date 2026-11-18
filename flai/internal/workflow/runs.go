package workflow

import "slices"

// The methods here read what a state runs when a story enters it from the
// state's runs (ADR-0312 §4), for flai serve to read in place of a state's
// name or kind.

// StateRuns is what flai runs when a story enters state, named by its name
// or an alias, nil for nothing. It is the state's runs, or, in a file that
// gives none of ADR-0312's keys (UsesSuiteKeys), the baseline's review's on
// each state a story's acceptance leaves from, so that such a file runs what
// it ran before (ADR-0312 §8, TH-1289); those runs have no Line, as the file
// gives no runs key. It is a copy: a change to it leaves w as it was.
func (w Workflow) StateRuns(state string) *Runs {
	s, ok := w.State(state)
	switch {
	case !ok:
		return nil
	case w.UsesSuiteKeys():
		return cloneRuns(s.Runs)
	case !slices.Contains(w.AcceptFrom(TypeStory), s.Name):
		return nil
	}
	b := baseline()
	for _, name := range b.AcceptFrom(TypeStory) {
		if bs, _ := b.State(name); bs.Runs != nil {
			r := cloneRuns(bs.Runs)
			r.Line = 0
			return r
		}
	}
	return nil
}

// KlusterStates are the states, in board order, whose StateRuns names a
// cluster step: where a story entering deploys to the cluster or runs its
// checks.
func (w Workflow) KlusterStates() []string {
	var out []string
	for _, s := range w.States {
		if r := w.StateRuns(s.Name); r != nil && len(r.Kluster) > 0 {
			out = append(out, s.Name)
		}
	}
	return out
}

// SuiteNames are the suites r runs: all for every suite, or else the list
// names; neither when r runs no suite, as a nil r runs none.
func (r *Runs) SuiteNames() (names []string, all bool) {
	if r == nil || r.Suites == nil {
		return nil, false
	}
	return slices.Clone(r.Suites.Names), r.Suites.All
}

// RunsStep reports whether r names the cluster step step, KlusterDeploy or
// KlusterChecks; a nil r names none.
func (r *Runs) RunsStep(step string) bool {
	return r != nil && slices.Contains(r.Kluster, step)
}

// cloneRuns is a copy of r that shares no list with it.
func cloneRuns(r *Runs) *Runs {
	if r == nil {
		return nil
	}
	c := *r
	c.Suites, c.Kluster = cloneSuites(r.Suites), slices.Clone(r.Kluster)
	return &c
}

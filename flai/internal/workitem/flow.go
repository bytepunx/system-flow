package workitem

import "github.com/bytepunx/system-flow/flai/internal/workflow"

// The helpers here read the board's roles from a workflow, so that the code
// that lays out the board, holds the pull, and follows an epic reads a
// state's kind, a transition's effect, a WIP limit, or an order, and never a
// state's name (ADR-0301 §2, §3). A role the workflow does not fill is the
// empty name, which no item's status is.

// pullSource is the state the pull takes a story from: the source of the
// story transition with the pull effect (ADR-0301 §3), the first it names
// when it names several.
func pullSource(w workflow.Workflow) string {
	t, ok := w.WithEffect(workflow.TypeStory, workflow.EffectPull)
	if !ok || len(t.From) == 0 {
		return ""
	}
	return t.From[0]
}

// pullTarget is the state the pull puts a story in, the active state its
// work starts in (ADR-0301 §3).
func pullTarget(w workflow.Workflow) string {
	t, _ := w.WithEffect(workflow.TypeStory, workflow.EffectPull)
	return t.To
}

// acceptSource is the waiting state the acceptance takes a story from: the
// source of the story transition with the accept effect (ADR-0301 §3), the
// first it names when it names several.
func acceptSource(w workflow.Workflow) string {
	t, ok := w.WithEffect(workflow.TypeStory, workflow.EffectAccept)
	if !ok || len(t.From) == 0 {
		return ""
	}
	return t.From[0]
}

// acceptedState is the workflow's state of the accepted kind, where a
// delivered item closes (ADR-0301 §2).
func acceptedState(w workflow.Workflow) string { return firstOfKind(w, workflow.KindAccepted) }

// cancelledState is the workflow's state of the cancelled kind, where an
// item closes undelivered (ADR-0301 §2).
func cancelledState(w workflow.Workflow) string { return firstOfKind(w, workflow.KindCancelled) }

// initialState is the state an item is created in, the workflow's initial
// queue state (ADR-0301 §2).
func initialState(w workflow.Workflow) string {
	s, _ := w.Initial()
	return s.Name
}

// orderedStates are the names of the states that read in the pull order, in
// board order (ADR-0301 §2).
func orderedStates(w workflow.Workflow) []string {
	var out []string
	for _, s := range w.States {
		if s.Order {
			out = append(out, s.Name)
		}
	}
	return out
}

// limitedStates are the states that take a WIP limit, in board order
// (ADR-0301 §2).
func limitedStates(w workflow.Workflow) []workflow.State {
	var out []workflow.State
	for _, s := range w.States {
		if s.Limit {
			out = append(out, s)
		}
	}
	return out
}

// stateNames are the names of states, in the order given.
func stateNames(states []workflow.State) []string {
	out := make([]string, len(states))
	for i, s := range states {
		out[i] = s.Name
	}
	return out
}

// closedKind reports whether a state of kind closes an item: accepted, or
// cancelled (ADR-0301 §2).
func closedKind(kind string) bool {
	return kind == workflow.KindAccepted || kind == workflow.KindCancelled
}

// startedKind reports whether an item in a state of kind has started, as an
// epic follows its stories: active, waiting, or accepted (ADR-0301 §2).
func startedKind(kind string) bool {
	return kind == workflow.KindActive || kind == workflow.KindWaiting || kind == workflow.KindAccepted
}

// underwayKind reports whether an item in a state of kind is pulled and not
// yet closed: active, or waiting (ADR-0301 §2).
func underwayKind(kind string) bool {
	return kind == workflow.KindActive || kind == workflow.KindWaiting
}

// firstOfKind is the name of the first state of kind in board order, empty
// when there is none.
func firstOfKind(w workflow.Workflow, kind string) string {
	if states := w.OfKind(kind); len(states) > 0 {
		return states[0].Name
	}
	return ""
}

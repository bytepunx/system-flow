package workflow

import "slices"

// The methods here read the roles flai serve's review passes act on, the
// kluster deploys, the forge pass, and the merge queue, from a state's kind
// and a transition's effect, never from a state's name (ADR-0301 §2, §3).

// AcceptFrom are the states the transitions of type typ with the accept
// effect leave from, in board order: where an item waits for its acceptance.
func (w Workflow) AcceptFrom(typ string) []string {
	var from []string
	for _, t := range w.Transitions {
		if t.Takes(typ) && t.HasEffect(EffectAccept) {
			from = append(from, t.From...)
		}
	}
	var out []string
	for _, s := range w.States {
		if slices.Contains(from, s.Name) {
			out = append(out, s.Name)
		}
	}
	return out
}

// AwaitsAcceptance reports whether an item of type typ in state waits for
// its acceptance: whether a transition with the accept effect leaves from
// state.
func (w Workflow) AwaitsAcceptance(typ, state string) bool {
	return state != "" && slices.Contains(w.AcceptFrom(typ), w.canonical(state))
}

// SendBack is the transition of type typ out of state from with the
// send-back effect: where an item flai sends back goes.
func (w Workflow) SendBack(typ, from string) (Transition, bool) {
	for _, t := range w.From(typ, from) {
		if t.HasEffect(EffectSendBack) {
			return t, true
		}
	}
	return Transition{}, false
}

// Gives reports whether the move of an item of type typ from one state to
// the other is a transition that gives effect.
func (w Workflow) Gives(typ, from, to, effect string) bool {
	t, ok := w.Find(typ, from, to)
	return ok && t.HasEffect(effect)
}

// Closed reports whether state closes an item: its kind is accepted or
// cancelled.
func (w Workflow) Closed(state string) bool {
	k := w.Kind(state)
	return k == KindAccepted || k == KindCancelled
}

// Underway reports whether an item in state is pulled and not yet closed:
// its kind is active or waiting.
func (w Workflow) Underway(state string) bool {
	k := w.Kind(state)
	return k == KindActive || k == KindWaiting
}

package serve

import (
	"slices"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The helpers here read what flai serve's agent, planner, and orchestrator
// passes act on from a state's kind and a transition's effect, never from a
// state's name (ADR-0301 §2, §3; S-0742). Each takes the state as an item's
// status gives it, a state's name or one of its aliases.

// open opens the project at root as workitem.Open does, running the
// launcher's workflow, so that the launcher's reads and moves follow the
// board its options give.
func (l *launcher) open(root string) (*workitem.Repo, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	repo.Workflow = l.flow
	return repo, nil
}

// pullTransition is the transition a story makes when its agent pulls it:
// its From are the states the pull takes from, its To the state the story
// then works in. ok is false for a workflow with none, which no valid
// workflow is.
func pullTransition(wf workflow.Workflow) (t workflow.Transition, ok bool) {
	return wf.WithEffect(workflow.TypeStory, workflow.EffectPull)
}

// isPullable reports whether a story in state waits to be pulled: whether the
// pull takes from state.
func isPullable(wf workflow.Workflow, state string) bool {
	t, ok := pullTransition(wf)
	return ok && slices.Contains(t.From, stateName(wf, state))
}

// isPulledInto reports whether state is the one the pull moves a story into.
func isPulledInto(wf workflow.Workflow, state string) bool {
	t, ok := pullTransition(wf)
	return ok && t.To == stateName(wf, state)
}

// pullQueue is the first state the pull takes from, in board order: where
// a story waits to be pulled, ready in the baseline.
func pullQueue(wf workflow.Workflow) string {
	t, ok := pullTransition(wf)
	if !ok {
		return ""
	}
	for _, s := range wf.States {
		if slices.Contains(t.From, s.Name) {
			return s.Name
		}
	}
	return ""
}

// stateName is the current name of the state named state or listing it
// among its aliases, else state as it is.
func stateName(wf workflow.Workflow, state string) string {
	if s, ok := wf.State(state); ok {
		return s.Name
	}
	return state
}

// isActive reports whether state is one an agent works a story in: its kind
// is active.
func isActive(wf workflow.Workflow, state string) bool {
	return wf.Kind(state) == workflow.KindActive
}

// isWaiting reports whether state is one a story waits in for its acceptance
// or a send back: its kind is waiting.
func isWaiting(wf workflow.Workflow, state string) bool {
	return wf.Kind(state) == workflow.KindWaiting
}

// isUnderway reports whether a story in state is pulled and not yet closed:
// its kind is active or waiting.
func isUnderway(wf workflow.Workflow, state string) bool {
	return isActive(wf, state) || isWaiting(wf, state)
}

// isQueued reports whether state is one an item waits in before it is pulled:
// its kind is queue.
func isQueued(wf workflow.Workflow, state string) bool {
	return wf.Kind(state) == workflow.KindQueue
}

// isAccepted reports whether state is the one an acceptance closes an item in:
// its kind is accepted.
func isAccepted(wf workflow.Workflow, state string) bool {
	return wf.Kind(state) == workflow.KindAccepted
}

// isClosed reports whether state closes an item: its kind is accepted or
// cancelled.
func isClosed(wf workflow.Workflow, state string) bool {
	k := wf.Kind(state)
	return k == workflow.KindAccepted || k == workflow.KindCancelled
}

// activeRoom is the room the WIP limits of the active-work states leave:
// limit is the sum of the limits above zero that limits sets on the states
// of the active kind, and room is limit less the stories counts holds in
// every state of that kind, never below zero. With no such limit, room is -1
// and limit 0: nothing limits the pull. limits and counts are by state
// name, as the board and its view keep them.
func activeRoom(wf workflow.Workflow, limits, counts map[string]int) (room, limit int) {
	in := 0
	for _, s := range wf.OfKind(workflow.KindActive) {
		if n := limits[s.Name]; n > 0 {
			limit += n
		}
		in += counts[s.Name]
	}
	if limit == 0 {
		return -1, 0
	}
	return max(0, limit-in), limit
}

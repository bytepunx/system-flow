package workitem

import (
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// The fixture's states that the baseline lacks: in-progress renamed doing,
// with no alias, and testing, a second active state with a limit, between
// it and review (ADR-0301 §2, §6).
const (
	fixtureDoing   = "doing"
	fixtureTesting = "testing"
)

// testFlowFile is the workflow the package's tests run beside the baseline:
// the baseline with in-progress renamed doing and testing added between
// doing and review, every transition the baseline's under the new names,
// with doing to testing, testing to review, and testing back to doing, so
// that an epic can walk it in board order and back (S-0739).
const testFlowFile = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: doing, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: doing, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - {from: doing, to: testing}
  - from: [doing, testing]
    to: review
    rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: [doing, testing], to: done, types: [task], rules: [children-closed]}
  - {from: review, to: doing, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: testing, to: doing}
  - {from: doing, to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, doing, testing], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: doing}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: doing}}
  - {name: review, title: Review time, sum: [review]}
---
`

// testFlow is the fixture workflow, parsed afresh, so a test may change it
// without changing another's.
func testFlow(t *testing.T) *workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", testFlowFile)
	if p != nil {
		t.Fatalf("the fixture workflow does not parse: %s", p)
	}
	return &w
}

// withTestFlow sets r to run the fixture workflow and returns it, for a test
// that opens its repository with newProject or another helper.
func withTestFlow(t *testing.T, r *Repo) *Repo {
	t.Helper()
	r.Workflow = testFlow(t)
	return r
}

// inTestFlow is the fixture's name for the baseline state named state: doing
// for in-progress, the same name for every other. A case written for the
// baseline runs under the fixture through it.
func inTestFlow(state string) string {
	if state == InProgress {
		return fixtureDoing
	}
	return state
}

// S-0739: the board's roles are read from the workflow by kind and effect,
// so the baseline and the fixture agree on each where their kinds agree and
// the fixture's renamed and added states take the roles they have.
func TestFlowRolesFollowTheWorkflow(t *testing.T) {
	same := func(s string) string { return s }
	for _, tc := range []struct {
		name  string
		w     workflow.Workflow
		state func(string) string
	}{
		{"baseline", workflow.Baseline(), same},
		{"fixture", *testFlow(t), inTestFlow},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, st := tc.w, tc.state
			for role, got := range map[string][2]string{
				"pull source":   {pullSource(w), st(Ready)},
				"pull target":   {pullTarget(w), st(InProgress)},
				"accept source": {acceptSource(w), st(Review)},
				"accepted":      {acceptedState(w), st(Done)},
				"cancelled":     {cancelledState(w), st(Cancelled)},
				"initial":       {initialState(w), st(Backlog)},
			} {
				if got[0] != got[1] {
					t.Errorf("%s: %q, want %q", role, got[0], got[1])
				}
			}
			if got, want := orderedStates(w), []string{st(Backlog), st(Ready)}; !slices.Equal(got, want) {
				t.Errorf("ordered: %v, want %v", got, want)
			}
			want := []string{st(Ready), st(InProgress), st(Review)}
			if tc.name == "fixture" {
				want = slices.Insert(want, 2, fixtureTesting)
			}
			if got := stateNames(limitedStates(w)); !slices.Equal(got, want) {
				t.Errorf("limited: %v, want %v", got, want)
			}
			for _, name := range w.Names() {
				kind := w.Kind(name)
				closed := name == Done || name == Cancelled
				started := !closed && name != Backlog && name != Ready || name == Done
				if closedKind(kind) != closed || startedKind(kind) != started || underwayKind(kind) != (started && !closed) {
					t.Errorf("%s (%s): closed %v, started %v, underway %v", name, kind, closedKind(kind), startedKind(kind), underwayKind(kind))
				}
			}
		})
	}
}

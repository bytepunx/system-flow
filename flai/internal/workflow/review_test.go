package workflow

import (
	"slices"
	"testing"
)

// renamedFlow is a board whose in-progress is doing and whose review is
// awaiting, neither with an alias, with testing, of the active kind, added
// between them: a name read anywhere in place of a kind or an effect gives
// the wrong answer under it (S-0743).
const renamedFlow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: doing, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: awaiting, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: doing, effects: [pull]}
  - {from: doing, to: testing}
  - {from: testing, to: awaiting}
  - {from: awaiting, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: awaiting, to: done, types: [task]}
  - {from: doing, to: done, types: [task]}
  - {from: awaiting, to: doing, effects: [send-back]}
  - {from: testing, to: doing}
  - {from: doing, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, doing, testing], to: cancelled, effects: [cancel]}
  - {from: awaiting, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// S-0743: the roles the review passes read come from kinds and effects: the
// baseline's review, in-progress, done, and cancelled, and the renamed
// board's awaiting, doing, and testing in their places.
func TestReviewRoles(t *testing.T) {
	for _, c := range []struct {
		name                 string
		w                    Workflow
		waiting, back, added string
	}{
		{"baseline", Baseline(), "review", "in-progress", ""},
		{"renamed", mustParse(t, renamedFlow), "awaiting", "doing", "testing"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := c.w
			if got := w.AcceptFrom(TypeStory); !slices.Equal(got, []string{c.waiting}) {
				t.Errorf("AcceptFrom(story) = %v, want [%s]", got, c.waiting)
			}
			if got := w.AcceptFrom(TypeTask); got != nil {
				t.Errorf("AcceptFrom(task) = %v: no task transition gives accept", got)
			}
			if !w.AwaitsAcceptance(TypeStory, c.waiting) {
				t.Errorf("%s does not await the acceptance", c.waiting)
			}
			for _, s := range []string{c.back, c.added, "ready", "done", "cancelled", ""} {
				if s != "" && w.AwaitsAcceptance(TypeStory, s) {
					t.Errorf("%s awaits the acceptance", s)
				}
			}
			tr, ok := w.SendBack(TypeStory, c.waiting)
			if !ok || tr.To != c.back {
				t.Errorf("SendBack(story, %s) = %v %v, want a move to %s", c.waiting, tr, ok, c.back)
			}
			if _, ok := w.SendBack(TypeStory, c.back); ok {
				t.Errorf("SendBack(story, %s) found a move: no send-back leaves from it", c.back)
			}
			if !w.Gives(TypeStory, c.waiting, c.back, EffectSendBack) {
				t.Errorf("%s to %s gives no send-back", c.waiting, c.back)
			}
			if w.Gives(TypeStory, c.back, "ready", EffectSendBack) {
				t.Errorf("%s to ready gives a send-back", c.back)
			}
			if !w.Gives(TypeStory, c.waiting, "done", EffectAccept) || !w.Gives(TypeStory, c.back, "cancelled", EffectCancel) {
				t.Error("the acceptance or the cancellation gives no effect of its own")
			}
			if w.Gives(TypeStory, "done", c.back, EffectSendBack) {
				t.Error("a move the workflow lacks gives an effect")
			}
			for _, s := range []string{"done", "cancelled"} {
				if !w.Closed(s) {
					t.Errorf("%s is not closed", s)
				}
			}
			for _, s := range []string{"backlog", "ready", c.back, c.waiting, c.added, "nowhere"} {
				if s != "" && w.Closed(s) {
					t.Errorf("%s is closed", s)
				}
			}
			for _, s := range []string{c.back, c.added, c.waiting} {
				if s != "" && !w.Underway(s) {
					t.Errorf("%s is not underway", s)
				}
			}
			for _, s := range []string{"backlog", "ready", "done", "cancelled", "nowhere"} {
				if w.Underway(s) {
					t.Errorf("%s is underway", s)
				}
			}
		})
	}
}

// S-0743: a state named by an alias awaits the acceptance as its state does.
func TestAwaitsAcceptanceByAlias(t *testing.T) {
	w := mustParse(t, renamedFlow)
	for i := range w.States {
		if w.States[i].Name == "awaiting" {
			w.States[i].Aliases = []string{"review"}
		}
	}
	if !w.AwaitsAcceptance(TypeStory, "review") {
		t.Error("review, awaiting's alias, does not await the acceptance")
	}
	if tr, ok := w.SendBack(TypeStory, "review"); !ok || tr.To != "doing" {
		t.Errorf("SendBack(story, review) = %v %v, want a move to doing", tr, ok)
	}
}

package release

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// renamedWorkflow is a board with done renamed shipped, cancelled renamed
// dropped, and a testing state added between in-progress and review.
const renamedWorkflow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: shipped, kind: accepted}
  - {name: dropped, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: testing}
  - {from: testing, to: review}
  - {from: review, to: shipped, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: review, to: shipped, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: [backlog, ready, in-progress, testing], to: dropped, effects: [cancel]}
  - {from: review, to: dropped, by: [cascade], effects: [cancel]}
  - {from: dropped, to: backlog, effects: [reopen]}
---
`

func renamedFlow(t *testing.T) workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", renamedWorkflow)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %v", p)
	}
	return w
}

// S-0749: under a workflow that renames done and cancelled, a story is
// accepted by its state's kind, the theme leaves out the stories in a
// cancelled state, and whole_epics releases an epic's stories once it is in
// a waiting or accepted state, naming those states.
func TestEvaluateReadsTheWorkflowsKinds(t *testing.T) {
	flow := renamedFlow(t)
	items := []*workitem.Item{
		{ID: "E-0001", Type: workitem.Epic, Status: "shipped"},
		{ID: "E-0002", Type: workitem.Epic, Status: "testing"},
		evalStory("S-0001", "shipped", "E-0001", ptr(10.0)),
		evalStory("S-0002", "dropped", "E-0001", nil),
		evalStory("S-0003", "testing", "E-0001", nil),
		evalStory("S-0004", "shipped", "E-0002", ptr(5.0)),
	}
	pending := map[string]bool{"S-0001": true, "S-0004": true}

	ev := EvaluateJudgedIn(flow, withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0001"}), items, pending, "")
	var theme []string
	for _, s := range ev.Theme {
		theme = append(theme, s.ID)
	}
	if !slices.Equal(theme, []string{"S-0001", "S-0003"}) || !slices.Equal(ev.NotAccepted, []string{"S-0003"}) {
		t.Errorf("the theme leaves out the dropped story and counts the shipped one accepted: theme %v, not accepted %v", theme, ev.NotAccepted)
	}
	if !ev.Theme[0].Accepted || ev.Theme[1].Accepted {
		t.Errorf("a story is accepted in a state of the accepted kind only: %+v", ev.Theme)
	}

	ev = EvaluateJudgedIn(flow, withRelease(manifest.Release{Policy: manifest.ReleaseThreshold, Count: ptr(1), WholeEpics: true}), items, pending, "")
	var held []string
	for _, h := range ev.HeldByEpic {
		held = append(held, h.ID)
	}
	if ev.Met || !slices.Equal(held, []string{"S-0004"}) {
		t.Fatalf("only the story whose epic is testing is held: held %v, %+v", held, ev)
	}
	if !strings.Contains(ev.Reason, "until its epic is in review or shipped: E-0002 is testing, with S-0004") {
		t.Errorf("the reason names the workflow's waiting and accepted states: %s", ev.Reason)
	}

	if ev := EvaluateJudgedIn(flow, withRelease(manifest.Release{Policy: manifest.ReleaseTheme, Epic: "E-0001"}), []*workitem.Item{evalStory("S-0005", workitem.Done, "E-0001", nil)}, nil, ""); ev.Theme[0].Accepted {
		t.Errorf("done, which the workflow does not have, is not accepted: %+v", ev.Theme)
	}
}

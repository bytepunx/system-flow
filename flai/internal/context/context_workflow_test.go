package context

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// renamedWorkflow is a board with backlog renamed icebox, cancelled renamed
// dropped, and a testing state added between in-progress and review.
const renamedWorkflow = `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: icebox, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: dropped, kind: cancelled}
transitions:
  - {from: icebox, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: testing}
  - {from: testing, to: review}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: review, to: done, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: [icebox, ready, in-progress, testing], to: dropped, effects: [cancel]}
  - {from: review, to: dropped, by: [cascade], effects: [cancel]}
  - {from: dropped, to: icebox, effects: [reopen]}
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

// S-0749: the ready reason's board lays out the workflow's columns in its
// order, the initial state's as a count, and a story's claim leaves out its
// tasks in a state of the cancelled kind.
func TestTheReasonsFiguresAndAStorysClaimReadTheWorkflowsKinds(t *testing.T) {
	flow := renamedFlow(t)
	card := func(id string) workitem.BoardCard {
		return workitem.BoardCard{ID: id, Title: "Story " + id, Nature: "feature", Age: "1h"}
	}
	view := workitem.BoardView{
		Columns:   map[string][]workitem.BoardCard{"icebox": {card("S-0001")}, "testing": {card("S-0002")}},
		WIPLimits: map[string]int{"testing": 2},
		Counts:    map[string]int{"icebox": 1, "testing": 1},
	}
	text := boardText(flow, view)
	want := "icebox: stories 1, each under promote_candidates\nready\n  -\nin-progress\n  -\ntesting (1/2)\n  S-0002  Story S-0002  feature  1h\nreview\n  -\ndone\n  -\ndropped\n  -\n"
	if text != want {
		t.Errorf("board under the workflow:\n%s\nwant:\n%s", text, want)
	}

	repo := &workitem.Repo{Workflow: &flow}
	story := &workitem.Item{ID: "S-0003", Type: workitem.Story, Touches: []string{"a.go"}}
	items := []*workitem.Item{
		{ID: "T-0001", Type: workitem.Task, Parent: "S-0003", Status: "testing", Touches: []string{"b.go"}},
		{ID: "T-0002", Type: workitem.Task, Parent: "S-0003", Status: "dropped", Touches: []string{"c.go"}},
	}
	claim := storyClaim(repo, story, items)
	if !slices.Contains(claim, "b.go") || slices.ContainsFunc(claim, func(p string) bool { return strings.Contains(p, "c.go") }) {
		t.Errorf("the claim takes the testing task's touches and leaves out the dropped one's: %v", claim)
	}
}

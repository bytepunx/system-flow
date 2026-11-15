package statsread

import (
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
)

// S-0749: the metrics' WIP limit is that of the state a pull moves a story
// into, by the pull effect, whatever the workflow names it; in-progress on
// the baseline.
func TestTheWIPLimitIsThatOfTheStateAPullMovesInto(t *testing.T) {
	if got := pulledInto(workflow.Baseline()); got != "in-progress" {
		t.Errorf("baseline: %q, want in-progress", got)
	}
	w, p := workflow.Parse("design/workflows/board.md", `---
title: Board
updated: 2026-11-15
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: building, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: shipped, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: building, effects: [pull]}
  - {from: building, to: testing}
  - {from: testing, to: review}
  - {from: review, to: shipped, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: review, to: shipped, types: [task]}
  - {from: review, to: building, effects: [send-back]}
  - {from: [backlog, ready, building, testing], to: cancelled, effects: [cancel]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %v", p)
	}
	if got := pulledInto(w); got != "building" {
		t.Errorf("renamed: %q, want building", got)
	}
}

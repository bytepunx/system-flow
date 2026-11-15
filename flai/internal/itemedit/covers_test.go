package itemedit

import (
	"reflect"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Who a story's change reaches (S-0333): each other open story whose claim
// covers a changed path, by the prefix rule, a component by its path, a
// task's touches as its story's, a shared path like any other, and an empty
// claim as everything.
func TestCovering(t *testing.T) {
	story := func(id, status string, touches ...string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status, Title: "title " + id, Touches: touches}
	}
	archived := story("S-0009", workitem.InProgress, "flai/cmd")
	archived.Archived = true
	items := []*workitem.Item{
		story("S-0001", workitem.InProgress, "flai/cmd"),           // the changer
		story("S-0002", workitem.InProgress, "flai/cmd/accept.go"), // a file it changed
		story("S-0003", workitem.Review, "cli"),                    // the component, by tag
		story("S-0004", workitem.InProgress, "flaiover"),           // elsewhere
		story("S-0005", workitem.InProgress, "flaiover"),           // its task claims design/system
		story("S-0006", workitem.Review),                           // no touches: may change anything
		story("S-0007", workitem.Ready, "flai/cmd"),                // not open
		story("S-0008", workitem.InProgress, "CHANGELOG.md"),       // told whatever the shared paths hold
		archived,
		{ID: "T-0001", Type: workitem.Task, Status: workitem.InProgress, Parent: "S-0005", Touches: []string{"design/system"}},
	}
	projects := []manifest.Project{{Name: "flai", Path: "flai", Tags: []string{"cli"}}, {Name: "flaiover", Path: "flaiover"}}
	changed := []string{"flai/cmd/accept.go", "flai/cmd/move.go", "design/system/workflow.md", "CHANGELOG.md"}
	got := Covering(items, projects, "S-1", changed)
	want := []Covered{
		{ID: "S-0002", Title: "title S-0002", Paths: []string{"flai/cmd/accept.go"}},
		{ID: "S-0003", Title: "title S-0003", Paths: []string{"flai/cmd/accept.go", "flai/cmd/move.go"}},
		{ID: "S-0005", Title: "title S-0005", Paths: []string{"design/system/workflow.md"}},
		{ID: "S-0006", Title: "title S-0006", Paths: changed},
		{ID: "S-0008", Title: "title S-0008", Paths: []string{"CHANGELOG.md"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("covering:\n got %+v\nwant %+v", got, want)
	}
	if got := Covering(items, projects, "S-0001", nil); got != nil {
		t.Errorf("a change of nothing reaches nobody, got %+v", got)
	}
	if got := Covering(items, projects, "S-0001", []string{"README.md"}); got[0].ID != "S-0006" || len(got) != 1 {
		t.Errorf("a path no claim covers reaches only the empty claim, got %+v", got)
	}
}

// editOtherWorkflow is a board unlike the baseline: review is renamed
// acceptance-review, keeping review as its alias, and testing, of kind
// active, sits between in-progress and it (ADR-0301 §6).
const editOtherWorkflow = `---
title: Board
updated: 2026-11-14
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, kind: active, limit: true}
  - {name: acceptance-review, kind: waiting, limit: true, aliases: [review]}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit]}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: testing}
  - {from: testing, to: acceptance-review}
  - {from: acceptance-review, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: acceptance-review, to: done, types: [task]}
  - {from: in-progress, to: done, types: [task]}
  - {from: acceptance-review, to: in-progress, effects: [send-back]}
  - {from: in-progress, to: ready}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress, testing], to: cancelled, effects: [cancel]}
  - {from: acceptance-review, to: cancelled, by: [cascade], effects: [cancel]}
  - {from: cancelled, to: backlog, effects: [reopen]}
---
`

// editWorkflow is editOtherWorkflow parsed.
func editWorkflow(t *testing.T) workflow.Workflow {
	t.Helper()
	w, p := workflow.Parse("design/workflows/board.md", editOtherWorkflow)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	return w
}

// S-0746: under a board with a state renamed and one added, a change reaches
// the stories open by the workflow's kinds, in testing and in
// acceptance-review, and not one in ready or done.
func TestCoveringInReachesTheStoriesOpenByTheWorkflowsKinds(t *testing.T) {
	story := func(id, status string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status, Title: "title " + id, Touches: []string{"flai/cmd"}}
	}
	items := []*workitem.Item{
		story("S-0001", "in-progress"),
		story("S-0002", "testing"),
		story("S-0003", "acceptance-review"),
		story("S-0004", "ready"),
		story("S-0005", "done"),
	}
	changed := []string{"flai/cmd/accept.go"}
	got := CoveringIn(editWorkflow(t), items, nil, "S-0001", changed)
	want := []Covered{
		{ID: "S-0002", Title: "title S-0002", Paths: changed},
		{ID: "S-0003", Title: "title S-0003", Paths: changed},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("covering:\n got %+v\nwant %+v", got, want)
	}
	// the baseline has no testing or acceptance-review: neither is open there
	if got := Covering(items, nil, "S-0001", changed); got != nil {
		t.Errorf("the baseline reached %+v", got)
	}
}

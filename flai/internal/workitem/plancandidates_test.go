package workitem

import (
	"reflect"
	"strings"
	"testing"
)

// S-0219: an epic is a candidate for the planner when it is in the backlog
// with no story, or open with every story done or cancelled and one done;
// the caller's epics are left out with its reason.
func TestPlanCandidates(t *testing.T) {
	epic := func(id, status string) *Item { return &Item{ID: id, Type: Epic, Status: status, Title: "Epic " + id} }
	story := func(id, parent, status string) *Item {
		return &Item{ID: id, Type: Story, Status: status, Parent: parent, Title: "Story " + id}
	}
	archivedDone := story("S-0002", "E-0002", Done)
	archivedDone.Archived = true
	archivedEpic := epic("E-0009", Backlog)
	archivedEpic.Archived = true
	draft := story("S-0009", "E-0010", Backlog)
	draft.Draft = true
	items := []*Item{
		epic("E-0001", Backlog),    // no stories: a candidate
		epic("E-0002", InProgress), // every story done or cancelled: a candidate
		story("S-0001", "E-0002", Done), archivedDone, story("S-0003", "E-0002", Cancelled),
		epic("E-0003", InProgress), // an open story: not one
		story("S-0004", "E-0003", Done), story("S-0005", "E-0003", Review),
		epic("E-0004", Done), // done: not one
		story("S-0006", "E-0004", Done),
		epic("E-0005", Cancelled), // cancelled, with no stories: not one
		epic("E-0006", Backlog),   // every story cancelled, none done: not one
		story("S-0007", "E-0006", Cancelled),
		epic("E-0007", InProgress), // no stories, but not in the backlog: not one
		epic("E-0008", Backlog),    // a candidate the caller leaves out
		archivedEpic,               // archived: never judged
		epic("E-0010", Backlog),    // a draft story is a story: not one
		draft,
		story("S-0008", "E-0011", Done), // a story whose epic is not among the items
	}
	got := PlanCandidatesOf(items, map[string]string{"E-0008": "a planner runs for it now", "E-0003": "ignored: not a candidate"})
	var ids []string
	for _, c := range got.Candidates {
		ids = append(ids, c.ID)
	}
	if !reflect.DeepEqual(ids, []string{"E-0001", "E-0002", "S-0009"}) {
		t.Fatalf("candidates %v, want E-0001, E-0002, and the draft S-0009: %+v", ids, got)
	}
	if c := got.Candidates[0]; c.Type != Epic || c.Status != Backlog || c.Title != "Epic E-0001" || !strings.HasPrefix(c.Reason, "in the backlog with no stories") {
		t.Errorf("E-0001: %+v", c)
	}
	if c := got.Candidates[1]; c.Status != InProgress || !strings.HasPrefix(c.Reason, "every story done or cancelled (2 done, 1 cancelled) while the epic is in-progress") {
		t.Errorf("E-0002, its archived story counted: %+v", c)
	}
	// a draft story without a plan is a candidate itself (S-0328)
	if c := got.Candidates[2]; c.Type != Story || c.Reason != "in the backlog without a plan: no touches; no forecast duration; no cost of delay value; no tasks" {
		t.Errorf("S-0009, a draft: %+v", c)
	}
	if want := []PlanCandidate{{ID: "E-0008", Type: Epic, Title: "Epic E-0008", Status: Backlog, Reason: "a planner runs for it now"}}; !reflect.DeepEqual(got.LeftOut, want) {
		t.Errorf("left out %+v, want %+v", got.LeftOut, want)
	}

	none := PlanCandidatesOf([]*Item{epic("E-0001", Done)}, nil)
	if none.Candidates == nil || none.LeftOut == nil || len(none.Candidates)+len(none.LeftOut) != 0 {
		t.Errorf("no candidates are empty lists, not null: %+v", none)
	}
}

// S-0328: a story is a candidate when it is in the backlog, not archived,
// and lacks touches, a forecast duration, a cost of delay value, or a task
// that is not cancelled, archived tasks counting; the reason names only what
// it lacks, and the caller's stories are left out with its reason.
func TestPlanCandidatesStories(t *testing.T) {
	value := 900.0
	planned := func(id, status string) *Item {
		return &Item{ID: id, Type: Story, Status: status, Title: "Story " + id, Touches: []string{"flai/cmd"},
			Forecast: &Forecast{Duration: "6h"}, CostOfDelay: &CostOfDelay{Value: &value}}
	}
	task := func(id, story, status string) *Item { return &Item{ID: id, Type: Task, Status: status, Parent: story} }
	bare := func(id, status string) *Item { return &Item{ID: id, Type: Story, Status: status, Title: "Story " + id} }

	inputsOnly := planned("S-0003", Backlog) // inputs but no value
	inputsOnly.CostOfDelay = &CostOfDelay{Inputs: &CostInputs{RevenuePerWeek: &value}}
	zeroForecast := planned("S-0004", Backlog) // a forecast of no duration
	zeroForecast.Forecast = &Forecast{Duration: "0s"}
	archived := bare("S-0008", Backlog)
	archived.Archived = true
	doneTask := task("T-0005", "S-0005", Done)
	doneTask.Archived = true
	items := []*Item{
		bare("S-0010", Backlog),    // nothing: a candidate, after S-0002 in ID order
		planned("S-0001", Backlog), // a plan, with a task: not one
		task("T-0001", "S-0001", InProgress),
		planned("S-0002", Backlog), // its one task cancelled: a candidate
		task("T-0002", "S-0002", Cancelled),
		inputsOnly, task("T-0003", "S-0003", Backlog),
		zeroForecast, task("T-0004", "S-0004", Backlog),
		planned("S-0005", Backlog), // an archived task counts: not one
		doneTask,
		bare("S-0006", Ready),   // not in the backlog: not one
		bare("S-0007", Backlog), // left out by the caller
		archived,                // archived: never judged
	}
	got := PlanCandidatesOf(items, map[string]string{"S-0007": "a planner runs for it now", "S-0001": "ignored: not a candidate"})
	want := []PlanCandidate{
		{ID: "S-0002", Type: Story, Title: "Story S-0002", Status: Backlog, Reason: "in the backlog without a plan: no tasks"},
		{ID: "S-0003", Type: Story, Title: "Story S-0003", Status: Backlog, Reason: "in the backlog without a plan: no cost of delay value"},
		{ID: "S-0004", Type: Story, Title: "Story S-0004", Status: Backlog, Reason: "in the backlog without a plan: no forecast duration"},
		{ID: "S-0010", Type: Story, Title: "Story S-0010", Status: Backlog, Reason: "in the backlog without a plan: no touches; no forecast duration; no cost of delay value; no tasks"},
	}
	if !reflect.DeepEqual(got.Candidates, want) {
		t.Errorf("candidates\n%+v\nwant\n%+v", got.Candidates, want)
	}
	if want := []PlanCandidate{{ID: "S-0007", Type: Story, Title: "Story S-0007", Status: Backlog, Reason: "a planner runs for it now"}}; !reflect.DeepEqual(got.LeftOut, want) {
		t.Errorf("left out %+v, want %+v", got.LeftOut, want)
	}
}

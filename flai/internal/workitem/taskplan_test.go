package workitem

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func planTask(id, status string, after ...string) *Item {
	return &Item{ID: id, Type: Task, Parent: "S-0001", Status: status, After: after}
}

// S-0176: a task is ready to start when every task of its after is done or
// cancelled, waiting while one is open, and in progress in review too; an
// entry that names no task of the story holds nothing.
func TestPlanStates(t *testing.T) {
	items := []*Item{
		{ID: "S-0001", Type: Story, Status: InProgress},
		planTask("T-0005", Ready, "T-0001", "T-0002", "T-0004"),
		planTask("T-0001", Done),
		planTask("T-0002", InProgress),
		planTask("T-0003", Cancelled),
		planTask("T-0004", Review, "t-1"),
		planTask("T-0006", Backlog, "T-0001", "T-0003"),
		planTask("T-0007", Ready, "T-0009"),
		{ID: "T-0008", Type: Task, Parent: "S-0002", Status: Ready},
	}
	p := PlanOf(items, "S-0001")
	var got []string
	for _, pt := range p.Tasks {
		got = append(got, fmt.Sprintf("%s %s %v", pt.ID, pt.State, pt.WaitingFor))
	}
	want := []string{
		"T-0001 done []",
		"T-0002 in-progress []",
		"T-0003 cancelled []",
		"T-0004 in-progress []",
		"T-0005 waiting [T-0002 T-0004]",
		"T-0006 ready []",
		"T-0007 ready []",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("states:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if s := *p.Summary(); s != (TaskSummary{Ready: 2, Waiting: 1, InProgress: 2, Done: 1, Layers: 3}) {
		t.Errorf("summary: %+v", s)
	}
	if PlanOf(items, "S-0003") != nil {
		t.Error("a story without tasks has no plan")
	}
}

// S-0176: layers group the open and done tasks by the longest chain of after
// steps before each, in ID order; a cancelled task is in none and holds no
// one up; a task on a cycle, or after one, is in none.
func TestPlanLayers(t *testing.T) {
	items := []*Item{
		planTask("T-0010", Ready, "T-0002"),
		planTask("T-0002", Done),
		planTask("T-0003", Ready),
		planTask("T-0004", Ready, "T-0002", "T-0010"),
		planTask("T-0005", Cancelled),
		planTask("T-0006", Ready, "T-0005"),
		planTask("T-0007", Ready, "T-0008"),
		planTask("T-0008", Ready, "T-0007"),
		planTask("T-0009", Ready, "T-0008", "T-0003"),
	}
	p := PlanOf(items, "S-0001")
	data, _ := json.Marshal(p.Layers)
	if string(data) != `[["T-0002","T-0003","T-0006"],["T-0010"],["T-0004"]]` {
		t.Errorf("layers: %s", data)
	}
	if got := PlanOf([]*Item{planTask("T-0001", Ready, "T-0001")}, "S-0001").Layers; len(got) != 0 {
		t.Errorf("a task after itself: %v", got)
	}
}

// S-0176: the plan's JSON is the contract flaiover reads: after and
// waiting_for left out when empty, layers always present.
func TestPlanJSON(t *testing.T) {
	p := PlanOf([]*Item{planTask("T-0001", Ready, "T-0002"), planTask("T-0002", Ready, "T-0001")}, "S-0001")
	data, _ := json.Marshal(p)
	want := `{"tasks":[{"id":"T-0001","state":"waiting","after":["T-0002"],"waiting_for":["T-0002"]},{"id":"T-0002","state":"waiting","after":["T-0001"],"waiting_for":["T-0001"]}],"layers":[]}`
	if string(data) != want {
		t.Errorf("json:\n%s\nwant:\n%s", data, want)
	}
	p = PlanOf([]*Item{planTask("T-0001", Done)}, "S-0001")
	data, _ = json.Marshal(p)
	if string(data) != `{"tasks":[{"id":"T-0001","state":"done"}],"layers":[["T-0001"]]}` {
		t.Errorf("json: %s", data)
	}
}

// S-0176: moving a task to in progress while a task of its after is open
// warns in the voice of a story's hold, and moves it.
func TestMovingAWaitingTaskWarns(t *testing.T) {
	r := &Repo{}
	first, second := planTask("T-0001", InProgress), planTask("T-0002", Ready)
	third := planTask("T-0003", Ready, "T-0001", "T-0002")
	items := []*Item{{ID: "S-0001", Type: Story, Status: InProgress}, first, second, third}
	warnings, err := r.Move(third, InProgress, MoveOptions{By: "agent", Now: time.Now(), Items: items})
	if err != nil {
		t.Fatal(err)
	}
	want := "T-0003 is waiting (after): waits for T-0001 (in progress) and T-0002 (ready); ready to start when T-0001 and T-0002 are done or cancelled"
	if strings.Join(warnings, "\n") != want || third.Status != InProgress {
		t.Errorf("warnings: %q, status %s", warnings, third.Status)
	}

	first.Status, second.Status = Done, Cancelled
	fourth := planTask("T-0004", Ready, "T-0001", "T-0002")
	items = append(items, fourth)
	if warnings, err := r.Move(fourth, InProgress, MoveOptions{By: "agent", Now: time.Now(), Items: items}); err != nil || len(warnings) != 0 {
		t.Errorf("waits for nothing open: %q %v", warnings, err)
	}
}

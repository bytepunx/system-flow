package workitem

import (
	"os"
	"strings"
	"testing"
)

// S-0199: a story is made a draft when it is created; nothing else is.
func TestCreateDraft(t *testing.T) {
	r := newProject(t)
	s, err := r.Create(NewOptions{Type: Story, Title: "Drafted", Draft: true, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(s.Path)
	if !s.Draft || !strings.Contains(string(data), "\ndraft: true\n") {
		t.Errorf("the story is a draft:\n%s", data)
	}
	if plain, err := r.Create(NewOptions{Type: Story, Title: "Plain", Now: t0}); err != nil || plain.Draft {
		t.Errorf("a story is not a draft unless asked: %v", err)
	}
	for _, typ := range []string{Epic, Task} {
		if _, err := r.Create(NewOptions{Type: typ, Title: "No", Parent: map[string]string{Task: s.ID}[typ], Draft: true, Now: t0}); err == nil || !strings.Contains(err.Error(), "only a story is a draft") {
			t.Errorf("a draft %s must be refused: %v", typ, err)
		}
	}
}

// S-0203: a story is created with a cost of delay, as one made from an issue
// is; a task is not, and a block without who set it is refused.
func TestCreateCostOfDelay(t *testing.T) {
	r := newProject(t)
	cod := &CostOfDelay{Inputs: &CostInputs{TimeLostPerCycle: "2h", By: "flai", At: t0.Format(TimeFormat)}}
	s, err := r.Create(NewOptions{Type: Story, Title: "Costed", CostOfDelay: cod, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	cod.Inputs.TimeLostPerCycle = "9h"
	got, err := r.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if c := got.CostOfDelay; c.IsZero() || c.Inputs.TimeLostPerCycle != "2h" || c.Inputs.By != "flai" || c.Inputs.At != t0.Format(TimeFormat) || c.By != "" {
		t.Errorf("the story carries the cost of delay it was given: %+v", c)
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "No", Parent: s.ID, CostOfDelay: cod, Now: t0}); err == nil || !strings.Contains(err.Error(), "not a task") {
		t.Errorf("a task with a cost of delay must be refused: %v", err)
	}
	unsigned := &CostOfDelay{Inputs: &CostInputs{TimeLostPerCycle: "2h", At: t0.Format(TimeFormat)}}
	if _, err := r.Create(NewOptions{Type: Story, Title: "Unsigned", CostOfDelay: unsigned, Now: t0}); err == nil || !strings.Contains(err.Error(), "cost_of_delay.inputs.by is required") {
		t.Errorf("a cost of delay without who set it must be refused: %v", err)
	}
}

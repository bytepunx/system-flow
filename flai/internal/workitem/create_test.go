package workitem

import (
	"os"
	"slices"
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

// I-0071: a touch is a path in the repository or a component name; one that
// starts with a dot is a path like any other, and one that leaves the
// repository, or would split or read as a flag, is refused.
func TestCleanTouches(t *testing.T) {
	got, err := CleanTouches([]string{" .claude/agents/planner.md ", ".github/workflows/", ".dockerignore", "flai/cmd", "", ".github/workflows", "a..b/c", "flai"})
	want := []string{".claude/agents/planner.md", ".github/workflows", ".dockerignore", "flai/cmd", "a..b/c", "flai"}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("CleanTouches = %q, %v; want %q", got, err, want)
	}
	if got, err := CleanTouches([]string{" ", "/"}); err != nil || got != nil {
		t.Errorf("nothing given is nil: %q, %v", got, err)
	}
	for _, c := range []struct{ touch, says string }{
		{"../x", "climbs out of the repository"},
		{"a/../b", "climbs out of the repository"},
		{"..", "climbs out of the repository"},
		{"a/..", "climbs out of the repository"},
		{"/etc", "is absolute"},
		{"a,b", "holds a comma"},
		{"-x", "starts with a dash"},
		{"a\x07b", "holds a control character"},
		{"a\nb", "holds a control character"},
	} {
		_, err := CleanTouches([]string{"flai", c.touch})
		if err == nil || !strings.Contains(err.Error(), c.says) || !strings.Contains(err.Error(), "give a repository path or component name") {
			t.Errorf("touch %q: want an error saying %q and what a touch may be, got %v", c.touch, c.says, err)
		}
	}
}

// I-0071: a story or a task is created keeping a touch that starts with a
// dot, and refused with one that leaves the repository, writing nothing.
func TestCreateTouches(t *testing.T) {
	r := newProject(t)
	s, err := r.Create(NewOptions{Type: Story, Title: "Dotted", Touches: []string{".claude/agents/", "flai/cmd"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := r.Get(s.ID); err != nil || !slices.Equal(got.Touches, []string{".claude/agents", "flai/cmd"}) {
		t.Errorf("the story keeps its touches: %v %v", got, err)
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "Outside", Parent: s.ID, Touches: []string{"../elsewhere"}, Now: t0}); err == nil || !strings.Contains(err.Error(), `touch "../elsewhere"`) {
		t.Errorf("a touch outside the repository must be refused: %v", err)
	}
	if items, err := r.List(false); err != nil || len(items) != 1 {
		t.Errorf("nothing is written for a refused touch: %d items, %v", len(items), err)
	}
}

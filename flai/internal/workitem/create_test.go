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

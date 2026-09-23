package workitem

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0103: a story created while the project has a default agent gets a copy
// in its front matter; what is given for the story wins over it.
func TestAStoryGetsTheProjectsAgent(t *testing.T) {
	r := newProject(t)
	plain := mustCreate(t, r, Story, "Before any default", "")
	if plain.Agent != nil {
		t.Fatalf("no default, no agent: %+v", plain.Agent)
	}
	data, _ := os.ReadFile(plain.Path)
	if strings.Contains(string(data), "agent:") {
		t.Errorf("an item with no agent writes no agent key, so an older flai still reads it:\n%s", data)
	}

	r.Manifest.Agent = &manifest.Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high"}}
	s := mustCreate(t, r, Story, "With the default", "")
	if s.Agent.Harness != "claude-code" || s.Agent.Model != "claude-opus-5-5" || s.Agent.Config["effort"] != "high" {
		t.Errorf("copied: %+v", s.Agent)
	}
	over, err := r.Create(NewOptions{Type: Story, Title: "Overridden", Owner: "alex", Now: t0, Agent: &manifest.Agent{Model: "claude-sonnet-5", Config: map[string]string{"max_turns": "20"}}})
	if err != nil {
		t.Fatal(err)
	}
	if over.Agent.Harness != "claude-code" || over.Agent.Model != "claude-sonnet-5" || over.Agent.Config["effort"] != "high" || over.Agent.Config["max_turns"] != "20" {
		t.Errorf("overridden: %+v", over.Agent)
	}
	// read back from the file, as every other flai command sees it
	back, err := r.Get(over.ID)
	if err != nil {
		t.Fatal(err)
	}
	if back.Agent.Model != "claude-sonnet-5" || back.Agent.Config["max_turns"] != "20" {
		t.Errorf("read back: %+v", back.Agent)
	}
	// a changed default does not change stories already made
	r.Manifest.Agent = &manifest.Agent{Model: "other"}
	if again, _ := r.Get(s.ID); again.Agent.Model != "claude-opus-5-5" {
		t.Errorf("an existing story changed with the default: %+v", again.Agent)
	}

	// epics and tasks carry none
	if _, err := r.Create(NewOptions{Type: Epic, Title: "E", Owner: "alex", Now: t0, Agent: &manifest.Agent{Model: "x"}}); err == nil {
		t.Error("an epic was given an agent")
	}
	e := mustCreate(t, r, Epic, "Epic with a project default", "")
	if e.Agent != nil {
		t.Errorf("an epic got the default: %+v", e.Agent)
	}
	e.Agent = &manifest.Agent{Model: "x"}
	if err := e.Validate(); err == nil || !strings.Contains(err.Error(), "only a story carries an agent") {
		t.Errorf("validate an epic with an agent: %v", err)
	}
}

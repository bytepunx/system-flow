package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentWithMergesAndOverrides(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high", "max_turns": "50"}}
	got := def.With(&Agent{Model: "claude-sonnet-5", Config: map[string]string{"effort": "low"}})
	if got.Harness != "claude-code" || got.Model != "claude-sonnet-5" || got.Config["effort"] != "low" || got.Config["max_turns"] != "50" {
		t.Errorf("merged: %+v", got)
	}
	if def.Config["effort"] != "high" {
		t.Error("the default was changed by a merge")
	}
	var none *Agent
	if none.With(nil) != nil || none.With(&Agent{}) != nil {
		t.Error("nothing with nothing is nil")
	}
	if got := none.With(&Agent{Harness: "codex"}); got == nil || got.Harness != "codex" {
		t.Errorf("an override with no default: %+v", got)
	}
}

func TestAgentValidate(t *testing.T) {
	good := &Agent{Harness: "claude-code", Model: "anthropic/claude-opus-5-5@2026", Config: map[string]string{"max_turns": "50", "permission.mode": "plan"}}
	if err := good.Validate(); err != nil {
		t.Errorf("good: %v", err)
	}
	for _, bad := range []*Agent{
		{Harness: "Claude Code"},
		{Model: "has space"},
		{Config: map[string]string{"Bad Key": "x"}},
		{Config: map[string]string{"k": "two\nlines"}},
	} {
		if err := bad.Validate(); err == nil {
			t.Errorf("accepted %+v", bad)
		}
	}
}

// flai owns the agent block of system-flow.yaml and nothing else of it.
func TestWriteAgentKeepsEverythingElse(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	orig := "version: 1\nname: Harbour\n# a comment the operator wrote\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 4242\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	a := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high", "note": "yes"}}
	if err := WriteAgent(path, a); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Agent.Harness != "claude-code" || m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || m.Agent.Config["note"] != "yes" || m.Dashboard.Port != 4242 {
		t.Errorf("read back: %+v %+v", m.Agent, m.Dashboard)
	}
	data, _ := os.ReadFile(path)
	if !strings.HasPrefix(string(data), orig) || !strings.Contains(string(data), `note: "yes"`) {
		t.Errorf("written:\n%s", data)
	}
	// replaced, not appended twice, and removed when cleared
	if err := WriteAgent(path, &Agent{Model: "claude-sonnet-5"}); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if strings.Count(string(data), "agent:") != 1 || strings.Contains(string(data), "claude-code") {
		t.Errorf("replaced:\n%s", data)
	}
	if err := WriteAgent(path, nil); err != nil {
		t.Fatal(err)
	}
	if data, _ = os.ReadFile(path); string(data) != orig {
		t.Errorf("cleared:\n%s", data)
	}
	if err := WriteAgent(path, &Agent{Harness: "Not Valid"}); err == nil {
		t.Error("an invalid agent was written")
	}
}

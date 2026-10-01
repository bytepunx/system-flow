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

// S-0189: an agent's roles merge role by role and key by key, compare, check,
// print, and read back as they were written.
func TestAgentRoles(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Roles: map[string]Role{
		"explore": {Model: "haiku"},
		"verify":  {Model: "sonnet", Config: map[string]string{"effort": "medium"}},
	}}
	got := def.With(&Agent{Roles: map[string]Role{"verify": {Config: map[string]string{"effort": "high"}}, "plan": {Model: "opus"}}})
	if got.Roles["explore"].Model != "haiku" || got.Roles["verify"].Model != "sonnet" || got.Roles["verify"].Config["effort"] != "high" || got.Roles["plan"].Model != "opus" {
		t.Errorf("merged: %+v", got.Roles)
	}
	if def.Roles["verify"].Config["effort"] != "medium" {
		t.Error("the default's role was changed by a merge")
	}
	if def.Same(got) || !got.Same(def.With(got)) {
		t.Error("Same does not compare roles")
	}
	if (&Agent{Roles: map[string]Role{"verify": {Model: "sonnet"}}}).IsZero() {
		t.Error("an agent with only roles is not empty")
	}
	if s := got.String(); s != "claude-code, claude-opus-5-5; explore: haiku; plan: opus; verify: sonnet, effort=high" {
		t.Errorf("String: %s", s)
	}
	for _, bad := range []*Agent{
		{Roles: map[string]Role{"Verify": {Model: "sonnet"}}},
		{Roles: map[string]Role{"verify": {}}},
		{Roles: map[string]Role{"verify": {Model: "has space"}}},
		{Roles: map[string]Role{"verify": {Config: map[string]string{"Bad": "x"}}}},
	} {
		if err := bad.Validate(); err == nil {
			t.Errorf("accepted %+v", bad.Roles)
		}
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, []byte("version: 1\nname: Harbour\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteAgent(path, got); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Agent.Same(got) {
		data, _ := os.ReadFile(path)
		t.Errorf("read back %+v from:\n%s", m.Agent, data)
	}
}

func TestParseRoles(t *testing.T) {
	got, err := ParseRoles([]string{"verify=claude-code"}, []string{"explore=haiku", "verify = sonnet"}, []string{"verify.effort=medium"})
	if err != nil {
		t.Fatal(err)
	}
	if got["explore"].Model != "haiku" || got["verify"].Harness != "claude-code" || got["verify"].Model != "sonnet" || got["verify"].Config["effort"] != "medium" {
		t.Errorf("parsed: %+v", got)
	}
	if got, err := ParseRoles(nil, nil, nil); got != nil || err != nil {
		t.Errorf("nothing: %v %v", got, err)
	}
	for _, bad := range [][]string{{"", "verify", ""}, {"", "", "verify=x"}, {"", "", ".effort=x"}} {
		if _, err := ParseRoles(split(bad[0]), split(bad[1]), split(bad[2])); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
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

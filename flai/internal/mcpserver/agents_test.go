package mcpserver

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// S-0177, ADR-0064: an operator's own agent can have a story's agent started
// or restarted on the host, which judges whether it may.
func TestAStorysAgentIsStartedThroughTheServer(t *testing.T) {
	type call struct{ verb, root, story, by string }
	var got []call
	f := setupWith(t, func(o *Options) {
		o.Agents = func(_ context.Context, verb, root, story, by string) (AgentStarted, error) {
			got = append(got, call{verb, root, story, by})
			if verb == "start" {
				return AgentStarted{}, errors.New("the agent host action is off for this project: flai serve enable agent")
			}
			return AgentStarted{Story: story, Agent: "agent-" + story, Harness: "claude-code", PID: 42, Started: "2026-10-01T10:00:00Z"}, nil
		}
	})
	out, failed := f.call(t, "agent_restart", map[string]any{"story": "s-1"})
	if failed != "" || out["agent"] != "agent-S-0001" || out["pid"] != float64(42) {
		t.Fatalf("restart: %v %s", out, failed)
	}
	if _, failed := f.call(t, "agent_start", map[string]any{"story": "S-0001"}); !strings.Contains(failed, "the agent host action is off") {
		t.Errorf("a refusal reaches the agent as the tool's error: %q", failed)
	}
	root := projectRoot(f.repo)
	if len(got) != 2 || got[0] != (call{"restart", root, "S-0001", "claude"}) || got[1].verb != "start" {
		t.Errorf("asked: %+v", got)
	}
	if _, failed := f.call(t, "agent_restart", map[string]any{"story": "T-0001"}); !strings.Contains(failed, "not a story's ID") || len(got) != 2 {
		t.Errorf("a task has no agent: %q", failed)
	}
}

func TestAServerThatCannotStartAgentsSaysHow(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "agent_restart", map[string]any{"story": "S-0001"}); !strings.Contains(failed, "flai serve agent restart S-0001") {
		t.Errorf("without a starter: %q", failed)
	}
}

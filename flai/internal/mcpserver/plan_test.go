package mcpserver

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// S-0208: an operator's own agent can have the planner started on the host
// for an epic or a story, and the host judges whether it may.
func TestThePlannerIsStartedThroughTheServer(t *testing.T) {
	type call struct{ root, item, by string }
	var got []call
	f := setupWith(t, func(o *Options) {
		o.Plans = func(_ context.Context, root, item, by string) (PlanStarted, error) {
			got = append(got, call{root, item, by})
			switch item {
			case "E-0002":
				return PlanStarted{}, errors.New("the plan host action is off for this project: flai serve enable plan")
			case "S-0003":
				return PlanStarted{}, errors.New("the planner is already running for S-0003 (pid 7, started 2026-10-03T09:00:00Z); one item has one planner at a time")
			}
			return PlanStarted{Item: item, Agent: "planner-" + item, Harness: "claude-code", Command: "claude", PID: 42, Session: "abc", Started: "2026-10-03T10:00:00Z"}, nil
		}
	})
	out, failed := f.call(t, "plan", map[string]any{"id": "e-1"})
	if failed != "" || out["item"] != "E-0001" || out["agent"] != "planner-E-0001" || out["pid"] != float64(42) || out["session"] != "abc" {
		t.Fatalf("an epic: %v %s", out, failed)
	}
	if out, failed := f.call(t, "plan", map[string]any{"id": "S-0001"}); failed != "" || out["item"] != "S-0001" {
		t.Errorf("a story: %v %s", out, failed)
	}
	if _, failed := f.call(t, "plan", map[string]any{"id": "E-0002"}); !strings.Contains(failed, "the plan host action is off") {
		t.Errorf("a refusal reaches the agent as the tool's error: %q", failed)
	}
	if _, failed := f.call(t, "plan", map[string]any{"id": "S-0003"}); !strings.Contains(failed, "already running for S-0003") {
		t.Errorf("a second run: %q", failed)
	}
	root := projectRoot(f.repo)
	if len(got) != 4 || got[0] != (call{root, "E-0001", "claude"}) || got[1].item != "S-0001" {
		t.Errorf("asked: %+v", got)
	}
	for id, want := range map[string]string{
		"T-0001": "T-0001 is a task; the planner plans an epic or a story",
		"X-1":    `"X-1" is not an epic's or a story's ID`,
		"":       "the planner plans an E-nnnn or an S-nnnn",
	} {
		if _, failed := f.call(t, "plan", map[string]any{"id": id}); !strings.Contains(failed, want) {
			t.Errorf("%q: %q", id, failed)
		}
	}
	if len(got) != 4 {
		t.Errorf("an ID the planner does not plan reached the host: %+v", got)
	}
}

func TestAServerThatCannotStartThePlannerSaysHow(t *testing.T) {
	f := setup(t)
	if _, failed := f.call(t, "plan", map[string]any{"id": "e-16"}); !strings.Contains(failed, "this flai mcp cannot start the planner; on the host run flai plan E-0016") {
		t.Errorf("without a starter: %q", failed)
	}
}

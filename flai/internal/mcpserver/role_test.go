package mcpserver

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// S-0417: an operator's own agent can have a role the project defines
// started on the host, on an item or on none, with inputs, and the host
// judges whether it may, its refusals coming back as the tool's errors.
func TestADefinedRoleIsStartedThroughTheServer(t *testing.T) {
	t.Parallel()
	type call struct {
		root, name, item string
		inputs           []harness.Input
		by               string
	}
	var got []call
	f := setupWith(t, func(o *Options) {
		o.Roles = func(_ context.Context, root, name, item string, inputs []harness.Input, by string) (RoleStarted, error) {
			got = append(got, call{root, name, item, inputs, by})
			if name != "reviewer" {
				return RoleStarted{}, errors.New("the project defines no role " + name + "; the roles it defines are reviewer")
			}
			return RoleStarted{Role: name, Item: "S-0042", Action: "review", Agent: "reviewer-S-0042", Harness: "claude-code", Command: "claude", PID: 42, Session: "abc", Started: "2026-10-28T10:00:00Z", Trigger: "asked"}, nil
		}
	})
	out, failed := f.call(t, "role_run", map[string]any{"name": " reviewer ", "item": " S-42 ", "inputs": []map[string]any{{"name": "focus", "value": "docs only"}, {"name": "depth", "value": "2"}}})
	if failed != "" || out["role"] != "reviewer" || out["item"] != "S-0042" || out["action"] != "review" || out["agent"] != "reviewer-S-0042" || out["pid"] != float64(42) || out["trigger"] != "asked" {
		t.Fatalf("a role on an item: %v %s", out, failed)
	}
	if _, failed := f.call(t, "role_run", map[string]any{"name": "auditor"}); !strings.Contains(failed, "the project defines no role auditor; the roles it defines are reviewer") {
		t.Errorf("a refusal: %q", failed)
	}
	root := projectRoot(f.repo)
	want := []harness.Input{{Name: "focus", Value: "docs only"}, {Name: "depth", Value: "2"}}
	if len(got) != 2 || got[0].root != root || got[0].name != "reviewer" || got[0].item != "S-42" || !slices.Equal(got[0].inputs, want) || got[0].by != "claude" ||
		got[1].name != "auditor" || got[1].item != "" || len(got[1].inputs) != 0 {
		t.Errorf("asked: %+v", got)
	}
	if _, failed := f.call(t, "role_run", map[string]any{"name": "  "}); !strings.Contains(failed, "give name, the role's name") {
		t.Errorf("no name: %q", failed)
	}
	if _, failed := f.call(t, "role_run", map[string]any{"name": "reviewer", "inputs": []map[string]any{{"name": " ", "value": "docs"}}}); !strings.Contains(failed, `an input with the value "docs" has no name`) {
		t.Errorf("an input with no name: %q", failed)
	}
	if len(got) != 2 {
		t.Errorf("a call refused by the server reached the host: %+v", got)
	}
}

// S-0422: the agent of a formula's instance asks with its own name, which
// the host judges it by, and the run's trigger, formula <name> <S-nnnn>,
// comes back in the answer.
func TestAnInstancesAgentIsAnsweredWithTheFormulaTrigger(t *testing.T) {
	t.Parallel()
	var by string
	f := setupWith(t, func(o *Options) {
		o.Agent = "agent-S-0042"
		o.Roles = func(_ context.Context, _, name, item string, _ []harness.Input, asker string) (RoleStarted, error) {
			by = asker
			return RoleStarted{Role: name, Item: item, Action: "review", Agent: "reviewer-" + item, PID: 42, Trigger: "formula release-check " + item}, nil
		}
	})
	out, failed := f.call(t, "role_run", map[string]any{"name": "reviewer", "item": "S-0042", "inputs": []map[string]any{{"name": "prompt", "value": "Report the risk."}}})
	if failed != "" || out["trigger"] != "formula release-check S-0042" || by != "agent-S-0042" {
		t.Errorf("role_run by the instance's agent: %v %q, asked by %q", out, failed, by)
	}
	if !strings.Contains(roleRunDescription, "save the agent of a formula's instance") || !strings.Contains(roleRunDescription, "formula <name> <S-nnnn>") {
		t.Errorf("the tool's description does not say an instance's agent may start a role step's role")
	}
}

func TestAServerThatCannotStartARoleSaysHow(t *testing.T) {
	t.Parallel()
	f := setup(t)
	if _, failed := f.call(t, "role_run", map[string]any{"name": "reviewer", "item": "S-0042", "inputs": []map[string]any{{"name": "focus", "value": "docs"}}}); !strings.Contains(failed, "this flai mcp cannot start a role; on the host run flai role run reviewer S-0042 --input focus=docs") {
		t.Errorf("without a starter: %q", failed)
	}
}

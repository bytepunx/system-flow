package guard

import (
	"strings"
	"testing"
)

// S-0577: flai kluster up, down, deploy, teardown, and promote are the
// operator's (ADR-0254 §7). Every session flai serve starts, and every
// sub-agent, is refused them, however the command line runs flai, with a
// reason naming the shared infrastructure; flai kluster status reads and
// passes; the operator's own session runs them all.
func TestFlaiKlusterWritesAreTheOperators(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "verifier"; return e }
	sessions := []struct {
		name  string
		guard Guard
		sub   bool   // the calls are a sub-agent's
		who   string // how the refusal names the caller
		ask   string // what the refusal says to do instead
	}{
		{"story's agent", Guard{Commands: commands, Story: "S-0001", Served: true}, false, "the story's agent cannot", "thread_open"},
		{"story's agent's sub-agent", Guard{Commands: commands, Story: "S-0001", Served: true}, true, "a sub-agent (verifier) cannot", "final message"},
		{"planner", Guard{Commands: commands, Role: RolePlan}, false, "the planner cannot", "thread_open"},
		{"planner's sub-agent", Guard{Commands: commands, Role: RolePlan}, true, "a sub-agent (verifier) cannot", "final message"},
		{"orchestrator", Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}, false, "the orchestrator cannot", "thread_open"},
		{"analyzer", Guard{Commands: commands, Role: RoleAnalyze}, false, "the analyzer cannot", "thread_open"},
		{"session flai serve started", Guard{Commands: commands, Served: true}, false, "a session flai serve started cannot", "thread_open"},
		{"operator's sub-agent", Guard{Commands: commands}, true, "a sub-agent (verifier) cannot", "final message"},
	}
	writes := []Event{
		bash("", "flai kluster up"),
		bash("", "flai kluster down --yes"),
		bash("", "cd /w && scripts/flai.sh kluster down -y"),
		bash("", "/home/alex/.flai/bin/flai kluster up --json"),
		bash("", "bin/flai --config c.json kluster deploy S-0001"),
		bash("", "KUBECONFIG=x env flai kluster teardown S-0001"),
		bash("", "bash -c 'flai kluster promote'"),
		bash("", "flai board --json; flai kluster up"),
	}
	reads := []Event{
		bash("", "flai kluster status"),
		bash("", "scripts/flai.sh kluster status --json"),
		bash("", "flai kluster up --help"),
		bash("", "flai kluster --help"),
	}
	for _, s := range sessions {
		as := func(e Event) Event {
			if s.sub {
				return sub(e)
			}
			return e
		}
		for _, e := range writes {
			r := s.guard.Decide(as(e))
			if !strings.Contains(r.Why, s.who) || !strings.Contains(r.Why, "shared infrastructure on the operator's machine") || !strings.Contains(r.Why, s.ask) {
				t.Errorf("%s: %q: %q", s.name, e.ToolInput.Command, r.Why)
			}
			if r.Needs != "" {
				t.Errorf("%s: no permission allows it, yet the refusal needs %s", s.name, r.Needs)
			}
		}
		for _, e := range reads {
			r := s.guard.Decide(as(e))
			if strings.Contains(r.Why, "shared infrastructure on the operator's machine") {
				t.Errorf("%s: %q refused as a write: %s", s.name, e.ToolInput.Command, r.Why)
			}
		}
	}
	// a story's agent is let through the reads altogether
	for _, e := range reads {
		if why := (Guard{Commands: commands, Story: "S-0001", Served: true}).Check(e); why != "" {
			t.Errorf("the story's agent: %q refused: %s", e.ToolInput.Command, why)
		}
	}
	operator := Guard{Commands: commands}
	for _, e := range append(append([]Event{}, writes...), reads...) {
		if why := operator.Check(e); why != "" {
			t.Errorf("the operator's session: %q refused: %s", e.ToolInput.Command, why)
		}
	}
	// the orchestrator's refusal names the call, so that its activity
	// document logs it; a story's agent's does not
	if r := (Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}).Decide(bash("", "flai kluster up")); r.Call != "flai kluster up" {
		t.Errorf("the orchestrator's refusal: %+v", r)
	}
	if r := (Guard{Commands: commands, Story: "S-0001", Served: true}).Decide(bash("", "flai kluster up")); r.Call != "" {
		t.Errorf("a story's agent's refusal is logged as the orchestrator's: %+v", r)
	}
}

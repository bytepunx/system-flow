package guard

import (
	"slices"
	"strings"
	"testing"
)

// S-0295, ADR-0096: the shared paths are the operator's. Every session flai
// serve starts, and every sub-agent, is refused shared_paths_edit and flai
// shared add and remove, however the command line runs flai; each lists and
// checks them; the operator's own session changes them.
func TestTheSharedPathsAreTheOperatorsToChange(t *testing.T) {
	commands := append(slices.Clone(g.Commands), "shared")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "verifier"; return e }
	sessions := []struct {
		name  string
		guard Guard
		sub   bool   // the calls are a sub-agent's
		who   string // how the refusal names the caller
		ask   string // what the refusal says to do instead
	}{
		{"story's agent", Guard{Commands: commands, Story: "S-0001"}, false, "the story's agent cannot", "thread_open"},
		{"story's agent's sub-agent", Guard{Commands: commands, Story: "S-0001"}, true, "a sub-agent (verifier) cannot", "final message"},
		{"planner", Guard{Commands: commands, Role: RolePlan}, false, "the planner cannot", "thread_open"},
		{"planner's sub-agent", Guard{Commands: commands, Role: RolePlan}, true, "a sub-agent (verifier) cannot", "final message"},
		{"orchestrator", Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}, false, "the orchestrator cannot", "thread_open"},
		{"analyzer", Guard{Commands: commands, Role: RoleAnalyze}, false, "the analyzer cannot", "thread_open"},
		{"session flai serve started with another role", Guard{Commands: commands, Role: "verify"}, false, "a session flai serve started with the role verify cannot", "thread_open"},
		{"session flai serve started", Guard{Commands: commands, Served: true}, false, "a session flai serve started cannot", "thread_open"},
		{"operator's sub-agent", Guard{Commands: commands}, true, "a sub-agent (verifier) cannot", "final message"},
	}
	edits := []Event{
		{ToolName: MCPPrefix + SharedPathsEdit},
		bash("", "flai shared add design/adrs"),
		bash("", "flai shared remove 'docs/users/*.md'"),
		bash("", "cd /w && scripts/flai.sh shared add docs"),
		bash("", "FLAI_AGENT=x env flai --config c.json shared add docs"),
		bash("", "bash -c 'flai shared remove docs'"),
		bash("", "flai board --json; flai shared add docs"),
	}
	reads := []Event{
		{ToolName: MCPPrefix + "shared_paths"},
		bash("", "flai shared list --json"),
		bash("", "flai shared check S-0001 docs/users/flai.md --json"),
		bash("", "scripts/flai.sh shared check flai/cmd"),
		bash("", "flai shared add --help"),
	}
	for _, s := range sessions {
		as := func(e Event) Event {
			if s.sub {
				return sub(e)
			}
			return e
		}
		for _, e := range edits {
			r := s.guard.Decide(as(e))
			if !strings.Contains(r.Why, s.who) || !strings.Contains(r.Why, "ADR-0096") || !strings.Contains(r.Why, "operator's") || !strings.Contains(r.Why, s.ask) {
				t.Errorf("%s: %s %q: %q", s.name, e.ToolName, e.ToolInput.Command, r.Why)
			}
			if r.Needs != "" {
				t.Errorf("%s: no permission allows it, yet the refusal needs %s", s.name, r.Needs)
			}
		}
		for _, e := range reads {
			if why := s.guard.Check(as(e)); why != "" {
				t.Errorf("%s: %s %q refused: %s", s.name, e.ToolName, e.ToolInput.Command, why)
			}
		}
	}

	operator := Guard{Commands: commands}
	for _, e := range append(slices.Clone(edits), reads...) {
		if why := operator.Check(e); why != "" {
			t.Errorf("the operator's session: %s %q refused: %s", e.ToolName, e.ToolInput.Command, why)
		}
	}
}

// The orchestrator's refusal names the call, so that its activity document
// logs it, with no permission that would allow it.
func TestTheOrchestratorsSharedRefusalIsLogged(t *testing.T) {
	o := Guard{Commands: append(slices.Clone(g.Commands), "shared"), Role: RoleOrchestrate, Permissions: allOn}
	for _, c := range []struct {
		e    Event
		call string
	}{
		{Event{ToolName: MCPPrefix + SharedPathsEdit}, SharedPathsEdit},
		{bash("", "flai shared add 'docs/*.md'"), "flai shared add docs/*.md"},
	} {
		if r := o.Decide(c.e); r.Call != c.call || r.Needs != "" {
			t.Errorf("%s: %+v", c.e.ToolName, r)
		}
	}
	if r := (Guard{Commands: o.Commands, Story: "S-0001"}).Decide(Event{ToolName: MCPPrefix + SharedPathsEdit}); r.Call != "" {
		t.Errorf("a story's agent's refusal is logged as the orchestrator's: %+v", r)
	}
}

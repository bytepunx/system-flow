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

// klusterTool is a call of the MCP tool tool naming story, "" for none.
func klusterTool(tool, story string) Event {
	e := Event{ToolName: MCPPrefix + tool}
	e.ToolInput.Story = story
	return e
}

// ADR-0292 §1, §2: a story's agent, served or not, asks flai serve to deploy
// or tear down its own story's namespace, with flai kluster deploy or
// teardown --job in any flag form or the MCP tools kluster_deploy and
// kluster_teardown; a request for another story, or for none, is refused it
// with what to ask for instead.
func TestAStorysAgentRequestsItsOwnStorysNamespace(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	agents := []struct {
		name  string
		guard Guard
	}{
		{"served", Guard{Commands: commands, Story: "S-0001", Served: true}},
		{"not served", Guard{Commands: commands, Story: "S-0001"}},
		{"named agent-S-0001", Guard{Commands: commands, Agent: "agent-S-0001", Served: true}},
		{"named agent-S-1, not served", Guard{Commands: commands, Agent: "agent-S-1"}},
	}
	own := []string{
		"flai kluster deploy S-0001 --job",
		"flai kluster teardown S-0001 --job",
		"scripts/flai.sh kluster deploy --job S-1",
		"flai kluster deploy S-0001 --job=true --json",
		"flai --config c.json kluster deploy S-0001 --job --wait 5m",
		"flai kluster --json teardown --wait=5m --job s-0001",
		"flai kluster deploy S-0001 --job=false --job",
		"cd /w && /home/alex/.flai/bin/flai kluster deploy S-0001 --job=1 2>&1",
		"bash -c 'flai kluster teardown S-0001 --job --yes'",
	}
	others := []string{
		"flai kluster deploy S-0002 --job",
		"flai kluster teardown --job",
		"flai kluster deploy S-0001 S-0002 --job",
		"flai kluster deploy S-0001 --job; flai kluster teardown S-0002 --job",
	}
	for _, a := range agents {
		for _, line := range own {
			if why := a.guard.Check(bash("", line)); why != "" {
				t.Errorf("%s: %q refused: %s", a.name, line, why)
			}
		}
		for _, line := range others {
			why := a.guard.Check(bash("", line))
			for _, want := range []string{"the story's agent cannot", "own story's namespace alone, S-0001", "ADR-0292", "--job"} {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %q: %q wants %q", a.name, line, why, want)
				}
			}
		}
		for _, tool := range []string{KlusterDeploy, KlusterTeardown} {
			for _, story := range []string{"", "S-0001", "S-1"} {
				if why := a.guard.Check(klusterTool(tool, story)); why != "" {
					t.Errorf("%s: %s for %q refused: %s", a.name, tool, story, why)
				}
			}
			why := a.guard.Check(klusterTool(tool, "S-0002"))
			for _, want := range []string{"the story's agent cannot call " + tool + " for S-0002", "S-0001", "ADR-0292"} {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %s for S-0002: %q wants %q", a.name, tool, why, want)
				}
			}
		}
	}
	// the reason names the story asked for, or that none was
	if why := agents[0].guard.Check(bash("", "flai kluster teardown --job")); !strings.Contains(why, "names no single story") {
		t.Errorf("no story: %q", why)
	}
}

// ADR-0292 §2: every sub-agent is refused the request, by command and by
// tool, its own story's included, with the reason that the request is the
// story's agent's; and every other session flai serve started is refused the
// command as it is every flai kluster command that writes.
func TestOnlyTheStorysAgentRequestsANamespace(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "task"; return e }
	parents := []Guard{
		{Commands: commands, Story: "S-0001", Served: true},
		{Commands: commands, Story: "S-0001"},
		{Commands: commands, Agent: "agent-S-0001", Served: true},
		{Commands: commands},
		{Commands: commands, Role: RolePlan},
	}
	requests := []Event{
		bash("", "flai kluster deploy S-0001 --job"),
		bash("", "flai kluster teardown S-0001 --job --json"),
		bash("", "flai kluster deploy S-0002 --job=true"),
		klusterTool(KlusterDeploy, ""),
		klusterTool(KlusterTeardown, "S-0001"),
		klusterTool(KlusterDeploy, "S-0002"),
	}
	for i, p := range parents {
		for _, e := range requests {
			why := p.Check(sub(e))
			if !strings.Contains(why, "a sub-agent (task) cannot") || !strings.Contains(why, "is the story's agent's") || !strings.Contains(why, "final message") {
				t.Errorf("parent %d: %s %q: %q", i, e.ToolName, e.ToolInput.Command, why)
			}
		}
	}
	served := []struct {
		guard Guard
		who   string
	}{
		{Guard{Commands: commands, Role: RolePlan}, "the planner cannot"},
		{Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}, "the orchestrator cannot"},
		{Guard{Commands: commands, Served: true}, "a session flai serve started cannot"},
	}
	for _, s := range served {
		why := s.guard.Check(bash("", "flai kluster deploy S-0001 --job"))
		if !strings.Contains(why, s.who) || !strings.Contains(why, "shared infrastructure on the operator's machine") {
			t.Errorf("%s: %q", s.who, why)
		}
	}
	// the operator's own session runs the request, as it runs the writes
	for _, line := range []string{"flai kluster deploy S-0001 --job", "flai kluster teardown S-0002 --job"} {
		if why := (Guard{Commands: commands}).Check(bash("", line)); why != "" {
			t.Errorf("the operator's session: %q refused: %s", line, why)
		}
	}
}

// ADR-0292 §2: up, down, and promote in any form, and deploy and teardown
// without --job set true, stay refused to every story's agent and its
// sub-agents, however the flags are written, so that no form reaches a write
// as a request.
func TestNoFlagFormMakesAKlusterWriteARequest(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "task"; return e }
	writes := []string{
		"flai kluster deploy S-0001",
		"flai kluster teardown S-0001 --yes",
		"flai kluster deploy S-0001 --job=false",
		"flai kluster deploy S-0001 --job --job=false",
		"flai kluster teardown S-0001 --job=no",
		"flai kluster deploy S-0001 --json",
		"flai kluster up --job",
		"flai kluster up --job=true",
		"flai kluster down --job -y",
		"flai kluster promote --job",
		"flai kluster promote S-0001 --job",
		"flai --json kluster up",
		"flai -y kluster down --job",
		"flai --config=c.json kluster deploy S-0001",
		"flai kluster --wait 5m deploy S-0001 --job",
		"flai --wait 5m kluster deploy S-0001",
		"flai --wait 5m kluster teardown S-0001 --yes",
		"flai kluster --job deploy S-0001",
		"flai kluster deploy S-0001 -- --job",
		"flai kluster deploy S-0001 --job; flai kluster up",
		"bash -c 'flai kluster teardown S-0001'",
	}
	agents := []Guard{
		{Commands: commands, Story: "S-0001", Served: true},
		{Commands: commands, Story: "S-0001"},
		{Commands: commands, Agent: "agent-S-0001", Served: true},
	}
	for i, a := range agents {
		for _, line := range writes {
			why := a.Check(bash("", line))
			if !strings.Contains(why, "the story's agent cannot") || !strings.Contains(why, "shared infrastructure on the operator's machine") || !strings.Contains(why, "thread_open") {
				t.Errorf("agent %d: %q: %q", i, line, why)
			}
			why = a.Check(sub(bash("", line)))
			if !strings.Contains(why, "a sub-agent (task) cannot") || !strings.Contains(why, "shared infrastructure on the operator's machine") && !strings.Contains(why, "is the story's agent's") {
				t.Errorf("agent %d's sub-agent: %q: %q", i, line, why)
			}
		}
	}
	// a line the shell could expand is no request: the guard cannot tell
	// what flai is given
	expanded := []string{
		"flai kluster deploy S-0001 --job $EXTRA",
		"flai kluster deploy S-0001 --job --wait $W",
		"flai kluster deploy S-0001 --job{,=false}",
		"flai kluster deploy \"$FLAI_STORY\" --job",
		"flai kluster deploy S-0001 --job=`echo false`",
		"flai kluster deploy S-0001 --job --wait=$(echo 5m --job=false)",
		"flai kluster teardown S-0001 --job --jo\\b=false",
		"flai kluster deploy S-0001 --job ./--job=*",
	}
	for i, a := range agents {
		for _, line := range expanded {
			if why := a.Check(bash("", line)); !strings.Contains(why, "the story's agent cannot") {
				t.Errorf("agent %d: %q: %q", i, line, why)
			}
		}
	}
	for _, line := range expanded[:4] {
		if why := agents[0].Check(bash("", line)); !strings.Contains(why, "written plainly") || !strings.Contains(why, "flai kluster deploy S-0001 --job") {
			t.Errorf("%q: %q", line, why)
		}
	}
	// status and deployments read, and stay as they were
	for _, line := range []string{"flai kluster status", "flai kluster deployments --json"} {
		if why := agents[0].Check(bash("", line)); why != "" {
			t.Errorf("%q refused: %s", line, why)
		}
	}
}

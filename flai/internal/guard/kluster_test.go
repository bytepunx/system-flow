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

// ADR-0292 §2: a story's agent and its sub-agents run flai kluster exec and
// log for their own story's namespace, however the line runs flai; another
// story's, the default one, none, one not written plainly, and one flai could
// read as another subcommand are refused them, with the form to write
// instead. The operator's own session runs them all.
func TestAStorysAgentReachesItsOwnNamespace(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "task"; return e }
	agents := []struct {
		name  string
		guard Guard
		sub   bool
		who   string
	}{
		{"served", Guard{Commands: commands, Story: "S-0001", Served: true}, false, "the story's agent cannot"},
		{"not served", Guard{Commands: commands, Story: "S-0001"}, false, "the story's agent cannot"},
		{"named agent-S-0001", Guard{Commands: commands, Agent: "agent-S-0001", Served: true}, false, "the story's agent cannot"},
		{"named agent-S-1, not served", Guard{Commands: commands, Agent: "agent-S-1"}, false, "the story's agent cannot"},
		{"served's sub-agent", Guard{Commands: commands, Story: "S-0001", Served: true}, true, "a sub-agent (task) cannot"},
		{"named's sub-agent", Guard{Commands: commands, Agent: "agent-S-0001"}, true, "a sub-agent (task) cannot"},
		{"read-only sub-agent", Guard{Commands: commands, Story: "S-0001", Served: true, ReadOnly: func(string) bool { return true }}, true, "a sub-agent (task) cannot"},
	}
	own := []string{
		"flai kluster exec S-0001 -- kubectl get pods",
		"flai kluster exec s-1 -- kubectl logs deploy/api",
		"scripts/flai.sh kluster exec S-0001 -- kubectl get svc",
		"flai --config c.json kluster exec S-0001 -- sh -c 'kubectl get pods -n $FLAI_KLUSTER_NAMESPACE'",
		"cd /w && /home/alex/.flai/bin/flai kluster exec S-0001 -- kubectl get deploy 2>&1",
		"flai kluster exec S-0001 -- kluster status",
		"flai kluster log S-0001 --step deploy",
		"flai kluster log --step checks S-0001",
		"flai kluster log S-0001 --step=build --json",
		"bash -c 'flai kluster log S-0001 --step deploy'",
		"KUBECONFIG=.flai-cache/kluster/agents/s-0001.kubeconfig kubectl get pods",
	}
	others := []struct{ line, names string }{
		{"flai kluster exec S-0002 -- kubectl get pods", "names S-0002, not this story's namespace"},
		{"flai kluster exec default -- kubectl get pods", "names default, not this story's namespace"},
		{"flai kluster exec -- kubectl get pods", "names no single namespace"},
		{"flai kluster log S-0002 --step deploy", "names S-0002"},
		{"flai kluster log --step S-0001 S-0002", "names S-0002"},
		{"flai kluster log S-0001 S-0002 --step deploy", "names no single namespace"},
		{"flai kluster exec S-0001 -- true; flai kluster exec S-0002 -- true", "names S-0002"},
		{"flai kluster exec S-0001 -- flai kluster exec S-0002 -- true", "names S-0002"},
		{"bash -c 'flai kluster log S-0002 --step deploy'", "names S-0002"},
		{"flai kluster exec $FLAI_STORY -- ls", "written plainly"},
		{"flai kluster log \"$S\" --step deploy", "written plainly"},
		{"flai kluster exec S-0001 --wait `echo 1 S-0002` -- ls", "written plainly"},
		{"flai kluster --config exec log S-0001 --step deploy", "another subcommand"},
	}
	for _, a := range agents {
		as := func(e Event) Event {
			if a.sub {
				return sub(e)
			}
			return e
		}
		for _, line := range own {
			if why := a.guard.Check(as(bash("", line))); why != "" {
				t.Errorf("%s: %q refused: %s", a.name, line, why)
			}
		}
		for _, o := range others {
			why := a.guard.Check(as(bash("", o.line)))
			for _, want := range []string{a.who, "a story's agent's and its sub-agents', for their own story's namespace alone", "ADR-0292", o.names, "Write it as flai kluster "} {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %q: %q wants %q", a.name, o.line, why, want)
				}
			}
		}
	}
	// the form to write names the session's own story and the command
	why := agents[0].guard.Check(bash("", "flai kluster log S-0002 --step deploy"))
	if !strings.Contains(why, "Write it as flai kluster log S-0001 --step <step>") {
		t.Errorf("log: %q", why)
	}
	why = agents[0].guard.Check(bash("", "flai kluster exec S-0002 -- ls"))
	if !strings.Contains(why, "Write it as flai kluster exec S-0001 -- <command>") {
		t.Errorf("exec: %q", why)
	}
	// a read-only sub-agent's exec writes no file but under the temporary
	// folder, as its other commands do (ADR-0141)
	readOnly := agents[len(agents)-1].guard
	for _, line := range []string{"flai kluster exec S-0001 -- rm -rf build", "flai kluster exec S-0001 -- bash -c 'touch x'"} {
		if why := readOnly.Check(sub(bash("", line))); !strings.Contains(why, "writes no file but under the temporary folder") {
			t.Errorf("read-only sub-agent: %q: %q", line, why)
		}
	}
	// a session with no story of its own is refused them all, its own
	// sub-agents too, whichever story they name
	none := []struct {
		name  string
		guard Guard
		sub   bool
		who   string
		ask   string
	}{
		{"planner", Guard{Commands: commands, Role: RolePlan}, false, "the planner cannot", "thread_open"},
		{"orchestrator", Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}, false, "the orchestrator cannot", "thread_open"},
		{"analyzer", Guard{Commands: commands, Role: RoleAnalyze}, false, "the analyzer cannot", "thread_open"},
		{"session flai serve started", Guard{Commands: commands, Served: true}, false, "a session flai serve started cannot", "thread_open"},
		{"planner's sub-agent", Guard{Commands: commands, Role: RolePlan}, true, "a sub-agent (task) cannot", "final message"},
		{"operator's sub-agent", Guard{Commands: commands}, true, "a sub-agent (task) cannot", "final message"},
	}
	for _, n := range none {
		for _, line := range []string{"flai kluster exec S-0001 -- kubectl get pods", "flai kluster log S-0001 --step deploy", "flai kluster exec default -- ls"} {
			e := bash("", line)
			if n.sub {
				e = sub(e)
			}
			why := n.guard.Check(e)
			for _, want := range []string{n.who, "for their own story's namespace alone", "no story of its own", n.ask} {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %q: %q wants %q", n.name, line, why, want)
				}
			}
		}
	}
	operator := Guard{Commands: commands}
	for _, line := range own {
		if why := operator.Check(bash("", line)); why != "" {
			t.Errorf("the operator's session: %q refused: %s", line, why)
		}
	}
	for _, o := range others {
		if why := operator.Check(bash("", o.line)); why != "" {
			t.Errorf("the operator's session: %q refused: %s", o.line, why)
		}
	}
}

// ADR-0292 §2: every agent and sub-agent is refused the kluster program but
// kluster status, in every form a line runs it, and a line naming flai's
// kubeconfig for the cluster, whole or as the end of a path, each with what to
// run instead; flai kluster, a word naming kluster that runs nothing, and the
// namespace's own kubeconfig pass. The operator's own session runs them all.
func TestNoAgentRunsTheKlusterProgramOrFlaisKubeconfig(t *testing.T) {
	commands := append(append([]string{}, g.Commands...), "kluster")
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "task"; return e }
	programs := []string{
		"kluster up",
		"kluster",
		"kluster --json status",
		"kluster delete flai",
		"/usr/local/bin/kluster down",
		"cd /w && kluster list",
		"flai board; kluster up",
		"bash -c 'kluster up'",
		"eval kluster down",
		"KUBECONFIG=x kluster up",
		"sudo -u alex kluster up",
		"timeout 5 kluster stop",
		"echo flai | xargs kluster up",
		"flai kluster exec S-0001 -- kluster up",
		"flai kluster exec S-0001 -- bash -c 'kluster down'",
	}
	kubeconfig := []string{
		"cat .flai-cache/kluster/kubeconfig",
		"KUBECONFIG=.flai-cache/kluster/kubeconfig kubectl get pods -A",
		"kubectl --kubeconfig=/home/alex/p/.flai-cache/kluster/kubeconfig get ns",
		"kubectl --kubeconfig \"/home/alex/p/.flai-cache/kluster/kubeconfig\" get ns",
		"export KUBECONFIG=$PWD/.flai-cache/kluster/kubeconfig",
		"KUBECONFIG=a.yaml:.flai-cache/kluster/kubeconfig kubectl get ns",
		"cp ./.flai-cache/kluster/kubeconfig /tmp/k",
		"cat .flai-cache/kluster/agents/../kubeconfig",
		"bash -c 'cat .flai-cache//kluster/kubeconfig | head'",
		"kubectl get ns > /tmp/x 2>.flai-cache/kluster/kubeconfig",
		"flai kluster exec S-0001 -- kubectl --kubeconfig .flai-cache/kluster/kubeconfig get ns",
	}
	pass := []string{
		"kluster status",
		"kluster status --json",
		"flai kluster status",
		"scripts/flai.sh kluster status --json",
		"cat design/system/kluster.md",
		"ls .flai-cache/kluster",
		"grep -n kluster flai/cmd/kluster.go",
		"echo kluster up",
		"cat .flai-cache/kluster/agents/s-0001.kubeconfig",
		"ls .flai-cache/kluster/agents/",
	}
	sessions := []struct {
		name  string
		guard Guard
		sub   bool
		who   string
		story string // the story the refusal's exec names
	}{
		{"story's agent", Guard{Commands: commands, Story: "S-0001", Served: true}, false, "the story's agent cannot", "S-0001"},
		{"story's agent, not served", Guard{Commands: commands, Story: "S-0001"}, false, "the story's agent cannot", "S-0001"},
		{"named agent-S-0001, not served", Guard{Commands: commands, Agent: "agent-S-0001"}, false, "the story's agent cannot", "S-0001"},
		{"story's agent's sub-agent", Guard{Commands: commands, Story: "S-0001", Served: true}, true, "a sub-agent (task) cannot", "S-0001"},
		{"planner", Guard{Commands: commands, Role: RolePlan}, false, "the planner cannot", "S-nnnn"},
		{"orchestrator", Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}, false, "the orchestrator cannot", "S-nnnn"},
		{"session flai serve started", Guard{Commands: commands, Served: true}, false, "a session flai serve started cannot", "S-nnnn"},
		{"operator's sub-agent", Guard{Commands: commands}, true, "a sub-agent (task) cannot", "S-nnnn"},
	}
	for _, s := range sessions {
		as := func(e Event) Event {
			if s.sub {
				return sub(e)
			}
			return e
		}
		// a session with no story of its own is refused flai kluster exec
		// before what it runs is read
		reason := func(line string, wants ...string) []string {
			if s.story == "S-nnnn" && strings.HasPrefix(line, "flai kluster exec") {
				return []string{s.who, "no story of its own"}
			}
			return append([]string{s.who}, wants...)
		}
		for _, line := range programs {
			why := s.guard.Check(as(bash("", line)))
			for _, want := range reason(line, "the project's cluster is flai's", "only as kluster status", "flai kluster status", "flai kluster exec "+s.story+" -- <command>", "ADR-0292") {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %q: %q wants %q", s.name, line, why, want)
				}
			}
		}
		for _, line := range kubeconfig {
			why := s.guard.Check(as(bash("", line)))
			for _, want := range reason(line, "is flai's kubeconfig for the whole cluster", "flai kluster exec "+s.story+" -- <command>", "namespace's own kubeconfig", "ADR-0292") {
				if !strings.Contains(why, want) {
					t.Errorf("%s: %q: %q wants %q", s.name, line, why, want)
				}
			}
		}
	}
	// the story's agents and their sub-agents run the rest
	for _, s := range sessions[:4] {
		for _, line := range pass {
			e := bash("", line)
			if s.sub {
				e = sub(e)
			}
			if why := s.guard.Check(e); strings.Contains(why, "flai's kubeconfig") || strings.Contains(why, "only as kluster status") {
				t.Errorf("%s: %q refused: %s", s.name, line, why)
			}
		}
	}
	for _, line := range pass {
		if why := sessions[0].guard.Check(bash("", line)); why != "" {
			t.Errorf("the story's agent: %q refused: %s", line, why)
		}
	}
	// the orchestrator's refusal names the call, so that its activity
	// document logs it
	if r := (Guard{Commands: commands, Role: RoleOrchestrate, Permissions: allOn}).Decide(bash("", "flai board; kluster up")); r.Call != "kluster up" {
		t.Errorf("the orchestrator's refusal: %+v", r)
	}
	operator := Guard{Commands: commands}
	for _, line := range append(append(append([]string{}, programs...), kubeconfig...), pass...) {
		if why := operator.Check(bash("", line)); why != "" {
			t.Errorf("the operator's session: %q refused: %s", line, why)
		}
	}
}

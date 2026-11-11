package guard

import (
	"slices"
	"strings"
	"testing"
)

// hostCommands are the test guard's commands with those that write the host
// configuration.
var hostCommands = append(slices.Clone(g.Commands), "serve", "config", "template")

// hostWritten is what each refusal of a write of the host configuration says.
const hostWritten = "only the operator changes it"

// I-0210: a story's agent, whose role has no write rule of its own, set the
// reminder channel in the operator's host configuration.
func TestAStorysAgentCannotSetTheReminderChannel(t *testing.T) {
	gs := Guard{Commands: hostCommands, Story: "S-0001", Served: true}
	v := Decide(gs, Call{Story: "S-0001", Kind: Shell, Command: "flai serve notify set --url x"})
	if v.Outcome != Deny || !strings.Contains(v.Reason, "the story's agent cannot") || !strings.Contains(v.Reason, "ADR-0220") || !strings.Contains(v.Reason, "thread_open") {
		t.Errorf("the story's agent set the reminder channel: %+v", v)
	}
}

// S-0511, ADR-0220: no session flai serve starts, and no sub-agent, writes
// the host configuration, in any form a shell line takes; its reads are left
// to the other rules, and the operator's own session is held by none.
func TestTheHostConfigurationIsTheOperators(t *testing.T) {
	sub := func(e Event) Event { e.AgentID, e.AgentType = "a1", "verifier"; return e }
	served := []struct {
		name  string
		guard Guard
	}{
		{"story's agent", Guard{Commands: hostCommands, Story: "S-0001", Served: true}},
		{"planner", Guard{Commands: hostCommands, Role: RolePlan}},
		{"orchestrator", Guard{Commands: hostCommands, Role: RoleOrchestrate, Permissions: allOn}},
		{"analyzer", Guard{Commands: hostCommands, Role: RoleAnalyze}},
		{"curator", Guard{Commands: hostCommands, Role: RoleCurate}},
		{"session flai serve started", Guard{Commands: hostCommands, Served: true}},
	}
	type session struct {
		name  string
		guard Guard
		sub   bool   // the calls are a sub-agent's
		who   string // how the refusal names the caller
		ask   string // what the refusal says to do instead
	}
	whos := map[string]string{"story's agent": "the story's agent", "planner": "the planner", "orchestrator": "the orchestrator", "analyzer": "the analyzer", "curator": "the curator", "session flai serve started": "a session flai serve started"}
	var sessions []session
	for _, s := range served {
		sessions = append(sessions,
			session{s.name, s.guard, false, whos[s.name] + " cannot", "thread_open"},
			session{s.name + "'s sub-agent", s.guard, true, "a sub-agent (verifier) cannot", "final message"})
	}
	sessions = append(sessions, session{"operator's sub-agent", Guard{Commands: hostCommands}, true, "a sub-agent (verifier) cannot", "final message"})
	writes := []string{
		"serve enable auto-approve",
		"serve disable push --all-projects",
		"serve agent set -- /bin/sh -c x",
		"serve agent set --name me",
		"serve agent clear",
		"serve agent harness claude-code --program /tmp/claude",
		"serve agent harness claude-code --program=/tmp/claude",
		"serve agent harness command --guard none",
		"serve agent harness command --deny-protected",
		"serve agent harness command --deny-protected=false",
		"serve agent harness claude-code -- --permission-mode bypassPermissions",
		"serve agent harness claude-code --reset",
		// S-0680: the container harness's settings write it too
		"serve agent harness container --image claude-code=x:1",
		"serve agent harness container --image=claude-code=x:1",
		"serve agent harness container --network flai",
		"serve agent harness container --login /home/x/.claude",
		"serve agent harness container --login=",
		"serve agent harness claude-code --outside /home/x/src/other",
		"serve agent harness claude-code --outside=",
		"serve agent provider litellm --api anthropic-messages --base-url http://h --key-env K",
		"serve agent provider litellm --base-url=http://h",
		"serve agent provider litellm --key-env K",
		"serve agent provider litellm --spend-key-env M",
		"serve agent provider litellm --model opus=x",
		"serve agent provider litellm --reset",
		// S-0720: a clear of the budget state lifts every project's hold
		"serve budget clear",
		"serve budget clear 'claude-code meclod'",
		"serve checks set --name lint -- make lint",
		"serve checks timeout 30",
		"serve checks clear lint",
		"serve components idle 5m",
		"serve components max-servers 3",
		"serve forge github --token-env T",
		"serve forge github --token-env=T",
		"serve forge github --reset",
		"serve notify set --url x",
		"serve notify set --command x",
		"serve notify clear",
		"serve mcp bind 0.0.0.0",
		"serve mcp clear",
		"serve import add /tmp/repos",
		"serve import remove /tmp/repos",
		"config set agent.command x",
		"template use github.com/x/y",
	}
	forms := []func(string) string{
		func(w string) string { return "flai " + w },
		func(w string) string { return "cd /w && scripts/flai.sh " + w },
		func(w string) string { return "/home/alex/.flai/bin/flai " + w + " --json" },
		func(w string) string { return "bin/flai --config c.json " + w },
		func(w string) string { return "FLAI_CONFIG=x env flai " + w },
		func(w string) string { return "bash -c 'flai " + w + "'" },
		func(w string) string { return "flai serve status; flai " + w },
		// I-0230, S-0613: a redirection, wherever it stands, is no word of
		// the command, so one before the subcommand is not read as it
		func(w string) string { return "flai " + w + " 2>&1" },
		func(w string) string { return "flai 2>/dev/null " + w },
		func(w string) string { return "scripts/flai.sh >/dev/null 2>&1 " + w },
	}
	reads := []string{
		"flai serve status",
		"flai serve actions",
		"flai serve journal -n 5",
		"flai serve agent",
		"flai serve agent show",
		"flai serve agent harness",
		"flai serve agent harness claude-code",
		"flai serve agent provider",
		"flai serve agent provider litellm",
		"flai serve agent provider --json",
		"flai serve budget",
		"flai serve budget --json",
		"flai serve checks show",
		"flai serve components show",
		"flai serve forge",
		"flai serve forge github",
		"scripts/flai.sh serve forge --json",
		"flai serve notify show",
		"flai serve mcp show",
		"flai serve mcp",
		"flai serve import list",
		"flai serve agent start S-0001",
		"flai serve stop",
		"flai config get",
		"flai config path",
		"flai template show",
		"flai serve notify set --help",
	}
	for _, s := range sessions {
		as := func(e Event) Event {
			if s.sub {
				return sub(e)
			}
			return e
		}
		for _, w := range writes {
			for _, form := range forms {
				line := form(w)
				r := s.guard.Decide(as(bash("", line)))
				if !strings.Contains(r.Why, s.who) || !strings.Contains(r.Why, "host configuration") || !strings.Contains(r.Why, hostWritten) || !strings.Contains(r.Why, "ADR-0220") || !strings.Contains(r.Why, s.ask) {
					t.Errorf("%s: %q: %q", s.name, line, r.Why)
				}
				if r.Needs != "" {
					t.Errorf("%s: no permission allows it, yet the refusal needs %s", s.name, r.Needs)
				}
			}
		}
		for _, line := range reads {
			if r := s.guard.Decide(as(bash("", line))); strings.Contains(r.Why, hostWritten) {
				t.Errorf("%s: %q refused as a write of the host configuration: %s", s.name, line, r.Why)
			}
		}
	}
	// a story's agent and a session with no role, which run flai freely, are
	// let through the reads altogether
	for _, gs := range []Guard{served[0].guard, served[5].guard} {
		for _, line := range reads {
			if why := gs.Check(bash("", line)); why != "" {
				t.Errorf("%+v: %q refused: %s", gs, line, why)
			}
		}
	}
	operator := Guard{Commands: hostCommands}
	for _, w := range writes {
		for _, form := range forms {
			if why := operator.Check(bash("", form(w))); why != "" {
				t.Errorf("the operator's session: %q refused: %s", form(w), why)
			}
		}
	}
	for _, line := range reads {
		if why := operator.Check(bash("", line)); why != "" {
			t.Errorf("the operator's session: %q refused: %s", line, why)
		}
	}
	// the orchestrator's refusal names the call, so that its activity
	// document logs it; a story's agent's does not
	if r := served[2].guard.Decide(bash("", "flai serve enable auto-approve")); r.Call != "flai serve enable auto-approve" {
		t.Errorf("the orchestrator's refusal: %+v", r)
	}
	if r := served[0].guard.Decide(bash("", "flai serve enable auto-approve")); r.Call != "" {
		t.Errorf("a story's agent's refusal is logged as the orchestrator's: %+v", r)
	}
}

// I-0318, S-0680: the forms of flai serve that ADR-0220 leaves open as reads
// pass for every agent the guard holds to its reads: the orchestrator, under
// its rules in Go and its baseline file, the planner, the analyzer, the
// curator, and a story's sub-agents, read-only or not, with a redirection or
// without. Main's guard, with no serve entry among its reads, refused them
// all. The forms that run, start, stop, or commit, and those that write the
// host configuration, stay refused to each.
func TestEveryAgentReadsTheReadFormsOfFlaiServe(t *testing.T) {
	withServe := func(gr Guard) Guard {
		gr.Commands = append(slices.Clone(gr.Commands), "serve", "config", "template")
		return gr
	}
	readOnly := func(agentType string) bool { return agentType == "verifier" }
	type agent struct {
		name  string
		guard Guard
		as    string // the sub-agent's type, "" for the session's own call
	}
	story := withServe(Guard{Commands: g.Commands, Story: "S-0001", Served: true, ReadOnly: readOnly})
	agents := []agent{
		{"orchestrator, its rules in Go", withServe(orchestrator(allOn)), ""},
		{"orchestrator, its baseline file", withServe(orchestrating(allOn, nil)), ""},
		{"planner", withServe(planGuard), ""},
		{"planner, its baseline file", withServe(Guard{Commands: g.Commands, Role: RolePlan, Design: "design"}), ""},
		{"analyzer", withServe(Guard{Commands: g.Commands, Role: RoleAnalyze}), ""},
		{"analyzer, its baseline file", withServe(Guard{Commands: g.Commands, Role: RoleAnalyze, Design: "design"}), ""},
		{"curator", withServe(Guard{Commands: g.Commands, Role: RoleCurate}), ""},
		{"curator, its baseline file", withServe(Guard{Commands: g.Commands, Role: RoleCurate, Design: "design"}), ""},
		{"a read-only sub-agent", story, "verifier"},
		{"a task sub-agent", story, "general-purpose"},
	}
	reads := []string{
		"flai serve status",
		"flai serve status --json",
		"flai serve status 2>&1",
		"cd /w && scripts/flai.sh serve status >/dev/null 2>&1",
		"flai serve actions",
		"flai serve journal",
		"flai serve journal -n 5",
		"flai serve journal --last 5 --json",
		"flai serve agent",
		"flai serve agent show",
		"flai serve agent harness",
		"flai serve agent harness claude-code",
		"flai serve agent harness claude-code --project /w --json",
		"flai serve agent provider",
		"flai serve agent provider litellm",
		"flai serve agent provider --json",
		"flai serve checks show",
		"flai serve components show",
		"flai serve forge",
		"flai serve forge github",
		"scripts/flai.sh serve forge --json",
		"flai serve notify show",
		"flai serve mcp",
		"flai serve mcp show",
		"flai serve import list",
	}
	refused := []string{
		"flai serve",
		"flai serve --json",
		"flai serve start",
		"flai serve stop",
		"flai serve status; flai serve stop",
		"flai serve agent start S-0001",
		"flai serve agent restart S-0001",
		"flai serve agent stop S-0001",
		"flai serve agent commit S-0001",
		"flai serve agent stream S-0001",
		"flai serve agent set --name me",
		"flai serve agent harness container --image claude-code=x:1",
		"flai serve agent harness container --network=flai",
		"flai serve agent harness claude-code -- --verbose",
		"flai serve agent provider litellm --key-env K",
		"flai serve orchestrate start",
		"flai serve project add /w",
		"flai serve usage S-0001",
		"flai serve settings",
		"flai serve checks set --name lint -- make lint",
		"flai serve components idle 5m",
		"flai serve forge github --token-env T",
		"flai serve notify set --url x",
		"flai serve mcp bind 0.0.0.0",
		"flai serve import add /w",
		"flai serve enable auto-approve",
	}
	for _, a := range agents {
		event := func(line string) Event { return bash(a.as, line) }
		for _, line := range reads {
			if r := a.guard.Decide(event(line)); r.Why != "" {
				t.Errorf("%s: %q refused: %s", a.name, line, r.Why)
			}
		}
		for _, line := range refused {
			if r := a.guard.Decide(event(line)); r.Why == "" {
				t.Errorf("%s: %q passes", a.name, line)
			}
		}
	}
}

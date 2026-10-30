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
		"serve agent provider litellm --api anthropic-messages --base-url http://h --key-env K",
		"serve agent provider litellm --base-url=http://h",
		"serve agent provider litellm --key-env K",
		"serve agent provider litellm --spend-key-env M",
		"serve agent provider litellm --model opus=x",
		"serve agent provider litellm --reset",
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
		"flai serve checks show",
		"flai serve components show",
		"flai serve forge",
		"flai serve forge github",
		"scripts/flai.sh serve forge --json",
		"flai serve notify show",
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

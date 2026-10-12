package guard

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// S-0353: Decide holds every caller's rules on a Call that names no harness:
// its kind, its bare tool, whether a sub-agent makes it, and the role and
// story it carries, rather than the guard's, decide.
func TestDecideHoldsEachCallerOnANeutralCall(t *testing.T) {
	running := Guard{Commands: g.Commands, Running: func() ([]SubAgent, error) {
		return []SubAgent{{ID: "a1", Type: "explorer"}}, nil
	}}
	off := Guard{Commands: g.Commands, Permissions: except(manifest.PermitPromoteToReady)}
	sub := func(c Call) Call { c.Subagent, c.Agent = true, "verifier"; return c }
	for _, tc := range []struct {
		name  string
		guard Guard
		call  Call
		want  Outcome
		why   string // what the reason says, for a Deny
		needs string
	}{
		{"a sub-agent's write through flai", g, sub(Call{Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: "review"}}), Deny, "a sub-agent (verifier) cannot call item_move: it reads", ""},
		{"a sub-agent's read through flai", g, sub(Call{Kind: FlaiTool, Tool: "item_get", Args: Args{ID: "S-0001"}}), Allow, "", ""},
		{"a sub-agent's ADR without commit", g, sub(Call{Kind: FlaiTool, Tool: ADRNew}), Allow, "", ""},
		{"a sub-agent's ADR with commit", g, sub(Call{Kind: FlaiTool, Tool: ADRNew, Args: Args{Commit: true}}), Deny, "cannot call adr_new with commit true", ""},
		{"a sub-agent's commit", g, sub(Call{Kind: Shell, Command: "git add -A && git commit -m x"}), Deny, `a sub-agent (verifier) cannot run "git add -A"`, ""},
		{"a sub-agent's move through flai's command line", g, sub(Call{Kind: Shell, Command: "flai move S-0001 review"}), Deny, "flai commands that change work items", ""},
		{"a sub-agent's test run", g, sub(Call{Kind: Shell, Command: "scripts/flai.sh test flai/internal/guard"}), Allow, "", ""},
		{"a sub-agent's protected write", g, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/.claude/settings.json"}), Deny, "a sub-agent (verifier) cannot use write_file on /w/.claude/settings.json", ""},
		{"a sub-agent's protected write with auto-approve on", Guard{Commands: g.Commands, AutoApprove: true}, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/.claude/settings.json"}), Allow, "", ""},
		{"a sub-agent's write of code", g, sub(Call{Kind: FileWrite, Name: "write_file", Path: "/w/flai/x.go"}), Allow, "", ""},
		{"the story's agent's wait while a sub-agent runs", running, Call{Story: "S-0001", Kind: FlaiTool, Tool: "wait_for_events"}, Deny, "the story's agent cannot call wait_for_events while its sub-agents run (explorer a1) and no thread on S-0001", ""},
		{"a wait outside a story", running, Call{Kind: FlaiTool, Tool: "wait_for_events"}, Allow, "", ""},
		{"the planner's file write", g, Call{Role: RolePlan, Kind: FileWrite, Name: "write_file", Path: "flai/x.go"}, Deny, "the planner cannot use write_file: it plans through flai", ""},
		{"the planner's move to backlog", g, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Backlog}}, Allow, "", ""},
		{"the planner's story that is not a draft", g, Call{Role: RolePlan, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story"}}, Deny, "the planner cannot create a story that is not a draft", ""},
		{"the orchestrator's promotion, its permission off", off, Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Ready}}, Deny, "the orchestrator cannot move S-0001 to ready: it needs orchestration.permissions." + manifest.PermitPromoteToReady + ", which is off", "orchestration.permissions." + manifest.PermitPromoteToReady},
		{"the orchestrator's promotion, its permission on", Guard{Commands: g.Commands, Permissions: allOn}, Call{Role: RoleOrchestrate, Kind: FlaiTool, Tool: "item_move", Args: Args{ID: "S-0001", To: Ready}}, Allow, "", ""},
		{"the orchestrator's file write", off, Call{Role: RoleOrchestrate, Kind: FileWrite, Name: "write_file", Path: "x.md"}, Deny, "the orchestrator cannot use write_file: the orchestrator never does it", ""},
		{"the analyzer's story", g, Call{Role: RoleAnalyze, Kind: FlaiTool, Tool: "item_new", Args: Args{Type: "story"}}, Deny, "the analyzer cannot call item_new: the analyzer authors no stories", ""},
		{"the analyzer's commit", g, Call{Role: RoleAnalyze, Kind: Shell, Command: "git commit -m x"}, Deny, `the analyzer cannot run "git commit -m x": it runs only git's reads`, ""},
		{"the analyzer's issue", g, Call{Role: RoleAnalyze, Kind: FlaiTool, Tool: "issue_new"}, Allow, "", ""},
		{"a call of no kind the guard reads", g, sub(Call{Kind: Other, Name: "read_file"}), Allow, "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := Decide(tc.guard, tc.call)
			if v.Outcome != tc.want || !strings.Contains(v.Reason, tc.why) || v.Needs != tc.needs {
				t.Errorf("%+v: got %+v, want outcome %d with %q, needs %q", tc.call, v, tc.want, tc.why, tc.needs)
			}
			if v.Outcome == Allow && v.Reason != "" {
				t.Errorf("allowed with a reason: %q", v.Reason)
			}
		})
	}
}

// S-0353: a served agent writes only as itself whatever the guard's own role
// and story say: the Call's decide whose names it writes as (ADR-0143).
func TestDecideTakesTheRoleAndStoryFromTheCall(t *testing.T) {
	agent := Guard{Commands: g.Commands, Agent: "system-flow"}
	run := Call{Story: "S-0001", Kind: Shell, Command: "flai move S-0001 review --by agent-S-0001"}
	if v := Decide(agent, run); v.Outcome != Allow {
		t.Errorf("the story's agent as its run's name: %+v", v)
	}
	run.Command = "flai move S-0001 review --by alex"
	if v := Decide(agent, run); v.Outcome != Deny || !strings.HasPrefix(v.Reason, "the story's agent cannot run ") || !strings.Contains(v.Reason, "ADR-0143") {
		t.Errorf("the story's agent as another: %+v", v)
	}
	// the operator's own session, with neither role nor story, is not held
	run.Story = ""
	if v := Decide(agent, run); v.Outcome != Allow {
		t.Errorf("the operator's session held: %+v", v)
	}
}

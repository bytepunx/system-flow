package guard

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// subprojectCommands are flai's commands as the cmd package gives them,
// flai subproject among them (S-0491).
var subprojectCommands = append(slices.Clone(g.Commands), "subproject")

// subprojectEdits are the command lines that add, change, or remove a
// sub-project, however they run flai.
var subprojectEdits = []string{
	"flai subproject add web --kind sveltekit",
	"scripts/flai.sh subproject set web --tags ui,app --json",
	"bin/flai subproject remove web",
	"cd /w && flai --config c.json subproject add api --path services/api --kind go",
}

// subprojectReads are the command lines that list the sub-projects or print
// the command's help.
var subprojectReads = []string{"flai subproject list", "scripts/flai.sh subproject list --json", "flai subproject"}

// S-0491: the sub-projects are the project's to shape, so a story's agent,
// served or not, and the operator's own session call subproject_edit and run
// flai subproject add, set, and remove, as well as listing them.
func TestTheStorysAgentEditsTheSubprojects(t *testing.T) {
	for name, gd := range map[string]Guard{
		"story's agent":           {Commands: subprojectCommands, Story: "S-0001", Served: true, Agent: "system-flow", Session: "20261012T000000"},
		"story's agent, unserved": {Commands: subprojectCommands, Story: "S-0001"},
		"operator's own session":  {Commands: subprojectCommands},
	} {
		for _, tool := range []string{SubprojectEdit, Subprojects} {
			if why := gd.Check(Event{ToolName: MCPPrefix + tool}); why != "" {
				t.Errorf("%s: %s refused: %s", name, tool, why)
			}
		}
		for _, c := range append(slices.Clone(subprojectEdits), subprojectReads...) {
			if why := gd.Check(bash("", c)); why != "" {
				t.Errorf("%s: %q refused: %s", name, c, why)
			}
		}
	}
}

// S-0491: every sub-agent lists the sub-projects, with subprojects and flai
// subproject list, but none changes them, a task sub-agent no more than an
// explorer or a verifier, as none runs flai components config: they rewrite
// system-flow.yaml, which is the story's agent's to change.
func TestASubAgentListsTheSubprojectsButDoesNotEditThem(t *testing.T) {
	gd := Guard{Commands: subprojectCommands, Story: "S-0001", ReadOnly: func(agentType string) bool { return agentType == "verifier" || agentType == "explorer" }}
	for _, who := range []string{"explorer", "verifier", "general-purpose"} {
		if why := gd.Check(Event{ToolName: MCPPrefix + Subprojects, AgentID: "a1", AgentType: who}); why != "" {
			t.Errorf("%s: subprojects refused: %s", who, why)
		}
		for _, c := range subprojectReads {
			if why := gd.Check(bash(who, c)); why != "" {
				t.Errorf("%s: %q refused: %s", who, c, why)
			}
		}
		why := gd.Check(Event{ToolName: MCPPrefix + SubprojectEdit, AgentID: "a1", AgentType: who})
		if !strings.Contains(why, "a sub-agent ("+who+") cannot call "+SubprojectEdit) || !strings.Contains(why, "final message") {
			t.Errorf("%s: subproject_edit: %q", who, why)
		}
		for _, c := range subprojectEdits {
			why := gd.Check(bash(who, c))
			if !strings.Contains(why, fmt.Sprintf("a sub-agent (%s) cannot run", who)) || !strings.Contains(why, "rewrite system-flow.yaml") || !strings.Contains(why, "story's agent's") {
				t.Errorf("%s: %q: %q", who, c, why)
			}
		}
	}
}

// S-0491: the planner, the analyzer, and the orchestrator list the
// sub-projects whatever their permissions. The planner's and the analyzer's
// rules, which allow no write they do not name, refuse them subproject_edit
// and flai subproject add, set, and remove; the orchestrator is refused them
// only while orchestration.permissions.edit_subprojects is off, and its
// refusal names the call and that permission (ADR-0195), as
// TestTheOrchestratorEditsSubprojectsByItsPermission shows.
func TestTheStrategicAgentsListTheSubprojectsButDoNotEditThem(t *testing.T) {
	for _, s := range []struct {
		name string
		gd   Guard
	}{
		{"planner", Guard{Commands: subprojectCommands, Role: RolePlan}},
		{"analyzer", Guard{Commands: subprojectCommands, Role: RoleAnalyze}},
		{"orchestrator", Guard{Commands: subprojectCommands, Role: RoleOrchestrate, Permissions: except(manifest.PermitEditSubprojects)}},
		{"orchestrator, every permission off", Guard{Commands: subprojectCommands, Role: RoleOrchestrate}},
	} {
		who, _, _ := strings.Cut(s.name, ",")
		if why := s.gd.Check(Event{ToolName: MCPPrefix + Subprojects}); why != "" {
			t.Errorf("the %s: subprojects refused: %s", s.name, why)
		}
		for _, c := range subprojectReads {
			if why := s.gd.Check(bash("", c)); why != "" {
				t.Errorf("the %s: %q refused: %s", s.name, c, why)
			}
		}
		needs := ""
		if s.gd.Role == RoleOrchestrate {
			needs = "orchestration.permissions." + manifest.PermitEditSubprojects
		}
		r := s.gd.Decide(Event{ToolName: MCPPrefix + SubprojectEdit})
		if !strings.Contains(r.Why, fmt.Sprintf("the %s cannot ", who)) || r.Needs != needs {
			t.Errorf("the %s: subproject_edit: %+v", s.name, r)
		}
		if s.gd.Role == RoleOrchestrate && r.Call != SubprojectEdit {
			t.Errorf("the orchestrator's refusal should name the call: %+v", r)
		}
		for _, c := range subprojectEdits {
			if r := s.gd.Decide(bash("", c)); !strings.Contains(r.Why, fmt.Sprintf("the %s cannot run", who)) || r.Needs != needs {
				t.Errorf("the %s: %q: %+v", s.name, c, r)
			}
		}
	}
}

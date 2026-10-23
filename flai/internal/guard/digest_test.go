package guard

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// S-0428: flai digest list and show are reads for every agent the guard holds
// to reads, and flai digest new and edit, which write an entry, are refused
// to every sub-agent and to the strategic agents alike.
func TestEveryAgentReadsTheDigestButNoneHeldToReadsWritesIt(t *testing.T) {
	cmds := append(slices.Clone(g.Commands), "digest")
	sub := Guard{Commands: cmds, ReadOnly: func(agentType string) bool { return agentType == "verifier" || agentType == "explorer" }}
	plan, orch := planGuard, orchestrator(allOn)
	plan.Commands, orch.Commands = cmds, cmds
	reads := []string{"flai digest list", "scripts/flai.sh digest list --status active --json", "bin/flai digest show D-0001 --json"}
	writes := []string{"flai digest new Fact --kind answer --source TH-0001", "scripts/flai.sh digest edit D-0001 --topic git"}
	for _, who := range []string{"general-purpose", "verifier", "explorer"} {
		for _, c := range reads {
			if why := sub.Check(bash(who, c)); why != "" {
				t.Errorf("%s, %q refused: %s", who, c, why)
			}
		}
		for _, c := range writes {
			if why := sub.Check(bash(who, c)); !strings.Contains(why, fmt.Sprintf("a sub-agent (%s) cannot run %q", who, c)) {
				t.Errorf("%s, %q: %q", who, c, why)
			}
		}
	}
	for _, s := range []struct {
		name string
		gd   Guard
	}{
		{"planner", plan},
		{"orchestrator", orch},
		{"analyzer", Guard{Commands: cmds, Role: RoleAnalyze}},
	} {
		for _, c := range reads {
			if why := s.gd.Check(bash("", c)); why != "" {
				t.Errorf("the %s, %q refused: %s", s.name, c, why)
			}
		}
		for _, c := range writes {
			if why := s.gd.Check(bash("", c)); !strings.Contains(why, fmt.Sprintf("the %s cannot run %q", s.name, c)) {
				t.Errorf("the %s, %q: %q", s.name, c, why)
			}
		}
	}
}

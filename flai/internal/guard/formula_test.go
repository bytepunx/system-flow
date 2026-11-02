package guard

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// formulaCommands are flai's commands as the cmd package gives them, flai
// formula among them (S-0421).
var formulaCommands = append(slices.Clone(g.Commands), "formula")

// S-0421: every sub-agent reads the formulas with flai formula list and show,
// but none runs one: flai formula run makes a story and its tasks, which are
// the story's agent's to write.
func TestASubAgentReadsTheFormulasButDoesNotRunOne(t *testing.T) {
	gd := Guard{Commands: formulaCommands, Story: "S-0001", ReadOnly: func(agentType string) bool { return agentType == "verifier" || agentType == "explorer" }}
	for _, who := range []string{"explorer", "verifier", "general-purpose"} {
		for _, c := range []string{"flai formula", "flai formula list --json", "scripts/flai.sh formula show release-check"} {
			if why := gd.Check(bash(who, c)); why != "" {
				t.Errorf("%s: %q refused: %s", who, c, why)
			}
		}
		why := gd.Check(bash(who, "flai formula run release-check --input version=1.4.0"))
		if !strings.Contains(why, fmt.Sprintf("a sub-agent (%s) cannot run", who)) {
			t.Errorf("%s: flai formula run: %q", who, why)
		}
	}
}

// S-0421: the story's agent, served or not, and the operator's own session
// run flai formula run as well as reading the formulas.
func TestTheStorysAgentRunsAFormula(t *testing.T) {
	for name, gd := range map[string]Guard{
		"story's agent":           {Commands: formulaCommands, Story: "S-0001", Served: true, Agent: "system-flow", Session: "20261012T000000"},
		"story's agent, unserved": {Commands: formulaCommands, Story: "S-0001"},
		"operator's own session":  {Commands: formulaCommands},
	} {
		for _, c := range []string{"flai formula list", "flai formula run release-check --input version=1.4.0"} {
			if why := gd.Check(bash("", c)); why != "" {
				t.Errorf("%s: %q refused: %s", name, c, why)
			}
		}
	}
}

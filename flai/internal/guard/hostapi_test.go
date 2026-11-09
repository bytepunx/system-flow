package guard

import (
	"slices"
	"strings"
	"testing"
)

// S-0612, I-0307: every sub-agent, read-only or not, and every strategic
// agent's own call run flai hostapi alone and with a method on HostAPIReads,
// --timing included, as docs/users/flai.md § Request timing names it; every
// other method is refused, and a sub-agent's refusal names it and says the
// host API's write methods are the story's agent's.
func TestEveryAgentRunsTheHostAPIsReadsButNoneItsWrites(t *testing.T) {
	cmds := append(slices.Clone(g.Commands), "hostapi")
	sub := Guard{Commands: cmds, Story: "S-0001", ReadOnly: func(agentType string) bool { return agentType == "verifier" || agentType == "explorer" }}
	plan, orch, curator := planGuard, orchestrator(allOn), itsCurator
	plan.Commands, orch.Commands, curator.Commands = cmds, cmds, append(slices.Clone(itsCurator.Commands), "hostapi")
	reads := []string{
		"flai hostapi",
		"flai hostapi --json",
		`flai hostapi --timing board.get '{"all":true}' >/dev/null`,
		"flai hostapi --timing inbox.designer 2>&1 >/dev/null",
		`scripts/flai.sh hostapi item.get '{"id":"S-0042"}'`,
		"flai --config x.json hostapi stats.get",
		"flai hostapi > /dev/null",
	}
	// each command line, with the method it names
	writes := map[string]string{}
	writes[`flai hostapi item.move '{"id":"S-0001","to":"review"}'`] = "item.move"
	writes["flai hostapi --timing settings.action"] = "settings.action"
	writes["flai hostapi >/dev/null doc.save"] = "doc.save"
	writes["flai hostapi components.status"] = "components.status"
	for _, who := range []string{"explorer", "verifier", "general-purpose"} {
		for _, c := range reads {
			if why := sub.Check(bash(who, c)); why != "" {
				t.Errorf("%s: %q refused: %s", who, c, why)
			}
		}
		for c, method := range writes {
			why := sub.Check(bash(who, c))
			if !strings.HasPrefix(why, "a sub-agent ("+who+") cannot run \"flai hostapi ") || !strings.Contains(why, method+" is not one") || !strings.Contains(why, "write methods are the story's agent's (ADR-0060)") {
				t.Errorf("%s: %q: %q", who, c, why)
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
		{"curator", curator},
	} {
		for _, c := range reads {
			if why := s.gd.Check(bash("", c)); why != "" {
				t.Errorf("the %s: %q refused: %s", s.name, c, why)
			}
		}
		for c := range writes {
			if why := s.gd.Check(bash("", c)); !strings.HasPrefix(why, "the "+s.name+" cannot run \"flai hostapi ") {
				t.Errorf("the %s: %q: %q", s.name, c, why)
			}
		}
	}
}

// HostAPIReads is sorted and holds no name twice, so that a method is found
// on it by reading it.
func TestHostAPIReadsIsSorted(t *testing.T) {
	if !slices.IsSorted(HostAPIReads) || len(slices.Compact(slices.Clone(HostAPIReads))) != len(HostAPIReads) {
		t.Errorf("HostAPIReads is not sorted, or names a method twice: %v", HostAPIReads)
	}
}

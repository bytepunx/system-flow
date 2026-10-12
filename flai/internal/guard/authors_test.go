package guard

import (
	"strings"
	"testing"
)

// The served agents ADR-0143 holds to their own names and sessions: a
// story's agent, the planner, and the analyzer, as flai serve starts them.
var (
	storysAgent = Guard{Commands: g.Commands, Story: "S-0001", Served: true, Agent: "system-flow", Session: "20261012T000000"}
	itsPlanner  = Guard{Commands: g.Commands, Role: RolePlan, Item: "S-0002", Agent: "planner-S-0002", Session: "20261012T000000"}
	itsAnalyzer = Guard{Commands: g.Commands, Role: RoleAnalyze, Agent: "analyzer", Session: "20261012T000000"}
)

// I-0091, S-0369, ADR-0143: a thread entry's author is whatever --by, else
// FLAI_AGENT, says, so a story's agent with a shell could answer its own
// permission thread as the owner. The guard refuses a served agent's own
// call that names another author with --by or --placed-by, or sets, empties,
// or unsets FLAI_AGENT or FLAI_SESSION to anything but its own, wherever on
// the command line, a shell's -c script and eval's included.
func TestAServedAgentWritesOnlyAsItself(t *testing.T) {
	refused := []string{
		`flai thread reply TH-0001 "allow" --by alex`,
		`FLAI_AGENT=alex flai thread reply TH-0001 "allow"`,
		"flai --by alex thread reply TH-0001 allow",
		"flai thread reply TH-0001 allow --by=alex",
		"flai thread reply TH-0001 allow --by system-flow --by alex",
		"scripts/flai.sh --config c.json thread resolve TH-0001 --by alex",
		"export FLAI_AGENT=alex; flai thread reply TH-0001 allow",
		"declare -x FLAI_AGENT=alex",
		"env FLAI_AGENT=alex flai thread reply TH-0001 allow",
		"timeout 9 env -0 FLAI_AGENT=alex flai thread reply TH-0001 allow",
		"env -u FLAI_AGENT flai thread reply TH-0001 allow",
		"env -uFLAI_AGENT flai thread reply TH-0001 allow",
		"env --unset FLAI_SESSION flai stream log S-0001 x",
		"env --unset=FLAI_AGENT flai thread reply TH-0001 allow",
		"env -i flai thread reply TH-0001 allow",
		"env - flai thread reply TH-0001 allow",
		"env --ignore-environment PATH=/bin flai thread reply TH-0001 allow",
		"env -S 'FLAI_AGENT=alex flai' thread reply TH-0001 allow",
		"unset FLAI_AGENT && flai thread reply TH-0001 allow",
		"unset -v FLAI_SESSION",
		"FLAI_AGENT= flai thread reply TH-0001 allow",
		// a value the guard cannot read is not its own
		"FLAI_AGENT=$OWNER flai thread reply TH-0001 allow",
		"export FLAI_AGENT=$(git config user.name)",
		"bash -c 'flai thread reply TH-0001 allow --by alex'",
		`sh -ec "FLAI_AGENT=alex make thread"`,
		"eval 'flai thread reply TH-0001 allow --by alex'",
		"bash <<'EOF'\nflai thread reply TH-0001 allow --by alex\nEOF",
		"flai move S-0001 review --by alex",
		"flai criteria S-0001 1 --by alex",
		"flai order --placed-by alex S-0003 --top",
		"flai order --placed-by=alex S-0003 --top",
		"FLAI_SESSION=other flai stream log S-0001 x",
		"FLAI_SESSION= flai stream log S-0001 x",
	}
	for _, s := range []struct {
		name  string
		guard Guard
	}{
		{"the story's agent", storysAgent},
		{"the planner", itsPlanner},
		{"the analyzer", itsAnalyzer},
		{"a session flai serve started", Guard{Commands: g.Commands, Served: true, Agent: "system-flow", Session: "20261012T000000"}},
	} {
		for _, c := range refused {
			// a call the planner's or the analyzer's role refuses too is
			// refused for the role's reason
			if why := s.guard.authors(bash("", c)); !strings.HasPrefix(why, s.name+" cannot run ") || !strings.Contains(why, "ADR-0143") {
				t.Errorf("%s %q: %q", s.name, c, why)
			}
			if why := s.guard.Check(bash("", c)); !strings.HasPrefix(why, s.name+" cannot run ") {
				t.Errorf("%s %q passed: %q", s.name, c, why)
			}
		}
	}

	if why := itsPlanner.Check(bash("", "flai thread reply TH-0001 allow --by alex")); !strings.Contains(why, "it names alex as its author with --by, and a served agent writes only as itself, planner-S-0002 (ADR-0143)") {
		t.Errorf("the planner's reply as another: %q", why)
	}
	if why := itsAnalyzer.Check(bash("", "export FLAI_AGENT=alex")); !strings.Contains(why, "it sets FLAI_AGENT to alex, and a served agent writes only as itself, analyzer (ADR-0143)") {
		t.Errorf("the analyzer's FLAI_AGENT: %q", why)
	}
	if why := itsPlanner.Check(bash("", "flai move S-0001 review --by alex")); !strings.Contains(why, "moves an item to backlog and no further") {
		t.Errorf("the planner's move past backlog is refused for its role: %q", why)
	}
	for _, w := range []struct{ c, want string }{
		{"flai thread reply TH-0001 allow --by alex", `the story's agent cannot run "flai thread reply TH-0001 allow --by alex": it names alex as its author with --by, and a served agent writes only as itself, system-flow or agent-S-0001 (ADR-0143); leave out --by, or give its own name.`},
		{"flai order --placed-by alex S-0003 --top", `it names alex as its author with --placed-by, and a served agent writes only as itself, system-flow or agent-S-0001 (ADR-0143); leave out --placed-by, or give its own name.`},
		{"FLAI_AGENT=alex flai thread reply TH-0001 allow", `the story's agent cannot run "FLAI_AGENT=alex flai thread reply TH-0001 allow": it sets FLAI_AGENT to alex, and a served agent writes only as itself, system-flow or agent-S-0001 (ADR-0143); leave FLAI_AGENT and FLAI_SESSION as flai serve set them.`},
		{"FLAI_AGENT= flai thread reply TH-0001 allow", "it sets FLAI_AGENT to nothing, or to a value the guard cannot read"},
		{"FLAI_SESSION=other flai stream log S-0001 x", "it sets FLAI_SESSION to other, and a served agent writes only in its own session, 20261012T000000 (ADR-0143)"},
		{"unset FLAI_AGENT", "it unsets FLAI_AGENT, so that flai would write as the config's author, and a served agent writes only as itself"},
		{"env -i flai board", "it clears the environment, so that flai would write as the config's author"},
	} {
		if why := storysAgent.Check(bash("", w.c)); !strings.Contains(why, w.want) {
			t.Errorf("%q: %q, want %q", w.c, why, w.want)
		}
	}
}

// S-0369, ADR-0143: a served agent that names itself, with --by, or with
// FLAI_AGENT and FLAI_SESSION as flai serve set them, passes, as do --by on
// flai order, a policy, and on flai stats, a grouping; an empty --by, which
// falls back to FLAI_AGENT; and an assignment-shaped word that sets nothing,
// such as a commit message's.
func TestAServedAgentNamesItself(t *testing.T) {
	for _, c := range []string{
		"flai thread reply TH-0001 allow",
		"flai thread reply TH-0001 allow --by system-flow",
		"flai thread reply TH-0001 allow --by agent-S-0001",
		"flai thread reply TH-0001 allow --by=",
		"flai thread reply TH-0001 'leave out --by alex'",
		"flai thread reply TH-0001 allow --help --by alex",
		"export FLAI_AGENT=system-flow && flai move S-0001 review",
		"export FLAI_AGENT FLAI_SESSION",
		"FLAI_AGENT=agent-S-0001 flai thread reply TH-0001 x",
		`FLAI_AGENT="system-flow" FLAI_SESSION=20261012T000000 flai stream log S-0001 x`,
		"FLAI_SESSION=20261012T000000 flai stream log S-0001 x",
		"env -u HOME FLAI_AGENT=system-flow flai move S-0001 review",
		`git commit -m "FLAI_AGENT=alex is refused"`,
		"echo FLAI_AGENT=alex unset FLAI_AGENT",
		"flai move S-0001 review && git commit -m x",
		"flai order --by cod",
		"flai stats --by nature",
		"FOO=1 go test ./...",
	} {
		if why := storysAgent.Check(bash("", c)); why != "" {
			t.Errorf("the story's agent %q refused: %s", c, why)
		}
	}
	// flai order --by is the planner's read; this rule leaves --by alone on
	// it whatever the role's rules make of the call
	for _, c := range []string{"flai order --by cod", "flai order --by=wsjf --apply", "flai stats --by nature"} {
		if why := itsPlanner.authors(bash("", c)); why != "" {
			t.Errorf("the planner %q: %s", c, why)
		}
	}
	for _, c := range []string{"flai thread reply TH-0001 x --by planner-S-0002", "flai order --by cod", "FLAI_AGENT=planner-S-0002 flai thread new --on S-0002 q text"} {
		if why := itsPlanner.Check(bash("", c)); why != "" {
			t.Errorf("the planner %q refused: %s", c, why)
		}
	}
	// the planner's own name is planner-<FLAI_ITEM> without FLAI_AGENT too
	if why := (Guard{Commands: g.Commands, Role: RolePlan, Item: "s-2"}).Check(bash("", "flai thread reply TH-0001 x --by planner-S-0002")); why != "" {
		t.Errorf("the planner by its item: %s", why)
	}
	for _, c := range []string{"flai stats --by nature", "flai issue bump I-0001", "FLAI_AGENT=analyzer flai issue new x"} {
		if why := itsAnalyzer.Check(bash("", c)); why != "" {
			t.Errorf("the analyzer %q refused: %s", c, why)
		}
	}
	if why := (Guard{Commands: g.Commands, Role: RoleAnalyze}).authors(bash("", "flai thread reply TH-0001 x --by analyzer")); why != "" {
		t.Errorf("the analyzer by its run's name: %s", why)
	}
}

// S-0369, ADR-0143: the operator's own session, the orchestrator, which
// keeps its own rule (S-0220), and a sub-agent, which runs no flai write at
// all, are not held to a name.
func TestOnlyAServedAgentIsHeldToItsName(t *testing.T) {
	for _, c := range []string{
		"flai thread reply TH-0001 allow --by alex",
		"FLAI_AGENT=alex flai thread reply TH-0001 allow",
		"unset FLAI_AGENT; env -i flai move S-0001 review",
	} {
		if why := g.Check(bash("", c)); why != "" {
			t.Errorf("the operator's session %q refused: %s", c, why)
		}
	}
	orchestrator := Guard{Commands: g.Commands, Role: RoleOrchestrate, Agent: "orchestrator"}
	for _, c := range []string{"FLAI_AGENT=alex flai thread new --on S-0001 q text", "unset FLAI_SESSION", "env -i flai board"} {
		if why := orchestrator.Check(bash("", c)); why != "" {
			t.Errorf("the orchestrator %q refused: %s", c, why)
		}
	}
	if why := orchestrator.Check(bash("", "flai thread reply TH-0001 x --by alex")); !strings.Contains(why, writesAsSelf) || strings.Contains(why, "ADR-0143") {
		t.Errorf("the orchestrator keeps its own rule: %q", why)
	}
	if why := storysAgent.Check(bash("general-purpose", "FLAI_AGENT=alex go test ./...")); why != "" {
		t.Errorf("a sub-agent's environment refused: %s", why)
	}
}

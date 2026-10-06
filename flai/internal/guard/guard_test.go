package guard

import (
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// g knows flai's commands as the cmd package gives them.
var g = Guard{Commands: []string{"accept", "adr", "archive", "block", "board", "check", "cod", "criteria", "doc", "edit", "epic", "forecast", "guard", "help", "issue", "move", "order", "prime", "promote", "push", "release", "show", "stats", "story", "stream", "task", "thread", "touches", "unblock", "version"}}

func bash(agent, cmd string) Event {
	e := Event{ToolName: "Bash", AgentType: agent}
	if agent != "" {
		e.AgentID = "a1b2"
	}
	e.ToolInput.Command = cmd
	return e
}

func TestTheStorysAgentIsNeverRefused(t *testing.T) {
	main := Event{ToolName: "mcp__flai__item_move", AgentType: "planner"} // a session started with --agent has a type, not an ID
	for _, e := range []Event{main, bash("", "flai move S-0001 review && git commit -m x")} {
		if why := g.Check(e); why != "" {
			t.Errorf("%+v refused: %s", e, why)
		}
	}
}

func TestASubAgentReadsFlaiOverMCP(t *testing.T) {
	for _, tool := range MCPReads {
		if why := g.Check(Event{ToolName: MCPPrefix + tool, AgentID: "a1", AgentType: "explorer"}); why != "" {
			t.Errorf("%s refused: %s", tool, why)
		}
	}
	for _, tool := range []string{"item_move", "item_edit", "item_new", "inbox", "wait_for_work", "wait_for_events", "thread_open", "thread_reply", "thread_resolve", "activity_log"} {
		why := g.Check(Event{ToolName: MCPPrefix + tool, AgentID: "a1", AgentType: "verifier"})
		if !strings.Contains(why, "a sub-agent (verifier) cannot call "+tool) || !strings.Contains(why, "final message") {
			t.Errorf("%s: %q", tool, why)
		}
	}
	if why := g.Check(Event{ToolName: MCPPrefix + "item_move", AgentID: "a1"}); !strings.Contains(why, "a sub-agent (unnamed)") {
		t.Errorf("no agent type: %q", why)
	}
	if why := g.Check(Event{ToolName: "Read", AgentID: "a1", AgentType: "explorer"}); why != "" {
		t.Errorf("Read refused: %s", why)
	}
}

func TestASubAgentRunsChecksButNotWrites(t *testing.T) {
	allowed := []string{
		"make test",
		"make flai-test",
		"cd flai && go test -race -short ./... 2>&1 | tail -60",
		"ls flai internal",
		"scripts/flai.sh check --strict",
		"flai board --json",
		"flai --config /tmp/c.json show S-0001",
		"flai doc search 'sub-agents'",
		"flai thread show TH-0042",
		"flai stream diff S-0175",
		"flai stream --help",
		"git -C /w diff main...HEAD --stat",
		"git log --oneline -5 && git status --short",
		"git",
		"grep -rn 'git commit' .",
		"echo 'flai move S-1 review'",
		"FOO=1 go vet ./...",
		"bash -c 'go test ./...'",
	}
	for _, c := range allowed {
		if why := g.Check(bash("verifier", c)); why != "" {
			t.Errorf("%q refused: %s", c, why)
		}
	}
	refused := []string{
		"flai move S-0175 review",
		"cd /w && scripts/flai.sh block S-1 --reason x",
		"flai board limit in-progress 3",
		"flai thread new --on S-1 'q'",
		"flai --config c.json accept S-1",
		"FLAI_AGENT=x env flai stream log S-1 hi",
		"echo $(flai issue bump I-1)",
		"git commit -am wip",
		"git -C /w stash push -u",
		"true; git checkout main",
		"bash -c 'git reset --hard'",
		"sh -c \"flai move S-1 done\"",
		// found by the verifier before review
		"HOME=/tmp flai move S-1 done",
		"GIT_DIR=/w/.git git commit -m x",
		"env -i git commit -m x",
		"sudo -u me git commit -m x",
		"timeout 30 git commit",
		"git ls-files | xargs git rm",
		`find . -exec git checkout {} ;`,
		"bash -lc 'git commit -m x'",
		"sh -ec 'flai move S-1 done'",
		"eval git push",
	}
	for _, c := range refused {
		why := g.Check(bash("verifier", c))
		if !strings.Contains(why, "a sub-agent (verifier) cannot run") || !strings.Contains(why, "ADR-0060") {
			t.Errorf("%q: %q", c, why)
		}
	}
}

// S-0282: a sub-agent reads a story's criteria and does not tick them; the
// story's agent ticks one once it has verified it, and neither the planner
// nor the orchestrator ticks any.
func TestOnlyTheStorysAgentTicksCriteria(t *testing.T) {
	if why := g.Check(bash("verifier", "flai criteria list S-0001 --json")); why != "" {
		t.Errorf("a sub-agent lists criteria: %s", why)
	}
	for _, c := range []string{"flai criteria tick S-0001 1", "flai criteria untick S-0001 1,2"} {
		if why := g.Check(bash("verifier", c)); !strings.Contains(why, "a sub-agent (verifier) cannot run") {
			t.Errorf("a sub-agent %q: %q", c, why)
		}
		if why := g.Check(bash("", c)); why != "" {
			t.Errorf("the story's agent %q: %s", c, why)
		}
		if why := planGuard.Check(bash("", c)); !strings.Contains(why, "the planner cannot") {
			t.Errorf("the planner %q: %q", c, why)
		}
		if r := orchestrator(allOn).Decide(bash("", c)); !strings.Contains(r.Why, "the orchestrator never does it, whatever its permissions") || r.Needs != "" {
			t.Errorf("the orchestrator %q: %+v", c, r)
		}
	}
	if why := planGuard.Check(bash("", "flai criteria list S-0001")); why != "" {
		t.Errorf("the planner lists criteria: %s", why)
	}
	if r := orchestrator(allOn).Decide(bash("", "flai criteria list S-0001")); r.Why != "" {
		t.Errorf("the orchestrator lists criteria: %+v", r)
	}
}

// S-0189: a verifier runs a cheaper model than the story's agent, so the
// corrections it finds are the story's agent's to make. Whatever sub-agent
// it is, the guard refuses it the commands that would make or record one.
func TestASubAgentCannotMakeTheCorrectionsItFinds(t *testing.T) {
	for _, who := range []string{"verifier", "explorer", "general-purpose"} {
		for _, c := range []string{"git add -A", "git apply fix.patch", "git restore flai/x.go", "git checkout -- flai/x.go", "git commit -m fix", "flai move T-1 done"} {
			if why := g.Check(bash(who, c)); !strings.Contains(why, "a sub-agent ("+who+") cannot run") {
				t.Errorf("%s %q: %q", who, c, why)
			}
		}
		e := Event{ToolName: MCPPrefix + "item_edit", AgentID: "a1", AgentType: who}
		if why := g.Check(e); !strings.Contains(why, "cannot call item_edit") {
			t.Errorf("%s item_edit: %q", who, why)
		}
	}
}

// planGuard is a guard in a planner session (S-0208).
var planGuard = Guard{Commands: g.Commands, Role: RolePlan}

func mcp(tool, to string) Event {
	e := Event{ToolName: MCPPrefix + tool}
	e.ToolInput.ID, e.ToolInput.To = "S-0001", to
	return e
}

// S-0208: the planner writes work items and threads through flai and moves
// an item to backlog, and nothing further.
func TestThePlannerPlansThroughFlai(t *testing.T) {
	allowed := []Event{mcp("item_move", "backlog"), {ToolName: "Read"}, {ToolName: "Grep"}, {ToolName: "Agent"}}
	for _, tool := range append(append([]string{}, MCPReads...), MCPPlans...) {
		allowed = append(allowed, mcp(tool, ""))
	}
	for _, c := range []string{
		"flai story new --epic E-0001 --draft 'Title'",
		"scripts/flai.sh epic new 'Outcome'",
		"flai edit S-0001 --touches flai/cmd --forecast-duration 4h",
		"flai touches S-0001 flai/internal/guard",
		"flai thread new --on E-0001 'Plan' 'text' && flai thread reply TH-0001 'more'",
		"flai move S-0001 backlog",
		"flai --config c.json move --reason 'back' S-0001 backlog --yes",
		"flai issue new 'friction' && flai issue bump I-0001",
		"flai board --json; flai show S-0001; flai stats",
		"flai move --help",
		"git log --oneline -5 && git diff main",
		"ls design/system",
	} {
		allowed = append(allowed, bash("", c))
	}
	for _, e := range allowed {
		if why := planGuard.Check(e); why != "" {
			t.Errorf("%s %q refused: %s", e.ToolName, e.ToolInput.Command, why)
		}
	}
}

// S-0217: the orchestrator's operations are reads in the form that reads: a
// sub-agent and the planner call their tools and run order --by, promote
// --candidates, and release --evaluate, and neither applies an order,
// places a story, or releases.
func TestTheOrchestrationReadsPassAndTheirWritesDoNot(t *testing.T) {
	reads := []string{
		"flai order --by wsjf",
		"flai order --by=cod --json",
		"flai --json order --by throughput",
		"scripts/flai.sh promote --candidates --limit 3",
		"flai promote --candidates=true --json",
		"flai release --evaluate",
		"flai --config c.json release --evaluate --json",
	}
	writes := []string{
		"flai order --by cod --apply",
		"flai order --apply --by fifo",
		"flai order --by=wsjf --apply=true",
		// --apply refuses whatever its value
		"flai order --by fifo --apply=false",
		"flai order S-0001 --top",
		"flai order S-0001 --by=",
		"flai promote",
		"flai promote --candidates=false",
		"flai release",
		"flai release --pending",
		"flai release S-0001 --apply",
		"flai release --evaluate=false --pending",
		"flai release --evaluate; flai release --pending",
	}
	for _, tool := range []string{"order_by_policy", "promote_candidates", "release_evaluate"} {
		if why := g.Check(Event{ToolName: MCPPrefix + tool, AgentID: "a1", AgentType: "explorer"}); why != "" {
			t.Errorf("sub-agent %s refused: %s", tool, why)
		}
		if why := planGuard.Check(mcp(tool, "")); why != "" {
			t.Errorf("planner %s refused: %s", tool, why)
		}
	}
	for _, c := range reads {
		if why := g.Check(bash("verifier", c)); why != "" {
			t.Errorf("sub-agent %q refused: %s", c, why)
		}
		if why := planGuard.Check(bash("", c)); why != "" {
			t.Errorf("planner %q refused: %s", c, why)
		}
	}
	for _, c := range writes {
		if why := g.Check(bash("verifier", c)); !strings.Contains(why, "a sub-agent (verifier) cannot run") {
			t.Errorf("sub-agent %q: %q", c, why)
		}
		if why := planGuard.Check(bash("", c)); !strings.Contains(why, "the planner cannot run") {
			t.Errorf("planner %q: %q", c, why)
		}
	}
}

func TestThePlannerNeverMovesAnItemPastBacklog(t *testing.T) {
	for _, to := range []string{"ready", "in-progress", "review", "done", "cancelled", ""} {
		why := planGuard.Check(mcp("item_move", to))
		if !strings.Contains(why, "the planner cannot move S-0001 to "+to+":") || !strings.Contains(why, "thread_open") {
			t.Errorf("item_move to %q: %q", to, why)
		}
	}
	for _, c := range []string{"flai move S-0001 ready", "flai move S-0001 ready --reason backlog", "flai move S-0001", "flai move --by backlog S-0001 in-progress"} {
		if why := planGuard.Check(bash("", c)); !strings.Contains(why, "moves an item to backlog and no further") {
			t.Errorf("%q: %q", c, why)
		}
	}
}

func TestThePlannerNeitherAcceptsNorPublishesNorEditsCode(t *testing.T) {
	for _, tool := range []string{"thread_resolve", "wait_for_work", "agent_start", "agent_restart", "issue_story"} {
		if why := planGuard.Check(mcp(tool, "")); !strings.Contains(why, "the planner cannot call "+tool+": it plans through flai") {
			t.Errorf("%s: %q", tool, why)
		}
	}
	for _, tool := range fileEdits {
		if why := planGuard.Check(Event{ToolName: tool}); !strings.Contains(why, "the planner cannot use "+tool+":") || !strings.Contains(why, "strategic-agents.md") {
			t.Errorf("%s: %q", tool, why)
		}
	}
	for _, c := range []string{
		"flai accept S-0001",
		"flai release --pending",
		"flai push --pending",
		"flai archive S-0001",
		"flai stream sync S-0001",
		"flai stream open S-0001",
		"flai block S-0001 --reason x",
		"flai thread resolve TH-0001",
		"flai edit S-0001 --no-draft",
		"git commit -m plan",
		"git push",
		"bash -c 'git checkout -b plan'",
	} {
		why := planGuard.Check(bash("", c))
		if !strings.HasPrefix(why, "the planner cannot run ") || !strings.Contains(why, "Ask the operator with thread_open on the item") {
			t.Errorf("%q: %q", c, why)
		}
	}
}

// itemNew is the planner's item_new of an item of a type, a draft or not.
func itemNew(typ string, draft bool) Event {
	e := Event{ToolName: MCPPrefix + "item_new"}
	e.ToolInput.Type, e.ToolInput.Draft = typ, draft
	return e
}

// S-0209: the stories the planner writes are drafts for the operator to
// finalize; a task or an epic it creates as before.
func TestThePlannerCreatesAStoryOnlyAsADraft(t *testing.T) {
	for _, e := range []Event{itemNew("story", true), itemNew("task", false), itemNew("epic", false)} {
		if why := planGuard.Check(e); why != "" {
			t.Errorf("item_new %s draft %v refused: %s", e.ToolInput.Type, e.ToolInput.Draft, why)
		}
	}
	if why := planGuard.Check(itemNew("story", false)); !strings.HasPrefix(why, "the planner cannot create a story that is not a draft: ") || !strings.Contains(why, "give item_new draft true, or flai story new --draft") || !strings.Contains(why, "Ask the operator with thread_open on the item") {
		t.Errorf("item_new story without draft: %q", why)
	}
	for _, c := range []string{
		"flai story new --epic E-0001 'T' --draft",
		"flai story new --draft=true --epic E-0001 'T'",
		"flai story new --draft=false --draft 'T'",
		"scripts/flai.sh --config c.json story new --epic E-0001 --draft 'T'",
		"flai epic new 'Outcome'",
		"flai story new --help",
	} {
		if why := planGuard.Check(bash("", c)); why != "" {
			t.Errorf("%q refused: %s", c, why)
		}
	}
	for _, c := range []string{
		"flai story new --epic E-0001 'T'",
		"flai story new --epic E-0001 'T' --draft=false",
		"flai story new --draft --draft=false 'T'",
		"flai story new --draft=maybe 'T'",
		"scripts/flai.sh story new --epic E-0001 'T'",
		"flai story new --epic E-0001 --draft 'A' && flai story new 'B'",
	} {
		why := planGuard.Check(bash("", c))
		if !strings.HasPrefix(why, "the planner cannot run ") || !strings.Contains(why, ": the stories it writes are drafts for the operator to finalize") || !strings.Contains(why, "Ask the operator with thread_open on the item") {
			t.Errorf("%q: %q", c, why)
		}
	}
	// a sub-agent is held as before, and a session not the planner's not at all
	sub := itemNew("story", true)
	sub.AgentID, sub.AgentType = "a1", "explorer"
	if why := planGuard.Check(sub); !strings.Contains(why, "a sub-agent (explorer) cannot call item_new") {
		t.Errorf("sub-agent item_new: %q", why)
	}
	if why := g.Check(itemNew("story", false)); why != "" {
		t.Errorf("the story's agent's item_new refused: %s", why)
	}
}

// S-0255: the planner drafts a story's tasks with flai task new, and runs
// no other task command.
func TestThePlannerDraftsTasksThroughFlai(t *testing.T) {
	for _, c := range []string{
		`flai task new --story S-0001 "x"`,
		"scripts/flai.sh task new --story S-0001 'x'",
	} {
		if why := planGuard.Check(bash("", c)); why != "" {
			t.Errorf("%q refused: %s", c, why)
		}
	}
	for _, c := range []string{"flai task foo", "flai task", "flai task move T-0001 ready"} {
		why := planGuard.Check(bash("", c))
		if !strings.Contains(why, "it runs only story new with --draft, epic new, task new, edit,") || !strings.Contains(why, "Ask the operator with thread_open on the item") {
			t.Errorf("%q: %q", c, why)
		}
	}
	if why := planGuard.Check(bash("explorer", "flai task new --story S-0001 'x'")); !strings.Contains(why, "a sub-agent (explorer) cannot run") {
		t.Errorf("sub-agent task new: %q", why)
	}
}

// S-0210: touches suggest, forecast, and cod print and write nothing, so a
// sub-agent and the planner run them; flai touches with a story's ID sets its
// touches, which the planner may and a sub-agent may not.
func TestPlanningReadsAreReads(t *testing.T) {
	for _, c := range []string{
		"flai touches suggest S-0001",
		"flai touches suggest S-0001 flai/cmd design/system/flai-cli.md --min 2 --limit 10 --json",
		"scripts/flai.sh --config c.json touches suggest S-0001",
		"flai forecast S-0001 --json",
		"flai cod E-0001 && flai cod S-0001",
	} {
		for _, gr := range []Guard{g, planGuard} {
			if why := gr.Check(bash("explorer", c)); why != "" {
				t.Errorf("role %q sub-agent %q refused: %s", gr.Role, c, why)
			}
		}
		if why := planGuard.Check(bash("", c)); why != "" {
			t.Errorf("planner %q refused: %s", c, why)
		}
	}
	for _, c := range []string{"flai touches S-0001 flai/cmd", "flai touches S-0001 suggest", "flai --config c.json touches S-0001 a/b"} {
		if why := g.Check(bash("verifier", c)); !strings.Contains(why, "a sub-agent (verifier) cannot run") || !strings.Contains(why, "ADR-0060") {
			t.Errorf("sub-agent %q: %q", c, why)
		}
		if why := planGuard.Check(bash("", c)); why != "" {
			t.Errorf("planner %q refused: %s", c, why)
		}
	}
}

// The planner's explorer is a sub-agent like any other.
func TestThePlannersSubAgentIsHeldAsASubAgent(t *testing.T) {
	sub := mcp("item_move", "backlog")
	sub.AgentID, sub.AgentType = "a1", "explorer"
	if why := planGuard.Check(sub); !strings.Contains(why, "a sub-agent (explorer) cannot call item_move") {
		t.Errorf("item_move: %q", why)
	}
	if why := planGuard.Check(bash("explorer", "flai story new 'x'")); !strings.Contains(why, "a sub-agent (explorer) cannot run") {
		t.Errorf("story new: %q", why)
	}
	if why := planGuard.Check(Event{ToolName: "Edit", AgentID: "a1", AgentType: "explorer"}); why != "" {
		t.Errorf("Edit refused: %s", why)
	}
}

// orchestrator is a guard in an orchestrator session with permissions p
// (S-0218).
func orchestrator(p manifest.Permissions) Guard {
	return Guard{Commands: append(append([]string{}, g.Commands...), "plan"), Role: RoleOrchestrate, Permissions: p}
}

// allOn are every permission on; except turns one of them off.
var allOn = manifest.Permissions{PlanBacklogEpics: true, FinalizeDrafts: true, PromoteToReady: true, OrderReady: true, AnswerThreads: manifest.AnswerAutonomous, AcceptReviews: true, Publish: true}

func except(name string) manifest.Permissions {
	p := allOn
	switch name {
	case manifest.PermitPlanBacklogEpics:
		p.PlanBacklogEpics = false
	case manifest.PermitFinalizeDrafts:
		p.FinalizeDrafts = false
	case manifest.PermitPromoteToReady:
		p.PromoteToReady = false
	case manifest.PermitOrderReady:
		p.OrderReady = false
	case manifest.PermitAnswerThreads:
		p.AnswerThreads = manifest.AnswerOff
	case manifest.PermitAcceptReviews:
		p.AcceptReviews = false
	case manifest.PermitPublish:
		p.Publish = false
	}
	return p
}

// itemOf is the orchestrator's call of an MCP tool on an item.
func itemOf(tool, id, to string) Event {
	e := Event{ToolName: MCPPrefix + tool}
	e.ToolInput.ID, e.ToolInput.To = id, to
	return e
}

// S-0218: each of the orchestrator's calls that a permission allows passes
// while the permission is on, and is refused, naming it, while it is off,
// whatever the other permissions are.
func TestTheOrchestratorsPermissionsAllowItsCalls(t *testing.T) {
	for _, c := range []struct {
		permit string
		on     manifest.Permissions
		calls  []Event
	}{
		{manifest.PermitPlanBacklogEpics, manifest.Permissions{PlanBacklogEpics: true}, []Event{itemOf("plan", "E-0016", ""), itemOf("plan", "e-16", ""), bash("", "flai plan E-0016"), bash("", "scripts/flai.sh --config c.json plan E-1 --json")}},
		{manifest.PermitFinalizeDrafts, manifest.Permissions{FinalizeDrafts: true}, []Event{bash("", "flai edit S-0001 --no-draft"), bash("", "flai edit --by orchestrator S-0001 --no-draft=true --hash=abc --json")}},
		{manifest.PermitPromoteToReady, manifest.Permissions{PromoteToReady: true}, []Event{itemOf("item_move", "S-0001", "ready"), bash("", "flai move S-0001 ready"), bash("", "flai move --reason 'top of the backlog' S-0001 ready")}},
		{manifest.PermitOrderReady, manifest.Permissions{OrderReady: true}, []Event{bash("", "flai order --by wsjf --apply")}},
		{manifest.PermitAnswerThreads, manifest.Permissions{AnswerThreads: manifest.AnswerRecommend}, []Event{recommendation(itemOf("thread_reply", "TH-0001", "")), bash("", "flai thread reply --recommend TH-0001 'I recommend S-0002'")}},
		{manifest.PermitAnswerThreads, manifest.Permissions{AnswerThreads: manifest.AnswerAutonomous}, []Event{recommendation(itemOf("thread_reply", "TH-0001", ""))}},
		{manifest.PermitAcceptReviews, manifest.Permissions{AcceptReviews: true}, []Event{bash("", "flai accept S-0001")}},
		{manifest.PermitPublish, manifest.Permissions{Publish: true}, []Event{bash("", "flai release --pending"), bash("", "flai push --pending"), bash("", "flai accept --help; flai push")}},
	} {
		needs := "orchestration.permissions." + c.permit
		for _, e := range c.calls {
			for _, on := range []manifest.Permissions{c.on, allOn} {
				if r := orchestrator(on).Decide(e); r.Why != "" {
					t.Errorf("%s on: %s %q %s refused: %s", c.permit, e.ToolName, e.ToolInput.Command, e.ToolInput.ID, r.Why)
				}
			}
			for _, off := range []manifest.Permissions{{}, {AnswerThreads: manifest.AnswerOff}, except(c.permit)} {
				r := orchestrator(off).Decide(e)
				if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, ": it needs "+needs+", which is off. Ask the operator with thread_open on the item") || r.Needs != needs || r.Call == "" {
					t.Errorf("%s off: %s %q %s: %+v", c.permit, e.ToolName, e.ToolInput.Command, e.ToolInput.ID, r)
				}
			}
		}
	}
}

// S-0218: the orchestrator reads, opens threads, records issues, and logs its
// activities with no permission at all, as a sub-agent reads.
func TestTheOrchestratorAlwaysReadsAndAsks(t *testing.T) {
	none := orchestrator(manifest.Permissions{})
	allowed := []Event{{ToolName: "Read"}, {ToolName: "Grep"}, {ToolName: "Agent"}}
	for _, tool := range append(append([]string{}, MCPReads...), MCPOrchestrates...) {
		allowed = append(allowed, itemOf(tool, "", ""))
	}
	for _, c := range []string{
		"flai order --by wsjf",
		"flai promote --candidates --json",
		"flai release --evaluate",
		"flai thread new --on S-0001 'Ready?' 'S-0001 has no forecast'",
		"flai issue new 'friction' && flai issue bump I-0001",
		"flai board --json; flai show S-0001; flai stats",
		"flai move --help",
		"git log --oneline -5 && git status --short",
		"ls wip/agents",
	} {
		allowed = append(allowed, bash("", c))
	}
	for _, e := range allowed {
		if why := none.Check(e); why != "" {
			t.Errorf("%s %q refused: %s", e.ToolName, e.ToolInput.Command, why)
		}
	}
}

// S-0218: what no permission allows the orchestrator never does, and its
// refusal says so and names no permission.
func TestTheOrchestratorNeverEditsFilesOrWritesBeyondItsPermissions(t *testing.T) {
	all := orchestrator(allOn)
	refused := []Event{itemOf("plan", "S-0001", ""), itemOf("item_move", "S-0001", "in-progress"), itemOf("item_move", "S-0001", "backlog"), itemOf("item_move", "S-0001", "done")}
	for _, tool := range []string{"item_new", "item_edit", "thread_resolve", "wait_for_work", "agent_start", "agent_restart", "issue_story"} {
		refused = append(refused, itemOf(tool, "", ""))
	}
	for _, tool := range fileEdits {
		refused = append(refused, Event{ToolName: tool})
	}
	for _, c := range []string{
		"flai plan S-0001",
		"flai move S-0001 in-progress",
		"flai move S-0001 backlog",
		"flai edit S-0001 --no-draft --title 'Better'",
		"flai edit S-0001 --no-draft --autocommit",
		"flai edit S-0001 --draft",
		"flai story new --epic E-0001 --draft 'x'",
		"flai thread resolve TH-0001",
		"flai release S-0001 --apply",
		"flai stream open S-0001",
		"flai archive S-0001",
		"flai touches S-0001 flai/cmd",
		"git commit -m order",
		"git push",
		"bash -c 'git checkout main'",
	} {
		refused = append(refused, bash("", c))
	}
	for _, e := range refused {
		r := all.Decide(e)
		if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, ": the orchestrator never does it, whatever its permissions: ") || !strings.Contains(r.Why, "thread_open on the item") || r.Needs != "" || r.Call == "" {
			t.Errorf("%s %q %s: %+v", e.ToolName, e.ToolInput.Command, e.ToolInput.ID, r)
		}
	}
	if r := all.Decide(itemOf("plan", "S-0001", "")); !strings.Contains(r.Why, "the orchestrator cannot plan S-0001: ") || !strings.Contains(r.Why, plansEpics) {
		t.Errorf("plan a story: %q", r.Why)
	}
}

// S-0218: a refusal carries the call, in one line, for the orchestrator's
// activity document.
func TestTheOrchestratorsRefusalNamesTheCall(t *testing.T) {
	edit := Event{ToolName: "Edit"}
	edit.ToolInput.FilePath = "flai/cmd/x.go"
	for _, c := range []struct {
		e    Event
		call string
	}{
		{itemOf("item_move", "S-0001", "ready"), "item_move S-0001 ready"},
		{itemOf("thread_reply", "", ""), "thread_reply"},
		{bash("", "cd /w &&\n  flai move S-0001 ready"), "flai move S-0001 ready"},
		{edit, "Edit flai/cmd/x.go"},
	} {
		if r := orchestrator(manifest.Permissions{}).Decide(c.e); r.Call != c.call {
			t.Errorf("%s: call %q, want %q", c.e.ToolName, r.Call, c.call)
		}
	}
	if r := orchestrator(manifest.Permissions{}).Decide(bash("", "flai move S-0001 ready")); r.Why != `the orchestrator cannot run "flai move S-0001 ready": it needs orchestration.permissions.promote_to_ready, which is off. `+askOperator {
		t.Errorf("message: %q", r.Why)
	}
}

// S-0218: the orchestrator's sub-agents are held as every sub-agent is,
// whatever the orchestrator's permissions.
func TestTheOrchestratorsSubAgentIsHeldAsASubAgent(t *testing.T) {
	all := orchestrator(allOn)
	sub := itemOf("item_move", "S-0001", "ready")
	sub.AgentID, sub.AgentType = "a1", "explorer"
	if why := all.Check(sub); !strings.Contains(why, "a sub-agent (explorer) cannot call item_move") {
		t.Errorf("item_move: %q", why)
	}
	if why := all.Check(bash("explorer", "flai accept S-0001")); !strings.Contains(why, "a sub-agent (explorer) cannot run") {
		t.Errorf("accept: %q", why)
	}
	if r := all.Decide(Event{ToolName: "Edit", AgentID: "a1", AgentType: "explorer"}); r.Why != "" || r.Call != "" {
		t.Errorf("Edit refused: %+v", r)
	}
}

// Without the role plan or orchestrate the guard decides as it did before
// the planner.
func TestOtherRolesAreDecidedAsBefore(t *testing.T) {
	for _, role := range []string{"", "analyze", "story"} {
		gr := Guard{Commands: g.Commands, Role: role}
		for _, e := range []Event{{ToolName: "Edit"}, {ToolName: "Write"}, mcp("item_move", "review"), bash("", "flai accept S-1 && git commit -m x"), {ToolName: "Edit", AgentID: "a1", AgentType: "verifier"}} {
			if why := gr.Check(e); why != "" {
				t.Errorf("role %q %s refused: %s", role, e.ToolName, why)
			}
		}
		if why := gr.Check(bash("verifier", "flai move S-1 backlog")); !strings.Contains(why, "a sub-agent (verifier) cannot run") || !strings.Contains(why, "story's agent's") {
			t.Errorf("role %q sub-agent move: %q", role, why)
		}
	}
}

func TestCommands(t *testing.T) {
	got := commands(`a "b c" 'd'; e|f && g $(h i) ` + "`j`")
	want := [][]string{{"a", "b c", "d"}, {"e"}, {"f"}, {"g"}, {"h", "i"}, {"j"}}
	if len(got) != len(want) {
		t.Fatalf("got %q", got)
	}
	for i := range want {
		if strings.Join(got[i], "|") != strings.Join(want[i], "|") {
			t.Errorf("%d: got %q want %q", i, got[i], want[i])
		}
	}
}

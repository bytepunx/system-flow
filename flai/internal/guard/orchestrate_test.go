package guard

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// call is a tool call as Claude Code hands it to the guard, decoded as flai
// guard decodes it.
func call(t *testing.T, tool, input string) Event {
	t.Helper()
	var e Event
	if err := json.Unmarshal([]byte(`{"tool_name":"`+MCPPrefix+tool+`","tool_input":`+input+`}`), &e); err != nil {
		t.Fatalf("%s %s: %v", tool, input, err)
	}
	return e
}

// described is a call as a test failure names it.
func described(e Event) string {
	return strings.TrimSpace(e.ToolName + " " + e.ToolInput.Command + " " + e.ToolInput.ID + " " + strings.Join(e.ToolInput.Fields, ","))
}

func TestAnInputKnowsTheFieldsItGives(t *testing.T) {
	e := call(t, "item_edit", `{"id":"S-0001","draft":false,"hash":"abc"}`)
	if !slices.Equal(e.ToolInput.Fields, []string{"draft", "hash", "id"}) || e.ToolInput.ID != "S-0001" || e.ToolInput.Draft {
		t.Errorf("input: %+v", e.ToolInput)
	}
	if e := call(t, "item_new", `{"type":"story","draft":true}`); !e.ToolInput.Draft || e.ToolInput.Type != "story" {
		t.Errorf("item_new: %+v", e.ToolInput)
	}
	if e := call(t, "item_edit", `{}`); len(e.ToolInput.Fields) != 0 || e.ToolInput.CostOfDelay != nil {
		t.Errorf("no fields: %+v", e.ToolInput)
	}
	// S-0328: the keys of cost_of_delay, as given; one that is not an object
	// leaves them nil and still decodes, so that the guard reads the call
	e = call(t, "item_edit", `{"id":"S-0001","cost_of_delay":{"revenue_per_week":500,"value":null}}`)
	if got := e.ToolInput.CostOfDelay; len(got) != 2 || string(got["revenue_per_week"]) != "500" || string(got["value"]) != "null" {
		t.Errorf("cost_of_delay: %+v", e.ToolInput.CostOfDelay)
	}
	if e := call(t, "item_edit", `{"id":"S-0001","cost_of_delay":"x"}`); e.ToolInput.CostOfDelay != nil || !slices.Equal(e.ToolInput.Fields, []string{"cost_of_delay", "id"}) {
		t.Errorf("cost_of_delay not an object: %+v", e.ToolInput)
	}
}

// S-0219, S-0328: each of the five permissions of the orchestrator's
// planning and ordering passes its calls, by MCP and by command line, while
// it is on, and while it is off refuses them naming it, whatever the others
// are.
func TestTheOrchestratorPlansFinalizesPromotesAndOrdersByItsPermissions(t *testing.T) {
	for _, c := range []struct {
		permit string
		on     manifest.Permissions
		calls  []Event
	}{
		{manifest.PermitPlanBacklogEpics, manifest.Permissions{PlanBacklogEpics: true}, []Event{
			call(t, "plan", `{"id":"E-0016"}`),
			bash("", "flai plan E-0016"),
			bash("", "flai plan --json e-16"),
			bash("", "flai plan --candidates=false E-0016"),
		}},
		{manifest.PermitPlanBacklogStories, manifest.Permissions{PlanBacklogStories: true}, []Event{
			call(t, "plan", `{"id":"S-0328"}`),
			call(t, "item_edit", `{"id":"S-0328","cost_of_delay":{"revenue_per_week":500}}`),
			call(t, "item_edit", `{"project":"system-flow","id":"s-328","hash":"abc","cost_of_delay":{"revenue_per_week":500,"penalty_per_week":0,"time_lost_per_cycle":"4h"}}`),
			bash("", "flai plan S-0328"),
			bash("", "scripts/flai.sh --config c.json plan --json s-328"),
			bash("", "flai edit S-0328 --revenue-per-week 500"),
			bash("", "flai edit --by orchestrator S-0328 --penalty-per-week=50 --time-lost-per-cycle 4h --hash abc --json"),
		}},
		{manifest.PermitFinalizeDrafts, manifest.Permissions{FinalizeDrafts: true}, []Event{
			call(t, "item_edit", `{"id":"S-0001","draft":false}`),
			call(t, "item_edit", `{"project":"system-flow","id":"S-0001","hash":"abc","draft":false}`),
			bash("", "flai edit S-0001 --no-draft"),
			bash("", "flai edit S-0001 --no-draft --hash abc --json"),
		}},
		{manifest.PermitPromoteToReady, manifest.Permissions{PromoteToReady: true}, []Event{
			call(t, "item_move", `{"id":"S-0001","to":"ready"}`),
			bash("", "flai move S-0001 ready"),
			bash("", "flai move --by orchestrator S-0001 ready --reason 'next by wsjf'"),
		}},
		{manifest.PermitOrderReady, manifest.Permissions{OrderReady: true}, []Event{
			bash("", "flai order --by wsjf --apply"),
			bash("", "flai order --apply --by=cod --keep-placed 2h"),
			bash("", "scripts/flai.sh --config c.json order --by fifo --apply --keep-placed=0 --json"),
		}},
	} {
		needs := "orchestration.permissions." + c.permit
		for _, e := range c.calls {
			for _, on := range []manifest.Permissions{c.on, allOn} {
				if r := orchestrator(on).Decide(e); r.Why != "" {
					t.Errorf("%s on: %s refused: %s", c.permit, described(e), r.Why)
				}
			}
			for _, off := range []manifest.Permissions{{}, except(c.permit)} {
				r := orchestrator(off).Decide(e)
				if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, ": it needs "+needs+", which is off. "+askOperator) || r.Needs != needs || r.Call == "" {
					t.Errorf("%s off: %s: %+v", c.permit, described(e), r)
				}
			}
		}
	}
}

// S-0219: a refusal of item_edit says it would finalize the draft, and
// names the permission.
func TestTheOrchestratorsFinalizingRefusalNamesItsPermission(t *testing.T) {
	r := orchestrator(except(manifest.PermitFinalizeDrafts)).Decide(call(t, "item_edit", `{"id":"S-0001","draft":false}`))
	if want := "the orchestrator cannot finalize S-0001: it needs orchestration.permissions.finalize_drafts, which is off. " + askOperator; r.Why != want || r.Call != "item_edit S-0001" || r.Needs != "orchestration.permissions.finalize_drafts" {
		t.Errorf("got %+v, want %q", r, want)
	}
}

// S-0219, S-0328: the orchestrator never plans an item but an epic or a
// story, edits an item beyond finalizing a draft or giving a story cost of
// delay inputs, or places a story by hand, whatever its permissions.
func TestTheOrchestratorNeverPlansATaskEditsAnItemOrPlacesOneByHand(t *testing.T) {
	edits := []string{
		`{"id":"S-0001"}`,
		`{"id":"S-0001","draft":true}`,
		`{"id":"S-0001","draft":null,"title":"Better"}`,
		// a cost of delay beyond its inputs, or beside another field
		`{"id":"S-0001","cost_of_delay":{"value":300}}`,
		`{"id":"S-0001","cost_of_delay":{"revenue_per_week":500,"value":300}}`,
		`{"id":"S-0001","cost_of_delay":{"revenue_per_week":500},"clear_cost_of_delay":true}`,
		`{"id":"S-0001","clear_cost_of_delay":true}`,
		`{"id":"S-0001","cost_of_delay":{"time_lost_per_cycle":""}}`,
		`{"id":"S-0001","cost_of_delay":{"revenue_per_week":null}}`,
		`{"id":"S-0001","cost_of_delay":{}}`,
		`{"id":"S-0001","cost_of_delay":null}`,
		`{"id":"S-0001","cost_of_delay":"x"}`,
		`{"id":"S-0001","cost_of_delay":{"revenue_per_week":500},"draft":false}`,
		`{"id":"S-0001","cost_of_delay":{"revenue_per_week":500},"forecast":{"duration":"2h"}}`,
		`{"id":"E-0001","cost_of_delay":{"revenue_per_week":500}}`,
		`{"cost_of_delay":{"revenue_per_week":500}}`,
	}
	for _, f := range []string{"title", "nature", "tags", "topics", "touches", "after", "parent", "agent", "clear_agent", "body", "cost_of_delay", "clear_cost_of_delay", "forecast", "clear_forecast"} {
		edits = append(edits, `{"id":"S-0001","draft":false,"`+f+`":"x"}`)
	}
	for _, c := range []struct {
		why   string
		calls []Event
	}{
		{plansCandidates, []Event{
			call(t, "plan", `{"id":"T-0001"}`),
			call(t, "plan", `{}`),
			bash("", "flai plan T-0001"),
			bash("", "flai plan --json"),
		}},
		{editsMCP, nil},
		{editsCLI, []Event{
			bash("", "flai edit S-0001 --title Better"),
			bash("", "flai edit S-0001 --cost-of-delay-value 300"),
			bash("", "flai edit S-0001 --revenue-per-week 500 --cost-of-delay-value=300"),
			bash("", "flai edit S-0001 --clear-cost-of-delay"),
			bash("", "flai edit S-0001 --penalty-per-week ''"),
			bash("", "flai edit S-0001 --revenue-per-week 500 --penalty-per-week="),
			bash("", "flai edit S-0001 --revenue-per-week"),
			bash("", "flai edit S-0001 --revenue-per-week 500 --touches flai"),
			bash("", "flai edit S-0001 --revenue-per-week 500 --no-draft"),
			bash("", "flai edit E-0001 --revenue-per-week 500"),
			bash("", "flai edit S-0001"),
		}},
		{byHand, []Event{
			bash("", "flai order S-0001 --top"),
			bash("", "flai order S-0001 --after S-0002 --placed-by orchestrator"),
			bash("", "flai order --before S-0002 S-0001"),
			bash("", "flai order S-0001 --bottom --by=wsjf --apply"),
			bash("", "flai order --by wsjf --apply S-0001"),
			bash("", "flai order --by '' S-0001 --top"),
			bash("", "flai order --keep-placed 2h S-0001"),
			bash("", "flai order --by cod --apply --keep-placed 2h S-0001"),
		}},
	} {
		if c.calls == nil {
			for _, in := range edits {
				c.calls = append(c.calls, call(t, "item_edit", in))
			}
		}
		for _, e := range c.calls {
			r := orchestrator(allOn).Decide(e)
			if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, ": the orchestrator never does it, whatever its permissions: "+c.why+". "+askOperator) || r.Needs != "" || r.Call == "" {
				t.Errorf("%s: %+v", described(e), r)
			}
		}
	}
}

// S-0328: a refusal of a story's plan, or of its cost of delay inputs, says
// what the call would do and names plan_backlog_stories.
func TestTheOrchestratorsStoryPlanningRefusalsNamePlanBacklogStories(t *testing.T) {
	const needs = "orchestration.permissions.plan_backlog_stories"
	for _, c := range []struct {
		e    Event
		why  string
		call string
	}{
		{call(t, "plan", `{"id":"S-0001"}`), "the orchestrator cannot plan S-0001: it needs " + needs + ", which is off. " + askOperator, "plan S-0001"},
		{call(t, "item_edit", `{"id":"S-0001","cost_of_delay":{"revenue_per_week":500}}`), "the orchestrator cannot give S-0001 cost of delay inputs: it needs " + needs + ", which is off. " + askOperator, "item_edit S-0001"},
		{bash("", "flai edit S-0001 --revenue-per-week 500"), `the orchestrator cannot run "flai edit S-0001 --revenue-per-week 500": it needs ` + needs + ", which is off. " + askOperator, "flai edit S-0001 --revenue-per-week 500"},
	} {
		if r := orchestrator(except(manifest.PermitPlanBacklogStories)).Decide(c.e); r.Why != c.why || r.Call != c.call || r.Needs != needs {
			t.Errorf("%s:\n got %+v\nwant %q, call %q", described(c.e), r, c.why, c.call)
		}
	}
}

// S-0219: the orchestrator's new reads pass with every permission off, for
// the orchestrator and its sub-agents alike, and their forms that write do
// not pass as reads.
func TestTheOrchestratorsNewReadsPassWithEveryPermissionOff(t *testing.T) {
	none := orchestrator(manifest.Permissions{})
	for _, c := range []string{
		"flai plan --candidates",
		"flai plan --candidates --json",
		"flai promote --candidates",
		"flai promote --drafts",
		"scripts/flai.sh --config c.json promote --drafts=true --json",
		"flai order --by wsjf",
		"flai order --by cod --keep-placed 2h --json",
		"flai release --evaluate",
	} {
		if why := none.Check(bash("", c)); why != "" {
			t.Errorf("orchestrator %q refused: %s", c, why)
		}
		if why := none.Check(bash("explorer", c)); why != "" {
			t.Errorf("sub-agent %q refused: %s", c, why)
		}
	}
	for _, c := range []string{
		"flai plan E-0016 --candidates=false",
		"flai promote --drafts=false",
		"flai order --by wsjf S-0001",
		"flai order --by wsjf --top",
		"flai order --by '' S-0001",
	} {
		if why := none.Check(bash("explorer", c)); !strings.Contains(why, "a sub-agent (explorer) cannot run") {
			t.Errorf("sub-agent %q: %q", c, why)
		}
		if r := none.Decide(bash("", c)); r.Why == "" {
			t.Errorf("orchestrator %q passed", c)
		}
	}
}

// queueing is the orchestrator's guard with permissions p and flai queue
// among the commands it knows.
func queueing(p manifest.Permissions) Guard {
	o := orchestrator(p)
	o.Commands = append(o.Commands, "queue")
	return o
}

// S-0403 criterion 2: while accept_reviews is on, the orchestrator enqueues a
// story, with queue_add or flai queue add, only as itself, and lists the
// queue whatever its permissions; it never lands a story with flai accept
// --batch, nor takes one out of the queue, holds it, or resumes it, whatever
// its permissions, and each refusal names queue_add.
func TestTheOrchestratorOnlyEnqueuesInTheMergeQueue(t *testing.T) {
	on, off := queueing(manifest.Permissions{AcceptReviews: true}), queueing(except(manifest.PermitAcceptReviews))
	needs := "orchestration.permissions." + manifest.PermitAcceptReviews
	for _, c := range []string{"flai queue", "flai queue --json", "scripts/flai.sh --config c.json queue --json"} {
		if r := queueing(manifest.Permissions{}).Decide(bash("", c)); r.Why != "" {
			t.Errorf("%q refused with every permission off: %+v", c, r)
		}
	}
	for _, e := range []Event{
		call(t, QueueAdd, `{"story":"S-0001","verified":"abc","evidence":"Verdict: pass\n- 1: x"}`),
		bash("", "flai queue add S-0001 --by orchestrator --verified abc --evidence -"),
		bash("", "scripts/flai.sh --config c.json queue add --by=orchestrator s-1 --verified=abc --evidence=evidence.md --json"),
	} {
		if r := on.Decide(e); r.Why != "" {
			t.Errorf("on: %s refused: %+v", described(e), r)
		}
		if r := off.Decide(e); !strings.Contains(r.Why, ": it needs "+needs+", which is off. "+askOperator) || r.Needs != needs || r.Call == "" {
			t.Errorf("off: %s: %+v", described(e), r)
		}
	}
	if r := off.Decide(call(t, QueueAdd, `{"story":"S-0001"}`)); r.Call != "queue_add S-0001" || !strings.HasPrefix(r.Why, "the orchestrator cannot enqueue S-0001: ") {
		t.Errorf("the refusal of queue_add names the story: %+v", r)
	}
	for _, c := range []string{"flai queue add S-0001", "flai queue add S-0001 --by alex --verified abc --evidence -"} {
		if r := on.Decide(bash("", c)); r.Why != fmt.Sprintf("the orchestrator cannot run %q: %s. %s", c, enqueuesAsSelf, askOperator) || r.Needs != needs {
			t.Errorf("on %q: %+v", c, r)
		}
	}
	for _, c := range []string{
		"flai queue remove S-0001",
		"flai queue hold",
		"flai queue resume",
		"scripts/flai.sh --config c.json queue remove s-1 --json",
		"flai accept S-0001 --by orchestrator --batch 20261019-041200",
		"flai accept S-0001 --by=orchestrator --batch=20261019-041200 --dry-run",
		"flai queue add S-0001 --by orchestrator --verified abc --evidence - && flai queue hold",
	} {
		for _, gr := range []Guard{on, queueing(allOn)} {
			r := gr.Decide(bash("", c))
			if !strings.Contains(r.Why, ": the orchestrator never does it, whatever its permissions: "+enqueuesOnly+". "+askOperator) || !strings.Contains(r.Why, "queue_add") || r.Needs != "" || r.Call == "" {
				t.Errorf("%q: %+v", c, r)
			}
		}
	}
	if r := on.Decide(bash("", "flai accept S-0001 --by orchestrator --verified abc --evidence -")); r.Why != "" {
		t.Errorf("the orchestrator's acceptance without --batch is flai accept's to judge by the queue: %+v", r)
	}
}

// S-0491 (ADR-0195): while edit_subprojects is on, the orchestrator adds,
// changes, and removes a sub-project, with flai subproject and its
// --autocommit, which commits the change on main, or with subproject_edit;
// while it is off, whatever the others are, each is refused naming the
// permission. Listing them passes with every permission off.
func TestTheOrchestratorEditsSubprojectsByItsPermission(t *testing.T) {
	subprojecting := func(p manifest.Permissions) Guard {
		o := orchestrator(p)
		o.Commands = append(o.Commands, "subproject")
		return o
	}
	needs := "orchestration.permissions." + manifest.PermitEditSubprojects
	for _, e := range []Event{
		call(t, SubprojectEdit, `{"op":"add","name":"web","kind":"sveltekit"}`),
		call(t, SubprojectEdit, `{"op":"remove","name":"web"}`),
		bash("", "flai subproject add web --kind sveltekit --autocommit --trailer 'Co-Authored-By: someone <s@s>'"),
		bash("", "flai subproject set web --tags ui,app --autocommit"),
		bash("", "scripts/flai.sh --config c.json subproject remove web --autocommit --json"),
	} {
		for _, on := range []manifest.Permissions{{EditSubprojects: true}, allOn} {
			if r := subprojecting(on).Decide(e); r.Why != "" {
				t.Errorf("on: %s refused: %+v", described(e), r)
			}
		}
		for _, off := range []manifest.Permissions{{}, except(manifest.PermitEditSubprojects)} {
			r := subprojecting(off).Decide(e)
			if !strings.HasPrefix(r.Why, "the orchestrator cannot ") || !strings.Contains(r.Why, ": it needs "+needs+", which is off. "+askOperator) || r.Needs != needs || r.Call == "" {
				t.Errorf("off: %s: %+v", described(e), r)
			}
		}
	}
	if r := subprojecting(manifest.Permissions{}).Decide(call(t, SubprojectEdit, `{"op":"add","name":"web"}`)); r.Why != "the orchestrator cannot edit sub-projects: it needs "+needs+", which is off. "+askOperator || r.Call != SubprojectEdit {
		t.Errorf("the refusal of subproject_edit: %+v", r)
	}
	for _, e := range []Event{call(t, Subprojects, `{}`), bash("", "flai subproject list --json"), bash("", "flai subproject")} {
		if r := subprojecting(manifest.Permissions{}).Decide(e); r.Why != "" {
			t.Errorf("%s refused with every permission off: %+v", described(e), r)
		}
	}
}

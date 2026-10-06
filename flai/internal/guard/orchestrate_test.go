package guard

import (
	"encoding/json"
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
	if e := call(t, "item_edit", `{}`); len(e.ToolInput.Fields) != 0 {
		t.Errorf("no fields: %+v", e.ToolInput)
	}
}

// S-0219: each of the four permissions of the orchestrator's planning and
// ordering passes its calls, by MCP and by command line, while it is on, and
// while it is off refuses them naming it, whatever the others are.
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

// S-0219: the orchestrator never plans a story, edits an item beyond
// finalizing it, or places a story by hand, whatever its permissions.
func TestTheOrchestratorNeverPlansAStoryEditsAnItemOrPlacesOneByHand(t *testing.T) {
	edits := []string{
		`{"id":"S-0001"}`,
		`{"id":"S-0001","draft":true}`,
		`{"id":"S-0001","draft":null,"title":"Better"}`,
	}
	for _, f := range []string{"title", "nature", "tags", "topics", "touches", "after", "parent", "agent", "clear_agent", "body", "cost_of_delay", "clear_cost_of_delay", "forecast", "clear_forecast"} {
		edits = append(edits, `{"id":"S-0001","draft":false,"`+f+`":"x"}`)
	}
	for _, c := range []struct {
		why   string
		calls []Event
	}{
		{plansEpics, []Event{
			call(t, "plan", `{"id":"S-0001"}`),
			call(t, "plan", `{}`),
			bash("", "flai plan S-0001"),
			bash("", "flai plan --json T-0001"),
		}},
		{"it edits an item only to finalize a draft, with item_edit draft false and nothing else", nil},
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

package roles

import (
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func TestConditionsTableIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range Conditions {
		if seen[c.Name] {
			t.Errorf("condition %s is listed twice", c.Name)
		}
		seen[c.Name] = true
		if !kebab.MatchString(c.Name) || c.Allows == "" {
			t.Errorf("condition %+v needs a kebab-case name and what it allows", c)
		}
		// A condition qualifies calls, triggers, or both (ADR-0211 §3), or
		// else each files entry and nothing more.
		if c.Files == (c.OnCall() || c.OnTrigger()) {
			t.Errorf("condition %s qualifies calls %v, a trigger %v, files %v; want calls or a trigger, or files alone", c.Name, c.OnCall(), c.OnTrigger(), c.Files)
		}
		if c.Schedule && len(c.Events) > 0 {
			t.Errorf("condition %s qualifies both a schedule and events", c.Name)
		}
		if c.Permission != "" && !slices.Contains(Permissions, c.Permission) {
			t.Errorf("condition %s is gated by unknown permission %s", c.Name, c.Permission)
		}
		for _, p := range c.Also {
			if c.Permission == "" || p == c.Permission || !slices.Contains(Permissions, p) {
				t.Errorf("condition %s may also stand under %s, unknown, its own, or with no permission gating it", c.Name, p)
			}
		}
		for _, e := range c.Events {
			if !slices.Contains(Events, e) {
				t.Errorf("condition %s qualifies unknown event %s", c.Name, e)
			}
		}
	}
	for _, g := range Gated {
		if !slices.Contains(Permissions, g.Permission) {
			t.Errorf("gate %+v names an unknown permission", g)
		}
		if g.Condition == "" {
			continue
		}
		c, ok := Condition(g.Condition)
		if !ok || !c.GatedBy(g.Permission) || !c.Qualifies(g.Call) {
			t.Errorf("gate %+v names a condition that does not qualify its call under its permission", g)
		}
	}
	for name := range matches {
		if c, ok := Condition(name); !ok || !c.OnTrigger() {
			t.Errorf("match %s is for no condition that qualifies a trigger", name)
		}
	}
	// ADR-0250 §5: every trigger condition is judged against an occurrence
	for _, c := range Conditions {
		if _, ok := matches[c.Name]; ok != c.OnTrigger() {
			t.Errorf("condition %s qualifies a trigger %v, has a match %v", c.Name, c.OnTrigger(), ok)
		}
	}
	for _, f := range Never {
		for _, w := range f.Unless {
			if _, ok := Condition(w); !ok {
				t.Errorf("never-granted %s names unknown condition %s", f.Call, w)
			}
		}
	}
}

func TestCondition(t *testing.T) {
	c, ok := Condition("to-ready")
	if !ok || c.Permission != PromoteToReady || !c.Qualifies("flai move") || c.Qualifies("item_new") {
		t.Errorf("Condition(to-ready) = %+v, %v", c, ok)
	}
	// S-0523, ADR-0224: cost-inputs stands under either permission that
	// lets the orchestrator set cost of delay inputs, and no other
	cost := cond(t, "cost-inputs")
	if !cost.GatedBy(PlanBacklogStories) || !cost.GatedBy(SetCostOfDelay) || cost.GatedBy(FinalizeDrafts) || cost.GatedBy("") || c.GatedBy(SetCostOfDelay) {
		t.Errorf("cost-inputs = %+v, to-ready = %+v", cost, c)
	}
	back, ok := Condition("to-backlog")
	if !ok || !back.OnCall() || !back.OnTrigger() || !back.Qualifies("item_move") || !back.QualifiesTrigger(Trigger{Kind: TriggerEvent, Event: "moved"}) {
		t.Errorf("Condition(to-backlog) = %+v, %v; want it on item_move and on {event: moved}", back, ok)
	}
	for name, event := range map[string]string{"to-review": "moved", "draft-changed": "edited", "awaits-orchestrator": "thread"} {
		c, ok := Condition(name)
		if !ok || c.OnCall() || !c.OnTrigger() || !c.QualifiesTrigger(Trigger{Kind: TriggerEvent, Event: event}) || c.Allows == "" {
			t.Errorf("Condition(%s) = %+v, %v; want a trigger condition on {event: %s}", name, c, ok, event)
		}
	}
	if _, ok := Condition("anything-goes"); ok {
		t.Error("Condition(anything-goes) found a condition")
	}
}

func TestQualifiesTrigger(t *testing.T) {
	stale, _ := Condition("stale-figures")
	each, _ := Condition("each-ready-story")
	draft, _ := Condition("draft-only")
	cases := []struct {
		cond Cond
		t    Trigger
		want bool
	}{
		{stale, Trigger{Kind: TriggerEvent, Event: "edited"}, true},
		{stale, Trigger{Kind: TriggerEvent, Event: "accepted"}, false},
		{stale, Trigger{Kind: TriggerSchedule, Schedule: "planning.schedule"}, false},
		{each, Trigger{Kind: TriggerSchedule, Schedule: "planning.schedule"}, true},
		{each, Trigger{Kind: TriggerAsked}, false},
		{draft, Trigger{Kind: TriggerEvent, Event: "edited"}, false},
		{cond(t, "to-review"), Trigger{Kind: TriggerEvent, Event: "moved"}, true},
		{cond(t, "to-review"), Trigger{Kind: TriggerEvent, Event: "review"}, false},
		{cond(t, "to-backlog"), Trigger{Kind: TriggerEvent, Event: "moved"}, true},
		{cond(t, "to-backlog"), Trigger{Kind: TriggerEvent, Event: "edited"}, false},
		{cond(t, "draft-changed"), Trigger{Kind: TriggerEvent, Event: "edited"}, true},
		{cond(t, "draft-changed"), Trigger{Kind: TriggerEvent, Event: "moved"}, false},
		{cond(t, "awaits-orchestrator"), Trigger{Kind: TriggerEvent, Event: "thread"}, true},
		{cond(t, "awaits-orchestrator"), Trigger{Kind: TriggerEvent, Event: "thread-opened"}, false},
	}
	for _, tc := range cases {
		if got := tc.cond.QualifiesTrigger(tc.t); got != tc.want {
			t.Errorf("%s.QualifiesTrigger(%+v) = %v, want %v", tc.cond.Name, tc.t, got, tc.want)
		}
	}
}

// cond is the condition named name, failing the test when there is none.
func cond(t *testing.T, name string) Cond {
	t.Helper()
	c, ok := Condition(name)
	if !ok {
		t.Fatalf("no condition %s", name)
	}
	return c
}

// ADR-0211 §3: each trigger condition with a match holds on an occurrence of
// its event that matches, and not on one that does not, nor on another event.
func TestMatches(t *testing.T) {
	cases := []struct {
		name string
		o    Occurrence
		want bool
	}{
		{"to-review", Occurrence{Event: "moved", To: "review"}, true},
		{"to-review", Occurrence{Event: "moved", To: "done"}, false},
		{"to-review", Occurrence{Event: "edited", To: "review"}, false},
		{"to-backlog", Occurrence{Event: "moved", To: "backlog"}, true},
		{"to-backlog", Occurrence{Event: "moved", To: "ready"}, false},
		{"to-backlog", Occurrence{Event: "cancelled", To: "backlog"}, false},
		{"draft-changed", Occurrence{Event: "edited", Draft: true}, true},
		{"draft-changed", Occurrence{Event: "edited"}, false},
		{"draft-changed", Occurrence{Event: "moved", Draft: true}, false},
		{"awaits-orchestrator", Occurrence{Event: "thread", Awaits: "you"}, true},
		{"awaits-orchestrator", Occurrence{Event: "thread", Awaits: "operator"}, true},
		{"awaits-orchestrator", Occurrence{Event: "thread", Awaits: "planner-S-0001"}, false},
		{"awaits-orchestrator", Occurrence{Event: "thread"}, false},
		{"awaits-orchestrator", Occurrence{Event: "thread-opened", Awaits: "you"}, false},
	}
	for _, tc := range cases {
		holds, implemented := cond(t, tc.name).Matches(tc.o)
		if !implemented || holds != tc.want {
			t.Errorf("%s.Matches(%+v) = %v, %v; want %v, true", tc.name, tc.o, holds, implemented, tc.want)
		}
	}
	for _, name := range []string{"draft-only", "under-folder", "cost-inputs"} {
		if holds, implemented := cond(t, name).Matches(Occurrence{Event: "edited", Draft: true}); holds || implemented {
			t.Errorf("%s.Matches = %v, %v; want no match implemented", name, holds, implemented)
		}
	}
}

// staleStory is a ready story planned at 09:00, changed by change when it is
// not nil.
func staleStory(change func(*workitem.Item)) *workitem.Item {
	it := &workitem.Item{ID: "S-0001", Type: workitem.Story, Status: workitem.Ready,
		Forecast: &workitem.Forecast{Duration: "2h", Delivery: "2026-10-05T00:00:00Z", By: "planner-S-0001", At: "2026-10-04T09:00:00Z"}}
	if change != nil {
		change(it)
	}
	return it
}

// staleEdit is an edit at 10:00 by by of changed, by a story whose own agent
// is agent-S-0001.
func staleEdit(by string, changed ...string) Edit {
	return Edit{Changed: changed, By: by, At: "2026-10-04T10:00:00Z", Own: "agent-S-0001"}
}

// ADR-0084, ADR-0250 §5: stale-figures names the fields of an edit that
// leave a planned backlog or ready story's figures stale, as the replanner's
// edit trigger does, and none for an edit that does not.
func TestStaleFigures(t *testing.T) {
	value := 1000.0
	cases := []struct {
		name string
		it   *workitem.Item
		e    Edit
		want []string
	}{
		{"goal", staleStory(nil), staleEdit("alex", "goal"), []string{"goal"}},
		{"criteria", staleStory(nil), staleEdit("alex", "criteria"), []string{"criteria"}},
		{"touches", staleStory(nil), staleEdit("alex", "touches"), []string{"touches"}},
		{"several, in their order", staleStory(nil), staleEdit("alex", "touches", "goal", "notes"), []string{"goal", "touches"}},
		{"no one named", staleStory(nil), staleEdit("", "goal"), []string{"goal"}},
		{"in backlog", staleStory(func(it *workitem.Item) { it.Status = workitem.Backlog }), staleEdit("alex", "goal"), []string{"goal"}},
		{"a value older than the edit, no forecast", staleStory(func(it *workitem.Item) {
			it.Forecast, it.CostOfDelay = nil, &workitem.CostOfDelay{Value: &value, By: "alex", At: "2026-10-04T09:00:00Z"}
		}), staleEdit("alex", "criteria"), []string{"criteria"}},
		{"inputs changed after the value", staleStory(func(it *workitem.Item) {
			it.CostOfDelay = &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T10:00:00Z"}, Value: &value, By: "planner-S-0001", At: "2026-10-04T09:00:00Z"}
		}), staleEdit("alex", "cost_of_delay"), []string{"cost_of_delay"}},
		{"a value written alone", staleStory(func(it *workitem.Item) {
			it.CostOfDelay = &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T08:00:00Z"}, Value: &value, By: "alex", At: "2026-10-04T10:00:00Z"}
		}), staleEdit("alex", "cost_of_delay"), nil},
		{"a field outside the three", staleStory(nil), staleEdit("alex", "title", "notes", "body", "forecast"), nil},
		{"a planner's edit", staleStory(nil), staleEdit("planner-S-0001", "goal"), nil},
		{"another story's planner's edit", staleStory(nil), staleEdit("planner-S-0002", "goal"), nil},
		{"the story's own agent's edit", staleStory(nil), staleEdit("agent-S-0001", "touches"), nil},
		{"another story's agent's edit", staleStory(nil), staleEdit("agent-S-0002", "touches"), []string{"touches"}},
		{"no own agent named", staleStory(nil), Edit{Changed: []string{"touches"}, By: "agent-S-0001", At: "2026-10-04T10:00:00Z"}, []string{"touches"}},
		{"never planned", staleStory(func(it *workitem.Item) { it.Forecast = nil }), staleEdit("alex", "goal"), nil},
		{"never planned, inputs only", staleStory(func(it *workitem.Item) {
			it.Forecast, it.CostOfDelay = nil, &workitem.CostOfDelay{Inputs: &workitem.CostInputs{RevenuePerWeek: &value, By: "alex", At: "2026-10-04T10:00:00Z"}}
		}), staleEdit("alex", "cost_of_delay"), nil},
		{"a fresh forecast", staleStory(func(it *workitem.Item) { it.Forecast.At = "2026-10-04T11:00:00Z" }), staleEdit("alex", "goal"), nil},
		{"a forecast in the edit's second", staleStory(func(it *workitem.Item) { it.Forecast.At = "2026-10-04T10:00:00Z" }), staleEdit("alex", "goal"), nil},
		{"a forecast with no stamp", staleStory(func(it *workitem.Item) { it.Forecast.At = "" }), staleEdit("alex", "goal"), nil},
		{"an edit with no stamp", staleStory(nil), Edit{Changed: []string{"goal"}, By: "alex"}, nil},
		{"in progress", staleStory(func(it *workitem.Item) { it.Status = workitem.InProgress }), staleEdit("alex", "goal"), nil},
		{"archived", staleStory(func(it *workitem.Item) { it.Archived = true }), staleEdit("alex", "goal"), nil},
		{"a task", staleStory(func(it *workitem.Item) { it.Type = workitem.Task }), staleEdit("alex", "touches"), nil},
		{"an item not there", nil, staleEdit("alex", "goal"), nil},
	}
	stale := cond(t, "stale-figures")
	for _, c := range cases {
		if got := StaleFigures(c.it, c.e); !slices.Equal(got, c.want) {
			t.Errorf("%s: StaleFigures = %q, want %q", c.name, got, c.want)
		}
		holds, implemented := stale.Matches(Occurrence{Event: "edited", Item: c.it, Edit: c.e})
		if !implemented || holds != (len(c.want) > 0) {
			t.Errorf("%s: stale-figures.Matches = %v, %v; want %v, true", c.name, holds, implemented, len(c.want) > 0)
		}
	}
	// only on an edit
	if holds, _ := stale.Matches(Occurrence{Event: "moved", Item: staleStory(nil), Edit: staleEdit("alex", "goal")}); holds {
		t.Error("stale-figures holds on a move")
	}
}

// ADR-0084, ADR-0250 §5: replan-agent holds on an acceptance, a
// cancellation, or a reorder while planning.replan is agent, and gives each
// story whose delivery moved, once, rather than the event's item.
func TestReplanAgent(t *testing.T) {
	c := cond(t, "replan-agent")
	for _, event := range []string{"accepted", "cancelled", "reordered"} {
		for replan, want := range map[string]bool{"agent": true, "deterministic": false, "never": false, "": false} {
			if holds, implemented := c.Matches(Occurrence{Event: event, Replan: replan}); !implemented || holds != want {
				t.Errorf("replan-agent.Matches(%s, replan %q) = %v, %v; want %v, true", event, replan, holds, implemented, want)
			}
		}
	}
	for _, event := range []string{"moved", "edited", "ready", ""} {
		if holds, _ := c.Matches(Occurrence{Event: event, Replan: "agent"}); holds {
			t.Errorf("replan-agent holds on %q", event)
		}
	}
	moved := []string{"S-0003", "S-0002", "S-0003"}
	if got := ReplanAgent("agent", moved); !slices.Equal(got, []string{"S-0003", "S-0002"}) {
		t.Errorf("ReplanAgent(agent) = %v, want the stories moved in their order, once each", got)
	}
	if got := ReplanAgent("agent", nil); len(got) != 0 {
		t.Errorf("ReplanAgent(agent) with nothing moved = %v, want none", got)
	}
	for _, replan := range []string{"deterministic", "never", ""} {
		if got := ReplanAgent(replan, moved); len(got) != 0 {
			t.Errorf("ReplanAgent(%q) = %v, want none", replan, got)
		}
	}
}

// ADR-0250 §5: each-ready-story holds on a schedule and gives the ready
// stories in pull order, once each, and none when the ready column is empty.
func TestEachReadyStory(t *testing.T) {
	c := cond(t, "each-ready-story")
	if holds, implemented := c.Matches(Occurrence{Schedule: "planning.schedule"}); !holds || !implemented {
		t.Errorf("each-ready-story.Matches(schedule) = %v, %v; want true, true", holds, implemented)
	}
	for _, o := range []Occurrence{{}, {Event: "edited"}, {Event: "moved", Schedule: "planning.schedule"}} {
		if holds, implemented := c.Matches(o); holds || !implemented {
			t.Errorf("each-ready-story.Matches(%+v) = %v, %v; want false, true", o, holds, implemented)
		}
	}
	// a condition on events holds on no schedule
	if holds, _ := cond(t, "to-review").Matches(Occurrence{Schedule: "planning.schedule", To: "review"}); holds {
		t.Error("to-review holds on a schedule")
	}
	story := func(id, status string) *workitem.Item {
		return &workitem.Item{ID: id, Type: workitem.Story, Status: status}
	}
	archived := story("S-0006", workitem.Ready)
	archived.Archived = true
	items := []*workitem.Item{
		story("S-0001", workitem.Ready), story("S-0002", workitem.Ready), story("S-0003", workitem.Backlog),
		story("S-0004", workitem.InProgress), story("S-0005", workitem.Ready), archived,
		{ID: "T-0001", Type: workitem.Task, Status: workitem.Ready},
	}
	order := []string{"S-0002", "S-0003", "S-0001", "S-0002"}
	if got := EachReadyStory(order, items); !slices.Equal(got, []string{"S-0002", "S-0001", "S-0005"}) {
		t.Errorf("EachReadyStory = %v, want the ready stories in pull order, then those it does not name", got)
	}
	if got := EachReadyStory(order, []*workitem.Item{story("S-0003", workitem.Backlog), story("S-0004", workitem.Review)}); len(got) != 0 {
		t.Errorf("EachReadyStory on an empty ready column = %v, want none", got)
	}
	if got := EachReadyStory(nil, nil); len(got) != 0 {
		t.Errorf("EachReadyStory with no items = %v, want none", got)
	}
}

func TestGrantsPlain(t *testing.T) {
	cases := []struct {
		call, perm string
		want       bool
	}{
		{"thread_open", "", true},
		{"thread_reply", "", true},
		{"item_edit", "", false},
		{"item_edit", FinalizeDrafts, false},
		{"item_new", FinalizeDrafts, true},
		{"item_new", "", false},
		{"queue_add", AcceptReviews, true},
		{"queue_add", Publish, false},
		{"flai queue add", AcceptReviews, false},
		{"release_publish", Publish, true},
		{"flai subproject set", EditSubprojects, true},
	}
	for _, tc := range cases {
		if got := GrantsPlain(tc.call, tc.perm); got != tc.want {
			t.Errorf("GrantsPlain(%q, %q) = %v, want %v", tc.call, tc.perm, got, tc.want)
		}
	}
}

func TestForbids(t *testing.T) {
	cases := []struct {
		call, when string
		worksStory bool
		want       string // the forbidden call reached, empty for none
	}{
		{"flai push", "", false, "flai push"},
		{"flai release", "", false, "flai release"},
		{"flai queue *", "", false, "flai queue remove"},
		{"flai queue", "", false, ""},
		{"flai accept", "", false, "flai accept"},
		{"flai accept", "accept-as-self", false, ""},
		{"flai order", "by-policy", false, ""},
		{"item_move", "to-ready", false, ""},
		{"item_move", "", false, "item_move"},
		{"story_start", "", false, "story_start"},
		{"story_start", "", true, ""},
		{"flai stream *", "", false, "flai stream answer"},
		{"flai stream *", "", true, ""},
		{"flai move", "accept-as-self", false, ""},
		{"thread_open", "", false, ""},
	}
	for _, tc := range cases {
		f, ok := Forbids(tc.call, tc.when, tc.worksStory)
		if got := map[bool]string{true: f.Call}[ok]; got != tc.want {
			t.Errorf("Forbids(%q, %q, %v) = %q, want %q", tc.call, tc.when, tc.worksStory, got, tc.want)
		}
	}
}

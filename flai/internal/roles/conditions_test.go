package roles

import (
	"slices"
	"testing"
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
		if c, ok := Condition(name); !ok || len(c.Events) == 0 {
			t.Errorf("match %s is for no condition that qualifies an event", name)
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
	for _, name := range []string{"stale-figures", "replan-agent", "each-ready-story", "draft-only", "under-folder"} {
		if holds, implemented := cond(t, name).Matches(Occurrence{Event: "edited", Draft: true}); holds || implemented {
			t.Errorf("%s.Matches = %v, %v; want no match implemented", name, holds, implemented)
		}
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

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
		targets := 0
		for _, has := range []bool{len(c.Calls) > 0, c.OnTrigger(), c.Files} {
			if has {
				targets++
			}
		}
		if targets != 1 {
			t.Errorf("condition %s qualifies %d kinds of entry, want one", c.Name, targets)
		}
		if c.Permission != "" && !slices.Contains(Permissions, c.Permission) {
			t.Errorf("condition %s is gated by unknown permission %s", c.Name, c.Permission)
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
		if !ok || c.Permission != g.Permission || !c.Qualifies(g.Call) {
			t.Errorf("gate %+v names a condition that does not qualify its call under its permission", g)
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
	}
	for _, tc := range cases {
		if got := tc.cond.QualifiesTrigger(tc.t); got != tc.want {
			t.Errorf("%s.QualifiesTrigger(%+v) = %v, want %v", tc.cond.Name, tc.t, got, tc.want)
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

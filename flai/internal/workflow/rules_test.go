package workflow

import (
	"slices"
	"strings"
	"testing"
)

// S-0738, S-0807: the vocabulary names ADR-0301 §4's rules, each with the
// types it applies to and its default severity, in the ADR's order, then
// ADR-0312 §2's verified, for stories, refusing by default.
func TestVocabularyIsTheADRs(t *testing.T) {
	want := []struct {
		name, severity string
		types          []string
	}{
		{"criteria-present", SeverityRefuse, []string{TypeStory}},
		{"not-draft", SeverityRefuse, []string{TypeStory}},
		{"has-task", SeverityRefuse, []string{TypeStory}},
		{"no-open-question", SeverityRefuse, []string{TypeStory}},
		{"clean-worktree", SeverityRefuse, []string{TypeStory}},
		{"criteria-ticked", SeverityRefuse, []string{TypeStory}},
		{"kluster-verified", SeverityWarn, []string{TypeStory}},
		{"children-closed", SeverityRefuse, []string{TypeEpic, TypeStory}},
		{"reason-given", SeverityRefuse, Types},
		{"parent-open", SeverityRefuse, Types},
		{"not-blocked", SeverityWarn, Types},
		{"not-held", SeverityWarn, []string{TypeStory}},
		{"after-done", SeverityWarn, []string{TypeStory}},
		{"task-after", SeverityWarn, []string{TypeTask}},
		{"verified", SeverityRefuse, []string{TypeStory}},
	}
	if len(Vocabulary) != len(want) {
		t.Fatalf("vocabulary has %d rules, want %d: %v", len(Vocabulary), len(want), RuleNames())
	}
	for i, w := range want {
		r := Vocabulary[i]
		if r.Name != w.name || r.Severity != w.severity || !slices.Equal(r.Types, w.types) || r.Checks == "" {
			t.Errorf("rule %d = %+v, want %s, %s, %v, saying what it checks", i, r, w.name, w.severity, w.types)
		}
		got, ok := LookupRule(w.name)
		if !ok || got.Name != w.name {
			t.Errorf("LookupRule(%q) = %+v, %v", w.name, got, ok)
		}
	}
	if _, ok := LookupRule("no-such-rule"); ok {
		t.Error("LookupRule found a rule the vocabulary lacks")
	}
	if r, _ := LookupRule(RuleChildrenClosed); !r.Applies(TypeEpic) || r.Applies(TypeTask) {
		t.Errorf("children-closed applies to %v, want epic and story only", r.Types)
	}
}

// S-0738: an entry's severity is the one it sets, else its rule's default.
func TestRuleRefSeverityOf(t *testing.T) {
	cases := []struct {
		ref  RuleRef
		want string
	}{
		{RuleRef{Rule: RuleHasTask}, SeverityRefuse},
		{RuleRef{Rule: RuleNotBlocked}, SeverityWarn},
		{RuleRef{Rule: RuleCriteriaTicked, Severity: SeverityWarn}, SeverityWarn},
		{RuleRef{Rule: RuleNotHeld, Severity: SeverityRefuse}, SeverityRefuse},
	}
	for _, tc := range cases {
		if got := tc.ref.SeverityOf(); got != tc.want {
			t.Errorf("%+v.SeverityOf() = %q, want %q", tc.ref, got, tc.want)
		}
	}
}

// S-0738: a transition naming a rule the vocabulary lacks is refused when the
// file is parsed, under workflow.rule, naming the rule and the transition.
func TestParseRefusesAnUnknownRule(t *testing.T) {
	content := strings.Replace(minimal, "to: done}", "to: done, rules: [has-task, sign-off]}", 1)
	_, p := Parse("design/workflows/board.md", content)
	if p == nil {
		t.Fatal("parsed, want a problem")
	}
	if p.Code != CodeRule || p.Line != 9 || p.Field != "transitions[0].rules[1]" {
		t.Errorf("got %s (code %s, line %d, field %q), want workflow.rule on line 9, field transitions[0].rules[1]", p, p.Code, p.Line, p.Field)
	}
	for _, says := range []string{"sign-off", "transitions[0] (backlog to done)", "criteria-present", "task-after"} {
		if !strings.Contains(p.Message, says) {
			t.Errorf("message %q does not say %q", p.Message, says)
		}
	}
	if _, p := Parse("design/workflows/board.md", strings.Replace(minimal, "to: done}", "to: done, rules: [has-task, {rule: not-held, severity: refuse}]}", 1)); p != nil {
		t.Errorf("a file naming the vocabulary's rules was refused: %s", p)
	}
}

// reviewBoard is a board with a move to review and an acceptance, and no
// key of ADR-0312's: the acceptance lists verified.
const reviewBoard = `---
title: Board
updated: 2026-11-18
status: active
states:
  - {name: backlog, kind: queue, initial: true}
  - {name: in-progress, kind: active}
  - {name: review, kind: waiting}
  - {name: done, kind: accepted}
transitions:
  - {from: backlog, to: in-progress}
  - {from: in-progress, to: review, rules: [has-task]}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, verified]}
---
`

// parsed is content parsed, failing the test when it does not parse.
func parsed(t *testing.T, content string) Workflow {
	t.Helper()
	w, p := Parse("design/workflows/board.md", content)
	if p != nil {
		t.Fatalf("the test's workflow does not parse: %s", p)
	}
	return w
}

// moveRules is w's MoveRules for a story's move from one state to the other.
func moveRules(t *testing.T, w Workflow, from, to string) []RuleRef {
	t.Helper()
	tr, ok := w.Find(TypeStory, from, to)
	if !ok {
		t.Fatalf("no move of a story from %s to %s", from, to)
	}
	return w.MoveRules(TypeStory, tr)
}

// S-0807: the verify a move reads is its transition's own; a file that gives
// none of ADR-0312's keys reads the baseline's on the same move (D-0669),
// and a file that gives one reads nothing on a move without verify.
func TestVerifyOf(t *testing.T) {
	verifyOf := func(w Workflow, from, to string) *Verify {
		tr, ok := w.Find(TypeStory, from, to)
		if !ok {
			t.Fatalf("no move of a story from %s to %s", from, to)
		}
		return w.VerifyOf(TypeStory, tr)
	}
	base := Baseline()
	if v := verifyOf(base, "in-progress", "review"); v == nil || !v.Suites.All || v.Severity != SeverityWarn {
		t.Errorf("the baseline's move to review reads %+v, want every suite at warn", v)
	}
	if v := verifyOf(base, "ready", "in-progress"); v != nil {
		t.Errorf("the baseline's pull reads %+v, want none", v)
	}
	old := parsed(t, reviewBoard)
	if v := verifyOf(old, "in-progress", "review"); v == nil || !v.Suites.All || v.Severity != SeverityWarn {
		t.Errorf("a file without the suite keys reads %+v on the move to review, want the baseline's, every suite at warn", v)
	}
	if v := verifyOf(old, "backlog", "in-progress"); v != nil {
		t.Errorf("a move the baseline lacks reads %+v, want none", v)
	}
	own := parsed(t, strings.Replace(reviewBoard, "to: done, types", "to: done, verify: [static], types", 1))
	if v := verifyOf(own, "in-progress", "review"); v != nil {
		t.Errorf("a file that gives verify on another move reads %+v on the move to review, want none", v)
	}
	if v := verifyOf(own, "review", "done"); v == nil || !slices.Equal(v.Suites.Names, []string{"static"}) {
		t.Errorf("the acceptance reads %+v, want its own [static]", v)
	}
}

// S-0807: a move by a transition with verify runs verified last at verify's
// severity, or at the place the transition lists it with the entry's
// severity, else verify's; a listed verified with no verify needs every
// suite at refuse; and the acceptance runs none, which S-0808's gates own.
func TestMoveRules(t *testing.T) {
	base := Baseline()
	got := moveRules(t, base, "in-progress", "review")
	want := []RuleRef{{Rule: RuleHasTask}, {Rule: RuleNoOpenQuestion}, {Rule: RuleCleanWorktree},
		{Rule: RuleCriteriaTicked, Severity: SeverityWarn}, {Rule: RuleKlusterVerified}, {Rule: RuleVerified, Severity: SeverityWarn}}
	if !slices.Equal(got, want) {
		t.Errorf("the baseline's move to review runs %v, want %v", got, want)
	}
	if tr, _ := base.Find(TypeStory, "in-progress", "review"); len(tr.Rules) != len(want)-1 {
		t.Errorf("MoveRules changed the transition's own rules: %v", tr.Rules)
	}
	if got := moveRules(t, base, "review", "done"); slices.ContainsFunc(got, func(r RuleRef) bool { return r.Rule == RuleVerified }) {
		t.Errorf("the baseline's acceptance runs verified: %v", got)
	}
	if got := moveRules(t, base, "ready", "in-progress"); slices.ContainsFunc(got, func(r RuleRef) bool { return r.Rule == RuleVerified }) {
		t.Errorf("a move without verify runs verified: %v", got)
	}

	refuse := parsed(t, strings.Replace(reviewBoard, "rules: [has-task]}", "rules: [has-task], verify: [static, behavior]}", 1))
	if got, want := moveRules(t, refuse, "in-progress", "review"), []RuleRef{{Rule: RuleHasTask}, {Rule: RuleVerified}}; !slices.Equal(got, want) {
		t.Errorf("verify at refuse runs %v, want %v", got, want)
	}
	if got := moveRules(t, refuse, "in-progress", "review")[1].SeverityOf(); got != SeverityRefuse {
		t.Errorf("verify without a severity refuses, got %s", got)
	}
	if got, want := moveRules(t, refuse, "review", "done"), []RuleRef{{Rule: RuleChildrenClosed}}; !slices.Equal(got, want) {
		t.Errorf("the acceptance that lists verified runs %v, want %v", got, want)
	}

	placed := parsed(t, strings.Replace(reviewBoard, "rules: [has-task]}", "rules: [verified, has-task], verify: {suites: all, severity: warn}}", 1))
	if got, want := moveRules(t, placed, "in-progress", "review"), []RuleRef{{Rule: RuleVerified, Severity: SeverityWarn}, {Rule: RuleHasTask}}; !slices.Equal(got, want) {
		t.Errorf("a listed verified runs %v, want %v", got, want)
	}
	entry := parsed(t, strings.Replace(reviewBoard, "rules: [has-task]}", "rules: [has-task, {rule: verified, severity: refuse}], verify: {suites: all, severity: warn}}", 1))
	if got, want := moveRules(t, entry, "in-progress", "review"), []RuleRef{{Rule: RuleHasTask}, {Rule: RuleVerified, Severity: SeverityRefuse}}; !slices.Equal(got, want) {
		t.Errorf("a listed verified with a severity of its own runs %v, want %v", got, want)
	}
	bare := parsed(t, strings.Replace(strings.Replace(reviewBoard, "to: done, types", "to: done, verify: all, types", 1), "rules: [has-task]}", "rules: [has-task, verified]}", 1))
	got = moveRules(t, bare, "in-progress", "review")
	if want := []RuleRef{{Rule: RuleHasTask}, {Rule: RuleVerified}}; !slices.Equal(got, want) || got[1].SeverityOf() != SeverityRefuse {
		t.Errorf("a listed verified with no verify runs %v, want %v at refuse", got, want)
	}
	if tr, _ := bare.Find(TypeStory, "in-progress", "review"); bare.VerifyOf(TypeStory, tr) != nil {
		t.Error("a listed verified with no verify reads a verify; the move's check needs every suite")
	}
}

package workflow

import (
	"slices"
	"strings"
	"testing"
)

// S-0738: the vocabulary names ADR-0301 §4's rules, each with the types it
// applies to and its default severity, in the ADR's order.
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

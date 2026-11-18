package workflow

import "slices"

// CodeRule is the code flai reports a transition naming a rule the
// vocabulary lacks under (ADR-0301 §4).
const CodeRule = "workflow.rule"

// The rules a transition may list (ADR-0301 §4), and verified, which a
// transition's verify gives (ADR-0312 §2). flai implements each; a workflow
// file names one and cannot define one.
const (
	RuleCriteriaPresent = "criteria-present"
	RuleNotDraft        = "not-draft"
	RuleHasTask         = "has-task"
	RuleNoOpenQuestion  = "no-open-question"
	RuleCleanWorktree   = "clean-worktree"
	RuleCriteriaTicked  = "criteria-ticked"
	RuleKlusterVerified = "kluster-verified"
	RuleChildrenClosed  = "children-closed"
	RuleReasonGiven     = "reason-given"
	RuleParentOpen      = "parent-open"
	RuleNotBlocked      = "not-blocked"
	RuleNotHeld         = "not-held"
	RuleAfterDone       = "after-done"
	RuleTaskAfter       = "task-after"
	RuleVerified        = "verified"
)

// Rule is one rule of the vocabulary: its name, what it checks, the item
// types it applies to, and the severity it has when a transition sets none.
type Rule struct {
	Name     string   `json:"name"`
	Checks   string   `json:"checks"`
	Types    []string `json:"types"`
	Severity string   `json:"severity"`
}

// Applies reports whether the rule checks an item of type typ.
func (r Rule) Applies(typ string) bool {
	return slices.Contains(r.Types, typ)
}

// Vocabulary is every rule flai implements, in ADR-0301 §4's order, then
// ADR-0312's.
var Vocabulary = []Rule{
	{RuleCriteriaPresent, "The body has ## Acceptance criteria with at least one checkbox (ADR-0021).", []string{TypeStory}, SeverityRefuse},
	{RuleNotDraft, "The story is not a draft, or the move finalizes it.", []string{TypeStory}, SeverityRefuse},
	{RuleHasTask, "The story has at least one task.", []string{TypeStory}, SeverityRefuse},
	{RuleNoOpenQuestion, "The story's narrative has no open question.", []string{TypeStory}, SeverityRefuse},
	{RuleCleanWorktree, "Nothing is uncommitted in the story's worktree; a worktree git cannot read warns (I-0017).", []string{TypeStory}, SeverityRefuse},
	{RuleCriteriaTicked, "Every acceptance criterion is ticked.", []string{TypeStory}, SeverityRefuse},
	{RuleKlusterVerified, "The story's namespace is deployed at its branch's head and its checks passed there (ADR-0292).", []string{TypeStory}, SeverityWarn},
	{RuleChildrenClosed, "Every child is in an accepted or cancelled state.", []string{TypeEpic, TypeStory}, SeverityRefuse},
	{RuleReasonGiven, "The move gives a reason; a cascade gives its parent's.", Types, SeverityRefuse},
	{RuleParentOpen, "The item's parent is not in a cancelled state.", Types, SeverityRefuse},
	{RuleNotBlocked, "The item has no open blocked interval.", Types, SeverityWarn},
	{RuleNotHeld, "The story is not held by an overlap or by after (ADR-0046, ADR-0230).", []string{TypeStory}, SeverityWarn},
	{RuleAfterDone, "Every story the story names in after is accepted (ADR-0279).", []string{TypeStory}, SeverityWarn},
	{RuleTaskAfter, "Every task the task names in after is closed (S-0176).", []string{TypeTask}, SeverityWarn},
	{RuleVerified, "The story's last flai verify passed, at its branch's head, every tier its paths select of the suites the transition's verify lists; a transition with verify runs it at verify's severity (ADR-0312 §2).", []string{TypeStory}, SeverityRefuse},
}

// LookupRule is the vocabulary's rule of that name.
func LookupRule(name string) (Rule, bool) {
	i := slices.IndexFunc(Vocabulary, func(r Rule) bool { return r.Name == name })
	if i < 0 {
		return Rule{}, false
	}
	return Vocabulary[i], true
}

// RuleNames are the vocabulary's rule names, in its order.
func RuleNames() []string {
	out := make([]string, len(Vocabulary))
	for i, r := range Vocabulary {
		out[i] = r.Name
	}
	return out
}

// SeverityOf is the severity the entry gives its rule: the one it sets, else
// the rule's default.
func (r RuleRef) SeverityOf() string {
	if r.Severity != "" {
		return r.Severity
	}
	if rule, ok := LookupRule(r.Rule); ok {
		return rule.Severity
	}
	return SeverityRefuse
}

// VerifyOf is the verify a move by t of an item of type typ reads: t's own,
// or, in a file that gives none of ADR-0312's keys, the baseline's on the
// same move, so that such a file runs what it ran before (ADR-0312 §8,
// D-0669); nil when neither gives one.
func (w Workflow) VerifyOf(typ string, t Transition) *Verify {
	if t.Verify != nil || w.UsesSuiteKeys() {
		return t.Verify
	}
	for _, from := range t.From {
		if b, ok := baseline().Find(typ, from, t.To); ok && b.Verify != nil {
			v := *b.Verify
			v.Suites.Names = slices.Clone(v.Suites.Names)
			return &v
		}
	}
	return nil
}

// MoveRules are the rules a move by t of an item of type typ runs, in their
// order: those t lists, and verified where VerifyOf gives a verify, at the
// place t lists it or else last, with verify's severity unless the entry sets
// its own. A file that lists verified on a move with no verify needs every
// suite. The transition with the accept effect runs no verified, listed or
// not: what the acceptance needs of the record is its own gates', which
// S-0808 makes read its verify (ADR-0312 §3).
func (w Workflow) MoveRules(typ string, t Transition) []RuleRef {
	isVerified := func(r RuleRef) bool { return r.Rule == RuleVerified }
	refs := slices.Clone(t.Rules)
	if t.HasEffect(EffectAccept) {
		return slices.DeleteFunc(refs, isVerified)
	}
	v, listed := w.VerifyOf(typ, t), slices.IndexFunc(refs, isVerified)
	switch {
	case v == nil:
	case listed < 0:
		refs = append(refs, RuleRef{Rule: RuleVerified, Severity: v.Severity})
	case refs[listed].Severity == "":
		refs[listed].Severity = v.Severity
	}
	return refs
}

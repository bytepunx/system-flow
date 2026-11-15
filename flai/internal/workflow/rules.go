package workflow

import "slices"

// CodeRule is the code flai reports a transition naming a rule the
// vocabulary lacks under (ADR-0301 §4).
const CodeRule = "workflow.rule"

// The rules a transition may list (ADR-0301 §4). flai implements each; a
// workflow file names one and cannot define one.
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

// Vocabulary is every rule flai implements, in ADR-0301 §4's order.
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

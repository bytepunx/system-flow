package roles

import (
	"slices"
	"strings"
)

// The operator's permissions, the <p> of orchestration.permissions.<p> and
// roles.<name>.permissions.<p> (ADR-0207 §5).
const (
	PlanBacklogEpics   = "plan_backlog_epics"
	PlanBacklogStories = "plan_backlog_stories"
	FinalizeDrafts     = "finalize_drafts"
	PromoteToReady     = "promote_to_ready"
	OrderReady         = "order_ready"
	AnswerThreads      = "answer_threads"
	AcceptReviews      = "accept_reviews"
	Publish            = "publish"
	EditSubprojects    = "edit_subprojects"
)

// Permissions are every permission a role file's setting may name, in the
// Settings page's order (ADR-0207 §5).
var Permissions = []string{PlanBacklogEpics, PlanBacklogStories, FinalizeDrafts, PromoteToReady, OrderReady, AnswerThreads, AcceptReviews, Publish, EditSubprojects}

// Cond is a named condition flai implements, which a role file names on a
// call entry, a trigger entry, or, for under-folder, each files entry
// (ADR-0207 §5).
type Cond struct {
	Name string `json:"name"`
	// Calls are the calls it qualifies, as a role file writes them: an MCP
	// tool's name (item_new) or a flai command path (flai story new).
	Calls []string `json:"calls,omitempty"`
	// Events are the trigger events it qualifies, for a trigger condition.
	Events []string `json:"events,omitempty"`
	// Schedule is whether it qualifies a schedule trigger.
	Schedule bool `json:"schedule,omitempty"`
	// Files is whether it holds each files entry rather than a call.
	Files bool `json:"files,omitempty"`
	// Allows says in one line what the call or trigger may do under it.
	Allows string `json:"allows"`
	// Permission is the permission that gates it, empty when none does.
	Permission string `json:"permission,omitempty"`
}

// Qualifies reports whether the condition may stand on a call entry naming
// call.
func (c Cond) Qualifies(call string) bool { return slices.Contains(c.Calls, call) }

// OnTrigger reports whether the condition qualifies a trigger entry. A
// condition may qualify calls too, as to-backlog does (ADR-0211 §3).
func (c Cond) OnTrigger() bool { return len(c.Events) > 0 || c.Schedule }

// OnCall reports whether the condition qualifies a call entry.
func (c Cond) OnCall() bool { return len(c.Calls) > 0 }

// Occurrence is an event that happened, in roles' own terms, which a
// trigger condition matches (ADR-0211 §3).
type Occurrence struct {
	// Event is the event's name, one of Events.
	Event string
	// To is the state a moved item moved to.
	To string
	// Draft is whether an edited item is a story that is a draft after the
	// edit.
	Draft bool
	// Awaits is whom a thread event's thread awaits: "you" for the reader,
	// or "operator" (ADR-0181).
	Awaits string
}

// matches are the trigger conditions whose match flai implements against an
// occurrence, by name; the event is checked against the condition's Events
// beside it.
var matches = map[string]func(Occurrence) bool{
	"to-review":           func(o Occurrence) bool { return o.To == "review" },
	"to-backlog":          func(o Occurrence) bool { return o.To == "backlog" },
	"draft-changed":       func(o Occurrence) bool { return o.Draft },
	"awaits-orchestrator": func(o Occurrence) bool { return o.Awaits == "you" || o.Awaits == "operator" },
}

// Matches reports whether the condition holds for the occurrence o: o's
// event is one the condition qualifies and the condition's match holds.
// implemented is false for a condition with no match against an occurrence:
// a call or files condition, and the trigger conditions stale-figures,
// replan-agent, and each-ready-story, which flai judges otherwise.
// awaits-orchestrator reads the orchestrator as the reader, so a thread that
// awaits it ("you") or the operator holds.
func (c Cond) Matches(o Occurrence) (holds, implemented bool) {
	match, ok := matches[c.Name]
	if !ok {
		return false, false
	}
	return slices.Contains(c.Events, o.Event) && match(o), true
}

// QualifiesTrigger reports whether the condition may stand on the trigger
// entry t.
func (c Cond) QualifiesTrigger(t Trigger) bool {
	switch t.Kind {
	case TriggerEvent:
		return slices.Contains(c.Events, t.Event)
	case TriggerSchedule:
		return c.Schedule
	}
	return false
}

// Conditions are every condition a role file may name (ADR-0207 §5).
var Conditions = []Cond{
	{Name: "draft-only", Calls: []string{"item_new", "flai story new"}, Allows: "a story only as a draft (draft true, --draft); an epic or a task as asked"},
	{Name: "to-backlog", Calls: []string{"item_move", "flai move"}, Events: []string{"moved"}, Allows: "on a call, a move to backlog and nowhere else; on a trigger, only a move to backlog"},
	{Name: "ungated-edit", Calls: []string{"item_edit", "flai edit"}, Allows: "any edit that neither finalizes a draft (draft false, --no-draft) nor sets a cost of delay input"},
	{Name: "own-thread", Calls: []string{"thread_reply", "thread_resolve", "flai thread reply", "flai thread resolve"}, Allows: "on a thread the role opened, a follow-up, never a recommendation or a source, and resolving it"},
	{Name: "under-folder", Files: true, Allows: "an edit of a file under the folder, after .. and symbolic links are resolved, never of the folder itself"},
	{Name: "plan-epic", Calls: []string{"plan", "flai plan"}, Allows: "starting the planner on an epic", Permission: PlanBacklogEpics},
	{Name: "plan-story", Calls: []string{"plan", "flai plan"}, Allows: "starting the planner on a story", Permission: PlanBacklogStories},
	{Name: "cost-inputs", Calls: []string{"item_edit", "flai edit"}, Allows: "giving a story cost of delay inputs and changing nothing else: no value, no clearing (ADR-0119)", Permission: PlanBacklogStories},
	{Name: "planner-thread", Calls: []string{"thread_reply", "thread_resolve", "flai thread reply", "flai thread resolve"}, Allows: "on a thread a story's planner opened, an answer, with or without a source, or a recommendation, and resolving it (ADR-0119)", Permission: PlanBacklogStories},
	{Name: "finalize-draft", Calls: []string{"item_edit", "flai edit"}, Allows: "finalizing a draft and changing nothing else", Permission: FinalizeDrafts},
	{Name: "to-ready", Calls: []string{"item_move", "flai move"}, Allows: "a move to ready and nowhere else", Permission: PromoteToReady},
	{Name: "by-policy", Calls: []string{"flai order"}, Allows: "flai order --by <policy> --apply, placing no story by hand", Permission: OrderReady},
	{Name: "answer-threads", Calls: []string{"thread_reply", "flai thread reply"}, Allows: "on a thread the role did not open, a recommendation while the setting is recommend, or an answer citing a source while it is autonomous; never resolving it (ADR-0090)", Permission: AnswerThreads},
	{Name: "accept-as-self", Calls: []string{"flai accept", "flai move"}, Allows: "accepting a story in review, or moving it to done, as the role itself (--by <name>), never with --batch (ADR-0093)", Permission: AcceptReviews},
	{Name: "enqueue-as-self", Calls: []string{"flai queue add"}, Allows: "enqueuing a story for the merge queue as the role itself (--by <name>)", Permission: AcceptReviews},
	{Name: "stale-figures", Events: []string{"edited"}, Allows: "only an edit of a backlog or ready story's goal, criteria, or touches by someone other than a planner, when its forecast or cost of delay value is older than the edit"},
	{Name: "replan-agent", Events: []string{"accepted", "cancelled", "reordered"}, Allows: "only while planning.replan is agent, for each story whose delivery moved"},
	{Name: "each-ready-story", Schedule: true, Allows: "once for every ready story"},
	{Name: "to-review", Events: []string{"moved"}, Allows: "only a move to review"},
	{Name: "draft-changed", Events: []string{"edited"}, Allows: "only an edit of a story that is a draft after the edit"},
	{Name: "awaits-orchestrator", Events: []string{"thread"}, Allows: "only a thread that awaits the orchestrator or the operator (ADR-0181)"},
}

// Condition returns the condition named name.
func Condition(name string) (Cond, bool) {
	i := slices.IndexFunc(Conditions, func(c Cond) bool { return c.Name == name })
	if i < 0 {
		return Cond{}, false
	}
	return Conditions[i], true
}

// ConditionNames are the names of every condition, in the table's order.
func ConditionNames() []string {
	out := make([]string, len(Conditions))
	for i, c := range Conditions {
		out[i] = c.Name
	}
	return out
}

// Gate is one way a call only the operator's settings allow is granted: the
// call, the permission whose setting allows it, and the condition it is
// named under, empty for the call named plainly (ADR-0207 §5).
type Gate struct {
	Call       string `json:"call"`
	Permission string `json:"permission"`
	Condition  string `json:"condition,omitempty"`
}

// Gated are the gated calls and how a role file grants each, the ADR's
// gated-calls table row by row (ADR-0207 §5).
var Gated = []Gate{
	{"plan", PlanBacklogEpics, "plan-epic"},
	{"flai plan", PlanBacklogEpics, "plan-epic"},
	{"plan", PlanBacklogStories, "plan-story"},
	{"flai plan", PlanBacklogStories, "plan-story"},
	{"item_edit", PlanBacklogStories, "cost-inputs"},
	{"flai edit", PlanBacklogStories, "cost-inputs"},
	{"item_edit", FinalizeDrafts, "finalize-draft"},
	{"flai edit", FinalizeDrafts, "finalize-draft"},
	{"item_new", FinalizeDrafts, ""},
	{"flai story new", FinalizeDrafts, ""},
	{"item_move", PromoteToReady, "to-ready"},
	{"flai move", PromoteToReady, "to-ready"},
	{"flai order", OrderReady, "by-policy"},
	{"thread_reply", AnswerThreads, "answer-threads"},
	{"flai thread reply", AnswerThreads, "answer-threads"},
	{"flai accept", AcceptReviews, "accept-as-self"},
	{"flai move", AcceptReviews, "accept-as-self"},
	{"queue_add", AcceptReviews, ""},
	{"flai queue add", AcceptReviews, "enqueue-as-self"},
	{"release_publish", Publish, ""},
	{"subproject_edit", EditSubprojects, ""},
	{"flai subproject add", EditSubprojects, ""},
	{"flai subproject set", EditSubprojects, ""},
	{"flai subproject remove", EditSubprojects, ""},
}

// PlainGated are the calls that, named plainly, would allow a gated call,
// so a role file names each plainly only inside a permission whose Gated row
// grants it plain; a thread reply is not among them, since a role replies on
// threads plainly (ADR-0207 §5).
var PlainGated = []string{
	"item_edit", "flai edit", "item_move", "flai move", "item_new", "flai story new",
	"plan", "flai plan", "flai order", "flai accept", "release_publish", "queue_add", "flai queue add",
	"subproject_edit", "flai subproject add", "flai subproject set", "flai subproject remove",
}

// GrantsPlain reports whether a role file may name call plainly inside a
// permission entry for p: call is not plain-gated, or a Gated row grants it
// plain under p.
func GrantsPlain(call, p string) bool {
	if !slices.Contains(PlainGated, call) {
		return true
	}
	return slices.Contains(Gated, Gate{Call: call, Permission: p})
}

// Forbidden is a call, or a form of one, that no setting allows (ADR-0207
// §5).
type Forbidden struct {
	// Call is the call as a role file writes it: an MCP tool, a flai command
	// path, or a git command.
	Call string `json:"call"`
	// Form is the form of the call no setting allows; empty for every form.
	Form string `json:"form,omitempty"`
	// Unless are the conditions that hold the call clear of Form.
	Unless []string `json:"unless,omitempty"`
	// StoryAgent marks a story's agent's call, forbidden only while the
	// role's works_story is false.
	StoryAgent bool `json:"story_agent,omitempty"`
}

// Never are the calls no setting allows, and a role file that names them
// anywhere is roles.grant (ADR-0207 §5). The command forms of a story's
// agent's calls stand beside its MCP tools.
var Never = []Forbidden{
	{Call: "flai release", Form: "any form but --evaluate"},
	{Call: "flai push"},
	{Call: "git push"},
	{Call: "git tag"},
	{Call: "flai accept", Form: "--batch", Unless: []string{"accept-as-self"}},
	{Call: "flai queue remove"},
	{Call: "flai queue hold"},
	{Call: "flai queue resume"},
	{Call: "flai order", Form: "placing a story by hand", Unless: []string{"by-policy"}},
	{Call: "flai thread confirm"},
	{Call: "item_move", Form: "to done", Unless: []string{"to-backlog", "to-ready"}},
	{Call: "shared_paths_edit"},
	{Call: "flai shared add"},
	{Call: "flai shared remove"},
	{Call: "flai forge push"},
	{Call: "story_start", StoryAgent: true},
	{Call: "task_done", StoryAgent: true},
	{Call: "criteria_tick", StoryAgent: true},
	{Call: "stream_state", StoryAgent: true},
	{Call: "message_get", StoryAgent: true},
	{Call: "message_send", StoryAgent: true},
	{Call: "message_reply", StoryAgent: true},
	{Call: "message_escalate", StoryAgent: true},
	{Call: "message_share", StoryAgent: true},
	{Call: "flai story start", StoryAgent: true},
	{Call: "flai task done", StoryAgent: true},
	{Call: "flai criteria tick", StoryAgent: true},
	{Call: "flai criteria untick", StoryAgent: true},
	{Call: "flai stream answer", StoryAgent: true},
	{Call: "flai stream log", StoryAgent: true},
	{Call: "flai stream open", StoryAgent: true},
	{Call: "flai stream restore", StoryAgent: true},
	{Call: "flai stream state", StoryAgent: true},
	{Call: "flai stream sync", StoryAgent: true},
	{Call: "flai message escalate", StoryAgent: true},
	{Call: "flai message reply", StoryAgent: true},
	{Call: "flai message send", StoryAgent: true},
	{Call: "flai message share", StoryAgent: true},
	{Call: "item_move", Form: "to in-progress or review", Unless: []string{"to-backlog", "to-ready"}, StoryAgent: true},
	{Call: "flai move", Form: "to in-progress or review", Unless: []string{"to-backlog", "to-ready", "accept-as-self"}, StoryAgent: true},
}

// covers reports whether a call entry naming call reaches target: the same
// call, or a flai <command> * entry and one of the command's subcommands.
func covers(call, target string) bool {
	if call == target {
		return true
	}
	prefix, ok := strings.CutSuffix(call, " *")
	return ok && strings.HasPrefix(target, prefix+" ")
}

// Forbids returns the never-granted entry a call entry naming call under
// condition when (empty for plain) reaches, for a role whose works_story is
// worksStory.
func Forbids(call, when string, worksStory bool) (Forbidden, bool) {
	for _, f := range Never {
		if !covers(call, f.Call) || (f.StoryAgent && worksStory) || (when != "" && slices.Contains(f.Unless, when)) {
			continue
		}
		return f, true
	}
	return Forbidden{}, false
}

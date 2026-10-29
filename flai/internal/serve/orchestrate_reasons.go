package serve

import (
	"fmt"
	"slices"
	"strings"
	"time"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator's reasons (ADR-0209 §1, §2, ADR-0211): what flai serve
// starts it for, the queue they wait in until a run takes them, and which
// of the inbox's events are reasons, by the orchestrator's role file's
// triggers. Nothing here reads or writes a file: the queue is kept by Dir
// (orchestrateQueues), and the events come from inbox.CatchUp.

// The kinds of a reason that is no event of the inbox's.
const (
	// reasonAsked is the operator's: orchestrate turned on, or flai serve
	// orchestrate start.
	reasonAsked = roles.TriggerAsked
	// reasonSchedule is orchestration.schedule coming round.
	reasonSchedule = roles.TriggerSchedule
)

// What an asked reason says.
const (
	askedTurnedOn = "orchestrate turned on"
	askedOperator = "the operator asked"
)

// serveCursor names the inbox cursor flai serve reads the orchestrator's
// events under, so that the orchestrator's own cursor is left where it is.
const serveCursor = "flai-serve-orchestrator"

// orchestrationSchedule is the setting a schedule trigger names for the
// orchestrator's periodic look.
const orchestrationSchedule = "orchestration.schedule"

// The most reasons a journal entry or an agent event's why names, and the
// most characters of each, before the rest are counted.
const (
	reasonsNamed  = 5
	reasonSaidMax = 120
)

// builtinTriggers are the orchestrator's triggers when its role file is
// missing, does not parse, or still says continuous: ADR-0209 §1 as
// ADR-0211 refines it.
var builtinTriggers = []roles.Trigger{
	{Kind: roles.TriggerAsked},
	{Kind: roles.TriggerEvent, Event: inbox.Thread, When: "awaits-orchestrator"},
	{Kind: roles.TriggerEvent, Event: workitem.Moved, When: "to-review"},
	{Kind: roles.TriggerEvent, Event: workitem.Moved, When: "to-backlog"},
	{Kind: roles.TriggerEvent, Event: inbox.Edited, When: "draft-changed"},
	{Kind: roles.TriggerEvent, Event: inbox.Planned},
	{Kind: roles.TriggerEvent, Event: inbox.ReleaseDue},
	{Kind: roles.TriggerEvent, Event: inbox.ReadyEmpty},
	{Kind: roles.TriggerEvent, Event: reordered},
	{Kind: roles.TriggerSchedule, Schedule: orchestrationSchedule},
}

// reordered is the kind of the inbox's event that the pull order changed.
const reordered = "reordered"

// orchestrateReason is one reason flai serve has to start the orchestrator.
type orchestrateReason struct {
	// Kind is the event's kind, reasonAsked, or reasonSchedule.
	Kind string `json:"kind"`
	// ID is the item's, the thread's, or "board"; the schedule's expression
	// for a schedule; empty for asked.
	ID    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
	// To is the state a moved item moved to, whom a thread awaits (you or
	// operator), what an edit changed, the items a plan wrote.
	To string `json:"to,omitempty"`
	// Summary is the reason in one line, as the inbox words it.
	Summary string `json:"summary"`
	// At is when it happened, and Queued when flai serve queued it.
	At     string `json:"at"`
	Queued string `json:"queued"`
}

// key is what makes two reasons one: the same kind, ID, and to.
func (r orchestrateReason) key() string { return r.Kind + "|" + r.ID + "|" + r.To }

// line is the reason in one line, for the run's FLAI_REASONS, its journal
// entry, and its started event: "moved S-0123 to review: <title>".
func (r orchestrateReason) line() string {
	title := ""
	if r.Title != "" {
		title = ": " + r.Title
	}
	var s string
	switch r.Kind {
	case reasonAsked:
		s = "asked: " + r.Summary
	case reasonSchedule:
		s = "schedule " + r.ID + " came round"
	case workitem.Moved:
		s = "moved " + r.ID + " to " + r.To + title
	case inbox.Thread:
		whom := "the orchestrator"
		if r.To == "operator" {
			whom = "the operator"
		}
		s = "thread " + r.ID + " awaits " + whom + title
	case inbox.Edited:
		s = "edited draft " + r.ID + " (" + strings.ReplaceAll(r.To, ",", ", ") + ")" + title
	case inbox.Planned:
		wrote := ""
		if r.To != "" {
			wrote = ", wrote " + strings.ReplaceAll(r.To, ",", ", ")
		}
		s = "planned " + r.ID + wrote + title
	case inbox.ReadyEmpty:
		s = "ready-empty: the ready column holds " + r.To + ", below the in-progress room"
	case reordered:
		s = "reordered: the pull order changed"
	default:
		s = r.Kind + " " + r.ID + title
	}
	return oneLine(s) // one line of FLAI_REASONS
}

// reasonLines are the reasons' lines, in turn.
func reasonLines(rs []orchestrateReason) []string {
	out := make([]string, len(rs))
	for i, r := range rs {
		out[i] = r.line()
	}
	return out
}

// prime is the reason in the form flai prime --reason and MCP prime's
// reasons take (ADR-0217 §1, §6), built from ctxpack's kinds: a story's move
// to review is review:S-nnnn; a story's move to backlog and a story's edit,
// which the trigger takes only for a draft, are draft:S-nnnn; a plan is
// planned:<id>; a thread is thread:TH-nnnn; a release due is release; the
// ready column draining and the order changing are ready. It is "" for a
// reason with no form: asked and schedule, which look at the whole project,
// any other event a project's trigger may name, and a form flai prime would
// refuse, such as a move of an epic, so that the run primes with the general
// pack rather than be refused.
func (r orchestrateReason) prime() string {
	var p ctxpack.Reason
	switch {
	case r.Kind == workitem.Moved && r.Type == workitem.Story && r.To == workitem.Review:
		p = ctxpack.Reason{Kind: ctxpack.ReasonReview, ID: r.ID}
	case r.Kind == workitem.Moved && r.Type == workitem.Story && r.To == workitem.Backlog,
		r.Kind == inbox.Edited && r.Type == workitem.Story:
		p = ctxpack.Reason{Kind: ctxpack.ReasonDraft, ID: r.ID}
	case r.Kind == inbox.Planned:
		p = ctxpack.Reason{Kind: ctxpack.ReasonPlanned, ID: r.ID}
	case r.Kind == inbox.Thread:
		p = ctxpack.Reason{Kind: ctxpack.ReasonThread, ID: r.ID}
	case r.Kind == inbox.ReleaseDue:
		p = ctxpack.Reason{Kind: ctxpack.ReasonRelease}
	case r.Kind == inbox.ReadyEmpty, r.Kind == reordered:
		p = ctxpack.Reason{Kind: ctxpack.ReasonReady}
	default:
		return ""
	}
	if _, err := ctxpack.ParseReasons([]string{p.String()}); err != nil {
		return ""
	}
	return p.String()
}

// primeReasons are the reasons' forms for prime, one for each reason in
// turn, so that the start prompt writes each beside its line, and the
// harness drops a form given twice when it names them as flags. They are
// nil when any reason has no form, since such a run looks at the whole
// project and primes with the general pack (ADR-0217 §6).
func primeReasons(rs []orchestrateReason) []string {
	if len(rs) == 0 {
		return nil
	}
	out := make([]string, len(rs))
	for i, r := range rs {
		if out[i] = r.prime(); out[i] == "" {
			return nil
		}
	}
	return out
}

// reasonsSaid names the reasons for a journal entry or a started event:
// "3 reasons: moved S-0001 to review: <title>; ...", each line cut at
// reasonSaidMax characters, the first reasonsNamed of them, and the rest
// counted.
func reasonsSaid(rs []orchestrateReason) string {
	if len(rs) == 0 {
		return "no reason"
	}
	head := fmt.Sprintf("%d reasons: ", len(rs))
	if len(rs) == 1 {
		head = "1 reason: "
	}
	var said []string
	for _, r := range rs[:min(len(rs), reasonsNamed)] {
		said = append(said, cut(r.line(), reasonSaidMax))
	}
	if more := len(rs) - reasonsNamed; more > 0 {
		said = append(said, fmt.Sprintf("and %d more", more))
	}
	return head + strings.Join(said, "; ")
}

// reasonQueue is what flai serve keeps of the orchestrator for one project
// between looks, and between one flai serve and the next: the reasons that
// wait for a run, and whether the orchestrate action was on at the last look.
type reasonQueue struct {
	// On is whether the orchestrate action was on at the last look; a look
	// that finds it on when this is false finds it turned on.
	On bool `json:"on,omitempty"`
	// First is when the oldest reason queued was queued; empty with none.
	First string `json:"first,omitempty"`
	// Reasons wait for a run, oldest first.
	Reasons []orchestrateReason `json:"reasons,omitempty"`
	// Taken are the reasons the run of Session was started for, kept until
	// it ends: a run that fails gives them back, so that they wait for the
	// next (ADR-0209 §5).
	Taken   []orchestrateReason `json:"taken,omitempty"`
	Session string              `json:"session,omitempty"`
}

// add queues each reason the queue does not hold yet, stamped queued at
// now, and returns how many it queued.
func (q *reasonQueue) add(now time.Time, rs ...orchestrateReason) int {
	n := 0
	for _, r := range rs {
		if slices.ContainsFunc(q.Reasons, func(o orchestrateReason) bool { return o.key() == r.key() }) {
			continue
		}
		r.Queued = now.UTC().Format(time.RFC3339)
		q.Reasons = append(q.Reasons, r)
		n++
	}
	q.first()
	return n
}

// remove takes the reasons a run was started with out of the queue; those
// queued since stay, for the next run.
func (q *reasonQueue) remove(done []orchestrateReason) {
	q.Reasons = slices.DeleteFunc(q.Reasons, func(r orchestrateReason) bool {
		return slices.ContainsFunc(done, func(d orchestrateReason) bool { return d.key() == r.key() })
	})
	q.first()
}

// clear empties the queue, the reasons taken included.
func (q *reasonQueue) clear() { q.Reasons, q.First, q.Taken, q.Session = nil, "", nil, "" }

// take takes the reasons the run of session was started with out of the
// queue and keeps them as taken until it ends; those queued since stay, for
// the next run.
func (q *reasonQueue) take(session string, rs []orchestrateReason) {
	q.remove(rs)
	q.Taken, q.Session = slices.Clone(rs), session
}

// settle lets go of the reasons the run of session took, now that it has
// ended: a run that failed gives them back to the queue, as they were
// queued, so that they are due at once, and settle returns how many; any
// other end drops them. A run that took none, or another session, settles
// nothing.
func (q *reasonQueue) settle(session string, failed bool) int {
	if session == "" || q.Session != session {
		return 0
	}
	var given []orchestrateReason
	if failed {
		for _, r := range q.Taken {
			if !slices.ContainsFunc(q.Reasons, func(o orchestrateReason) bool { return o.key() == r.key() }) {
				given = append(given, r)
			}
		}
		// queued before any waiting now, so first, oldest first
		q.Reasons = append(given, q.Reasons...)
		q.first()
	}
	q.Taken, q.Session = nil, ""
	return len(given)
}

// first sets First from the reasons queued.
func (q *reasonQueue) first() {
	q.First = ""
	for _, r := range q.Reasons {
		if q.First == "" || r.Queued < q.First {
			q.First = r.Queued
		}
	}
}

// due says whether a run should start for the queue at now: a reason
// queued, and either an asked one among them or coalesce passed since the
// first was queued. When one is queued and not due, at is when it will be.
func (q reasonQueue) due(now time.Time, coalesce time.Duration) (ready bool, at time.Time) {
	if len(q.Reasons) == 0 {
		return false, time.Time{}
	}
	if slices.ContainsFunc(q.Reasons, func(r orchestrateReason) bool { return r.Kind == reasonAsked }) {
		return true, now
	}
	first, err := time.Parse(time.RFC3339, q.First)
	if err != nil {
		return true, now
	}
	at = first.Add(coalesce)
	return !now.Before(at), at
}

// askedReason is the operator's reason, saying what they did, at now.
func askedReason(said string, now time.Time) orchestrateReason {
	return orchestrateReason{Kind: reasonAsked, Summary: said, At: now.UTC().Format(time.RFC3339)}
}

// scheduleReason is orchestration.schedule, sched, coming round at now.
func scheduleReason(sched *cron.Schedule, now time.Time) orchestrateReason {
	return orchestrateReason{Kind: reasonSchedule, ID: sched.String(), Summary: orchestrationSchedule + " " + sched.String() + " came round", At: now.UTC().Format(time.RFC3339)}
}

// orchestratorTriggers are the triggers of the orchestrator's role file
// among all, or builtinTriggers when there is none, as when it does not
// parse, or when it still names continuous (ADR-0209 §1).
func orchestratorTriggers(all []roles.Role) []roles.Trigger {
	i := slices.IndexFunc(all, func(r roles.Role) bool { return r.Name == roles.Orchestrator })
	if i < 0 || slices.ContainsFunc(all[i].Trigger, func(t roles.Trigger) bool { return t.Kind == roles.TriggerContinuous }) {
		return builtinTriggers
	}
	return all[i].Trigger
}

// scheduled says whether triggers start the orchestrator on
// orchestration.schedule.
func scheduled(triggers []roles.Trigger) bool {
	return slices.ContainsFunc(triggers, func(t roles.Trigger) bool {
		return t.Kind == roles.TriggerSchedule && t.Schedule == orchestrationSchedule
	})
}

// reasonsOf are the events, read as the orchestrator, that a
// trigger among triggers matches, as reasons, in turn: an event no trigger
// matches is none. A task's changes are none, since the orchestrator works
// no task; nor is an epic's move following its story's, or an item's
// cancellation with its parent, which the story's or the parent's own move
// tells; nor is the replan's forecast edit by flai, which changes no draft
// (ADR-0209 §1). items are the project's, for whether an edited story
// is a draft.
func reasonsOf(events []inbox.Event, triggers []roles.Trigger, items []*workitem.Item) []orchestrateReason {
	byID := itemsByID(items)
	orchestrator := func(awaits string) bool { return awaits == workitem.ActivityOrchestrator }
	var out []orchestrateReason
	for _, e := range events {
		if r, ok := reasonOf(e, triggers, byID, orchestrator); ok {
			out = append(out, r)
		}
	}
	return out
}

// itemsByID are items by their IDs.
func itemsByID(items []*workitem.Item) map[string]*workitem.Item {
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	return byID
}

// reasonOf is e as a reason, and whether a trigger among triggers matches
// it, read by the reader you says a thread's awaits names; reasonsOf says
// which events are never reasons. byID are the project's items, for whether
// an edited story is a draft. A defined role's events are matched here too,
// with its own agents as the reader (S-0418).
func reasonOf(e inbox.Event, triggers []roles.Trigger, byID map[string]*workitem.Item, you func(awaits string) bool) (orchestrateReason, bool) {
	if e.Type == workitem.Task || e.Follows != "" || (e.Kind == workitem.Moved && e.Cause != "") {
		return orchestrateReason{}, false
	}
	if !matched(triggers, e, occurrenceOf(e, byID, you)) {
		return orchestrateReason{}, false
	}
	return orchestrateReason{Kind: e.Kind, ID: e.ID, Type: e.Type, Title: oneLine(e.Title), To: e.To, Summary: oneLine(e.Summary), At: e.At}, true
}

// matched says whether a trigger among triggers matches e, as o, under any
// of its names.
func matched(triggers []roles.Trigger, e inbox.Event, o roles.Occurrence) bool {
	for _, name := range eventNames(e) {
		o.Event = name
		if triggered(triggers, o) {
			return true
		}
	}
	return false
}

// eventNames are the trigger events an inbox event is, in roles' terms: its
// kind for the kinds ADR-0211 names, and, for a move, ADR-0207's name for
// the state it moved to beside it; none for another kind.
func eventNames(e inbox.Event) []string {
	switch e.Kind {
	case workitem.Moved:
		names := []string{workitem.Moved}
		switch {
		case e.To == workitem.Done && e.Type == workitem.Story:
			names = append(names, "accepted")
		case e.To == workitem.Cancelled, e.To == workitem.Ready, e.To == workitem.Review:
			names = append(names, e.To)
		}
		return names
	case inbox.Edited, inbox.Thread, inbox.Planned, inbox.ReadyEmpty, inbox.ReleaseDue, reordered:
		return []string{e.Kind}
	}
	return nil
}

// replanBy is who the replan's forecast edits are by.
const replanBy = "flai"

// occurrenceOf is e in roles' terms, without its event's name: the state a
// move went to, whether an edited item is a story that is a draft after the
// edit, and whom a thread awaits, "you" for the reader, whose names you
// says, or "operator".
func occurrenceOf(e inbox.Event, byID map[string]*workitem.Item, you func(awaits string) bool) roles.Occurrence {
	var o roles.Occurrence
	switch e.Kind {
	case workitem.Moved:
		o.To = e.To
	case inbox.Edited:
		it := byID[e.ID]
		o.Draft = it != nil && it.Type == workitem.Story && it.Draft && (e.By != replanBy || e.To != "forecast")
	case inbox.Thread:
		switch awaits := e.Fields[inbox.FieldAwaits]; {
		case awaits == "operator":
			o.Awaits = "operator"
		case awaits != "" && you(awaits):
			o.Awaits = "you"
		}
	}
	return o
}

// triggered says whether an event trigger among triggers matches o: its
// event is o's, and its condition, when it names one, holds for o. A
// condition flai does not judge against an event matches nothing.
func triggered(triggers []roles.Trigger, o roles.Occurrence) bool {
	return slices.ContainsFunc(triggers, func(t roles.Trigger) bool {
		if t.Kind != roles.TriggerEvent || t.Event != o.Event {
			return false
		}
		if t.When == "" {
			return true
		}
		c, ok := roles.Condition(t.When)
		if !ok {
			return false
		}
		holds, implemented := c.Matches(o)
		return implemented && holds
	})
}

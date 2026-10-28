package inbox

import (
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Thread is the kind, and the type, of the change that a thread's newest entry
// makes it await the reader, or, read in the orchestrator's role, the operator
// (ADR-0211 §1, §5): ID and Title are the thread's, To is whom it awaits as
// Summarize says it, "you" or "operator", By the entry's author, Cause the
// thread's story, or its anchor's path when it is on none, and At the entry's
// time. Fields hold FieldAwaits and FieldOn. Threads have no log: CatchUp
// derives the change from the thread files, as it does messages.
const Thread = "thread"

// The keys of a thread change's Fields.
const (
	FieldAwaits = "awaits" // thread: whom it awaits by name, the reader's agent or "operator"
	FieldOn     = "on"     // thread: the item it is anchored to, or its path when on none
)

// threadEvents are the changes of kind thread that the open threads read as
// now: one for each whose newest entry is someone else's than r's and makes it
// await r, or, with operator, the operator, under Summarize's rule
// (ADR-0181). Whether that entry is news is the caller's to judge by its time.
// An entry stamped after now, by a writer whose clock is ahead of the
// reader's, waits until the reader's clock reaches it: a cursor set at now
// would let it through on every look until then, where it is told once.
func threadEvents(repo *workitem.Repo, r Reader, operator bool, now time.Time) ([]Event, error) {
	if r.Agent == "" {
		r.Agent = "agent"
	}
	all, err := threads.List(repo)
	if err != nil {
		return nil, err
	}
	var out []Event
	for _, th := range all {
		if !th.Open() {
			continue
		}
		sum := Summarize(repo, r, th)
		awaits := ""
		// the reader's own entry is no news to it, whomever the thread awaits
		switch {
		case sum.LastBy == r.Agent:
			continue
		case sum.Awaiting == "you":
			awaits = r.Agent
		case operator && sum.Awaiting == "operator":
			awaits = "operator"
		default:
			continue
		}
		entries := th.Entries() // at least one: a thread with none awaits another
		at := entries[len(entries)-1].At
		if t, err := time.Parse(workitem.TimeFormat, at); err == nil && t.After(now) {
			continue
		}
		on := th.Anchor.Item
		if on == "" {
			on = th.Anchor.Path
		}
		cause := sum.Story
		if cause == "" {
			cause = th.Anchor.Path
		}
		c := workitem.Change{ID: th.ID, Type: Thread, Title: th.Title, Kind: Thread, To: sum.Awaiting, By: sum.LastBy, Cause: cause, At: at}
		out = append(out, Event{Change: c, Summary: describeThread(c), Fields: map[string]string{FieldAwaits: awaits, FieldOn: on}})
	}
	return out, nil
}

// describeThread says a change of kind thread: whom the thread awaits, who
// wrote on what, and what to do with it.
func describeThread(c workitem.Change) string {
	whom := "you"
	if c.To == "operator" {
		whom = "the operator"
	}
	wrote := ""
	if c.By != "" {
		wrote = ": " + c.By + " wrote"
		if c.Cause != "" {
			wrote += " on " + c.Cause
		}
	}
	return c.ID + " " + c.Title + " awaits " + whom + wrote + "; read it with thread_get and answer it with thread_reply"
}

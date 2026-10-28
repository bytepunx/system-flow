package inbox

import (
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// threadChanges are the changes of kind thread that CatchUp reports to r, read
// in the orchestrator's role with operator, at the fixture's clock.
func (f *fixture) threadChanges(t *testing.T, r Reader, operator bool) []Event {
	t.Helper()
	items, err := f.repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	events, _, err := CatchUp(f.repo, r, "", operator, f.clock, items)
	if err != nil {
		t.Fatal(err)
	}
	var out []Event
	for _, e := range events {
		if e.Kind == Thread {
			out = append(out, e)
		}
	}
	return out
}

// openAt opens a thread on on by author at at.
func (f *fixture) openAt(t *testing.T, title, on, author string, at time.Time) *threads.Thread {
	t.Helper()
	th, err := threads.New(f.repo, threads.NewOptions{Title: title, On: on, Author: author, Text: "Which way?", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// replyAt adds a plain entry by author to the thread id at at.
func (f *fixture) replyAt(t *testing.T, id, author string, at time.Time) {
	t.Helper()
	if _, err := threads.Reply(f.repo, id, author, "This way.", at); err != nil {
		t.Fatal(err)
	}
}

// ADR-0211 §5: a reply on a thread the reader takes part in is one thread
// event, naming the thread, its story, whom it awaits, and the reply's author
// and time, in words that say what to do; a second read after the cursor
// moved reports it no more.
func TestAReplyOnAThreadTheReaderTakesPartInIsOneThreadEvent(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	claude := Reader{Agent: "claude", Own: f.story.ID}
	asked := f.openAt(t, "May I?", other.ID, "claude", t0)
	f.clock = t0.Add(10 * time.Minute)
	if got := f.threadChanges(t, claude, false); len(got) != 0 {
		t.Fatalf("claude's own question is no news to it: %+v", got)
	}

	at := t0.Add(15 * time.Minute)
	f.replyAt(t, asked.ID, "alex", at)
	f.clock = t0.Add(20 * time.Minute)
	got := f.threadChanges(t, claude, false)
	if len(got) != 1 {
		t.Fatalf("one thread event for the reply: %+v", got)
	}
	e := got[0]
	if e.ID != asked.ID || e.Type != Thread || e.Title != "May I?" || e.To != "you" || e.By != "alex" || e.Cause != other.ID || e.At != at.Format(workitem.TimeFormat) {
		t.Errorf("the event names the thread, its story, whom it awaits, and the reply: %+v", e.Change)
	}
	if e.Fields[FieldAwaits] != "claude" || e.Fields[FieldOn] != other.ID {
		t.Errorf("fields say whom it awaits by name and what it is on: %v", e.Fields)
	}
	want := asked.ID + " May I? awaits you: alex wrote on " + other.ID + "; read it with thread_get and answer it with thread_reply"
	if e.Summary != want || Describe(e.Change) != want {
		t.Errorf("summary:\n%s\nwant\n%s\nDescribe: %s", e.Summary, want, Describe(e.Change))
	}
	if again := f.threadChanges(t, claude, false); len(again) != 0 {
		t.Errorf("a second read reports it no more: %+v", again)
	}
	f.clock = t0.Add(30 * time.Minute)
	if later := f.threadChanges(t, claude, false); len(later) != 0 {
		t.Errorf("nor does one after the cursor moved on: %+v", later)
	}
}

// ADR-0181, ADR-0211 §5: a new thread on the reader's own story awaits it, in
// which it takes no part yet; another reader with no story of its own hears
// nothing of it.
func TestANewThreadOnTheReadersOwnStoryIsOneThreadEvent(t *testing.T) {
	f := setup(t)
	claude := Reader{Agent: "claude", Own: f.story.ID}
	f.threadChanges(t, claude, false) // the first look sets the cursor
	f.threadChanges(t, Reader{Agent: "codex"}, false)

	f.clock = t0.Add(10 * time.Minute)
	th := f.openAt(t, "Which way here?", f.story.ID, "alex", t0.Add(5*time.Minute))
	got := f.threadChanges(t, claude, false)
	if len(got) != 1 || got[0].ID != th.ID || got[0].To != "you" || got[0].By != "alex" || got[0].Cause != f.story.ID || got[0].Fields[FieldAwaits] != "claude" {
		t.Fatalf("one thread event for the new thread on claude's story: %+v", got)
	}
	if none := f.threadChanges(t, Reader{Agent: "codex"}, false); len(none) != 0 {
		t.Errorf("a reader with no part in it hears nothing: %+v", none)
	}
}

// ADR-0211 §5: read in the orchestrator's role, a planner's question that
// awaits the operator is one thread event saying so; read outside it, it is
// none; and the orchestrator's own entry on it is no news to it.
func TestAPlannersQuestionAwaitsTheOperatorInTheOrchestratorsRole(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Owner = "alex"
	other := f.storyIn(t, "Other", []string{"docs/other"})
	orch := Reader{Agent: "orchestrator"}
	f.threadChanges(t, orch, true) // the first look sets the cursor
	f.threadChanges(t, Reader{Agent: "claude"}, false)

	f.clock = t0.Add(10 * time.Minute)
	planner := "planner-" + other.ID
	th := f.openAt(t, "Split it?", other.ID, planner, t0.Add(5*time.Minute))
	got := f.threadChanges(t, orch, true)
	if len(got) != 1 || got[0].ID != th.ID || got[0].To != "operator" || got[0].By != planner || got[0].Fields[FieldAwaits] != "operator" || got[0].Fields[FieldOn] != other.ID {
		t.Fatalf("one thread event awaiting the operator: %+v", got)
	}
	want := th.ID + " Split it? awaits the operator: " + planner + " wrote on " + other.ID + "; read it with thread_get and answer it with thread_reply"
	if got[0].Summary != want {
		t.Errorf("summary:\n%s\nwant\n%s", got[0].Summary, want)
	}
	if none := f.threadChanges(t, Reader{Agent: "claude"}, false); len(none) != 0 {
		t.Errorf("outside the orchestrator's role the operator's thread is none of the reader's: %+v", none)
	}

	f.clock = t0.Add(20 * time.Minute)
	f.replyAt(t, th.ID, "orchestrator", t0.Add(15*time.Minute))
	if Summarize(f.repo, orch, mustThread(t, f, th.ID)).Awaiting != "operator" {
		t.Fatal("the orchestrator's escalation still awaits the operator")
	}
	if own := f.threadChanges(t, orch, true); len(own) != 0 {
		t.Errorf("the orchestrator's own entry is no news to it: %+v", own)
	}
}

// ADR-0211 §5: no thread event for the reader's own entry, a resolved
// thread, or a thread that awaits another agent.
func TestNoThreadEventForTheReadersOwnEntryAResolvedThreadOrAnothers(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Owner = "alex"
	other := f.storyIn(t, "Other", []string{"docs/other"})
	claude := Reader{Agent: "claude", Own: f.story.ID}
	f.threadChanges(t, claude, false) // the first look sets the cursor

	f.clock = t0.Add(20 * time.Minute)
	answered := f.openAt(t, "Answered", f.story.ID, "alex", t0.Add(11*time.Minute))
	f.replyAt(t, answered.ID, "claude", t0.Add(12*time.Minute))
	resolved := f.openAt(t, "Resolved", f.story.ID, "alex", t0.Add(13*time.Minute))
	if _, err := threads.Resolve(f.repo, resolved.ID, "alex", "settled", t0.Add(14*time.Minute)); err != nil {
		t.Fatal(err)
	}
	another := f.openAt(t, "Another's", other.ID, "alex", t0.Add(15*time.Minute))
	f.replyAt(t, f.openAt(t, "Codex's", other.ID, "codex", t0.Add(16*time.Minute)).ID, "alex", t0.Add(17*time.Minute))
	if got := Summarize(f.repo, claude, mustThread(t, f, another.ID)).Awaiting; got != "other" {
		t.Fatalf("the thread on another story awaits another: %s", got)
	}

	if got := f.threadChanges(t, claude, false); len(got) != 0 {
		t.Errorf("no thread event for claude's own entry, a resolved thread, or another's: %+v", got)
	}
}

// An entry stamped after the reader's now, by a writer whose clock is ahead,
// is told once, when the reader's clock reaches it, not on every look before.
func TestAThreadEntryAheadOfTheReadersClockIsToldOnceItsTimeComes(t *testing.T) {
	f := setup(t)
	claude := Reader{Agent: "claude", Own: f.story.ID}
	f.clock = t0.Add(10 * time.Minute)
	f.threadChanges(t, claude, false) // the first look sets the cursor
	th := f.openAt(t, "Ahead", f.story.ID, "alex", t0.Add(12*time.Minute))
	for range 2 {
		if early := f.threadChanges(t, claude, false); len(early) != 0 {
			t.Fatalf("not yet: %+v", early)
		}
	}
	f.clock = t0.Add(13 * time.Minute)
	if got := f.threadChanges(t, claude, false); len(got) != 1 || got[0].ID != th.ID {
		t.Fatalf("told once its time has come: %+v", got)
	}
	if again := f.threadChanges(t, claude, false); len(again) != 0 {
		t.Errorf("and not again: %+v", again)
	}
}

// mustThread reads the thread id again.
func mustThread(t *testing.T, f *fixture, id string) *threads.Thread {
	t.Helper()
	th, err := threads.Get(f.repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

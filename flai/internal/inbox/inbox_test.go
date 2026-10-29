package inbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC)

type fixture struct {
	repo  *workitem.Repo
	epic  *workitem.Item
	story *workitem.Item // in progress, touching flai/internal/inbox
	clock time.Time
}

func setup(t *testing.T) *fixture {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"design/system", "docs/users", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	f := &fixture{repo: repo, clock: t0.Add(time.Minute)}
	if f.epic, err = repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: t0}); err != nil {
		t.Fatal(err)
	}
	f.story = f.storyIn(t, "Story", []string{"flai/internal/inbox"}, workitem.Ready, workitem.InProgress)
	return f
}

// storyIn creates a story with one criterion and moves it through states.
func (f *fixture) storyIn(t *testing.T, title string, touches []string, states ...string) *workitem.Item {
	t.Helper()
	s, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: f.epic.ID, Owner: "alex", Touches: touches, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(s.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	if s, err = f.repo.Get(s.ID); err != nil {
		t.Fatal(err)
	}
	for _, st := range states {
		if _, err := f.repo.Transition(s, st, "alex", "", t0); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func (f *fixture) read(t *testing.T, opt Options) Inbox {
	t.Helper()
	opt.Repo = f.repo
	opt.Now = func() time.Time { return f.clock }
	out, err := Read(context.Background(), opt)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func threadIDs(in Inbox) []string {
	var ids []string
	for _, th := range in.Threads {
		ids = append(ids, th.ID)
	}
	return ids
}

// newThread opens a thread on on by author.
func (f *fixture) newThread(t *testing.T, title, on, author, text string) *threads.Thread {
	t.Helper()
	th, err := threads.New(f.repo, threads.NewOptions{Title: title, On: on, Author: author, Text: text, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// reply adds an entry by author to the thread id, a recommendation when rec.
func (f *fixture) reply(t *testing.T, id, author, text string, rec bool) *threads.Thread {
	t.Helper()
	th, err := threads.ReplyWith(f.repo, id, author, text, t0.Add(time.Minute), threads.Marks{Recommendation: rec})
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// ADR-0181: a thread on the agent's own story whose last entry is someone
// else's awaits the agent; its own question there awaits the operator, listed
// with Operator or All; a thread on another story that it takes no part in
// awaits it not at all, listed only with All; Story keeps the threads on that
// story.
func TestThreadsAwaitingTheAgent(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	asked := f.newThread(t, "Which way?", f.story.ID, "alex", "A or B?")
	mine := f.newThread(t, "May I?", f.story.ID, "claude", "Shall I?")
	elsewhere := f.newThread(t, "There?", other.ID, "alex", "And there?")

	in := f.read(t, Options{Agent: "claude", Own: f.story.ID})
	if in.Agent != "claude" || in.AwaitingYou != 1 || strings.Join(threadIDs(in), " ") != asked.ID {
		t.Errorf("awaiting claude: %d %v, want %s", in.AwaitingYou, threadIDs(in), asked.ID)
	}
	if in.Threads[0].Awaiting != "you" || in.Threads[0].LastBy != "alex" || in.Threads[0].Story != f.story.ID {
		t.Errorf("summary: %+v", in.Threads[0])
	}
	op := f.read(t, Options{Agent: "claude", Own: f.story.ID, Operator: true})
	if strings.Join(threadIDs(op), " ") != asked.ID+" "+mine.ID || op.AwaitingYou != 1 || op.Threads[1].Awaiting != "operator" {
		t.Errorf("operator lists claude's own question too, awaiting the operator, without all: %v, awaiting %d", threadIDs(op), op.AwaitingYou)
	}
	all := f.read(t, Options{Agent: "claude", Own: f.story.ID, All: true})
	if got := threadIDs(all); len(got) != 3 || all.AwaitingYou != 1 || all.Threads[2].Awaiting != "other" {
		t.Errorf("all lists the thread on another story too, awaiting another: %v, awaiting %d", got, all.AwaitingYou)
	}
	if none := f.read(t, Options{Agent: "claude"}); len(none.Threads) != 0 || none.AwaitingYou != 0 {
		t.Errorf("with no story of its own, claude takes part only in its question, which awaits the operator: %v", threadIDs(none))
	}
	if got := threadIDs(f.read(t, Options{Agent: "claude", Story: other.ID, Own: other.ID})); len(got) != 1 || got[0] != elsewhere.ID {
		t.Errorf("story keeps its own threads: %v", got)
	}
}

// ADR-0181: awaiting is "you" for a participant, on another story too, and on
// the agent's own story or a task of it; "operator" for an opener whose
// recommendation awaits the operator, and for a thread whose last entry is by
// none of the operator's names, the project's owner, the story's owner, or
// "designer" when neither is named; else "other".
func TestSummarizeAwaitingByTheRule(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Owner = "alex"
	bobs, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Bob's", Parent: f.epic.ID, Owner: "bob", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "A task", Parent: f.story.ID, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	claude := Reader{Agent: "claude", Own: f.story.ID}
	awaiting := func(r Reader, th *threads.Thread) string {
		t.Helper()
		return Summarize(f.repo, r, th).Awaiting
	}

	// a participant on another story's thread is awaited when another wrote last
	there := f.newThread(t, "There?", bobs.ID, "claude", "May I?")
	if got := awaiting(claude, there); got != "operator" {
		t.Errorf("claude's own question on another story awaits the operator: %s", got)
	}
	there = f.reply(t, there.ID, "bob", "Yes.", false)
	if got := awaiting(claude, there); got != "you" {
		t.Errorf("a participant is awaited on another story's thread: %s", got)
	}
	if got := awaiting(Reader{Agent: "codex"}, there); got != "other" {
		t.Errorf("the story owner's name counts as the operator's: bob wrote last, so it awaits another: %s", got)
	}
	carol := f.reply(t, f.newThread(t, "Carol's", bobs.ID, "planner", "Q?").ID, "carol", "Hm.", false)
	if got := awaiting(Reader{Agent: "codex"}, carol); got != "operator" {
		t.Errorf("a thread whose last entry is by none of the operator's names awaits the operator: %s", got)
	}

	// the agent's own story and a task of it, in any zero padding
	onTask := f.newThread(t, "On the task", task.ID, "alex", "And this?")
	if got := awaiting(Reader{Agent: "claude", Own: strings.Replace(f.story.ID, "S-", "S-0", 1)}, onTask); got != "you" {
		t.Errorf("a thread on a task of the agent's own story awaits it: %s", got)
	}
	if got := awaiting(Reader{Agent: "claude"}, onTask); got != "other" {
		t.Errorf("without its own story claude takes no part in it, and alex wrote last: %s", got)
	}

	// ADR-0090: the opener's question with a recommendation pending awaits
	// the operator, who confirms it, and so does it for anyone else
	rec := f.newThread(t, "Which?", f.story.ID, "claude", "A or B?")
	rec = f.reply(t, rec.ID, "codex", "A, because.", true)
	if Summarize(f.repo, claude, rec).PendingRecommendation == nil {
		t.Fatal("the recommendation is pending")
	}
	if got := awaiting(claude, rec); got != "operator" {
		t.Errorf("an opener with a pending recommendation reads operator: %s", got)
	}
	if got := awaiting(Reader{Agent: "alex"}, rec); got != "operator" {
		t.Errorf("a pending recommendation awaits the operator: %s", got)
	}

	// with no owner named anywhere the operator is "designer"
	f.repo.Manifest.Owner = ""
	onPath := f.newThread(t, "On a path", "design/system", "designer", "Look here.")
	if got := awaiting(claude, onPath); got != "other" {
		t.Errorf("with no owner named, designer is the operator: %s", got)
	}
	onPath = f.reply(t, onPath.ID, "planner", "Looked.", false)
	if got := awaiting(claude, onPath); got != "operator" {
		t.Errorf("with no owner named, a planner's last entry awaits the operator: %s", got)
	}
	if got := awaiting(Reader{Agent: "planner"}, onPath); got != "operator" {
		t.Errorf("the planner wrote last, so it awaits the operator, not the planner: %s", got)
	}
}

// ADR-0181 reproduces the measurement of 2026-10-20: a new story's agent read
// 31 threads awaiting it, none on its story, all planners' questions, the
// orchestrator's escalations, and flai serve's notices awaiting the operator.
// Its inbox, story-scoped or not, lists none and counts none, while the
// orchestrator's role lists them all, and a thread on the agent's story or a
// task of it whose last entry is the operator's awaits the agent.
func TestTheOperatorsThreadsDoNotAwaitAStorysAgent(t *testing.T) {
	f := setup(t)
	f.repo.Manifest.Owner = "alex"
	other := f.storyIn(t, "Other", []string{"docs/other"})
	agent := "agent-" + f.story.ID
	authors := []string{"planner-" + other.ID, "orchestrator", "flai"}
	for i := range 31 {
		on := []string{other.ID, f.epic.ID, "design/system"}[i%3]
		th := f.newThread(t, fmt.Sprintf("Question %d", i), on, authors[i%3], "Operator, which way?")
		if i%5 == 0 {
			f.reply(t, th.ID, "alex", "Tell me more.", false)
			f.reply(t, th.ID, authors[i%3], "Here is more.", false)
		}
	}

	scoped := Options{Agent: agent, Story: f.story.ID, Own: f.story.ID}
	unscoped := Options{Agent: agent, Own: f.story.ID}
	for name, opt := range map[string]Options{"story-scoped": scoped, "unscoped": unscoped} {
		in := f.read(t, opt)
		if len(in.Threads) != 0 || in.AwaitingYou != 0 || in.ThreadsOmitted != 0 {
			t.Errorf("the %s inbox lists and counts none of the operator's threads: %d listed, awaiting %d, omitted %d", name, len(in.Threads), in.AwaitingYou, in.ThreadsOmitted)
		}
	}
	orch := f.read(t, Options{Agent: "orchestrator", Operator: true})
	if len(orch.Threads) != 31 {
		t.Fatalf("the orchestrator's role lists the 31 threads awaiting the operator: %d", len(orch.Threads))
	}
	for _, th := range orch.Threads {
		if th.Awaiting != "operator" {
			t.Errorf("%s awaits the operator, as the orchestrator reads it: %s", th.ID, th.Awaiting)
		}
	}

	task, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "A task", Parent: f.story.ID, Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	onStory := f.newThread(t, "Which way here?", f.story.ID, "alex", "A or B?")
	onTask := f.reply(t, f.newThread(t, "And on the task?", task.ID, "planner-"+f.story.ID, "Split it?").ID, "alex", "Yes, split it.", false)
	for name, opt := range map[string]Options{"story-scoped": scoped, "unscoped": unscoped} {
		in := f.read(t, opt)
		if strings.Join(threadIDs(in), " ") != onStory.ID+" "+onTask.ID || in.AwaitingYou != 2 || in.ThreadsOmitted != 0 {
			t.Errorf("the %s inbox lists the story's two threads, awaiting the agent: %v, awaiting %d, omitted %d", name, threadIDs(in), in.AwaitingYou, in.ThreadsOmitted)
			continue
		}
		for _, th := range in.Threads {
			if th.Awaiting != "you" || th.LastBy != "alex" || th.Story != f.story.ID {
				t.Errorf("%s awaits the agent on its story: %+v", th.ID, th)
			}
		}
	}
}

// Ready stories come in the board's pull order, a held one keeping its place
// with why it is held.
func TestReadyStoriesInPullOrderWithAHeldOne(t *testing.T) {
	f := setup(t)
	held := f.storyIn(t, "Overlaps", []string{"flai/internal/inbox/inbox.go"}, workitem.Ready)
	clear := f.storyIn(t, "Clear", []string{"docs/clear"}, workitem.Ready)
	board, err := f.repo.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	board.Order = []string{clear.ID, held.ID}
	if err := board.Save("2026-09-18"); err != nil {
		t.Fatal(err)
	}

	in := f.read(t, Options{Agent: "claude"})
	if len(in.Ready) != 2 || in.Ready[0].ID != clear.ID || in.Ready[1].ID != held.ID {
		t.Fatalf("ready in pull order: %+v", in.Ready)
	}
	if in.Ready[0].Held != nil || in.Ready[1].Held == nil || !strings.Contains(in.Ready[1].Held.Reason, f.story.ID) {
		t.Errorf("the overlapping story is held by %s: %+v, %+v", f.story.ID, in.Ready[0].Held, in.Ready[1].Held)
	}
	if !in.CanPull || in.PullHold != "" {
		t.Errorf("no limit is full: can_pull %v, pull_hold %q", in.CanPull, in.PullHold)
	}
}

// A change is reported once: a second reader of the same agent's inbox, here
// the changes alone as wait_for_events reads them, does not report it again,
// while another agent's cursor is its own.
func TestAChangeIsReportedOnceAcrossCalls(t *testing.T) {
	f := setup(t)
	f.read(t, Options{Agent: "claude"}) // the first look sets the cursor

	f.clock = t0.Add(10 * time.Minute)
	moved := f.storyIn(t, "Moved", []string{"docs/moved"})
	if _, err := f.repo.Transition(moved, workitem.Ready, "alex", "", t0.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	first := f.read(t, Options{Agent: "claude"})
	if len(first.Changes) != 1 || first.Changes[0].ID != moved.ID || first.Changes[0].Summary != moved.ID+" Moved moved to ready by alex" {
		t.Fatalf("the move is reported: %+v", first.Changes)
	}

	items, err := f.repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	again, omitted, err := CatchUp(f.repo, Reader{Agent: "claude"}, "", false, f.clock, items)
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != 0 || omitted != 0 {
		t.Errorf("a change already reported is not reported again: %+v, %d omitted", again, omitted)
	}
	if other := f.read(t, Options{Agent: "codex"}); len(other.Changes) == 0 {
		t.Error("another agent's cursor is its own: it hears of the move")
	}
}

// send starts a conversation from one story to another, by from's agent.
func (f *fixture) send(t *testing.T, from, to *workitem.Item, text string, at time.Time) *messages.Conversation {
	t.Helper()
	c, err := messages.Send(f.repo, messages.SendOptions{From: from.ID, To: to.ID, Author: "agent-" + from.ID, Text: text, About: []string{"design/system"}, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// S-0331: the inbox lists the open conversations of the agent's own story,
// each with the other story, its about, its last entry, and whether it awaits
// this story; a closed one, and one between two other stories, are left out;
// awaiting_you counts threads only; and with no story there are none.
func TestTheInboxListsTheOpenConversationsOfTheAgentsStory(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"}, workitem.Ready, workitem.InProgress)
	third := f.storyIn(t, "Third", []string{"docs/third"}, workitem.Ready, workitem.InProgress)
	toMe := f.send(t, other, f.story, "Will you leave plan.md to me?", t0)
	fromMe := f.send(t, f.story, third, "I am changing the inbox.", t0)
	f.send(t, other, third, "Not ours.", t0)
	closed := f.send(t, f.story, other, "Done with this.", t0)
	if _, err := messages.Close(f.repo, closed.ID, "alex", "settled", t0); err != nil {
		t.Fatal(err)
	}

	in := f.read(t, Options{Agent: "claude", Own: f.story.ID})
	if len(in.Messages) != 2 || in.Messages[0].ID != toMe.ID || in.Messages[1].ID != fromMe.ID {
		t.Fatalf("the story's two open conversations, in ID order: %+v", in.Messages)
	}
	mine, sent := in.Messages[0], in.Messages[1]
	if mine.Awaiting != "you" || mine.With != other.ID || len(mine.About) != 1 || mine.About[0] != "design/system" || mine.Entries != 1 {
		t.Errorf("a message to the story awaits it, with the other story and about: %+v", mine)
	}
	if mine.Last.Author != "agent-"+other.ID || mine.Last.Story != other.ID || mine.Last.Text != "Will you leave plan.md to me?" || mine.Last.At != t0.Format(workitem.TimeFormat) {
		t.Errorf("the last entry: %+v", mine.Last)
	}
	if sent.Awaiting != "other" || sent.With != third.ID {
		t.Errorf("a message the story sent awaits the other: %+v", sent)
	}
	if in.AwaitingYou != 0 || len(in.Threads) != 0 {
		t.Errorf("awaiting_you counts threads only: %d %v", in.AwaitingYou, in.Threads)
	}
	if none := f.read(t, Options{Agent: "claude"}); none.Messages != nil {
		t.Errorf("an agent with no story of its own has no messages: %+v", none.Messages)
	}
}

// S-0331: a message to the agent's story is a change of kind message, told
// once per conversation and look, naming the conversation and the sender's
// story; the story's own entries are not told back, nor are any without a
// story.
func TestAMessageToTheStoryIsToldOnce(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"}, workitem.Ready, workitem.InProgress)
	catchUp := func(story string) []Event {
		t.Helper()
		items, err := f.repo.List(true)
		if err != nil {
			t.Fatal(err)
		}
		events, _, err := CatchUp(f.repo, Reader{Agent: "claude", Own: story}, story, false, f.clock, items)
		if err != nil {
			t.Fatal(err)
		}
		var out []Event
		for _, e := range events {
			if e.Kind == Message {
				out = append(out, e)
			}
		}
		return out
	}
	catchUp(f.story.ID) // the first look sets the cursor

	f.clock = t0.Add(10 * time.Minute)
	c := f.send(t, other, f.story, "Will you leave plan.md to me?", t0.Add(5*time.Minute))
	got := catchUp(f.story.ID)
	if len(got) != 1 || got[0].ID != c.ID || got[0].Type != Message || got[0].Cause != other.ID || got[0].By != "agent-"+other.ID || got[0].To != "design/system" || got[0].Title != c.Title {
		t.Fatalf("one message event naming the conversation and the sender's story: %+v", got)
	}
	if s := got[0].Summary; !strings.Contains(s, other.ID+"'s agent, agent-"+other.ID+", wrote in "+c.ID) || !strings.Contains(s, "about design/system") || !strings.Contains(s, "message_reply") {
		t.Errorf("summary: %s", s)
	}
	if again := catchUp(f.story.ID); len(again) != 0 {
		t.Errorf("told once: %+v", again)
	}

	f.clock = t0.Add(20 * time.Minute)
	if _, err := messages.Reply(f.repo, c.ID, f.story.ID, "claude", "Yes.", t0.Add(15*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if mine := catchUp(f.story.ID); len(mine) != 0 {
		t.Errorf("the story's own reply is not told back: %+v", mine)
	}

	f.clock = t0.Add(30 * time.Minute)
	if _, err := messages.Reply(f.repo, c.ID, other.ID, "agent-"+other.ID, "Thanks.", t0.Add(25*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if none := catchUp(""); len(none) != 0 {
		t.Errorf("no story, no messages: %+v", none)
	}
	f.clock = t0.Add(40 * time.Minute)
	if _, err := messages.Reply(f.repo, c.ID, other.ID, "agent-"+other.ID, "And one more.", t0.Add(35*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if reply := catchUp(f.story.ID); len(reply) != 1 || reply[0].At != t0.Add(35*time.Minute).Format(workitem.TimeFormat) {
		t.Errorf("a reply is told too, as the newest entry: %+v", reply)
	}
}

// S-0372, ADR-0154: a conversation whose newest entry from a story is settled
// awaits neither story, so its summary's awaiting is empty for both, and the
// settled entry is told as a message whose summary says no reply is owed.
func TestASettledConversationAwaitsNeitherStory(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"}, workitem.Ready, workitem.InProgress)
	messageChanges := func() []Event {
		t.Helper()
		items, err := f.repo.List(true)
		if err != nil {
			t.Fatal(err)
		}
		events, _, err := CatchUp(f.repo, Reader{Agent: "claude", Own: f.story.ID}, f.story.ID, false, f.clock, items)
		if err != nil {
			t.Fatal(err)
		}
		var out []Event
		for _, e := range events {
			if e.Kind == Message {
				out = append(out, e)
			}
		}
		return out
	}
	messageChanges() // the first look sets the cursor

	f.clock = t0.Add(10 * time.Minute)
	c := f.send(t, f.story, other, "I changed the inbox.", t0.Add(time.Minute))
	if _, err := messages.ReplyWith(f.repo, c.ID, other.ID, "agent-"+other.ID, "Agreed; nothing breaks here.", t0.Add(5*time.Minute), messages.Marks{Settled: true}); err != nil {
		t.Fatal(err)
	}
	got := messageChanges()
	want := other.ID + "'s agent, agent-" + other.ID + ", wrote in " + c.ID + " " + c.Title + " about design/system, settled: read it with message_get; no reply is owed"
	if len(got) != 1 || got[0].ID != c.ID || got[0].Summary != want {
		t.Fatalf("one message event saying no reply is owed:\n%+v\nwant summary %q", got, want)
	}

	for _, as := range []struct{ agent, story string }{{"claude", f.story.ID}, {"agent-" + other.ID, other.ID}} {
		in := f.read(t, Options{Agent: as.agent, Own: as.story})
		if len(in.Messages) != 1 || in.Messages[0].Awaiting != "" || !in.Messages[0].Last.Settled {
			t.Errorf("a settled conversation awaits neither story, as %s's agent sees it: %+v", as.story, in.Messages)
		}
	}
}

// jsonSize is how many bytes v encodes to as JSON.
func jsonSize(t *testing.T, v any) int {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return len(b)
}

// ADR-0179: a story-scoped read lists the threads on the story and counts,
// not lists, those on other items that it would otherwise have listed;
// awaiting_you counts the listed ones. A thread elsewhere that does not await
// the agent is not counted (ADR-0181). Without a story the same repository
// lists both, and threads_omitted is zero and absent from the JSON.
func TestAStoryScopedReadCountsTheThreadsOnOtherItems(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	own := f.newThread(t, "Which way?", f.story.ID, "alex", "A or B?")
	f.newThread(t, "There?", other.ID, "alex", "And there?")
	answered := f.reply(t, f.newThread(t, "Mine there", other.ID, "claude", "Shall I?").ID, "alex", "Yes.", false)
	f.newThread(t, "Asked there", other.ID, "claude", "May I?")

	in := f.read(t, Options{Agent: "claude", Story: f.story.ID, Own: f.story.ID})
	if got := threadIDs(in); len(got) != 1 || got[0] != own.ID || in.AwaitingYou != 1 || in.ThreadsOmitted != 1 {
		t.Errorf("the story's thread listed, the one awaiting claude elsewhere counted: %v, awaiting %d, omitted %d", got, in.AwaitingYou, in.ThreadsOmitted)
	}
	if op := f.read(t, Options{Agent: "claude", Story: f.story.ID, Own: f.story.ID, Operator: true}); op.ThreadsOmitted != 2 {
		t.Errorf("with operator, claude's question awaiting the operator elsewhere is counted too: %d", op.ThreadsOmitted)
	}
	if all := f.read(t, Options{Agent: "claude", Story: f.story.ID, Own: f.story.ID, All: true}); all.ThreadsOmitted != 3 {
		t.Errorf("with all, every open thread elsewhere is counted: %d", all.ThreadsOmitted)
	}

	whole := f.read(t, Options{Agent: "claude", Own: f.story.ID})
	if got := threadIDs(whole); strings.Join(got, " ") != own.ID+" "+answered.ID || whole.ThreadsOmitted != 0 {
		t.Errorf("without a story both are listed and none counted: %v, omitted %d", got, whole.ThreadsOmitted)
	}
	b, err := json.Marshal(whole)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "threads_omitted") || strings.Contains(string(b), "entries_shortened") {
		t.Errorf("the standalone inbox's JSON is as before: %s", b)
	}
}

// ADR-0179, ADR-0211 §5: a story's inbox reports as a change only a thread on
// that story coming to await the agent; one on another item is counted among
// the threads, not listed as a change either.
func TestAStoryScopedReadReportsOnlyItsStorysThreadChanges(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	f.read(t, Options{Agent: "claude", Own: f.story.ID}) // a cursor
	own := f.openAt(t, "Which way?", f.story.ID, "alex", t0.Add(time.Minute))
	there := f.openAt(t, "Shall I?", other.ID, "claude", t0.Add(time.Minute))
	f.replyAt(t, there.ID, "alex", t0.Add(2*time.Minute))
	f.clock = t0.Add(3 * time.Minute)

	in := f.read(t, Options{Agent: "claude", Story: f.story.ID, Own: f.story.ID})
	var got []string
	for _, e := range in.Changes {
		if e.Kind == Thread {
			got = append(got, e.ID)
		}
	}
	if len(got) != 1 || got[0] != own.ID || in.ThreadsOmitted != 1 {
		t.Errorf("only the story's thread is a change, the other counted: %v, omitted %d", got, in.ThreadsOmitted)
	}
}

// I-0167: an inbox with 33 threads awaiting the agent on other items, none on
// the story, each with a long last entry, and 50 changes, was over Claude
// Code's tool-output limit. Read for the story, it lists none of the 33 and
// counts them, and fitted it is at most ctxpack.PartLimit bytes of JSON. The
// agent takes part in each, as ADR-0181 has a thread elsewhere await it.
func TestTheStorysInboxWithManyThreadsElsewhereFitsThePartLimit(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	long := strings.Repeat("The operator's question, at length. ", 42) // about 1,500 bytes
	for i := range 33 {
		on := other.ID
		if i%2 == 1 {
			on = f.epic.ID
		}
		f.reply(t, f.newThread(t, fmt.Sprintf("Question %d", i), on, "claude", "Which way?").ID, "alex", long, false)
	}
	changes := make([]Event, 0, MaxEvents)
	for i := range MaxEvents {
		c := workitem.Change{ID: fmt.Sprintf("S-%04d", 1000+i), Type: workitem.Story, Title: "A story with a long title", Kind: Edited, To: "body", By: "alex", At: t0.Add(time.Duration(i) * time.Second).Format(workitem.TimeFormat)}
		changes = append(changes, Event{Change: c, Summary: Describe(c) + " " + strings.Repeat("x", 700)})
	}

	whole := f.read(t, Options{Agent: "claude"})
	whole.Changes = changes
	if len(whole.Threads) != 33 || jsonSize(t, whole) <= ctxpack.PartLimit {
		t.Fatalf("the agent's whole inbox reproduces I-0167: %d threads, %d bytes", len(whole.Threads), jsonSize(t, whole))
	}

	in := f.read(t, Options{Agent: "claude", Story: f.story.ID})
	if len(in.Threads) != 0 || in.AwaitingYou != 0 || in.ThreadsOmitted != 33 {
		t.Fatalf("the story's inbox lists none of the 33 and counts them: %d listed, awaiting %d, omitted %d", len(in.Threads), in.AwaitingYou, in.ThreadsOmitted)
	}
	in.Changes = changes
	if !Fit(&in, ctxpack.PartLimit, func() int { return jsonSize(t, &in) }) {
		t.Fatal("the story's inbox fits")
	}
	if n := jsonSize(t, in); n > ctxpack.PartLimit || len(in.Changes) == 0 || in.Omitted != MaxEvents-len(in.Changes) {
		t.Errorf("fitted to %d bytes with %d changes, %d omitted", n, len(in.Changes), in.Omitted)
	}
}

// fitFixture is an inbox of ten changes of about 1,000 bytes each and two
// threads, one with a long last entry and a long pending recommendation.
func fitFixture() Inbox {
	in := Inbox{Agent: "claude", Omitted: 3, Ready: []workitem.BoardCard{}}
	for i := range 10 {
		in.Changes = append(in.Changes, Event{Change: workitem.Change{ID: fmt.Sprintf("S-%04d", i), Kind: "moved", At: t0.Add(time.Duration(i) * time.Second).Format(workitem.TimeFormat)}, Summary: strings.Repeat("c", 1000)})
	}
	rec := &threads.Entry{Author: "codex", Text: "x" + strings.Repeat("é", 1500), Recommendation: true}
	in.Threads = []ThreadSummary{
		{ID: "TH-0001", LastEntry: strings.Repeat("e", 3000), PendingRecommendation: rec},
		{ID: "TH-0002", LastEntry: "short"},
	}
	return in
}

// ADR-0179: Fit drops the fewest changes that fit, oldest first, counting
// them, and leaves the threads whole while that is enough.
func TestFitDropsTheOldestChangesFirst(t *testing.T) {
	in := fitFixture()
	want := in
	want.Changes, want.Omitted = in.Changes[4:], in.Omitted+4
	limit := jsonSize(t, want)

	kept := in
	if !Fit(&in, limit, func() int { return jsonSize(t, &in) }) {
		t.Fatal("dropping four changes fits")
	}
	if len(in.Changes) != 6 || in.Changes[0].ID != "S-0004" || in.Omitted != 7 || in.EntriesShortened != 0 {
		t.Errorf("the four oldest dropped and counted: %d changes from %s, omitted %d, shortened %d", len(in.Changes), in.Changes[0].ID, in.Omitted, in.EntriesShortened)
	}
	if in.Threads[0].LastEntry != kept.Threads[0].LastEntry {
		t.Error("the threads are whole while dropping changes is enough")
	}
	if len(kept.Changes) != 10 || kept.Omitted != 3 {
		t.Errorf("a copy of the inbox kept before stays whole: %d changes, %d omitted", len(kept.Changes), kept.Omitted)
	}
}

// ADR-0179: when dropping every change is not enough, Fit shortens each
// listed thread's last entry and pending recommendation at a rune boundary
// with an ellipsis, counting each thread once, and never the caller's
// entries; it says false when even that does not fit.
func TestFitShortensEntriesOnlyWhenTheChangesAreNotEnough(t *testing.T) {
	in := fitFixture()
	none := in
	none.Changes = []Event{}
	limit := jsonSize(t, none) - 100

	kept := in
	rec := in.Threads[0].PendingRecommendation
	recText := rec.Text
	if !Fit(&in, limit, func() int { return jsonSize(t, &in) }) {
		t.Fatal("shortening the entries fits")
	}
	if len(in.Changes) != 0 || in.Omitted != 13 || in.EntriesShortened != 1 {
		t.Errorf("every change dropped, one thread shortened: %d changes, omitted %d, shortened %d", len(in.Changes), in.Omitted, in.EntriesShortened)
	}
	th := in.Threads[0]
	if th.LastEntry != strings.Repeat("e", ShortEntry)+"…" {
		t.Errorf("last entry shortened to %d bytes: %q", len(th.LastEntry), th.LastEntry)
	}
	if s := th.PendingRecommendation.Text; !utf8.ValidString(s) || !strings.HasSuffix(s, "…") || len(s) > ShortEntry+len("…") || !strings.HasPrefix(recText, strings.TrimSuffix(s, "…")) {
		t.Errorf("recommendation shortened at a rune boundary: %d bytes, valid %v", len(s), utf8.ValidString(s))
	}
	if th.PendingRecommendation.Author != "codex" || !th.PendingRecommendation.Recommendation {
		t.Errorf("the recommendation keeps the rest of its entry: %+v", th.PendingRecommendation)
	}
	if in.Threads[1].LastEntry != "short" {
		t.Errorf("a short entry is kept: %q", in.Threads[1].LastEntry)
	}
	if rec.Text != recText || kept.Threads[0].LastEntry != strings.Repeat("e", 3000) || len(kept.Changes) != 10 {
		t.Error("the caller's threads, entries, and changes are not written to")
	}

	tiny := fitFixture()
	if Fit(&tiny, 100, func() int { return jsonSize(t, &tiny) }) {
		t.Error("an inbox that cannot fit 100 bytes says so")
	}
	if tiny.EntriesShortened != 1 || len(tiny.Changes) != 0 {
		t.Errorf("it is cut as far as it goes: %d changes, %d shortened", len(tiny.Changes), tiny.EntriesShortened)
	}
}

// S-0534: the inbox carries the host's budget state that Options.Budget
// reads, read afresh at each call, and says nothing of it, not even the key,
// once none is set or with no reader.
func TestTheInboxCarriesTheHostsBudgetWhileItIsSet(t *testing.T) {
	f := setup(t)
	if in := f.read(t, Options{Agent: "claude"}); in.Budget != nil {
		t.Errorf("no reader, no budget: %+v", in.Budget)
	}
	set := []Budget{{Provider: "anthropic", Limit: "usage_limit_reached", Status: 429, Since: "2026-09-18T09:00:00Z", Until: "2026-09-18T14:00:00Z"}}
	state := set
	read := func() []Budget { return state }
	in := f.read(t, Options{Agent: "claude", Budget: read})
	if len(in.Budget) != 1 || in.Budget[0] != set[0] {
		t.Fatalf("the state set: %+v", in.Budget)
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if want := `"budget":[{"provider":"anthropic","limit":"usage_limit_reached","status":429,"since":"2026-09-18T09:00:00Z","until":"2026-09-18T14:00:00Z"}]`; !strings.Contains(string(data), want) {
		t.Errorf("the JSON shape, %s, in %s", want, data)
	}
	state = nil
	in = f.read(t, Options{Agent: "claude", Budget: read})
	if data, err = json.Marshal(in); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"budget"`) {
		t.Errorf("a cleared state leaves budget out: %s", data)
	}
}

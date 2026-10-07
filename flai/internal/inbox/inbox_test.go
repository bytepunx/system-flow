package inbox

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

// A thread whose last entry is someone else's awaits the agent; its own
// question awaits someone else and is listed only with All; Story keeps the
// threads on that story.
func TestThreadsAwaitingTheAgent(t *testing.T) {
	f := setup(t)
	other := f.storyIn(t, "Other", []string{"docs/other"})
	asked, err := threads.New(f.repo, threads.NewOptions{Title: "Which way?", On: f.story.ID, Author: "alex", Text: "A or B?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	mine, err := threads.New(f.repo, threads.NewOptions{Title: "May I?", On: f.story.ID, Author: "claude", Text: "Shall I?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	elsewhere, err := threads.New(f.repo, threads.NewOptions{Title: "There?", On: other.ID, Author: "alex", Text: "And there?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}

	in := f.read(t, Options{Agent: "claude"})
	if in.Agent != "claude" || in.AwaitingYou != 2 || strings.Join(threadIDs(in), " ") != asked.ID+" "+elsewhere.ID {
		t.Errorf("awaiting claude: %d %v, want %s and %s", in.AwaitingYou, threadIDs(in), asked.ID, elsewhere.ID)
	}
	if in.Threads[0].Awaiting != "you" || in.Threads[0].LastBy != "alex" || in.Threads[0].Story != f.story.ID {
		t.Errorf("summary: %+v", in.Threads[0])
	}
	if got := threadIDs(f.read(t, Options{Agent: "claude", All: true})); len(got) != 3 {
		t.Errorf("all lists the thread awaiting someone else too: %v (%s)", got, mine.ID)
	}
	if got := threadIDs(f.read(t, Options{Agent: "claude", Story: other.ID})); len(got) != 1 || got[0] != elsewhere.ID {
		t.Errorf("story keeps its own threads: %v", got)
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
	again, omitted, err := CatchUp(f.repo, "claude", "", f.clock, items)
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
		events, _, err := CatchUp(f.repo, "claude", story, f.clock, items)
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

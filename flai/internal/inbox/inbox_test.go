package inbox

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	again, omitted, err := CatchUp(f.repo, "claude", f.clock, items)
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

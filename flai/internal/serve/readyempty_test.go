package serve

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ADR-0211 §1: ready-empty is told on the look that finds ready below the
// room after one that did not, and with no limit or no room ready is never
// below.
func TestReadyDrainedDecidesOnceADraining(t *testing.T) {
	cases := []struct {
		name           string
		ready, room    int
		wasBelow       bool
		emit, nowBelow bool
	}{
		{"drains below the room", 1, 2, false, true, true},
		{"empty column below the room", 0, 3, false, true, true},
		{"stays below", 0, 2, true, false, true},
		{"at the room", 2, 2, false, false, false},
		{"reaches the room again", 2, 2, true, false, false},
		{"above the room", 4, 2, true, false, false},
		{"room zero", 0, 0, false, false, false},
		{"room zero after below", 0, 0, true, false, false},
		{"no limit", 0, -1, false, false, false},
		{"no limit after below", 0, -1, true, false, false},
	}
	for _, c := range cases {
		emit, below := readyDrained(c.ready, c.room, c.wasBelow)
		if emit != c.emit || below != c.nowBelow {
			t.Errorf("%s: readyDrained(%d, %d, %v) = %v, %v; want %v, %v", c.name, c.ready, c.room, c.wasBelow, emit, below, c.emit, c.nowBelow)
		}
	}
}

// readyEmptyLab is a project whose launcher has the agent action off, so that
// its looks start nothing and the ready column holds what the test puts in it.
type readyEmptyLab struct {
	t    *testing.T
	root string
	repo *workitem.Repo
	epic *workitem.Item
	o    Options
	l    *launcher
	logs bytes.Buffer
}

func newReadyEmptyLab(t *testing.T, limit int) *readyEmptyLab {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	board := fmt.Sprintf("---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  in-progress: %d\n  review: 3\n---\n\n# Board\n", limit)
	if err := os.WriteFile(filepath.Join(root, "wip/kanban/board.md"), []byte(board), 0o644); err != nil {
		t.Fatal(err)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	lab := &readyEmptyLab{t: t, root: root, repo: repo}
	if lab.epic, err = repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	dir := Dir(filepath.Join(t.TempDir(), "serve"))
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	lab.o = Options{
		Dir:    dir,
		Logger: slog.New(slog.NewTextHandler(&lab.logs, nil)),
		Now:    time.Now,
		Agent:  func(string) AgentConfig { return AgentConfig{} },
		Host:   hostapi.Host{Record: func(hostapi.Entry) {}},
	}
	lab.restart()
	return lab
}

// restart is flai serve starting again on the same project: a new launcher,
// the same .flai-cache.
func (lab *readyEmptyLab) restart() {
	lab.l = newLauncher(lab.o, Entry{Key: "t", Name: "t", Root: lab.root})
}

// looks looks n times, as flai serve does each minute.
func (lab *readyEmptyLab) looks(n int) {
	for range n {
		lab.l.look(context.Background(), false)
	}
}

// ready makes a story with a criterion and a path of its own and moves it to
// ready.
func (lab *readyEmptyLab) ready(title string) string {
	lab.t.Helper()
	return lab.readyWith(title, []string{"docs/" + strings.ToLower(title)}, nil)
}

// readyWith makes a story with a criterion, touches, and after, and moves it
// to ready.
func (lab *readyEmptyLab) readyWith(title string, touches, after []string) string {
	lab.t.Helper()
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex",
		Touches: touches, After: after, Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, err := os.ReadFile(st.Path)
	if err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Ready)
	return st.ID
}

func (lab *readyEmptyLab) move(id, to string) {
	lab.t.Helper()
	it, err := lab.repo.Get(id)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(it, to, "alex", "", time.Now()); err != nil {
		lab.t.Fatal(err)
	}
}

// events are the ready-empty events in the project's event log.
func (lab *readyEmptyLab) events() []inbox.LoggedEvent {
	var out []inbox.LoggedEvent
	for _, e := range inbox.LoggedEvents(lab.repo) {
		if e.Kind == inbox.ReadyEmpty {
			out = append(out, e)
		}
	}
	return out
}

func (lab *readyEmptyLab) want(n int, when string) {
	lab.t.Helper()
	if got := lab.events(); len(got) != n {
		lab.t.Fatalf("%s: %d ready-empty events, want %d: %+v", when, len(got), n, got)
	}
	if strings.Contains(lab.logs.String(), "level=WARN") {
		lab.t.Fatalf("%s: warned:\n%s", when, lab.logs.String())
	}
}

// S-0535 criterion 2, ADR-0211 §1: flai serve tells one ready-empty event
// when the ready column drains below the in-progress room, none while it stays
// below, a restart included, and one again after ready reaches the room and
// drains again.
func TestTheLauncherTellsOneReadyEmptyADraining(t *testing.T) {
	lab := newReadyEmptyLab(t, 2)
	a, b := lab.ready("A"), lab.ready("B")
	lab.looks(1)
	lab.want(0, "ready 2 at the room 2")

	lab.move(a, workitem.Backlog)
	lab.looks(1)
	lab.want(1, "ready 1 below the room 2")
	e := lab.events()[0]
	if e.ID != "board" || e.To != "1" || e.By != "flai serve" || e.Fields[inbox.FieldRoom] != "2" || e.Fields[inbox.FieldLimit] != "2" {
		t.Errorf("the event: %+v", e)
	}

	lab.looks(3)
	lab.move(b, workitem.Backlog)
	lab.looks(3)
	lab.want(1, "ready staying below over several looks, down to an empty column")

	lab.restart()
	lab.looks(2)
	lab.want(1, "a restart with ready still below")

	lab.move(a, workitem.Ready)
	lab.move(b, workitem.Ready)
	lab.looks(2)
	lab.want(1, "ready back at the room")

	lab.move(b, workitem.Backlog)
	lab.looks(2)
	lab.want(2, "ready drained again")
	if e := lab.events()[1]; e.To != "1" || e.Fields[inbox.FieldRoom] != "2" {
		t.Errorf("the second event: %+v", e)
	}
}

// ADR-0211 §1: with the in-progress limit full there is no room to fall
// below, and with no limit there is no room at all, so no event.
func TestNoReadyEmptyWithoutRoom(t *testing.T) {
	t.Run("the limit full", func(t *testing.T) {
		lab := newReadyEmptyLab(t, 1)
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		lab.looks(3)
		lab.want(0, "ready 0 with room 0")
	})
	t.Run("no limit", func(t *testing.T) {
		lab := newReadyEmptyLab(t, 0)
		lab.looks(3)
		lab.want(0, "ready 0 with no in-progress limit")
	})
}

// I-0285, ADR-0279: a ready column whose stories are all held on after is
// drained, since none of them can fill a lane, so ready-empty is told and the
// orchestrator is started to promote. Before ADR-0279 the held stories were
// counted, two at a room of two, and nothing was told while lanes stood empty.
// A story held on overlap alone still counts: a trial merge may clear it at the
// next look.
func TestReadyEmptyLeavesOutStoriesHeldOnAfter(t *testing.T) {
	t.Run("held on after", func(t *testing.T) {
		lab := newReadyEmptyLab(t, 3)
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		lab.readyWith("A", []string{"docs/a"}, []string{busy})
		lab.readyWith("B", []string{"docs/b"}, []string{busy})
		lab.looks(2)
		lab.want(1, "two ready stories held on after at a room of two")
		if e := lab.events()[0]; e.To != "0" || e.Fields[inbox.FieldRoom] != "2" || e.Fields[inbox.FieldLimit] != "3" {
			t.Errorf("the event: %+v", e)
		}
	})
	t.Run("held on overlap", func(t *testing.T) {
		lab := newReadyEmptyLab(t, 3)
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		lab.readyWith("A", []string{"docs/busy"}, nil)
		lab.readyWith("B", []string{"docs/busy"}, nil)
		lab.looks(2)
		lab.want(0, "two ready stories held on overlap at a room of two")
	})
}

// ADR-0279: the count leaves out a story held on after, alone or with an
// overlap, and keeps one held on overlap alone and one not held.
func TestReadyNotHeldOnAfterCounts(t *testing.T) {
	lab := newReadyEmptyLab(t, 5)
	busy := lab.ready("Busy")
	lab.move(busy, workitem.InProgress)
	lab.readyWith("Free", []string{"docs/free"}, nil)
	lab.readyWith("Overlap", []string{"docs/busy"}, nil)
	lab.readyWith("After", []string{"docs/after"}, []string{busy})
	lab.readyWith("Both", []string{"docs/busy"}, []string{busy})
	stories, holds, _, _, _, err := readyStories(lab.repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(stories) != 4 {
		t.Fatalf("%d ready stories, want 4", len(stories))
	}
	if got := readyNotHeldOnAfter(stories, holds); got != 2 {
		t.Errorf("readyNotHeldOnAfter = %d, want 2 (Free and Overlap)", got)
	}
	if got := readyNotHeldOnAfter(stories, nil); got != 4 {
		t.Errorf("readyNotHeldOnAfter with no holds = %d, want 4", got)
	}
}

// t3215Lab is a starved-board lab running kindsWorkflow, whose board limits
// doing at 2 and testing at 1, both of the active kind, and awaiting, of the
// waiting kind, at 3. The board names no in-progress limit, so the one flai
// keeps by default, 2, is all a pass that reads the state's name finds.
func t3215Lab(t *testing.T) *readyEmptyLab {
	t.Helper()
	lab := newStarvedLab(t, 0)
	lab.t3215Limits()
	lab.o.Workflow = kindsWorkflow(t)
	lab.repo.Workflow = lab.o.Workflow
	lab.restart()
	return lab
}

// t3215Limits writes the board's limits on doing, testing, and awaiting
// again, keeping its pull order. Board.Save, which each move of a story
// runs, writes the limits of the baseline's columns alone, so a move drops
// them; the lab writes them again before each look.
func (lab *readyEmptyLab) t3215Limits() {
	lab.t.Helper()
	board, err := lab.repo.LoadBoard()
	if err != nil {
		lab.t.Fatal(err)
	}
	text := "---\ntitle: Board\nstatus: active\nwip_limits:\n  ready: 5\n  doing: 2\n  testing: 1\n  awaiting: 3\norder: [" + strings.Join(board.Order, ", ") + "]\n---\n\n# Board\n"
	if err := os.WriteFile(filepath.Join(lab.root, "wip/kanban/board.md"), []byte(text), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// t3215Look is the ready-empty and starved-board part of a look, on the
// project as the lab's launcher opens it, under its workflow: the ready
// stories and their holds as readyStories reads them, against the room
// activeRoom gives from the board's limits and counts. It answers that room
// and the limit it is of.
func t3215Look(lab *readyEmptyLab) (room, limit int) {
	lab.t.Helper()
	lab.t3215Limits()
	repo, err := lab.l.open(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	stories, holds, _, _, _, err := readyStories(repo)
	if err != nil {
		lab.t.Fatal(err)
	}
	items, err := repo.List(false)
	if err != nil {
		lab.t.Fatal(err)
	}
	board, err := repo.LoadBoard()
	if err != nil {
		lab.t.Fatal(err)
	}
	view := workitem.NewBoardView(items, board, time.Now(), false, nil, holds)
	room, limit = activeRoom(repo.Flow(), view.WIPLimits, view.Counts)
	lab.l.noteReadyEmpty(repo, readyNotHeldOnAfter(stories, holds), room)
	lab.l.noteStarved(repo, stories, holds, room)
	return room, limit
}

// t3215Through moves the story id from ready through each of states in turn.
func (lab *readyEmptyLab) t3215Through(id string, states ...string) {
	lab.t.Helper()
	for _, s := range states {
		lab.move(id, s)
	}
}

// S-0742 criteria 3 and 4, ADR-0279, ADR-0285: under a workflow with
// in-progress renamed doing, testing added after it, and review renamed
// awaiting, the ready-empty event and the starved-board thread count the
// stories in the queue the pull takes from against the room the limits of
// doing and testing leave together. A story in doing and one in testing both
// fill it, and one in awaiting does not; the event's limit is their sum.
func TestReadyEmptyAndTheStarvedBoardCountTheActiveWorkRoom(t *testing.T) {
	lab := t3215Lab(t)
	lab.backlog("Count the turns", true)
	waits := lab.ready("Waits")
	lab.t3215Through(waits, "doing", "testing", "awaiting")
	tests := lab.ready("Tests")
	lab.t3215Through(tests, "doing", "testing")
	works := lab.ready("Works")
	lab.t3215Through(works, "doing")

	room, limit := t3215Look(lab)
	if room != 1 || limit != 3 {
		t.Fatalf("one story in doing and one in testing of limits 2 and 1 leave room %d of %d, want 1 of 3", room, limit)
	}
	lab.want(1, "an empty ready column below an active-work room of 1")
	if e := lab.events()[0]; e.To != "0" || e.Fields[inbox.FieldRoom] != "1" || e.Fields[inbox.FieldLimit] != "3" {
		t.Errorf("the event: %+v, want 0 ready of room 1, limit 3", e)
	}
	th := lab.starved(1, "a free active-work lane, the story in awaiting filling none, and nothing to start")

	// sent back from awaiting into doing, the story fills the last lane
	lab.move(waits, "doing")
	if room, _ := t3215Look(lab); room != 0 {
		t.Fatalf("two stories in doing and one in testing leave room %d, want 0", room)
	}
	lab.want(1, "no room left, so no second event")
	for _, b := range lab.boardThreads() {
		if b.Open() {
			t.Errorf("the starved-board thread %s is open with every active-work lane taken (%s was opened)", b.ID, th.ID)
		}
	}

	// a story pulled into ready fills the room the move to testing frees
	lab.move(tests, "awaiting")
	ready := lab.ready("Next")
	if room, _ := t3215Look(lab); room != 1 {
		t.Fatalf("after the story in testing moved to awaiting, room %d, want 1", room)
	}
	lab.want(1, "ready at the room of 1")
	if again := lab.boardThreads(); len(again) != 1 {
		t.Errorf("%d threads on the board with %s in ready able to start, want the one resolved", len(again), ready)
	}
}

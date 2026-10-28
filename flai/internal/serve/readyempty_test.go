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
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex",
		Touches: []string{"docs/" + strings.ToLower(title)}, Now: time.Now()})
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

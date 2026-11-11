package serve

import (
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newStarvedLab is a readyEmptyLab whose clock moves on a minute at each
// reading, as flai serve's looks do, so that each entry flai serve writes on
// a thread is one of its own rather than part of the last.
func newStarvedLab(t *testing.T, limit int) *readyEmptyLab {
	t.Helper()
	lab := newReadyEmptyLab(t, limit)
	at := time.Now()
	lab.o.Now = func() time.Time {
		at = at.Add(time.Minute)
		return at
	}
	lab.restart()
	return lab
}

// backlog makes a story in the backlog, a draft when draft is set, held by
// after behind the stories after names.
func (lab *readyEmptyLab) backlog(title string, draft bool, after ...string) string {
	lab.t.Helper()
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: lab.epic.ID, Owner: "alex",
		Touches: []string{"docs/" + strings.ToLower(title)}, After: after, Draft: draft, Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	return st.ID
}

// ask opens a thread on story by its planner, which awaits the operator.
func (lab *readyEmptyLab) ask(story, title string) string {
	lab.t.Helper()
	th, err := threads.New(lab.repo, threads.NewOptions{Title: title, On: story, Author: "planner-" + story, Text: "Which way?", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	return th.ID
}

// boardThreads are the threads on the board file, open or resolved.
func (lab *readyEmptyLab) boardThreads() []*threads.Thread {
	lab.t.Helper()
	all, err := threads.For(lab.repo, "wip/kanban/board.md")
	if err != nil {
		lab.t.Fatal(err)
	}
	return all
}

// starved is the one open thread flai serve keeps on the starved board, with
// want entries, and no warning logged.
func (lab *readyEmptyLab) starved(want int, when string) *threads.Thread {
	lab.t.Helper()
	if strings.Contains(lab.logs.String(), "level=WARN") {
		lab.t.Fatalf("%s: warned:\n%s", when, lab.logs.String())
	}
	var open []*threads.Thread
	for _, th := range lab.boardThreads() {
		if th.Open() {
			open = append(open, th)
		}
	}
	if len(open) != 1 {
		lab.t.Fatalf("%s: %d open threads on the board, want 1: %+v", when, len(open), open)
	}
	th := open[0]
	if th.Title != starvedTitle || th.Opener() != "flai serve" {
		lab.t.Fatalf("%s: the thread is %q by %q", when, th.Title, th.Opener())
	}
	if got := len(th.Entries()); got != want {
		lab.t.Fatalf("%s: %d entries on %s, want %d:\n%s", when, got, th.ID, want, th.Body)
	}
	return th
}

// I-0288, S-0691 criterion 1, ADR-0285: the ready column empty, the
// in-progress limit with room, and the backlog only two drafts, each with a
// thread awaiting the operator, and stories held by after behind them. One
// look opens one thread on the board naming both drafts and their threads;
// later looks, a restart among them, write nothing; a change to the list adds
// one entry; and a story that enters ready and can start resolves it, after
// which the next starving opens a new one.
func TestAStarvedBoardOfI0288TellsTheOperatorOnce(t *testing.T) {
	lab := newStarvedLab(t, 2)
	first := lab.backlog("Count the turns", true)
	asked := lab.ask(first, "Which turns count?")
	second := lab.backlog("Price the tokens", true)
	priced := lab.ask(second, "Which price list")
	lab.backlog("Chart the turns", false, first)
	lab.backlog("Chart the price", false, second)

	lab.looks(1)
	th := lab.starved(1, "the first look on the starved board")
	text := th.Entries()[0].Text
	for _, want := range []string{
		"2 in-progress lanes stand idle",
		"- " + first + " Count the turns: a draft; finalize it. It waits on " + asked + " Which turns count?\n",
		"- " + second + " Price the tokens: a draft; finalize it. It waits on " + priced + " Which price list.",
		"Chart the turns: held by `after`; the root of its chain is " + first + " Count the turns (backlog).",
		"`flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("the thread's text lacks %q:\n%s", want, text)
		}
	}

	lab.looks(2)
	lab.restart()
	lab.looks(2)
	lab.starved(1, "looks and a restart with the list unchanged")

	lab.backlog("Alert on turns", false, first)
	lab.looks(1)
	th = lab.starved(2, "a story added to the list")
	if e := th.Entries()[1]; !strings.HasPrefix(e.Text, "The list changed.") || !strings.Contains(e.Text, "Alert on turns") {
		t.Errorf("the second entry:\n%s", e.Text)
	}
	lab.looks(2)
	lab.starved(2, "looks after the change")

	ship := lab.ready("Ship")
	lab.looks(1)
	resolved, err := threads.Get(lab.repo, th.ID)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Open() {
		t.Fatalf("a story that can start leaves %s open:\n%s", th.ID, resolved.Body)
	}
	if !strings.Contains(resolved.Body, ship+" Ship in ready can start") {
		t.Errorf("the resolve does not name %s:\n%s", ship, resolved.Body)
	}
	if st, err := readStarved(starvedPath(lab.repo)); err != nil || st != (starvedState{}) {
		t.Errorf("the state after the resolve: %+v, %v", st, err)
	}

	lab.move(ship, workitem.Backlog)
	lab.looks(2)
	if again := lab.starved(1, "the next starving"); again.ID == th.ID {
		t.Errorf("the next starving reopened %s rather than opening a new thread", th.ID)
	}
	if n := len(lab.boardThreads()); n != 2 {
		t.Errorf("%d threads on the board, want 2", n)
	}
}

// ADR-0285 §3: a host with a reminder channel is not told to set one; a host
// without one is (above).
func TestAStarvedBoardWithAReminderChannelNamesNone(t *testing.T) {
	lab := newStarvedLab(t, 1)
	lab.o.Notify = func() config.Notify { return config.Notify{Command: "notify-send flai"} }
	lab.restart()
	lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
	lab.looks(1)
	text := lab.starved(1, "a starved board on a host with a channel").Entries()[0].Text
	if strings.Contains(text, "notify set") || !strings.Contains(text, "1 in-progress lane stands idle") {
		t.Errorf("the thread's text:\n%s", text)
	}
}

// ADR-0285 §1, §5: a full in-progress limit, or none, starves no board; a
// lane taken resolves the open thread, saying so; and a thread the operator
// resolved while the board stayed starved is not opened again until its
// list changes.
func TestABoardWithNoFreeLaneIsNotStarved(t *testing.T) {
	t.Run("no limit", func(t *testing.T) {
		lab := newStarvedLab(t, 0)
		lab.backlog("Count the turns", true)
		lab.looks(2)
		if n := len(lab.boardThreads()); n != 0 {
			t.Errorf("%d threads with no in-progress limit", n)
		}
	})
	t.Run("a lane taken", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.backlog("Count the turns", true)
		lab.looks(1)
		th := lab.starved(1, "a free lane and nothing to start")
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		lab.looks(1)
		got, err := threads.Get(lab.repo, th.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Open() || !strings.Contains(got.Body, "Every in-progress lane is taken") {
			t.Errorf("a full limit leaves:\n%s", got.Body)
		}
	})
	t.Run("resolved by the operator", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.backlog("Count the turns", true)
		lab.looks(1)
		th := lab.starved(1, "a free lane and nothing to start")
		if _, err := threads.Resolve(lab.repo, th.ID, "alex", "seen", time.Now()); err != nil {
			t.Fatal(err)
		}
		lab.looks(2)
		if n := len(lab.boardThreads()); n != 1 {
			t.Fatalf("%d threads after the operator resolved the first, want 1", n)
		}
		lab.backlog("Price the tokens", true)
		lab.looks(1)
		if again := lab.starved(1, "the list changed after the operator resolved it"); again.ID == th.ID {
			t.Errorf("the changed list reopened %s", th.ID)
		}
	})
}

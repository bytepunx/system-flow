package serve

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newStarvedLab is a readyEmptyLab whose clock moves on a minute at each
// reading, as flai serve's looks do, so that each entry flai serve writes on
// a thread is one of its own rather than part of the last, and further when
// the test waits.
func newStarvedLab(t *testing.T, limit int) *readyEmptyLab {
	t.Helper()
	lab := newReadyEmptyLab(t, limit)
	lab.clock = time.Now()
	lab.o.Now = func() time.Time {
		lab.clock = lab.clock.Add(time.Minute)
		return lab.clock
	}
	lab.restart()
	return lab
}

// wait moves a newStarvedLab's clock on by d, as time passing between looks.
func (lab *readyEmptyLab) wait(d time.Duration) {
	lab.clock = lab.clock.Add(d)
}

// settle looks, waits the settle window, and looks again, so that the
// board's state has held for the window at the second look (ADR-0321).
func (lab *readyEmptyLab) settle() {
	lab.looks(1)
	lab.wait(starvedSettle)
	lab.looks(1)
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
// want entries, or any number when want is below zero, and no warning logged.
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
	if got := len(th.Entries()); want >= 0 && got != want {
		lab.t.Fatalf("%s: %d entries on %s, want %d:\n%s", when, got, th.ID, want, th.Body)
	}
	return th
}

// I-0288, S-0691 criterion 1, ADR-0285, ADR-0321: the ready column empty,
// the in-progress limit with room, and the backlog only two drafts, each with
// a thread awaiting the operator, and stories held by after behind them. The
// first look opens nothing; a look once the board has been starved for the
// settle window opens one thread on the board naming both drafts and their
// threads; later looks, a restart among them, write nothing; a change to the
// list adds one entry once the window has passed since the last; and a story
// that enters ready and can start resolves it once it has been able to for
// the window, after which the next starving opens a new one.
func TestAStarvedBoardOfI0288TellsTheOperatorOnce(t *testing.T) {
	lab := newStarvedLab(t, 2)
	first := lab.backlog("Count the turns", true)
	asked := lab.ask(first, "Which turns count?")
	second := lab.backlog("Price the tokens", true)
	priced := lab.ask(second, "Which price list")
	lab.backlog("Chart the turns", false, first)
	lab.backlog("Chart the price", false, second)

	lab.looks(1)
	if n := len(lab.boardThreads()); n != 0 {
		t.Fatalf("%d threads at the first look on the starved board, want none within the settle window", n)
	}
	lab.wait(starvedSettle)
	lab.looks(1)
	th := lab.starved(1, "a look once the board has been starved for the window")
	text := th.Entries()[0].Text
	for _, want := range []string{
		"2 in-progress lanes stand idle",
		"- " + first + " Count the turns: a draft; finalize it. It waits on " + asked + " Which turns count?\n",
		"- " + second + " Price the tokens: a draft; finalize it. It waits on " + priced + " Which price list.",
		"Chart the turns: held by `after`; the root of its chain is " + first + " Count the turns (backlog).",
		"`flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`",
		"at most once every 15 minutes",
		"once a story has been able to start, or be promoted, for 15 minutes.",
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
	lab.wait(starvedSettle)
	lab.looks(1)
	th = lab.starved(2, "a story added to the list")
	if e := th.Entries()[1]; !strings.HasPrefix(e.Text, "The list changed.") || !strings.Contains(e.Text, "Alert on turns") {
		t.Errorf("the second entry:\n%s", e.Text)
	}
	lab.looks(2)
	lab.starved(2, "looks after the change")

	ship := lab.ready("Ship")
	lab.looks(1)
	lab.starved(2, "a story that can start within the settle window")
	lab.wait(starvedSettle)
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
	lab.settle()
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
	lab.settle()
	text := lab.starved(1, "a starved board on a host with a channel").Entries()[0].Text
	if strings.Contains(text, "notify set") || !strings.Contains(text, "1 in-progress lane stands idle") {
		t.Errorf("the thread's text:\n%s", text)
	}
}

// ADR-0285 §1, §5, ADR-0321: a full in-progress limit, or none, starves no
// board; a lane taken for the settle window resolves the open thread, saying
// so; and a thread the operator resolved while the board stayed starved is
// not opened again until its list changes, and then at once.
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
		lab.settle()
		th := lab.starved(1, "a free lane and nothing to start")
		busy := lab.ready("Busy")
		lab.move(busy, workitem.InProgress)
		lab.settle()
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
		lab.settle()
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

// entriesApart fails the test when two of flai serve's entries on th are
// less than the settle window apart (ADR-0321 §3).
func (lab *readyEmptyLab) entriesApart(th *threads.Thread) {
	lab.t.Helper()
	var last time.Time
	for _, e := range th.Entries() {
		if e.Author != readyEmptyAuthor {
			continue
		}
		at, err := time.Parse(workitem.TimeFormat, e.At)
		if err != nil {
			lab.t.Fatal(err)
		}
		if !last.IsZero() && at.Sub(last) < starvedSettle {
			lab.t.Errorf("entries on %s at %s and %s, less than %s apart", th.ID, last.Format(time.RFC3339), e.At, starvedSettle)
		}
		last = at
	}
}

// I-0389, S-0847 criterion 1, ADR-0321: while drafts were finalized and
// promoted one at a time the board flipped between starved and not from one
// look to the next, and flai serve opened a thread at each starving and
// resolved it at the next look, 48 in five days, and wrote an entry at each
// change of the list, 12 in four minutes on TH-1134. A board that flips
// faster than the settle window opens no thread; a board that stays starved
// for it opens one and keeps it through a flip back shorter than the window;
// and a list that changes at each look gains at most one entry a window, the
// newest list.
func TestAFlippingBoardOfI0389OpensOneThreadAndWritesOnceAWindow(t *testing.T) {
	t.Run("flips faster than the window", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
		ship := lab.ready("Ship")
		start := lab.clock
		for i := 0; lab.clock.Sub(start) < 3*starvedSettle; i++ {
			if i%2 == 0 {
				lab.move(ship, workitem.Backlog)
			} else {
				lab.move(ship, workitem.Ready)
			}
			lab.looks(1)
		}
		if strings.Contains(lab.logs.String(), "level=WARN") {
			t.Fatalf("warned:\n%s", lab.logs.String())
		}
		if all := lab.boardThreads(); len(all) != 0 {
			t.Errorf("a board flipping at each look for %s opened %d threads, want none", 3*starvedSettle, len(all))
		}
	})
	t.Run("a flip back shorter than the window", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
		lab.settle()
		th := lab.starved(1, "a board starved for the window")
		ship := lab.ready("Ship")
		for range 3 {
			lab.looks(1)
			lab.wait(starvedSettle / 3)
			lab.looks(1)
			lab.starved(-1, "a board not starved for less than the window")
			lab.move(ship, workitem.Backlog)
			lab.looks(1)
			if again := lab.starved(-1, "starved again within the window"); again.ID != th.ID {
				t.Fatalf("a flip back opened %s rather than keeping %s", again.ID, th.ID)
			}
			lab.move(ship, workitem.Ready)
		}
		if all := lab.boardThreads(); len(all) != 1 {
			t.Errorf("%d threads on the board after flips back, want the one", len(all))
		}
	})
	t.Run("a list that changes at each look", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
		lab.settle()
		th := lab.starved(1, "a board starved for the window")
		lab.backlog("Alert 0", true)
		lab.looks(1)
		lab.starved(1, "a changed list within the window of the open")
		opened, err := time.Parse(workitem.TimeFormat, th.Entries()[0].At)
		if err != nil {
			t.Fatal(err)
		}
		start, last := lab.clock, ""
		for i := 1; lab.clock.Sub(start) < 3*starvedSettle; i++ {
			last = fmt.Sprintf("Alert %d", i)
			lab.backlog(last, true)
			lab.looks(1)
		}
		th = lab.starved(-1, "a list changing at each look")
		lab.entriesApart(th)
		most := 1 + int(lab.clock.Sub(opened)/starvedSettle)
		if n := len(th.Entries()); n < 3 || n > most {
			t.Errorf("%d entries on %s over %s of a list changing at each look, want 3 to %d, one a window:\n%s", n, th.ID, lab.clock.Sub(opened), most, th.Body)
		}
		lab.wait(starvedSettle)
		lab.looks(1)
		th = lab.starved(-1, "a look a window after the last entry")
		lab.entriesApart(th)
		if e := th.Entries()[len(th.Entries())-1]; !strings.Contains(e.Text, last+":") {
			t.Errorf("the last entry lacks the newest story %q:\n%s", last, e.Text)
		}
	})
}

// ADR-0321 §2, §4: the times the settle window runs from are kept in
// .flai-cache/starved.json, so a restart neither opens nor resolves the
// thread early, and the thread is resolved once the board has stayed not
// starved for the window. When the file cannot be read, the open thread is
// found by its title and the window starts from that look.
func TestTheStarvedBoardsSettleWindowSurvivesARestart(t *testing.T) {
	t.Run("kept", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
		lab.looks(1)
		lab.restart()
		if st, err := readStarved(starvedPath(lab.repo)); err != nil || st.Starved == "" || st.Thread != "" {
			t.Fatalf("the state of a board starved within the window: %+v, %v", st, err)
		}
		lab.wait(starvedSettle)
		lab.looks(1)
		th := lab.starved(1, "a look a window after the starving, across a restart")
		st, err := readStarved(starvedPath(lab.repo))
		if err != nil || st.Thread != th.ID || st.Wrote == "" || st.Starved == "" {
			t.Fatalf("the state once the thread opened: %+v, %v", st, err)
		}

		lab.backlog("Price the tokens", true)
		lab.restart()
		lab.looks(1)
		lab.starved(1, "a changed list within the window of the open, across a restart")

		ship := lab.ready("Ship")
		lab.looks(1)
		lab.restart()
		if st, err := readStarved(starvedPath(lab.repo)); err != nil || st.Clear == "" || st.Starved != "" || st.Thread != th.ID {
			t.Fatalf("the state of a board not starved within the window: %+v, %v", st, err)
		}
		lab.looks(1)
		lab.starved(1, "a story that can start within the window, across a restart")
		lab.wait(starvedSettle)
		lab.looks(1)
		got, err := threads.Get(lab.repo, th.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Open() || !strings.Contains(got.Body, ship+" Ship in ready can start") {
			t.Errorf("a board not starved for the window leaves:\n%s", got.Body)
		}
		if st, err := readStarved(starvedPath(lab.repo)); err != nil || st != (starvedState{}) {
			t.Errorf("the state after the resolve: %+v, %v", st, err)
		}
	})
	t.Run("lost", func(t *testing.T) {
		lab := newStarvedLab(t, 1)
		lab.ask(lab.backlog("Count the turns", true), "Which turns count?")
		lab.settle()
		th := lab.starved(1, "a board starved for the window")
		if err := os.WriteFile(starvedPath(lab.repo), []byte("{"), 0o600); err != nil {
			t.Fatal(err)
		}
		lab.ready("Ship")
		lab.looks(1)
		if !strings.Contains(lab.logs.String(), "could not be read") {
			t.Errorf("an unreadable state is not warned of:\n%s", lab.logs.String())
		}
		got, err := threads.Get(lab.repo, th.ID)
		if err != nil {
			t.Fatal(err)
		}
		if !got.Open() {
			t.Fatalf("an unreadable state resolved %s at once:\n%s", th.ID, got.Body)
		}
		if st, err := readStarved(starvedPath(lab.repo)); err != nil || st.Thread != th.ID || st.Clear == "" {
			t.Fatalf("the state after an unreadable one: %+v, %v", st, err)
		}
		lab.wait(starvedSettle)
		lab.looks(1)
		if got, err = threads.Get(lab.repo, th.ID); err != nil {
			t.Fatal(err)
		}
		if got.Open() {
			t.Errorf("a board not starved for the window after a lost state leaves %s open:\n%s", th.ID, got.Body)
		}
	})
}

package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A starved board (ADR-0285, I-0288): no story in ready can start, the
// in-progress limit has room, and no backlog story is a promotion candidate.
// flai serve tells the operator on one thread of its own, on the board file,
// each backlog story that cannot be promoted and the decision it waits on,
// adds an entry only when that list changes, and resolves the thread once a
// story can start. It does each only once the board has held its state for
// starvedSettle (ADR-0321, I-0389), so a board that flips between starved and
// not from one look to the next opens no thread, keeps the one it has, and
// gains at most one entry a window.

// starvedSettle is the settle window of ADR-0321: how long the board must
// stay starved before its thread opens, how long it must stay not starved
// before the thread is resolved, and the least time between two entries
// flai serve writes there.
const starvedSettle = 15 * time.Minute

// starvedTitle is the title of flai serve's thread on a starved board. It is
// fixed, so that a restart finds the open thread by it; the counts are in the
// text.
const starvedTitle = "The board is starved: nothing in ready can start and no backlog story can be promoted"

// starvedReader is who flai serve reads the threads as to list a starved
// board's decisions: the reminder's reader, a name no author takes, so that no
// thread awaits it and a thread on a draft awaits the operator exactly when
// the operator's own view says so (ADR-0181). Reading as flai serve would
// read a thread it took part in as awaiting it, not the operator.
var starvedReader = inbox.Reader{Agent: reminderReader}

// starvedState is what flai serve last wrote on its thread on a starved
// board, and when the board's state began, kept between looks and across a
// restart, so that the same list is written once and the settle window runs
// on. The times are RFC 3339: Starved is when a look first found the board
// starved, for the current stretch; Clear is when a look first found it not
// starved while the state names a thread; and Wrote is when flai serve last
// wrote on the thread, the open or a reply.
type starvedState struct {
	Thread  string `json:"thread,omitempty"`
	List    string `json:"list,omitempty"`
	Starved string `json:"starved,omitempty"`
	Clear   string `json:"clear,omitempty"`
	Wrote   string `json:"wrote,omitempty"`
}

// starvedLook is what one look at the board starts from: where the state is
// kept, what it held, whether it could not be read, and the look's time.
type starvedLook struct {
	path string
	was  starvedState
	lost bool
	now  time.Time
}

// starvedPath is where a project's starvedState is kept: in its main
// checkout's .flai-cache, beside the ready-empty flag.
func starvedPath(repo *workitem.Repo) string {
	return filepath.Join(repo.CacheDir(), "starved.json")
}

// noteStarved judges, at a look, whether the board is starved, from the
// ready stories as readyStories answers them, their holds, and the room the
// in-progress limit leaves (below zero with no limit). A board starved for
// the settle window gets flai serve's thread, or an entry on it when the list
// changed and the window has passed since it last wrote; a board not starved
// for the window has the thread resolved. The backlog is read only when ready
// holds no story that can start and a lane is free, so a busy board pays
// nothing. Every failure is warned of and the look goes on.
func (l *launcher) noteStarved(repo *workitem.Repo, ready []readyStory, holds *workitem.Holds, room int) {
	k := starvedLook{path: starvedPath(repo), now: l.now()}
	var err error
	k.was, err = readStarved(k.path)
	k.lost = err != nil
	if k.lost {
		l.warn("starved board: what flai serve last wrote could not be read; its thread is looked for by its title, and the settle window starts again", "path", k.path, "err", err)
	}
	if why := notStarved(ready, holds, room); why != "" {
		l.unstarve(repo, k, why)
		return
	}
	on, err := boardPath(repo)
	if err != nil {
		l.warn("starved board: the board file's path could not be named, so the board is not judged", "err", err)
		return
	}
	if _, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(on))); errors.Is(err, fs.ErrNotExist) {
		// a project with no board file has no thread to anchor there, and
		// its limits are flai's defaults rather than the operator's
		return
	}
	items, err := repo.List(false)
	if err != nil {
		l.warn("starved board: the work items could not be read, so the board is not judged", "err", err)
		return
	}
	s, err := inbox.ReadStarved(repo, items, starvedReader)
	if err != nil {
		l.warn("starved board: the backlog's decisions could not be listed, so the board is not judged", "err", err)
		return
	}
	if len(s.Candidates) > 0 {
		c := s.Candidates[0]
		l.unstarve(repo, k, fmt.Sprintf("%s %s can be promoted from the backlog, so the board is no longer starved.", c.ID, oneLine(c.Title)))
		return
	}
	l.starve(repo, on, k, starvedList(s), room)
}

// notStarved says why ready and room starve no board, or "" when they would:
// no lane is free, there is no in-progress limit, or a ready story not held
// on after can start, counted as ready-empty counts it (ADR-0279).
func notStarved(ready []readyStory, holds *workitem.Holds, room int) string {
	switch {
	case room == 0:
		return "Every in-progress lane is taken, so the board is no longer starved."
	case room < 0:
		return "The board has no in-progress limit, so it is never starved."
	}
	for _, s := range ready {
		if readyNotHeldOnAfter([]readyStory{s}, holds) == 1 {
			name := s.ID
			if s.item != nil {
				name += " " + oneLine(s.item.Title)
			}
			return name + " in ready can start, so the board is no longer starved."
		}
	}
	return ""
}

// starve opens flai serve's thread on the starved board, on the board file
// on, with list, once the board has been starved for the settle window, or
// adds list to the open one when it differs from what flai serve last wrote
// there and the window has passed since it wrote, and keeps what it wrote. A
// changed list that comes sooner waits for a later look, which writes the
// newest. A thread the operator resolved while the board stayed starved is
// not opened again until the list changes.
func (l *launcher) starve(repo *workitem.Repo, on string, k starvedLook, list string, room int) {
	st := k.was
	st.Clear = ""
	if _, ok := sinceKept(st.Starved, k.now); !ok {
		st.Starved = keptTime(k.now)
	}
	open, err := starvedThreads(repo)
	if err != nil {
		l.warn("starved board: its thread could not be looked for, so the operator is not told; the next look tries again", "err", err)
		return
	}
	if len(open) == 0 {
		if starved, _ := sinceKept(st.Starved, k.now); starved < starvedSettle {
			l.keepStarved(k, st)
			return
		}
		if k.was.Thread != "" && k.was.List == list {
			l.keepStarved(k, st)
			return
		}
		th, err := threads.New(repo, threads.NewOptions{Title: starvedTitle, On: on, Author: readyEmptyAuthor, Text: l.starvedText(list, room, false), Now: k.now})
		if err != nil {
			l.warn("starved board: its thread could not be opened, so the operator is not told; the next look tries again", "on", on, "err", err)
			return
		}
		l.keepStarved(k, starvedState{Thread: th.ID, List: list, Starved: st.Starved, Wrote: keptTime(k.now)})
		l.log("starved board thread opened", "thread", th.ID, "room", room)
		return
	}
	th := open[0]
	if th.ID != k.was.Thread {
		st = starvedState{Thread: th.ID, Starved: st.Starved}
		if lastWrote(th, list) {
			st.List = list
		}
	}
	if st.Wrote == "" {
		st.Wrote = lastWroteAt(th)
	}
	if st.List == list {
		l.keepStarved(k, st)
		return
	}
	if wrote, ok := sinceKept(st.Wrote, k.now); ok && wrote < starvedSettle {
		l.keepStarved(k, st)
		return
	}
	if _, err := threads.Reply(repo, th.ID, readyEmptyAuthor, l.starvedText(list, room, true), k.now); err != nil {
		l.warn("starved board: the changed list could not be written on its thread; the next look tries again", "thread", th.ID, "err", err)
		return
	}
	st.List, st.Wrote = list, keptTime(k.now)
	l.keepStarved(k, st)
	l.log("starved board thread told a changed list", "thread", th.ID, "room", room)
}

// unstarve resolves flai serve's open thread on a starved board, saying why,
// once the board has not been starved for the settle window, and forgets what
// it wrote there, so that the next starving opens a new one. A board that
// starves again within the window keeps the thread. It looks for the thread
// only when the state names one, or could not be read; then it names the
// thread it finds, and the window starts at this look.
func (l *launcher) unstarve(repo *workitem.Repo, k starvedLook, why string) {
	if k.was.Thread == "" && !k.lost {
		l.keepStarved(k, starvedState{})
		return
	}
	st := k.was
	st.Starved = ""
	if k.lost {
		open, err := starvedThreads(repo)
		if err != nil {
			l.warn("starved board: its thread could not be looked for, so it is not resolved; the next look tries again", "err", err)
			return
		}
		st = starvedState{}
		if len(open) > 0 {
			st = starvedState{Thread: open[0].ID, Clear: keptTime(k.now)}
		}
		l.keepStarved(k, st)
		return
	}
	calm, ok := sinceKept(st.Clear, k.now)
	if !ok {
		st.Clear = keptTime(k.now)
	}
	if calm < starvedSettle {
		l.keepStarved(k, st)
		return
	}
	open, err := starvedThreads(repo)
	if err != nil {
		l.warn("starved board: its thread could not be looked for, so it is not resolved; the next look tries again", "err", err)
		return
	}
	var resolved []string
	for _, th := range open {
		if _, err := threads.Resolve(repo, th.ID, readyEmptyAuthor, why, k.now); err != nil {
			l.warn("starved board: its thread could not be resolved; the next look tries again", "thread", th.ID, "resolved", resolved, "err", err)
			return
		}
		resolved = append(resolved, th.ID)
	}
	l.keepStarved(k, starvedState{})
	if len(resolved) > 0 {
		l.log("starved board thread resolved", "threads", resolved, "why", why)
	}
}

// keepStarved keeps st at k's path when it differs from what the look read,
// or what was kept could not be read, warning when it cannot: the next look
// then finds the thread by its title.
func (l *launcher) keepStarved(k starvedLook, st starvedState) {
	if st == k.was && !k.lost {
		return
	}
	if err := writeStarved(k.path, st); err != nil {
		l.warn("starved board: what flai serve wrote on its thread could not be kept", "path", k.path, "thread", st.Thread, "err", err)
	}
}

// keptTime is t as a starvedState keeps it, RFC 3339 in UTC.
func keptTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// sinceKept is how long before now the kept time at was, and false when at
// is unset or is not RFC 3339, which the state treats as unset.
func sinceKept(at string, now time.Time) (time.Duration, bool) {
	if at == "" {
		return 0, false
	}
	t, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return 0, false
	}
	return now.Sub(t), true
}

// starvedText is an entry of flai serve's thread on a starved board: the
// lanes that stand idle, for a changed list that it changed, the missing
// reminder channel when the host has none, and list last, so that lastWrote
// finds it.
func (l *launcher) starvedText(list string, room int, changed bool) string {
	var b strings.Builder
	if changed {
		b.WriteString("The list changed. ")
	}
	lanes := fmt.Sprintf("%d in-progress lanes stand idle", room)
	if room == 1 {
		lanes = "1 in-progress lane stands idle"
	}
	b.WriteString("No story in ready can start, no backlog story can be promoted, and " + lanes + ".")
	if !changed {
		window := fmt.Sprintf("%d minutes", int(starvedSettle/time.Minute))
		b.WriteString(" flai serve adds an entry here when the list below changes, at most once every " + window + ", and resolves this thread once a story has been able to start, or be promoted, for " + window + ".")
	}
	b.WriteString("\n\n")
	if !l.hasReminderChannel() {
		b.WriteString("This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. ")
		b.WriteString("Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.\n\n")
	}
	b.WriteString(list)
	return b.String()
}

// hasReminderChannel says whether the host has a reminder channel set
// (ADR-0205), as the reminder judges it.
func (l *launcher) hasReminderChannel() bool {
	if l.notify == nil {
		return false
	}
	ch := l.notify()
	return ch.URL != "" || strings.TrimSpace(ch.Command) != ""
}

// starvedList is the list a starved board's thread gives: each backlog story
// and the decision it waits on, under a line that introduces them, or that
// the backlog holds none. The same backlog gives the same list.
func starvedList(s inbox.Starved) string {
	text := s.Text()
	if text == "" {
		return "The backlog holds no story that is not in ready: write or plan one to start work again.\n"
	}
	return "Each backlog story, in the pull order, and the decision it waits on:\n\n" + text
}

// lastWrote says whether flai serve's newest entry on th ends with list, for
// a thread the kept state does not name.
func lastWrote(th *threads.Thread, list string) bool {
	entries := th.Entries()
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Author == readyEmptyAuthor {
			return strings.HasSuffix(strings.TrimSpace(entries[i].Text), strings.TrimSpace(list))
		}
	}
	return false
}

// lastWroteAt is when flai serve wrote its newest entry on th, as the thread
// keeps it, or "" when it wrote none, for a state that does not say.
func lastWroteAt(th *threads.Thread) string {
	entries := th.Entries()
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Author == readyEmptyAuthor {
			return entries[i].At
		}
	}
	return ""
}

// starvedThreads are the open threads flai serve opened on the board file
// for a starved board, by their title.
func starvedThreads(repo *workitem.Repo) ([]*threads.Thread, error) {
	on, err := boardPath(repo)
	if err != nil {
		return nil, err
	}
	all, err := threads.For(repo, on)
	if err != nil {
		return nil, fmt.Errorf("read the threads on %s: %w; fix the thread file it names", on, err)
	}
	title := workitem.CleanTitle(starvedTitle)
	var out []*threads.Thread
	for _, th := range all {
		if th.Open() && th.Opener() == readyEmptyAuthor && th.Title == title {
			out = append(out, th)
		}
	}
	return out, nil
}

// boardPath is the board file's path in repo, wip/kanban/board.md under the
// project's layout, slash-separated.
func boardPath(repo *workitem.Repo) (string, error) {
	rel, err := filepath.Rel(repo.Root, filepath.Join(repo.KanbanDir(), workitem.BoardFile))
	if err != nil {
		return "", fmt.Errorf("name the board file under %s: %w; check layout.wip in system-flow.yaml", repo.Root, err)
	}
	return filepath.ToSlash(rel), nil
}

// readStarved reads the state kept at path; none kept is none written.
func readStarved(path string) (starvedState, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return starvedState{}, nil
	}
	if err != nil {
		return starvedState{}, err
	}
	var st starvedState
	if err := json.Unmarshal(data, &st); err != nil {
		return starvedState{}, fmt.Errorf("%s: %w; delete it, and flai serve keeps it afresh", path, err)
	}
	return st, nil
}

// writeStarved keeps st at path in one step.
func writeStarved(path string, st starvedState) error {
	data, err := json.Marshal(st)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return atomicfile.WriteFile(path, append(data, '\n'), 0o600)
}

package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
// story can start.

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
// board, kept between looks and across a restart, so that the same list is
// written once.
type starvedState struct {
	Thread string `json:"thread,omitempty"`
	List   string `json:"list,omitempty"`
}

// starvedPath is where a project's starvedState is kept: in its main
// checkout's .flai-cache, beside the ready-empty flag.
func starvedPath(repo *workitem.Repo) string {
	return filepath.Join(repo.CacheDir(), "starved.json")
}

// noteStarved judges, at a look, whether the board is starved, from the
// ready stories as readyStories answers them, their holds, and the room the
// in-progress limit leaves (below zero with no limit). A starved board gets
// flai serve's thread, or an entry on it when the list changed since it last
// wrote; a board that is not starved has the thread resolved. The backlog is
// read only when ready holds no story that can start and a lane is free, so a
// busy board pays nothing. Every failure is warned of and the look goes on.
func (l *launcher) noteStarved(repo *workitem.Repo, ready []readyStory, holds *workitem.Holds, room int) {
	path := starvedPath(repo)
	was, err := readStarved(path)
	lost := err != nil
	if lost {
		l.warn("starved board: what flai serve last wrote could not be read; its thread is looked for by its title", "path", path, "err", err)
	}
	if why := notStarved(ready, holds, room); why != "" {
		l.unstarve(repo, path, was, lost, why)
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
		l.unstarve(repo, path, was, lost, fmt.Sprintf("%s %s can be promoted from the backlog, so the board is no longer starved.", c.ID, oneLine(c.Title)))
		return
	}
	l.starve(repo, on, path, was, starvedList(s), room)
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
// on, with list, or adds list to the open one when it differs from what flai
// serve last wrote there, and keeps what it wrote. A thread the operator
// resolved while the board stayed starved is not opened again until the list
// changes.
func (l *launcher) starve(repo *workitem.Repo, on, path string, was starvedState, list string, room int) {
	open, err := starvedThreads(repo)
	if err != nil {
		l.warn("starved board: its thread could not be looked for, so the operator is not told; the next look tries again", "err", err)
		return
	}
	if len(open) == 0 {
		if was.Thread != "" && was.List == list {
			return
		}
		th, err := threads.New(repo, threads.NewOptions{Title: starvedTitle, On: on, Author: readyEmptyAuthor, Text: l.starvedText(list, room, false), Now: l.now()})
		if err != nil {
			l.warn("starved board: its thread could not be opened, so the operator is not told; the next look tries again", "on", on, "err", err)
			return
		}
		l.keepStarved(path, starvedState{Thread: th.ID, List: list})
		l.log("starved board thread opened", "thread", th.ID, "room", room)
		return
	}
	th := open[0]
	if th.ID == was.Thread && was.List == list {
		return
	}
	if th.ID != was.Thread && lastWrote(th, list) {
		l.keepStarved(path, starvedState{Thread: th.ID, List: list})
		return
	}
	if _, err := threads.Reply(repo, th.ID, readyEmptyAuthor, l.starvedText(list, room, true), l.now()); err != nil {
		l.warn("starved board: the changed list could not be written on its thread; the next look tries again", "thread", th.ID, "err", err)
		return
	}
	l.keepStarved(path, starvedState{Thread: th.ID, List: list})
	l.log("starved board thread told a changed list", "thread", th.ID, "room", room)
}

// unstarve resolves flai serve's open thread on a starved board, saying why,
// and forgets what it wrote there, so that the next starving opens a new one.
// It looks for the thread only when the state names one, or could not be
// read.
func (l *launcher) unstarve(repo *workitem.Repo, path string, was starvedState, lost bool, why string) {
	if was.Thread == "" && !lost {
		return
	}
	open, err := starvedThreads(repo)
	if err != nil {
		l.warn("starved board: its thread could not be looked for, so it is not resolved; the next look tries again", "err", err)
		return
	}
	var resolved []string
	for _, th := range open {
		if _, err := threads.Resolve(repo, th.ID, readyEmptyAuthor, why, l.now()); err != nil {
			l.warn("starved board: its thread could not be resolved; the next look tries again", "thread", th.ID, "resolved", resolved, "err", err)
			return
		}
		resolved = append(resolved, th.ID)
	}
	l.keepStarved(path, starvedState{})
	if len(resolved) > 0 {
		l.log("starved board thread resolved", "threads", resolved, "why", why)
	}
}

// keepStarved keeps st at path, warning when it cannot: the next look then
// finds the thread by its title.
func (l *launcher) keepStarved(path string, st starvedState) {
	if err := writeStarved(path, st); err != nil {
		l.warn("starved board: what flai serve wrote on its thread could not be kept", "path", path, "thread", st.Thread, "err", err)
	}
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
		b.WriteString(" flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.")
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

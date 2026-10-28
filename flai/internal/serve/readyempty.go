package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// readyEmptyAuthor is who a ready-empty event is by (ADR-0211 §1).
const readyEmptyAuthor = "flai serve"

// readyEmptyState is what flai serve keeps between looks, and across a
// restart, of the ready column against the in-progress room: whether the last
// look found it below, so that one draining is told once.
type readyEmptyState struct {
	Below bool `json:"below"`
}

// readyEmptyPath is where a project's readyEmptyState is kept: in its main
// checkout's .flai-cache, beside the event log it is the memory of.
func readyEmptyPath(repo *workitem.Repo) string {
	return filepath.Join(repo.CacheDir(), "ready-empty.json")
}

// readyDrained decides, from the ready stories, held ones counted, the room
// the in-progress limit leaves (the limit less the stories in progress, or
// below zero when there is no limit), and whether the last look found ready
// below the room, whether this look tells a ready-empty event and whether
// ready is below the room now (ADR-0211 §1). Ready is below when it is fewer
// than the room; with no limit, or no room, it is not. The event is told on
// the look that finds it below after one that did not, and not again until a
// look finds ready at the room or above.
func readyDrained(ready, room int, wasBelow bool) (emit, below bool) {
	below = room > 0 && ready < room
	return below && !wasBelow, below
}

// noteReadyEmpty appends a ready-empty event to the project's event log when
// the ready column, ready stories in it, has drained below room, the
// in-progress limit's room as readyStories counts it, since the last look,
// and keeps whether it is below in .flai-cache, so that neither a later look
// nor a restart of flai serve tells the same draining again. It keeps the
// flag before it appends, so that a flag it cannot keep never tells one
// draining at every look. A flag it cannot read or keep, a board it cannot
// read for the limit, or an event it cannot append, is warned of and the look
// goes on: the event log is news, not state (ADR-0211 §4).
func (l *launcher) noteReadyEmpty(repo *workitem.Repo, ready, room int) {
	path := readyEmptyPath(repo)
	was, err := readReadyEmpty(path)
	if err != nil {
		l.warn("ready-empty: the flag of the last look could not be read; taken as not below", "path", path, "err", err)
	}
	emit, below := readyDrained(ready, room, was.Below)
	if below == was.Below {
		return
	}
	if err := writeReadyEmpty(path, readyEmptyState{Below: below}); err != nil {
		l.warn("ready-empty: the flag could not be kept, so no event is told until it can be", "path", path, "below", below, "err", err)
		return
	}
	if !emit {
		return
	}
	board, err := repo.LoadBoard()
	if err != nil {
		l.warn("ready-empty: the board could not be read for the in-progress limit, so the event is not told", "ready", ready, "room", room, "err", err)
		return
	}
	limit := board.WIPLimits[workitem.InProgress]
	if err := inbox.AppendEvent(repo, inbox.ReadyEmptyEvent(ready, room, limit, readyEmptyAuthor)); err != nil {
		l.warn("ready-empty: the event could not be appended to the event log", "ready", ready, "room", room, "limit", limit, "err", err)
		return
	}
	l.log("ready column below the in-progress room", "ready", ready, "room", room, "limit", limit)
}

// readReadyEmpty reads the flag kept at path; none kept is not below.
func readReadyEmpty(path string) (readyEmptyState, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return readyEmptyState{}, nil
	}
	if err != nil {
		return readyEmptyState{}, err
	}
	var st readyEmptyState
	if err := json.Unmarshal(data, &st); err != nil {
		return readyEmptyState{}, fmt.Errorf("%s: %w; delete it, and flai serve keeps it afresh", path, err)
	}
	return st, nil
}

// writeReadyEmpty keeps st at path in one step.
func writeReadyEmpty(path string, st readyEmptyState) error {
	data, err := json.Marshal(st)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return atomicfile.WriteFile(path, append(data, '\n'), 0o600)
}

package mcpserver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A cursor is how far one agent has read the repository's changes. It is a
// read marker under .flai-cache, outside git: tooling never owns state
// (overview.md), and losing a cursor only repeats or skips the report of a
// change. Ready work and open threads are state and are always listed.
type cursor struct {
	Seen  string          `json:"seen"`            // TimeFormat; changes up to here were reported
	Keys  map[string]bool `json:"keys,omitempty"`  // changes stamped with that very second, already reported
	Order []string        `json:"order,omitempty"` // the pull order as last reported
	known bool            // a cursor file existed
}

// firstLook is how far back a new agent is told about.
const firstLook = 24 * time.Hour

// maxEvents is the most changes one look reports; the newest are kept and the
// rest are counted, not listed. It is a constant, not a setting: the bound
// exists so that a look always fits a tool result. The first look of a new
// agent on 2026-09-19 was 209 changes and 68 KB, about 330 characters each,
// and did not fit (S-0061); fifty is about 16 KB, which leaves room for
// threads and ready work, and is more than an agent acts on in one look.
const maxEvents = 50

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func (s *server) cursorPath() string {
	name := strings.Trim(unsafeName.ReplaceAllString(s.agent, "-"), "-.")
	if name == "" {
		name = "agent"
	}
	return filepath.Join(s.repo.CacheDir(), "mcp", name+".json")
}

func (s *server) loadCursor() cursor {
	var c cursor
	if data, err := os.ReadFile(s.cursorPath()); err == nil && json.Unmarshal(data, &c) == nil && c.Seen != "" {
		c.known = true
	}
	return c
}

func (c cursor) since(now time.Time) time.Time {
	if t, err := time.Parse(workitem.TimeFormat, c.Seen); err == nil {
		return t
	}
	return now.Add(-firstLook)
}

// saveCursor is best effort: a cursor that cannot be written costs a
// repeated report, not a failed call.
func (s *server) saveCursor(c cursor) {
	path := s.cursorPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	tmp := path + ".tmp"
	if os.WriteFile(tmp, data, 0o600) == nil {
		_ = os.Rename(tmp, path)
	}
}

// Event is one thing that changed since the agent last looked.
type Event struct {
	workitem.Change
	Summary string `json:"summary" jsonschema:"the change in words"`
}

// catchUp returns what others changed since the cursor, the newest maxEvents
// of them, and how many older ones it left out, and advances the cursor past
// all of them: what a look leaves out never comes back in a later one. A
// first look, with no cursor, is told of stories and epics only; a day of
// task transitions is history to an agent that has just arrived, not news.
func (s *server) catchUp() ([]Event, int, error) {
	items, err := s.repo.List(true)
	if err != nil {
		return nil, 0, err
	}
	board, err := s.repo.LoadBoard()
	if err != nil {
		return nil, 0, err
	}
	now := s.now().UTC().Truncate(time.Second)
	cur := s.loadCursor()
	events := []Event{}
	next := cursor{Seen: now.Format(workitem.TimeFormat), Keys: map[string]bool{}, Order: board.Order}
	for _, c := range workitem.Changes(items, cur.since(now), s.agent, cur.Keys) {
		// the cursor passes every change, reported or not
		if c.At == next.Seen {
			next.Keys[c.Key()] = true
		}
		if !cur.known && c.Type == workitem.Task {
			continue
		}
		events = append(events, Event{Change: c, Summary: describe(c)})
	}
	omitted := 0
	if len(events) > maxEvents {
		omitted = len(events) - maxEvents
		events = events[omitted:] // oldest first, so the newest are at the end
	}
	// keys already reported for a second that is still the current one
	if cur.Seen == next.Seen {
		for k := range cur.Keys {
			next.Keys[k] = true
		}
	}
	if cur.known && reordered(cur.Order, board.Order) {
		events = append(events, Event{
			Change:  workitem.Change{ID: "board", Type: "board", Title: "pull order", Kind: "reordered", At: next.Seen},
			Summary: "the pull order is now " + strings.Join(board.Order, ", "),
		})
	}
	s.saveCursor(next)
	return events, omitted, nil
}

// reordered reports whether the stories two pull orders share come in a
// different sequence. Stories entering and leaving the order are moves,
// reported as such (flai move appends a ready story and drops a started
// one); only a change of priority among them is news of its own.
func reordered(before, after []string) bool {
	in := func(list []string) map[string]bool {
		m := map[string]bool{}
		for _, id := range list {
			m[id] = true
		}
		return m
	}
	inBefore, inAfter := in(before), in(after)
	var a, b []string
	for _, id := range before {
		if inAfter[id] {
			a = append(a, id)
		}
	}
	for _, id := range after {
		if inBefore[id] {
			b = append(b, id)
		}
	}
	return strings.Join(a, ",") != strings.Join(b, ",")
}

func describe(c workitem.Change) string {
	who := ""
	if c.By != "" {
		who = " by " + c.By
	}
	switch c.Kind {
	case workitem.Moved:
		if c.Cause != "" {
			return c.ID + " " + c.Title + " was cancelled with " + c.Cause + who
		}
		return c.ID + " " + c.Title + " moved to " + c.To + who
	case workitem.WasBlocked:
		return c.ID + " " + c.Title + " was blocked: " + c.Reason
	case workitem.Unblocked:
		return c.ID + " " + c.Title + " was unblocked"
	}
	return c.ID + " " + c.Kind
}

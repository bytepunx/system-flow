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

// catchUp returns what others changed since the cursor and advances it.
func (s *server) catchUp() ([]Event, error) {
	items, err := s.repo.List(true)
	if err != nil {
		return nil, err
	}
	board, err := s.repo.LoadBoard()
	if err != nil {
		return nil, err
	}
	now := s.now().UTC().Truncate(time.Second)
	cur := s.loadCursor()
	events := []Event{}
	next := cursor{Seen: now.Format(workitem.TimeFormat), Keys: map[string]bool{}, Order: board.Order}
	for _, c := range workitem.Changes(items, cur.since(now), s.agent, cur.Keys) {
		events = append(events, Event{Change: c, Summary: describe(c)})
		if c.At == next.Seen {
			next.Keys[c.Key()] = true
		}
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
	return events, nil
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
		return c.ID + " " + c.Title + " moved to " + c.To + who
	case workitem.WasBlocked:
		return c.ID + " " + c.Title + " was blocked: " + c.Reason
	case workitem.Unblocked:
		return c.ID + " " + c.Title + " was unblocked"
	}
	return c.ID + " " + c.Kind
}

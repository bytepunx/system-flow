package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
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

// Edited is the kind of a change someone made to an item's own fields or
// body with flai edit; To names what of it changed.
const Edited = "edited"

// Overlapped is the kind of the change to an open story when a story whose
// changes its claim covers is accepted: Cause is the accepted story, To the
// paths, comma separated (S-0132). It is also the kind of the change to each
// of two stories in progress when a write grows the claim of one into the
// other's: Cause is the other story, To the paths the claim gained (I-0059).
const Overlapped = "overlapped"

// maxPaths is how many paths an overlap's summary names; To has them all.
const maxPaths = 10

// Event is one thing that changed since the agent last looked.
type Event struct {
	workitem.Change
	Summary string `json:"summary" jsonschema:"the change in words"`
	Project string `json:"project,omitempty" jsonschema:"the project it happened in, when the server serves more than one"`
}

// catchUp returns what others changed since the cursor, the newest maxEvents
// of them, and how many older ones it left out, and advances the cursor past
// all of them: what a look leaves out never comes back in a later one. A
// first look, with no cursor, is told of stories and epics only; a day of
// task transitions is history to an agent that has just arrived, not news.
func (s *server) catchUp(ctx context.Context) ([]Event, int, error) {
	items, err := s.listItems(ctx)
	if err != nil {
		return nil, 0, err
	}
	return s.catchUpWith(items)
}

// catchUpWith is catchUp over items already listed, archive included.
func (s *server) catchUpWith(items []*workitem.Item) ([]Event, int, error) {
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
	// Edits someone else made with flai edit or from the dashboard (S-0085).
	// They are not in the items' front matter, which is strict and shared
	// with older flai; flai edit notes them beside the cursors.
	since := cur.since(now).UTC().Truncate(time.Second)
	noticed := func(c workitem.Change, summary string) {
		at, err := time.Parse(workitem.TimeFormat, c.At)
		news := at.After(since) || (at.Equal(since) && !cur.Keys[c.Key()])
		if err != nil || !news {
			return
		}
		if c.At == next.Seen {
			next.Keys[c.Key()] = true
		}
		if !cur.known && c.Type == workitem.Task {
			return
		}
		events = append(events, Event{Change: c, Summary: summary})
	}
	for _, n := range itemedit.Notices(s.repo) {
		if n.By != s.agent {
			c := workitem.Change{ID: n.ID, Type: n.Type, Title: n.Title, Kind: Edited, To: strings.Join(n.Changed, ","), By: n.By, At: n.At}
			noticed(c, describe(c))
		}
	}
	// Open stories that an acceptance changed paths under (S-0132), and
	// stories in progress whose claims a write grew to overlap (I-0059): the
	// change is on the story told, caused by the accepted or the other story,
	// and names the paths. Whoever accepted or wrote, it is news to the agent
	// of the story told.
	for _, o := range itemedit.Overlaps(s.repo) {
		c := workitem.Change{ID: o.ID, Type: workitem.Story, Title: o.Title, Kind: Overlapped, To: strings.Join(o.Paths, ","), By: o.By, Cause: o.Cause(), At: o.At}
		summary := describe(c)
		if o.Accepted == "" {
			summary = describeGrown(c, o.Grew, o.Reached)
		}
		noticed(c, summary)
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].At < events[j].At })
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
	if cur.known && workitem.Reordered(cur.Order, board.Order) {
		events = append(events, Event{
			Change:  workitem.Change{ID: "board", Type: "board", Title: "pull order", Kind: "reordered", At: next.Seen},
			Summary: "the pull order is now " + strings.Join(board.Order, ", "),
		})
	}
	s.saveCursor(next)
	return events, omitted, nil
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
		if c.Follows != "" {
			return c.ID + " " + c.Title + " moved to " + c.To + who + ", following " + c.Follows
		}
		return c.ID + " " + c.Title + " moved to " + c.To + who
	case workitem.WasBlocked:
		return c.ID + " " + c.Title + " was blocked: " + c.Reason
	case workitem.Unblocked:
		return c.ID + " " + c.Title + " was unblocked"
	case Edited:
		return c.ID + " " + c.Title + " was edited" + who + ": " + strings.ReplaceAll(c.To, ",", ", ") + ". Read it again before you go on"
	case Overlapped:
		return c.Cause + " was accepted" + who + " and changed " + namePaths(c.To) + ", which " + c.ID + " " + c.Title + " claims. Run flai stream sync " + c.ID + " and the tests before you go on"
	}
	return c.ID + " " + c.Kind
}

// describeGrown says an overlapped change whose cause is not an acceptance:
// a write grew the claim of the story grew into that of reached, and both
// stories now claim the paths (I-0059).
func describeGrown(c workitem.Change, grew, reached string) string {
	who := ""
	if c.By != "" {
		who = ", written by " + c.By
	}
	return grew + "'s claim grew to overlap " + reached + "'s on " + namePaths(c.To) + who + ". Both stories claim them now: coordinate with " + c.Cause + "'s agent before " + c.ID + " " + c.Title + " changes them"
}

// namePaths names the comma-separated paths of an overlap, at most maxPaths.
func namePaths(to string) string {
	paths := strings.Split(to, ",")
	if len(paths) > maxPaths {
		return strings.Join(paths[:maxPaths], ", ") + fmt.Sprintf(" and %d more", len(paths)-maxPaths)
	}
	return strings.Join(paths, ", ")
}

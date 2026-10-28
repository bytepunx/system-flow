package inbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Planned is the kind of the event that a planner's run ended having written
// items beyond the item it planned (ADR-0211): ID, Type, and Title are the
// item planned, To the IDs written, comma separated, and By the planner.
const Planned = "planned"

// ReadyEmpty is the kind of the event that the ready column fell below the
// in-progress limit's room, once per draining (ADR-0211): ID and Type are
// "board", To is the ready count, Fields hold FieldRoom and FieldLimit, and By
// is flai serve.
const ReadyEmpty = "ready-empty"

// ReleaseDue is the kind of the event that an acceptance changed what flai
// release --evaluate answers (ADR-0211): ID, Type, and Title are the story
// accepted, Fields hold FieldBefore and FieldAfter, and By is the acceptor.
const ReleaseDue = "release-due"

// The keys of a logged event's Fields.
const (
	FieldRoom   = "room"   // ready-empty: the in-progress limit less the stories in progress
	FieldLimit  = "limit"  // ready-empty: the in-progress limit
	FieldBefore = "before" // release-due: what flai release --evaluate answered before the acceptance
	FieldAfter  = "after"  // release-due: what it answers after it
)

// A LoggedEvent is one line of the event log: something that happened in
// flai serve or flai accept and leaves no trace in a work item file
// (ADR-0211 §4). The log is under .flai-cache, outside git, as the cursors
// are: losing it loses news, never state.
type LoggedEvent struct {
	At     string            `json:"at"` // workitem.TimeFormat
	Kind   string            `json:"kind"`
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Title  string            `json:"title"`
	To     string            `json:"to,omitempty"`
	By     string            `json:"by,omitempty"`
	Cause  string            `json:"cause,omitempty"`
	Fields map[string]string `json:"fields,omitempty"` // the kind's own, under the Field keys
}

// Change is e as the inbox reports it, without its Fields.
func (e LoggedEvent) Change() workitem.Change {
	return workitem.Change{ID: e.ID, Type: e.Type, Title: e.Title, Kind: e.Kind, To: e.To, By: e.By, Cause: e.Cause, At: e.At}
}

// PlannedEvent is the planned event of a planner's run, by, on the item id of
// type typ and title, that wrote the items written.
func PlannedEvent(id, typ, title, by string, written []string) LoggedEvent {
	return LoggedEvent{Kind: Planned, ID: id, Type: typ, Title: title, To: strings.Join(written, ","), By: by}
}

// ReadyEmptyEvent is the ready-empty event of a ready column holding ready
// stories, below room, the in-progress limit less the stories in progress,
// told by by.
func ReadyEmptyEvent(ready, room, limit int, by string) LoggedEvent {
	return LoggedEvent{Kind: ReadyEmpty, ID: "board", Type: "board", Title: "ready column", To: strconv.Itoa(ready), By: by,
		Fields: map[string]string{FieldRoom: strconv.Itoa(room), FieldLimit: strconv.Itoa(limit)}}
}

// ReleaseDueEvent is the release-due event of by's acceptance of the story id
// with title, which changed what flai release --evaluate answers from before
// to after, each that answer in words.
func ReleaseDueEvent(id, title, by, before, after string) LoggedEvent {
	return LoggedEvent{Kind: ReleaseDue, ID: id, Type: workitem.Story, Title: title, By: by,
		Fields: map[string]string{FieldBefore: before, FieldAfter: after}}
}

// EventsPath is where the event log of a repository is: events.jsonl in its
// main checkout's .flai-cache, beside edits.jsonl and overlaps.jsonl.
func EventsPath(repo *workitem.Repo) string {
	return filepath.Join(repo.CacheDir(), "events.jsonl")
}

// AppendEvent appends e to the event log as one JSON line, stamped now when
// its At is empty: a writer leaves it so, so that a reader whose cursor stands
// at the second of the write still hears of it. Several writers may append at
// once: each line is one write to a file opened for appending. The log is
// appended to and never rewritten (ADR-0211 §4). A failure loses the news,
// not what happened; the writer decides whether to say so.
func AppendEvent(repo *workitem.Repo, e LoggedEvent) error {
	if e.At == "" {
		e.At = time.Now().UTC().Format(workitem.TimeFormat)
	}
	path := EventsPath(repo)
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("encode the %s event on %s for %s: %w", e.Kind, e.ID, path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create the folder of the event log %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open the event log %s to append the %s event on %s: %w", path, e.Kind, e.ID, err)
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		_ = f.Close()
		return fmt.Errorf("append the %s event on %s to %s: %w", e.Kind, e.ID, path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close the event log %s after the %s event on %s: %w", path, e.Kind, e.ID, err)
	}
	return nil
}

// LoggedEvents returns every event in the log, oldest first, passing over a
// line that does not parse or lacks its time, kind, or ID; none when there is
// no log.
func LoggedEvents(repo *workitem.Repo) []LoggedEvent {
	f, err := os.Open(EventsPath(repo))
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()
	var out []LoggedEvent
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		var e LoggedEvent
		if json.Unmarshal(sc.Bytes(), &e) == nil && e.At != "" && e.Kind != "" && e.ID != "" {
			out = append(out, e)
		}
	}
	return out
}

// letPass ends the summary of an event the orchestrator acts on and a story's
// agent does not.
const letPass = ". A story's agent may let it pass"

// describeLogged says an event of the log in words, with its kind's fields;
// without them it says what the change alone tells.
func describeLogged(c workitem.Change, fields map[string]string) string {
	who := ""
	if c.By != "" {
		who = " by " + c.By
	}
	switch c.Kind {
	case Planned:
		wrote := ""
		if c.To != "" {
			wrote = ": wrote " + strings.ReplaceAll(c.To, ",", ", ")
		}
		return c.ID + " " + c.Title + " was planned" + who + wrote + letPass
	case ReadyEmpty:
		holds := "the ready column fell below the in-progress room"
		if n, err := strconv.Atoi(c.To); err == nil {
			holds = "the ready column holds " + stories(n) + ", below the in-progress room"
		}
		if room := fields[FieldRoom]; room != "" {
			holds += " of " + room
		}
		if limit := fields[FieldLimit]; limit != "" {
			holds += " (limit " + limit + ")"
		}
		return holds + letPass
	case ReleaseDue:
		answer := ": run it to see what it answers now"
		if after := fields[FieldAfter]; after != "" {
			answer = ": it answered " + orNothing(fields[FieldBefore]) + " and now answers " + after
		}
		return c.ID + " " + c.Title + "'s acceptance" + who + " changed what flai release --evaluate answers" + answer + letPass
	}
	return c.ID + " " + c.Kind
}

// stories counts n stories in words.
func stories(n int) string {
	if n == 1 {
		return "1 story"
	}
	return strconv.Itoa(n) + " stories"
}

// orNothing is s, or "nothing" when it is empty.
func orNothing(s string) string {
	if s == "" {
		return "nothing"
	}
	return s
}

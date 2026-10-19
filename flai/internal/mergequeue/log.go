package mergequeue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The merge queue's event log (S-0404): one line for each event the queue
// acts on, so that flai stats can measure the queue. It is
// .flai-cache/merge-queue/events.jsonl, beside the state, append-only and
// never trimmed. Update, the state's only writer, derives each write's events
// by comparing the state before and after it, and appends them; the log
// measures the queue and never stops a write.

// The events of the log, as Event.Event names them.
const (
	// EventEnqueued is a story enqueued.
	EventEnqueued = "enqueued"
	// EventBatchBuilt is a batch built: its stories rebased, those whose
	// rebase stopped left out, and its verification begun.
	EventBatchBuilt = "batch_built"
	// EventVerified is a batch's verification finished, passed or not.
	EventVerified = "verified"
	// EventBisected is the bisection of a batch that failed begun.
	EventBisected = "bisected"
	// EventLanded is a story landed by the queue.
	EventLanded = "landed"
	// EventSentBack is the story a failed batch was bisected to, sent back.
	EventSentBack = "sent_back"
	// EventDequeued is a story that left the queue other than by landing or
	// being sent back.
	EventDequeued = "dequeued"
)

// eventKinds are the events of the log, in the order a write's lines are.
var eventKinds = []string{EventEnqueued, EventBatchBuilt, EventVerified, EventBisected, EventLanded, EventSentBack, EventDequeued}

// ReasonRemoved is the reason of a story dequeued with no outcome recorded
// for it, as one taken out with flai queue remove.
const ReasonRemoved = "removed"

// Event is one line of the merge queue's event log.
type Event struct {
	At      time.Time `json:"at"`                 // when it happened, UTC to the second
	Event   string    `json:"event"`              // one of the Event constants
	Batch   string    `json:"batch,omitempty"`    // the batch it is of, when one is
	Stories []string  `json:"stories"`            // the stories it is on
	By      string    `json:"by,omitempty"`       // who enqueued the story, on an enqueued line
	Passed  *bool     `json:"passed,omitempty"`   // whether the batch passed, on a verified line
	LeftOut []string  `json:"left_out,omitempty"` // the stories left out, on a batch_built line
	Reason  string    `json:"reason,omitempty"`   // why the story left, on a dequeued line
}

// EventsPath is where the merge queue's event log is in the cache folder
// cacheDir, beside StatePath's state.
func EventsPath(cacheDir string) string {
	return filepath.Join(cacheDir, "merge-queue", "events.jsonl")
}

// Events are the events of a write that took the queue's state from before
// to after at now, in the order of eventKinds and, within an event, in the
// state's order.
func Events(before, after State, now time.Time) []Event {
	var out []Event
	for _, e := range after.Enqueued {
		if _, ok := before.EnqueuedAs(e.Story); !ok {
			out = append(out, Event{At: e.At, Event: EventEnqueued, Stories: []string{e.Story}, By: e.By})
		}
	}
	out = append(out, batchEvents(before.Batch, after.Batch, now)...)
	recorded := map[string]Outcome{}
	var landed, sentBack []Event
	for _, o := range after.Outcomes {
		if prev, ok := lastOutcome(before, o.Story); ok && prev.Outcome == o.Outcome && prev.At.Equal(o.At) {
			continue
		}
		recorded[workitem.CanonicalID(o.Story)] = o
		ev := Event{At: o.At, Stories: []string{o.Story}, Batch: batchHolding(o.Story, before, after)}
		switch o.Outcome {
		case OutcomeLanded:
			ev.Event = EventLanded
			landed = append(landed, ev)
		case OutcomeFailing:
			ev.Event = EventSentBack
			sentBack = append(sentBack, ev)
		}
	}
	out = append(append(out, landed...), sentBack...)
	for _, e := range before.Enqueued {
		if _, ok := after.EnqueuedAs(e.Story); ok {
			continue
		}
		ev := Event{At: now, Event: EventDequeued, Stories: []string{e.Story}, Reason: ReasonRemoved}
		if o, ok := recorded[workitem.CanonicalID(e.Story)]; ok {
			if o.Outcome == OutcomeLanded || o.Outcome == OutcomeFailing {
				continue
			}
			ev.At, ev.Reason = o.At, o.Outcome
		}
		out = append(out, ev)
	}
	return out
}

// batchEvents are the events of the batch after, set against prev, the batch
// before the write when it is the same batch: built once it leaves building
// and rebasing, verified once it has a result, bisected once its bisection
// begins.
func batchEvents(prev, b *Batch, now time.Time) []Event {
	if b == nil {
		return nil
	}
	if prev != nil && prev.ID != b.ID {
		prev = nil
	}
	stories := func() []string { return append([]string{}, b.Stories...) }
	var out []Event
	if !building(b.Step) && (prev == nil || building(prev.Step)) {
		var left []string
		for _, l := range b.LeftOut {
			left = append(left, l.Story)
		}
		out = append(out, Event{At: now, Event: EventBatchBuilt, Batch: b.ID, Stories: stories(), LeftOut: left})
	}
	if b.Result != nil && (prev == nil || prev.Result == nil) {
		passed, at := b.Result.Passed, b.Result.FinishedAt
		if at.IsZero() {
			at = now
		}
		out = append(out, Event{At: at, Event: EventVerified, Batch: b.ID, Stories: stories(), Passed: &passed})
	}
	if b.Bisect != nil && (prev == nil || prev.Bisect == nil) {
		out = append(out, Event{At: now, Event: EventBisected, Batch: b.ID, Stories: stories()})
	}
	return out
}

// building reports whether a batch at step is still being built: at
// StepBuilding, rebasing a story, or at no step yet.
func building(step string) bool {
	return step == "" || step == StepBuilding || step == StepRebasing || strings.HasPrefix(step, StepRebasing+" ")
}

// lastOutcome is story's last outcome in st, and whether it has one.
func lastOutcome(st State, story string) (Outcome, bool) {
	for _, o := range st.Outcomes {
		if sameStory(o.Story, story) {
			return o, true
		}
	}
	return Outcome{}, false
}

// batchHolding is the ID of the batch of the first of states that holds
// story, "" when none does.
func batchHolding(story string, states ...State) string {
	for _, st := range states {
		if st.HoldsStory(story) {
			return st.Batch.ID
		}
	}
	return ""
}

// logEvents appends to the event log in cacheDir the events of the write
// that took the state from before, as it was encoded, to after at now.
func logEvents(cacheDir string, before []byte, after State, now time.Time) error {
	var prev State
	if err := json.Unmarshal(before, &prev); err != nil {
		return fmt.Errorf("decode the merge queue's state as it was before the write, to log its events: %w", err)
	}
	return AppendEvents(cacheDir, Events(prev, after, now))
}

// AppendEvents appends events to the merge queue's event log in cacheDir in
// one write, a line each, their times in UTC to the second. A line a write
// before cut off is dropped first, so that it is not joined to the first of
// these.
func AppendEvents(cacheDir string, events []Event) error {
	if len(events) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, e := range events {
		if !slices.Contains(eventKinds, e.Event) {
			return fmt.Errorf("log a merge queue event %q on %s: an event is one of %s", e.Event, orNone(e.Stories), strings.Join(eventKinds, ", "))
		}
		e.At = e.At.UTC().Truncate(time.Second)
		e.Stories = append([]string{}, e.Stories...)
		line, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("encode the merge queue's %s event on %s for its event log: %w", e.Event, orNone(e.Stories), err)
		}
		buf.Write(append(line, '\n'))
	}
	path := EventsPath(cacheDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("make the folder %s for the merge queue's event log: %w; check that the project's .flai-cache is writable", filepath.Dir(path), err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("open the merge queue's event log %s: %w; check that the project's .flai-cache is writable", path, err)
	}
	if err := dropCutOff(f); err != nil {
		_ = f.Close()
		return fmt.Errorf("drop the line cut off at the end of the merge queue's event log %s: %w; remove that line by hand", path, err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		_ = f.Close()
		return fmt.Errorf("append to the merge queue's event log %s: %w; check that the disk has room", path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close the merge queue's event log %s: %w", path, err)
	}
	return nil
}

// dropCutOff truncates f after its last newline when it does not end with
// one: a line a write was cut off in. It reads f whole only then.
func dropCutOff(f *os.File) error {
	info, err := f.Stat()
	if err != nil || info.Size() == 0 {
		return err
	}
	last := make([]byte, 1)
	if _, err := f.ReadAt(last, info.Size()-1); err != nil {
		return err
	}
	if last[0] == '\n' {
		return nil
	}
	data, err := io.ReadAll(io.NewSectionReader(f, 0, info.Size()))
	if err != nil {
		return err
	}
	return f.Truncate(int64(bytes.LastIndexByte(data, '\n') + 1))
}

// ReadEvents returns every line of the merge queue's event log in cacheDir,
// oldest first: none when there is no log yet, and an error when it cannot be
// read or a line cannot be parsed. Blank lines are skipped, and so is a last
// line that does not parse and no newline ends: a write cut off.
func ReadEvents(cacheDir string) ([]Event, error) {
	path := EventsPath(cacheDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read the merge queue's event log %s: %w; check that the project's .flai-cache is readable", path, err)
	}
	lines := strings.Split(string(data), "\n")
	var out []Event
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			if i == len(lines)-1 {
				break
			}
			return nil, fmt.Errorf("parse line %d of the merge queue's event log %s: %w; correct or remove that line", i+1, path, err)
		}
		out = append(out, e)
	}
	return out, nil
}

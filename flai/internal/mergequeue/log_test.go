package mergequeue

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// copyState is a deep copy of st, as Update keeps the state before a write.
func copyState(t *testing.T, st State) State {
	t.Helper()
	data, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var c State
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// eventLog reads the event log in dir, failing the test on an error.
func eventLog(t *testing.T, dir string) []Event {
	t.Helper()
	got, err := ReadEvents(dir)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

// S-0404: each write's events come from the state before and after it, in
// the order of their kinds, and a write that changes nothing the log measures
// writes none.
func TestEventsOfAWrite(t *testing.T) {
	t1, t2, now := t0.Add(time.Minute), t0.Add(2*time.Minute), t0.Add(time.Hour)
	yes, no := true, false
	enqueued := func(stories ...string) []Enqueued {
		out := []Enqueued{}
		for _, s := range stories {
			out = append(out, Enqueued{Story: s, By: "alex", At: t0})
		}
		return out
	}
	batch := func(step string, stories ...string) *Batch {
		return &Batch{ID: "B-1", Stories: stories, Step: step}
	}
	failed := func() *Batch {
		b := batch(StepFailed, "S-0001", "S-0002")
		b.Result = &Result{Passed: false, StoppedAt: "unit", FinishedAt: t1}
		return b
	}
	cases := []struct {
		name   string
		before State
		change func(*State)
		want   []Event
	}{
		{
			name:   "a story enqueued",
			before: State{Enqueued: enqueued("S-0001")},
			change: func(s *State) { s.Enqueue("S-0002", "orchestrator", t1) },
			want:   []Event{{At: t1, Event: EventEnqueued, Stories: []string{"S-0002"}, By: "orchestrator"}},
		},
		{
			name:   "a hold",
			before: State{Enqueued: enqueued("S-0001")},
			change: func(s *State) { s.Hold("alex", t1) },
		},
		{
			name:   "a batch from building to rebasing",
			before: State{Enqueued: enqueued("S-0001"), Batch: batch(StepBuilding, "S-0001")},
			change: func(s *State) { s.Batch.Step = StepRebasing + " S-0001" },
		},
		{
			name:   "a new batch at building",
			before: State{Enqueued: enqueued("S-0001")},
			change: func(s *State) { s.Batch = batch(StepBuilding, "S-0001") },
		},
		{
			name:   "a batch from rebasing to verifying",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: batch(StepRebasing+" S-0002", "S-0001")},
			change: func(s *State) {
				s.Batch.LeftOut = []LeftOut{{Story: "S-0002", Paths: []string{"a.md"}}}
				s.Batch.Step = StepVerifying
			},
			want: []Event{{At: now, Event: EventBatchBuilt, Batch: "B-1", Stories: []string{"S-0001"}, LeftOut: []string{"S-0002"}}},
		},
		{
			name:   "a batch at verifying again",
			before: State{Enqueued: enqueued("S-0001"), Batch: batch(StepVerifying, "S-0001")},
			change: func(s *State) { s.Batch.Step = StepVerifying + " unit" },
		},
		{
			name:   "a batch replaced by one at verifying",
			before: State{Enqueued: enqueued("S-0001"), Batch: batch(StepVerifying, "S-0001")},
			change: func(s *State) {
				s.Batch = &Batch{ID: "B-2", Stories: []string{"S-0001"}, Step: StepVerifying}
			},
			want: []Event{{At: now, Event: EventBatchBuilt, Batch: "B-2", Stories: []string{"S-0001"}}},
		},
		{
			name:   "a batch verified",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: batch(StepVerifying+" unit", "S-0001", "S-0002")},
			change: func(s *State) {
				s.Batch.Result, s.Batch.Step = &Result{Passed: false, FinishedAt: t1}, StepFailed
			},
			want: []Event{{At: t1, Event: EventVerified, Batch: "B-1", Stories: []string{"S-0001", "S-0002"}, Passed: &no}},
		},
		{
			name:   "a batch verified with no finish time",
			before: State{Enqueued: enqueued("S-0001"), Batch: batch(StepVerifying, "S-0001")},
			change: func(s *State) { s.Batch.Result, s.Batch.Step = &Result{Passed: true}, StepVerified },
			want:   []Event{{At: now, Event: EventVerified, Batch: "B-1", Stories: []string{"S-0001"}, Passed: &yes}},
		},
		{
			name:   "a batch bisected",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: failed()},
			change: func(s *State) {
				s.Batch.Bisect, s.Batch.Step = &Bisect{Steps: []BisectStep{}, Bad: 2}, StepBisecting+" 1 of 2"
			},
			want: []Event{{At: now, Event: EventBisected, Batch: "B-1", Stories: []string{"S-0001", "S-0002"}}},
		},
		{
			name: "a bisection step",
			before: func() State {
				b := failed()
				b.Bisect = &Bisect{Steps: []BisectStep{}, Bad: 2}
				return State{Enqueued: enqueued("S-0001", "S-0002"), Batch: b}
			}(),
			change: func(s *State) { s.Batch.Bisect.Good = 1 },
		},
		{
			name:   "a story landed",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: batch(StepLanding+" S-0001", "S-0001", "S-0002")},
			change: func(s *State) {
				s.Dequeue("S-0001")
				s.RecordOutcome(Outcome{Story: "S-0001", Outcome: OutcomeLanded, At: t2, Detail: "landed by merge queue batch B-1 as alex"})
			},
			want: []Event{{At: t2, Event: EventLanded, Batch: "B-1", Stories: []string{"S-0001"}}},
		},
		{
			name: "a landed story's landing recorded again",
			before: State{Enqueued: enqueued("S-0002"), Batch: batch(StepLanding+" S-0001", "S-0001", "S-0002"),
				Outcomes: []Outcome{{Story: "S-0001", Outcome: OutcomeLanded, At: t2}}},
			change: func(s *State) {
				s.RecordOutcome(Outcome{Story: "S-0001", Outcome: OutcomeLanded, At: t2})
				s.Batch.Landing = &Landing{Stories: []string{"S-0001"}, Landed: []string{"S-0001"}, Head: "abc"}
			},
		},
		{
			name:   "a failing story sent back",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: failed()},
			change: func(s *State) {
				s.Dequeue("s-2")
				s.RecordOutcome(Outcome{Story: "S-0002", Outcome: OutcomeFailing, At: t2})
				s.Batch.Bisect = &Bisect{Failing: "S-0002", SentBack: true}
			},
			want: []Event{
				{At: now, Event: EventBisected, Batch: "B-1", Stories: []string{"S-0001", "S-0002"}},
				{At: t2, Event: EventSentBack, Batch: "B-1", Stories: []string{"S-0002"}},
			},
		},
		{
			name:   "a story removed",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: batch(StepVerifying, "S-0001")},
			change: func(s *State) {
				s.Dequeue("S-0001")
				s.DropBatchHolding("S-0001", "S-0001 was removed from the queue by alex", "alex", t1)
			},
			want: []Event{{At: now, Event: EventDequeued, Stories: []string{"S-0001"}, Reason: ReasonRemoved}},
		},
		{
			name:   "a story accepted by hand",
			before: State{Enqueued: enqueued("S-0001")},
			change: func(s *State) {
				s.Dequeue("S-0001")
				s.RecordOutcome(Outcome{Story: "S-0001", Outcome: OutcomeAcceptedByHand, At: t1})
			},
			want: []Event{{At: t1, Event: EventDequeued, Stories: []string{"S-0001"}, Reason: OutcomeAcceptedByHand}},
		},
		{
			name:   "a story left out stays enqueued",
			before: State{Enqueued: enqueued("S-0001"), Batch: batch(StepRebasing+" S-0001", "S-0001")},
			change: func(s *State) {
				s.Batch.Stories = nil
				s.RecordOutcome(Outcome{Story: "S-0001", Outcome: OutcomeLeftOut, At: t1})
			},
		},
		{
			name:   "several events in one write, in the order of their kinds",
			before: State{Enqueued: enqueued("S-0001", "S-0002"), Batch: batch(StepRebasing+" S-0002", "S-0001", "S-0002")},
			change: func(s *State) {
				s.Dequeue("S-0002")
				s.Enqueue("S-0003", "alex", t1)
				s.Batch.Step = StepVerifying
			},
			want: []Event{
				{At: t1, Event: EventEnqueued, Stories: []string{"S-0003"}, By: "alex"},
				{At: now, Event: EventBatchBuilt, Batch: "B-1", Stories: []string{"S-0001", "S-0002"}},
				{At: now, Event: EventDequeued, Stories: []string{"S-0002"}, Reason: ReasonRemoved},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			after := copyState(t, c.before)
			c.change(&after)
			if got := Events(c.before, after, now); !reflect.DeepEqual(got, c.want) {
				t.Errorf("events:\n%+v\nwant:\n%+v", got, c.want)
			}
		})
	}
}

// Each event is one line, its time in UTC to the second, read back as it was
// appended, oldest first; an event of no known kind is not appended.
func TestAppendAndReadEvents(t *testing.T) {
	dir := filepath.Join(t.TempDir(), ".flai-cache")
	at := time.Date(2026, 10, 19, 10, 30, 15, 500, time.FixedZone("CEST", 2*3600))
	passed := true
	if err := AppendEvents(dir, []Event{
		{At: at, Event: EventEnqueued, Stories: []string{"S-0001"}, By: "alex"},
		{At: at, Event: EventVerified, Batch: "B-1", Stories: []string{"S-0001"}, Passed: &passed},
	}); err != nil {
		t.Fatal(err)
	}
	if err := AppendEvents(dir, []Event{{At: at, Event: EventDequeued, Stories: []string{"S-0002"}, Reason: ReasonRemoved}}); err != nil {
		t.Fatal(err)
	}
	if err := AppendEvents(dir, nil); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(EventsPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	want := `{"at":"2026-10-19T08:30:15Z","event":"enqueued","stories":["S-0001"],"by":"alex"}` + "\n" +
		`{"at":"2026-10-19T08:30:15Z","event":"verified","batch":"B-1","stories":["S-0001"],"passed":true}` + "\n" +
		`{"at":"2026-10-19T08:30:15Z","event":"dequeued","stories":["S-0002"],"reason":"removed"}` + "\n"
	if string(data) != want {
		t.Errorf("log:\n%s\nwant:\n%s", data, want)
	}
	utc := time.Date(2026, 10, 19, 8, 30, 15, 0, time.UTC)
	wantRead := []Event{
		{At: utc, Event: EventEnqueued, Stories: []string{"S-0001"}, By: "alex"},
		{At: utc, Event: EventVerified, Batch: "B-1", Stories: []string{"S-0001"}, Passed: &passed},
		{At: utc, Event: EventDequeued, Stories: []string{"S-0002"}, Reason: ReasonRemoved},
	}
	if got := eventLog(t, dir); !reflect.DeepEqual(got, wantRead) {
		t.Errorf("read: %+v, want %+v", got, wantRead)
	}
	if err := AppendEvents(dir, []Event{{At: at, Event: "merged", Stories: []string{"S-0001"}}}); err == nil {
		t.Error("an event of no known kind was appended")
	}
	if got := eventLog(t, dir); len(got) != 3 {
		t.Errorf("a refused append wrote lines: %+v", got)
	}
}

// The log is events.jsonl beside the state; a missing log is no events.
func TestAMissingEventLogReadsAsNone(t *testing.T) {
	if got, want := EventsPath("/p/.flai-cache"), filepath.Join("/p/.flai-cache", "merge-queue", "events.jsonl"); got != want {
		t.Errorf("events path %s, want %s", got, want)
	}
	if got, err := ReadEvents(t.TempDir()); err != nil || got != nil {
		t.Errorf("missing log: %+v %v", got, err)
	}
}

// A last line a write cut off is skipped, and the next append drops it
// rather than joining it; a bad line before the last is an error naming it,
// and blank lines are skipped.
func TestACutOffLastLineIsSkippedAndABadLineIsAnError(t *testing.T) {
	dir := t.TempDir()
	good := `{"at":"2026-10-19T08:30:15Z","event":"enqueued","stories":["S-0001"],"by":"alex"}`
	if err := os.MkdirAll(filepath.Dir(EventsPath(dir)), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(EventsPath(dir), []byte(good+"\n\n"+`{"at":"2026-10-19T0`), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := eventLog(t, dir); len(got) != 1 {
		t.Errorf("cut-off line not skipped: %+v", got)
	}
	if err := AppendEvents(dir, []Event{{At: t0, Event: EventDequeued, Stories: []string{"S-0001"}, Reason: ReasonRemoved}}); err != nil {
		t.Fatal(err)
	}
	if got := eventLog(t, dir); len(got) != 2 || got[1].Event != EventDequeued {
		t.Errorf("the append after a cut-off line: %+v", got)
	}
	if err := os.WriteFile(EventsPath(dir), []byte(good+"\nnot json\n"+good+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadEvents(dir)
	if err == nil || got != nil || !strings.Contains(err.Error(), "line 2 of the merge queue's event log") {
		t.Errorf("bad line: %+v %v", got, err)
	}
}

// Update appends each write's events to the log, none for a write that
// changes nothing; Remove logs its story dequeued as removed, and
// AcceptedByHand as accepted by hand.
func TestUpdateLogsEachWritesEvents(t *testing.T) {
	repo := project(t, "")
	dir := repo.CacheDir()
	for _, s := range []string{"S-0001", "S-0002"} {
		if _, err := Update(dir, func(st *State) error { st.Enqueue(s, "alex", t0); return nil }); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Update(dir, func(st *State) error { st.Enqueue("S-0001", "alex", t0); return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err := Remove(repo, "S-0001", "alex", t0); err != nil {
		t.Fatal(err)
	}
	if _, err := AcceptedByHand(dir, "S-0002", "alex", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	got := eventLog(t, dir)
	var lines []string
	for _, e := range got {
		lines = append(lines, e.Event+" "+strings.Join(e.Stories, ",")+" "+e.Reason)
	}
	want := []string{"enqueued S-0001 ", "enqueued S-0002 ", "dequeued S-0001 removed", "dequeued S-0002 accepted_by_hand"}
	if !reflect.DeepEqual(lines, want) {
		t.Errorf("log %q, want %q", lines, want)
	}
	if !got[3].At.Equal(t0.Add(time.Minute)) {
		t.Errorf("accepted by hand at %s, want the outcome's time", got[3].At)
	}
}

// An event log that cannot be appended to is logged at warn, and the write
// to the state stands.
func TestAnEventLogNotWrittenWarnsAndTheWriteStands(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(EventsPath(dir), 0o755); err != nil {
		t.Fatal(err)
	}
	var logged bytes.Buffer
	was := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logged, nil)))
	t.Cleanup(func() { slog.SetDefault(was) })
	if _, err := Update(dir, func(st *State) error { st.Enqueue("S-0001", "alex", t0); return nil }); err != nil {
		t.Fatalf("the write failed on the event log: %v", err)
	}
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := st.EnqueuedAs("S-0001"); !ok {
		t.Errorf("the write did not stand: %+v", st)
	}
	for _, want := range []string{"level=WARN", "component=mergequeue", "detail=", "err="} {
		if !strings.Contains(logged.String(), want) {
			t.Errorf("the failed append's warning lacks %s: %s", want, logged.String())
		}
	}
}

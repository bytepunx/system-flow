package mergequeue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// writeState writes data as the queue's state in cacheDir.
func writeState(t *testing.T, cacheDir, data string) {
	t.Helper()
	path := StatePath(cacheDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0173: the state is one file in the main checkout's .flai-cache.
func TestTheStateIsInTheCachesMergeQueueFolder(t *testing.T) {
	if got, want := StatePath("/p/.flai-cache"), filepath.Join("/p/.flai-cache", "merge-queue", "state.json"); got != want {
		t.Errorf("state path %s, want %s", got, want)
	}
}

// With no file the queue has nothing enqueued, no batch, and no outcomes,
// and that is no error.
func TestAMissingStateReadsAsEmpty(t *testing.T) {
	st, err := ReadState(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if st.Batch != nil || st.Enqueued == nil || len(st.Enqueued) != 0 || st.Outcomes == nil || len(st.Outcomes) != 0 {
		t.Errorf("state %+v, want empty lists and no batch", st)
	}
}

// The state holds the stories enqueued with who and when, the batch and its
// step, and each story's last outcome.
func TestTheStateReadsEnqueuedTheBatchAndTheOutcomes(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{
  "enqueued": [{"story": "S-0002", "by": "orchestrator", "at": "2026-10-19T08:00:00Z"}, {"story": "S-0003", "by": "alex", "at": "2026-10-19T08:05:00Z"}],
  "batch": {"stories": ["S-0002", "S-0003"], "step": "verifying"},
  "outcomes": [{"story": "S-0001", "outcome": "landed", "at": "2026-10-19T07:00:00Z", "detail": "merged with S-0004"}]
}`)
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 10, 19, 8, 5, 0, 0, time.UTC)
	if len(st.Enqueued) != 2 || st.Enqueued[1] != (Enqueued{Story: "S-0003", By: "alex", At: at}) {
		t.Errorf("enqueued %+v", st.Enqueued)
	}
	if st.Batch == nil || strings.Join(st.Batch.Stories, ",") != "S-0002,S-0003" || st.Batch.Step != "verifying" {
		t.Errorf("batch %+v", st.Batch)
	}
	want := Outcome{Story: "S-0001", Outcome: "landed", At: time.Date(2026, 10, 19, 7, 0, 0, 0, time.UTC), Detail: "merged with S-0004"}
	if len(st.Outcomes) != 1 || st.Outcomes[0] != want {
		t.Errorf("outcomes %+v, want %+v", st.Outcomes, want)
	}
}

// Lists written as null, or left out, read as empty lists, so that every
// surface shows [] rather than null.
func TestTheStatesListsAreNeverNull(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"enqueued": null}`)
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st.Enqueued == nil || st.Outcomes == nil || st.Batch != nil {
		t.Errorf("state %+v, want empty lists and no batch", st)
	}
}

// A malformed state is an error naming the file and what to do.
func TestAMalformedStateIsAnErrorNamingTheFile(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"enqueued": [`)
	_, err := ReadState(dir)
	if err == nil {
		t.Fatal("a malformed state read without error")
	}
	for _, says := range []string{StatePath(dir), "merge queue's state", "delete it"} {
		if !strings.Contains(err.Error(), says) {
			t.Errorf("error %q does not say %q", err, says)
		}
	}
}

// Update writes what fn leaves, and ReadState reads it back with the batch's
// branch, worktree, result, and drop, and the hold, under snake_case names.
func TestUpdateWritesTheStateReadStateReads(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.UTC)
	batch := &Batch{ID: "B-1", Branch: "merge-queue/B-1", Worktree: "/p/.flai-cache/merge-queue/B-1", Stories: []string{"S-0002"}, StartedAt: at, Step: "verify"}
	batch.Result = &Result{Passed: false, StoppedAt: "integration", Summary: "TestX failed", FinishedAt: at}
	batch.Dropped = &Dropped{Reason: "S-0002 was accepted by hand by alex", By: "alex", At: at}
	want := State{Enqueued: []Enqueued{{Story: "S-0002", By: "alex", At: at}}, Held: true, HeldBy: "alex", HeldAt: at, Batch: batch, Outcomes: []Outcome{}}
	got, err := Update(dir, func(s *State) error { *s = want; return nil })
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(StatePath(dir))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{`"held": true`, `"held_by"`, `"held_at"`, `"branch"`, `"worktree"`, `"started_at"`, `"stopped_at"`, `"finished_at"`, `"dropped"`, `"reason"`} {
		if !strings.Contains(string(data), name) {
			t.Errorf("the state file has no %s:\n%s", name, data)
		}
	}
	read, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, st := range []State{got, read} {
		if !st.Held || st.HeldBy != "alex" || !st.HeldAt.Equal(at) || st.Batch == nil || *st.Batch.Result != *want.Batch.Result || *st.Batch.Dropped != *want.Batch.Dropped || st.Batch.Branch != "merge-queue/B-1" || st.Batch.Worktree == "" || !st.Batch.StartedAt.Equal(at) {
			t.Errorf("state %+v, batch %+v", st, st.Batch)
		}
	}
}

// A state with no hold and no result leaves those fields out of its file.
func TestAnUnheldStateLeavesTheHoldOut(t *testing.T) {
	dir := t.TempDir()
	if _, err := Update(dir, func(s *State) error { s.Enqueue("S-0001", "alex", time.Now()); return nil }); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(StatePath(dir))
	for _, name := range []string{"held", "batch", "dropped"} {
		if strings.Contains(string(data), name) {
			t.Errorf("the state file has %s:\n%s", name, data)
		}
	}
}

// An error from fn writes nothing, and a state left as it was is not
// written: no file appears for either.
func TestUpdateWritesNothingOnAnErrorOrNoChange(t *testing.T) {
	dir := t.TempDir()
	refused := errors.New("refused")
	if _, err := Update(dir, func(s *State) error { s.Enqueue("S-0001", "alex", time.Now()); return refused }); !errors.Is(err, refused) {
		t.Errorf("err %v, want fn's", err)
	}
	if _, err := Update(dir, func(s *State) error { s.Dequeue("S-0001"); return nil }); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(StatePath(dir)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("a state file was written: %v", err)
	}
}

// Concurrent writers each take the lock in turn, so that none loses
// another's write.
func TestConcurrentUpdatesLoseNoWrite(t *testing.T) {
	dir := t.TempDir()
	const writers = 16
	var wg sync.WaitGroup
	errs := make(chan error, writers)
	for i := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := Update(dir, func(s *State) error {
				s.Enqueue(fmt.Sprintf("S-%04d", i+1), "alex", time.Now())
				return nil
			})
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Enqueued) != writers {
		t.Errorf("%d enqueued, want %d: %+v", len(st.Enqueued), writers, st.Enqueued)
	}
}

// A malformed state is not overwritten: Update stops with ReadState's error.
func TestUpdateOfAMalformedStateIsAnError(t *testing.T) {
	dir := t.TempDir()
	writeState(t, dir, `{"enqueued": [`)
	if _, err := Update(dir, func(*State) error { return nil }); err == nil || !strings.Contains(err.Error(), "delete it") {
		t.Errorf("err %v, want ReadState's", err)
	}
}

// Enqueue keeps the first enqueuer, Dequeue reports a story not enqueued,
// and both match an ID however it is written.
func TestEnqueueAndDequeue(t *testing.T) {
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.FixedZone("x", 3600))
	var s State
	if !s.Enqueue("S-0001", "alex", at) || s.Enqueue("s-1", "orchestrator", at.Add(time.Hour)) {
		t.Fatal("enqueue twice")
	}
	if e, ok := s.EnqueuedAs("S-1"); !ok || e.By != "alex" || e.At.Location() != time.UTC {
		t.Errorf("enqueued as %+v %v", e, ok)
	}
	if s.Dequeue("S-0002") || !s.Dequeue("s-0001") || len(s.Enqueued) != 0 {
		t.Errorf("dequeue: %+v", s.Enqueued)
	}
}

// S-0403: EnqueueEntry records the orchestrator's verified commit and
// evidence with the story, in UTC, and a second entry for it keeps the
// first's place and what it recorded; both are written to the state file.
func TestEnqueueEntryRecordsTheOrchestratorsEvidence(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.FixedZone("x", 3600))
	first := Enqueued{Story: "S-0001", By: "orchestrator", At: at, Verified: "4f1c2a9", Evidence: "Verdict: pass\n- 1: `cli/feature.go`"}
	if _, err := Update(dir, func(s *State) error {
		if !s.EnqueueEntry(first) || s.EnqueueEntry(Enqueued{Story: "s-1", By: "alex", At: at, Verified: "other", Evidence: "other"}) {
			t.Error("enqueue twice")
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	first.At = at.UTC()
	if len(st.Enqueued) != 1 || st.Enqueued[0] != first || st.Enqueued[0].At.Location() != time.UTC {
		t.Errorf("enqueued %+v, want %+v", st.Enqueued, first)
	}
	data, err := os.ReadFile(StatePath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"verified": "4f1c2a9"`) || !strings.Contains(string(data), `"evidence": "Verdict: pass\n- 1: `+"`cli/feature.go`"+`"`) {
		t.Errorf("the state file does not record the verified commit and the evidence:\n%s", data)
	}
	var plain State
	plain.Enqueue("S-0002", "alex", at)
	if data, _ := json.Marshal(plain.Enqueued[0]); strings.Contains(string(data), "verified") || strings.Contains(string(data), "evidence") {
		t.Errorf("an enqueue by anyone else records neither: %s", data)
	}
}

// Hold records who held the queue and when, a second hold keeps the first,
// and Resume clears it.
func TestHoldAndResume(t *testing.T) {
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.UTC)
	var s State
	if s.Resume() {
		t.Error("resumed a queue not held")
	}
	if !s.Hold("alex", at) || s.Hold("orchestrator", at.Add(time.Hour)) || s.HeldBy != "alex" || !s.HeldAt.Equal(at) {
		t.Errorf("hold: %+v", s)
	}
	if !s.Resume() || s.Held || s.HeldBy != "" || !s.HeldAt.IsZero() {
		t.Errorf("resume: %+v", s)
	}
}

// RecordOutcome keeps one outcome per story, the last.
func TestRecordOutcomeReplacesTheStorysLast(t *testing.T) {
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.UTC)
	s := State{Outcomes: []Outcome{{Story: "S-0001", Outcome: "landed", At: at}, {Story: "S-0002", Outcome: "landed", At: at}}}
	s.RecordOutcome(Outcome{Story: "S-1", Outcome: "failed", At: at.Add(time.Hour)})
	if len(s.Outcomes) != 2 || s.Outcomes[0].Story != "S-0002" || s.Outcomes[1].Outcome != "failed" {
		t.Errorf("outcomes %+v", s.Outcomes)
	}
}

// DropBatchHolding drops only a batch that holds the story, and only once.
func TestDropBatchHolding(t *testing.T) {
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.UTC)
	s := State{Batch: &Batch{Stories: []string{"S-0001", "S-0002"}, Step: "verify"}}
	if s.DropBatchHolding("S-0003", "x", "alex", at) || s.Batch.Dropped != nil {
		t.Error("dropped a batch that does not hold the story")
	}
	if !s.DropBatchHolding("S-2", "why", "alex", at) || *s.Batch.Dropped != (Dropped{Reason: "why", By: "alex", At: at}) {
		t.Errorf("dropped %+v", s.Batch.Dropped)
	}
	if s.DropBatchHolding("S-0001", "again", "bo", at) || s.Batch.Dropped.Reason != "why" {
		t.Error("dropped a batch twice")
	}
	var none State
	if none.DropBatchHolding("S-0001", "x", "alex", at) {
		t.Error("dropped no batch")
	}
}

// ADR-0173: a story accepted by hand leaves the queue, the batch that holds
// it is dropped saying so, and its outcome is recorded; a story not in the
// queue writes no state.
func TestAcceptedByHandLeavesTheQueueAndDropsTheBatch(t *testing.T) {
	at := time.Date(2026, 10, 19, 9, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	if _, err := AcceptedByHand(dir, "S-0001", "alex", at); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(StatePath(dir)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("a story not in the queue wrote a state: %v", err)
	}
	writeState(t, dir, `{"enqueued": [{"story": "S-0001", "by": "orchestrator", "at": "2026-10-19T08:00:00Z"}, {"story": "S-0002", "by": "alex", "at": "2026-10-19T08:05:00Z"}],
  "batch": {"stories": ["S-0001", "S-0002"], "step": "verify"}}`)
	st, err := AcceptedByHand(dir, "S-0001", "alex", at)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Enqueued) != 1 || st.Enqueued[0].Story != "S-0002" {
		t.Errorf("enqueued %+v", st.Enqueued)
	}
	if st.Batch == nil || st.Batch.Dropped == nil || *st.Batch.Dropped != (Dropped{Reason: "S-0001 was accepted by hand by alex", By: "alex", At: at}) {
		t.Errorf("batch %+v", st.Batch)
	}
	want := Outcome{Story: "S-0001", Outcome: OutcomeAcceptedByHand, At: at, Detail: "accepted by hand by alex"}
	if len(st.Outcomes) != 1 || st.Outcomes[0] != want {
		t.Errorf("outcomes %+v, want %+v", st.Outcomes, want)
	}
	if read, _ := ReadState(dir); len(read.Enqueued) != 1 || read.Batch.Dropped == nil {
		t.Errorf("not written: %+v", read)
	}
}

// A refusal names every reason by kind after its sentence.
func TestARefusalNamesEachReason(t *testing.T) {
	r := &Refusal{Story: "S-0002", Text: "S-0002 is not enqueued", Reasons: []Reason{{Kind: KindVerifyMissing, Text: "no verify"}, {Kind: KindThreadOpen, Text: "TH-1 open"}}}
	if got := r.Error(); got != "S-0002 is not enqueued: verify_missing: no verify; thread_open: TH-1 open" {
		t.Errorf("refusal %q", got)
	}
	if got := (&Refusal{Text: "nothing to remove"}).Error(); got != "nothing to remove" {
		t.Errorf("refusal %q", got)
	}
}

package mergequeue

import (
	"os"
	"path/filepath"
	"strings"
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

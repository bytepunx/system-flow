package mergequeue

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// State is the merge queue's state file (ADR-0173): what is enqueued, the
// batch under verification, and each story's last outcome. S-0398 reads it
// only; the stories that enqueue, batch, and land write it.
type State struct {
	// Enqueued are the stories enqueued, in the order they were.
	Enqueued []Enqueued `json:"enqueued"`
	// Batch is the batch under verification; nil when there is none.
	Batch *Batch `json:"batch,omitempty"`
	// Outcomes are each story's last outcome, one per story.
	Outcomes []Outcome `json:"outcomes"`
}

// Enqueued is a story enqueued: who enqueued it and when, in UTC.
type Enqueued struct {
	Story string    `json:"story"`
	By    string    `json:"by"`
	At    time.Time `json:"at"`
}

// Batch is the batch under verification: its stories, in the order they are
// rebased, and the step it is at.
type Batch struct {
	Stories []string `json:"stories"`
	Step    string   `json:"step"`
}

// Outcome is a story's last outcome in the queue, when it came, in UTC, and
// what it came to in a sentence.
type Outcome struct {
	Story   string    `json:"story"`
	Outcome string    `json:"outcome"`
	At      time.Time `json:"at"`
	Detail  string    `json:"detail,omitempty"`
}

// StatePath is where the queue's state is in the cache folder cacheDir, the
// main checkout's .flai-cache.
func StatePath(cacheDir string) string {
	return filepath.Join(cacheDir, "merge-queue", "state.json")
}

// ReadState is the queue's state in cacheDir: empty, with no batch and no
// outcomes, when there is no file.
func ReadState(cacheDir string) (State, error) {
	empty := State{Enqueued: []Enqueued{}, Outcomes: []Outcome{}}
	path := StatePath(cacheDir)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil
	}
	if err != nil {
		return empty, fmt.Errorf("read the merge queue's state at %s: %w; check that the project's .flai-cache is readable", path, err)
	}
	st := empty
	if err := json.Unmarshal(data, &st); err != nil {
		return empty, fmt.Errorf("read the merge queue's state at %s: %w; fix the file, or delete it to empty the queue, and enqueue its stories again", path, err)
	}
	if st.Enqueued == nil {
		st.Enqueued = []Enqueued{}
	}
	if st.Outcomes == nil {
		st.Outcomes = []Outcome{}
	}
	return st, nil
}

package mergequeue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/filelock"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// State is the merge queue's state file (ADR-0173): what is enqueued,
// whether the queue is held, the batch under verification, and each story's
// last outcome. Update is its only writer.
type State struct {
	// Enqueued are the stories enqueued, in the order they were.
	Enqueued []Enqueued `json:"enqueued"`
	// Held says the queue builds no batch until it is resumed, and HeldBy
	// and HeldAt, in UTC, who held it and when.
	Held   bool      `json:"held,omitempty"`
	HeldBy string    `json:"held_by,omitempty"`
	HeldAt time.Time `json:"held_at,omitzero"`
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
// rebased, the branch and worktree it is built in, when it started, the step
// it is at, its result once verified, and why it was dropped, when it was.
type Batch struct {
	ID       string   `json:"id,omitempty"`
	Branch   string   `json:"branch,omitempty"`
	Worktree string   `json:"worktree,omitempty"`
	Stories  []string `json:"stories"`
	// StartedAt is when the batch was built, in UTC.
	StartedAt time.Time `json:"started_at,omitzero"`
	Step      string    `json:"step"`
	// Result is the batch's verification once it finished; nil before.
	Result *Result `json:"result,omitempty"`
	// Dropped is why the batch was dropped, by whom and when; nil while it
	// was not.
	Dropped *Dropped `json:"dropped,omitempty"`
}

// Result is a batch's verification: whether it passed, the step or test
// tier it stopped at when it did not, its findings in short, and when it
// finished, in UTC.
type Result struct {
	Passed     bool      `json:"passed"`
	StoppedAt  string    `json:"stopped_at,omitempty"`
	Summary    string    `json:"summary,omitempty"`
	FinishedAt time.Time `json:"finished_at"`
}

// Dropped is why a batch was dropped, who dropped it, and when, in UTC.
type Dropped struct {
	Reason string    `json:"reason"`
	By     string    `json:"by"`
	At     time.Time `json:"at"`
}

// Outcome is a story's last outcome in the queue, when it came, in UTC, and
// what it came to in a sentence.
type Outcome struct {
	Story   string    `json:"story"`
	Outcome string    `json:"outcome"`
	At      time.Time `json:"at"`
	Detail  string    `json:"detail,omitempty"`
}

// OutcomeAcceptedByHand is the outcome of a story in the queue that someone
// accepted with flai accept rather than the queue landing it.
const OutcomeAcceptedByHand = "accepted_by_hand"

// lockWait is how long Update waits for another writer of the state, flai
// serve's loop or a shell's flai queue, to finish with it.
const lockWait = 10 * time.Second

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
	st.fill()
	return st, nil
}

// Update reads the queue's state in cacheDir under a lock beside it, lets fn
// change it, and writes what fn leaves in one step, so that flai serve and a
// shell never write the state at once and no write is lost. An error from fn
// writes nothing and is returned as it is; a state fn leaves as it was is
// not written.
func Update(cacheDir string, fn func(*State) error) (State, error) {
	path := StatePath(cacheDir)
	unlock, err := filelock.Lock(path+".lock", lockWait)
	if err != nil {
		return State{}, fmt.Errorf("write the merge queue's state at %s: %w", path, err)
	}
	defer unlock()
	st, err := ReadState(cacheDir)
	if err != nil {
		return State{}, err
	}
	before, err := json.Marshal(st)
	if err != nil {
		return State{}, fmt.Errorf("write the merge queue's state at %s: %w", path, err)
	}
	if err := fn(&st); err != nil {
		return State{}, err
	}
	st.fill()
	after, err := json.Marshal(st)
	if err != nil {
		return State{}, fmt.Errorf("write the merge queue's state at %s: %w", path, err)
	}
	if bytes.Equal(before, after) {
		return st, nil
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return State{}, fmt.Errorf("write the merge queue's state at %s: %w", path, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return State{}, fmt.Errorf("write the merge queue's state at %s: %w; check that the project's .flai-cache is writable", path, err)
	}
	return st, nil
}

// fill makes the state's lists empty rather than null, so that every
// surface shows [].
func (s *State) fill() {
	if s.Enqueued == nil {
		s.Enqueued = []Enqueued{}
	}
	if s.Outcomes == nil {
		s.Outcomes = []Outcome{}
	}
}

// EnqueuedAs is story's place in the queue, who enqueued it and when, and
// whether it is enqueued.
func (s *State) EnqueuedAs(story string) (Enqueued, bool) {
	for _, e := range s.Enqueued {
		if sameStory(e.Story, story) {
			return e, true
		}
	}
	return Enqueued{}, false
}

// Enqueue puts story at the end of the queue as enqueued by by at at, and
// reports false, changing nothing, when it is enqueued already.
func (s *State) Enqueue(story, by string, at time.Time) bool {
	if _, ok := s.EnqueuedAs(story); ok {
		return false
	}
	s.Enqueued = append(s.Enqueued, Enqueued{Story: story, By: by, At: at.UTC()})
	return true
}

// Dequeue takes story out of the queue, and reports false when it was not
// enqueued.
func (s *State) Dequeue(story string) bool {
	n := len(s.Enqueued)
	s.Enqueued = slices.DeleteFunc(s.Enqueued, func(e Enqueued) bool { return sameStory(e.Story, story) })
	return len(s.Enqueued) != n
}

// Hold holds the queue as by at at, and reports false, keeping who held it
// first and when, when it is held already.
func (s *State) Hold(by string, at time.Time) bool {
	if s.Held {
		return false
	}
	s.Held, s.HeldBy, s.HeldAt = true, by, at.UTC()
	return true
}

// Resume lets the queue build batches again, and reports false when it was
// not held.
func (s *State) Resume() bool {
	if !s.Held {
		return false
	}
	s.Held, s.HeldBy, s.HeldAt = false, "", time.Time{}
	return true
}

// RecordOutcome records o as its story's last outcome, in place of the one
// before.
func (s *State) RecordOutcome(o Outcome) {
	o.At = o.At.UTC()
	s.Outcomes = slices.DeleteFunc(s.Outcomes, func(x Outcome) bool { return sameStory(x.Story, o.Story) })
	s.Outcomes = append(s.Outcomes, o)
}

// HoldsStory reports whether the batch under verification holds story.
func (s *State) HoldsStory(story string) bool {
	return s.Batch != nil && slices.ContainsFunc(s.Batch.Stories, func(x string) bool { return sameStory(x, story) })
}

// DropBatchHolding marks the batch that holds story dropped, for reason, by
// by at at, and reports false when no batch holds it or the batch was
// dropped already.
func (s *State) DropBatchHolding(story, reason, by string, at time.Time) bool {
	if !s.HoldsStory(story) || s.Batch.Dropped != nil {
		return false
	}
	s.Batch.Dropped = &Dropped{Reason: reason, By: by, At: at.UTC()}
	return true
}

// sameStory reports whether a and b name the same story, however each is
// written.
func sameStory(a, b string) bool {
	return workitem.CanonicalID(a) == workitem.CanonicalID(b)
}

// Change is what a write to the queue did: its action, the story it was on,
// whether it changed the state, who made it and when, in UTC, and the state
// it left.
type Change struct {
	Action  string    `json:"action"`
	Story   string    `json:"story,omitempty"`
	Changed bool      `json:"changed"`
	By      string    `json:"by"`
	At      time.Time `json:"at"`
	State   State     `json:"state"`
}

// The actions a Change names.
const (
	ActionAdd    = "add"
	ActionRemove = "remove"
	ActionHold   = "hold"
	ActionResume = "resume"
)

// Refusal is a write to the queue refused: the story it was on, a sentence
// saying why and what to do, and, for a story not eligible, every reason it
// is not.
type Refusal struct {
	Story   string   `json:"story,omitempty"`
	Text    string   `json:"text"`
	Reasons []Reason `json:"reasons,omitempty"`
}

// Error is the refusal's sentence, with each reason by kind.
func (r *Refusal) Error() string {
	if len(r.Reasons) == 0 {
		return r.Text
	}
	parts := make([]string, 0, len(r.Reasons))
	for _, x := range r.Reasons {
		parts = append(parts, x.Kind+": "+x.Text)
	}
	return r.Text + ": " + strings.Join(parts, "; ")
}

// Add enqueues story in repo's queue as by at at, once it is a story in
// review and eligible, reading eligibility with Compute through r. A story
// enqueued already is left in its place, with who enqueued it first, and the
// change says nothing changed. Under pull-request acceptance, and for a
// story not in review or not eligible, it is refused with a *Refusal naming
// why.
func Add(r execx.Runner, repo *workitem.Repo, story, by string, at time.Time) (Change, error) {
	ch := Change{Action: ActionAdd, Story: story, By: by, At: at.UTC()}
	if repo.Manifest.Acceptance.PullRequest() {
		return ch, &Refusal{Story: story, Text: fmt.Sprintf("%s is not enqueued: %s; accept it with flai accept %s", story, OffUnderPullRequest, story)}
	}
	it, err := repo.Get(story)
	if err != nil {
		return ch, fmt.Errorf("find %s to enqueue it in the merge queue: %w", story, err)
	}
	ch.Story = it.ID
	if it.Type != workitem.Story {
		return ch, &Refusal{Story: it.ID, Text: fmt.Sprintf("%s is not enqueued: the merge queue takes stories only, and %s is a %s", it.ID, it.ID, it.Type)}
	}
	if it.Status != workitem.Review {
		return ch, &Refusal{Story: it.ID, Text: fmt.Sprintf("%s is not enqueued: the merge queue takes a story in review, and %s is %s; move it to review once its work is done", it.ID, it.ID, it.Status)}
	}
	q, err := Compute(r, repo)
	if err != nil {
		return ch, err
	}
	if _, ok := (&State{Enqueued: q.Enqueued}).EnqueuedAs(it.ID); !ok {
		for _, e := range q.Ineligible {
			if sameStory(e.Story, it.ID) {
				return ch, &Refusal{Story: it.ID, Text: fmt.Sprintf("%s is not enqueued: it is not eligible for the merge queue", it.ID), Reasons: e.Reasons}
			}
		}
	}
	ch.State, err = Update(repo.CacheDir(), func(s *State) error {
		ch.Changed = s.Enqueue(it.ID, by, at)
		return nil
	})
	return ch, err
}

// Remove takes story out of repo's queue as by at at, and marks dropped the
// batch that holds it, so that the queue does not land a story no one has
// enqueued. A story not enqueued is refused with a *Refusal.
func Remove(repo *workitem.Repo, story, by string, at time.Time) (Change, error) {
	id := workitem.CanonicalID(story)
	ch := Change{Action: ActionRemove, Story: id, By: by, At: at.UTC()}
	var err error
	ch.State, err = Update(repo.CacheDir(), func(s *State) error {
		if !s.Dequeue(id) {
			return &Refusal{Story: id, Text: fmt.Sprintf("%s is not enqueued in the merge queue, so there is nothing to remove; flai queue lists what is", id)}
		}
		s.DropBatchHolding(id, id+" was removed from the queue by "+by, by, at)
		ch.Changed = true
		return nil
	})
	return ch, err
}

// Hold holds repo's queue as by at at: it builds no batch until resumed. A
// queue held already keeps who held it first, and the change says nothing
// changed.
func Hold(repo *workitem.Repo, by string, at time.Time) (Change, error) {
	ch := Change{Action: ActionHold, By: by, At: at.UTC()}
	var err error
	ch.State, err = Update(repo.CacheDir(), func(s *State) error {
		ch.Changed = s.Hold(by, at)
		return nil
	})
	return ch, err
}

// Resume lets repo's queue build batches again, as by at at. A queue not
// held is left as it is, and the change says nothing changed.
func Resume(repo *workitem.Repo, by string, at time.Time) (Change, error) {
	ch := Change{Action: ActionResume, By: by, At: at.UTC()}
	var err error
	ch.State, err = Update(repo.CacheDir(), func(s *State) error {
		ch.Changed = s.Resume()
		return nil
	})
	return ch, err
}

// AcceptedByHand takes story, which by accepted with flai accept at at, out
// of the queue in cacheDir, marks dropped the batch that holds it, and
// records the acceptance as its outcome. A story neither enqueued nor in the
// batch leaves the state as it was.
func AcceptedByHand(cacheDir, story, by string, at time.Time) (State, error) {
	return Update(cacheDir, func(s *State) error {
		dequeued, batched := s.Dequeue(story), s.HoldsStory(story)
		s.DropBatchHolding(story, story+" was accepted by hand by "+by, by, at)
		if dequeued || batched {
			s.RecordOutcome(Outcome{Story: story, Outcome: OutcomeAcceptedByHand, At: at, Detail: "accepted by hand by " + by})
		}
		return nil
	})
}

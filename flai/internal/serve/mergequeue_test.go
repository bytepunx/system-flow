package serve

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeBatches stands in for mergequeue.RunBatch: it counts its calls and how
// many run at once, says on entered when one begins, and, while block is
// open, waits for it to close or for its ctx to end.
type fakeBatches struct {
	mu      sync.Mutex
	calls   int
	running int
	most    int
	deps    []mergequeue.Deps
	repos   []*workitem.Repo
	ctxErr  error
	err     error
	block   chan struct{}
	entered chan struct{}
}

func newFakeBatches() *fakeBatches {
	return &fakeBatches{entered: make(chan struct{}, 64)}
}

func (f *fakeBatches) run(ctx context.Context, d mergequeue.Deps, repo *workitem.Repo) (bool, error) {
	f.mu.Lock()
	f.calls++
	f.running++
	f.most = max(f.most, f.running)
	f.deps, f.repos = append(f.deps, d), append(f.repos, repo)
	block, err := f.block, f.err
	f.mu.Unlock()
	defer func() {
		f.mu.Lock()
		f.running--
		f.mu.Unlock()
	}()
	select {
	case f.entered <- struct{}{}:
	default:
	}
	if block != nil {
		select {
		case <-block:
		case <-ctx.Done():
			f.mu.Lock()
			f.ctxErr = ctx.Err()
			f.mu.Unlock()
			return true, ctx.Err()
		}
	}
	return true, err
}

func (f *fakeBatches) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// queueLab is a project with no story and its merge queue pass, whose
// batches are fakeBatches'. The merge-queue host action is on, acceptance is
// local, and accept_reviews is on.
type queueLab struct {
	t       *testing.T
	root    string
	on      bool
	logs    *syncBuffer
	batches *fakeBatches
	pass    *queuePass
	runs    *queueRuns
}

func newQueueLab(t *testing.T) *queueLab {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	lab := &queueLab{t: t, root: root, on: true, logs: &syncBuffer{}, batches: newFakeBatches(), runs: newQueueRuns()}
	lab.manifest("local", true)
	o := Options{Logger: slog.New(slog.NewJSONHandler(lab.logs, nil)), Now: time.Now, Version: "test",
		Host: hostapi.Host{Enabled: func(action, r string) bool {
			return lab.on && action == hostapi.ActionMergeQueue && r == root
		}}}
	lab.pass = newQueuePass(o, Entry{Key: "t", Name: "t", Root: root}, lab.runs)
	lab.pass.run = lab.batches.run
	return lab
}

// manifest writes the project's manifest with acceptance mode and
// accept_reviews on or off.
func (lab *queueLab) manifest(mode string, acceptReviews bool) {
	lab.t.Helper()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nacceptance:\n  mode: " + mode + "\n"
	if acceptReviews {
		m += "orchestration:\n  permissions:\n    accept_reviews: true\n"
	}
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

func (lab *queueLab) repo() *workitem.Repo {
	lab.t.Helper()
	repo, err := workitem.Open(lab.root)
	if err != nil {
		lab.t.Fatal(err)
	}
	return repo
}

func (lab *queueLab) look() { lab.pass.look(context.Background()) }

// records are the log's records whose field key is value; any value when
// value is empty.
func (lab *queueLab) records(key, value string) []map[string]any {
	lab.t.Helper()
	var out []map[string]any
	for line := range strings.Lines(lab.logs.String()) {
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			lab.t.Fatalf("log line %q: %v", line, err)
		}
		if v, ok := rec[key]; ok && (value == "" || v == value) {
			out = append(out, rec)
		}
	}
	return out
}

// S-0399: with the merge-queue host action off, the queue builds no batch
// and says nothing, whatever the manifest allows.
func TestTheMergeQueueRunsNoBatchWithTheActionOff(t *testing.T) {
	lab := newQueueLab(t)
	lab.on = false
	lab.look()
	lab.look()
	if n := lab.batches.count(); n != 0 {
		t.Fatalf("batches run = %d, want 0", n)
	}
	if logs := lab.logs.String(); logs != "" {
		t.Fatalf("logged with the action off: %s", logs)
	}
}

// S-0399: with accept_reviews off, the queue builds no batch, and says why
// once across looks.
func TestTheMergeQueueRunsNoBatchWithAcceptReviewsOff(t *testing.T) {
	lab := newQueueLab(t)
	lab.manifest("local", false)
	lab.look()
	lab.look()
	if n := lab.batches.count(); n != 0 {
		t.Fatalf("batches run = %d, want 0", n)
	}
	if recs := lab.records("reason", queueNoAcceptReviews); len(recs) != 1 {
		t.Fatalf("said accept_reviews is off %d times, want once:\n%s", len(recs), lab.logs.String())
	}
}

// S-0399 criterion 2: under pull-request acceptance the queue builds no
// batch and says why once across looks; back under local acceptance it runs,
// and under pull-request acceptance again it says so again.
func TestTheMergeQueueRunsNoBatchUnderPullRequestAcceptanceAndSaysSoOnce(t *testing.T) {
	lab := newQueueLab(t)
	lab.manifest("pull-request", true)
	lab.look()
	lab.look()
	lab.look()
	if n := lab.batches.count(); n != 0 {
		t.Fatalf("batches run = %d, want 0", n)
	}
	recs := lab.records("reason", queuePullRequest)
	if len(recs) != 1 || recs[0]["level"] != "INFO" || recs[0]["detail"] != mergequeue.OffUnderPullRequest || recs[0]["root"] != lab.root {
		t.Fatalf("said pull-request acceptance stops it %d times, want once with the detail:\n%s", len(recs), lab.logs.String())
	}
	lab.manifest("local", true)
	lab.look()
	if n := lab.batches.count(); n != 1 {
		t.Fatalf("batches run under local acceptance = %d, want 1", n)
	}
	if recs := lab.records("was", queuePullRequest); len(recs) != 1 {
		t.Fatalf("said it may run again %d times, want once:\n%s", len(recs), lab.logs.String())
	}
	lab.manifest("pull-request", true)
	lab.look()
	if recs := lab.records("reason", queuePullRequest); len(recs) != 2 {
		t.Fatalf("said pull-request acceptance stops it %d times after the mode came back, want twice:\n%s", len(recs), lab.logs.String())
	}
}

// S-0399 criterion 4: while the queue is held no batch starts, and once it
// is resumed one does.
func TestTheMergeQueueRunsNoBatchWhileHeld(t *testing.T) {
	lab := newQueueLab(t)
	if _, err := mergequeue.Hold(lab.repo(), "alex", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.look()
	lab.look()
	if n := lab.batches.count(); n != 0 {
		t.Fatalf("batches run while held = %d, want 0", n)
	}
	if recs := lab.records("reason", queueHeld); len(recs) != 1 || recs[0]["held_by"] != "alex" {
		t.Fatalf("said it is held %d times, want once with who held it:\n%s", len(recs), lab.logs.String())
	}
	if _, err := mergequeue.Resume(lab.repo(), "alex", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.look()
	if n := lab.batches.count(); n != 1 {
		t.Fatalf("batches run once resumed = %d, want 1", n)
	}
}

// S-0399 criteria 2 and 3: with the action, accept_reviews, and local
// acceptance on, a look runs one batch of the project's queue in flai
// serve's process, with its git, clock, verifier, and logger, and logs it.
func TestTheMergeQueueRunsABatchWhenEverythingIsOn(t *testing.T) {
	lab := newQueueLab(t)
	lab.look()
	if n := lab.batches.count(); n != 1 {
		t.Fatalf("batches run = %d, want 1", n)
	}
	d, repo := lab.batches.deps[0], lab.batches.repos[0]
	if d.Runner == nil || d.Now == nil || d.Verify == nil || d.Log == nil {
		t.Fatalf("deps = %+v, want each set", d)
	}
	if repo == nil || repo.Manifest.Name != "t" {
		t.Fatalf("repo = %+v, want the project's", repo)
	}
	recs := lab.records("took", "")
	if len(recs) != 1 || recs[0]["level"] != "INFO" || recs[0]["root"] != lab.root || recs[0]["key"] != "t" {
		t.Fatalf("logged the batch run as %v, want one info line with the project:\n%s", recs, lab.logs.String())
	}
	if recs := lab.records("reason", ""); len(recs) != 0 {
		t.Fatalf("said the queue did not run:\n%s", lab.logs.String())
	}
}

// S-0399: a batch that fails is logged as an error with the project and the
// error, and the next look tries again.
func TestAFailedMergeQueueBatchIsLoggedAndTriedAgain(t *testing.T) {
	lab := newQueueLab(t)
	lab.batches.err = errors.New("build merge queue batch 1: git refused")
	lab.look()
	lab.look()
	recs := lab.records("err", "build merge queue batch 1: git refused")
	if len(recs) != 2 || recs[0]["level"] != "ERROR" || recs[0]["root"] != lab.root {
		t.Fatalf("logged the failure as %v, want an error line a look:\n%s", recs, lab.logs.String())
	}
	if n := lab.batches.count(); n != 2 {
		t.Fatalf("batches run = %d, want 2", n)
	}
}

// S-0399 criterion 3: the queue runs in its own goroutine, one batch at a
// time per project: looks and wakes while a batch runs start no other, and
// once it ends the next look runs one.
func TestTheMergeQueueRunsOneBatchAtATimePerProject(t *testing.T) {
	lab := newQueueLab(t)
	lab.pass.every = time.Millisecond
	lab.batches.block = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		lab.runs.wait()
	}()
	lab.pass.start(ctx)
	waitEntered(t, lab.batches)
	for range 5 {
		lab.pass.wake()
		lab.pass.look(ctx) // another pass for the project, as one still finishing would be
	}
	time.Sleep(50 * time.Millisecond) // ticks of the loop, which waits on the batch
	if n := lab.batches.count(); n != 1 {
		t.Fatalf("batches run while one runs = %d, want 1", n)
	}
	lab.batches.mu.Lock()
	close(lab.batches.block)
	lab.batches.mu.Unlock()
	waitEntered(t, lab.batches)
	lab.batches.mu.Lock()
	defer lab.batches.mu.Unlock()
	if lab.batches.most != 1 {
		t.Fatalf("batches run at once = %d, want 1", lab.batches.most)
	}
}

// S-0399: stopping the pass, as flai serve stopping does, cancels the batch
// running, waits for it, and logs it as stopped rather than failed.
func TestStoppingTheMergeQueueCancelsTheBatchRunning(t *testing.T) {
	lab := newQueueLab(t)
	lab.batches.block = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	lab.pass.start(ctx)
	waitEntered(t, lab.batches)
	cancel()
	done := make(chan struct{})
	go func() {
		lab.runs.wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the pass did not end once stopped")
	}
	lab.batches.mu.Lock()
	ctxErr := lab.batches.ctxErr
	lab.batches.mu.Unlock()
	if !errors.Is(ctxErr, context.Canceled) {
		t.Fatalf("the batch's ctx ended with %v, want it cancelled", ctxErr)
	}
	recs := lab.records("err", context.Canceled.Error())
	if len(recs) != 1 || recs[0]["level"] != "INFO" {
		t.Fatalf("logged the stop as %v, want one info line:\n%s", recs, lab.logs.String())
	}
	if recs := lab.records("level", "ERROR"); len(recs) != 0 {
		t.Fatalf("logged the stop as a failure:\n%s", lab.logs.String())
	}
}

// waitEntered waits for a batch to begin.
func waitEntered(t *testing.T, f *fakeBatches) {
	t.Helper()
	select {
	case <-f.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("no batch began")
	}
}

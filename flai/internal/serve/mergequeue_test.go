package serve

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
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
	// idle has a batch run begin nothing, as RunBatch does while a verified
	// batch waits to be landed or bisected; order says, in order, "run" for
	// each batch run and "advance" for each advance, the lab's fakeAdvances'.
	idle  bool
	order []string
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
	f.order = append(f.order, "run")
	block, err, idle := f.block, f.err, f.idle
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
	return !idle, err
}

func (f *fakeBatches) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

// fakeAdvances stands in for mergequeue.Advance: it notes each call in the
// batches' order, keeps its Deps and Accepter, runs then, when set, as the
// advance's work, and answers did and err.
type fakeAdvances struct {
	batches *fakeBatches
	deps    []mergequeue.Deps
	accepts []mergequeue.Accepter
	then    func(ctx context.Context, repo *workitem.Repo)
	did     bool
	err     error
}

func (f *fakeAdvances) advance(ctx context.Context, d mergequeue.Deps, accept mergequeue.Accepter, repo *workitem.Repo) (bool, error) {
	f.batches.mu.Lock()
	f.batches.order = append(f.batches.order, "advance")
	f.deps, f.accepts = append(f.deps, d), append(f.accepts, accept)
	then, did, err := f.then, f.did, f.err
	f.batches.mu.Unlock()
	if then != nil {
		then(ctx, repo)
	}
	return did, err
}

func (f *fakeAdvances) count() int {
	f.batches.mu.Lock()
	defer f.batches.mu.Unlock()
	return len(f.deps)
}

// queueLab is a project with no story and its merge queue pass, whose
// batches are fakeBatches' and advances fakeAdvances'. The merge-queue host
// action is on, acceptance is local, and accept_reviews is on. The host's
// journal is kept in journal.
type queueLab struct {
	t        *testing.T
	root     string
	on       bool
	logs     *syncBuffer
	batches  *fakeBatches
	advances *fakeAdvances
	pass     *queuePass
	runs     *queueRuns

	mu      sync.Mutex
	journal []hostapi.Entry
}

func newQueueLab(t *testing.T) *queueLab {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"design", "docs", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	lab := &queueLab{t: t, root: root, on: true, logs: &syncBuffer{}, batches: newFakeBatches(), runs: newQueueRuns()}
	lab.advances = &fakeAdvances{batches: lab.batches}
	lab.manifest("local", true)
	o := Options{Logger: slog.New(slog.NewJSONHandler(lab.logs, nil)), Now: time.Now, Version: "test",
		Host: hostapi.Host{Enabled: func(action, r string) bool {
			return lab.on && action == hostapi.ActionMergeQueue && r == root
		}, Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.journal = append(lab.journal, e)
		}}}
	lab.pass = newQueuePass(o, Entry{Key: "t", Name: "t", Root: root}, lab.runs)
	lab.pass.run = lab.batches.run
	lab.pass.advance = lab.advances.advance
	return lab
}

// entries are the host's journal so far.
func (lab *queueLab) entries() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	return slices.Clone(lab.journal)
}

// verifiedBatch writes the queue's state with batch b1, holding S-0001 and
// S-0002, verified and passed, as RunBatch leaves one to be landed.
func (lab *queueLab) verifiedBatch() {
	lab.t.Helper()
	if _, err := mergequeue.Update(lab.repo().CacheDir(), func(s *mergequeue.State) error {
		s.Batch = &mergequeue.Batch{ID: "b1", Stories: []string{"S-0001", "S-0002"}, Step: mergequeue.StepVerified, Result: &mergequeue.Result{Passed: true}}
		return nil
	}); err != nil {
		lab.t.Fatal(err)
	}
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
	if n, m := lab.batches.count(), lab.advances.count(); n != 0 || m != 0 {
		t.Fatalf("batches run = %d, advanced = %d, want 0 and 0", n, m)
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
	if n, m := lab.batches.count(), lab.advances.count(); n != 0 || m != 0 {
		t.Fatalf("batches run = %d, advanced = %d, want 0 and 0", n, m)
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
	if n, m := lab.batches.count(), lab.advances.count(); n != 0 || m != 0 {
		t.Fatalf("batches run = %d, advanced = %d, want 0 and 0", n, m)
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
	if n, m := lab.batches.count(), lab.advances.count(); n != 0 || m != 0 {
		t.Fatalf("batches run while held = %d, advanced = %d, want 0 and 0", n, m)
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
	if n := lab.advances.count(); n != 0 {
		t.Fatalf("advanced %d times after the batch's run stopped, want 0", n)
	}
}

// S-0400 criteria 1, 3, and 5: a look takes the batch on after the batch's
// run, with the same Deps and an Accepter, also when the run began nothing,
// as it does while a verified batch waits; an advance that did something has
// the pass look again at once and is logged once with the project.
func TestALookAdvancesTheBatchAfterItsRun(t *testing.T) {
	lab := newQueueLab(t)
	lab.batches.idle = true
	lab.advances.did = true
	lab.look()
	lab.batches.idle = false
	lab.advances.did = false
	lab.look()
	if got := strings.Join(lab.batches.order, ","); got != "run,advance,run,advance" {
		t.Fatalf("calls = %s, want each look's run before its advance", got)
	}
	for i, d := range lab.advances.deps {
		run := lab.batches.deps[i]
		if d.Runner != run.Runner || d.Log != run.Log || d.Now == nil || d.Verify == nil || lab.advances.accepts[i] == nil {
			t.Fatalf("advance %d had deps %+v and an accepter %v, want the run's deps and one", i, d, lab.advances.accepts[i] != nil)
		}
	}
	recs := lab.records("msg", "merge queue batch advanced")
	if len(recs) != 1 || recs[0]["level"] != "INFO" || recs[0]["root"] != lab.root || recs[0]["key"] != "t" || recs[0]["took"] == nil {
		t.Fatalf("logged the advance as %v, want one info line, for the look that did something:\n%s", recs, lab.logs.String())
	}
	select {
	case <-lab.pass.woken:
	default:
		t.Fatal("an advance that did something did not wake the pass")
	}
}

// S-0400: an advance that fails is logged as an error, with the project and
// the error, and the next look advances again; one that dropped the batch
// as it failed is logged once as a warning with why, and journalled as a
// failed drop.
func TestAFailedAdvanceIsLoggedAndTheNextLookGoesOn(t *testing.T) {
	lab := newQueueLab(t)
	lab.advances.err = errors.New("send S-0002 back: thread refused")
	lab.look()
	lab.look()
	recs := lab.records("err", "send S-0002 back: thread refused")
	if len(recs) != 2 || recs[0]["level"] != "ERROR" || recs[0]["msg"] != "merge queue batch not advanced" || recs[0]["root"] != lab.root {
		t.Fatalf("logged the failure as %v, want an error line a look:\n%s", recs, lab.logs.String())
	}
	if n := lab.advances.count(); n != 2 {
		t.Fatalf("advances = %d, want 2", n)
	}
	select {
	case <-lab.pass.woken:
		t.Fatal("a failed advance woke the pass, which would build the batch again at once")
	default:
	}

	lab.verifiedBatch()
	lab.advances.err, lab.advances.did = errors.New("land S-0001 from merge queue batch b1 as alex: exit status 1"), true
	lab.advances.then = func(_ context.Context, repo *workitem.Repo) {
		if _, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
			s.Batch.Dropped = &mergequeue.Dropped{Reason: "landing S-0001 failed", By: "flai serve", At: time.Now()}
			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	lab.look()
	recs = lab.records("reason", "landing S-0001 failed")
	if len(recs) != 1 || recs[0]["level"] != "WARN" || recs[0]["msg"] != "merge queue batch dropped" || recs[0]["batch"] != "b1" || recs[0]["err"] == nil {
		t.Fatalf("logged the drop as %v, want one warning with the batch, why, and the error:\n%s", recs, lab.logs.String())
	}
	if recs := lab.records("level", "ERROR"); len(recs) != 2 {
		t.Fatalf("logged the drop as an error too:\n%s", lab.logs.String())
	}
	got := lab.entries()
	if len(got) != 1 || got[0].Action != hostapi.ActionMergeQueue || got[0].Outcome != "failed" || got[0].By != queueAuthor ||
		!strings.Contains(got[0].Detail, "b1") || !strings.Contains(got[0].Detail, "landing S-0001 failed") {
		t.Fatalf("journal = %+v, want the drop, failed", got)
	}
}

// S-0400: an advance stopped as flai serve stops is logged as stopped, not
// failed.
func TestAnAdvanceStoppedIsLoggedAsStopped(t *testing.T) {
	lab := newQueueLab(t)
	ctx, cancel := context.WithCancel(context.Background())
	lab.advances.then = func(context.Context, *workitem.Repo) { cancel() }
	lab.advances.err = context.Canceled
	lab.pass.look(ctx)
	recs := lab.records("err", context.Canceled.Error())
	if len(recs) != 1 || recs[0]["level"] != "INFO" {
		t.Fatalf("logged the stop as %v, want one info line:\n%s", recs, lab.logs.String())
	}
	if recs := lab.records("level", "ERROR"); len(recs) != 0 {
		t.Fatalf("logged the stop as a failure:\n%s", lab.logs.String())
	}
}

// S-0400 criteria 1 and 4: each story an advance landed or sent back is
// journalled as flai serve's, with its outcome, and logged on the one line.
func TestEachStoryLandedOrSentBackIsJournalled(t *testing.T) {
	lab := newQueueLab(t)
	lab.verifiedBatch()
	lab.advances.did = true
	at := time.Date(2026, 10, 19, 22, 0, 0, 0, time.UTC)
	lab.advances.then = func(_ context.Context, repo *workitem.Repo) {
		if _, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
			s.RecordOutcome(mergequeue.Outcome{Story: "S-0001", Outcome: mergequeue.OutcomeLanded, At: at, Detail: "landed by merge queue batch b1 as alex"})
			s.RecordOutcome(mergequeue.Outcome{Story: "S-0002", Outcome: mergequeue.OutcomeFailing, At: at, Detail: "sent back to in-progress: merge queue batch b1 fails with it at unit"})
			s.Batch = nil
			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	lab.look()
	got := lab.entries()
	if len(got) != 2 || got[0].Detail != "S-0001: landed by merge queue batch b1 as alex" || got[1].Detail != "S-0002: sent back to in-progress: merge queue batch b1 fails with it at unit" ||
		got[0].Outcome != "done" || got[1].At != "2026-10-19T22:00:00Z" || got[1].Method != "serve.merge-queue" || got[1].Project != "t" {
		t.Fatalf("journal = %+v, want the landing and the send-back", got)
	}
	recs := lab.records("msg", "merge queue batch advanced")
	if len(recs) != 1 || recs[0]["batch"] != "b1" || recs[0]["landed"] != "S-0001" || recs[0]["sent_back"] != "S-0002" {
		t.Fatalf("logged the advance as %v, want the batch, the story landed, and the one sent back:\n%s", recs, lab.logs.String())
	}
	lab.look() // nothing new: nothing journalled again
	if got := lab.entries(); len(got) != 2 {
		t.Fatalf("journal after a look that changed nothing = %+v, want the two entries alone", got)
	}
}

// holdFindings are a preflight's findings on the main checkout: a task note
// whose link and message carry what Markdown would read as code, a link, and
// a tab, and an error of flai check.
var holdFindings = []verify.FileFinding{
	{Tier: "markdown", Checkout: verify.CheckoutMain, Path: "wip/kanban/tasks/T-0001-x.md", Line: 12, Rule: "MD051", Message: "Link fragments should be valid [Context: \"[see](#nowhere)\"]\tat https://example.com `x`"},
	{Tier: mergequeue.PreflightCheckTier, Checkout: verify.CheckoutMain, Path: "wip/kanban/stories/S-0001-x.md", Rule: "item.front-matter", Message: "status is not one of the statuses"},
}

// holdingRun has the lab's batch runs hold the queue on holdFindings since
// since while held is set, and clear the hold otherwise, as RunBatch's
// preflight does, beginning no batch.
func (lab *queueLab) holdingRun(held *atomic.Bool, since time.Time) {
	lab.pass.run = func(_ context.Context, _ mergequeue.Deps, repo *workitem.Repo) (bool, error) {
		lab.batches.mu.Lock()
		lab.batches.calls++
		lab.batches.mu.Unlock()
		_, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
			switch {
			case !held.Load():
				s.Preflight = nil
			case s.Preflight == nil:
				s.Preflight = &mergequeue.PreflightHold{Since: since, Main: "0123456789abcdef", Findings: holdFindings}
			}
			return nil
		})
		return false, err
	}
}

// holdThreads are the threads on the lab's manifest.
func (lab *queueLab) holdThreads() []*threads.Thread {
	lab.t.Helper()
	all, err := threads.For(lab.repo(), "system-flow.yaml")
	if err != nil {
		lab.t.Fatal(err)
	}
	return all
}

// state is the lab's merge queue state.
func (lab *queueLab) state() mergequeue.State {
	lab.t.Helper()
	st, err := mergequeue.ReadState(lab.repo().CacheDir())
	if err != nil {
		lab.t.Fatal(err)
	}
	return st
}

// S-0549 criteria 1, 2, and 3, ADR-0234: a look that finds the queue held by
// its preflight journals the hold once, naming each file and rule and the
// thread, and opens one thread to the operator on the manifest, as flai
// serve, naming each file, line, and rule in text the project's markdown
// lint passes, saying the queue builds no batch until they are fixed and
// changes none of them, and records it on the hold; a second look opens none
// and journals nothing. The look that finds the hold gone resolves the
// thread, saying main is clean, and journals the clearing, once.
func TestAPreflightHoldIsJournalledOnceWithOneThreadToTheOperatorResolvedWhenItClears(t *testing.T) {
	lab := newQueueLab(t)
	if err := os.WriteFile(filepath.Join(lab.root, ".markdownlint.yaml"), []byte(projectLint), 0o644); err != nil {
		t.Fatal(err)
	}
	var held atomic.Bool
	held.Store(true)
	since := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	lab.holdingRun(&held, since)
	lab.look()
	all := lab.holdThreads()
	if len(all) != 1 || !all[0].Open() || all[0].Opener() != queueAuthor || all[0].Title != holdTitle {
		t.Fatalf("threads on the manifest = %+v, want one open, from flai serve, titled %q:\n%s", all, holdTitle, lab.logs.String())
	}
	th := all[0]
	body := th.Marshal()
	for _, want := range []string{"`wip/kanban/tasks/T-0001-x.md`, line 12, `MD051`: `` Link fragments should be valid [Context: \"[see](#nowhere)\"] at https://example.com `x` ``",
		"`wip/kanban/stories/S-0001-x.md`, `item.front-matter`: `status is not one of the statuses` (flai check)", "at main's head `0123456789ab`",
		"builds no batch", "changes none of these files", "clears the hold, resolves this thread"} {
		if !strings.Contains(body, want) {
			t.Errorf("the thread does not say %q:\n%s", want, body)
		}
	}
	if got := lab.state().Preflight; got == nil || got.Thread != th.ID {
		t.Fatalf("the hold = %+v, want its thread %s recorded", got, th.ID)
	}
	got := lab.entries()
	if len(got) != 1 || got[0].Action != hostapi.ActionMergeQueue || got[0].By != queueAuthor || got[0].At != "2026-11-01T09:00:00Z" ||
		!strings.Contains(got[0].Detail, "wip/kanban/tasks/T-0001-x.md (MD051); wip/kanban/stories/S-0001-x.md (item.front-matter)") || !strings.Contains(got[0].Detail, "thread "+th.ID) {
		t.Fatalf("journal = %+v, want the hold, with each file and rule and the thread", got)
	}

	lab.look()
	if all := lab.holdThreads(); len(all) != 1 {
		t.Fatalf("threads on the manifest after a second held look = %d, want the one", len(all))
	}
	if got := lab.entries(); len(got) != 1 {
		t.Fatalf("journal after a second held look = %+v, want the hold alone", got)
	}

	held.Store(false)
	lab.look()
	all = lab.holdThreads()
	if len(all) != 1 || all[0].Open() || !strings.Contains(all[0].Marshal(), "Resolved: main is clean: the merge queue's preflight found nothing at ") {
		t.Fatalf("threads on the manifest once the hold cleared = %+v, want the one resolved as main is clean", all)
	}
	got = lab.entries()
	if len(got) != 2 || !strings.Contains(got[1].Detail, "merge queue preflight hold cleared: main is clean") || !strings.Contains(got[1].Detail, "thread "+th.ID+" resolved") {
		t.Fatalf("journal = %+v, want the hold, then its clearing with the thread", got)
	}
	lab.look()
	if got := lab.entries(); len(got) != 2 {
		t.Fatalf("journal after a clean look = %+v, want nothing new", got)
	}
}

// S-0549, ADR-0234: a thread that cannot be opened is logged at warn and
// fails no look; the hold is journalled once, saying no thread is open yet,
// and the next look opens the thread and journals nothing again. A hold whose
// thread was resolved has a new one opened.
func TestAPreflightThreadNotOpenedIsTriedAgainAtTheNextLook(t *testing.T) {
	lab := newQueueLab(t)
	threadsDir := filepath.Join(lab.root, "wip", "threads")
	if err := os.WriteFile(threadsDir, []byte("not a folder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var held atomic.Bool
	held.Store(true)
	lab.holdingRun(&held, time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC))
	lab.look()
	if recs := lab.records("msg", "merge queue preflight thread not opened; the next look tries again"); len(recs) != 1 || recs[0]["level"] != "WARN" {
		t.Fatalf("logged the thread not opened as %v, want one warning:\n%s", recs, lab.logs.String())
	}
	if got := lab.entries(); len(got) != 1 || !strings.Contains(got[0].Detail, "no thread to the operator yet") {
		t.Fatalf("journal = %+v, want the hold, with no thread yet", got)
	}
	if err := os.Remove(threadsDir); err != nil {
		t.Fatal(err)
	}
	lab.look()
	all := lab.holdThreads()
	if len(all) != 1 || lab.state().Preflight.Thread != all[0].ID {
		t.Fatalf("threads on the manifest = %+v, hold %+v; want one, recorded on the hold", all, lab.state().Preflight)
	}
	if got := lab.entries(); len(got) != 1 {
		t.Fatalf("journal after the thread opened = %+v, want the hold alone", got)
	}
	if _, err := threads.Resolve(lab.repo(), all[0].ID, "alex", "seen", time.Now()); err != nil {
		t.Fatal(err)
	}
	lab.look()
	all = lab.holdThreads()
	if len(all) != 2 || !all[1].Open() || lab.state().Preflight.Thread != all[1].ID {
		t.Fatalf("threads on the manifest = %+v, want a second, open and recorded, once the first was resolved", all)
	}
}

// S-0549 criterion 3, ADR-0234: a batch whose landing waits on main is
// journalled once, with the story and the reason, and its landing again,
// once main is clean, is journalled once, with the story it resumed at, how
// long it waited, the commits, and that no tier ran.
func TestALandingWaitingOnMainAndItsLandingAgainAreJournalledOnce(t *testing.T) {
	lab := newQueueLab(t)
	var held atomic.Bool
	lab.holdingRun(&held, time.Now())
	stalledAt := time.Now().Add(-time.Minute).UTC().Truncate(time.Second)
	if _, err := mergequeue.Update(lab.repo().CacheDir(), func(s *mergequeue.State) error {
		s.Batch = &mergequeue.Batch{ID: "b1", Stories: []string{"S-0001", "S-0002"}, Step: mergequeue.StepVerified, Result: &mergequeue.Result{Passed: true},
			Heads: []mergequeue.Head{{Story: "S-0001", Verified: "1111111111111111", Rebased: "aaaaaaaaaaaaaaaa"}, {Story: "S-0002", Verified: "2222222222222222", Rebased: "bbbbbbbbbbbbbbbb"}}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	stall := func(_ context.Context, repo *workitem.Repo) {
		if _, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
			s.Batch.Landing = &mergequeue.Landing{Stories: []string{"S-0001", "S-0002"}, Landed: []string{"S-0001"}}
			s.Batch.Stalled = &mergequeue.Stalled{Story: "S-0002", Reason: "flai accept S-0002: exit status 1", At: stalledAt}
			s.Preflight = &mergequeue.PreflightHold{Since: stalledAt, Findings: holdFindings}
			return nil
		}); err != nil {
			t.Error(err)
		}
		held.Store(true)
	}
	lab.advances.then = stall
	lab.look()
	lab.advances.then = nil
	lab.look()
	var stalls []hostapi.Entry
	for _, e := range lab.entries() {
		if strings.Contains(e.Detail, "waits on main") {
			stalls = append(stalls, e)
		}
	}
	if len(stalls) != 1 || stalls[0].Outcome != "failed" || !strings.Contains(stalls[0].Detail, "landing S-0002 from merge queue batch b1 waits on main: flai accept S-0002: exit status 1") {
		t.Fatalf("journal = %+v, want the wait on main once, failed, with the story and the reason", lab.entries())
	}
	if recs := lab.records("stalled", "S-0002"); len(recs) != 1 || recs[0]["msg"] != "merge queue batch advanced" {
		t.Fatalf("logged the wait as %v, want one line:\n%s", recs, lab.logs.String())
	}

	before := len(lab.entries())
	held.Store(false)
	lab.advances.then = func(_ context.Context, repo *workitem.Repo) {
		if _, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
			s.RecordOutcome(mergequeue.Outcome{Story: "S-0002", Outcome: mergequeue.OutcomeLanded, At: time.Now(), Detail: "landed by merge queue batch b1 as alex"})
			s.Batch = nil
			return nil
		}); err != nil {
			t.Error(err)
		}
	}
	lab.look()
	lab.advances.then = nil
	lab.look()
	got := lab.entries()[before:]
	if len(got) != 3 || got[0].Outcome != "done" ||
		!strings.Contains(got[0].Detail, "merge queue batch b1 lands again at S-0002 after waiting 1m") ||
		!strings.Contains(got[0].Detail, "(S-0002 at bbbbbbbbbbbb)") || !strings.Contains(got[0].Detail, "with no new verification: no tier ran") ||
		got[1].Detail != "S-0002: landed by merge queue batch b1 as alex" || !strings.Contains(got[2].Detail, "merge queue preflight hold cleared") {
		t.Fatalf("journal after main came clean = %+v, want the landing again, the story landed, and the hold cleared, once each", got)
	}
}

// S-0549 criterion 2, ADR-0234: while the preflight holds the queue, a
// change to any file of the main checkout the watcher sees, outside kanban
// and threads too, and a move of main's head wake the pass; while it does
// not, only a change under kanban or threads does.
func TestAChangeInTheMainCheckoutWakesTheQueueWhileTheHoldStands(t *testing.T) {
	lab := newQueueLab(t)
	woke := func(rel string) bool {
		t.Helper()
		lab.pass.changed(rel)
		select {
		case <-lab.pass.woken:
			return true
		default:
			return false
		}
	}
	if woke("wip/agents/S-0001.md") || woke("docs/x.md") || woke(".git/refs/heads/main") {
		t.Fatal("a change outside kanban and threads woke a queue not held")
	}
	if !woke("wip/kanban/tasks/T-0001-x.md") || !woke("wip/threads/TH-0001-x.md") {
		t.Fatal("a change under kanban or threads did not wake the queue")
	}
	var held atomic.Bool
	held.Store(true)
	lab.holdingRun(&held, time.Now())
	lab.look()
	for _, rel := range []string{"wip/agents/S-0001.md", "docs/x.md", ".git/refs/heads/main", "wip/kanban/tasks/T-0001-x.md"} {
		if !woke(rel) {
			t.Errorf("a change to %s did not wake the queue the preflight holds", rel)
		}
	}
	held.Store(false)
	lab.look()
	if woke("wip/agents/S-0001.md") {
		t.Fatal("a change outside kanban and threads woke the queue once its hold cleared")
	}
}

// S-0549 criterion 2, ADR-0234: a running pass watches the main checkout's
// refs, and a move of main's head has it look again while held, and not once
// the hold cleared.
func TestAMoveOfMainsHeadWakesARunningPassWhileHeld(t *testing.T) {
	lab := newQueueLab(t)
	ref := filepath.Join(lab.root, ".git", "refs", "heads", "main")
	if err := os.MkdirAll(filepath.Dir(ref), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ref, []byte("1111111111111111111111111111111111111111\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.pass.every = time.Hour
	lab.pass.o.WatchEvery = 5 * time.Millisecond
	var held atomic.Bool
	held.Store(true)
	lab.holdingRun(&held, time.Now())
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		lab.runs.wait()
	}()
	lab.pass.start(ctx)
	waitFor(t, "the first look to see the hold", lab.pass.held)
	n := lab.batches.count()
	if err := os.WriteFile(ref, []byte("2222222222222222222222222222222222222222\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "a look after main's head moved", func() bool { return lab.batches.count() > n })
	held.Store(false)
	lab.pass.wake()
	waitFor(t, "a look that clears the hold", func() bool { return !lab.pass.held() })
	time.Sleep(50 * time.Millisecond) // the wake's look ends
	n = lab.batches.count()
	if err := os.WriteFile(ref, []byte("3333333333333333333333333333333333333333\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond) // ticks of the watcher
	if got := lab.batches.count(); got != n {
		t.Fatalf("looks after main's head moved with no hold = %d, want %d", got, n)
	}
}

// S-0604: a running pass tells the release package of every move of the main
// checkout's refs, a tag's and a branch's, loose or packed, so that the
// board's pending read takes no git process while none came.
func TestAPassTellsReleaseOfAMovedTagOrBranch(t *testing.T) {
	lab := newQueueLab(t)
	git := filepath.Join(lab.root, ".git")
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(git, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("HEAD", "ref: refs/heads/main\n")
	write("refs/heads/main", "1111111111111111111111111111111111111111\n")
	var moves atomic.Int64
	lab.on = false
	lab.pass.every = time.Hour
	lab.pass.o.WatchEvery = 5 * time.Millisecond
	lab.pass.refsMoved = func(root string) {
		if root == lab.root {
			moves.Add(1)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		lab.runs.wait()
	}()
	lab.pass.start(ctx)
	waitFor(t, "the watcher's word to count from after its first look", func() bool { return moves.Load() >= 1 })
	time.Sleep(50 * time.Millisecond) // ticks of the watcher with nothing moved
	for _, move := range []struct{ what, rel, body string }{
		{"a new tag", "refs/tags/cli/v1.0.0", "2222222222222222222222222222222222222222\n"},
		{"main's head", "refs/heads/main", "3333333333333333333333333333333333333333\n"},
		{"the packed refs", "packed-refs", "4444444444444444444444444444444444444444 refs/tags/cli/v0.9.0\n"},
	} {
		n := moves.Load()
		write(move.rel, move.body)
		waitFor(t, "release told of "+move.what, func() bool { return moves.Load() > n })
	}
	n := moves.Load()
	time.Sleep(100 * time.Millisecond) // ticks of the watcher with nothing moved
	if got := moves.Load(); got != n {
		t.Errorf("moves told with nothing moved = %d, want %d", got-n, 0)
	}
}

// fakeProcess stands in for runCommand, keeping what it was asked to run
// and answering stderr and err.
type fakeProcess struct {
	dir, name string
	env, args []string
	stderr    string
	err       error
}

func (f *fakeProcess) run(_ context.Context, dir string, env []string, name string, args ...string) ([]byte, error) {
	f.dir, f.env, f.name, f.args = dir, env, name, args
	return []byte(f.stderr), f.err
}

// S-0400 criterion 1 (ADR-0176): the Accepter runs this flai's own binary as
// flai accept <story> --by <by> --batch <batch> in the main checkout, as
// flai serve and in no agent's run; a failure carries the command and its
// last lines on standard error.
func TestTheAccepterRunsFlaiAcceptWithTheBatch(t *testing.T) {
	t.Setenv("FLAI_AGENT", "agent-S-0400")
	t.Setenv("FLAI_ROLE", "orchestrate")
	t.Setenv("FLAI_STORY", "S-0400")
	proc := &fakeProcess{}
	accept := flaiAccepter(proc.run, func() (string, error) { return "/opt/flai/bin/flai", nil })("/home/p")
	if err := accept(context.Background(), "S-0401", "alex", "20261019-041200"); err != nil {
		t.Fatal(err)
	}
	if proc.dir != "/home/p" || proc.name != "/opt/flai/bin/flai" || strings.Join(proc.args, " ") != "accept S-0401 --by=alex --batch=20261019-041200" {
		t.Fatalf("ran %s %v in %s, want this flai's flai accept with the batch in the main checkout", proc.name, proc.args, proc.dir)
	}
	env := map[string]string{}
	for _, kv := range proc.env {
		k, v, _ := strings.Cut(kv, "=")
		env[k] = v
	}
	if env["FLAI_AGENT"] != queueAgent || env["FLAI_SESSION"] != "" || env["FLAI_ROLE"] != "" || env["FLAI_STORY"] != "" || env["PATH"] == "" {
		t.Fatalf("env = %v, want FLAI_AGENT %s, no session, no run marks, and flai serve's own", proc.env, queueAgent)
	}

	proc.stderr, proc.err = "one\ntwo\n\nthree\nfour\nfive\nError: rule: S-0401 is not the next story\n", errors.New("exit status 1")
	err := accept(context.Background(), "S-0401", "alex", "b1")
	if err == nil || !strings.Contains(err.Error(), "flai accept S-0401 --by=alex --batch=b1: exit status 1: two; three; four; five; Error: rule: S-0401 is not the next story") ||
		strings.Contains(err.Error(), "one") {
		t.Fatalf("err = %v, want the command, its exit, and its last five lines", err)
	}

	broken := flaiAccepter(proc.run, func() (string, error) { return "", errors.New("no /proc") })("/home/p")
	if err := broken(context.Background(), "S-0401", "alex", "b1"); err == nil || !strings.Contains(err.Error(), "no /proc") {
		t.Fatalf("err = %v, want why the executable is not known", err)
	}
}

// S-0400: runCommand runs the process in its folder with its environment
// and answers what it wrote on standard error with its exit.
func TestRunCommandAnswersStandardErrorAndTheExit(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	dir := t.TempDir()
	stderr, err := runCommand(context.Background(), dir, []string{"X=from-env"}, sh, "-c", `echo "$X in $(pwd)" >&2; exit 3`)
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 3 {
		t.Fatalf("err = %v, want exit status 3", err)
	}
	if got := strings.TrimSpace(string(stderr)); !strings.HasPrefix(got, "from-env in ") || !strings.HasSuffix(got, filepath.Base(dir)) {
		t.Fatalf("stderr = %q, want the variable and the folder", got)
	}
}

// S-0400 criteria 1, 3, and 4, through a pass on a real git project: one look
// builds and verifies a batch of two stories with RunBatch, which fails with
// the second, and mergequeue.Advance bisects it to that story, sends it back
// to in-progress by flai with a thread from flai serve, lands the first with
// the Accepter as whoever enqueued it, and leaves the queue with no batch.
// The tiers are a fake that fails wherever b.txt is. It skips under -short.
func TestAPassLandsAndSendsBackThroughTheRealAdvance(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	lab := newQueueLab(t)
	git := func(dir string, args ...string) string {
		t.Helper()
		out, err := (execx.System{}).Run(dir, "git", args...)
		if err != nil {
			t.Fatal(out, err)
		}
		return strings.TrimSpace(out)
	}
	write := func(dir, rel, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, rel), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write(lab.root, ".gitignore", ".flai-cache/\nwip/\n")
	write(lab.root, "x.txt", "one\n")
	git(lab.root, "init", "-q", "-b", "main")
	git(lab.root, "config", "user.email", "t@t")
	git(lab.root, "config", "user.name", "t")
	git(lab.root, "config", "commit.gpgsign", "false")
	git(lab.root, "add", "-A")
	git(lab.root, "commit", "-q", "-m", "init")
	repo := lab.repo()
	review := time.Now().Add(-time.Hour)
	var stories []string
	for _, file := range []string{"a.txt", "b.txt"} {
		s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Add " + file, Nature: "feature", Owner: "alex", Now: review})
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(s.Path)
		if err != nil {
			t.Fatal(err)
		}
		write(filepath.Dir(s.Path), filepath.Base(s.Path), strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [x] it adds "+file+"\n", 1))
		if s, err = repo.Get(s.ID); err != nil {
			t.Fatal(err)
		}
		for _, to := range []string{workitem.Ready, workitem.InProgress} {
			if _, err := repo.Transition(s, to, "alex", "", review); err != nil {
				t.Fatal(err)
			}
		}
		task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Write " + file, Parent: s.ID, Owner: "alex", Now: review})
		if err != nil {
			t.Fatal(err)
		}
		for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
			if _, err := repo.Transition(task, to, "alex", "", review); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := repo.Transition(s, workitem.Review, "alex", "", review); err != nil {
			t.Fatal(err)
		}
		// a story in review has its narrative, or flai check's narrative.missing
		// error on main holds the queue in the landing's preflight (ADR-0234)
		if _, err := repo.OpenStream(s, workitem.StreamOptions{Agent: "agent-" + s.ID, Now: review}); err != nil {
			t.Fatal(err)
		}
		wt := repo.WorktreePath(s.ID)
		git(lab.root, "worktree", "add", "-q", "-b", storygit.Branch(s.ID), wt, "main")
		write(wt, file, file+"\n")
		git(wt, "add", "-A")
		git(wt, "commit", "-q", "-m", s.ID+"'s change")
		if err := verify.SaveReport(repo, verify.Report{Story: s.ID, Commit: git(wt, "rev-parse", "HEAD"), Passed: true}); err != nil {
			t.Fatal(err)
		}
		stories = append(stories, s.ID)
	}
	if _, err := mergequeue.Update(repo.CacheDir(), func(s *mergequeue.State) error {
		for _, id := range stories {
			s.Enqueue(id, "alex", review)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	tiers := func(_ context.Context, dir, _, _ string, _ []string, started func(string)) (verify.Report, error) {
		started("unit")
		if _, err := os.Stat(filepath.Join(dir, "b.txt")); err == nil {
			return verify.Report{StoppedAt: "unit", Steps: []verify.Step{{Name: "unit", Tier: true, State: verify.Failed}}}, nil
		}
		return verify.Report{Passed: true, Steps: []verify.Step{{Name: "unit", Tier: true, State: verify.Passed}}}, nil
	}
	lab.pass.run = func(ctx context.Context, d mergequeue.Deps, repo *workitem.Repo) (bool, error) {
		d.Verify = tiers
		return mergequeue.RunBatch(ctx, d, repo)
	}
	lab.pass.advance = func(ctx context.Context, d mergequeue.Deps, accept mergequeue.Accepter, repo *workitem.Repo) (bool, error) {
		d.Verify = tiers
		return mergequeue.Advance(ctx, d, accept, repo)
	}
	var accepted []string
	lab.pass.accepter = func(dir string) mergequeue.Accepter {
		return func(_ context.Context, story, by, batch string) error {
			accepted = append(accepted, story+" "+by+" "+dir)
			// what flai accept --batch records once it has merged the story
			_, err := mergequeue.Landed(repo.CacheDir(), story, batch, by, time.Now())
			return err
		}
	}
	lab.look()
	if want := []string{stories[0] + " alex " + lab.root}; !slices.Equal(accepted, want) {
		t.Fatalf("accepted %v, want %v: the story before the failing one, as its enqueuer, in the main checkout", accepted, want)
	}
	failing, err := repo.Get(stories[1])
	if err != nil {
		t.Fatal(err)
	}
	if last := failing.Transitions[len(failing.Transitions)-1]; failing.Status != workitem.InProgress || last.By != "flai" {
		t.Fatalf("%s is %s, last moved %+v; want it sent back to in-progress by flai", failing.ID, failing.Status, last)
	}
	if data, err := os.ReadFile(failing.Path); err != nil || !strings.Contains(string(data), "fails with "+failing.ID) {
		t.Fatalf("%s's notes have no reason for the send-back (%v):\n%s", failing.ID, err, data)
	}
	all, err := threads.For(repo, failing.ID)
	if err != nil || len(all) != 1 || !all[0].Open() || all[0].Opener() != queueAuthor {
		t.Fatalf("threads on %s = %v, %v; want one open, from flai serve", failing.ID, all, err)
	}
	st, err := mergequeue.ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if st.Batch != nil || len(st.Enqueued) != 0 {
		t.Fatalf("state = %+v, want no batch and nothing enqueued", st)
	}
	recs := lab.records("msg", "merge queue batch advanced")
	if len(recs) != 1 || recs[0]["landed"] != stories[0] || recs[0]["sent_back"] != stories[1] {
		t.Fatalf("logged the advance as %v, want the story landed and the one sent back:\n%s", recs, lab.logs.String())
	}
	if recs := lab.records("level", "ERROR"); len(recs) != 0 {
		t.Fatalf("logged a failure:\n%s", lab.logs.String())
	}
	if got := lab.entries(); len(got) != 2 {
		t.Fatalf("journal = %+v, want the landing and the send-back", got)
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

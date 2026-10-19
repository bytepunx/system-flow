package serve

import (
	"context"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The merge queue, run for a project with no model (S-0399, ADR-0173).
//
// A queue pass per served project looks every queueEvery, and at once when
// the project's work items or threads change. While the merge-queue host
// action is on, orchestration.permissions.accept_reviews is on, acceptance is
// local, and the queue is not held, it builds and verifies one batch with
// mergequeue.RunBatch, in its own goroutine, so that no agent's turn and no
// request waits on a batch's tiers. A batch already running finishes when the
// queue is held; the next look sees the hold and starts none. Why the queue
// does not run is logged once, when it changes, not at every look. When flai
// serve stops, or stops serving the project, the batch running is cancelled
// and its worktree left for the next run to build again.

// queueEvery is how often a project's merge queue is looked at when nothing
// changed: a story enqueued waits at most this long for a batch to start.
const queueEvery = 30 * time.Second

// Why a project's merge queue does not run, as the reason field of the log
// line that says so.
const (
	queueNoAcceptReviews = "accept_reviews_off"
	queuePullRequest     = "pull_request"
	queueHeld            = "held"
	queueUnreadable      = "unreadable"
)

// batchRunner builds and verifies one batch of repo's merge queue, as
// mergequeue.RunBatch does.
type batchRunner func(ctx context.Context, d mergequeue.Deps, repo *workitem.Repo) (bool, error)

// queueRuns are the merge queue passes of one flai serve: the projects whose
// batch is running, so that no two run for one project at once, even while a
// pass that stopped serving a project still finishes its batch, and the
// passes to wait for when flai serve stops.
type queueRuns struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	busy map[string]bool
}

func newQueueRuns() *queueRuns {
	return &queueRuns{busy: map[string]bool{}}
}

// claim marks root's batch as running and reports whether it was not already.
func (q *queueRuns) claim(root string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.busy[root] {
		return false
	}
	q.busy[root] = true
	return true
}

// release marks root's batch as no longer running.
func (q *queueRuns) release(root string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.busy, root)
}

// wait waits for every pass started to end.
func (q *queueRuns) wait() { q.wg.Wait() }

// queuePass runs the merge queue of one served project.
type queuePass struct {
	o     Options
	e     Entry
	runs  *queueRuns
	every time.Duration
	run   batchRunner
	woken chan struct{}

	mu sync.Mutex
	// said is the reason last logged for the queue not running, "" when it
	// last could run or the action was off.
	said string
}

// newQueuePass is project e's merge queue pass, one of runs.
func newQueuePass(o Options, e Entry, runs *queueRuns) *queuePass {
	return &queuePass{o: o, e: e, runs: runs, every: queueEvery, run: mergequeue.RunBatch, woken: make(chan struct{}, 1)}
}

// start looks now, then every p.every and whenever woken, in the background
// until ctx ends; a batch running is cancelled with ctx.
func (p *queuePass) start(ctx context.Context) {
	p.runs.wg.Add(1)
	go func() {
		defer p.runs.wg.Done()
		for {
			p.look(ctx)
			select {
			case <-ctx.Done():
				return
			case <-p.woken:
			case <-time.After(p.every):
			}
		}
	}()
}

// wake has the pass look again as soon as it is not running a batch; it
// never waits.
func (p *queuePass) wake() {
	select {
	case p.woken <- struct{}{}:
	default:
	}
}

// look runs one batch of the project's merge queue when the action, the
// permission, the acceptance mode, and the hold let it, and logs what came
// of it. A batch already running for the project is left to finish.
func (p *queuePass) look(ctx context.Context) {
	if p.o.Host.Enabled == nil || !p.o.Host.Enabled(hostapi.ActionMergeQueue, p.e.Root) {
		p.mu.Lock()
		p.said = ""
		p.mu.Unlock()
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.notRun(queueUnreadable, "flai serve could not open the project, so its merge queue starts no batch: fix what the error names", err)
		return
	}
	switch {
	case !repo.Manifest.Orchestration.Permissions.AcceptReviews:
		p.notRun(queueNoAcceptReviews, "orchestration.permissions.accept_reviews is off in system-flow.yaml, so the merge queue starts no batch: turn it on to let the queue run", nil)
		return
	case repo.Manifest.Acceptance.PullRequest():
		p.notRun(queuePullRequest, mergequeue.OffUnderPullRequest, nil)
		return
	}
	st, err := mergequeue.ReadState(repo.CacheDir())
	if err != nil {
		p.notRun(queueUnreadable, "flai serve could not read the merge queue's state, so it starts no batch: fix what the error names", err)
		return
	}
	if st.Held {
		p.notRun(queueHeld, "the merge queue is held, so it starts no batch: flai queue resume lets it build batches again", nil, "held_by", st.HeldBy)
		return
	}
	p.mayRun()
	if !p.runs.claim(p.e.Root) {
		return
	}
	defer p.runs.release(p.e.Root)
	git, started := gitOf(p.o), p.o.Now()
	ran, err := p.run(ctx, mergequeue.Deps{Runner: git, Now: p.o.Now, Verify: mergequeue.CheckoutVerifier(git, p.o.Version),
		Log: p.o.Logger.With("root", p.e.Root)}, repo)
	took := p.o.Now().Sub(started).Round(time.Second).String()
	switch {
	case err != nil && ctx.Err() != nil:
		p.o.Logger.Info("merge queue batch stopped; the next run builds it again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took, "err", err.Error())
	case err != nil:
		p.o.Logger.Error("merge queue batch failed", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took, "err", err.Error())
	case ran:
		p.o.Logger.Info("merge queue batch run", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took)
	}
}

// notRun logs once, until the reason changes, that the queue does not run
// and why: a warning with err when the project or its queue could not be
// read, a lifecycle event otherwise.
func (p *queuePass) notRun(reason, detail string, err error, fields ...any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.said == reason {
		return
	}
	p.said = reason
	args := append([]any{"component", "serve", "root", p.e.Root, "key", p.e.Key, "reason", reason, "detail", detail}, fields...)
	if err != nil {
		p.o.Logger.Warn("merge queue not run", append(args, "err", err.Error())...)
		return
	}
	p.o.Logger.Info("merge queue not run", args...)
}

// mayRun logs, once, that a queue said not to run may run again.
func (p *queuePass) mayRun() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.said == "" {
		return
	}
	p.o.Logger.Info("merge queue may run again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "was", p.said)
	p.said = ""
}

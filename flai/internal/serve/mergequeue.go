package serve

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
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
//
// After the batch's turn, in the same goroutine and under the same claim, the
// pass takes a verified batch on with mergequeue.Advance (S-0400, ADR-0176):
// it lands each story of a batch that passed, in queue order, as whoever
// enqueued it, each through flai accept <story> --by <by> --batch <batch>, a
// process of this flai's own binary in the main checkout; it bisects a batch
// that failed to the first story with which it fails, sends that story back
// to in-progress with a reason and a thread to its agent, and lands the
// stories before it; and it drops a batch when the main branch moved other
// than by the batch's landing, so that the next look builds it again from the
// new main. What came of it is logged once a look and kept in the host's
// journal, a story landed or sent back and a batch dropped an entry each.

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

// queueAuthor is who flai serve's merge queue acts as in the host's journal;
// queueAgent is FLAI_AGENT in the flai accept a landing runs, naming flai
// serve as the dashboard's acceptance names flaiover.
const (
	queueAuthor = "flai serve"
	queueAgent  = "flai-serve"
)

// acceptSaid is how many of a failed landing's last lines on standard error
// its error carries.
const acceptSaid = 5

// batchRunner builds and verifies one batch of repo's merge queue, as
// mergequeue.RunBatch does.
type batchRunner func(ctx context.Context, d mergequeue.Deps, repo *workitem.Repo) (bool, error)

// batchAdvancer takes repo's verified merge queue batch one step on, landing
// its stories with accept, as mergequeue.Advance does.
type batchAdvancer func(ctx context.Context, d mergequeue.Deps, accept mergequeue.Accepter, repo *workitem.Repo) (bool, error)

// commandRunner runs name with args in dir, with env as its whole
// environment, until it ends or ctx does, and answers what it wrote on
// standard error, with an error when it could not run or exited non-zero.
type commandRunner func(ctx context.Context, dir string, env []string, name string, args ...string) ([]byte, error)

// runCommand is the commandRunner that runs the process: ctx ending
// interrupts it, and kills it ten seconds later, as hostapi.ExecRunner does.
func runCommand(ctx context.Context, dir string, env []string, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir, cmd.Env = dir, env
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	cmd.WaitDelay = 10 * time.Second
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stderr.Bytes(), err
}

// flaiAccepter is, for the main checkout dir, the Accepter flai serve lands
// a story with (ADR-0176): exe, this flai's own executable, run with run as
// flai accept <story> --by=<by> --batch=<batch> in dir, with FLAI_AGENT
// naming flai serve, no session, and none of the variables that say what an
// agent's run is for. Its error, when the acceptance fails, carries the
// command and its last lines on standard error.
func flaiAccepter(run commandRunner, exe func() (string, error)) func(dir string) mergequeue.Accepter {
	return func(dir string) mergequeue.Accepter {
		return func(ctx context.Context, story, by, batch string) error {
			flai, err := exe()
			if err != nil {
				return fmt.Errorf("this flai's executable is not known: %w", err)
			}
			args := []string{"accept", story, "--by=" + by, "--batch=" + batch}
			stderr, err := run(ctx, dir, acceptEnv(os.Environ()), flai, args...)
			if err == nil {
				return nil
			}
			if said := lastLines(string(stderr), acceptSaid); said != "" {
				return fmt.Errorf("flai %s: %w: %s", strings.Join(args, " "), err, said)
			}
			return fmt.Errorf("flai %s: %w", strings.Join(args, " "), err)
		}
	}
}

// acceptEnv is env as a landing's flai accept has it: without the variables
// that say what an agent's run is for, FLAI_AGENT naming flai serve, and no
// session.
func acceptEnv(env []string) []string {
	out := slices.DeleteFunc(slices.Clone(env), func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return name == "FLAI_AGENT" || name == "FLAI_SESSION" || slices.Contains(runMarks, name)
	})
	return append(out, "FLAI_AGENT="+queueAgent, "FLAI_SESSION=")
}

// lastLines are the last n lines of s that are not blank, joined with "; ".
func lastLines(s string, n int) string {
	var lines []string
	for line := range strings.Lines(s) {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines[max(0, len(lines)-n):], "; ")
}

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
	// advance takes the verified batch on after run; accepter is the
	// Accepter it lands a story with in the main checkout dir.
	advance  batchAdvancer
	accepter func(dir string) mergequeue.Accepter
	woken    chan struct{}

	mu sync.Mutex
	// said is the reason last logged for the queue not running, "" when it
	// last could run or the action was off.
	said string
}

// newQueuePass is project e's merge queue pass, one of runs.
func newQueuePass(o Options, e Entry, runs *queueRuns) *queuePass {
	return &queuePass{o: o, e: e, runs: runs, every: queueEvery, run: mergequeue.RunBatch, advance: mergequeue.Advance,
		accepter: flaiAccepter(runCommand, os.Executable), woken: make(chan struct{}, 1)}
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
// permission, the acceptance mode, and the hold let it, then takes the
// verified batch on, and logs what came of each. A batch already running for
// the project is left to finish.
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
	d := mergequeue.Deps{Runner: git, Now: p.o.Now, Verify: mergequeue.CheckoutVerifier(git, p.o.Version), Log: p.o.Logger.With("root", p.e.Root)}
	// a story the queue sent back that is in review again at a new commit
	// has its send back's thread resolved, so that it can be batched again
	if resolved, err := mergequeue.ResolveSentBack(d, repo); err != nil {
		p.o.Logger.Warn("merge queue send back threads not resolved; the next look tries again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "err", err.Error())
	} else if len(resolved) > 0 {
		p.o.Logger.Info("merge queue send back threads resolved", "component", "serve", "root", p.e.Root, "key", p.e.Key, "stories", resolved)
	}
	ran, err := p.run(ctx, d, repo)
	took := p.o.Now().Sub(started).Round(time.Second).String()
	switch {
	case err != nil && ctx.Err() != nil:
		p.o.Logger.Info("merge queue batch stopped; the next run builds it again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took, "err", err.Error())
		return
	case err != nil:
		p.o.Logger.Error("merge queue batch failed", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took, "err", err.Error())
	case ran:
		p.o.Logger.Info("merge queue batch run", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took)
	}
	// a batch verified now, or at an earlier look, waits here to be landed or
	// bisected: RunBatch begins none while it does
	p.advanceBatch(ctx, d, repo)
}

// advanced is what one advance did to the queue, read from its state before
// and after: the batch it took on, the stories it landed or sent back, each
// with its outcome, and why it dropped the batch, "" when it did not.
type advanced struct {
	batch    string
	outcomes []mergequeue.Outcome
	dropped  string
}

// advancedBetween is what changed from was to now: the outcomes landed or
// failing that now has and was had not, and the drop of was's batch.
func advancedBetween(was, now mergequeue.State) advanced {
	var a advanced
	if b := was.Batch; b != nil {
		a.batch = b.ID
		if n := now.Batch; n != nil && n.ID == b.ID && n.Dropped != nil && b.Dropped == nil {
			a.dropped = n.Dropped.Reason
		}
	}
	for _, o := range now.Outcomes {
		if o.Outcome != mergequeue.OutcomeLanded && o.Outcome != mergequeue.OutcomeFailing {
			continue
		}
		if !slices.ContainsFunc(was.Outcomes, func(x mergequeue.Outcome) bool {
			return x.Story == o.Story && x.Outcome == o.Outcome && x.At.Equal(o.At)
		}) {
			a.outcomes = append(a.outcomes, o)
		}
	}
	return a
}

// stories are the stories of a's outcomes that came to outcome, joined with
// commas.
func (a advanced) stories(outcome string) string {
	var ids []string
	for _, o := range a.outcomes {
		if o.Outcome == outcome {
			ids = append(ids, o.Story)
		}
	}
	return strings.Join(ids, ",")
}

// advanceBatch takes the project's verified batch on with p.advance, landing
// with p.accepter in the main checkout, logs what came of it once, journals
// each story landed or sent back and the batch dropped, and has the pass look
// again at once when it did anything and failed in nothing, so that the
// stories still enqueued are batched without waiting for the next tick; after
// a failure the next tick builds the batch again.
func (p *queuePass) advanceBatch(ctx context.Context, d mergequeue.Deps, repo *workitem.Repo) {
	was, _ := mergequeue.ReadState(repo.CacheDir()) // Advance reads it again, and says when it cannot
	started := p.o.Now()
	did, err := p.advance(ctx, d, p.accepter(repo.MainRoot), repo)
	took := p.o.Now().Sub(started).Round(time.Second).String()
	if err != nil && ctx.Err() != nil {
		p.o.Logger.Info("merge queue batch stopped; the next run takes it on again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "took", took, "err", err.Error())
		return
	}
	now, rerr := mergequeue.ReadState(repo.CacheDir())
	if rerr != nil {
		now = was
	}
	a := advancedBetween(was, now)
	p.journal(a, err)
	fields := []any{"component", "serve", "root", p.e.Root, "key", p.e.Key, "batch", a.batch,
		"landed", a.stories(mergequeue.OutcomeLanded), "sent_back", a.stories(mergequeue.OutcomeFailing), "took", took}
	switch {
	case err != nil && a.dropped != "":
		p.o.Logger.Warn("merge queue batch dropped", append(fields, "reason", a.dropped, "err", err.Error())...)
	case err != nil:
		p.o.Logger.Error("merge queue batch not advanced", append(fields, "err", err.Error())...)
	case did:
		p.o.Logger.Info("merge queue batch advanced", append(fields, "dropped", a.dropped)...)
	}
	if did && err == nil {
		p.wake()
	}
}

// journal keeps in the host's journal each story a landed or sent back, as
// done, and the batch it dropped, as failed when err stopped the advance.
func (p *queuePass) journal(a advanced, err error) {
	if p.o.Host.Record == nil {
		return
	}
	entry := hostapi.Entry{Action: hostapi.ActionMergeQueue, Method: "serve.merge-queue", Project: p.e.Key, Root: p.e.Root, By: queueAuthor, Outcome: "done"}
	for _, o := range a.outcomes {
		entry.At, entry.Detail = o.At.UTC().Format(time.RFC3339), o.Story+": "+o.Detail
		p.o.Host.Record(entry)
	}
	if a.dropped == "" {
		return
	}
	entry.At, entry.Detail = p.o.Now().UTC().Format(time.RFC3339), fmt.Sprintf("dropped merge queue batch %s: %s", a.batch, a.dropped)
	if err != nil {
		entry.Outcome = "failed"
	}
	p.o.Host.Record(entry)
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

package serve

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/watch"
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
//
// While the landing's preflight holds the queue on what it found on the main
// checkout (ADR-0234), the pass keeps one thread to the operator open on the
// manifest, naming each file and finding, journals the hold once, and, at the
// look that finds main clean, resolves the thread and journals the clearing.
// A batch whose landing waits on main, and its landing again with no new
// verification, are journalled once each. While held, or while the landing
// waits, any change the watcher sees in the main checkout, and a move of
// main's head, wake the pass, so that the operator's fix resumes it within
// seconds.

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

// holdTitle is the title of the one thread flai serve opens to the operator
// for the preflight's hold, which is how a later look finds it again.
const holdTitle = "The merge queue is held: the main checkout has files a landing would fail on"

// mainRefs are the files of the main checkout's git that say where its
// branches' heads and its tags are, relative to its root: the pass watches
// them, so that a move of main's head wakes it while held, and so that the
// release package reads no git for the board while none moved (S-0604).
var mainRefs = []string{".git/HEAD", ".git/packed-refs", ".git/refs/heads", ".git/refs/tags"}

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
	// refsMoved tells the release package root's refs may have moved,
	// release.RefsMoved but in a test.
	refsMoved func(root string)

	mu sync.Mutex
	// said is the reason last logged for the queue not running, "" when it
	// last could run or the action was off.
	said string
	// holding says the preflight held the queue, or the batch's landing
	// waited on main, when the last look ended; heldSince is when the hold
	// last journalled began; stale says a thread of a hold that cleared may
	// still be open, as at the first look and after a resolve failed.
	holding   bool
	heldSince time.Time
	stale     bool
}

// newQueuePass is project e's merge queue pass, one of runs.
func newQueuePass(o Options, e Entry, runs *queueRuns) *queuePass {
	return &queuePass{o: o, e: e, runs: runs, every: queueEvery, run: mergequeue.RunBatch, advance: mergequeue.Advance,
		accepter: flaiAccepter(runCommand, os.Executable), woken: make(chan struct{}, 1), refsMoved: release.RefsMoved, stale: true}
}

// start looks now, then every p.every and whenever woken, in the background
// until ctx ends; a batch running is cancelled with ctx. It watches the main
// checkout's mainRefs meanwhile, so that a move of main's head wakes it
// while held.
func (p *queuePass) start(ctx context.Context) {
	p.runs.wg.Add(2)
	go func() {
		defer p.runs.wg.Done()
		p.watchRefs(ctx)
	}()
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

// watchRefs watches the main checkout's mainRefs until ctx ends, waking the
// pass on a change (changed). It also tells the release package of every
// move, so that the board's pending read starts no git process while none
// came (S-0604). A root whose .git is not a folder has none of these files to
// watch, and its board keeps asking git.
func (p *queuePass) watchRefs(ctx context.Context) {
	w := &watch.Watcher{Root: p.e.Root, Paths: mainRefs, Every: p.o.WatchEvery}
	emit := p.changed
	if info, err := os.Stat(filepath.Join(p.e.Root, ".git")); err == nil && info.IsDir() {
		root := p.e.Root
		release.Watch(root)
		defer release.Unwatch(root)
		emit = func(rel string) {
			p.refsMoved(root)
			p.changed(rel)
		}
		// The watcher takes its first look as Run starts. A move between Watch
		// and that look would go unseen, so its word counts only from a tick
		// after it.
		every := w.Every
		if every <= 0 {
			every = 300 * time.Millisecond
		}
		armed := time.AfterFunc(every, func() { p.refsMoved(root) })
		defer armed.Stop()
	}
	w.Run(ctx, emit)
}

// wake has the pass look again as soon as it is not running a batch; it
// never waits.
func (p *queuePass) wake() {
	select {
	case p.woken <- struct{}{}:
	default:
	}
}

// changed has the pass look again for a change to the project's file rel,
// relative to its root: one under a kanban or threads folder, where a
// story's state and an answer to a thread are, and, while the preflight
// holds the queue or the batch's landing waits on main, any, so that the
// fix the operator makes resumes it within seconds (ADR-0234).
func (p *queuePass) changed(rel string) {
	slashed := filepath.ToSlash(rel)
	if strings.Contains(slashed, "/kanban/") || strings.Contains(slashed, "/threads/") || p.held() {
		p.wake()
	}
}

// held reports whether the last look ended with the preflight holding the
// queue or the batch's landing waiting on main.
func (p *queuePass) held() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.holding
}

// look runs one batch of the project's merge queue when the action, the
// permission, the acceptance mode, and the hold let it, then takes the
// verified batch on, logs what came of each, and follows the preflight's
// hold. A batch already running for the project is left to finish.
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
	if ctx.Err() == nil {
		p.followHold(repo, st)
	}
}

// advanced is what one advance did to the queue, read from its state before
// and after: the batch it took on, the stories it landed or sent back, each
// with its outcome, and why it dropped the batch, "" when it did not. Stalled
// is the wait on main the batch's landing came to, and resumed the wait it
// went on from, with from the heads of the stories it had still to land;
// each nil when it did not (ADR-0234).
type advanced struct {
	batch    string
	outcomes []mergequeue.Outcome
	dropped  string
	stalled  *mergequeue.Stalled
	resumed  *mergequeue.Stalled
	from     []mergequeue.Head
}

// advancedBetween is what changed from was to now: the outcomes landed or
// failing that now has and was had not, the drop of was's batch, the wait on
// main its landing came to, and the wait it went on from: was's batch
// waited, not dropped, and now it is gone, landed out, or the same batch
// waits no longer or waits on another landing.
func advancedBetween(was, now mergequeue.State) advanced {
	var a advanced
	if b := was.Batch; b != nil {
		a.batch = b.ID
		n := now.Batch
		same := n != nil && n.ID == b.ID
		if same && n.Dropped != nil && b.Dropped == nil {
			a.dropped = n.Dropped.Reason
		}
		if b.Stalled != nil && b.Dropped == nil && (!same || !sameStall(n.Stalled, b.Stalled)) {
			a.resumed, a.from = b.Stalled, unlanded(b)
		}
		if same && n.Stalled != nil && !sameStall(n.Stalled, b.Stalled) {
			a.stalled = n.Stalled
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

// sameStall reports whether x and y are one wait on main: both set, on one
// story, from one moment.
func sameStall(x, y *mergequeue.Stalled) bool {
	return x != nil && y != nil && x.Story == y.Story && x.At.Equal(y.At)
}

// unlanded are the heads of b's stories its landing has not landed yet, in
// the order rebased.
func unlanded(b *mergequeue.Batch) []mergequeue.Head {
	var out []mergequeue.Head
	for _, h := range b.Heads {
		if b.Landing == nil || !slices.Contains(b.Landing.Landed, h.Story) {
			out = append(out, h)
		}
	}
	return out
}

// stallStory is the story s waits at, "" when s is nil.
func stallStory(s *mergequeue.Stalled) string {
	if s == nil {
		return ""
	}
	return s.Story
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
// each story landed or sent back, the batch dropped, its landing's wait on
// main, and its landing again, and has the pass look
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
		"landed", a.stories(mergequeue.OutcomeLanded), "sent_back", a.stories(mergequeue.OutcomeFailing),
		"stalled", stallStory(a.stalled), "resumed", stallStory(a.resumed), "took", took}
	switch {
	case err != nil && a.dropped != "":
		p.o.Logger.Warn("merge queue batch dropped", append(fields, "reason", a.dropped, "err", err.Error())...)
	case err != nil:
		p.o.Logger.Error("merge queue batch not advanced", append(fields, "err", err.Error())...)
	case did || a.stalled != nil || a.resumed != nil:
		p.o.Logger.Info("merge queue batch advanced", append(fields, "dropped", a.dropped)...)
	}
	if did && err == nil {
		p.wake()
	}
}

// journal keeps in the host's journal the landing a went on with after
// waiting on main, as done, saying no tier ran; each story a landed or sent
// back, as done; the wait on main its landing came to, as failed; and the
// batch it dropped, as failed when err stopped the advance.
func (p *queuePass) journal(a advanced, err error) {
	if r := a.resumed; r != nil {
		p.record(p.o.Now(), "done", fmt.Sprintf("merge queue batch %s lands again at %s after waiting %s on main, from the commits it verified%s, with no new verification: no tier ran",
			a.batch, r.Story, p.o.Now().Sub(r.At).Round(time.Second), headsText(a.from)))
	}
	for _, o := range a.outcomes {
		p.record(o.At, "done", o.Story+": "+o.Detail)
	}
	if s := a.stalled; s != nil {
		p.record(s.At, "failed", fmt.Sprintf("landing %s from merge queue batch %s waits on main: %s; the batch is kept, and lands on from the commits it verified once a look finds the main checkout clean", s.Story, a.batch, s.Reason))
	}
	if a.dropped == "" {
		return
	}
	outcome := "done"
	if err != nil {
		outcome = "failed"
	}
	p.record(p.o.Now(), outcome, fmt.Sprintf("dropped merge queue batch %s: %s", a.batch, a.dropped))
}

// record keeps one entry of the merge queue's in the host's journal, as flai
// serve's, at at.
func (p *queuePass) record(at time.Time, outcome, detail string) {
	if p.o.Host.Record == nil {
		return
	}
	p.o.Host.Record(hostapi.Entry{At: at.UTC().Format(time.RFC3339), Action: hostapi.ActionMergeQueue, Method: "serve.merge-queue", Project: p.e.Key, Root: p.e.Root,
		By: queueAuthor, Outcome: outcome, Detail: detail})
}

// headsText is " (" each head's story at its rebased commit ")", or "" for
// none.
func headsText(heads []mergequeue.Head) string {
	if len(heads) == 0 {
		return ""
	}
	parts := make([]string, 0, len(heads))
	for _, h := range heads {
		parts = append(parts, h.Story+" at "+short(h.Rebased))
	}
	return " (" + strings.Join(parts, ", ") + ")"
}

// followHold follows the preflight's hold from was, the queue's state when
// the look began, to the state the look left (ADR-0234). While the hold
// stands it keeps one thread to the operator open on it and journals it once,
// when it begins; once it is gone it resolves the thread, saying why, and
// journals that it cleared. It notes for changed whether the queue is held
// or its landing waits on main. A thread it cannot open or resolve is logged
// at warn and tried again at the next look; nothing here fails the look.
func (p *queuePass) followHold(repo *workitem.Repo, was mergequeue.State) {
	now, err := mergequeue.ReadState(repo.CacheDir())
	if err != nil {
		p.o.Logger.Warn("merge queue state not read after a look; the next look follows its preflight hold", "component", "serve", "root", p.e.Root, "key", p.e.Key, "err", err.Error())
		return
	}
	p.mu.Lock()
	p.holding = now.Preflight != nil || (now.Batch != nil && now.Batch.Stalled != nil)
	journalled, stale := p.heldSince, p.stale
	p.mu.Unlock()
	if hold := now.Preflight; hold != nil {
		thread, recorded := p.holdThread(repo, hold)
		if (was.Preflight == nil || recorded) && !hold.Since.Equal(journalled) {
			p.record(hold.Since, "done", holdDetail(hold, thread))
			p.mu.Lock()
			p.heldSince = hold.Since
			p.mu.Unlock()
		}
		return
	}
	if was.Preflight == nil && !stale {
		return
	}
	at := p.o.Now().UTC()
	resolved, err := resolveHoldThreads(repo, was.Preflight, at)
	p.mu.Lock()
	p.stale = err != nil
	p.mu.Unlock()
	if err != nil {
		p.o.Logger.Warn("merge queue preflight thread not resolved; the next look tries again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "err", err.Error())
	} else if len(resolved) > 0 {
		p.o.Logger.Info("merge queue preflight thread resolved", "component", "serve", "root", p.e.Root, "key", p.e.Key, "threads", resolved)
	}
	if prev := was.Preflight; prev != nil {
		detail := fmt.Sprintf("merge queue preflight hold cleared: main is clean after %s held since %s, so the queue builds batches again",
			at.Sub(prev.Since).Round(time.Second), prev.Since.UTC().Format(time.RFC3339))
		if len(resolved) > 0 {
			detail += "; thread " + strings.Join(resolved, ", ") + " resolved"
		}
		p.record(at, "done", detail)
	}
}

// holdThread is the open thread to the operator of hold, and whether it
// recorded one on it: the one the hold names while it is open; otherwise
// the one flai serve opened earlier and left open, or a new one, recorded on
// the hold. It is "" when none could be opened, which is logged at warn.
func (p *queuePass) holdThread(repo *workitem.Repo, hold *mergequeue.PreflightHold) (string, bool) {
	if hold.Thread != "" {
		if th, err := threads.Get(repo, hold.Thread); err == nil && th.Open() {
			return hold.Thread, false
		}
	}
	th, err := openHoldThread(repo, hold, p.o.Now())
	if err != nil {
		p.o.Logger.Warn("merge queue preflight thread not opened; the next look tries again", "component", "serve", "root", p.e.Root, "key", p.e.Key, "err", err.Error())
		return "", false
	}
	if _, err := mergequeue.SetPreflightThread(repo.CacheDir(), th.ID); err != nil {
		p.o.Logger.Warn("merge queue preflight thread not recorded on the hold; the next look records it", "component", "serve", "root", p.e.Root, "key", p.e.Key, "thread", th.ID, "err", err.Error())
	}
	p.o.Logger.Info("merge queue preflight thread recorded", "component", "serve", "root", p.e.Root, "key", p.e.Key, "thread", th.ID,
		"files", holdFiles(hold.Findings))
	return th.ID, true
}

// openHoldThread is the thread to the operator flai serve opened for a hold
// and left open, or, when there is none, a new one for hold, anchored on the
// manifest, as flai serve, at now.
func openHoldThread(repo *workitem.Repo, hold *mergequeue.PreflightHold, now time.Time) (*threads.Thread, error) {
	open, err := openHoldThreads(repo)
	if err != nil {
		return nil, err
	}
	if len(open) > 0 {
		return open[0], nil
	}
	return threads.New(repo, threads.NewOptions{Title: holdTitle, On: manifest.File, Author: queueAuthor, Text: holdText(hold), Now: now})
}

// openHoldThreads are the open threads flai serve opened on the manifest for
// a preflight hold, by their title.
func openHoldThreads(repo *workitem.Repo) ([]*threads.Thread, error) {
	all, err := threads.For(repo, manifest.File)
	if err != nil {
		return nil, fmt.Errorf("read the threads on %s: %w; fix the thread file it names", manifest.File, err)
	}
	title := workitem.CleanTitle(holdTitle)
	var out []*threads.Thread
	for _, th := range all {
		if th.Open() && th.Opener() == queueAuthor && th.Title == title {
			out = append(out, th)
		}
	}
	return out, nil
}

// resolveHoldThreads resolves, as flai serve at at, each open thread flai
// serve opened for a preflight hold, and prev's thread when it is open, as
// main is clean, and answers their IDs.
func resolveHoldThreads(repo *workitem.Repo, prev *mergequeue.PreflightHold, at time.Time) ([]string, error) {
	open, err := openHoldThreads(repo)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, th := range open {
		ids = append(ids, th.ID)
	}
	if prev != nil && prev.Thread != "" && !slices.Contains(ids, prev.Thread) {
		if th, err := threads.Get(repo, prev.Thread); err == nil && th.Open() {
			ids = append(ids, th.ID)
		}
	}
	why := fmt.Sprintf("main is clean: the merge queue's preflight found nothing at %s, so the queue builds batches again", at.Format(time.RFC3339))
	var resolved []string
	for _, id := range ids {
		if _, err := threads.Resolve(repo, id, queueAuthor, why, at); err != nil {
			return resolved, fmt.Errorf("resolve the merge queue's preflight thread %s: %w; resolve it with flai thread resolve %s", id, err, id)
		}
		resolved = append(resolved, id)
	}
	return resolved, nil
}

// holdText is the text of the thread to the operator for hold: each finding,
// with its file, line, rule, message, and what found it, that the queue
// builds no batch until they are fixed and changes none of them, and that the
// hold clears on its own. Code-like parts are in code spans, with no tab, so
// that the project's markdown lint passes the thread.
func holdText(hold *mergequeue.PreflightHold) string {
	var b strings.Builder
	b.WriteString("The merge queue's preflight found what a landing would fail on in the main checkout")
	if hold.Main != "" {
		fmt.Fprintf(&b, ", at main's head %s", codeSpan(short(hold.Main)))
	}
	b.WriteString(":\n\n")
	for _, f := range hold.Findings {
		b.WriteString("- " + codeSpan(spanText(f.Path)))
		if f.Line > 0 {
			fmt.Fprintf(&b, ", line %d", f.Line)
		}
		if f.Rule != "" {
			b.WriteString(", " + codeSpan(spanText(f.Rule)))
		}
		if msg := spanText(f.Message); msg != "" {
			b.WriteString(": " + codeSpan(msg))
		}
		fmt.Fprintf(&b, " (%s)\n", foundBy(f))
	}
	b.WriteString("\nThe queue builds no batch and runs no tier until these are fixed. It changes none of these files: they are yours. ")
	b.WriteString("Fix them in the main checkout. The queue looks again as the main checkout changes and as main's head moves, ")
	b.WriteString("and the look that finds main clean clears the hold, resolves this thread, and builds batches again, with nothing else asked of you.")
	return b.String()
}

// spanText is s on one line, with no tab, for a code span.
func spanText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// foundBy says what found f: flai check, or the lint of its tier.
func foundBy(f verify.FileFinding) string {
	if f.Tier == mergequeue.PreflightCheckTier {
		return "flai check"
	}
	return "the lint of tier " + codeSpan(spanText(f.Tier))
}

// holdDetail is the journal's entry for hold beginning: each file with its
// rules, and thread, or that none is open yet.
func holdDetail(hold *mergequeue.PreflightHold, thread string) string {
	detail := "merge queue held by its preflight, so it builds no batch until the main checkout is clean: " + strings.Join(holdFiles(hold.Findings), "; ")
	if thread == "" {
		return detail + "; no thread to the operator yet, the next look opens one"
	}
	return detail + "; thread " + thread
}

// holdFiles are the files findings name, in the order first named, each
// with the rules found in it.
func holdFiles(findings []verify.FileFinding) []string {
	var files []string
	rules := map[string][]string{}
	for _, f := range findings {
		if _, ok := rules[f.Path]; !ok {
			files = append(files, f.Path)
			rules[f.Path] = nil
		}
		if f.Rule != "" && !slices.Contains(rules[f.Path], f.Rule) {
			rules[f.Path] = append(rules[f.Path], f.Rule)
		}
	}
	out := make([]string, 0, len(files))
	for _, file := range files {
		if len(rules[file]) == 0 {
			out = append(out, file)
			continue
		}
		out = append(out, file+" ("+strings.Join(rules[file], ", ")+")")
	}
	return out
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

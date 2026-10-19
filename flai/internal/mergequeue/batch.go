package mergequeue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps a batch is at, as Batch.Step names them. StepRebasing is
// followed by the story being rebased, "rebasing S-0401", and StepVerifying,
// once a tier starts, by the tier, "verifying go-test".
const (
	StepBuilding  = "building"
	StepRebasing  = "rebasing"
	StepVerifying = "verifying"
	StepVerified  = "verified"
	StepFailed    = "failed"
)

// OutcomeLeftOut is the outcome of a story selected for a batch whose rebase
// onto it failed, or stopped on a path outside claims.shared that a story
// ahead in the batch changed: it is left out of the batch and stays
// enqueued (ADR-0175).
const OutcomeLeftOut = "left_out"

// StoppedAtVerify is the StoppedAt of a batch's result when its tiers could
// not be run at all, as when the batch's manifest declares tiers that are not
// valid.
const StoppedAtVerify = "verify"

// WorktreeName is the folder, under the main checkout's .flai-cache/worktrees,
// that the batch is built in; BranchPrefix begins the name of each batch's
// branch, merge-queue/<batch ID>. Both are flai's own.
const (
	WorktreeName = "merge-queue"
	BranchPrefix = "merge-queue/"
)

// summaryFindings is how many of the failing tier's findings a result's
// summary names; summaryMax caps its length.
const (
	summaryFindings = 3
	summaryMax      = 600
)

// Verifier verifies the batch's worktree at dir, built on base, with the
// tiers flai verify selects for what it changes against base, calling
// started with each tier's name as it starts. An error with ctx ended is a
// verification stopped; any other is one that could not run.
type Verifier func(ctx context.Context, dir, base string, started func(tier string)) (verify.Report, error)

// CheckoutVerifier is the Verifier flai serve runs: the test tiers of the
// manifest in the batch's worktree, read for the flai at version, run with
// verify.Tiers, git run through r. It leaves no tier out: the tiers marked
// queue, which flai verify leaves to the batch while the queue can batch,
// run here whenever the batch's paths select them (ADR-0174).
func CheckoutVerifier(r execx.Runner, version string) Verifier {
	return func(ctx context.Context, dir, base string, started func(string)) (verify.Report, error) {
		tiers, err := verify.CheckoutTiersFor(dir, r, version)
		if err != nil {
			return verify.Report{}, err
		}
		return verify.Tiers(ctx, verify.TierOptions{Dir: dir, Base: base, Tiers: tiers, Git: r, Started: started})
	}
}

// Deps are what RunBatch works with.
type Deps struct {
	// Runner runs git; nil is execx.System.
	Runner execx.Runner
	// Now stamps the batch, its steps' outcomes, and its result; nil is
	// time.Now.
	Now func() time.Time
	// Verify verifies the built batch; nil is CheckoutVerifier with this
	// flai's version.
	Verify Verifier
	// Log takes the batch's lifecycle; nil discards it.
	Log *slog.Logger
}

// RunBatch selects, builds, and verifies one batch of repo's merge queue,
// and reports whether it began one. It does nothing under pull-request
// acceptance, while the state's batch has a result and is not dropped (it
// waits to be landed or bisected), or when no story is to be batched: the
// enqueued stories, in the order they were enqueued, that Compute lists as
// eligible, but for one left out of a batch or left waiting at the main
// branch's commit and the head it has now, at most the manifest's
// merge_queue.batch_max. A batch with no result and not dropped, as a flai
// serve stopped mid-run leaves it, and a dropped one are replaced.
//
// It removes the last batch's worktree and branch, then builds the batch in
// its own worktree, WorktreeName, on its own branch, from the main branch's
// head, rebasing a copy of each story's verified head onto the batch in
// turn; the stories' own branches and worktrees are not touched. A rebase
// that stops runs the steps flai stream sync runs on a stop, in the batch's
// worktree: it merges the issue files and the ADR index, merges the other
// markdown files inside claims.shared by sections, writes the issue summary
// again, and continues, and a story whose every stop they settle stays in
// the batch; the fold after a clean rebase is left to the landing's
// acceptance (ADR-0175). A story whose rebase stops where they do not settle
// it is taken out of the batch, its conflicting paths on the batch's
// LeftOut, and the next goes on. Its outcome is by where the paths flai
// does not resolve conflict: all inside claims.shared, OutcomeWaiting, and
// it stays enqueued; any outside that a story ahead in the batch changed,
// OutcomeLeftOut, it stays enqueued, and flai messages it from each such
// story ahead about those paths; any outside that none changed, a conflict
// with the main branch itself, OutcomeSentBack: it is moved back to
// in-progress with the paths as the reason and dequeued, or, when the move
// fails, left out. With every story taken out the
// batch is removed and the state left with none. Otherwise the batch is
// verified once with d.Verify, and its result recorded; a batch dropped
// meanwhile keeps its result and stays dropped. Each step is written to the
// state as it begins. When ctx ends it stops with ctx's error and writes no
// result, so that the next run builds the batch again.
func RunBatch(ctx context.Context, d Deps, repo *workitem.Repo) (bool, error) {
	if d.Runner == nil {
		d.Runner = execx.System{}
	}
	if d.Now == nil {
		d.Now = time.Now
	}
	if d.Verify == nil {
		d.Verify = CheckoutVerifier(d.Runner, buildinfo.Version)
	}
	if d.Log == nil {
		d.Log = slog.New(slog.DiscardHandler)
	}
	q, err := Compute(d.Runner, repo)
	if err != nil {
		return false, err
	}
	if q.Off || (q.Batch != nil && q.Batch.Result != nil && q.Batch.Dropped == nil) || len(enqueuedEligible(q)) == 0 {
		return false, nil
	}
	main, err := storygit.MainBranch(d.Runner, repo.MainRoot)
	if err != nil {
		return false, fmt.Errorf("build a merge queue batch: find the main branch it is built from: %w; check out the main branch in the main checkout", err)
	}
	base, err := d.Runner.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+main+"^{commit}")
	if err != nil {
		return false, fmt.Errorf("build a merge queue batch: read the head of %s: %w; check that git can read the main branch", main, err)
	}
	picks := pick(q, strings.TrimSpace(base), repo.Manifest.BatchMax())
	if len(picks) == 0 {
		return false, nil
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	now := d.Now().UTC()
	b := &build{d: d, repo: repo, main: main, base: strings.TrimSpace(base), id: now.Format("20060102-150405")}
	b.branch = BranchPrefix + b.id
	b.path = filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)
	return true, b.run(ctx, q.Batch, picks, now)
}

// enqueuedEligible are q's stories enqueued, in the order they were, that
// are eligible.
func enqueuedEligible(q Queue) []Entry {
	eligible := map[string]Entry{}
	for _, e := range q.Eligible {
		eligible[workitem.CanonicalID(e.Story)] = e
	}
	var out []Entry
	for _, en := range q.Enqueued {
		if e, ok := eligible[workitem.CanonicalID(en.Story)]; ok {
			out = append(out, e)
		}
	}
	return out
}

// pick is the stories to batch from q: those enqueued and eligible, but for
// one whose last outcome left it out of a batch, or left it waiting, built
// on base with the head it has now, at most limit.
func pick(q Queue, base string, limit int) []Entry {
	last := map[string]Outcome{}
	for _, o := range q.Outcomes {
		last[workitem.CanonicalID(o.Story)] = o
	}
	var out []Entry
	for _, e := range enqueuedEligible(q) {
		if len(out) >= limit {
			break
		}
		if o, ok := last[workitem.CanonicalID(e.Story)]; ok && (o.Outcome == OutcomeLeftOut || o.Outcome == OutcomeWaiting) && o.Base == base && o.Head == e.Head {
			continue
		}
		out = append(out, e)
	}
	return out
}

// build is one batch being built and verified.
type build struct {
	d    Deps
	repo *workitem.Repo
	// id, branch, and path are the batch's ID, its branch, and its
	// worktree, absolute; main is the main branch, and base its commit the
	// batch is built on.
	id, branch, path, main, base string
}

// ahead is a story kept in the batch and the paths its rebase changed on
// it, as git names them.
type ahead struct {
	Story string
	Paths []string
}

// rebased is what rebasing a story onto the batch came to: the batch's new
// head and the paths the story changed on it when it is kept; its outcome
// and the paths its stop conflicted on when it is not.
type rebased struct {
	head    string
	changed []string
	out     *Outcome
	paths   []string
}

// errReplaced is a batch the state no longer holds, which no run but this
// one should have changed.
var errReplaced = errors.New("the merge queue's state no longer holds the batch being built; another flai serve may be running the queue: stop it, and let one run the queue")

// run builds the batch of picks in place of prev and verifies it.
func (b *build) run(ctx context.Context, prev *Batch, picks []Entry, now time.Time) error {
	if err := b.clean(prev); err != nil {
		return err
	}
	rel := filepath.ToSlash(filepath.Join(".flai-cache", "worktrees", WorktreeName))
	stories := make([]string, 0, len(picks))
	for _, e := range picks {
		stories = append(stories, e.Story)
	}
	if _, err := Update(b.repo.CacheDir(), func(s *State) error {
		s.Batch = &Batch{ID: b.id, Branch: b.branch, Worktree: rel, Base: b.base, Stories: slices.Clone(stories), StartedAt: now, Step: StepBuilding}
		return nil
	}); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(b.path), 0o755); err != nil {
		return fmt.Errorf("build merge queue batch %s: make %s: %w; check that the project's .flai-cache is writable", b.id, filepath.Dir(b.path), err)
	}
	if _, err := b.d.Runner.Run(b.repo.MainRoot, "git", "worktree", "add", "--quiet", "-b", b.branch, b.path, b.base); err != nil {
		return fmt.Errorf("build merge queue batch %s: check %s out in %s: %w; check that git can add a worktree there", b.id, b.branch, rel, err)
	}
	head := b.base
	var kept, left []string
	var before []ahead
	for _, e := range picks {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := b.step(StepRebasing + " " + e.Story); err != nil {
			return err
		}
		rb, err := b.rebase(e, head, before)
		if err != nil {
			return err
		}
		if rb.out == nil {
			head, kept = rb.head, append(kept, e.Story)
			before = append(before, ahead{Story: e.Story, Paths: rb.changed})
			if err := b.update(func(_ *State, batch *Batch) {
				batch.Heads = append(batch.Heads, Head{Story: e.Story, Verified: e.Head, Rebased: rb.head})
			}); err != nil {
				return err
			}
			continue
		}
		left = append(left, e.Story)
		if err := b.update(func(s *State, batch *Batch) {
			batch.Stories = slices.DeleteFunc(batch.Stories, func(x string) bool { return sameStory(x, e.Story) })
			batch.LeftOut = append(batch.LeftOut, LeftOut{Story: e.Story, Paths: rb.paths})
			s.RecordOutcome(*rb.out)
			if rb.out.Outcome == OutcomeSentBack {
				s.Dequeue(e.Story)
			}
		}); err != nil {
			return err
		}
	}
	if len(kept) == 0 {
		return b.close(left)
	}
	b.d.Log.Info("merge queue batch built", "component", "mergequeue", "batch", b.id, "branch", b.branch, "stories", kept, "left_out", left)
	if err := b.step(StepVerifying); err != nil {
		return err
	}
	rep, verr := b.d.Verify(ctx, b.path, b.base, func(tier string) {
		if ctx.Err() != nil {
			return
		}
		if err := b.step(StepVerifying + " " + tier); err != nil {
			b.d.Log.Warn("merge queue batch step not written; verification goes on", "component", "mergequeue", "batch", b.id, "tier", tier, "err", err.Error())
		}
	})
	if err := ctx.Err(); err != nil {
		return err
	}
	res := result(rep, verr, b.d.Now())
	if err := b.update(func(_ *State, batch *Batch) {
		batch.Result, batch.Step = res, StepVerified
		if !res.Passed {
			batch.Step = StepFailed
		}
	}); err != nil {
		return err
	}
	b.d.Log.Info("merge queue batch verified", "component", "mergequeue", "batch", b.id, "passed", res.Passed, "stopped_at", res.StoppedAt)
	return nil
}

// rebase rebases a copy of e's verified head, the commits it has beyond its
// merge base with the main branch, onto the batch's head onto in the
// batch's worktree, and moves the batch's branch to it: the new head, and
// the paths the story changed on the batch. A rebase that stops is
// continued over the stops the issue files, the ADR index, the issue
// summary, and the shared markdown files merged by sections settle, as flai
// stream sync continues it (ADR-0175), and one they settle is a rebase that
// finished. A stop they do not settle is aborted, and answers the story's
// outcome by where its paths conflict, as unsettled decides with before,
// the stories kept ahead of it, and the paths that stop conflicted on; a
// rebase git failed without stopping answers it left out. The batch's
// branch stays at onto. The story's own branch never moves. An error is a
// worktree git could not put back on the batch's branch, or whose paths
// changed it could not read.
func (b *build) rebase(e Entry, onto string, before []ahead) (rebased, error) {
	r, wt := b.d.Runner, b.path
	from, err := r.Run(wt, "git", "merge-base", b.base, e.Head)
	if err != nil {
		return rebased{out: b.leftOut(e, "failed: git found no merge base with the main branch ("+said(err)+")")}, nil
	}
	// update-refs off, so that no branch at a commit rebased, the story's
	// own among them, is moved with it
	_, err = r.Run(wt, "git", "-c", "rebase.updateRefs=false", "rebase", "--quiet", "--onto", onto, strings.TrimSpace(from), e.Head)
	files := issues.FilesIn(b.repo, wt, r, b.d.Now)
	var settled storygit.Settled
	if err != nil && storygit.RebaseInProgress(r, wt) {
		settled, err = storygit.ContinueOverResolved(r, wt, files, b.d.Log, e.Story)
	}
	if err == nil && !settled.Waits {
		if len(settled.Merged)+len(settled.Regenerated) > 0 {
			b.d.Log.Info("merge queue batch settled a story's stops by merging and writing files again", "component", "mergequeue",
				"batch", b.id, "story", e.Story, "merged", settled.Merged, "regenerated", settled.Regenerated)
		}
		if _, err := r.Run(wt, "git", "checkout", "--quiet", "-B", b.branch); err != nil {
			return rebased{}, fmt.Errorf("build merge queue batch %s: move %s to %s rebased: %w; remove the worktree %s, and the queue builds the batch again", b.id, b.branch, e.Story, err, b.path)
		}
		head, err := r.Run(wt, "git", "rev-parse", "HEAD")
		if err != nil {
			return rebased{}, fmt.Errorf("build merge queue batch %s: read its head after %s: %w", b.id, e.Story, err)
		}
		head = strings.TrimSpace(head)
		changed, err := r.Run(wt, "git", "-c", "core.quotePath=false", "diff", "--name-only", onto, head)
		if err != nil {
			return rebased{}, fmt.Errorf("build merge queue batch %s: read the paths %s changed on it: %w; check that git can diff in %s, and the queue builds the batch again", b.id, e.Story, err, b.path)
		}
		return rebased{head: head, changed: lines(changed)}, nil
	}
	if !settled.Waits {
		out := b.leftOut(e, "failed: git could not rebase it ("+said(err)+")")
		if _, err := r.Run(wt, "git", "checkout", "--quiet", b.branch); err != nil {
			return rebased{}, b.notBack(e, err)
		}
		return rebased{out: out}, nil
	}
	if _, aerr := r.Run(wt, "git", "rebase", "--abort"); aerr != nil {
		return rebased{}, fmt.Errorf("build merge queue batch %s: abort the rebase of %s that stopped: %w; run git rebase --abort in %s, and the queue builds the batch again", b.id, e.Story, aerr, b.path)
	}
	if _, err := r.Run(wt, "git", "checkout", "--quiet", b.branch); err != nil {
		return rebased{}, b.notBack(e, err)
	}
	return rebased{out: b.unsettled(e, settled.Conflicts, files, before), paths: settled.Conflicts}, nil
}

// notBack is the error of a batch worktree git could not check the batch's
// branch out in again after e was taken out of the batch.
func (b *build) notBack(e Entry, err error) error {
	return fmt.Errorf("build merge queue batch %s: check %s out again after leaving %s out: %w; remove the worktree %s, and the queue builds the batch again", b.id, b.branch, e.Story, err, b.path)
}

// lines are the non-empty lines git printed.
func lines(out string) []string {
	var ps []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			ps = append(ps, l)
		}
	}
	return ps
}

// stop is how the batch treats a stop it could not settle, by where the
// paths flai does not resolve conflict (ADR-0175).
type stop struct {
	// Outcome is OutcomeWaiting, OutcomeLeftOut, or OutcomeSentBack.
	Outcome string
	// Outside is the paths outside claims.shared.
	Outside []string
	// Ahead is each story ahead in the batch that changed paths of Outside,
	// with those paths, in batch order.
	Ahead []ahead
}

// classify is how the batch treats a stop whose paths flai does not resolve
// are unresolved, shared saying whether a path is inside claims.shared and
// before being the stories kept ahead of it with the paths each changed:
// waiting when every path is inside claims.shared; left out when a story
// ahead changed a path outside it; sent back, a conflict with the main
// branch itself, when none did.
func classify(unresolved []string, shared func(string) bool, before []ahead) stop {
	var s stop
	for _, p := range unresolved {
		if !shared(p) {
			s.Outside = append(s.Outside, p)
		}
	}
	if len(s.Outside) == 0 {
		s.Outcome = OutcomeWaiting
		return s
	}
	for _, a := range before {
		var both []string
		for _, p := range s.Outside {
			if slices.Contains(a.Paths, p) {
				both = append(both, p)
			}
		}
		if len(both) > 0 {
			s.Ahead = append(s.Ahead, ahead{Story: a.Story, Paths: both})
		}
	}
	s.Outcome = OutcomeSentBack
	if len(s.Ahead) > 0 {
		s.Outcome = OutcomeLeftOut
	}
	return s
}

// unsettled is e's outcome for a stop on paths that files did not settle,
// before being the stories kept ahead of it (ADR-0175): waiting on paths
// inside claims.shared alone; left out on paths outside it that a story
// ahead changed, each such story ahead messaging e about them; and sent
// back to in-progress on paths outside it that none changed. A stop on no
// path leaves e out.
func (b *build) unsettled(e Entry, paths []string, files storygit.Files, before []ahead) *Outcome {
	if len(paths) == 0 {
		return b.leftOut(e, "stopped")
	}
	claims := b.repo.SharedClaims()
	s := classify(files.Unresolved(paths), func(p string) bool { _, ok := claims.Match(p); return ok }, before)
	switch s.Outcome {
	case OutcomeWaiting:
		b.d.Log.Info("merge queue batch left a story waiting on shared paths", "component", "mergequeue", "batch", b.id, "story", e.Story, "paths", paths)
		return b.waiting(e, paths)
	case OutcomeLeftOut:
		stories := make([]string, 0, len(s.Ahead))
		for _, a := range s.Ahead {
			stories = append(stories, a.Story)
		}
		b.d.Log.Info("merge queue batch left out a story that conflicts with a story ahead", "component", "mergequeue", "batch", b.id, "story", e.Story, "paths", s.Outside, "ahead", stories)
		b.tell(e, s.Ahead)
		return b.leftOut(e, fmt.Sprintf("stopped on conflicts in %s with %s, ahead of it in the batch", workitem.Shorten(s.Outside, 10), strings.Join(stories, ", ")))
	default:
		return b.sendBack(e, s.Outside)
	}
}

// tell messages e from each story ahead in the batch about the paths both
// change, in the two stories' conversation, as a trial-merge conflict is
// told (ADR-0121). Messages that cannot be sent are logged and stop nothing.
func (b *build) tell(e Entry, conflicting []ahead) {
	var errs []error
	for _, a := range conflicting {
		if _, err := messages.Notify(b.repo, messages.SendOptions{From: a.Story, To: e.Story, Author: storygit.ConflictAuthor, Text: b.conflictText(a.Story, e.Story, a.Paths), About: a.Paths, Now: b.d.Now()}); err != nil {
			errs = append(errs, fmt.Errorf("message %s from %s: %w", e.Story, a.Story, err))
		}
	}
	if len(errs) > 0 {
		b.d.Log.Warn("merge queue batch conflict messages not sent", "component", "mergequeue", "batch", b.id, "story", e.Story, "err", errors.Join(errs...).Error())
	}
}

// conflictText is the message from the story ahead to story, whose rebase
// onto the batch stopped on paths both change.
func (b *build) conflictText(ahead, story string, paths []string) string {
	var list strings.Builder
	for _, p := range paths {
		fmt.Fprintf(&list, "- `%s`\n", p)
	}
	return fmt.Sprintf("%s does not rebase onto %s in merge queue batch %s.\n\n"+
		"The merge queue's batch %s could not rebase %s onto the stories ahead of it, %s among them: both change\n\n%s\n"+
		"%s is left out of the batch and stays enqueued; the queue batches it again once the main branch or %s moves. "+
		"Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. "+
		"Once %s lands, sync %s with `flai stream sync %s` and settle these paths. "+
		"When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.",
		story, ahead, b.id, b.id, storygit.Branch(story), ahead, list.String(), story, storygit.Branch(story), ahead, story, story)
}

// sendBack moves e back to in-progress, by flai, for a conflict with the
// main branch itself in paths, and answers its outcome, sent back; a move
// that fails is logged and leaves e out instead.
func (b *build) sendBack(e Entry, paths []string) *Outcome {
	reason := fmt.Sprintf("the merge queue's batch %s could not rebase it onto %s: conflicts in %s; sync it with flai stream sync %s",
		b.id, b.main, workitem.Shorten(paths, 10), e.Story)
	it, err := b.repo.Get(e.Story)
	if err == nil {
		_, err = b.repo.Transition(it, workitem.InProgress, sendBackBy, reason, b.d.Now())
	}
	if err != nil {
		b.d.Log.Warn("merge queue batch could not send back a story that conflicts with the main branch; it is left out", "component", "mergequeue", "batch", b.id, "story", e.Story, "paths", paths, "err", err.Error())
		return b.leftOut(e, "stopped on conflicts with "+b.main+" in "+workitem.Shorten(paths, 10)+", and could not be sent back to in-progress ("+said(err)+")")
	}
	b.d.Log.Info("merge queue batch sent back a story that conflicts with the main branch", "component", "mergequeue", "batch", b.id, "story", e.Story, "paths", paths)
	return &Outcome{Story: e.Story, Outcome: OutcomeSentBack, At: b.d.Now(), Base: b.base, Head: e.Head,
		Detail: "sent back to in-progress and out of the queue: " + reason}
}

// leftOut is e's outcome, left out of the batch for why.
func (b *build) leftOut(e Entry, why string) *Outcome {
	return &Outcome{
		Story: e.Story, Outcome: OutcomeLeftOut, At: b.d.Now(), Base: b.base, Head: e.Head,
		Detail: fmt.Sprintf("left out of merge queue batch %s: rebasing %s onto the batch %s; it stays enqueued, and is batched again once the main branch or %s moves",
			b.id, storygit.Branch(e.Story), why, storygit.Branch(e.Story)),
	}
}

// waiting is e's outcome, waiting on the shared paths its rebase onto the
// batch stopped on and the merge steps could not settle.
func (b *build) waiting(e Entry, paths []string) *Outcome {
	return &Outcome{
		Story: e.Story, Outcome: OutcomeWaiting, At: b.d.Now(), Base: b.base, Head: e.Head,
		Detail: fmt.Sprintf("waiting, out of merge queue batch %s: rebasing %s onto the batch stopped on shared paths the merge steps could not settle, %s; it stays enqueued, is not sent back for them, and is batched again once the main branch or %s moves",
			b.id, storygit.Branch(e.Story), workitem.Shorten(paths, 10), storygit.Branch(e.Story)),
	}
}

// close removes the batch whose every story was left out, its worktree and
// branch, and leaves the state with no batch: there is nothing to verify.
func (b *build) close(left []string) error {
	if err := b.clean(&Batch{Branch: b.branch}); err != nil {
		return err
	}
	if _, err := Update(b.repo.CacheDir(), func(s *State) error {
		if s.Batch == nil || s.Batch.ID != b.id {
			return errReplaced
		}
		s.Batch = nil
		return nil
	}); err != nil {
		return err
	}
	b.d.Log.Info("merge queue batch left every story out", "component", "mergequeue", "batch", b.id, "left_out", left)
	return nil
}

// clean removes the batch worktree, WorktreeName, whatever it holds, and the
// branch of prev and the branch checked out there, each when it is a batch's.
// It touches nothing else.
func (b *build) clean(prev *Batch) error {
	r, root := b.d.Runner, b.repo.MainRoot
	var branches []string
	if prev != nil {
		branches = append(branches, prev.Branch)
	}
	if _, err := os.Stat(b.path); err == nil {
		if storygit.IsWorktree(r, b.path) {
			if cur, err := r.Run(b.path, "git", "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
				branches = append(branches, strings.TrimSpace(cur))
			}
			if _, err := r.Run(root, "git", "worktree", "remove", "--force", b.path); err != nil {
				return fmt.Errorf("remove the merge queue's last batch worktree %s: %w; remove it with git worktree remove --force %s", b.path, err, b.path)
			}
		} else if err := os.RemoveAll(b.path); err != nil {
			return fmt.Errorf("remove the merge queue's last batch folder %s: %w; remove it by hand", b.path, err)
		}
	}
	if _, err := r.Run(root, "git", "worktree", "prune"); err != nil {
		return fmt.Errorf("prune the merge queue's last batch worktree: %w; run git worktree prune in %s", err, root)
	}
	for _, br := range slices.Compact(slices.Sorted(slices.Values(branches))) {
		if !strings.HasPrefix(br, BranchPrefix) || !storygit.BranchExists(r, root, br) {
			continue
		}
		if _, err := r.Run(root, "git", "branch", "-D", br); err != nil {
			return fmt.Errorf("delete the merge queue's last batch branch %s: %w; delete it with git branch -D %s", br, err, br)
		}
	}
	return nil
}

// step writes the batch's step.
func (b *build) step(step string) error {
	return b.update(func(_ *State, batch *Batch) { batch.Step = step })
}

// update lets fn change the state and the batch being built in it, and
// writes them, keeping whatever else the state holds, the batch's being
// dropped among it.
func (b *build) update(fn func(*State, *Batch)) error {
	_, err := Update(b.repo.CacheDir(), func(s *State) error {
		if s.Batch == nil || s.Batch.ID != b.id {
			return errReplaced
		}
		fn(s, s.Batch)
		return nil
	})
	return err
}

// result is the batch's result from its verification's report, or from the
// error that kept it from running, finished at at.
func result(rep verify.Report, err error, at time.Time) *Result {
	res := &Result{Passed: rep.Passed && err == nil, StoppedAt: rep.StoppedAt, FinishedAt: at.UTC()}
	switch {
	case err != nil:
		res.StoppedAt, res.Summary = StoppedAtVerify, clip("the batch's tiers could not run: "+err.Error())
	case rep.Passed:
		res.Summary = passedSummary(rep)
	default:
		res.Summary = failedSummary(rep)
	}
	return res
}

// passedSummary names the tiers a passing report ran.
func passedSummary(rep verify.Report) string {
	var ran []string
	for _, s := range rep.Steps {
		if s.State == verify.Passed {
			ran = append(ran, s.Name)
		}
	}
	if len(ran) == 0 {
		return "no tier selects the paths the batch changes"
	}
	return "passed " + strings.Join(ran, ", ")
}

// failedSummary is the failing tier and its first findings, a line each.
func failedSummary(rep verify.Report) string {
	for _, s := range rep.Steps {
		if s.Name != rep.StoppedAt || s.State != verify.Failed {
			continue
		}
		if len(s.Findings) == 0 {
			if s.ExitCode != nil {
				return clip(fmt.Sprintf("%s failed with exit status %d", s.Name, *s.ExitCode))
			}
			return clip(s.Name + " failed")
		}
		var lines []string
		for _, f := range s.Findings[:min(summaryFindings, len(s.Findings))] {
			lines = append(lines, strings.TrimSpace(strings.SplitN(f.Text(), "\n", 2)[0]))
		}
		more := len(s.Findings) - len(lines) + s.Omitted
		text := s.Name + " failed: " + strings.Join(lines, "; ")
		if more > 0 {
			text += fmt.Sprintf("; and %d more", more)
		}
		return clip(text)
	}
	return clip(rep.StoppedAt + " failed")
}

// clip is s cut to summaryMax characters, marked when it is cut.
func clip(s string) string {
	if r := []rune(s); len(r) > summaryMax {
		return string(r[:summaryMax-1]) + "…"
	}
	return s
}

// said is git's error on one line, short.
func said(err error) string {
	return clip(strings.Join(strings.Fields(err.Error()), " "))
}

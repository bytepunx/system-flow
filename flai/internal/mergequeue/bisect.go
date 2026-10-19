package mergequeue

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The steps of a batch that failed, as Batch.Step names them (ADR-0176).
// StepBisecting is followed by the stories a step verifies of the batch's,
// "bisecting 2 of 4", and, once a tier starts, by the tier, "bisecting 2 of
// 4, verifying unit"; StepBisected is a bisection that found its story.
const (
	StepBisecting = "bisecting"
	StepBisected  = "bisected"
)

// sendBackBy moves the failing story back to in-progress; sendBackAuthor
// opens its thread, as the forge's send-back does.
const (
	sendBackBy     = "flai"
	sendBackAuthor = "flai serve"
)

// withDefaults is d with what it leaves nil filled as RunBatch fills it.
func withDefaults(d Deps) Deps {
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
	return d
}

// RunBisect bisects the batch of repo's merge queue that failed, and reports
// whether it ran a step or found the story with which the batch fails. It
// does nothing under pull-request acceptance, with no batch, or with one
// dropped, not verified, passed, or bisected already.
//
// Main itself counts as passing and the whole batch as failing. Each step
// verifies main with the first stories of the batch, half those between the
// most known to pass and the fewest known to fail, checked out in the batch's
// worktree at the last one's rebased head, and halves again, until the first
// story with which the batch fails is found: its commit in the batch, the
// result there, and a landing of the stories before it. Each step is written
// to the state as it begins and once it has a result. When ctx ends it stops
// with ctx's error, leaving the step without a result, and the next run
// verifies that step again. A batch dropped meanwhile is left as it is.
func RunBisect(ctx context.Context, d Deps, repo *workitem.Repo) (bool, error) {
	d = withDefaults(d)
	if repo.Manifest.Acceptance.PullRequest() {
		return false, nil
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return false, err
	}
	batch := st.Batch
	if batch == nil || batch.Dropped != nil || batch.Result == nil || batch.Result.Passed || batch.Bisect.Found() {
		return false, nil
	}
	if err := headsMatch(batch); err != nil {
		return false, err
	}
	b := &build{d: d, repo: repo, id: batch.ID, branch: batch.Branch, base: batch.Base, path: filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)}
	n := len(batch.Stories)
	ran := false
	for {
		if err := ctx.Err(); err != nil {
			return ran, err
		}
		mid, commit, dropped, found := 0, "", false, false
		if err := b.update(func(_ *State, batch *Batch) {
			if batch.Dropped != nil {
				dropped = true
				return
			}
			if batch.Bisect == nil {
				batch.Bisect = &Bisect{Steps: []BisectStep{}, Good: 0, Bad: n}
			}
			bs := batch.Bisect
			if bs.Bad-bs.Good <= 1 {
				found = true
				return
			}
			mid = (bs.Good + bs.Bad) / 2
			commit = batch.Heads[mid-1].Rebased
			step := BisectStep{Stories: slices.Clone(batch.Stories[:mid]), Commit: commit, StartedAt: d.Now().UTC()}
			if last := len(bs.Steps) - 1; last >= 0 && bs.Steps[last].Result == nil {
				bs.Steps[last] = step // a step a run stopped is verified again in place
			} else {
				bs.Steps = append(bs.Steps, step)
			}
			batch.Step = fmt.Sprintf("%s %d of %d", StepBisecting, mid, n)
		}); err != nil {
			return ran, err
		}
		if dropped {
			return ran, nil
		}
		if found {
			return true, b.found()
		}
		ran = true
		if err := b.checkout(commit, "--detach", commit); err != nil {
			return ran, err
		}
		rep, verr := d.Verify(ctx, b.path, b.base, func(tier string) {
			if ctx.Err() != nil {
				return
			}
			if err := b.step(fmt.Sprintf("%s %d of %d, %s %s", StepBisecting, mid, n, StepVerifying, tier)); err != nil {
				d.Log.Warn("merge queue batch step not written; bisection goes on", "component", "mergequeue", "batch", b.id, "tier", tier, "err", err.Error())
			}
		})
		if err := ctx.Err(); err != nil {
			return ran, err
		}
		res := result(rep, verr, d.Now())
		if err := b.update(func(_ *State, batch *Batch) {
			bs := batch.Bisect
			if bs == nil || len(bs.Steps) == 0 {
				return
			}
			bs.Steps[len(bs.Steps)-1].Result = res
			if res.Passed {
				bs.Good = mid
			} else {
				bs.Bad = mid
			}
		}); err != nil {
			return ran, err
		}
		d.Log.Info("merge queue batch bisection step verified", "component", "mergequeue", "batch", b.id, "stories", mid, "of", n, "commit", commit, "passed", res.Passed)
	}
}

// headsMatch is an error unless batch's heads are its stories', one to one
// and in order, each with a rebased head, as a bisection checks them out.
func headsMatch(batch *Batch) error {
	ok := len(batch.Stories) > 0 && len(batch.Heads) == len(batch.Stories)
	for i := 0; ok && i < len(batch.Stories); i++ {
		ok = sameStory(batch.Heads[i].Story, batch.Stories[i]) && batch.Heads[i].Rebased != ""
	}
	if ok {
		return nil
	}
	heads := make([]string, 0, len(batch.Heads))
	for _, h := range batch.Heads {
		heads = append(heads, h.Story)
	}
	return fmt.Errorf("bisect merge queue batch %s: its heads, for %s, do not match its stories, %s, one to one and in order, so it cannot be bisected; the batch must be built again from the main branch",
		batch.ID, orNone(heads), orNone(batch.Stories))
}

// orNone is ids joined with commas, or "none".
func orNone(ids []string) string {
	if len(ids) == 0 {
		return "none"
	}
	return strings.Join(ids, ", ")
}

// found records the story the bisection found, the first of the fewest known
// to fail, its commit, the result there, and the landing of the stories
// before it, and checks the batch's branch out again in its worktree.
func (b *build) found() error {
	var failing string
	var steps int
	if err := b.update(func(_ *State, batch *Batch) {
		bs := batch.Bisect
		bad := bs.Bad
		if bad < 1 || bad > len(batch.Stories) || len(batch.Heads) != len(batch.Stories) {
			return // headsMatch and the halving keep Bad within the batch
		}
		bs.Failing, bs.Commit = batch.Stories[bad-1], batch.Heads[bad-1].Rebased
		bs.Result = copyResult(batch.Result)
		for _, s := range bs.Steps {
			if len(s.Stories) == bad && s.Result != nil {
				bs.Result = copyResult(s.Result)
			}
		}
		batch.Landing = &Landing{Stories: slices.Clone(batch.Stories[:bad-1]), Landed: []string{}, Head: batch.Base}
		batch.Step = StepBisected
		failing, steps = bs.Failing, len(bs.Steps)
	}); err != nil {
		return err
	}
	if err := b.checkout(b.branch, b.branch); err != nil {
		return err
	}
	b.d.Log.Info("merge queue batch bisected", "component", "mergequeue", "batch", b.id, "failing", failing, "steps", steps)
	return nil
}

// copyResult is a copy of r, nil for nil.
func copyResult(r *Result) *Result {
	if r == nil {
		return nil
	}
	c := *r
	return &c
}

// checkout checks what out in the batch's worktree, with args to git
// checkout; whatever a tier left changed there is discarded, as the worktree
// is flai's own.
func (b *build) checkout(what string, args ...string) error {
	argv := append([]string{"checkout", "--quiet", "--force"}, args...)
	if _, err := b.d.Runner.Run(b.path, "git", argv...); err != nil {
		return fmt.Errorf("bisect merge queue batch %s: check %s out in its worktree %s: %w; check that the worktree is there and git can check it out, or the batch must be built again from the main branch", b.id, what, b.path, err)
	}
	return nil
}

// SendBack sends back to in-progress, once, the story the bisection of repo's
// merge queue batch found, and reports whether it moved it. A story in review
// is moved by flai, the reason, naming the batch, the commit, and what failed
// there, appended to its notes, and a thread opened on it to its agent by flai
// serve, saying what failed, at which commit, with which stories, and what to
// do. A story no longer in review is not moved and gets no thread. Either way
// it leaves the queue with OutcomeFailing as its outcome, and the bisection
// is marked sent back. It does nothing under pull-request acceptance, with no
// batch, or with one dropped, not bisected, or sent back already.
func SendBack(d Deps, repo *workitem.Repo) (bool, error) {
	d = withDefaults(d)
	if repo.Manifest.Acceptance.PullRequest() {
		return false, nil
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return false, err
	}
	batch := st.Batch
	if batch == nil || batch.Dropped != nil || !batch.Bisect.Found() || batch.Bisect.SentBack {
		return false, nil
	}
	bs := batch.Bisect
	story, now := bs.Failing, d.Now().UTC()
	head, _ := batch.HeadOf(story)
	tier, summary := "verify", "the batch failed"
	if bs.Result != nil {
		tier, summary = orDefault(bs.Result.StoppedAt, tier), orDefault(bs.Result.Summary, summary)
	}
	it, err := repo.Get(story)
	if err != nil && !errors.Is(err, workitem.ErrNotFound) {
		return false, fmt.Errorf("send %s back to in-progress, as merge queue batch %s fails with it: read it: %w; fix the story's file, and flai serve sends it back at its next look", story, batch.ID, err)
	}
	moved := err == nil && it.Type == workitem.Story && it.Status == workitem.Review
	detail := fmt.Sprintf("sent back to in-progress: merge queue batch %s fails with it at %s", batch.ID, tier)
	if moved {
		th, err := openFailing(repo, batch, tier, summary, head, now)
		if err != nil {
			return false, fmt.Errorf("send %s back to in-progress, as merge queue batch %s fails with it: open its thread: %w; flai serve sends it back at its next look", story, batch.ID, err)
		}
		reason := fmt.Sprintf("merge queue batch %s fails with %s at %s: %s", batch.ID, story, short(bs.Commit), oneLine(summary))
		if _, err := repo.Transition(it, workitem.InProgress, sendBackBy, reason, now); err != nil {
			return false, fmt.Errorf("send %s back to in-progress, as merge queue batch %s fails with it: %w; move it with flai move %s in-progress --reason, and see thread %s", story, batch.ID, err, story, th.ID)
		}
		d.Log.Info("merge queue story sent back", "component", "mergequeue", "batch", batch.ID, "story", story, "commit", bs.Commit, "tier", tier, "thread", th.ID)
	} else {
		was := "not found"
		if it != nil {
			was = it.Status
		}
		detail = fmt.Sprintf("merge queue batch %s fails with it at %s; it was %s, so it was not sent back", batch.ID, tier, was)
		d.Log.Info("merge queue failing story not in review; not sent back", "component", "mergequeue", "batch", batch.ID, "story", story, "status", was)
	}
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		if s.Batch == nil || s.Batch.ID != batch.ID || s.Batch.Bisect == nil {
			return errReplaced
		}
		s.Dequeue(story)
		s.RecordOutcome(Outcome{Story: story, Outcome: OutcomeFailing, At: now, Detail: detail, Base: batch.Base, Head: head.Verified})
		s.Batch.Bisect.SentBack = true
		return nil
	}); err != nil {
		return moved, err
	}
	return moved, nil
}

// openFailing opens the thread on the failing story of batch to its agent,
// or answers the open one flai serve opened on it with the same title, as a
// send-back before this one opened and could not record.
func openFailing(repo *workitem.Repo, batch *Batch, tier, summary string, head Head, now time.Time) (*threads.Thread, error) {
	story := batch.Bisect.Failing
	opt := threads.NewOptions{
		Title: fmt.Sprintf("Merge queue batch %s fails with %s at %s", batch.ID, story, tier),
		On:    story, Author: sendBackAuthor, Text: failingText(batch, tier, summary, head), Now: now,
	}
	if all, err := threads.For(repo, story); err == nil {
		title := workitem.CleanTitle(opt.Title)
		for _, th := range all {
			if th.Open() && th.Opener() == sendBackAuthor && th.Title == title {
				return th, nil
			}
		}
	}
	return threads.New(repo, opt)
}

// failingText says what failed of batch's failing story, at which commit,
// with which stories, and what to do.
func failingText(batch *Batch, tier, summary string, head Head) string {
	bs := batch.Bisect
	story := bs.Failing
	with := "alone"
	before := batch.Stories[:bs.Bad-1]
	if len(before) > 0 {
		with = "with the stories before it, " + strings.Join(before, ", ") + ","
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Merge queue batch %s fails with %s at %s. The queue verified %s %s on main at %s, at %s, the copy of its verified head %s rebased in the batch, and %s failed:\n\n",
		codeSpan(batch.ID), story, codeSpan(tier), story, with, codeSpan(short(batch.Base)), codeSpan(short(bs.Commit)), codeSpan(short(head.Verified)), codeSpan(tier))
	fmt.Fprintf(&b, "%s\n\n", fenced(summary))
	if len(before) > 0 {
		b.WriteString("The stories before it land; those after it stay enqueued. ")
	}
	fmt.Fprintf(&b, "flai moved %s back to in-progress and took it out of the merge queue. Fix it on %s, run flai verify %s, move it to review, and enqueue it again with flai queue add %s.",
		story, codeSpan(storygit.Branch(story)), story, story)
	return b.String()
}

// fenced is s in a text code block whose fence is longer than any run of
// backticks in s.
func fenced(s string) string {
	fence := "```"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	return fence + "text\n" + strings.TrimSpace(s) + "\n" + fence
}

// codeSpan is s in a code span whose fence is longer than any run of
// backticks in s.
func codeSpan(s string) string {
	fence := "`"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		s = " " + s + " "
	}
	return fence + s + fence
}

// oneLine is s on one line.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// orDefault is s, or def when s is blank.
func orDefault(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

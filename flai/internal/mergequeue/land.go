package mergequeue

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// StepLanding is the step of a batch landing its stories, followed by the
// story it lands, "landing S-0401" (ADR-0176).
const StepLanding = "landing"

// queueBy drops a batch as flai serve takes it on.
const queueBy = "flai serve"

// Accepter accepts story as by through flai accept's flow, as the landing of
// merge queue batch batch. flai serve gives one that runs flai accept <story>
// --by <by> --batch <batch> with this flai's own binary in the main checkout.
type Accepter func(ctx context.Context, story, by, batch string) error

// Landable is whether repo's merge queue lets story land now as part of
// batch, as the preview of flai accept --batch checks it: the story's heads
// in the batch, whose Rebased head is the verified commit the acceptance
// takes, and every reason it may not land, each a sentence saying why and
// what clears it, none when it may. Git is read through r. It is an error
// only when the queue's state or the main branch cannot be read.
func Landable(r execx.Runner, repo *workitem.Repo, story, batch string) (Head, []string, error) {
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return Head{}, nil, err
	}
	b := st.Batch
	if b == nil {
		return Head{}, []string{fmt.Sprintf("the merge queue holds no batch, so %s cannot land as part of merge queue batch %s; flai serve lands the stories of the batch it verified, or accept %s by hand with flai accept %s", story, batch, story, story)}, nil
	}
	if b.ID != batch {
		return Head{}, []string{fmt.Sprintf("the merge queue's batch is %s, not %s, so %s cannot land as part of %s; flai serve lands the stories of the batch the queue holds", b.ID, batch, story, batch)}, nil
	}
	head, hasHead := b.HeadOf(story)
	var reasons []string
	if b.Dropped != nil {
		reasons = append(reasons, fmt.Sprintf("merge queue batch %s was dropped by %s: %s; its stories still enqueued are batched again from the main branch and land with the next batch", b.ID, b.Dropped.By, b.Dropped.Reason))
	}
	l := b.Landing
	switch {
	case l == nil && b.Result == nil:
		reasons = append(reasons, fmt.Sprintf("merge queue batch %s is not verified yet (it is at %s), so nothing of it lands; wait for flai serve to verify it", b.ID, orDefault(b.Step, StepBuilding)))
	case l == nil && !b.Result.Passed:
		reasons = append(reasons, fmt.Sprintf("merge queue batch %s failed and is not bisected yet, so it is not known which of its stories land; wait for flai serve to bisect it", b.ID))
	case l == nil:
		reasons = append(reasons, fmt.Sprintf("merge queue batch %s passed, but its landing has not begun; wait for flai serve to land its stories in queue order", b.ID))
	case !slices.ContainsFunc(l.Stories, func(x string) bool { return sameStory(x, story) }):
		reasons = append(reasons, fmt.Sprintf("%s is not among the stories merge queue batch %s lands (%s); a story the batch does not land stays enqueued, or was sent back, and lands with a batch that passes with it", story, b.ID, orNone(l.Stories)))
	case slices.ContainsFunc(l.Landed, func(x string) bool { return sameStory(x, story) }):
		reasons = append(reasons, fmt.Sprintf("%s has landed from merge queue batch %s already; there is nothing more to land", story, b.ID))
	default:
		if next, ok := l.Next(); ok && !sameStory(next, story) {
			reasons = append(reasons, fmt.Sprintf("%s is not the next story merge queue batch %s lands: %s is, as the batch lands its stories in queue order; land %s first", story, b.ID, next, next))
		}
	}
	if hasHead {
		if reason := branchMoved(r, repo, b, head); reason != "" {
			reasons = append(reasons, reason)
		}
	} else if l != nil && slices.ContainsFunc(l.Stories, func(x string) bool { return sameStory(x, story) }) {
		reasons = append(reasons, fmt.Sprintf("merge queue batch %s kept no heads for %s, so the commit it verified it at is not known; the batch must be built again from the main branch", b.ID, story))
	}
	if l != nil && l.Head != "" {
		main, at, err := mainHead(r, repo, fmt.Sprintf("check that %s may land from merge queue batch %s", story, b.ID))
		if err != nil {
			return head, reasons, err
		}
		if at != l.Head {
			reasons = append(reasons, fmt.Sprintf("the main branch %s is at %s, not at %s, where merge queue batch %s's landing left it; main moved other than by the batch's landing, so flai serve drops the batch and builds it again from the new main", main, short(at), short(l.Head), b.ID))
		}
	}
	return head, reasons, nil
}

// branchMoved is why the story of head may not land when its branch is no
// longer at the head batch b rebased a copy of, or cannot be read; "" when it
// is there.
func branchMoved(r execx.Runner, repo *workitem.Repo, b *Batch, head Head) string {
	branch := storygit.Branch(head.Story)
	out, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", branch+"^{commit}")
	at := strings.TrimSpace(out)
	if err != nil || at == "" {
		why := "git names no commit"
		if err != nil {
			why = said(err)
		}
		return fmt.Sprintf("cannot read the head of %s to check that it is the head merge queue batch %s verified a copy of: %s; check that git can read the branch", branch, b.ID, why)
	}
	if at != head.Verified {
		return fmt.Sprintf("%s is at %s, not at %s, the head merge queue batch %s rebased a copy of and verified; the branch moved since, so %s lands with a batch built from its new head once its flai verify passes there", branch, short(at), short(head.Verified), b.ID, head.Story)
	}
	return ""
}

// mainHead is the main branch's name and its head in repo's main checkout,
// read through r; doing says what it is read for, in an error.
func mainHead(r execx.Runner, repo *workitem.Repo, doing string) (string, string, error) {
	main, err := storygit.MainBranch(r, repo.MainRoot)
	if err != nil {
		return "", "", fmt.Errorf("%s: find the main branch: %w; check out the main branch in the main checkout", doing, err)
	}
	out, err := r.Run(repo.MainRoot, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+main+"^{commit}")
	if err != nil || strings.TrimSpace(out) == "" {
		if err == nil {
			err = errors.New("git names no commit")
		}
		return "", "", fmt.Errorf("%s: read the head of %s: %w; check that git can read the main branch", doing, main, err)
	}
	return main, strings.TrimSpace(out), nil
}

// Landed records in the merge queue's state in cacheDir that by accepted
// story at at as the landing of merge queue batch batch: what flai accept
// --batch calls, once it has merged the story, in place of AcceptedByHand.
// The story leaves the queue with OutcomeLanded as its outcome, and, when
// the state's batch is batch and lands it, it is landed in the batch's
// landing, once. It never drops the batch, and changes no other batch.
func Landed(cacheDir, story, batch, by string, at time.Time) (State, error) {
	return Update(cacheDir, func(s *State) error {
		s.Dequeue(story)
		s.RecordOutcome(Outcome{Story: story, Outcome: OutcomeLanded, At: at, Detail: fmt.Sprintf("landed by merge queue batch %s as %s", batch, by)})
		if b := s.Batch; b != nil && b.ID == batch {
			landInLanding(b.Landing, story)
		}
		return nil
	})
}

// landInLanding adds story, as l names it, to l's landed stories, when l
// lands it and it is not landed already.
func landInLanding(l *Landing, story string) {
	if l == nil || slices.ContainsFunc(l.Landed, func(x string) bool { return sameStory(x, story) }) {
		return
	}
	if i := slices.IndexFunc(l.Stories, func(x string) bool { return sameStory(x, story) }); i >= 0 {
		l.Landed = append(l.Landed, l.Stories[i])
	}
}

// DropIfMainMoved drops, as by at at, repo's merge queue batch that is
// verified and not dropped when the main branch's head is not where the
// batch expects it, the head its landing left or, before any landing, the
// commit it was built on, and reports whether it dropped it: the main branch
// moved other than by the batch's landing, as by an acceptance by hand, so
// what the batch verified is not what would land. Git is read through r. A
// batch replaced or dropped meanwhile is left as it is.
func DropIfMainMoved(r execx.Runner, repo *workitem.Repo, by string, at time.Time) (bool, error) {
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return false, err
	}
	b := st.Batch
	if b == nil || b.Result == nil || b.Dropped != nil {
		return false, nil
	}
	_, head, err := mainHead(r, repo, "check that the main branch is where merge queue batch "+b.ID+" expects it")
	if err != nil {
		return false, err
	}
	if head == expectedMain(b) {
		return false, nil
	}
	dropped := false
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		cur := s.Batch
		if cur == nil || cur.ID != b.ID || cur.Dropped != nil || cur.Result == nil {
			return nil
		}
		want := expectedMain(cur)
		if head == want {
			return nil
		}
		reason := fmt.Sprintf("the main branch moved from %s to %s while merge queue batch %s waited to land; it is built again from the new main", short(want), short(head), cur.ID)
		cur.Dropped, dropped = &Dropped{Reason: reason, By: by, At: at.UTC()}, true
		return nil
	}); err != nil {
		return false, err
	}
	return dropped, nil
}

// expectedMain is where b expects the main branch's head: where its landing
// left it, or the commit it was built on before its landing began.
func expectedMain(b *Batch) string {
	if b.Landing != nil && b.Landing.Head != "" {
		return b.Landing.Head
	}
	return b.Base
}

// dropBatch drops repo's merge queue batch id, for reason, as by at at,
// unless it was replaced or dropped already.
func dropBatch(repo *workitem.Repo, id, reason, by string, at time.Time) error {
	_, err := Update(repo.CacheDir(), func(s *State) error {
		if s.Batch == nil || s.Batch.ID != id || s.Batch.Dropped != nil {
			return nil
		}
		s.Batch.Dropped = &Dropped{Reason: reason, By: by, At: at.UTC()}
		return nil
	})
	return err
}

// buildOf is the build of batch, as the landing, the bisection, and the
// clean-up work on it.
func buildOf(d Deps, repo *workitem.Repo, batch *Batch) *build {
	return &build{d: d, repo: repo, id: batch.ID, branch: batch.Branch, base: batch.Base, path: filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)}
}

// Land lands the stories of repo's merge queue batch's landing, one at a
// time in queue order, each accepted with accept as whoever enqueued it, and
// reports whether it landed any. Before each it drops the batch when the
// main branch moved (DropIfMainMoved) and stops; it stops too when the batch
// is dropped, has no landing, or every story landed, or when the next story
// is no longer enqueued, as one removed from the queue, whose batch is
// dropped. Each landing is written as the batch's step as it begins, and once
// it is done the main branch's head is where the landing leaves it. A
// landing that fails drops the batch, so that the stories not landed are
// batched again from the main branch, and is returned. When ctx ends it
// stops with ctx's error, the batch kept, and the next run lands on. It does
// nothing under pull-request acceptance.
func Land(ctx context.Context, d Deps, accept Accepter, repo *workitem.Repo) (bool, error) {
	d = withDefaults(d)
	if repo.Manifest.Acceptance.PullRequest() {
		return false, nil
	}
	landed := false
	for {
		if err := ctx.Err(); err != nil {
			return landed, err
		}
		if dropped, err := DropIfMainMoved(d.Runner, repo, queueBy, d.Now()); err != nil || dropped {
			return landed, err
		}
		st, err := ReadState(repo.CacheDir())
		if err != nil {
			return landed, err
		}
		batch := st.Batch
		if batch == nil || batch.Dropped != nil {
			return landed, nil
		}
		story, ok := batch.Landing.Next()
		if !ok {
			return landed, nil
		}
		en, ok := st.EnqueuedAs(story)
		if !ok {
			// Remove drops the batch as it takes the story out; this keeps a
			// state edited by hand from waiting on a story no one enqueued
			reason := fmt.Sprintf("%s is no longer enqueued, so merge queue batch %s does not land it; the stories not landed are batched again from the main branch", story, batch.ID)
			return landed, dropBatch(repo, batch.ID, reason, queueBy, d.Now())
		}
		b := buildOf(d, repo, batch)
		if err := b.step(StepLanding + " " + story); err != nil {
			return landed, err
		}
		if err := accept(ctx, story, en.By, batch.ID); err != nil {
			if ctx.Err() != nil {
				return landed, ctx.Err()
			}
			reason := fmt.Sprintf("landing %s failed: %s; the stories not landed are batched again from the main branch", story, said(err))
			if derr := dropBatch(repo, batch.ID, reason, queueBy, d.Now()); derr != nil {
				return landed, fmt.Errorf("land %s from merge queue batch %s as %s: %w; and drop the batch: %w", story, batch.ID, en.By, err, derr)
			}
			d.Log.Warn("merge queue batch dropped; a landing failed", "component", "mergequeue", "batch", batch.ID, "story", story, "by", en.By)
			return landed, fmt.Errorf("land %s from merge queue batch %s as %s: %w; the batch is dropped and its stories not landed are batched again from the main branch: fix what stopped the acceptance, or accept %s by hand with flai accept %s", story, batch.ID, en.By, err, story, story)
		}
		landed = true
		_, head, err := mainHead(d.Runner, repo, fmt.Sprintf("record where landing %s from merge queue batch %s left the main branch", story, batch.ID))
		if err != nil {
			return landed, err
		}
		if err := b.update(func(_ *State, batch *Batch) {
			if batch.Landing != nil {
				batch.Landing.Head = head
				landInLanding(batch.Landing, story) // an older flai accept --batch does not
			}
		}); err != nil {
			return landed, err
		}
		d.Log.Debug("merge queue story landed", "component", "mergequeue", "batch", batch.ID, "story", story, "by", en.By, "head", head)
	}
}

// Advance takes repo's merge queue batch, once RunBatch has verified it, one
// step on, and reports whether it did anything. A batch whose main branch
// moved is dropped (DropIfMainMoved), and RunBatch builds it again. A batch
// that passed lands every story; one that failed is bisected (RunBisect), the
// story it fails with sent back (SendBack), and the stories before it
// landed; with Land, accepting each with accept. Once every story the batch
// lands has landed, its worktree and branch are removed and the state left
// with no batch, so that RunBatch batches the stories still enqueued, those
// after a failing one among them. A bisection that cannot run drops the
// batch and is returned. It does nothing under pull-request acceptance, with
// no batch, one dropped, which RunBatch replaces, or one not verified yet.
func Advance(ctx context.Context, d Deps, accept Accepter, repo *workitem.Repo) (bool, error) {
	d = withDefaults(d)
	if repo.Manifest.Acceptance.PullRequest() {
		return false, nil
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		return false, err
	}
	batch := st.Batch
	if batch == nil || batch.Dropped != nil || batch.Result == nil {
		return false, nil
	}
	if dropped, err := DropIfMainMoved(d.Runner, repo, queueBy, d.Now()); err != nil || dropped {
		return dropped, err
	}
	b := buildOf(d, repo, batch)
	did := false
	if batch.Result.Passed && batch.Landing == nil {
		if err := b.update(func(_ *State, batch *Batch) {
			if batch.Landing == nil {
				batch.Landing = &Landing{Stories: slices.Clone(batch.Stories), Landed: []string{}, Head: batch.Base}
			}
		}); err != nil {
			return did, err
		}
		did = true
	}
	sentBack := ""
	if !batch.Result.Passed {
		ran, err := RunBisect(ctx, d, repo)
		did = did || ran
		if err != nil {
			if ctx.Err() == nil {
				if derr := dropBatch(repo, batch.ID, said(err), queueBy, d.Now()); derr != nil {
					return did, fmt.Errorf("%w; and drop the batch: %w", err, derr)
				}
			}
			return did, err
		}
		if st, err = ReadState(repo.CacheDir()); err != nil {
			return did, err
		}
		cur := st.Batch
		if cur == nil || cur.ID != batch.ID || cur.Dropped != nil || !cur.Bisect.Found() {
			return did, nil
		}
		if !cur.Bisect.SentBack {
			if _, err := SendBack(d, repo); err != nil {
				return did, err
			}
			did = true
		}
		sentBack = cur.Bisect.Failing
	}
	landed, err := Land(ctx, d, accept, repo)
	did = did || landed
	if err != nil {
		return did, err
	}
	if st, err = ReadState(repo.CacheDir()); err != nil {
		return did, err
	}
	cur := st.Batch
	if cur == nil || cur.ID != batch.ID || cur.Dropped != nil || cur.Landing == nil {
		return did, nil
	}
	if _, more := cur.Landing.Next(); more {
		return did, nil
	}
	if err := b.clean(cur); err != nil {
		return did, err
	}
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		if s.Batch != nil && s.Batch.ID == batch.ID {
			s.Batch = nil
		}
		return nil
	}); err != nil {
		return did, err
	}
	d.Log.Info("merge queue batch landed", "component", "mergequeue", "batch", batch.ID, "landed", cur.Landing.Landed, "sent_back", sentBack)
	return true, nil
}

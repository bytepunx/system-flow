package mergequeue

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// t3 is when the verified batches are landed.
var t3 = t2.Add(time.Hour)

// acceptCall is one acceptance the queue asked for.
type acceptCall struct{ story, by, batch string }

// fakeAccept stands in for flai accept --batch in real git: it checks with
// Landable that the story may land, makes an empty commit on main for it,
// and records it with Landed, as the acceptance does once it has merged. It
// keeps each call and the heads Landable answered; fail, when it names the
// story, is the acceptance's error instead, and after, when set, runs once
// each story has landed.
type fakeAccept struct {
	t     *testing.T
	repo  *workitem.Repo
	fail  map[string]error
	after func(story string)
	calls []acceptCall
	heads []Head
}

func (a *fakeAccept) accept(_ context.Context, story, by, batch string) error {
	a.calls = append(a.calls, acceptCall{story, by, batch})
	if err := a.fail[story]; err != nil {
		return err
	}
	head, reasons, err := Landable(execx.System{}, a.repo, story, batch)
	if err != nil || len(reasons) > 0 {
		return fmt.Errorf("%s may not land: %v, %w", story, reasons, err)
	}
	a.heads = append(a.heads, head)
	runGit(a.t, a.repo.MainRoot, "commit", "-q", "--allow-empty", "-m", "accept "+story)
	if _, err := Landed(a.repo.CacheDir(), story, batch, by, t3); err != nil {
		return err
	}
	if a.after != nil {
		a.after(story)
	}
	return nil
}

// advance runs Advance, failing the test on an error, and answers whether it
// did anything and the state it left.
func advance(t *testing.T, d Deps, a *fakeAccept, repo *workitem.Repo) (bool, State) {
	t.Helper()
	did, err := Advance(context.Background(), d, a.accept, repo)
	if err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	return did, st
}

// enqueueAs enqueues story as by.
func enqueueAs(t *testing.T, repo *workitem.Repo, story, by string) {
	t.Helper()
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		s.Enqueue(story, by, t0)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

// gone fails the test unless batch's worktree and branch are removed.
func gone(t *testing.T, repo *workitem.Repo, batch *Batch) {
	t.Helper()
	if storygit.BranchExists(execx.System{}, repo.MainRoot, batch.Branch) {
		t.Errorf("the batch's branch %s is left", batch.Branch)
	}
	if _, err := os.Stat(filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the batch's worktree is left: %v", err)
	}
}

// outcomeOf is story's last outcome in st, and false when it has none.
func outcomeOf(st State, story string) (Outcome, bool) {
	for _, o := range st.Outcomes {
		if sameStory(o.Story, story) {
			return o, true
		}
	}
	return Outcome{}, false
}

// S-0400 criterion 1: a green batch lands each story in queue order, each
// accepted as whoever enqueued it as the landing of the batch, at its rebased
// head in the batch; each leaves the queue as landed, and once every one has
// the batch, its worktree, and its branch are gone.
func TestAGreenBatchLandsEachStoryInQueueOrderAsItsEnqueuer(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueueAs(t, repo, b, "orchestrator")
	enqueue(t, repo, a)
	_, built := runBatch(t, deps(proc, t1), repo)
	bt := built.Batch
	if bt == nil || bt.Step != StepVerified || !slices.Equal(bt.Stories, []string{b, a}) {
		t.Fatalf("batch %+v, want %s and %s verified", bt, b, a)
	}

	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, deps(proc, t2), fa, repo)
	if !did {
		t.Error("Advance did nothing")
	}
	want := []acceptCall{{b, "orchestrator", bt.ID}, {a, "alex", bt.ID}}
	if !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want %+v", fa.calls, want)
	}
	if !slices.Equal(fa.heads, bt.Heads) {
		t.Errorf("Landable answered heads %+v, want the batch's %+v", fa.heads, bt.Heads)
	}
	for _, c := range want {
		o, ok := outcomeOf(st, c.story)
		if !ok || o.Outcome != OutcomeLanded || o.Detail != "landed by merge queue batch "+bt.ID+" as "+c.by {
			t.Errorf("%s's outcome %+v, want landed as %s", c.story, o, c.by)
		}
	}
	if len(st.Enqueued) != 0 || st.Batch != nil {
		t.Errorf("enqueued %+v, batch %+v; want the queue empty and no batch", st.Enqueued, st.Batch)
	}
	gone(t, repo, bt)
	if got, wantSubjects := subjects(t, repo, "main"), []string{"accept " + a, "accept " + b, "init"}; !slices.Equal(got, wantSubjects) {
		t.Errorf("main has %v, want %v", got, wantSubjects)
	}
}

// S-0400 criterion 3: a red batch is bisected to the story it fails with,
// which is sent back; the stories before it land, and those after it stay
// enqueued, with no outcome, and are batched next once the batch is gone.
func TestARedBatchLandsTheStoriesBeforeTheFailingOne(t *testing.T) {
	repo, v, stories, _, before := failedBatch(t, 3, 2)
	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, fileDeps(v, t2), fa, repo)
	if !did {
		t.Error("Advance did nothing")
	}
	if want := []acceptCall{{stories[0], "alex", before.Batch.ID}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want %+v alone", fa.calls, want)
	}
	if o, ok := outcomeOf(st, stories[0]); !ok || o.Outcome != OutcomeLanded {
		t.Errorf("%s's outcome %+v, want landed", stories[0], o)
	}
	it, err := repo.Get(stories[1])
	if err != nil {
		t.Fatal(err)
	}
	if o, ok := outcomeOf(st, stories[1]); it.Status != workitem.InProgress || !ok || o.Outcome != OutcomeFailing {
		t.Errorf("%s is %s with outcome %+v, want it sent back as failing", stories[1], it.Status, o)
	}
	if _, ok := st.EnqueuedAs(stories[2]); !ok || len(st.Enqueued) != 1 {
		t.Errorf("enqueued %+v, want %s alone", st.Enqueued, stories[2])
	}
	if o, ok := outcomeOf(st, stories[2]); ok {
		t.Errorf("%s has outcome %+v, want none", stories[2], o)
	}
	if st.Batch != nil {
		t.Fatalf("batch %+v, want it cleared", st.Batch)
	}
	gone(t, repo, before.Batch)

	ran, next := runBatch(t, fileDeps(v, t3), repo)
	if !ran || next.Batch == nil || !slices.Equal(next.Batch.Stories, stories[2:]) {
		t.Errorf("ran %v, batch %+v; want one of %s", ran, next.Batch, stories[2])
	}
}

// S-0552 criterion 3: a red batch of four failing with its third story is
// bisected as before, its first half passing in the first step and the first
// three failing in the second, and the stories of that clean half land
// through the acceptance of the batch's landing with no verification of their
// own: the verifier runs for the batch and the two steps, and never again.
// The log holds one verified line, two step_verified lines, and two landings.
func TestACleanHalfLandsWithoutAVerificationOfItsOwn(t *testing.T) {
	repo, v, stories, _, before := failedBatch(t, 4, 3)
	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, fileDeps(v, t2), fa, repo)
	if !did {
		t.Error("Advance did nothing")
	}
	if want := [][]string{stories[:2], stories[:3]}; !reflect.DeepEqual(v.stories, want) {
		t.Errorf("verified after the batch given %v, want the two steps' %v alone", v.stories, want)
	}
	id := before.Batch.ID
	if want := []acceptCall{{stories[0], "alex", id}, {stories[1], "alex", id}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want the clean half %+v as the batch's landing", fa.calls, want)
	}
	for _, s := range stories[:2] {
		if o, ok := outcomeOf(st, s); !ok || o.Outcome != OutcomeLanded {
			t.Errorf("%s's outcome %+v, want landed", s, o)
		}
	}
	if o, ok := outcomeOf(st, stories[2]); !ok || o.Outcome != OutcomeFailing {
		t.Errorf("%s's outcome %+v, want it sent back as failing", stories[2], o)
	}
	if _, ok := st.EnqueuedAs(stories[3]); !ok || len(st.Enqueued) != 1 {
		t.Errorf("enqueued %+v, want %s alone", st.Enqueued, stories[3])
	}
	count := map[string]int{}
	for _, e := range eventLog(t, repo.CacheDir()) {
		count[e.Event]++
	}
	if count[EventVerified] != 1 || count[EventStepVerified] != 2 || count[EventLanded] != 2 {
		t.Errorf("logged %v, want 1 verified, 2 step_verified, and 2 landed", count)
	}
	if st.Batch != nil {
		t.Errorf("batch %+v, want it cleared once its landing is done", st.Batch)
	}
}

// S-0400 criterion 3: a batch that fails with its first story lands nothing;
// the story is sent back and the batch cleared.
func TestARedBatchFailingWithItsFirstStoryLandsNothing(t *testing.T) {
	repo, v, stories, _, before := failedBatch(t, 2, 1)
	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, fileDeps(v, t2), fa, repo)
	if !did || len(fa.calls) != 0 {
		t.Errorf("did %v, accepted %+v; want nothing landed", did, fa.calls)
	}
	if it, err := repo.Get(stories[0]); err != nil || it.Status != workitem.InProgress {
		t.Errorf("%s is %v (%v), want it sent back to in-progress", stories[0], it, err)
	}
	if st.Batch != nil {
		t.Errorf("batch %+v, want it cleared", st.Batch)
	}
	if _, ok := st.EnqueuedAs(stories[1]); !ok {
		t.Errorf("%s left the queue", stories[1])
	}
	gone(t, repo, before.Batch)
}

// S-0400 criterion 5: when main moves while a verified batch waits to land,
// as by an acceptance by hand, the batch is dropped, naming the move, and
// nothing lands; the next run builds it again from the new main.
func TestMainMovingDropsAVerifiedBatchWhichIsBuiltAgain(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	_, built := runBatch(t, deps(proc, t1), repo)
	if built.Batch == nil || built.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, want it verified", built.Batch)
	}
	commitOnMain(t, repo, "y.txt", "accepted by hand\n")
	moved := runGit(t, repo.MainRoot, "rev-parse", "main")

	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, deps(proc, t2), fa, repo)
	if !did || len(fa.calls) != 0 {
		t.Errorf("did %v, accepted %+v; want the batch dropped and nothing landed", did, fa.calls)
	}
	dr := st.Batch.Dropped
	if dr == nil || dr.By != "flai serve" || !dr.At.Equal(t2) {
		t.Fatalf("dropped %+v, want by flai serve at %s", dr, t2)
	}
	for _, says := range []string{"moved from " + short(built.Batch.Base) + " to " + short(moved), built.Batch.ID, "change paths outside wip/: y.txt", "built again from the new main"} {
		if !strings.Contains(dr.Reason, says) {
			t.Errorf("reason %q does not say %q", dr.Reason, says)
		}
	}

	ran, next := runBatch(t, deps(proc, t3), repo)
	if !ran || next.Batch == nil || next.Batch.ID == built.Batch.ID || next.Batch.Base != moved || next.Batch.Dropped != nil {
		t.Errorf("ran %v, batch %+v; want a new batch built on %s", ran, next.Batch, moved)
	}
}

// S-0400 criterion 5: main moving between two landings, as by an acceptance
// by hand while flai serve was stopped, drops the batch before the next
// landing; the main branch's head the first landing left is not a move.
func TestMainMovingBetweenLandingsDropsTheBatch(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	runBatch(t, deps(proc, t1), repo)
	ctx, cancel := context.WithCancel(context.Background())
	fa := &fakeAccept{t: t, repo: repo, after: func(string) { cancel() }}

	_, err := Advance(ctx, deps(proc, t2), fa.accept, repo)
	if !errors.Is(err, context.Canceled) || len(fa.calls) != 1 {
		t.Fatalf("%v, accepted %+v; want flai serve stopped after %s landed", err, fa.calls, a)
	}
	if dropped, err := DropIfMainMoved(execx.System{}, repo, "flai serve", t2); dropped || err != nil {
		t.Errorf("dropped %v, %v; want main at the head the landing left kept", dropped, err)
	}
	commitOnMain(t, repo, "y.txt", "accepted by hand\n")

	fa.after = nil
	_, st := advance(t, deps(proc, t3), fa, repo)
	if len(fa.calls) != 1 {
		t.Errorf("accepted %+v, want %s not landed", fa.calls, b)
	}
	if st.Batch == nil || st.Batch.Dropped == nil || !strings.Contains(st.Batch.Dropped.Reason, "the main branch moved") {
		t.Errorf("batch %+v, want it dropped as main moved", st.Batch)
	}
	if _, ok := st.EnqueuedAs(b); !ok {
		t.Errorf("%s left the queue", b)
	}
}

// commitOnMain writes content to rel in repo's main checkout and commits it
// alone on main, a path under the ignored wip folder too.
func commitOnMain(t *testing.T, repo *workitem.Repo, rel, content string) {
	t.Helper()
	writeFile(t, repo.MainRoot, rel, content)
	runGit(t, repo.MainRoot, "add", "-f", "--", rel)
	runGit(t, repo.MainRoot, "commit", "-q", "-m", "change "+rel, "--", rel)
}

// ADR-0282, I-0287: main moving while a verified batch waits to land, by a
// commit that changes only wip paths no story of the batch changes, such as
// flai's replan, keeps the batch, which lands its story over it with no new
// verification.
func TestAWipOnlyCommitOnMainKeepsAVerifiedBatchWhichLands(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	_, built := runBatch(t, deps(proc, t1), repo)
	if built.Batch == nil || built.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, want it verified", built.Batch)
	}
	commitOnMain(t, repo, "wip/kanban/board.md", "# Board\n\nreplanned\n")
	if dropped, err := DropIfMainMoved(execx.System{}, repo, "flai serve", t2); dropped || err != nil {
		t.Fatalf("dropped %v, %v; want the batch kept over a wip-only commit", dropped, err)
	}
	verified := len(proc.ran)

	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, deps(proc, t2), fa, repo)
	if !did || !slices.Equal(fa.calls, []acceptCall{{a, "alex", built.Batch.ID}}) {
		t.Errorf("did %v, accepted %+v; want %s landed from %s", did, fa.calls, a, built.Batch.ID)
	}
	if o, ok := outcomeOf(st, a); !ok || o.Outcome != OutcomeLanded {
		t.Errorf("%s's outcome %+v, want landed", a, o)
	}
	if st.Batch != nil || len(proc.ran) != verified {
		t.Errorf("batch %+v, verifications %d then %d; want the batch landed and gone with no new verification", st.Batch, verified, len(proc.ran))
	}
}

// ADR-0282: a wip-only commit on main between two landings, such as the
// replan an acceptance sets off, keeps the batch, and the next story lands
// over it.
func TestAWipOnlyCommitBetweenLandingsKeepsTheBatch(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	_, built := runBatch(t, deps(proc, t1), repo)
	ctx, cancel := context.WithCancel(context.Background())
	fa := &fakeAccept{t: t, repo: repo, after: func(string) { cancel() }}
	if _, err := Advance(ctx, deps(proc, t2), fa.accept, repo); !errors.Is(err, context.Canceled) || len(fa.calls) != 1 {
		t.Fatalf("%v, accepted %+v; want flai serve stopped after %s landed", err, fa.calls, a)
	}
	commitOnMain(t, repo, "wip/agents/"+a+".md", "# "+a+"\n")

	fa.after = nil
	_, st := advance(t, deps(proc, t3), fa, repo)
	if want := []acceptCall{{a, "alex", built.Batch.ID}, {b, "alex", built.Batch.ID}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want %+v", fa.calls, want)
	}
	if st.Batch != nil || len(st.Enqueued) != 0 {
		t.Errorf("batch %+v, enqueued %+v; want both landed and the batch gone", st.Batch, st.Enqueued)
	}
}

// ADR-0282 §2: main moving while a verified batch waits to land drops it, as
// before, when main does not descend from the batch's base, or its commits
// change a wip path a story of the batch changes too; the reason says which.
func TestMainMovingOtherThanByWipOnlyCommitsDropsTheBatch(t *testing.T) {
	for _, tc := range []struct {
		name string
		move func(t *testing.T, repo *workitem.Repo)
		says string
	}{
		{name: "not descending", says: "does not descend from", move: func(t *testing.T, repo *workitem.Repo) {
			orphan := runGit(t, repo.MainRoot, "commit-tree", "main^{tree}", "-m", "rewritten")
			runGit(t, repo.MainRoot, "reset", "-q", "--hard", orphan)
		}},
		{name: "a wip path the story changes", says: "change wip/ paths a story of the batch changes too: wip/agents/shared.md", move: func(t *testing.T, repo *workitem.Repo) {
			commitOnMain(t, repo, "wip/agents/shared.md", "main's\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, proc := gitProject(t, "")
			a, _ := readyStoryChanging(t, repo, "A", t0.Add(time.Hour), func(_, wt string) {
				writeFile(t, wt, "a.txt", "a\n")
				writeFile(t, wt, "wip/agents/shared.md", "the story's\n")
				runGit(t, wt, "add", "-f", "--", "wip/agents/shared.md")
			})
			enqueue(t, repo, a)
			_, built := runBatch(t, deps(proc, t1), repo)
			if built.Batch == nil || built.Batch.Step != StepVerified {
				t.Fatalf("batch %+v, want it verified", built.Batch)
			}
			tc.move(t, repo)

			fa := &fakeAccept{t: t, repo: repo}
			_, st := advance(t, deps(proc, t2), fa, repo)
			if len(fa.calls) != 0 {
				t.Errorf("accepted %+v, want nothing landed", fa.calls)
			}
			if st.Batch == nil || st.Batch.Dropped == nil || !strings.Contains(st.Batch.Dropped.Reason, tc.says) || !strings.Contains(st.Batch.Dropped.Reason, "the main branch moved") {
				t.Errorf("batch %+v, want it dropped as main moved, saying %q", st.Batch, tc.says)
			}
		})
	}
}

// S-0400 criterion 1: a landing that fails drops the batch, saying which
// story failed to land and that the rest are batched again; the error is
// returned and the story stays enqueued.
func TestALandingThatFailsDropsTheBatch(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	runBatch(t, deps(proc, t1), repo)
	fa := &fakeAccept{t: t, repo: repo, fail: map[string]error{a: errors.New("merge stopped\non x.txt")}}

	_, err := Advance(context.Background(), deps(proc, t2), fa.accept, repo)
	if err == nil || !strings.Contains(err.Error(), "merge stopped") || !strings.Contains(err.Error(), "flai accept "+a) {
		t.Errorf("error %v, want the acceptance's, saying what to do", err)
	}
	st, rerr := ReadState(repo.CacheDir())
	if rerr != nil {
		t.Fatal(rerr)
	}
	dr := st.Batch.Dropped
	if dr == nil || dr.Reason != "landing "+a+" failed: merge stopped on x.txt; the stories not landed are batched again from the main branch" || dr.By != "flai serve" {
		t.Errorf("dropped %+v, want it dropped as the landing of %s failed", dr, a)
	}
	if _, ok := st.EnqueuedAs(a); !ok || st.Batch.Step != "landing "+a {
		t.Errorf("enqueued %+v, step %q; want %s still enqueued and the step its landing", st.Enqueued, st.Batch.Step, a)
	}
}

// S-0400 criterion 2: flai accept --batch may take the batch's verification
// only for the next story of a batch whose landing lets it land, its branch at
// the head the batch verified a copy of and main where the landing left it;
// otherwise every reason says why and what clears it.
func TestLandableAnswersTheRebasedHeadOrEveryReason(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	_, built := runBatch(t, deps(proc, t1), repo)
	bt := built.Batch
	set := func(fn func(*Batch)) {
		t.Helper()
		if _, err := Update(repo.CacheDir(), func(s *State) error {
			fn(s.Batch)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	set(func(b *Batch) {
		b.Landing = &Landing{Stories: slices.Clone(b.Stories), Landed: []string{}, Head: b.Base}
	})
	landable := func(story, batch string) (Head, []string) {
		t.Helper()
		head, reasons, err := Landable(execx.System{}, repo, story, batch)
		if err != nil {
			t.Fatal(err)
		}
		return head, reasons
	}
	says := func(name string, reasons []string, want ...string) {
		t.Helper()
		if len(reasons) != 1 {
			t.Errorf("%s: reasons %q, want one", name, reasons)
			return
		}
		for _, w := range want {
			if !strings.Contains(reasons[0], w) {
				t.Errorf("%s: reason %q does not say %q", name, reasons[0], w)
			}
		}
	}

	head, reasons := landable(a, bt.ID)
	if len(reasons) != 0 || head != bt.Heads[0] {
		t.Errorf("head %+v, reasons %q; want %s's heads %+v and none", head, reasons, a, bt.Heads[0])
	}
	_, reasons = landable(a, "20991231-000000")
	says("another batch", reasons, bt.ID, "not 20991231-000000")
	_, reasons = landable(b, bt.ID)
	says("not the next", reasons, "not the next", a+" is")

	set(func(b *Batch) {
		b.Dropped = &Dropped{Reason: "S-0009 was accepted by hand by alex", By: "alex", At: t2}
	})
	_, reasons = landable(a, bt.ID)
	says("dropped", reasons, "dropped", "S-0009 was accepted by hand by alex")
	set(func(b *Batch) { b.Dropped = nil })

	set(func(b *Batch) { b.Landing, b.Result = nil, nil })
	_, reasons = landable(a, bt.ID)
	says("not verified", reasons, "not verified")
	set(func(b *Batch) { b.Result = &Result{StoppedAt: "unit", FinishedAt: t1} })
	_, reasons = landable(a, bt.ID)
	says("not bisected", reasons, "not bisected")
	set(func(b *Batch) {
		b.Result = &Result{Passed: true, FinishedAt: t1}
		b.Landing = &Landing{Stories: slices.Clone(b.Stories), Landed: []string{}, Head: b.Base}
	})

	wt := repo.WorktreePath(a)
	runGit(t, wt, "commit", "-q", "--allow-empty", "-m", "more")
	_, reasons = landable(a, bt.ID)
	says("the branch moved", reasons, storygit.Branch(a), short(headA))
	runGit(t, wt, "reset", "-q", "--hard", headA)

	commitOnMain(t, repo, "wip/kanban/board.md", "# Board\n")
	_, reasons = landable(a, bt.ID)
	if len(reasons) != 0 {
		t.Errorf("reasons %q after a wip-only commit on main, want none: the batch lands over it (ADR-0282)", reasons)
	}
	commitOnMain(t, repo, "y.txt", "accepted by hand\n")
	_, reasons = landable(a, bt.ID)
	says("main moved", reasons, "main branch main", short(bt.Base), "change paths outside wip/: y.txt")

	if _, err := Update(repo.CacheDir(), func(s *State) error {
		s.Batch = nil
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	_, reasons = landable(a, bt.ID)
	says("no batch", reasons, "no batch", "flai accept "+a)
}

// S-0400 criterion 1: flai accept --batch's record of a landing takes the
// story out of the queue with landed as its outcome and lands it in its
// batch's landing, once, leaving the batch in place; for another batch it
// changes no batch.
func TestLandedRecordsTheLandingAndLeavesTheBatch(t *testing.T) {
	dir := t.TempDir()
	batch := func() *Batch {
		return &Batch{ID: "b1", Stories: []string{"S-0001", "S-0002"}, Step: "landing S-0001", Result: &Result{Passed: true, FinishedAt: t1},
			Landing: &Landing{Stories: []string{"S-0001", "S-0002"}, Landed: []string{}, Head: "base"}}
	}
	if _, err := Update(dir, func(s *State) error {
		s.Enqueue("S-0001", "orchestrator", t0)
		s.Enqueue("S-0002", "alex", t0)
		s.Batch = batch()
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	st, err := Landed(dir, "S-1", "b1", "orchestrator", t3)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := st.EnqueuedAs("S-0001"); ok || len(st.Enqueued) != 1 {
		t.Errorf("enqueued %+v, want S-0001 taken out", st.Enqueued)
	}
	want := Outcome{Story: "S-1", Outcome: OutcomeLanded, At: t3, Detail: "landed by merge queue batch b1 as orchestrator"}
	if o, ok := outcomeOf(st, "S-0001"); !ok || o != want {
		t.Errorf("outcome %+v, want %+v", o, want)
	}
	if st.Batch == nil || st.Batch.Dropped != nil || !slices.Equal(st.Batch.Landing.Landed, []string{"S-0001"}) {
		t.Errorf("batch %+v, want it kept with S-0001 landed", st.Batch)
	}
	if st, err = Landed(dir, "S-0001", "b1", "orchestrator", t3); err != nil || !slices.Equal(st.Batch.Landing.Landed, []string{"S-0001"}) {
		t.Errorf("landed %v, %v; want S-0001 landed once", st.Batch.Landing.Landed, err)
	}

	before, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st, err = Landed(dir, "S-0002", "b0", "alex", t3); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", *st.Batch.Landing) != fmt.Sprintf("%+v", *before.Batch.Landing) || st.Batch.Dropped != nil {
		t.Errorf("batch %+v, want it unchanged by another batch's landing", st.Batch)
	}
	if o, ok := outcomeOf(st, "S-0002"); !ok || o.Outcome != OutcomeLanded || len(st.Enqueued) != 0 {
		t.Errorf("outcome %+v, enqueued %+v; want S-0002 landed and out of the queue", o, st.Enqueued)
	}
}

// Advance does nothing, running no git and accepting nothing, under
// pull-request acceptance, with no batch, one dropped, or one not verified.
func TestAdvanceDoesNothingWithNoVerifiedBatch(t *testing.T) {
	const batch = `"batch": {"id": "b1", "branch": "merge-queue/b1", "stories": ["S-0001"], "step": "verifying"`
	passed := `, "result": {"passed": true, "finished_at": "2026-10-19T17:00:00Z"}`
	for _, tc := range []struct {
		name, mode, state string
	}{
		{name: "no batch", state: `{}`},
		{name: "a batch not verified", state: `{` + batch + `}}`},
		{name: "a dropped batch", state: `{` + batch + passed + `, "dropped": {"reason": "removed", "by": "alex", "at": "2026-10-19T17:00:00Z"}}}`},
		{name: "pull-request acceptance", mode: "pull-request", state: `{` + batch + passed + `}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := project(t, tc.mode)
			writeState(t, repo.CacheDir(), tc.state)
			before, _ := os.ReadFile(StatePath(repo.CacheDir()))
			called := false
			accept := func(context.Context, string, string, string) error {
				called = true
				return nil
			}
			d := Deps{Runner: scriptedGit{}, Now: func() time.Time { return t3 }, Verify: func(context.Context, string, string, string, []string, func(string)) (verify.Report, error) {
				called = true
				return verify.Report{}, nil
			}}
			did, err := Advance(context.Background(), d, accept, repo)
			if did || err != nil || called {
				t.Errorf("did %v, %v, called %v; want nothing done", did, err, called)
			}
			if after, _ := os.ReadFile(StatePath(repo.CacheDir())); string(after) != string(before) {
				t.Errorf("state written:\n%s\nwas\n%s", after, before)
			}
		})
	}
}

// S-0402 criterion 2, ADR-0178, ADR-0248: files the main checkout has not
// committed that no story of the batch changes stop no landing, as the
// analyzer's report and issues did on TH-0462: those of the shared set,
// which the landing commits, pass the preflight's flai check, and the
// operator's are left as they were.
func TestUncommittedFilesNoStoryChangesStopNoLanding(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	writeFile(t, repo.MainRoot, "design/analysis/2026-10-19-bottlenecks.md", "---\ntitle: Bottlenecks\nupdated: 2026-10-19T08:00:00Z\nstatus: active\nfocus: bottlenecks\nfrom: 2026-09-01\nto: 2026-09-30\n---\n\n# Bottlenecks\n")
	writeFile(t, repo.MainRoot, "design/issues/I-0150-x.md", "---\nid: I-0150\ntitle: X\nclass: defect\nstatus: open\ncount: 1\nfirst_reported: 2026-10-19T12:00:00Z\nlast_reported: 2026-10-19T12:00:00Z\nupdated: 2026-10-19T12:00:00Z\n---\n\n# I-0150 X\n")
	writeFile(t, repo.MainRoot, "x.txt", "the operator's edit\n")
	before := runGit(t, repo.MainRoot, "status", "--porcelain", "--untracked-files=all")

	q, err := Compute(execx.System{}, repo)
	if err != nil || len(q.Eligible) != 2 {
		t.Fatalf("eligible %+v, ineligible %+v, %v; want both stories eligible", q.Eligible, q.Ineligible, err)
	}
	_, built := runBatch(t, deps(proc, t1), repo)
	fa := &fakeAccept{t: t, repo: repo}
	_, st := advance(t, deps(proc, t2), fa, repo)
	if want := []acceptCall{{a, "alex", built.Batch.ID}, {b, "alex", built.Batch.ID}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want %+v", fa.calls, want)
	}
	if len(st.Enqueued) != 0 || st.Batch != nil {
		t.Errorf("enqueued %+v, batch %+v; want both landed and no batch", st.Enqueued, st.Batch)
	}
	if after := runGit(t, repo.MainRoot, "status", "--porcelain", "--untracked-files=all"); after != before {
		t.Errorf("the main checkout's uncommitted files changed:\n%s\nwere\n%s", after, before)
	}
}

// S-0402 criteria 3 and 4, ADR-0178: a file a story of the batch changes,
// left uncommitted in the main checkout after the batch was built, holds that
// story back only: the rest land, the story stays enqueued with the file in
// its outcome, and the queue names it, the file, and the command that frees
// it, until the file is clean; then it is batched again.
func TestAnUncommittedFileAStoryChangesHoldsBackThatStoryOnly(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	_, built := runBatch(t, deps(proc, t1), repo)
	if bt := built.Batch; bt == nil || bt.Step != StepVerified || !slices.Equal(bt.Stories, []string{a, b}) {
		t.Fatalf("batch %+v, want %s and %s verified", bt, a, b)
	}
	writeFile(t, repo.MainRoot, "a.txt", "the operator's own a\n")

	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, deps(proc, t2), fa, repo)
	if !did {
		t.Error("Advance did nothing")
	}
	if want := []acceptCall{{b, "alex", built.Batch.ID}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want only %+v", fa.calls, want)
	}
	if len(st.Enqueued) != 1 || st.Enqueued[0].Story != a || st.Batch != nil {
		t.Errorf("enqueued %+v, batch %+v; want %s still enqueued and the batch gone", st.Enqueued, st.Batch, a)
	}
	o, ok := outcomeOf(st, a)
	if !ok || o.Outcome != OutcomeHeldBack || !o.At.Equal(t2) || o.Detail != HeldBackText(a, []string{"a.txt"}) {
		t.Errorf("%s's outcome %+v, want held back on a.txt", a, o)
	}
	if !strings.Contains(o.Detail, "a.txt") || !strings.Contains(o.Detail, "git stash push -- a.txt") {
		t.Errorf("the outcome must name the file and the command that frees it: %q", o.Detail)
	}
	if b, _ := os.ReadFile(filepath.Join(repo.MainRoot, "a.txt")); string(b) != "the operator's own a\n" {
		t.Errorf("the operator's file was changed: %q", b)
	}

	q, err := Compute(execx.System{}, repo)
	if err != nil {
		t.Fatal(err)
	}
	var reason *Reason
	for _, e := range q.Ineligible {
		for i := range e.Reasons {
			if e.Story == a && e.Reasons[i].Kind == KindDirtyMain {
				reason = &e.Reasons[i]
			}
		}
	}
	if reason == nil || !slices.Equal(reason.Paths, []string{"a.txt"}) || reason.Command != "git stash push -- a.txt" || reason.Text != HeldBackText(a, []string{"a.txt"}) {
		t.Fatalf("the queue must name %s, a.txt, and the command that frees it: ineligible %+v", a, q.Ineligible)
	}
	if ran, again := runBatch(t, deps(proc, t3), repo); ran || again.Batch != nil {
		t.Errorf("a held story is not batched while the file is dirty: %+v", again.Batch)
	}

	if err := os.Remove(filepath.Join(repo.MainRoot, "a.txt")); err != nil {
		t.Fatal(err)
	}
	// the fake acceptance leaves B in review, so it is listed too
	if q, err := Compute(execx.System{}, repo); err != nil || len(q.Ineligible) != 0 || !slices.ContainsFunc(q.Eligible, func(e Entry) bool { return e.Story == a }) {
		t.Fatalf("once the file is clean %s is eligible again: eligible %+v, ineligible %+v, %v", a, q.Eligible, q.Ineligible, err)
	}
	if _, again := runBatch(t, deps(proc, t3), repo); again.Batch == nil || !slices.Equal(again.Batch.Stories, []string{a}) {
		t.Errorf("once the file is clean %s is batched again: %+v", a, again.Batch)
	}
}

// ADR-0248: a file of the shared set a story of the batch changes, left
// uncommitted in the main checkout under design after the batch was built,
// holds nothing back: the landing commits it, so the story lands, and the
// queue's eligibility counts it as clean.
func TestAnUncommittedSharedFileAStoryChangesHoldsNothingBack(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "design/d.md", "# d\n")
	enqueue(t, repo, a)
	_, built := runBatch(t, deps(proc, t1), repo)
	if bt := built.Batch; bt == nil || bt.Step != StepVerified {
		t.Fatalf("batch %+v, want %s verified", bt, a)
	}
	writeFile(t, repo.MainRoot, "design/d.md", "# flai's own d\n")

	q, err := Compute(execx.System{}, repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.Ineligible) != 0 {
		t.Errorf("a dirty file under design holds no story back from the queue: ineligible %+v", q.Ineligible)
	}
	fa := &fakeAccept{t: t, repo: repo}
	did, st := advance(t, deps(proc, t2), fa, repo)
	if !did {
		t.Error("Advance did nothing")
	}
	if want := []acceptCall{{a, "alex", built.Batch.ID}}; !slices.Equal(fa.calls, want) {
		t.Errorf("accepted %+v, want %+v", fa.calls, want)
	}
	if len(st.Enqueued) != 0 || st.Batch != nil {
		t.Errorf("enqueued %+v, batch %+v; want %s landed and no batch", st.Enqueued, st.Batch, a)
	}
	if o, ok := outcomeOf(st, a); !ok || o.Outcome != OutcomeLanded {
		t.Errorf("%s's outcome %+v, want landed", a, o)
	}
}

// ADR-0178: the landing names no held story next, and the command that frees
// a held story quotes each path the shell would split.
func TestALandingSkipsAHeldStoryAndQuotesTheCommandThatFreesIt(t *testing.T) {
	l := &Landing{Stories: []string{"S-0001", "S-0002", "S-0003"}, Landed: []string{"S-0001"}, HeldBack: []LeftOut{{Story: "S-0002", Paths: []string{"a.txt"}}}}
	if next, ok := l.Next(); !ok || next != "S-0003" {
		t.Errorf("Next = %s, %v; want S-0003, past the held S-0002", next, ok)
	}
	l.Landed = append(l.Landed, "S-0003")
	if next, ok := l.Next(); ok {
		t.Errorf("Next = %s; want none once every story not held has landed", next)
	}
	if got, want := FreeCommand([]string{"design/a.md", "docs/it's here.md"}), `git stash push -- design/a.md 'docs/it'\''s here.md'`; got != want {
		t.Errorf("FreeCommand = %s, want %s", got, want)
	}
}

// keptBatch is a project with two stories enqueued, A and B, and a batch of
// them that passed its verification at t1 on main's head c0ffee, with each
// story's verified and rebased heads, its git scripted in g, whose main head
// a test moves by scripting it again. Its deps run g at the clock *at, with
// a preflight that answers *findings and counts its looks, a verifier that
// counts its runs, and a JSON log.
type keptBatch struct {
	repo     *workitem.Repo
	g        scriptedGit
	a, b     string
	id       string
	at       time.Time
	findings []verify.FileFinding
	looks    int
	verified int
	logs     bytes.Buffer
}

// mainHeadArgs are the arguments of the git command that reads main's head.
var mainHeadArgs = []string{"rev-parse", "--verify", "--quiet", "refs/heads/main^{commit}"}

func newKeptBatch(t *testing.T) *keptBatch {
	t.Helper()
	k := &keptBatch{repo: project(t, ""), g: scriptedGit{}, id: "20261019-160000", at: t2}
	k.g.script("main", "rev-parse", "--abbrev-ref", "HEAD")
	k.g.script("c0ffee", mainHeadArgs...)
	k.g.script("", "-c", "core.quotePath=false", "status", "--porcelain", "--no-renames", "--untracked-files=all")
	k.g.script("", "worktree", "prune")
	sa, _ := story(t, k.repo, "A", []bool{true}, t0.Add(time.Hour))
	sb, _ := story(t, k.repo, "B", []bool{true}, t0.Add(2*time.Hour))
	k.a, k.b = sa.ID, sb.ID
	passed(t, k.repo, k.g, k.a, "aaaa")
	passed(t, k.repo, k.g, k.b, "bbbb")
	enqueue(t, k.repo, k.a, k.b)
	if _, err := Update(k.repo.CacheDir(), func(s *State) error {
		s.Batch = &Batch{ID: k.id, Branch: BranchPrefix + k.id, Base: "c0ffee", Stories: []string{k.a, k.b}, StartedAt: t1, Step: StepVerified,
			Result: &Result{Passed: true, Summary: "passed unit", FinishedAt: t1},
			Heads:  []Head{{Story: k.a, Verified: "aaaa", Rebased: "a1a1"}, {Story: k.b, Verified: "bbbb", Rebased: "b1b1"}}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return k
}

func (k *keptBatch) deps() Deps {
	return Deps{Runner: k.g, Now: func() time.Time { return k.at }, Log: slog.New(slog.NewJSONHandler(&k.logs, nil)),
		Verify: func(context.Context, string, string, string, []string, func(string)) (verify.Report, error) {
			k.verified++
			return verify.Report{Passed: true}, nil
		},
		Preflight: func(context.Context, *workitem.Repo) ([]verify.FileFinding, error) {
			k.looks++
			return slices.Clone(k.findings), nil
		}}
}

// accepter is a fake flai accept --batch on k: it records each call, fails
// with the error fail names for the story, and otherwise checks with
// Landable that the story may land, moves main's head on, and records the
// landing with Landed.
func (k *keptBatch) accepter(calls *[]acceptCall, fail map[string]error) Accepter {
	return func(_ context.Context, story, by, batch string) error {
		*calls = append(*calls, acceptCall{story, by, batch})
		if err := fail[story]; err != nil {
			return err
		}
		if _, reasons, err := Landable(k.g, k.repo, story, batch); err != nil || len(reasons) > 0 {
			return fmt.Errorf("%s may not land: %q, %w", story, reasons, err)
		}
		k.g.script("landed-"+story, mainHeadArgs...)
		_, err := Landed(k.repo.CacheDir(), story, batch, by, k.at)
		return err
	}
}

// state is k's merge queue state, failing the test when it cannot be read.
func (k *keptBatch) state(t *testing.T) State {
	t.Helper()
	st, err := ReadState(k.repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	return st
}

// mainNote is the preflight's finding of an uncommitted task note in the
// main checkout whose link breaks MD051, which the lint fix does not clear.
var mainNote = verify.FileFinding{Tier: "markdown", Checkout: verify.CheckoutMain, Path: "wip/kanban/tasks/T-0009-a-note.md", Line: 7, Level: "error", Rule: "MD051", Message: "Link fragments should be valid"}

// errLintRefusal is flai accept --batch refusing a landing for the note's lint.
var errLintRefusal = errors.New("the markdown tier still fails on files the fix does not clear:\nwip/kanban/tasks/T-0009-a-note.md:7 MD051")

// ADR-0234, S-0549 criterion 3: a passing batch whose first landing is
// refused for an uncommitted note's lint, while the preflight finds the
// note, is kept with its landing and its heads, stalled at that story, with
// the queue held on the note; nothing more lands while the hold stands, and
// a hand acceptance is told why. Once the note is fixed, a look clears the
// hold with no batch built, and the landing goes on from the commits the
// batch verified, logging landing_resumed; no tier runs at any point.
func TestALandingThatFailedOnMainLandsAgainWithoutANewVerification(t *testing.T) {
	k := newKeptBatch(t)
	k.findings = []verify.FileFinding{mainNote}
	var calls []acceptCall
	accept := k.accepter(&calls, map[string]error{k.a: errLintRefusal})

	did, err := Advance(context.Background(), k.deps(), accept, k.repo)
	if err != nil || !did {
		t.Fatalf("did %v, %v; want the landing begun and no error, as a hold is a state", did, err)
	}
	st := k.state(t)
	bt := st.Batch
	if bt == nil || bt.Dropped != nil || bt.Landing == nil || bt.Landing.Head != "c0ffee" || len(bt.Landing.Landed) != 0 || len(bt.Heads) != 2 {
		t.Fatalf("batch %+v; want it kept with its landing at c0ffee and its heads", bt)
	}
	if want := (&Stalled{Story: k.a, Reason: said(errLintRefusal), At: t2}); !reflect.DeepEqual(bt.Stalled, want) {
		t.Errorf("stalled %+v, want %+v", bt.Stalled, want)
	}
	if want := (&PreflightHold{Since: t2, Main: "c0ffee", Findings: []verify.FileFinding{mainNote}}); !reflect.DeepEqual(st.Preflight, want) {
		t.Errorf("hold %+v, want %+v", st.Preflight, want)
	}
	if _, ok := st.EnqueuedAs(k.a); !ok || k.looks != 1 || k.verified != 0 || len(calls) != 1 {
		t.Errorf("%s enqueued %v, looks %d, verified %d, accepted %+v; want it still enqueued, one look, no tier run, one acceptance", k.a, ok, k.looks, k.verified, calls)
	}
	if !strings.Contains(k.logs.String(), "merge queue landing waits on main") || !strings.Contains(k.logs.String(), mainNote.Path) {
		t.Errorf("logs %s; want a warning that the landing waits on main, naming %s", k.logs.String(), mainNote.Path)
	}

	// the hold stands: nothing lands, and a hand acceptance is told why
	k.at = t2.Add(time.Minute)
	if ran, _ := runBatch(t, k.deps(), k.repo); ran || k.looks != 2 || k.verified != 0 {
		t.Fatalf("ran %v, looks %d, verified %d; want the look to keep the hold and build nothing", ran, k.looks, k.verified)
	}
	if _, err := Advance(context.Background(), k.deps(), accept, k.repo); err != nil || len(calls) != 1 {
		t.Fatalf("%v, accepted %+v; want nothing landed while the preflight holds the queue", err, calls)
	}
	if _, reasons, err := Landable(k.g, k.repo, k.a, k.id); err != nil || len(reasons) != 1 || !strings.Contains(reasons[0], "waits for main to be clean") || !strings.Contains(reasons[0], mainNote.Path) {
		t.Errorf("reasons %q, %v; want one saying the landing waits for main to be clean, naming %s", reasons, err, mainNote.Path)
	}

	// the note is fixed: the look clears the hold and builds nothing
	k.findings, k.at = nil, t3
	if ran, st := runBatch(t, k.deps(), k.repo); ran || st.Preflight != nil || st.Batch == nil || st.Batch.ID != k.id || st.Batch.Dropped != nil || k.verified != 0 {
		t.Fatalf("ran %v, hold %+v, batch %+v, verified %d; want the hold cleared, the batch kept, and no tier run", ran, st.Preflight, st.Batch, k.verified)
	}
	if _, reasons, err := Landable(k.g, k.repo, k.a, k.id); err != nil || len(reasons) != 0 {
		t.Errorf("reasons %q, %v; want %s landable once main is clean", reasons, err, k.a)
	}

	accept = k.accepter(&calls, nil)
	did, err = Advance(context.Background(), k.deps(), accept, k.repo)
	if err != nil || !did {
		t.Fatalf("did %v, %v; want the batch landed", did, err)
	}
	if want := []acceptCall{{k.a, "alex", k.id}, {k.a, "alex", k.id}, {k.b, "alex", k.id}}; !slices.Equal(calls, want) {
		t.Errorf("accepted %+v, want %+v", calls, want)
	}
	if st := k.state(t); st.Batch != nil || len(st.Enqueued) != 0 || k.verified != 0 {
		t.Errorf("batch %+v, enqueued %+v, verified %d; want both landed, the batch gone, and no tier run", st.Batch, st.Enqueued, k.verified)
	}
	var resumed []Event
	for _, e := range eventLog(t, k.repo.CacheDir()) {
		if e.Event == EventLandingResumed {
			resumed = append(resumed, e)
		}
	}
	// Update stamps a write's events with the wall clock, so the wait is
	// measured from the stall to now
	if len(resumed) != 1 || resumed[0].Batch != k.id || !slices.Equal(resumed[0].Stories, []string{k.a}) || resumed[0].Waited <= 0 {
		t.Errorf("landing_resumed events %+v; want one of %s at %s with how long it waited", resumed, k.id, k.a)
	}
	if !strings.Contains(k.logs.String(), "merge queue batch lands again without a new verification") {
		t.Errorf("logs %s; want that the batch lands again without a new verification", k.logs.String())
	}
}

// ADR-0234: a landing kept on main's state is dropped as before when main
// moved meanwhile, as by committing the fix on it, and nothing lands or is
// verified; a hand acceptance is told main moved.
func TestALandingKeptOnMainIsDroppedWhenMainMoved(t *testing.T) {
	k := newKeptBatch(t)
	k.findings = []verify.FileFinding{mainNote}
	var calls []acceptCall
	accept := k.accepter(&calls, map[string]error{k.a: errLintRefusal})
	if _, err := Advance(context.Background(), k.deps(), accept, k.repo); err != nil {
		t.Fatal(err)
	}
	k.g.script("beef", mainHeadArgs...)
	k.findings, k.at = nil, t3
	runBatch(t, k.deps(), k.repo)
	if _, reasons, err := Landable(k.g, k.repo, k.a, k.id); err != nil || len(reasons) != 1 || !strings.Contains(reasons[0], "main moved") {
		t.Errorf("reasons %q, %v; want one saying main moved", reasons, err)
	}

	did, err := Advance(context.Background(), k.deps(), k.accepter(&calls, nil), k.repo)
	if err != nil || !did {
		t.Fatalf("did %v, %v; want the batch dropped", did, err)
	}
	st := k.state(t)
	if st.Batch == nil || st.Batch.Dropped == nil || !strings.Contains(st.Batch.Dropped.Reason, "moved from c0ffee to beef") {
		t.Fatalf("batch %+v; want it dropped as main moved", st.Batch)
	}
	if len(calls) != 1 || k.verified != 0 {
		t.Errorf("accepted %+v, verified %d; want nothing more landed and no tier run", calls, k.verified)
	}
	for _, e := range eventLog(t, k.repo.CacheDir()) {
		if e.Event == EventLandingResumed {
			t.Errorf("logged %+v; a dropped landing does not resume", e)
		}
	}
}

// ADR-0176, ADR-0234: a landing that fails while the preflight finds the
// main checkout clean, or cannot run, failed on its stories: the batch is
// dropped as before, with the acceptance's error, and nothing holds the
// queue or stalls.
func TestALandingThatFailsWithMainCleanDropsTheBatch(t *testing.T) {
	for _, tc := range []struct {
		name      string
		preflight Preflighter
	}{
		{name: "main clean", preflight: func(context.Context, *workitem.Repo) ([]verify.FileFinding, error) { return nil, nil }},
		{name: "preflight cannot run", preflight: func(context.Context, *workitem.Repo) ([]verify.FileFinding, error) {
			return nil, errors.New("run the merge queue's preflight: no git")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			k := newKeptBatch(t)
			var calls []acceptCall
			d := k.deps()
			d.Preflight = tc.preflight
			_, err := Advance(context.Background(), d, k.accepter(&calls, map[string]error{k.a: errors.New("merge stopped\non x.txt")}), k.repo)
			if err == nil || !strings.Contains(err.Error(), "merge stopped") || !strings.Contains(err.Error(), "flai accept "+k.a) {
				t.Errorf("error %v, want the acceptance's, saying what to do", err)
			}
			st := k.state(t)
			dr := st.Batch.Dropped
			if dr == nil || dr.Reason != "landing "+k.a+" failed: merge stopped on x.txt; the stories not landed are batched again from the main branch" {
				t.Errorf("dropped %+v, want it dropped as the landing of %s failed", dr, k.a)
			}
			if st.Preflight != nil || st.Batch.Stalled != nil || k.verified != 0 {
				t.Errorf("hold %+v, stalled %+v, verified %d; want neither, and no tier run", st.Preflight, st.Batch.Stalled, k.verified)
			}
		})
	}
}

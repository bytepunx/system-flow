package mergequeue

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	runGit(t, repo.MainRoot, "commit", "-q", "--allow-empty", "-m", "accepted by hand")
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
	for _, says := range []string{"moved from " + short(built.Batch.Base) + " to " + short(moved), built.Batch.ID, "built again from the new main"} {
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
	runGit(t, repo.MainRoot, "commit", "-q", "--allow-empty", "-m", "accepted by hand")

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

	runGit(t, repo.MainRoot, "commit", "-q", "--allow-empty", "-m", "accepted by hand")
	_, reasons = landable(a, bt.ID)
	says("main moved", reasons, "main branch main", short(bt.Base))

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
			d := Deps{Runner: scriptedGit{}, Now: func() time.Time { return t3 }, Verify: func(context.Context, string, string, func(string)) (verify.Report, error) {
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

// S-0402 criterion 2, ADR-0178: files the main checkout has not committed
// that no story of the batch changes stop no landing, as the analyzer's
// report and issues did on TH-0462; they are left as they were.
func TestUncommittedFilesNoStoryChangesStopNoLanding(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	writeFile(t, repo.MainRoot, "design/analysis/2026-10-19-all.md", "# report\n")
	writeFile(t, repo.MainRoot, "design/issues/I-0150-x.md", "# issue\n")
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

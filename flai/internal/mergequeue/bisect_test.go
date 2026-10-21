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
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// t2 is when the failed batches are bisected, after they were verified.
var t2 = t1.Add(time.Hour)

// fileVerifier verifies a worktree by the files in it: its unit tier fails
// when every one of files is there, and logs the worktree's head to unit.log
// in the folder it is given, standing for the tier's output. It keeps the
// batch's step in the state as each verification runs; cancel, when set,
// ends the run's context as one runs, and drop, when set, runs once as one
// runs.
type fileVerifier struct {
	cacheDir string
	files    []string
	cancel   context.CancelFunc
	drop     func()
	steps    []string
}

func (v *fileVerifier) verify(_ context.Context, dir, _, logDir string, started func(string)) (verify.Report, error) {
	started("unit")
	if head, err := (execx.System{}).Run(dir, "git", "rev-parse", "HEAD"); err == nil {
		if err := os.MkdirAll(logDir, 0o755); err == nil {
			_ = os.WriteFile(filepath.Join(logDir, "unit.log"), []byte("unit at "+head), 0o644)
		}
	}
	if st, err := ReadState(v.cacheDir); err == nil && st.Batch != nil {
		v.steps = append(v.steps, st.Batch.Step)
	}
	if v.cancel != nil {
		v.cancel()
	}
	if v.drop != nil {
		v.drop()
		v.drop = nil
	}
	there := slices.IndexFunc(v.files, func(f string) bool {
		_, err := os.Stat(filepath.Join(dir, f))
		return err != nil
	}) < 0
	if !there {
		return verify.Report{Passed: true, Steps: []verify.Step{{Name: "unit", Tier: true, State: verify.Passed}}}, nil
	}
	finding := verify.Finding{Name: "TestFiles", Message: strings.Join(v.files, " and ") + " are there"}
	return verify.Report{StoppedAt: "unit", Steps: []verify.Step{{Name: "unit", Tier: true, State: verify.Failed, Findings: []verify.Finding{finding}}}}, nil
}

// fileDeps run real git at the clock at, verifying with v.
func fileDeps(v *fileVerifier, at time.Time) Deps {
	return Deps{Runner: execx.System{}, Now: func() time.Time { return at }, Verify: v.verify}
}

// failedBatch is a project in real git with a batch of n stories, the i-th
// adding f<i>.txt, enqueued in that order, verified and failed with a
// verifier that fails when the files of the stories at failWith, from 1, are
// all there. It answers the repo, the verifier, the stories in order, their
// heads, and the state. It skips under -short.
func failedBatch(t *testing.T, n int, failWith ...int) (*workitem.Repo, *fileVerifier, []string, map[string]string, State) {
	t.Helper()
	repo, _ := gitProject(t, "")
	v := &fileVerifier{cacheDir: repo.CacheDir()}
	for _, i := range failWith {
		v.files = append(v.files, fmt.Sprintf("f%d.txt", i))
	}
	var stories []string
	heads := map[string]string{}
	for i := 1; i <= n; i++ {
		id, head := readyStory(t, repo, fmt.Sprintf("Story %d", i), t0.Add(time.Duration(i)*time.Hour), fmt.Sprintf("f%d.txt", i), "f\n")
		stories, heads[id] = append(stories, id), head
	}
	enqueue(t, repo, stories...)
	_, st := runBatch(t, fileDeps(v, t1), repo)
	if st.Batch == nil || st.Batch.Step != StepFailed || len(st.Batch.Heads) != n {
		t.Fatalf("batch %+v, want %d stories failed", st.Batch, n)
	}
	v.steps = nil
	return repo, v, stories, heads, st
}

// runBisect runs RunBisect, failing the test on an error, and answers whether
// it ran and the state it left.
func runBisect(t *testing.T, d Deps, repo *workitem.Repo) (bool, State) {
	t.Helper()
	ran, err := RunBisect(context.Background(), d, repo)
	if err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	return ran, st
}

// prefixes are how many of the batch's stories each step verified.
func prefixes(bs *Bisect) []int {
	out := []int{}
	for _, s := range bs.Steps {
		out = append(out, len(s.Stories))
	}
	return out
}

// S-0400 criterion 3: a red batch is bisected to the first story with which
// it fails, main counting as passing: each step verifies main with the
// batch's first stories at the last one's rebased head, a batch of four its
// first half first, each step written as it begins and kept with its result.
// The stories before the failing one are to land; the batch's worktree is
// left on its branch and the stories' own branches are not touched. A
// failure that needs two stories together is found at the later one.
func TestARedBatchIsBisectedToTheFirstStoryWithWhichItFails(t *testing.T) {
	for _, tc := range []struct {
		name      string
		failWith  []int
		failing   int
		steps     []int
		passed    []bool
		good, bad int
	}{
		{name: "the first story", failWith: []int{1}, failing: 1, steps: []int{2, 1}, passed: []bool{false, false}, good: 0, bad: 1},
		{name: "the second story", failWith: []int{2}, failing: 2, steps: []int{2, 1}, passed: []bool{false, true}, good: 1, bad: 2},
		{name: "the third story", failWith: []int{3}, failing: 3, steps: []int{2, 3}, passed: []bool{true, false}, good: 2, bad: 3},
		{name: "the fourth story", failWith: []int{4}, failing: 4, steps: []int{2, 3}, passed: []bool{true, true}, good: 3, bad: 4},
		{name: "the first and third stories together", failWith: []int{1, 3}, failing: 3, steps: []int{2, 3}, passed: []bool{true, false}, good: 2, bad: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, v, stories, heads, before := failedBatch(t, 4, tc.failWith...)
			ran, st := runBisect(t, fileDeps(v, t2), repo)
			bt := st.Batch
			if !ran || bt == nil || bt.ID != before.Batch.ID || bt.Bisect == nil {
				t.Fatalf("ran %v, batch %+v; want the batch bisected", ran, bt)
			}
			bs := bt.Bisect
			if got := prefixes(bs); !slices.Equal(got, tc.steps) {
				t.Errorf("steps verified %v stories, want %v", got, tc.steps)
			}
			for i, s := range bs.Steps {
				k := len(s.Stories)
				if !slices.Equal(s.Stories, stories[:k]) || s.Commit != bt.Heads[k-1].Rebased || !s.StartedAt.Equal(t2) {
					t.Errorf("step %d %+v, want the first %d stories at %s's rebased head", i, s, k, stories[k-1])
				}
				if wantPassed := i < len(tc.passed) && tc.passed[i]; s.Result == nil || s.Result.Passed != wantPassed {
					t.Errorf("step %d result %+v, want passed %v", i, s.Result, wantPassed)
				}
				// S-0457 criterion 1 (I-0172): each step's output is kept in
				// its own folder below the batch's, which its result names
				wantLog := fmt.Sprintf(".flai-cache/merge-queue/batches/%s/step-%d", bt.ID, i+1)
				if s.Result == nil || s.Result.Log != wantLog {
					t.Errorf("step %d result %+v, want its log in %s", i, s.Result, wantLog)
				}
				if data, err := os.ReadFile(filepath.Join(repo.MainRoot, filepath.FromSlash(wantLog), "unit.log")); err != nil || !strings.Contains(string(data), s.Commit) {
					t.Errorf("step %d's log %q, %v; want unit's output at %s", i, data, err, s.Commit)
				}
			}
			batchLog := ".flai-cache/merge-queue/batches/" + bt.ID
			if bt.Result == nil || bt.Result.Log != batchLog {
				t.Errorf("batch's result %+v, want its log in %s", bt.Result, batchLog)
			}
			if data, err := os.ReadFile(filepath.Join(repo.MainRoot, filepath.FromSlash(batchLog), "unit.log")); err != nil || !strings.Contains(string(data), bt.Heads[3].Rebased) {
				t.Errorf("the batch's log %q, %v; want unit's output at the batch's head", data, err)
			}
			want := stories[tc.failing-1]
			if bs.Failing != want || bs.Commit != bt.Heads[tc.failing-1].Rebased || bs.Good != tc.good || bs.Bad != tc.bad || bs.SentBack {
				t.Errorf("bisection %+v, want %s found at its rebased head %s, good %d, bad %d", bs, want, bt.Heads[tc.failing-1].Rebased, tc.good, tc.bad)
			}
			if r := bs.Result; r == nil || r.Passed || r.StoppedAt != "unit" || !strings.Contains(r.Summary, "TestFiles") {
				t.Errorf("result at the failing story %+v, want unit failed", r)
			}
			if l := bt.Landing; l == nil || !slices.Equal(l.Stories, stories[:tc.failing-1]) || len(l.Landed) != 0 || l.Head != bt.Base {
				t.Errorf("landing %+v, want %v from %s", l, stories[:tc.failing-1], bt.Base)
			}
			if bt.Step != StepBisected || bt.Result == nil || bt.Result.Passed || bt.Dropped != nil {
				t.Errorf("batch at %q, result %+v, dropped %+v; want bisected and its own result kept", bt.Step, bt.Result, bt.Dropped)
			}
			wantSteps := []string{}
			for _, k := range tc.steps {
				wantSteps = append(wantSteps, fmt.Sprintf("bisecting %d of 4, verifying unit", k))
			}
			if !slices.Equal(v.steps, wantSteps) {
				t.Errorf("steps as each ran %v, want %v", v.steps, wantSteps)
			}
			wt := filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)
			if got := runGit(t, wt, "rev-parse", "--abbrev-ref", "HEAD"); got != bt.Branch {
				t.Errorf("the batch's worktree is on %s, want %s", got, bt.Branch)
			}
			if got := runGit(t, wt, "rev-parse", "HEAD"); got != bt.Heads[3].Rebased {
				t.Errorf("the batch's branch is at %s, want its head %s", got, bt.Heads[3].Rebased)
			}
			untouched(t, repo, heads)
			if len(st.Enqueued) != 4 || len(st.Outcomes) != 0 {
				t.Errorf("enqueued %+v, outcomes %+v; want every story still enqueued and no outcome", st.Enqueued, st.Outcomes)
			}

			// a bisection found is not run again
			v.steps = nil
			if again, _ := runBisect(t, fileDeps(v, t2.Add(time.Hour)), repo); again || len(v.steps) != 0 {
				t.Errorf("ran %v, verified %v; want a found bisection left as it is", again, v.steps)
			}
		})
	}
}

// S-0400 criterion 3: a bisection stopped while a step verifies leaves that
// step without a result; the next run verifies it again in place, rather
// than adding a second, and finishes.
func TestAStoppedBisectionIsResumedInPlace(t *testing.T) {
	repo, v, stories, _, _ := failedBatch(t, 4, 3)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	ran, err := RunBisect(ctx, fileDeps(v, t2), repo)
	if !ran || !errors.Is(err, context.Canceled) {
		t.Fatalf("ran %v, %v; want the bisection stopped", ran, err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if bs := st.Batch.Bisect; bs == nil || len(bs.Steps) != 1 || bs.Steps[0].Result != nil || bs.Found() || st.Batch.Step != "bisecting 2 of 4, verifying unit" {
		t.Fatalf("batch at %q, bisection %+v; want one step without a result", st.Batch.Step, st.Batch.Bisect)
	}

	v.cancel = nil
	later := t2.Add(time.Hour)
	ran, st = runBisect(t, fileDeps(v, later), repo)
	bs := st.Batch.Bisect
	if !ran || !slices.Equal(prefixes(bs), []int{2, 3}) || bs.Failing != stories[2] || st.Batch.Step != StepBisected {
		t.Fatalf("ran %v, bisection %+v; want steps of 2 and 3 and %s found", ran, bs, stories[2])
	}
	if !bs.Steps[0].StartedAt.Equal(later) || bs.Steps[0].Result == nil || !bs.Steps[0].Result.Passed {
		t.Errorf("first step %+v, want it verified again at %s and passed", bs.Steps[0], later)
	}
}

// S-0400 criterion 3: a batch of one that failed is bisected to its story
// with no step, the batch's own result the result there, and nothing to land.
func TestABatchOfOneIsBisectedWithNoStep(t *testing.T) {
	repo, v, stories, _, before := failedBatch(t, 1, 1)
	ran, st := runBisect(t, fileDeps(v, t2), repo)
	bt := st.Batch
	bs := bt.Bisect
	if !ran || bs == nil || bs.Failing != stories[0] || bs.Commit != bt.Heads[0].Rebased || len(bs.Steps) != 0 || bs.Good != 0 || bs.Bad != 1 {
		t.Fatalf("ran %v, bisection %+v; want %s found with no step", ran, bs, stories[0])
	}
	if bs.Result == nil || *bs.Result != *before.Batch.Result || len(v.steps) != 0 {
		t.Errorf("result %+v, verified %v; want the batch's own result %+v and nothing verified", bs.Result, v.steps, before.Batch.Result)
	}
	if bt.Landing == nil || bt.Landing.Stories == nil || len(bt.Landing.Stories) != 0 || bt.Step != StepBisected {
		t.Errorf("landing %+v at %q, want none to land", bt.Landing, bt.Step)
	}
}

// S-0400 criterion 3: a batch dropped while a step verifies keeps that step's
// result and is bisected no further.
func TestADroppedBatchIsBisectedNoFurther(t *testing.T) {
	repo, v, stories, _, _ := failedBatch(t, 4, 3)
	v.drop = func() {
		if _, err := AcceptedByHand(repo.CacheDir(), stories[0], "alex", t2); err != nil {
			t.Fatal(err)
		}
	}
	ran, st := runBisect(t, fileDeps(v, t2), repo)
	bs := st.Batch.Bisect
	if !ran || st.Batch.Dropped == nil || bs == nil || len(bs.Steps) != 1 || bs.Steps[0].Result == nil || bs.Found() {
		t.Errorf("ran %v, batch %+v; want it dropped after one step, with its result", ran, st.Batch)
	}
}

// RunBisect does nothing, running no git and verifying nothing, for a batch
// that passed, is still verifying, was dropped, or was bisected already, and
// under pull-request acceptance; a failed batch whose heads do not match its
// stories is an error saying it must be built again.
func TestRunBisectDoesNothingWhenThereIsNothingToBisect(t *testing.T) {
	const heads = `"heads": [{"story": "S-0001", "verified": "v1", "rebased": "r1"}, {"story": "S-0002", "verified": "v2", "rebased": "r2"}]`
	batch := func(extra string) string {
		return `{"batch": {"id": "b1", "branch": "merge-queue/b1", "stories": ["S-0001", "S-0002"], "step": "failed", ` + heads + extra + `}}`
	}
	failed := `, "result": {"passed": false, "stopped_at": "unit", "finished_at": "2026-10-19T17:00:00Z"}`
	for _, tc := range []struct {
		name, mode, state string
	}{
		{name: "no batch", state: `{}`},
		{name: "a batch that passed", state: batch(`, "result": {"passed": true, "finished_at": "2026-10-19T17:00:00Z"}`)},
		{name: "a batch still verifying", state: batch(``)},
		{name: "a dropped batch", state: batch(failed + `, "dropped": {"reason": "removed", "by": "alex", "at": "2026-10-19T17:00:00Z"}`)},
		{name: "a batch bisected already", state: batch(failed + `, "bisect": {"steps": [], "good": 1, "bad": 2, "failing": "S-0002", "commit": "r2"}`)},
		{name: "pull-request acceptance", mode: "pull-request", state: batch(failed)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := project(t, tc.mode)
			writeState(t, repo.CacheDir(), tc.state)
			before, _ := os.ReadFile(StatePath(repo.CacheDir()))
			called := false
			d := Deps{Runner: scriptedGit{}, Now: func() time.Time { return t2 }, Verify: func(context.Context, string, string, string, func(string)) (verify.Report, error) {
				called = true
				return verify.Report{}, nil
			}}
			ran, err := RunBisect(context.Background(), d, repo)
			if ran || err != nil || called {
				t.Errorf("ran %v, %v, verified %v; want nothing done", ran, err, called)
			}
			if after, _ := os.ReadFile(StatePath(repo.CacheDir())); string(after) != string(before) {
				t.Errorf("state written:\n%s\nwas\n%s", after, before)
			}
		})
	}

	repo := project(t, "")
	writeState(t, repo.CacheDir(), `{"batch": {"id": "b1", "stories": ["S-0001", "S-0002"], "step": "failed", "heads": [{"story": "S-0002", "rebased": "r2"}]`+failed+`}}`)
	_, err := RunBisect(context.Background(), Deps{Runner: scriptedGit{}}, repo)
	if err == nil || !strings.Contains(err.Error(), "b1") || !strings.Contains(err.Error(), "built again") {
		t.Errorf("error %v, want one saying batch b1 must be built again", err)
	}
}

// bisected is a project, not in git, with a story before and the failing
// story, both in review and enqueued, in a batch bisected to the failing one,
// failed at unit. It answers the repo and the two stories.
func bisected(t *testing.T) (*workitem.Repo, string, string) {
	t.Helper()
	repo := project(t, "")
	first, _ := story(t, repo, "Lands", []bool{true}, t0.Add(time.Hour))
	failing, _ := story(t, repo, "Fails", []bool{true}, t0.Add(2*time.Hour))
	enqueue(t, repo, first.ID, failing.ID)
	res := &Result{StoppedAt: "unit", Summary: "unit failed: TestFiles: f2.txt is there", FinishedAt: t1}
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		s.Batch = &Batch{
			ID: "20261019-160000", Branch: "merge-queue/20261019-160000", Base: "ba5e0000000000000000", Stories: []string{first.ID, failing.ID}, Step: StepBisected, Result: res,
			Heads: []Head{{Story: first.ID, Verified: "f1f1f1f1f1f1f1f1", Rebased: "r1r1r1r1r1r1r1r1"}, {Story: failing.ID, Verified: "fa11fa11fa11fa11", Rebased: "4e4e4e4e4e4e4e4e"}},
			Bisect: &Bisect{Steps: []BisectStep{{Stories: []string{first.ID}, Commit: "r1r1r1r1r1r1r1r1", StartedAt: t2, Result: &Result{Passed: true, FinishedAt: t2}}},
				Good: 1, Bad: 2, Failing: failing.ID, Commit: "4e4e4e4e4e4e4e4e", Result: res},
			Landing: &Landing{Stories: []string{first.ID}, Landed: []string{}, Head: "ba5e0000000000000000"},
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return repo, first.ID, failing.ID
}

// threadsBy are the threads on story that author opened.
func threadsBy(t *testing.T, repo *workitem.Repo, story, author string) []*threads.Thread {
	t.Helper()
	all, err := threads.For(repo, story)
	if err != nil {
		t.Fatal(err)
	}
	var out []*threads.Thread
	for _, th := range all {
		if th.Opener() == author {
			out = append(out, th)
		}
	}
	return out
}

// S-0400 criterion 4: the story a red batch is bisected to moves from review
// to in-progress by flai, the reason naming the batch, the commit, the
// failing tier, and its first findings appended to its notes; a thread on it
// to its agent by flai serve says what failed, at which commit, with which
// stories, and what to do. It leaves the queue as failing, the story before
// it stays enqueued to land, and it is sent back once.
func TestTheFailingStoryIsSentBackOnceWithAThread(t *testing.T) {
	repo, first, failing := bisected(t)
	d := Deps{Runner: scriptedGit{}, Now: func() time.Time { return t2 }}

	sent, err := SendBack(d, repo)
	if !sent || err != nil {
		t.Fatalf("sent %v, %v; want %s sent back", sent, err, failing)
	}
	it, err := repo.Get(failing)
	if err != nil {
		t.Fatal(err)
	}
	last := it.Transitions[len(it.Transitions)-1]
	if it.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != "flai" {
		t.Errorf("%s is %s, last moved %+v; want in-progress by flai", failing, it.Status, last)
	}
	reason := "moved to in-progress: merge queue batch 20261019-160000 fails with " + failing + " at 4e4e4e4e4e4e: unit failed: TestFiles: f2.txt is there"
	if !strings.Contains(it.Body, reason) {
		t.Errorf("%s's notes do not say %q:\n%s", failing, reason, it.Body)
	}
	ths := threadsBy(t, repo, failing, "flai serve")
	if len(ths) != 1 {
		t.Fatalf("threads by flai serve on %s: %d, want 1", failing, len(ths))
	}
	th := ths[0]
	if !th.Open() || th.Title != "Merge queue batch 20261019-160000 fails with "+failing+" at unit" {
		t.Errorf("thread %q open %v", th.Title, th.Open())
	}
	for _, says := range []string{"`unit`", "`4e4e4e4e4e4e`", "`fa11fa11fa11`", "`ba5e00000000`", first, "TestFiles: f2.txt is there", "story/" + failing, "flai verify " + failing, "flai queue add " + failing} {
		if !strings.Contains(th.Body, says) {
			t.Errorf("thread does not say %q:\n%s", says, th.Body)
		}
	}

	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := st.EnqueuedAs(failing); ok {
		t.Errorf("%s is still enqueued", failing)
	}
	if _, ok := st.EnqueuedAs(first); !ok {
		t.Errorf("%s, to land, left the queue", first)
	}
	if len(st.Outcomes) != 1 {
		t.Fatalf("outcomes %+v, want %s's alone", st.Outcomes, failing)
	}
	o := st.Outcomes[0]
	if o.Story != failing || o.Outcome != OutcomeFailing || !o.At.Equal(t2) || o.Base != "ba5e0000000000000000" || o.Head != "fa11fa11fa11fa11" ||
		o.Detail != "sent back to in-progress: merge queue batch 20261019-160000 fails with it at unit" {
		t.Errorf("outcome %+v", o)
	}
	if !st.Batch.Bisect.SentBack || st.Batch.Dropped != nil || st.Batch.Landing == nil {
		t.Errorf("batch %+v, want it sent back and its landing kept", st.Batch)
	}

	again, err := SendBack(Deps{Now: func() time.Time { return t2.Add(time.Hour) }}, repo)
	if again || err != nil {
		t.Errorf("sent again %v, %v; want nothing done", again, err)
	}
	if it2, _ := repo.Get(failing); len(it2.Transitions) != len(it.Transitions) || len(threadsBy(t, repo, failing, "flai serve")) != 1 {
		t.Errorf("a second send-back moved %s or opened a thread", failing)
	}
}

// S-0400 criterion 4: a failing story no longer in review, moved by hand, is
// not moved and gets no thread, but leaves the queue as failing.
func TestAFailingStoryNotInReviewIsNotMovedButLeavesTheQueue(t *testing.T) {
	repo, _, failing := bisected(t)
	it, err := repo.Get(failing)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Transition(it, workitem.InProgress, "alex", "rework", t2); err != nil {
		t.Fatal(err)
	}
	moves := len(it.Transitions)

	sent, err := SendBack(Deps{Now: func() time.Time { return t2 }}, repo)
	if sent || err != nil {
		t.Fatalf("sent %v, %v; want it not sent back", sent, err)
	}
	if it, _ = repo.Get(failing); len(it.Transitions) != moves || len(threadsBy(t, repo, failing, "flai serve")) != 0 {
		t.Errorf("%s moved %d times, want %d, or got a thread", failing, len(it.Transitions), moves)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := st.EnqueuedAs(failing); ok || !st.Batch.Bisect.SentBack {
		t.Errorf("%s still enqueued, or the bisection %+v not marked sent back", failing, st.Batch.Bisect)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Outcome != OutcomeFailing || !strings.Contains(st.Outcomes[0].Detail, "it was in-progress, so it was not sent back") {
		t.Errorf("outcomes %+v, want %s failing, not sent back", st.Outcomes, failing)
	}
}

// sentBackAt records story's outcome as failing in batch 20261019-160000 at
// head, as SendBack does, and opens its thread by flai serve; it answers the
// thread.
func sentBackAt(t *testing.T, repo *workitem.Repo, story, head string) *threads.Thread {
	t.Helper()
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		s.RecordOutcome(Outcome{Story: story, Outcome: OutcomeFailing, At: t2, Detail: "sent back to in-progress: merge queue batch 20261019-160000 fails with it at unit", Base: "ba5e0000000000000000", Head: head})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	th, err := threads.New(repo, threads.NewOptions{Title: "Merge queue batch 20261019-160000 fails with " + story + " at unit", On: story, Author: "flai serve", Text: "Fix it.", Now: t2})
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// thread is the thread id as it is now.
func thread(t *testing.T, repo *workitem.Repo, id string) *threads.Thread {
	t.Helper()
	th, err := threads.Get(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	return th
}

// S-0400 criterion 4: a story the merge queue sent back that is in review
// again has its send back's thread resolved by flai serve once its branch is
// at a commit other than the one that failed, and is answered; while its
// branch is still at that commit nothing is resolved. A thread on it that
// flai serve did not open, or with another title, is left open, and so is
// the send back's thread of a story still in progress.
func TestASentBackStoryInReviewAtANewCommitHasItsThreadResolved(t *testing.T) {
	repo, _ := gitProject(t, "")
	back, head := readyStory(t, repo, "Back in review", t0.Add(time.Hour), "a.txt", "a\n")
	th := sentBackAt(t, repo, back, head)
	asked, err := threads.New(repo, threads.NewOptions{Title: "Merge queue batch 20261019-160000 fails with " + back + " at unit", On: back, Author: "alex", Text: "Why?", Now: t2})
	if err != nil {
		t.Fatal(err)
	}
	other, err := threads.New(repo, threads.NewOptions{Title: "Pull request #12 of " + back + ": the required check ci failed", On: back, Author: "flai serve", Text: "Fix it.", Now: t2})
	if err != nil {
		t.Fatal(err)
	}
	working, workingHead := readyStory(t, repo, "Still in progress", t0.Add(2*time.Hour), "b.txt", "b\n")
	it, err := repo.Get(working)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Transition(it, workitem.InProgress, "flai", "merge queue batch 20261019-160000 fails with it", t2); err != nil {
		t.Fatal(err)
	}
	inProgress := sentBackAt(t, repo, working, workingHead)
	wt := repo.WorktreePath(working)
	writeFile(t, wt, "b.txt", "b, fixed\n")
	runGit(t, wt, "commit", "-q", "-am", working+"'s fix")
	at := t2.Add(time.Hour)
	d := Deps{Runner: execx.System{}, Now: func() time.Time { return at }}

	resolved, err := ResolveSentBack(d, repo)
	if err != nil || len(resolved) != 0 {
		t.Fatalf("resolved %v, %v; want nothing while %s is at the head that failed", resolved, err, back)
	}
	if !thread(t, repo, th.ID).Open() {
		t.Fatalf("%s resolved while %s is at the head that failed", th.ID, back)
	}

	wt = repo.WorktreePath(back)
	writeFile(t, wt, "a.txt", "a, fixed\n")
	runGit(t, wt, "commit", "-q", "-am", back+"'s fix")
	fixed := runGit(t, wt, "rev-parse", "HEAD")
	resolved, err = ResolveSentBack(d, repo)
	if err != nil || !slices.Equal(resolved, []string{back}) {
		t.Fatalf("resolved %v, %v; want %s", resolved, err, back)
	}
	got := thread(t, repo, th.ID)
	es := got.Entries()
	last := es[len(es)-1]
	if got.Open() || last.Author != "flai serve" || !strings.Contains(last.Text, back+" is back in review at "+fixed[:12]+", a new commit") || last.At != "2026-10-19T18:00:00Z" {
		t.Errorf("%s open %v, last entry %+v; want it resolved by flai serve at the new commit", th.ID, got.Open(), last)
	}
	for _, id := range []string{asked.ID, other.ID, inProgress.ID} {
		if !thread(t, repo, id).Open() {
			t.Errorf("%s was resolved, want it left open", id)
		}
	}

	resolved, err = ResolveSentBack(d, repo)
	if err != nil || len(resolved) != 0 {
		t.Errorf("resolved again %v, %v; want nothing", resolved, err)
	}
}

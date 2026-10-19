package mergequeue

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
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

// t1 is when the batches are built, after every fixture story entered
// review.
var t1 = time.Date(2026, 10, 19, 16, 0, 0, 0, time.UTC)

// batchTiers are a unit tier for the text files the batch changes, and a
// lint tier for them after it.
var batchTiers = []verify.Tier{
	{Name: "unit", Paths: []string{"**/*.txt"}, Command: []string{"unit", "{files}"}},
	{Name: "lint", Paths: []string{"**/*.txt"}, Command: []string{"lint"}},
}

// tierProc runs each tier's command by its first argument, failing those
// in fail, and keeps what ran and the batch's step in the state as each
// ran; cancel, when set, ends the run's context as a tier runs.
type tierProc struct {
	cacheDir string
	fail     map[string]bool
	cancel   context.CancelFunc
	ran      []string
	steps    []string
}

func (p *tierProc) Run(_ context.Context, _ string, argv, _ []string, _, stderr io.Writer) (int, error) {
	p.ran = append(p.ran, strings.Join(argv, " "))
	if st, err := ReadState(p.cacheDir); err == nil && st.Batch != nil {
		p.steps = append(p.steps, st.Batch.Step)
	}
	if p.cancel != nil {
		p.cancel()
	}
	if p.fail[argv[0]] {
		_, _ = io.WriteString(stderr, argv[0]+": boom\n")
		return 1, nil
	}
	return 0, nil
}

// verifier verifies a batch with flai verify's tier selection and runner, as
// CheckoutVerifier does, with batchTiers run by p and real git.
func verifier(p *tierProc) Verifier {
	return func(ctx context.Context, dir, base string, started func(string)) (verify.Report, error) {
		return verify.Tiers(ctx, verify.TierOptions{Dir: dir, Base: base, Tiers: batchTiers, Git: execx.System{}, Started: started, RunOptions: verify.RunOptions{Proc: p}})
	}
}

// deps run real git at the clock at, verifying with p.
func deps(p *tierProc, at time.Time) Deps {
	return Deps{Runner: execx.System{}, Now: func() time.Time { return at }, Verify: verifier(p)}
}

// gitProject is a project in a real git repository, accepting locally, with
// extra added to its manifest, and x.txt committed on main; wip and the
// cache are ignored, as the items need not be committed. It skips under
// -short.
func gitProject(t *testing.T, extra string) (*workitem.Repo, *tierProc) {
	t.Helper()
	if testing.Short() {
		t.Skip("integration: runs real git")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" + extra
	writeFile(t, root, "system-flow.yaml", manifest)
	writeFile(t, root, ".gitignore", ".flai-cache/\nwip/\n")
	writeFile(t, root, "x.txt", "one\ntwo\nthree\n")
	runGit(t, root, "init", "-q", "-b", "main")
	runGit(t, root, "config", "user.email", "t@t")
	runGit(t, root, "config", "user.name", "t")
	runGit(t, root, "config", "commit.gpgsign", "false")
	runGit(t, root, "add", "-A")
	runGit(t, root, "commit", "-q", "-m", "init")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo, &tierProc{cacheDir: repo.CacheDir(), fail: map[string]bool{}}
}

// writeFile writes rel below dir, making its folder.
func writeFile(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runGit runs git in dir and answers what it printed, failing the test when
// git fails.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := (execx.System{}).Run(dir, "git", args...)
	if err != nil {
		t.Fatal(out, err)
	}
	return strings.TrimSpace(out)
}

// readyStory makes a story in review at review whose branch, in its own
// worktree, commits content to rel, with a passing flai verify result at its
// head, enqueued; it answers the story and its head.
func readyStory(t *testing.T, repo *workitem.Repo, title string, review time.Time, rel, content string) (string, string) {
	t.Helper()
	s, _ := story(t, repo, title, []bool{true}, review)
	wt := repo.WorktreePath(s.ID)
	runGit(t, repo.MainRoot, "worktree", "add", "-q", "-b", storygit.Branch(s.ID), wt, "main")
	writeFile(t, wt, rel, content)
	runGit(t, wt, "add", "-A")
	runGit(t, wt, "commit", "-q", "-m", s.ID+"'s change")
	head := runGit(t, wt, "rev-parse", "HEAD")
	if err := verify.SaveReport(repo, verify.Report{Story: s.ID, Commit: head, Passed: true}); err != nil {
		t.Fatal(err)
	}
	return s.ID, head
}

// enqueue enqueues the stories in order.
func enqueue(t *testing.T, repo *workitem.Repo, stories ...string) {
	t.Helper()
	if _, err := Update(repo.CacheDir(), func(s *State) error {
		for _, id := range stories {
			s.Enqueue(id, "alex", t0)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

// runBatch runs RunBatch, failing the test on an error, and answers whether
// it began a batch and the state it left.
func runBatch(t *testing.T, d Deps, repo *workitem.Repo) (bool, State) {
	t.Helper()
	ran, err := RunBatch(context.Background(), d, repo)
	if err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	return ran, st
}

// subjects are the subjects of the commits on ref, newest first.
func subjects(t *testing.T, repo *workitem.Repo, ref string) []string {
	t.Helper()
	return strings.Split(runGit(t, repo.MainRoot, "log", "--format=%s", ref), "\n")
}

// untouched fails the test unless each story's branch is at its head, its
// worktree is checked out there and clean, and its last flai verify result
// is still the one at its head.
func untouched(t *testing.T, repo *workitem.Repo, heads map[string]string) {
	t.Helper()
	for id, head := range heads {
		if got := runGit(t, repo.MainRoot, "rev-parse", storygit.Branch(id)); got != head {
			t.Errorf("%s is at %s, want it left at %s", storygit.Branch(id), got, head)
		}
		wt := repo.WorktreePath(id)
		if got, dirty := runGit(t, wt, "rev-parse", "HEAD"), runGit(t, wt, "status", "--porcelain"); got != head || dirty != "" {
			t.Errorf("%s's worktree at %s with %q, want it left at %s and clean", id, got, dirty, head)
		}
		if rep, ok, err := verify.LastReport(repo, id); err != nil || !ok || rep.Commit != head || !rep.Passed || len(rep.Steps) != 0 {
			t.Errorf("%s's last flai verify result %+v, %v, %v; want it left as it was", id, rep, ok, err)
		}
	}
}

// S-0399 criteria 2 and 3: two clean stories are batched in the order they
// were enqueued, each rebased onto the one before on the batch's own branch
// from main, and the batch is verified once with the tiers its paths select,
// each step written as it begins and the result kept in the state; the
// stories' branches, worktrees, and flai verify results are not touched.
func TestABatchOfTwoCleanStoriesIsBuiltInOrderAndVerifiedOnce(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, headB := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	main := runGit(t, repo.MainRoot, "rev-parse", "main")
	enqueue(t, repo, b, a)

	ran, st := runBatch(t, deps(proc, t1), repo)
	if !ran || st.Batch == nil {
		t.Fatalf("ran %v, state %+v; want a batch", ran, st)
	}
	bt := st.Batch
	if bt.ID != "20261019-160000" || bt.Branch != "merge-queue/20261019-160000" || bt.Worktree != ".flai-cache/worktrees/merge-queue" ||
		bt.Base != main || !slices.Equal(bt.Stories, []string{b, a}) || len(bt.LeftOut) != 0 || !bt.StartedAt.Equal(t1) || bt.Step != StepVerified || bt.Dropped != nil {
		t.Errorf("batch %+v", bt)
	}
	if r := bt.Result; r == nil || !r.Passed || r.StoppedAt != "" || r.Summary != "passed unit, lint" || !r.FinishedAt.Equal(t1) {
		t.Errorf("result %+v, want passed unit and lint at %s", r, t1)
	}
	if want := []string{"verifying unit", "verifying lint"}; !slices.Equal(proc.steps, want) {
		t.Errorf("steps as each tier ran %v, want %v", proc.steps, want)
	}
	if want := []string{"unit a.txt b.txt", "lint"}; !slices.Equal(proc.ran, want) {
		t.Errorf("ran %v, want %v, once", proc.ran, want)
	}
	if got, want := subjects(t, repo, bt.Branch), []string{a + "'s change", b + "'s change", "init"}; !slices.Equal(got, want) {
		t.Errorf("the batch's branch has %v, want %v", got, want)
	}
	// S-0400: each story's verified head and its rebased head in the batch
	wantHeads := []Head{
		{Story: b, Verified: headB, Rebased: runGit(t, repo.MainRoot, "rev-parse", bt.Branch+"~1")},
		{Story: a, Verified: headA, Rebased: runGit(t, repo.MainRoot, "rev-parse", bt.Branch)},
	}
	if !slices.Equal(bt.Heads, wantHeads) {
		t.Errorf("heads %+v, want %+v", bt.Heads, wantHeads)
	}
	wt := filepath.Join(repo.MainRoot, filepath.FromSlash(bt.Worktree))
	if got := runGit(t, wt, "rev-parse", "--abbrev-ref", "HEAD"); got != bt.Branch {
		t.Errorf("the batch's worktree is on %s, want %s", got, bt.Branch)
	}
	if got := runGit(t, repo.MainRoot, "rev-parse", "main"); got != main {
		t.Errorf("main moved to %s", got)
	}
	untouched(t, repo, map[string]string{a: headA, b: headB})
	if len(st.Enqueued) != 2 || len(st.Outcomes) != 0 {
		t.Errorf("enqueued %+v, outcomes %+v; want both still enqueued and no outcome", st.Enqueued, st.Outcomes)
	}

	// a verified batch waits in the state to be landed
	proc.ran = nil
	again, st2 := runBatch(t, deps(proc, t1.Add(time.Hour)), repo)
	if again || st2.Batch == nil || st2.Batch.ID != bt.ID || len(proc.ran) != 0 {
		t.Errorf("ran %v, batch %+v, tiers %v; want the verified batch left waiting", again, st2.Batch, proc.ran)
	}
}

// S-0399 criterion 2: a story whose rebase onto the batch stops is left out,
// with its conflicting paths on the batch, the stop and the paths as its
// outcome, and stays enqueued; the stories after it go on. Its branch is not
// touched.
func TestAStoryWhoseRebaseStopsIsLeftOutAndTheRestGoOn(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "x.txt", "one\ntwo by A\nthree\n")
	b, headB := readyStory(t, repo, "B", t0.Add(2*time.Hour), "x.txt", "one\ntwo by B\nthree\n")
	c, headC := readyStory(t, repo, "C", t0.Add(3*time.Hour), "c.txt", "c\n")
	main := runGit(t, repo.MainRoot, "rev-parse", "main")
	enqueue(t, repo, a, b, c)

	ran, st := runBatch(t, deps(proc, t1), repo)
	if !ran || st.Batch == nil {
		t.Fatalf("ran %v, state %+v; want a batch", ran, st)
	}
	if !slices.Equal(st.Batch.Stories, []string{a, c}) || len(st.Batch.LeftOut) != 1 || st.Batch.LeftOut[0].Story != b ||
		!slices.Equal(st.Batch.LeftOut[0].Paths, []string{"x.txt"}) || st.Batch.Step != StepVerified {
		t.Errorf("batch %+v, want %s and %s with %s left out", st.Batch, a, c, b)
	}
	if got, want := subjects(t, repo, st.Batch.Branch), []string{c + "'s change", a + "'s change", "init"}; !slices.Equal(got, want) {
		t.Errorf("the batch's branch has %v, want %v", got, want)
	}
	// S-0400: a story left out has no heads in the batch
	if _, ok := st.Batch.HeadOf(b); ok || len(st.Batch.Heads) != 2 || st.Batch.Heads[1].Story != c || st.Batch.Heads[1].Rebased != runGit(t, repo.MainRoot, "rev-parse", st.Batch.Branch) {
		t.Errorf("heads %+v, want %s's and %s's alone", st.Batch.Heads, a, c)
	}
	if len(st.Outcomes) != 1 {
		t.Fatalf("outcomes %+v, want %s's alone", st.Outcomes, b)
	}
	o := st.Outcomes[0]
	if o.Story != b || o.Outcome != OutcomeLeftOut || o.Base != main || o.Head != headB || !o.At.Equal(t1) {
		t.Errorf("outcome %+v", o)
	}
	for _, says := range []string{"conflicts in x.txt", storygit.Branch(b), "stays enqueued", st.Batch.ID} {
		if !strings.Contains(o.Detail, says) {
			t.Errorf("detail %q does not say %q", o.Detail, says)
		}
	}
	wt := filepath.Join(repo.MainRoot, filepath.FromSlash(st.Batch.Worktree))
	if storygit.RebaseInProgress(execx.System{}, wt) || runGit(t, wt, "status", "--porcelain") != "" {
		t.Error("the batch's worktree is left with a rebase or changes")
	}
	if _, ok := st.EnqueuedAs(b); !ok {
		t.Errorf("%s is no longer enqueued", b)
	}
	untouched(t, repo, map[string]string{a: headA, b: headB, c: headC})
}

// S-0399 criterion 2: with every story left out there is nothing to verify:
// the batch, its worktree, and its branch are removed, the outcomes kept,
// and the queue does not build the same batch again while neither main nor
// the story's branch moves.
func TestABatchThatLeavesEveryStoryOutIsRemoved(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "x.txt", "one\ntwo by A\nthree\n")
	writeFile(t, repo.MainRoot, "x.txt", "one\ntwo on main\nthree\n")
	runGit(t, repo.MainRoot, "commit", "-q", "-am", "main's change")
	enqueue(t, repo, a)

	ran, st := runBatch(t, deps(proc, t1), repo)
	if !ran || st.Batch != nil || len(proc.ran) != 0 {
		t.Fatalf("ran %v, batch %+v, tiers %v; want no batch left and nothing verified", ran, st.Batch, proc.ran)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Outcome != OutcomeLeftOut || !strings.Contains(st.Outcomes[0].Detail, "conflicts in x.txt") {
		t.Errorf("outcomes %+v, want %s left out on x.txt", st.Outcomes, a)
	}
	if storygit.BranchExists(execx.System{}, repo.MainRoot, "merge-queue/20261019-160000") {
		t.Error("the batch's branch is left")
	}
	if _, err := os.Stat(filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the batch's worktree is left: %v", err)
	}
	if again, _ := runBatch(t, deps(proc, t1.Add(time.Minute)), repo); again {
		t.Error("the same story was batched again at the same main and head")
	}
	untouched(t, repo, map[string]string{a: headA})
}

// S-0399 criterion 3: a tier that fails fails the batch, its result naming
// the tier it stopped at and its findings, and the tiers after it do not
// run.
func TestAFailingTierFailsTheBatch(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	proc.fail["unit"] = true

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || st.Batch.Step != StepFailed {
		t.Fatalf("batch %+v, want it failed", st.Batch)
	}
	r := st.Batch.Result
	if r == nil || r.Passed || r.StoppedAt != "unit" || !strings.HasPrefix(r.Summary, "unit failed") || !r.FinishedAt.Equal(t1) {
		t.Errorf("result %+v, want stopped at unit", r)
	}
	if !slices.Equal(proc.ran, []string{"unit a.txt"}) {
		t.Errorf("ran %v, want unit alone", proc.ran)
	}
}

// S-0405 criterion 3: the batch's verification, as flai serve runs it with
// the manifest of the batch's worktree, runs the tiers marked queue that its
// paths select, those flai verify leaves to the queue while it can batch, so
// that no story lands without them (ADR-0174).
func TestTheBatchRunsTheTiersMarkedQueue(t *testing.T) {
	ran := filepath.Join(t.TempDir(), "ran")
	tier := func(name string) string {
		return "    command: [sh, -c, \"echo " + name + " >> '" + ran + "'\"]\n    paths: [\"**/*.txt\"]\n"
	}
	repo, _ := gitProject(t, "orchestration:\n  permissions:\n    accept_reviews: true\ntests:\n"+
		"  - name: unit\n"+tier("unit")+
		"  - name: smoke\n"+tier("smoke")+"    all_only: true\n    queue: true\n")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)

	d := Deps{Runner: execx.System{}, Now: func() time.Time { return t1 }, Verify: CheckoutVerifier(execx.System{}, "")}
	_, st := runBatch(t, d, repo)
	if st.Batch == nil || st.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, want it verified", st.Batch)
	}
	if r := st.Batch.Result; r == nil || !r.Passed || r.Summary != "passed unit, smoke" {
		t.Errorf("result %+v, want unit and smoke passed", r)
	}
	data, err := os.ReadFile(ran)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Fields(string(data)); !slices.Equal(got, []string{"unit", "smoke"}) {
		t.Errorf("ran %v, want unit and then smoke, the tier marked queue", got)
	}
}

// S-0399 criterion 2: a batch takes at most merge_queue.batch_max stories,
// in the order they were enqueued; the rest stay enqueued.
func TestBatchMaxCapsTheBatch(t *testing.T) {
	repo, proc := gitProject(t, "merge_queue:\n  batch_max: 1\n")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, b, a)

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{b}) || len(st.Enqueued) != 2 {
		t.Errorf("batch %+v, enqueued %+v; want %s alone, both still enqueued", st.Batch, st.Enqueued, b)
	}
}

// S-0399: a batch dropped while it is verified, as flai accept by hand drops
// it, keeps its result and stays dropped; the next run replaces it, removing
// its branch and worktree, and builds the next batch from main.
func TestADroppedBatchKeepsItsResultAndIsReplaced(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	b, _ := readyStory(t, repo, "B", t0.Add(2*time.Hour), "b.txt", "b\n")
	enqueue(t, repo, a, b)
	d := deps(proc, t1)
	inner := d.Verify
	d.Verify = func(ctx context.Context, dir, base string, started func(string)) (verify.Report, error) {
		if _, err := AcceptedByHand(repo.CacheDir(), a, "alex", t1); err != nil {
			t.Fatal(err)
		}
		return inner(ctx, dir, base, started)
	}

	_, st := runBatch(t, d, repo)
	first := st.Batch
	if first == nil || first.Dropped == nil || first.Result == nil || !first.Result.Passed || first.Step != StepVerified {
		t.Fatalf("batch %+v, want it dropped with its result", first)
	}

	ran, st := runBatch(t, deps(proc, t1.Add(time.Hour)), repo)
	next := st.Batch
	if !ran || next == nil || next.ID == first.ID || next.Dropped != nil || next.Step != StepVerified || !slices.Equal(next.Stories, []string{b}) {
		t.Fatalf("ran %v, batch %+v; want a new batch of %s", ran, next, b)
	}
	if storygit.BranchExists(execx.System{}, repo.MainRoot, first.Branch) || !storygit.BranchExists(execx.System{}, repo.MainRoot, next.Branch) {
		t.Errorf("branches: %s left, or %s missing", first.Branch, next.Branch)
	}
}

// S-0399: a context that ends while the batch is verified stops the run
// with no result, the state keeping the step it was at; the next run builds
// the batch again from scratch.
func TestARunStoppedMidwayIsBuiltAgain(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	ctx, cancel := context.WithCancel(context.Background())
	proc.cancel = cancel

	ran, err := RunBatch(ctx, deps(proc, t1), repo)
	if !ran || !errors.Is(err, context.Canceled) {
		t.Fatalf("ran %v, %v; want the run stopped", ran, err)
	}
	st, err := ReadState(repo.CacheDir())
	if err != nil {
		t.Fatal(err)
	}
	if st.Batch == nil || st.Batch.Result != nil || st.Batch.Step != "verifying unit" {
		t.Fatalf("batch %+v, want it at verifying unit with no result", st.Batch)
	}

	proc.cancel = nil
	ran, st = runBatch(t, deps(proc, t1.Add(time.Hour)), repo)
	if !ran || st.Batch == nil || st.Batch.ID != "20261019-170000" || st.Batch.Step != StepVerified || st.Batch.Result == nil {
		t.Errorf("ran %v, batch %+v; want it built again and verified", ran, st.Batch)
	}
}

// RunBatch builds nothing, running no git beyond what Compute reads, while
// the state's batch has a result and is not dropped, under pull-request
// acceptance, and with no eligible story enqueued.
func TestRunBatchDoesNothingWhenThereIsNothingToBuild(t *testing.T) {
	verified := `{"enqueued": [{"story": "S-0001", "by": "alex", "at": "2026-10-19T08:00:00Z"}],
  "batch": {"id": "b1", "stories": ["S-0001"], "step": "verified", "result": {"passed": true, "finished_at": "2026-10-19T09:00:00Z"}}}`
	for _, tc := range []struct {
		name, mode, state string
		enqueue           bool
	}{
		{name: "a verified batch waits", state: verified},
		{name: "pull-request acceptance", mode: "pull-request", enqueue: true},
		{name: "nothing enqueued"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := project(t, tc.mode)
			g := scriptedGit{}
			s, _ := story(t, repo, "Ready", []bool{true}, t0.Add(time.Hour))
			passed(t, repo, g, s.ID, "aaaa")
			if tc.state != "" {
				writeState(t, repo.CacheDir(), tc.state)
			}
			if tc.enqueue {
				enqueue(t, repo, s.ID)
			}
			before, _ := os.ReadFile(StatePath(repo.CacheDir()))
			called := false
			d := Deps{Runner: g, Now: func() time.Time { return t1 }, Verify: func(context.Context, string, string, func(string)) (verify.Report, error) {
				called = true
				return verify.Report{}, nil
			}}
			ran, err := RunBatch(context.Background(), d, repo)
			if ran || err != nil || called {
				t.Errorf("ran %v, %v, verified %v; want nothing done", ran, err, called)
			}
			if after, _ := os.ReadFile(StatePath(repo.CacheDir())); string(after) != string(before) {
				t.Errorf("state written:\n%s\nwas\n%s", after, before)
			}
		})
	}
}

// pick takes the eligible stories in the order they were enqueued, skips
// those not eligible and those left out at the same main commit and head,
// and stops at the limit.
func TestPickTakesTheEnqueuedEligibleStoriesInOrder(t *testing.T) {
	q := Queue{
		Eligible:   []Entry{{Story: "S-0001", Head: "h1"}, {Story: "S-0002", Head: "h2"}, {Story: "S-0003", Head: "h3"}, {Story: "S-0005", Head: "h5"}, {Story: "S-0006", Head: "h6"}},
		Ineligible: []Entry{{Story: "S-0004"}},
		Enqueued:   []Enqueued{{Story: "S-0006"}, {Story: "S-0004"}, {Story: "S-0003"}, {Story: "S-0002"}, {Story: "S-0005"}, {Story: "S-0001"}},
		Outcomes: []Outcome{
			{Story: "S-0003", Outcome: OutcomeLeftOut, Base: "m1", Head: "h3"},
			{Story: "S-0002", Outcome: OutcomeLeftOut, Base: "m0", Head: "h2"},
			{Story: "S-0005", Outcome: OutcomeLeftOut, Base: "m1", Head: "h4"},
		},
	}
	if got, want := ids(pick(q, "m1", 4)), []string{"S-0006", "S-0002", "S-0005", "S-0001"}; !slices.Equal(got, want) {
		t.Errorf("picked %v, want %v", got, want)
	}
	if got, want := ids(pick(q, "m1", 2)), []string{"S-0006", "S-0002"}; !slices.Equal(got, want) {
		t.Errorf("picked %v, want %v", got, want)
	}
}

// A batch's result names the tier it stopped at and its first findings, or
// that its tiers could not run, or the tiers it passed.
func TestTheResultSummarisesTheVerification(t *testing.T) {
	findings := []verify.Finding{{Path: "a.go", Line: 3, Name: "TestA", Message: "boom\nmore"}, {Name: "TestB", Message: "b"}, {Name: "TestC", Message: "c"}, {Name: "TestD", Message: "d"}}
	failed := verify.Report{StoppedAt: "unit", Steps: []verify.Step{{Name: "unit", Tier: true, State: verify.Failed, Findings: findings, Omitted: 2}}}
	if r := result(failed, nil, t1); r.Passed || r.StoppedAt != "unit" || r.Summary != "unit failed: a.go:3 TestA: boom; TestB: b; TestC: c; and 3 more" || !r.FinishedAt.Equal(t1) {
		t.Errorf("failed result %+v", r)
	}
	if r := result(verify.Report{Passed: true}, errors.New("the manifest's tests are not valid"), t1); r.Passed || r.StoppedAt != StoppedAtVerify || !strings.Contains(r.Summary, "not valid") {
		t.Errorf("result of an error %+v", r)
	}
	ok := verify.Report{Passed: true, Steps: []verify.Step{{Name: "unit", State: verify.Passed}, {Name: "lint", State: verify.Skipped}}}
	if r := result(ok, nil, t1); !r.Passed || r.Summary != "passed unit" {
		t.Errorf("passed result %+v", r)
	}
	if r := result(verify.Report{Passed: true, Steps: []verify.Step{}}, nil, t1); r.Summary != "no tier selects the paths the batch changes" {
		t.Errorf("empty result %+v", r)
	}
}

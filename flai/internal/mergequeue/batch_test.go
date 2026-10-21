package mergequeue

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/messages"
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
// CheckoutVerifier does, with batchTiers run by p and real git, each tier
// logged in the folder it is given.
func verifier(p *tierProc) Verifier {
	return func(ctx context.Context, dir, base, logDir string, started func(string)) (verify.Report, error) {
		return verify.Tiers(ctx, verify.TierOptions{Dir: dir, Base: base, Tiers: batchTiers, Git: execx.System{}, Started: started, RunOptions: verify.RunOptions{Proc: p, LogDir: logDir}})
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
// head; it answers the story and its head.
func readyStory(t *testing.T, repo *workitem.Repo, title string, review time.Time, rel, content string) (string, string) {
	t.Helper()
	return readyStoryChanging(t, repo, title, review, func(_, wt string) { writeFile(t, wt, rel, content) })
}

// readyStoryChanging is readyStory whose branch commits what change writes
// in its worktree wt for the story id.
func readyStoryChanging(t *testing.T, repo *workitem.Repo, title string, review time.Time, change func(id, wt string)) (string, string) {
	t.Helper()
	s, _ := story(t, repo, title, []bool{true}, review)
	wt := repo.WorktreePath(s.ID)
	runGit(t, repo.MainRoot, "worktree", "add", "-q", "-b", storygit.Branch(s.ID), wt, "main")
	change(s.ID, wt)
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
// touched. S-0401 criterion 3: the stop is on a path outside claims.shared
// that the story ahead changed, so flai messages the story from the story
// ahead about it in the two stories' conversation (ADR-0175).
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
	for _, says := range []string{"conflicts in x.txt", storygit.Branch(b), "stays enqueued", st.Batch.ID, a} {
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
	if it, err := repo.Get(b); err != nil || it.Status != workitem.Review {
		t.Errorf("%s is %+v, %v; want it left in review", b, it, err)
	}
	convs, err := messages.For(repo, b)
	if err != nil || len(convs) != 1 {
		t.Fatalf("conversations of %s: %+v, %v; want one", b, convs, err)
	}
	conv := convs[0]
	e := conv.Entries()
	if conv.From != a || conv.To != b || !slices.Equal(conv.About, []string{"x.txt"}) || conv.Awaiting() != b || len(e) != 1 ||
		e[0].Author != storygit.ConflictAuthor || e[0].Story != a {
		t.Fatalf("conversation %+v, entries %+v; want one from %s to %s about x.txt", conv, e, a, b)
	}
	for _, says := range []string{st.Batch.ID, "- `x.txt`", "flai stream sync " + b, "`flai message escalate`"} {
		if !strings.Contains(e[0].Text, says) {
			t.Errorf("message does not say %q:\n%s", says, e[0].Text)
		}
	}
	untouched(t, repo, map[string]string{a: headA, b: headB, c: headC})
}

// S-0401 criterion 3: a stop on a path outside claims.shared that no story
// ahead in the batch changed is a conflict with the main branch itself: the
// story is sent back to in-progress, by flai, with the paths as the reason,
// and leaves the queue; the stories after it go on (ADR-0175).
func TestABatchSendsBackAStoryThatConflictsWithTheMainBranch(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "x.txt", "one\ntwo by A\nthree\n")
	c, headC := readyStory(t, repo, "C", t0.Add(2*time.Hour), "c.txt", "c\n")
	writeFile(t, repo.MainRoot, "x.txt", "one\ntwo on main\nthree\n")
	runGit(t, repo.MainRoot, "commit", "-q", "-am", "main's change")
	main := runGit(t, repo.MainRoot, "rev-parse", "main")
	enqueue(t, repo, a, c)

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{c}) || len(st.Batch.LeftOut) != 1 || st.Batch.LeftOut[0].Story != a ||
		!slices.Equal(st.Batch.LeftOut[0].Paths, []string{"x.txt"}) || st.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, want %s alone with %s taken out on x.txt", st.Batch, c, a)
	}
	if len(st.Outcomes) != 1 {
		t.Fatalf("outcomes %+v, want %s's alone", st.Outcomes, a)
	}
	o := st.Outcomes[0]
	if o.Story != a || o.Outcome != OutcomeSentBack || o.Base != main || o.Head != headA || !o.At.Equal(t1) {
		t.Errorf("outcome %+v", o)
	}
	reason := "the merge queue's batch " + st.Batch.ID + " could not rebase it onto main: conflicts in x.txt; sync it with flai stream sync " + a
	if !strings.Contains(o.Detail, reason) {
		t.Errorf("detail %q does not say %q", o.Detail, reason)
	}
	if _, ok := st.EnqueuedAs(a); ok {
		t.Errorf("%s is still enqueued", a)
	}
	if _, ok := st.EnqueuedAs(c); !ok {
		t.Errorf("%s is no longer enqueued", c)
	}
	it, err := repo.Get(a)
	if err != nil {
		t.Fatal(err)
	}
	if last := it.Transitions[len(it.Transitions)-1]; it.Status != workitem.InProgress || last.To != workitem.InProgress || last.By != sendBackBy {
		t.Errorf("%s is %s, last moved %+v; want it sent back to in-progress by flai", a, it.Status, last)
	}
	if !strings.Contains(it.Body, "moved to in-progress: "+reason) {
		t.Errorf("%s's notes do not give the reason %q:\n%s", a, reason, it.Body)
	}
	if convs, err := messages.For(repo, a); err != nil || len(convs) != 0 {
		t.Errorf("conversations of %s: %+v, %v; want none", a, convs, err)
	}
	untouched(t, repo, map[string]string{a: headA, c: headC})
}

// sharedFlaiCLI is a manifest that shares the flai CLI's design.
const sharedFlaiCLI = "claims:\n  shared:\n    - design/system/flai-cli.md\n"

// cliDesign is the shared design document on main, with one line in one
// section.
const cliDesign = "# flai CLI\n\n## Commands\n\nline one\n"

// cliDesignOnMain commits cliDesign on main.
func cliDesignOnMain(t *testing.T, repo *workitem.Repo) {
	t.Helper()
	writeFile(t, repo.MainRoot, "design/system/flai-cli.md", cliDesign)
	runGit(t, repo.MainRoot, "add", "-A")
	runGit(t, repo.MainRoot, "commit", "-q", "-m", "docs: the CLI's design")
}

// S-0401 criterion 2: two stories that both add a section at the end of a
// markdown file inside claims.shared stop the second's rebase onto the batch
// on it; the batch merges it by sections, both sections kept, the batch's
// first, and both stories stay in the batch (ADR-0175).
func TestTheBatchMergesASharedMarkdownFileBySections(t *testing.T) {
	repo, proc := gitProject(t, sharedFlaiCLI)
	cliDesignOnMain(t, repo)
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "design/system/flai-cli.md", cliDesign+"\n## By A\n\nA's section\n")
	b, headB := readyStory(t, repo, "B", t0.Add(2*time.Hour), "design/system/flai-cli.md", cliDesign+"\n## By B\n\nB's section\n")
	enqueue(t, repo, a, b)
	var logged bytes.Buffer
	d := deps(proc, t1)
	d.Log = slog.New(slog.NewTextHandler(&logged, nil))

	_, st := runBatch(t, d, repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{a, b}) || len(st.Batch.LeftOut) != 0 || len(st.Outcomes) != 0 || st.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, outcomes %+v; want both stories kept", st.Batch, st.Outcomes)
	}
	want := cliDesign + "\n## By A\n\nA's section\n\n## By B\n\nB's section\n"
	if got := runGit(t, repo.MainRoot, "show", st.Batch.Branch+":design/system/flai-cli.md") + "\n"; got != want {
		t.Errorf("the shared design on the batch:\n%s\nwant:\n%s", got, want)
	}
	if !strings.Contains(logged.String(), "merged=[design/system/flai-cli.md]") {
		t.Errorf("the log does not say the file was merged:\n%s", logged.String())
	}
	untouched(t, repo, map[string]string{a: headA, b: headB})
}

// S-0401 criterion 2: a stop on a markdown file inside claims.shared that
// the section merge cannot settle, the same line changed two ways, here by
// main, leaves the story waiting with the path: it stays enqueued and in
// review, is never sent back, and is not picked again at the same main and
// head; the stories after it go on (ADR-0175).
func TestABatchLeavesAStoryWaitingOnASharedFileItCannotSettle(t *testing.T) {
	repo, proc := gitProject(t, sharedFlaiCLI)
	cliDesignOnMain(t, repo)
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "design/system/flai-cli.md", strings.Replace(cliDesign, "line one", "line by A", 1))
	c, headC := readyStory(t, repo, "C", t0.Add(2*time.Hour), "c.txt", "c\n")
	writeFile(t, repo.MainRoot, "design/system/flai-cli.md", strings.Replace(cliDesign, "line one", "line on main", 1))
	runGit(t, repo.MainRoot, "commit", "-q", "-am", "docs: main's line")
	main := runGit(t, repo.MainRoot, "rev-parse", "main")
	enqueue(t, repo, a, c)

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{c}) || len(st.Batch.LeftOut) != 1 || st.Batch.LeftOut[0].Story != a ||
		!slices.Equal(st.Batch.LeftOut[0].Paths, []string{"design/system/flai-cli.md"}) {
		t.Fatalf("batch %+v, want %s alone with %s waiting on the shared design", st.Batch, c, a)
	}
	if len(st.Outcomes) != 1 {
		t.Fatalf("outcomes %+v, want %s's alone", st.Outcomes, a)
	}
	o := st.Outcomes[0]
	if o.Story != a || o.Outcome != OutcomeWaiting || o.Base != main || o.Head != headA || !o.At.Equal(t1) {
		t.Errorf("outcome %+v", o)
	}
	for _, says := range []string{"design/system/flai-cli.md", "stays enqueued", st.Batch.ID, storygit.Branch(a)} {
		if !strings.Contains(o.Detail, says) {
			t.Errorf("detail %q does not say %q", o.Detail, says)
		}
	}
	if _, ok := st.EnqueuedAs(a); !ok {
		t.Errorf("%s is no longer enqueued", a)
	}
	if it, err := repo.Get(a); err != nil || it.Status != workitem.Review {
		t.Errorf("%s is %+v, %v; want it left in review", a, it, err)
	}
	q, err := Compute(execx.System{}, repo)
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(pick(q, main, 10)); !slices.Equal(got, []string{c}) {
		t.Errorf("picked %v at the same main and head, want %s alone", got, c)
	}
	untouched(t, repo, map[string]string{a: headA, c: headC})
}

// S-0399 criterion 2: with every story left out there is nothing to verify:
// the batch, its worktree, and its branch are removed, the outcomes kept,
// and the queue does not build the same batch again while neither main nor
// the story's branch moves. The story conflicts with main on a shared path
// that is not markdown, which is never merged, so it waits (ADR-0175).
func TestABatchThatLeavesEveryStoryOutIsRemoved(t *testing.T) {
	repo, proc := gitProject(t, "claims:\n  shared:\n    - x.txt\n")
	a, headA := readyStory(t, repo, "A", t0.Add(time.Hour), "x.txt", "one\ntwo by A\nthree\n")
	writeFile(t, repo.MainRoot, "x.txt", "one\ntwo on main\nthree\n")
	runGit(t, repo.MainRoot, "commit", "-q", "-am", "main's change")
	enqueue(t, repo, a)

	ran, st := runBatch(t, deps(proc, t1), repo)
	if !ran || st.Batch != nil || len(proc.ran) != 0 {
		t.Fatalf("ran %v, batch %+v, tiers %v; want no batch left and nothing verified", ran, st.Batch, proc.ran)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Outcome != OutcomeWaiting || !strings.Contains(st.Outcomes[0].Detail, "x.txt") {
		t.Errorf("outcomes %+v, want %s waiting on x.txt", st.Outcomes, a)
	}
	if _, ok := st.EnqueuedAs(a); !ok {
		t.Errorf("%s is no longer enqueued", a)
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

// adrIndexPath is the ADR index the stories add rows to.
const adrIndexPath = "design/adrs/README.md"

// adrIndex is an ADR index whose table holds the rows of the ADRs numbered
// in ns.
func adrIndex(ns ...int) string {
	var rows strings.Builder
	for _, n := range ns {
		fmt.Fprintf(&rows, "| [%04d](%04d-decision-%d.md) | Decision %d | accepted |\n", n, n, n, n)
	}
	return "# Architecture Decision Records\n\nOne decision per file.\n\n| ADR | Title | Status |\n|-----|-------|--------|\n" + rows.String()
}

// sharedRecords commits on main issue I-0001, the issue summary written from
// it, and an ADR index of ADRs 1 and 2, the files each story then changes.
func sharedRecords(t *testing.T, repo *workitem.Repo) {
	t.Helper()
	if _, err := issues.New(repo, issues.NewOptions{Title: "Shared friction", Class: "efficiency", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := issues.WriteSummary(repo, t0); err != nil {
		t.Fatal(err)
	}
	writeFile(t, repo.MainRoot, adrIndexPath, adrIndex(1, 2))
	runGit(t, repo.MainRoot, "add", "-A")
	runGit(t, repo.MainRoot, "commit", "-q", "-m", "docs: issues and ADRs")
}

// recording is a story's change that bumps I-0001 for the story at at,
// writes the issue summary again, and adds ADR adr to the index, as a story
// that records an issue and an ADR commits them.
func recording(t *testing.T, repo *workitem.Repo, at time.Time, adr int) func(id, wt string) {
	return func(id, _ string) {
		wt := issues.RepoFor(repo, id)
		is, err := issues.Get(wt, "I-0001")
		if err != nil {
			t.Fatal(err)
		}
		if err := issues.Bump(is, id, "", "Again in "+id+".", at); err != nil {
			t.Fatal(err)
		}
		if _, err := issues.WriteSummary(wt, at); err != nil {
			t.Fatal(err)
		}
		writeFile(t, wt.Root, adrIndexPath, adrIndex(1, 2, adr))
	}
}

// S-0401 criterion 1: two stories that both bump one issue, write the issue
// summary again, and add a row to the ADR index stop the second's rebase
// onto the batch on all three; the batch merges the issue and the index,
// writes the summary again from the merge, and continues, as flai stream
// sync does, so both stories stay in the batch (ADR-0175).
func TestTheBatchMergesTheIssuesTheSummaryAndTheADRIndexItStopsOn(t *testing.T) {
	repo, proc := gitProject(t, "")
	sharedRecords(t, repo)
	a, headA := readyStoryChanging(t, repo, "A", t0.Add(time.Hour), recording(t, repo, t0.Add(time.Hour), 3))
	b, headB := readyStoryChanging(t, repo, "B", t0.Add(2*time.Hour), recording(t, repo, t0.Add(2*time.Hour), 4))
	enqueue(t, repo, a, b)
	var logged bytes.Buffer
	d := deps(proc, t1)
	d.Log = slog.New(slog.NewTextHandler(&logged, nil))

	_, st := runBatch(t, d, repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{a, b}) || len(st.Batch.LeftOut) != 0 || len(st.Outcomes) != 0 || st.Batch.Step != StepVerified {
		t.Fatalf("batch %+v, outcomes %+v; want both stories kept", st.Batch, st.Outcomes)
	}
	if got, want := subjects(t, repo, st.Batch.Branch), []string{b + "'s change", a + "'s change", "docs: issues and ADRs", "init"}; !slices.Equal(got, want) {
		t.Errorf("the batch's branch has %v, want %v", got, want)
	}
	show := func(rel string) string { return runGit(t, repo.MainRoot, "show", st.Batch.Branch+":"+rel) + "\n" }
	if got, want := show(adrIndexPath), adrIndex(1, 2, 3, 4); got != want {
		t.Errorf("the ADR index on the batch:\n%s\nwant:\n%s", got, want)
	}
	issue := show("design/issues/I-0001-shared-friction.md")
	for _, says := range []string{"count: 3", "Again in " + a + ".", "Again in " + b + "."} {
		if !strings.Contains(issue, says) {
			t.Errorf("I-0001 on the batch does not say %q:\n%s", says, issue)
		}
	}
	summary := show("design/issues/summary.md")
	if !strings.Contains(summary, "[I-0001]") || strings.Contains(summary, "<<<<<<<") || !strings.Contains(summary, "| 3 |") {
		t.Errorf("the summary on the batch is not written from the merged issue:\n%s", summary)
	}
	wt := filepath.Join(repo.MainRoot, filepath.FromSlash(st.Batch.Worktree))
	if storygit.RebaseInProgress(execx.System{}, wt) || runGit(t, wt, "status", "--porcelain") != "" {
		t.Error("the batch's worktree is left with a rebase or changes")
	}
	for _, says := range []string{"component=mergequeue", "story=" + b, adrIndexPath, "design/issues/I-0001-shared-friction.md", "regenerated=[design/issues/summary.md]"} {
		if !strings.Contains(logged.String(), says) {
			t.Errorf("the log does not say %q:\n%s", says, logged.String())
		}
	}
	untouched(t, repo, map[string]string{a: headA, b: headB})
}

// S-0401 criterion 1: a stop on a path the merge steps do not settle beside
// one they do is still aborted and the story left out, with the stop's
// paths, and the batch's worktree is left clean on its branch.
func TestTheBatchLeavesOutAStoryWhoseStopTheMergeStepsDoNotSettle(t *testing.T) {
	repo, proc := gitProject(t, "")
	sharedRecords(t, repo)
	both := func(adr int, line string) func(id, wt string) {
		record := recording(t, repo, t0.Add(time.Duration(adr)*time.Hour), adr)
		return func(id, wt string) {
			record(id, wt)
			writeFile(t, wt, "x.txt", "one\n"+line+"\nthree\n")
		}
	}
	a, headA := readyStoryChanging(t, repo, "A", t0.Add(time.Hour), both(3, "two by A"))
	b, headB := readyStoryChanging(t, repo, "B", t0.Add(2*time.Hour), both(4, "two by B"))
	enqueue(t, repo, a, b)

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || !slices.Equal(st.Batch.Stories, []string{a}) || len(st.Batch.LeftOut) != 1 || st.Batch.LeftOut[0].Story != b {
		t.Fatalf("batch %+v, want %s alone with %s left out", st.Batch, a, b)
	}
	if got := st.Batch.LeftOut[0].Paths; !slices.Contains(got, "x.txt") {
		t.Errorf("left out on %v, want x.txt among them", got)
	}
	if len(st.Outcomes) != 1 || st.Outcomes[0].Outcome != OutcomeLeftOut || !strings.Contains(st.Outcomes[0].Detail, "x.txt") {
		t.Errorf("outcomes %+v, want %s left out on x.txt", st.Outcomes, b)
	}
	if got, want := subjects(t, repo, st.Batch.Branch), []string{a + "'s change", "docs: issues and ADRs", "init"}; !slices.Equal(got, want) {
		t.Errorf("the batch's branch has %v, want %v", got, want)
	}
	wt := filepath.Join(repo.MainRoot, filepath.FromSlash(st.Batch.Worktree))
	if storygit.RebaseInProgress(execx.System{}, wt) || runGit(t, wt, "status", "--porcelain") != "" {
		t.Error("the batch's worktree is left with a rebase or changes")
	}
	if got := runGit(t, wt, "rev-parse", "--abbrev-ref", "HEAD"); got != st.Batch.Branch {
		t.Errorf("the batch's worktree is on %s, want %s", got, st.Batch.Branch)
	}
	untouched(t, repo, map[string]string{a: headA, b: headB})
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

// S-0457 criterion 1 (I-0172): the tiers of a batch that fails keep their
// whole output in the batch's own folder beside the queue's state, which its
// result names, after the batch's worktree is cleaned.
func TestAFailingBatchKeepsItsTiersOutputAfterItsWorktreeIsCleaned(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	proc.fail["unit"] = true
	d := deps(proc, t1)

	_, st := runBatch(t, d, repo)
	if st.Batch == nil || st.Batch.Result == nil || st.Batch.Result.Passed {
		t.Fatalf("batch %+v, want it failed", st.Batch)
	}
	want := ".flai-cache/merge-queue/batches/" + st.Batch.ID
	if got := st.Batch.Result.Log; got != want {
		t.Errorf("result's log %q, want %q", got, want)
	}
	b := &build{d: withDefaults(d), repo: repo, path: filepath.Join(repo.CacheDir(), "worktrees", WorktreeName)}
	if err := b.clean(st.Batch); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(b.path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the batch's worktree is still there: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(repo.MainRoot, filepath.FromSlash(want), "unit.log"))
	if err != nil || !strings.Contains(string(data), "unit: boom") {
		t.Errorf("unit's log %q, %v; want the tier's whole output kept", data, err)
	}
	if _, err := os.Stat(filepath.Join(BatchLogs(repo.CacheDir(), st.Batch.ID), "lint.log")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("lint, which did not run, has a log: %v", err)
	}
}

// S-0457: a new batch keeps the newest logsKept batches' folders of tier
// output, its own among them, by batch ID, and removes the older ones; it
// touches no file there.
func TestANewBatchKeepsTheNewestBatchesLogs(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	var old []string
	for i := range 25 {
		id := t1.Add(-time.Duration(25-i) * time.Hour).Format("20060102-150405")
		old = append(old, id)
		writeFile(t, BatchLogs(repo.CacheDir(), id), "unit.log", "old\n")
	}
	stray := filepath.Join(filepath.Dir(BatchLogs(repo.CacheDir(), "x")), "notes.txt")
	writeFile(t, filepath.Dir(stray), filepath.Base(stray), "kept\n")

	_, st := runBatch(t, deps(proc, t1), repo)
	if st.Batch == nil || st.Batch.Result == nil || st.Batch.Result.Log == "" {
		t.Fatalf("batch %+v, want it verified and logged", st.Batch)
	}
	entries, err := os.ReadDir(filepath.Dir(BatchLogs(repo.CacheDir(), st.Batch.ID)))
	if err != nil {
		t.Fatal(err)
	}
	var kept []string
	for _, e := range entries {
		if e.IsDir() {
			kept = append(kept, e.Name())
		}
	}
	if want := append(slices.Clone(old[25-(logsKept-1):]), st.Batch.ID); !slices.Equal(kept, want) {
		t.Errorf("kept %v, want the newest %d, %v", kept, logsKept, want)
	}
	if _, err := os.Stat(stray); err != nil {
		t.Errorf("a file there was removed: %v", err)
	}
}

// S-0457: a folder of batch logs that cannot be pruned is logged at warn and
// stops no batch.
func TestABatchWhoseLogsCannotBePrunedGoesOn(t *testing.T) {
	repo, proc := gitProject(t, "")
	a, _ := readyStory(t, repo, "A", t0.Add(time.Hour), "a.txt", "a\n")
	enqueue(t, repo, a)
	// a file where the batches' folder should be: nothing can be pruned or
	// logged there
	writeFile(t, filepath.Dir(StatePath(repo.CacheDir())), "batches", "not a folder\n")
	var logged bytes.Buffer
	d := deps(proc, t1)
	d.Log = slog.New(slog.NewTextHandler(&logged, nil))

	_, st := runBatch(t, d, repo)
	if st.Batch == nil || st.Batch.Step != StepVerified || st.Batch.Result == nil || st.Batch.Result.Log != "" {
		t.Fatalf("batch %+v, want it verified with no log", st.Batch)
	}
	if !strings.Contains(logged.String(), "level=WARN") || !strings.Contains(logged.String(), "not pruned") {
		t.Errorf("logged %q, want a warning that the logs were not pruned", logged.String())
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
	d.Verify = func(ctx context.Context, dir, base, logDir string, started func(string)) (verify.Report, error) {
		if _, err := AcceptedByHand(repo.CacheDir(), a, "alex", t1); err != nil {
			t.Fatal(err)
		}
		return inner(ctx, dir, base, logDir, started)
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
			d := Deps{Runner: g, Now: func() time.Time { return t1 }, Verify: func(context.Context, string, string, string, func(string)) (verify.Report, error) {
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

// ADR-0175: pick skips a story left waiting at the same main commit and
// head, as it skips one left out, and takes it again once either moves; a
// story sent back is not skipped for it.
func TestPickSkipsAStoryWaitingAtTheSameMainAndHead(t *testing.T) {
	q := Queue{
		Eligible: []Entry{{Story: "S-0001", Head: "h1"}, {Story: "S-0002", Head: "h2"}, {Story: "S-0003", Head: "h3"}, {Story: "S-0004", Head: "h4"}},
		Enqueued: []Enqueued{{Story: "S-0001"}, {Story: "S-0002"}, {Story: "S-0003"}, {Story: "S-0004"}},
		Outcomes: []Outcome{
			{Story: "S-0001", Outcome: OutcomeWaiting, Base: "m1", Head: "h1"},
			{Story: "S-0002", Outcome: OutcomeWaiting, Base: "m0", Head: "h2"},
			{Story: "S-0003", Outcome: OutcomeWaiting, Base: "m1", Head: "h0"},
			{Story: "S-0004", Outcome: OutcomeSentBack, Base: "m1", Head: "h4"},
		},
	}
	if got, want := ids(pick(q, "m1", 10)), []string{"S-0002", "S-0003", "S-0004"}; !slices.Equal(got, want) {
		t.Errorf("picked %v, want %v", got, want)
	}
}

// ADR-0175: a stop the batch cannot settle waits when every path flai does
// not resolve is inside claims.shared, is left out when a story ahead
// changed a path outside it, naming each such story with the paths it
// shares, and is sent back when none did.
func TestClassifyAStopByWhereItsPathsConflict(t *testing.T) {
	shared := func(p string) bool { return strings.HasPrefix(p, "design/") }
	before := []ahead{{Story: "S-0001", Paths: []string{"a.go", "design/x.md"}}, {Story: "S-0002", Paths: []string{"b.go"}}, {Story: "S-0003", Paths: []string{"a.go", "b.go"}}}
	for _, tc := range []struct {
		name  string
		paths []string
		want  stop
	}{
		{name: "shared alone", paths: []string{"design/x.md", "design/y.md"}, want: stop{Outcome: OutcomeWaiting}},
		{name: "none", paths: nil, want: stop{Outcome: OutcomeWaiting}},
		{name: "outside, changed ahead", paths: []string{"design/x.md", "a.go", "c.go"}, want: stop{Outcome: OutcomeLeftOut, Outside: []string{"a.go", "c.go"},
			Ahead: []ahead{{Story: "S-0001", Paths: []string{"a.go"}}, {Story: "S-0003", Paths: []string{"a.go"}}}}},
		{name: "outside, changed by none ahead", paths: []string{"design/x.md", "c.go"}, want: stop{Outcome: OutcomeSentBack, Outside: []string{"c.go"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := classify(tc.paths, shared, before); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("classify %v: %+v, want %+v", tc.paths, got, tc.want)
			}
		})
	}
	if got := classify([]string{"c.go"}, shared, nil); got.Outcome != OutcomeSentBack {
		t.Errorf("with no story ahead: %+v, want sent back", got)
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
	if r := result(failed, nil, t1); r.Findings != "a.go:3 TestA: boom\n    more\nTestB: b\nTestC: c\nTestD: d\n… 2 more findings left out\n" {
		t.Errorf("failed result's findings:\n%s", r.Findings)
	}
	line := strings.Repeat("x", 99)
	long := verify.Report{StoppedAt: "smoke", Steps: []verify.Step{{Name: "smoke", Tier: true, State: verify.Failed, Findings: []verify.Finding{{Name: "smoke", Message: strings.TrimSuffix(strings.Repeat(line+"\n", 100), "\n")}}}}}
	r := result(long, nil, t1)
	if n := len([]rune(r.Findings)); n > findingsMax+100 || !strings.HasPrefix(r.Findings, "smoke: "+line+"\n    "+line+"\n") ||
		!strings.HasSuffix(r.Findings, "    "+line+"\n… 2496 more characters left out; the tier's whole output is in its log\n") {
		t.Errorf("findings of %d characters, want them cut at a line within %d and marked:\n%s", n, findingsMax, r.Findings)
	}
	if r := result(verify.Report{Passed: true}, errors.New("the manifest's tests are not valid"), t1); r.Passed || r.StoppedAt != StoppedAtVerify || !strings.Contains(r.Summary, "not valid") || r.Findings != "" {
		t.Errorf("result of an error %+v", r)
	}
	ok := verify.Report{Passed: true, Steps: []verify.Step{{Name: "unit", State: verify.Passed}, {Name: "lint", State: verify.Skipped}}}
	if r := result(ok, nil, t1); !r.Passed || r.Summary != "passed unit" || r.Findings != "" {
		t.Errorf("passed result %+v", r)
	}
	if r := result(verify.Report{Passed: true, Steps: []verify.Step{}}, nil, t1); r.Summary != "no tier selects the paths the batch changes" {
		t.Errorf("empty result %+v", r)
	}
}

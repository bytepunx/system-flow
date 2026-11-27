package serve

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// resolveGit answers the git a rebase and resolve job runs itself: the
// worktree's head from its HEAD file, the main checkout's branch, main, a
// rebase waiting while the worktree's rebase-merge folder is there, the
// paths in conflict, and git rebase --abort, which removes that folder. It
// records each command.
type resolveGit struct {
	mu        sync.Mutex
	conflicts []string
	ran       []string
}

func (g *resolveGit) Run(dir, name string, args ...string) (string, error) {
	cmd := strings.Join(args, " ")
	g.mu.Lock()
	g.ran = append(g.ran, cmd)
	conflicts := slices.Clone(g.conflicts)
	g.mu.Unlock()
	switch {
	case name != "git":
	case cmd == "rev-parse HEAD":
		data, err := os.ReadFile(filepath.Join(dir, "HEAD"))
		return string(data), err
	case cmd == "rev-parse --abbrev-ref HEAD":
		return "main", nil
	case strings.HasPrefix(cmd, "rev-parse --git-path "):
		return filepath.Join(dir, args[2]), nil
	case cmd == "diff --name-only --diff-filter=U":
		return strings.Join(conflicts, "\n"), nil
	case cmd == "rebase --abort":
		return "", os.RemoveAll(filepath.Join(dir, "rebase-merge"))
	case cmd == "symbolic-ref --quiet --short HEAD":
		return storygit.Branch(filepath.Base(dir)), nil
	case cmd == "status --porcelain":
		return "", nil
	}
	return "", fmt.Errorf("resolveGit does not answer %s %s", name, cmd)
}

func (g *resolveGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (*resolveGit) LookPath(name string) (string, error) { return name, nil }

// aborted says whether the job ran git rebase --abort.
func (g *resolveGit) aborted() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return slices.Contains(g.ran, "rebase --abort")
}

// newResolveLab is a verify lab whose story is in review, with the stub
// close-out script in its worktree, and whose job's processes a fake runner
// records rather than runs, with this flai at hostFlai and resolveGit as its
// git.
func newResolveLab(t *testing.T) (*verifyLab, *fakeCloseOut, *resolveGit) {
	t.Helper()
	lab := newVerifyLab(t, true)
	lab.release()
	lab.toReview()
	fake, git := &fakeCloseOut{}, &resolveGit{}
	lab.pass.run = fake.run
	lab.pass.exe = func() (string, error) { return hostFlai, nil }
	lab.pass.o.Git = git
	return lab, fake, git
}

// toReview moves the lab's story to review, with a criterion and a task, as
// its agent would.
func (lab *verifyLab) toReview() {
	lab.t.Helper()
	st, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: lab.story, Owner: "alex", Now: time.Now()}); err != nil {
		lab.t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Review} {
		lab.moveStory(to)
	}
}

// moveStory moves the lab's story to to, as its agent.
func (lab *verifyLab) moveStory(to string) {
	lab.t.Helper()
	st, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(st, to, "agent-"+lab.story, "the test moves it", time.Now()); err != nil {
		lab.t.Fatal(err)
	}
}

// askResolve asks for the story's rebase and resolve, as flai accept
// --resolve does.
func (lab *verifyLab) askResolve() closeout.Answer {
	lab.t.Helper()
	return lab.ask(closeout.Request{Kind: closeout.KindResolve, By: "alex"})
}

// resolveJournal is the resolve job's entries journaled so far.
func (lab *verifyLab) resolveJournal() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	var out []hostapi.Entry
	for _, e := range lab.entries {
		if e.Action == closeOutAction && e.Method == resolveMethod {
			out = append(out, e)
		}
	}
	return out
}

// endedJob is the story's job once it ended.
func (lab *verifyLab) endedJob() closeout.Job {
	lab.t.Helper()
	lab.awaitVerified(1)
	job, ok, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story)
	if err != nil || !ok || job.Running() {
		lab.t.Fatalf("the ended job: %+v %v %v", job, ok, err)
	}
	return job
}

// status is the story's status and who last moved it.
func (lab *verifyLab) status() (string, string) {
	lab.t.Helper()
	st, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return st.Status, st.Transitions[len(st.Transitions)-1].By
}

// threadsOn are the threads on the lab's story.
func (lab *verifyLab) threadsOn() []*threads.Thread {
	lab.t.Helper()
	all, err := threads.For(lab.repo, lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return all
}

// logged is the job's log.
func (lab *verifyLab) logged() string {
	lab.t.Helper()
	out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if err != nil {
		lab.t.Fatal(err)
	}
	return string(out)
}

// syncArgv is this flai's flai stream sync of story.
func syncArgv(story string) string { return hostFlai + " stream sync " + story }

// S-0804 criteria 1, 2, 5, 6: a rebase and resolve of a story in review
// starts a job, even when the verify record holds a pass, which rebases with
// this flai's flai stream sync in the worktree and then runs the close-out
// with its own message and FLAI_VERIFY_PROGRESS=1, all of it in the job's log
// between the job's progress lines. A clean rebase and a passing close-out
// leave the story in review with nothing sent back, journaled apart from the
// close-outs flai stats counts. Before S-0804 a request's kind was not read
// and it ran a close-out, or was answered from the record.
func TestARebaseAndResolveThatRebasesCleanlyAndPassesLeavesTheStoryInReview(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	lab.pass.recorded = func(*workitem.Repo, string, string, bool) (verify.Report, bool, error) {
		return verify.Report{Story: lab.story, Commit: "head-a", Passed: true, RanAt: time.Now()}, true, nil
	}
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			_, _ = io.WriteString(out, "rebased story/"+lab.story+" onto main\n")
			_ = os.WriteFile(filepath.Join(lab.worktree, "HEAD"), []byte("head-r\n"), 0o644)
			return 0
		}
		_, _ = io.WriteString(out, "close-out: "+lab.story+" passed every step; ready to move to review\n")
		return 0
	}
	lab.start()
	a := lab.askResolve()
	if !a.Started || a.Recorded || a.Refused() || !strings.HasPrefix(a.Job, lab.story+"-") {
		t.Fatalf("the answer: %+v", a)
	}
	job := lab.endedJob()
	if job.Kind != closeout.KindResolve || job.Outcome != closeout.OutcomeStayed || job.Thread != "" || job.Verdict != closeout.VerdictPassed ||
		job.Commit != "head-r" || len(job.Conflicts) != 0 || !strings.Contains(job.Reason, "stays in review") {
		t.Errorf("the ended job: %+v", job)
	}
	if status, by := lab.status(); status != workitem.Review || by != "agent-"+lab.story {
		t.Errorf("the story is %s, last moved by %s; want it left in review", status, by)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("nothing sent back, yet threads: %+v", th)
	}
	calls := fake.ran()
	want := []string{syncArgv(lab.story), "sh " + closeOutScript + " " + lab.story + " -m " + resolveMessage(lab.story)}
	if got := argvs(calls); len(got) != 2 || strings.Join(got[0], " ") != want[0] || strings.Join(got[1], " ") != want[1] {
		t.Fatalf("the job ran %q, want %q", got, want)
	}
	for _, c := range calls {
		if c.dir != lab.worktree || !envHas(c.env, closeout.ProgressEnv+"=1") || !envHas(c.env, "FLAI_AGENT=alex") {
			t.Errorf("%q ran in %s with %q, want the worktree, %s=1, and the asker", c.argv, c.dir, gitEnvSaid(c.env), closeout.ProgressEnv)
		}
	}
	log := lab.logged()
	for _, s := range []string{"resolve: " + lab.story + " ", " step started rebase\n", "rebased story/" + lab.story + " onto main\n", " step ended passed rebase\n",
		" step started close-out\n", "passed every step", " step ended passed close-out\n"} {
		if !strings.Contains(log, s) {
			t.Errorf("the job's log lacks %q:\n%s", s, log)
		}
	}
	st, ok, err := closeout.Status(lab.repo.CacheDir(), lab.story)
	if err != nil || !ok || len(st.Steps) != 2 || st.Steps[0].Name != closeout.StepRebase || st.Steps[0].Outcome != "passed" ||
		st.Steps[1].Name != closeout.StepCloseOut || st.Steps[1].Outcome != "passed" {
		t.Errorf("the job's steps: %+v %v %v", st.Steps, ok, err)
	}
	entries := lab.resolveJournal()
	if len(entries) != 2 || entries[0].Outcome != "started" || entries[1].Outcome != "done" || !strings.HasSuffix(entries[1].Detail, "; "+closeout.OutcomeStayed) {
		t.Errorf("the resolve job's journal: %+v", entries)
	}
	if n := len(lab.journal()); n != 0 {
		t.Errorf("a rebase and resolve is not journaled as a close-out flai stats counts: %d entries", n)
	}
}

// S-0804 criterion 3: a rebase that stops on conflicts is aborted, the
// close-out does not run, and the story is sent back to in-progress by flai
// with the send back's thread naming the paths, the branch, the base, and
// that flai aborted the rebase, which flai serve starts its agent on.
func TestARebaseAndResolveWhoseRebaseConflictsAbortsItAndSendsTheStoryBack(t *testing.T) {
	t.Parallel()
	lab, fake, git := newResolveLab(t)
	git.conflicts = []string{"flai/cmd/a.go", "docs/b.md"}
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			_, _ = io.WriteString(out, "flai stream sync: the rebase stopped on conflicts in flai/cmd/a.go, docs/b.md\n")
			_ = os.MkdirAll(filepath.Join(lab.worktree, "rebase-merge"), 0o755)
			return 1
		}
		return 0
	}
	lab.start()
	lab.askResolve()
	job := lab.endedJob()
	if got := argvs(fake.ran()); len(got) != 1 {
		t.Errorf("the close-out ran after the conflicts: %q", got)
	}
	if !git.aborted() {
		t.Error("the rebase was not aborted")
	}
	if _, err := os.Stat(filepath.Join(lab.worktree, "rebase-merge")); err == nil {
		t.Error("the rebase still waits in the worktree")
	}
	if job.Outcome != closeout.OutcomeSentBack || job.Thread == "" || !slices.Equal(job.Conflicts, git.conflicts) ||
		job.Verdict != closeout.VerdictStopped || job.StoppedAt != closeout.StepRebase || !strings.Contains(job.Reason, "flai/cmd/a.go, docs/b.md") {
		t.Errorf("the ended job: %+v", job)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != storygit.SentBackAuthor {
		t.Errorf("the story is %s, last moved by %s; want it sent back to in-progress by %s", status, by, storygit.SentBackAuthor)
	}
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].ID != job.Thread || ths[0].Opener() != storygit.SentBackAuthor || !strings.HasPrefix(ths[0].Title, storygit.SentBackPrefix) {
		t.Fatalf("the send back's thread: %+v", ths)
	}
	text := ths[0].Entries()[0].Text
	for _, s := range []string{"`flai/cmd/a.go`", "`docs/b.md`", "`story/" + lab.story + "`", "`main`", "flai aborted the rebase", "flai stream sync " + lab.story, job.ID} {
		if !strings.Contains(text, s) {
			t.Errorf("the thread's text lacks %q:\n%s", s, text)
		}
	}
	started := false
	for _, th := range acceptSentBackThreads(lab.repo) {
		started = started || th.ID == job.Thread
	}
	if !started {
		t.Error("flai serve does not start the story's agent on the send back's thread")
	}
	if !strings.Contains(lab.logged(), " step ended failed rebase\n") {
		t.Errorf("the job's log does not end the rebase failed:\n%s", lab.logged())
	}
}

// S-0804 criterion 4: a close-out after a clean rebase that stops at a step
// sends the story back with the thread naming the step, its first ten
// findings from the job's own verify report, how many more there are, and
// the job's log.
func TestARebaseAndResolveWhoseCloseOutStopsSendsTheStoryBackNamingTheStep(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	findings := make([]verify.Finding, 12)
	for i := range findings {
		findings[i] = verify.Finding{Name: fmt.Sprintf("TestCase%02d", i), Path: "flai/x_test.go", Line: i + 1, Message: "want 1, got 2"}
	}
	fake.answer = func(argv []string, _ io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			return 0
		}
		rep := verify.Report{Story: lab.story, Commit: "head-a", RanAt: time.Now().UTC(), StoppedAt: "go-test",
			Steps: []verify.Step{{Name: "go-test", Tier: true, State: verify.Failed, Findings: findings, Omitted: 3}}}
		if err := verify.SaveReport(lab.repo, rep); err != nil {
			return 2
		}
		return 1
	}
	lab.start()
	lab.askResolve()
	job := lab.endedJob()
	if job.Outcome != closeout.OutcomeSentBack || job.Verdict != closeout.VerdictStopped || job.StoppedAt != "go-test" || len(job.Conflicts) != 0 ||
		!strings.Contains(job.Reason, "stopped at go-test with 15 findings") {
		t.Errorf("the ended job: %+v", job)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != storygit.SentBackAuthor {
		t.Errorf("the story is %s, last moved by %s; want it sent back", status, by)
	}
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].ID != job.Thread || !strings.HasPrefix(ths[0].Title, storygit.SentBackPrefix) || !strings.Contains(ths[0].Title, "stopped at go-test") {
		t.Fatalf("the send back's thread: %+v", ths)
	}
	text := ths[0].Entries()[0].Text
	for _, s := range []string{"stopped at go-test", "flai/x_test.go:1 TestCase00: want 1, got 2", "flai/x_test.go:10 TestCase09", "- and 5 more", closeout.LogPath(lab.repo.CacheDir(), lab.story)} {
		if !strings.Contains(text, s) {
			t.Errorf("the thread's text lacks %q:\n%s", s, text)
		}
	}
	if strings.Contains(text, "TestCase10") {
		t.Errorf("the thread lists more than ten findings:\n%s", text)
	}
}

// S-0804: a story that left review while its job ran is not sent back, and
// the job's reason says so.
func TestARebaseAndResolveLeavesAStoryNoLongerInReviewAsItIs(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	fake.answer = func(argv []string, _ io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			return 0
		}
		st, err := lab.repo.Get(lab.story)
		if err == nil {
			_, err = lab.repo.Transition(st, workitem.InProgress, "alex", "alex takes it back", time.Now())
		}
		if err != nil {
			return 2
		}
		return 1
	}
	lab.start()
	lab.askResolve()
	job := lab.endedJob()
	if job.Outcome != "" || job.Thread != "" || !strings.Contains(job.Reason, "no longer in review, so flai left it as it is") {
		t.Errorf("the ended job: %+v", job)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != "alex" {
		t.Errorf("the story is %s, last moved by %s; want it as alex left it", status, by)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("nothing sent back, yet threads: %+v", th)
	}
}

// S-0804 criterion 8: a second rebase and resolve, or close-out, asked for
// while the job runs is answered with the job under way and starts nothing.
func TestASecondRequestWhileARebaseAndResolveRunsIsAnsweredWithIt(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	hold := make(chan struct{})
	fake.answer = func(argv []string, _ io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			<-hold
		}
		return 0
	}
	lab.start()
	a := lab.askResolve()
	again := lab.askResolve()
	closeOut := lab.ask(closeout.Request{})
	close(hold)
	if !a.Started || again.Started || again.Job != a.Job || closeOut.Started || closeOut.Job != a.Job {
		t.Errorf("the answers: %+v, then %+v and %+v; want the second and the third to name the first's job", a, again, closeOut)
	}
	lab.endedJob()
	if n := len(lab.verified()); n != 1 {
		t.Errorf("one job, one verified event: %d", n)
	}
	syncs := 0
	for _, c := range fake.ran() {
		if strings.Join(c.argv, " ") == syncArgv(lab.story) {
			syncs++
		}
	}
	if syncs != 1 {
		t.Errorf("the branch was synced %d times, want once", syncs)
	}
}

// S-0804: a rebase and resolve of a story not in review is refused, and
// nothing runs or is journaled.
func TestARebaseAndResolveOfAStoryNotInReviewIsRefused(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	lab.moveStory(workitem.InProgress)
	lab.start()
	a := lab.askResolve()
	if !a.Refused() || a.Job != "" || !strings.Contains(a.Reason, "is in in-progress, not in review") {
		t.Fatalf("the refusal: %+v", a)
	}
	if _, ok, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story); ok || err != nil {
		t.Errorf("a refusal writes no job: %v %v", ok, err)
	}
	if got := fake.ran(); len(got) != 0 {
		t.Errorf("a refusal ran %q", argvs(got))
	}
	if n := len(lab.resolveJournal()); n != 0 {
		t.Errorf("a refusal journals nothing: %d entries", n)
	}
}

// S-0804 criterion 8, ADR-0320 §5: a rebase and resolve flai serve's restart
// interrupted is ended interrupted, journaled as a resolve job, and sends
// nothing back; its story still in review, flai serve asks again for it,
// and the rebase and resolve that starts, marked again, runs.
func TestAnInterruptedRebaseAndResolveSendsNothingBack(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	left := lab.leftJob(closeout.KindResolve, "alex")
	lab.leave(left)
	lab.start()
	if e := lab.awaitVerified(1); e.To != closeout.VerdictInterrupted || e.Fields[inbox.FieldJob] != left.ID {
		t.Errorf("the verified event of the interrupted job: %+v", e)
	}
	again := lab.awaitVerified(2)
	job, _, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story)
	if err != nil || job.ID != again.Fields[inbox.FieldJob] || job.Running() || !job.Again || job.Kind != closeout.KindResolve || job.By != "alex" ||
		job.Message != resolveMessage(lab.story) || job.Outcome != closeout.OutcomeStayed || job.Thread != "" {
		t.Errorf("the rebase and resolve asked again: %+v %v", job, err)
	}
	if got := argvs(fake.ran()); len(got) != 2 || strings.Join(got[0], " ") != syncArgv(lab.story) {
		t.Errorf("the job asked again ran %q, want the sync, then the close-out", got)
	}
	if status, _ := lab.status(); status != workitem.Review {
		t.Errorf("the story is %s; want it left in review", status)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("an interrupted job sent the story back: %+v", th)
	}
	entries := lab.resolveJournal()
	if len(entries) != 3 || entries[0].Verdict != closeout.VerdictInterrupted || !strings.HasSuffix(entries[0].Detail, "; asked again") ||
		entries[1].Outcome != "started" || !strings.HasSuffix(entries[1].Detail, askedAgainByServe) {
		t.Errorf("the resolve jobs' journal: %+v", entries)
	}
}

// ADR-0320 §5: an interrupted rebase and resolve of a story that left review
// is not asked again, and nothing runs.
func TestAnInterruptedRebaseAndResolveOfAStoryThatLeftReviewIsNotAskedAgain(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newResolveLab(t)
	lab.moveStory(workitem.InProgress)
	left := lab.leftJob(closeout.KindResolve, "alex")
	lab.leave(left)
	lab.start()
	lab.awaitVerified(1)
	time.Sleep(50 * time.Millisecond) // looks of the pass after the sweep
	if job, _, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story); err != nil || job.ID != left.ID || job.Verdict != closeout.VerdictInterrupted {
		t.Errorf("a job started for a story out of review: %+v %v", job, err)
	}
	if got := fake.ran(); len(got) != 0 {
		t.Errorf("nothing runs, yet %q ran", argvs(got))
	}
	if _, err := os.Stat(closeout.File(closeout.Dir(lab.repo.CacheDir()), lab.story, closeout.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("a request was written: %v", err)
	}
	if entries := lab.resolveJournal(); len(entries) != 1 || strings.Contains(entries[0].Detail, "asked again") {
		t.Errorf("the interrupted resolve job's journal: %+v", entries)
	}
}

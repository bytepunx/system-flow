package serve

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/sendback"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// statesFlow is a board with two states after review that run something on
// entry: integration, which runs the behavior and smoke suites and the
// cluster's deploy and sends a story back to in-progress, and staging, which
// runs every suite and has no send-back transition; review runs nothing.
const statesFlow = `---
title: Board
updated: 2026-11-18
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: integration, kind: waiting, runs: {suites: [behavior, smoke], kluster: [deploy]}}
  - {name: staging, kind: waiting, runs: {suites: all}}
  - {name: done, kind: accepted}
transitions:
  - {from: backlog, to: ready}
  - {from: ready, to: in-progress, effects: [pull]}
  - {from: in-progress, to: review}
  - {from: review, to: integration}
  - {from: integration, to: staging}
  - {from: staging, to: done, types: [epic, story], by: [acceptance], effects: [accept]}
  - {from: [in-progress, review, integration, staging], to: done, types: [task]}
  - {from: review, to: in-progress, effects: [send-back]}
  - {from: integration, to: in-progress, effects: [send-back]}
---
`

// The states of statesFlow that run something.
const (
	integration = "integration"
	staging     = "staging"
)

// newStateLab is a verify lab on statesFlow whose story is in review, whose
// jobs' processes a fake runner records rather than runs, with this flai at
// hostFlai and resolveGit as its git, and whose pass looks for state jobs at
// every look.
func newStateLab(t *testing.T) (*verifyLab, *fakeCloseOut, *resolveGit) {
	t.Helper()
	lab := newVerifyLab(t, false)
	w, problem := workflow.Parse("design/workflows/board.md", statesFlow)
	if problem != nil {
		t.Fatalf("the board: %v", problem)
	}
	lab.repo.Workflow, lab.pass.o.Workflow = &w, &w
	lab.toReview()
	fake, git := &fakeCloseOut{}, &resolveGit{}
	lab.pass.run, lab.pass.o.Git, lab.pass.stateEvery = fake.run, git, 0
	lab.pass.exe = func() (string, error) { return hostFlai, nil }
	return lab, fake, git
}

// lookOnce runs one look of the pass and waits for the jobs it started.
func (lab *verifyLab) lookOnce() {
	lab.pass.look(context.Background())
	lab.wg.Wait()
}

// moveAt moves the lab's story to to at at, as by moves it.
func (lab *verifyLab) moveAt(to, by string, at time.Time) {
	lab.t.Helper()
	st, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	if _, err := lab.repo.Transition(st, to, by, "the test moves it", at); err != nil {
		lab.t.Fatal(err)
	}
}

// stateJob is the story's last state job, and whether it has had one.
func (lab *verifyLab) stateJob() (closeout.Job, bool) {
	lab.t.Helper()
	j, ok, err := closeout.ReadStateJob(lab.repo.CacheDir(), lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return j, ok
}

// entered is when the lab's story entered its state.
func (lab *verifyLab) entered() time.Time {
	lab.t.Helper()
	st, err := lab.repo.Get(lab.story)
	if err != nil {
		lab.t.Fatal(err)
	}
	return st.EnteredAt().UTC()
}

// stateJournal is the state job's entries journaled so far.
func (lab *verifyLab) stateJournal() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	var out []hostapi.Entry
	for _, e := range lab.entries {
		if e.Action == closeOutAction && e.Method == stateMethod {
			out = append(out, e)
		}
	}
	return out
}

// outcomes are a state job's steps as label:outcome.
func outcomes(j closeout.Job) []string {
	var out []string
	for _, s := range j.StateSteps {
		out = append(out, s.Label()+":"+s.Outcome)
	}
	return out
}

// suiteArgv is this flai's verify of story limited to suite, as a state job
// runs it.
func suiteArgv(story, suite string) string {
	return hostFlai + " verify " + story + " --all --suite " + suite
}

// passing answers the job's sync by rebasing the worktree's head to head-r
// and every verify by passing.
func passing(lab *verifyLab) func(argv []string, out io.Writer) int {
	return func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			_, _ = io.WriteString(out, "rebased story/"+lab.story+" onto main\n")
			_ = os.WriteFile(filepath.Join(lab.worktree, "HEAD"), []byte("head-r\n"), 0o644)
		}
		return 0
	}
}

// S-0810 criteria 1 and 2: a story entering a state whose runs name suites
// and a cluster step gets one job for the entry, which rebases its branch,
// then runs flai verify limited to each suite in its worktree with its
// progress written, and records each step as it starts and ends, the cluster
// step skipped with its reason; a job that passes leaves the story in the
// state with nothing sent back, its output and progress lines in its own log,
// journaled apart from the close-outs, and the story's last job as the
// dashboard reads it. A later look runs nothing for the same entry. Before
// S-0810 nothing ran on entering a state.
func TestAStateJobThatPassesLeavesTheStoryInTheState(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	fake.answer = passing(lab)
	lab.moveStory(integration)
	entry := lab.entered()
	lab.lookOnce()
	want := []string{syncArgv(lab.story), suiteArgv(lab.story, "behavior"), suiteArgv(lab.story, "smoke")}
	got := argvs(fake.ran())
	if len(got) != 3 || strings.Join(got[0], " ") != want[0] || strings.Join(got[1], " ") != want[1] || strings.Join(got[2], " ") != want[2] {
		t.Fatalf("the job ran %q, want %q", got, want)
	}
	for _, c := range fake.ran() {
		if c.dir != lab.worktree || !envHas(c.env, closeout.ProgressEnv+"=1") || !envHas(c.env, "FLAI_AGENT="+storygit.SentBackAuthor) {
			t.Errorf("%q ran in %s with %q, want the worktree, %s=1, and flai", c.argv, c.dir, gitEnvSaid(c.env), closeout.ProgressEnv)
		}
	}
	job, ok := lab.stateJob()
	wantSteps := []string{"rebase:passed", "suite behavior:passed", "suite smoke:passed", "kluster deploy:skipped"}
	if !ok || job.Kind != closeout.KindState || job.State != integration || !job.Entry.Equal(entry) || job.Running() || job.Verdict != closeout.VerdictPassed ||
		job.Outcome != closeout.OutcomeStayedInState || job.Thread != "" || job.Commit != "head-r" || !slices.Equal(outcomes(job), wantSteps) {
		t.Fatalf("the state job: %+v, steps %v; want %v", job, outcomes(job), wantSteps)
	}
	for _, s := range job.StateSteps[:3] {
		if s.Started.IsZero() || s.Ended.Before(s.Started) || s.Reason == "" {
			t.Errorf("step %s was not recorded as it started and ended: %+v", s.Label(), s)
		}
	}
	if deploy := job.StateSteps[3]; !deploy.Started.IsZero() || !strings.Contains(deploy.Reason, "kluster pass") {
		t.Errorf("the cluster step is not skipped with its reason: %+v", deploy)
	}
	if status, by := lab.status(); status != integration || by != "agent-"+lab.story {
		t.Errorf("the story is %s, last moved by %s; want it left in %s", status, by, integration)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("nothing sent back, yet threads: %+v", th)
	}
	log, err := os.ReadFile(closeout.File(closeout.Dir(lab.repo.CacheDir()), lab.story, closeout.StateLogSuffix))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"state: " + lab.story + " ", " step started rebase\n", "rebased story/" + lab.story + " onto main\n", " step ended passed rebase\n",
		" step started suite behavior\n", " step ended passed suite smoke\n", " step ended skipped kluster deploy\n"} {
		if !strings.Contains(string(log), s) {
			t.Errorf("the job's log lacks %q:\n%s", s, log)
		}
	}
	if st, ok, err := closeout.Status(lab.repo.CacheDir(), lab.story); err != nil || !ok || st.ID != job.ID || st.Log != closeout.JobLogPath(lab.repo.CacheDir(), job) || len(st.Steps) != 4 {
		t.Errorf("the story's last job as the dashboard reads it: %+v %v %v", st, ok, err)
	}
	if entries := lab.stateJournal(); len(entries) != 2 || entries[0].Outcome != "started" || entries[1].Outcome != "done" || !strings.HasSuffix(entries[1].Detail, "; "+closeout.OutcomeStayedInState) {
		t.Errorf("the state job's journal: %+v", entries)
	}
	if n := len(lab.journal()); n != 0 {
		t.Errorf("a state job is not journaled as a close-out flai stats counts: %d entries", n)
	}
	if n := len(lab.verified()); n != 0 {
		t.Errorf("a state job announced %d verified events, which a close-out's asker waits for", n)
	}

	lab.lookOnce()
	if n := len(fake.ran()); n != 3 {
		t.Errorf("the same entry ran again: %d calls", n)
	}
}

// S-0810 criterion 2: a rebase that stops on conflicts is aborted, the
// suites after it are skipped, and the story is sent back by flai along the
// state's send-back transition, with the send back's thread naming the
// state, the paths, and the job's log, which flai serve starts its agent on.
func TestAStateJobWhoseRebaseConflictsSendsTheStoryBack(t *testing.T) {
	t.Parallel()
	lab, fake, git := newStateLab(t)
	git.conflicts = []string{"flai/cmd/a.go", "docs/b.md"}
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == syncArgv(lab.story) {
			_, _ = io.WriteString(out, "flai stream sync: the rebase stopped on conflicts\n")
			_ = os.MkdirAll(filepath.Join(lab.worktree, "rebase-merge"), 0o755)
			return 1
		}
		return 0
	}
	lab.moveStory(integration)
	lab.lookOnce()
	if got := argvs(fake.ran()); len(got) != 1 {
		t.Errorf("a suite ran after the conflicts: %q", got)
	}
	if !git.aborted() {
		t.Error("the rebase was not aborted")
	}
	job, _ := lab.stateJob()
	wantSteps := []string{"rebase:failed", "suite behavior:skipped", "suite smoke:skipped", "kluster deploy:skipped"}
	if job.Verdict != closeout.VerdictStopped || job.StoppedAt != closeout.StepRebase || job.Outcome != closeout.OutcomeSentBack || job.Thread == "" ||
		!slices.Equal(job.Conflicts, git.conflicts) || !slices.Equal(outcomes(job), wantSteps) || !strings.Contains(job.StateSteps[1].Reason, "rebase failed before it") {
		t.Errorf("the state job: %+v, steps %v; want %v", job, outcomes(job), wantSteps)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != storygit.SentBackAuthor {
		t.Errorf("the story is %s, last moved by %s; want it sent back to in-progress by %s", status, by, storygit.SentBackAuthor)
	}
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].ID != job.Thread || ths[0].Opener() != storygit.SentBackAuthor ||
		ths[0].Title != sendback.StateJobTitle(lab.story, integration, "the rebase stopped on conflicts in flai/cmd/a.go, docs/b.md") {
		t.Fatalf("the send back's thread: %+v", ths)
	}
	text := ths[0].Entries()[0].Text
	for _, s := range []string{integration, "`flai/cmd/a.go`", "`docs/b.md`", "flai stream sync " + lab.story, closeout.StateLogSuffix, "once " + lab.story + " is back in " + integration} {
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
}

// S-0810 criterion 2: a suite that fails skips the suites after it and sends
// the story back, with the thread quoting the last lines it printed, but its
// progress lines, and naming the job's log.
func TestAStateJobWhoseSuiteFailsSendsTheStoryBack(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	pass := passing(lab)
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == suiteArgv(lab.story, "behavior") {
			_, _ = io.WriteString(out, "verify: "+lab.story+" 2026-11-18T10:00:00Z tier started go-test\n"+
				"failed tier go-test (3s)\n  x_test.go:12: the answer is not 42\n"+
				"verify: "+lab.story+" stopped at go-test (exit 1)\n")
			return 1
		}
		return pass(argv, out)
	}
	lab.moveStory(integration)
	lab.lookOnce()
	job, _ := lab.stateJob()
	wantSteps := []string{"rebase:passed", "suite behavior:failed", "suite smoke:skipped", "kluster deploy:skipped"}
	if job.StoppedAt != "suite behavior" || job.Outcome != closeout.OutcomeSentBack || !slices.Equal(outcomes(job), wantSteps) ||
		!strings.Contains(job.StateSteps[1].Reason, "exited 1: verify: "+lab.story+" stopped at go-test (exit 1)") {
		t.Errorf("the state job: %+v, steps %v; want %v", job, outcomes(job), wantSteps)
	}
	if got := argvs(fake.ran()); len(got) != 2 {
		t.Errorf("the job ran %q; want the sync and the failing suite alone", got)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != storygit.SentBackAuthor {
		t.Errorf("the story is %s, last moved by %s; want it sent back", status, by)
	}
	ths := lab.threadsOn()
	if len(ths) != 1 || ths[0].Title != sendback.StateJobTitle(lab.story, integration, "suite behavior failed") {
		t.Fatalf("the send back's thread: %+v", ths)
	}
	text := ths[0].Entries()[0].Text
	for _, s := range []string{"x_test.go:12: the answer is not 42", "stopped at its step suite behavior", closeout.StateLogSuffix} {
		if !strings.Contains(text, s) {
			t.Errorf("the thread's text lacks %q:\n%s", s, text)
		}
	}
	if strings.Contains(text, "tier started go-test") {
		t.Errorf("the thread quotes a progress line:\n%s", text)
	}
}

// S-0810 criterion 1: a second entry into the state, after the story left
// it, runs the job again, keyed by the new entry.
func TestAStateJobRunsAgainOnALaterEntry(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	fake.answer = passing(lab)
	at := time.Now().Add(time.Minute)
	lab.moveAt(integration, "agent-"+lab.story, at)
	lab.lookOnce()
	first, _ := lab.stateJob()
	lab.moveAt(workitem.InProgress, "alex", at.Add(time.Minute))
	lab.moveAt(workitem.Review, "agent-"+lab.story, at.Add(2*time.Minute))
	lab.lookOnce()
	if n := len(fake.ran()); n != 3 {
		t.Fatalf("a state with no runs ran something: %d calls", n)
	}
	lab.moveAt(integration, "agent-"+lab.story, at.Add(3*time.Minute))
	lab.lookOnce()
	second, _ := lab.stateJob()
	if n := len(fake.ran()); n != 6 {
		t.Errorf("the second entry ran %d calls in all, want 6", n)
	}
	if second.Entry.Equal(first.Entry) || !second.Entry.Equal(at.Add(3*time.Minute).UTC().Truncate(time.Second)) || second.Verdict != closeout.VerdictPassed {
		t.Errorf("the second entry's job %+v, the first's %+v", second, first)
	}
}

// S-0810 criterion 2: a story that left the state while its job ran is left
// where it is when the job fails, with nothing sent back and the job's
// reason saying so.
func TestAStateJobLeavesAStoryThatLeftTheStateAsItIs(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	pass := passing(lab)
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == suiteArgv(lab.story, "behavior") {
			st, err := lab.repo.Get(lab.story)
			if err == nil {
				_, err = lab.repo.Transition(st, workitem.InProgress, "alex", "the test moves it", time.Now().Add(time.Minute))
			}
			if err != nil {
				t.Error(err)
			}
			return 1
		}
		return pass(argv, out)
	}
	lab.moveStory(integration)
	lab.lookOnce()
	job, _ := lab.stateJob()
	if job.Outcome != closeout.OutcomeNotSentBack || job.Thread != "" || !strings.Contains(job.Reason, "flai left it as it is") {
		t.Errorf("the state job: %+v", job)
	}
	if status, by := lab.status(); status != workitem.InProgress || by != "alex" {
		t.Errorf("the story is %s, last moved by %s; want it where alex put it", status, by)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("a story that left the state got a send back's thread: %+v", th)
	}
}

// S-0810: a state job flai serve's restart left running is ended
// interrupted with nothing sent back and no verified event, and its entry
// runs again.
func TestAnInterruptedStateJobSendsNothingBackAndRunsAgain(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	fake.answer = passing(lab)
	lab.moveStory(integration)
	dir := closeout.Dir(lab.repo.CacheDir())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	left := closeout.Job{ID: lab.story + "-20261118T080000Z", Story: lab.story, Kind: closeout.KindState, By: storygit.SentBackAuthor, Started: time.Now().Add(-time.Hour).UTC(),
		State: integration, Entry: lab.entered(), StateSteps: []closeout.StateStep{{Kind: closeout.StateStepRebase, Name: closeout.StepRebase, Started: time.Now().Add(-time.Hour).UTC()}}}
	if err := closeout.WriteJob(dir, left); err != nil {
		t.Fatal(err)
	}
	lab.lookOnce()
	entries := lab.stateJournal()
	if len(entries) < 1 || entries[0].Verdict != closeout.VerdictInterrupted || !strings.Contains(entries[0].Detail, left.ID) {
		t.Errorf("the interrupted state job's journal: %+v", entries)
	}
	if status, _ := lab.status(); status != integration {
		t.Errorf("the story is %s; want it left in %s", status, integration)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("an interrupted state job sent the story back: %+v", th)
	}
	if n := len(lab.verified()); n != 0 {
		t.Errorf("an interrupted state job announced %d verified events", n)
	}
	if job, _ := lab.stateJob(); job.ID == left.ID || job.Verdict != closeout.VerdictPassed || len(fake.ran()) != 3 {
		t.Errorf("the interrupted entry did not run again: %+v, %d calls", job, len(fake.ran()))
	}
}

// S-0810: a state with no runs starts nothing and records nothing.
func TestAStateWithNoRunsStartsNothing(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	lab.lookOnce()
	if _, ok := lab.stateJob(); ok || len(fake.ran()) != 0 {
		t.Errorf("review, which runs nothing, started a job: %d calls", len(fake.ran()))
	}
}

// S-0810 criterion 2: a state that runs every suite runs flai verify with no
// suite named, and with no send-back transition records the failure and
// moves nothing.
func TestAStateJobWithNoSendBackTransitionMovesNothing(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	pass := passing(lab)
	every := hostFlai + " verify " + lab.story + " --all"
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == every {
			return 1
		}
		return pass(argv, out)
	}
	at := time.Now().Add(time.Minute)
	lab.moveAt(integration, "agent-"+lab.story, at)
	lab.lookOnce()
	lab.moveAt(staging, "alex", at.Add(time.Minute))
	lab.lookOnce()
	job, _ := lab.stateJob()
	if job.State != staging || job.StoppedAt != "every suite" || job.Outcome != closeout.OutcomeNotSentBack || job.Thread != "" ||
		!strings.Contains(job.Reason, "no send-back transition out of "+staging) {
		t.Errorf("the state job: %+v", job)
	}
	if got := argvs(fake.ran()); len(got) != 5 || strings.Join(got[4], " ") != every {
		t.Errorf("the job ran %q; want every suite last", got)
	}
	if status, by := lab.status(); status != staging || by != "alex" {
		t.Errorf("the story is %s, last moved by %s; want it left in %s", status, by, staging)
	}
	if th := lab.threadsOn(); len(th) != 0 {
		t.Errorf("a state with no send-back transition opened a thread: %+v", th)
	}
}

// S-0810 criterion 5: on the baseline, a story entering review, whose runs
// name only the cluster's steps, which the kluster pass runs, gets its entry
// recorded with every step skipped and why: nothing is rebased, nothing runs,
// nothing is journaled or sent back, and the story's last job, as the
// dashboard reads it, is not this record.
func TestTheBaselinesReviewRebasesNothingAndSendsNothingBack(t *testing.T) {
	t.Parallel()
	lab, fake, git := newResolveLab(t)
	lab.pass.stateEvery = 0
	lab.lookOnce()
	job, ok := lab.stateJob()
	wantSteps := []string{"rebase:skipped", "kluster deploy:skipped", "kluster checks:skipped"}
	if !ok || job.State != workitem.Review || job.Verdict != closeout.VerdictPassed || job.Outcome != closeout.OutcomeStayedInState || !slices.Equal(outcomes(job), wantSteps) {
		t.Fatalf("the baseline's review job: %+v, steps %v; want %v", job, outcomes(job), wantSteps)
	}
	for _, s := range job.StateSteps {
		if s.Reason == "" || !s.Started.IsZero() {
			t.Errorf("step %s is not skipped with its reason: %+v", s.Label(), s)
		}
	}
	if got := argvs(fake.ran()); len(got) != 0 || git.aborted() {
		t.Errorf("the baseline's review ran %q", got)
	}
	if status, by := lab.status(); status != workitem.Review || by != "agent-"+lab.story {
		t.Errorf("the story is %s, last moved by %s; want it left in review", status, by)
	}
	if th := lab.threadsOn(); len(th) != 0 || len(lab.stateJournal()) != 0 {
		t.Errorf("the baseline's review opened threads %+v or journaled %+v", th, lab.stateJournal())
	}
	if _, ok, err := closeout.Status(lab.repo.CacheDir(), lab.story); ok || err != nil {
		t.Errorf("a state job that ran nothing stands as the story's last job: %v %v", ok, err)
	}
}

// S-0810: a close-out asked for while a state job of the story runs is
// refused, naming the job and its log, rather than started beside it in the
// same worktree.
func TestACloseOutWhileAStateJobRunsIsRefused(t *testing.T) {
	t.Parallel()
	lab, fake, _ := newStateLab(t)
	hold := make(chan struct{})
	pass := passing(lab)
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == suiteArgv(lab.story, "behavior") {
			<-hold
		}
		return pass(argv, out)
	}
	lab.moveStory(integration)
	lab.pass.lookStates(context.Background())
	waitFor(t, "the first suite runs", func() bool { return len(fake.ran()) == 2 })
	dir := closeout.Dir(lab.repo.CacheDir())
	data, err := json.Marshal(closeout.Request{Story: lab.story, By: "agent-" + lab.story, At: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(closeout.File(dir, lab.story, closeout.RequestSuffix), data, 0o644); err != nil {
		t.Fatal(err)
	}
	lab.pass.look(context.Background())
	close(hold)
	lab.wg.Wait()
	answer, err := os.ReadFile(closeout.File(dir, lab.story, closeout.AnswerSuffix))
	if err != nil {
		t.Fatal(err)
	}
	var a closeout.Answer
	if err := json.Unmarshal(answer, &a); err != nil {
		t.Fatal(err)
	}
	if !a.Refused() || a.Job != "" || !strings.Contains(a.Reason, "entry into "+integration+" runs") || !strings.Contains(a.Reason, closeout.StateLogSuffix) {
		t.Errorf("the answer: %+v", a)
	}
	if n := len(fake.ran()); n != 3 {
		t.Errorf("a close-out ran beside the state job: %d calls", n)
	}
}

package serve

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// stubCloseOut stands in for scripts/close-out.sh in a story's worktree: it
// waits while the file hold is there, says its arguments and FLAI_AGENT,
// "commits" by moving HEAD to next-head when there is one, and exits with the
// code in exit-code, 0 when there is none.
const stubCloseOut = `while [ -e hold ]; do sleep 0.02; done
echo "close-out $*"
echo "agent=$FLAI_AGENT"
[ ! -e next-head ] || cp next-head HEAD
exit "$(cat exit-code 2>/dev/null || echo 0)"
`

// headGit answers git rev-parse HEAD in a worktree from its HEAD file.
type headGit struct{}

func (headGit) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && strings.Join(args, " ") == "rev-parse HEAD" {
		data, err := os.ReadFile(filepath.Join(dir, "HEAD"))
		return string(data), err
	}
	return "", errors.New("headGit answers only git rev-parse HEAD")
}

func (headGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return headGit{}.Run(dir, name, args...)
}

func (headGit) LookPath(name string) (string, error) { return name, nil }

// verifyLab is a project with one story in progress, its worktree at
// head-a, and flai serve's verify pass on it.
type verifyLab struct {
	t        *testing.T
	repo     *workitem.Repo
	story    string
	worktree string
	pass     *verifyPass
	wg       sync.WaitGroup
	stop     context.CancelFunc

	mu      sync.Mutex
	entries []hostapi.Entry
	kept    []string // story, commit, and the kept report's commit
}

// newVerifyLab makes the lab; with script, the worktree has the stub
// close-out, held until the test lets it go.
func newVerifyLab(t *testing.T, script bool) *verifyLab {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	st, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Close me", Parent: epic.ID, Owner: "alex", Now: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	lab := &verifyLab{t: t, repo: repo, story: st.ID, worktree: repo.WorktreePath(st.ID)}
	lab.write("HEAD", "head-a\n")
	if script {
		lab.write("hold", "")
		lab.write(closeOutScript, stubCloseOut)
	}
	o := Options{Logger: slog.New(slog.DiscardHandler), Now: time.Now, Git: headGit{},
		Host: hostapi.Host{Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.entries = append(lab.entries, e)
		}}}
	lab.pass = newVerifyPass(o, Entry{Key: "t", Name: "t", Root: root}, &lab.wg)
	lab.pass.every = 5 * time.Millisecond
	lab.pass.keepAt = func(_ *workitem.Repo, commit string, rep verify.Report) error {
		lab.mu.Lock()
		defer lab.mu.Unlock()
		lab.kept = append(lab.kept, rep.Story, commit, rep.Commit)
		return nil
	}
	return lab
}

// write writes a file of the worktree.
func (lab *verifyLab) write(rel, content string) {
	lab.t.Helper()
	path := filepath.Join(lab.worktree, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		lab.t.Fatal(err)
	}
}

// release lets the stub close-out run on.
func (lab *verifyLab) release() {
	lab.t.Helper()
	if err := os.Remove(filepath.Join(lab.worktree, "hold")); err != nil {
		lab.t.Fatal(err)
	}
}

// start starts the pass, stopped and waited for when the test ends.
func (lab *verifyLab) start() {
	ctx, cancel := context.WithCancel(context.Background())
	lab.stop = cancel
	lab.pass.start(ctx)
	lab.t.Cleanup(func() {
		cancel()
		lab.wg.Wait()
	})
}

// ask asks for the story's close-out as its agent does.
func (lab *verifyLab) ask(req closeout.Request) closeout.Answer {
	lab.t.Helper()
	req.Story = lab.story
	if req.By == "" {
		req.By = "agent-" + lab.story
	}
	a, err := closeout.Ask(context.Background(), lab.repo.CacheDir(), req, closeout.Wait{Take: 5 * time.Second, Answer: 5 * time.Second, Poll: 5 * time.Millisecond})
	if err != nil {
		lab.t.Fatal(err)
	}
	return a
}

// verified are the verified events in the project's event log.
func (lab *verifyLab) verified() []inbox.LoggedEvent {
	var out []inbox.LoggedEvent
	for _, e := range inbox.LoggedEvents(lab.repo) {
		if e.Kind == inbox.KindVerified {
			out = append(out, e)
		}
	}
	return out
}

// awaitVerified waits for the n-th verified event and returns it.
func (lab *verifyLab) awaitVerified(n int) inbox.LoggedEvent {
	lab.t.Helper()
	waitFor(lab.t, "the verified event", func() bool { return len(lab.verified()) >= n })
	return lab.verified()[n-1]
}

// journal is the close-out entries journaled so far.
func (lab *verifyLab) journal() []hostapi.Entry {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	var out []hostapi.Entry
	for _, e := range lab.entries {
		if e.Action == closeOutAction && e.Method == closeOutMethod {
			out = append(out, e)
		}
	}
	return out
}

// saveReport stores the story's newest verify report, run now.
func (lab *verifyLab) saveReport(commit string, passed bool, stoppedAt string) {
	lab.t.Helper()
	if err := verify.SaveReport(lab.repo, verify.Report{Story: lab.story, Commit: commit, RanAt: time.Now().UTC(), Passed: passed, StoppedAt: stoppedAt}); err != nil {
		lab.t.Fatal(err)
	}
}

// S-0546 criteria 1 and 4: a request starts a job that runs the worktree's
// close-out with the asker as FLAI_AGENT and the message; a second request
// while it runs is answered with its id and starts nothing; its end is a
// verified event with the verdict and a journal entry with its duration,
// beside the start's; and a pass whose script committed keeps the report
// under the new head.
func TestACloseOutRequestRunsAJobThatEndsWithAVerifiedEvent(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, true)
	lab.start()
	a := lab.ask(closeout.Request{Message: "close out"})
	if !a.Started || !strings.HasPrefix(a.Job, lab.story+"-") || len(a.Job) != len(lab.story)+len("-20261031T090000Z") || a.Refused() || a.Recorded {
		t.Fatalf("the answer that starts the job: %+v", a)
	}
	again := lab.ask(closeout.Request{})
	if again.Started || again.Job != a.Job {
		t.Fatalf("a second request while the job runs: %+v, want job %s not started", again, a.Job)
	}
	if job, ok, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story); err != nil || !ok || !job.Running() || job.ID != a.Job || job.By != "agent-"+lab.story {
		t.Fatalf("the running job's state: %+v %v %v", job, ok, err)
	}
	lab.saveReport("head-a", true, "")
	lab.write("next-head", "head-b\n")
	lab.release()

	e := lab.awaitVerified(1)
	want := map[string]string{inbox.FieldJob: a.Job, inbox.FieldCommit: "head-b", inbox.FieldVerdict: closeout.VerdictPassed, inbox.FieldStoppedAt: "", inbox.FieldRecorded: "false"}
	for k, v := range want {
		if e.Fields[k] != v {
			t.Errorf("the verified event's %s: %q, want %q (%+v)", k, e.Fields[k], v, e)
		}
	}
	if e.ID != lab.story || e.To != closeout.VerdictPassed || e.By != verifyServeAuthor || e.Fields[inbox.FieldDuration] == "" {
		t.Errorf("the verified event: %+v", e)
	}
	if n := len(lab.verified()); n != 1 {
		t.Errorf("one job, one verified event: %d", n)
	}
	entries := lab.journal()
	if len(entries) != 2 {
		t.Fatalf("the start and the end journaled: %+v", entries)
	}
	if s := entries[0]; s.Outcome != "started" || s.By != "agent-"+lab.story || s.Detail != lab.story+" job "+a.Job || s.Verdict != "" || s.Project != "t" {
		t.Errorf("the start entry: %+v", s)
	}
	if end := entries[1]; end.Outcome != "done" || end.Verdict != closeout.VerdictPassed || end.DurationMS <= 0 || !strings.HasPrefix(end.Detail, lab.story+" job "+a.Job+": passed in ") {
		t.Errorf("the end entry: %+v", end)
	}
	job, _, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story)
	if err != nil || job.Running() || job.ExitCode == nil || *job.ExitCode != 0 || job.Verdict != closeout.VerdictPassed || job.Commit != "head-b" || job.DurationMS <= 0 {
		t.Errorf("the ended job's state: %+v %v", job, err)
	}
	lab.mu.Lock()
	kept := strings.Join(lab.kept, " ")
	lab.mu.Unlock()
	if kept != lab.story+" head-b head-a" {
		t.Errorf("the passing report kept under the commit the close-out made: %q", kept)
	}
	out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if err != nil || !strings.Contains(string(out), "close-out "+lab.story+" -m close out") || !strings.Contains(string(out), "agent=agent-"+lab.story) {
		t.Errorf("the job's log: %q %v", out, err)
	}
}

// ADR-0229 §3: a close-out that exits non-zero is stopped at the step its own
// verify run names, journaled failed.
func TestACloseOutThatStopsIsStoppedAtTheStepTheRecordNames(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, true)
	lab.start()
	a := lab.ask(closeout.Request{})
	lab.saveReport("head-a", false, "go-test")
	lab.write("exit-code", "1\n")
	lab.release()
	e := lab.awaitVerified(1)
	if e.To != closeout.VerdictStopped || e.Fields[inbox.FieldStoppedAt] != "go-test" || e.Fields[inbox.FieldJob] != a.Job {
		t.Errorf("the verified event of a stop: %+v", e)
	}
	entries := lab.journal()
	if end := entries[len(entries)-1]; end.Outcome != "failed" || end.Verdict != closeout.VerdictStopped || !strings.Contains(end.Detail, ": stopped at go-test in ") {
		t.Errorf("the end entry of a stop: %+v", end)
	}
	lab.mu.Lock()
	defer lab.mu.Unlock()
	if len(lab.kept) != 0 {
		t.Errorf("a stop keeps nothing: %v", lab.kept)
	}
}

// ADR-0229 §5: a close-out asked at a commit the record holds a verdict for
// is answered with it and runs nothing, journaled as one end entry from the
// record and announced as recorded; fresh always runs.
func TestACloseOutTheRecordHoldsIsAnsweredWithoutAJob(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, true)
	ranAt := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	lab.pass.recorded = func(_ *workitem.Repo, story, worktree string, _ bool) (verify.Report, bool, error) {
		if story != lab.story || worktree != lab.worktree {
			return verify.Report{}, false, nil
		}
		return verify.Report{Story: story, Commit: "head-a", RanAt: ranAt, Passed: true}, true, nil
	}
	lab.start()
	a := lab.ask(closeout.Request{})
	if !a.Recorded || a.Started || a.Job != "" || a.Verdict != closeout.VerdictPassed || a.Commit != "head-a" || !a.RanAt.Equal(ranAt) {
		t.Fatalf("the answer from the record: %+v", a)
	}
	if _, ok, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story); ok || err != nil {
		t.Errorf("no job for an answer from the record: %v %v", ok, err)
	}
	e := lab.awaitVerified(1)
	if e.Fields[inbox.FieldRecorded] != "true" || e.To != closeout.VerdictPassed || e.Fields[inbox.FieldJob] != "" || e.Fields[inbox.FieldCommit] != "head-a" || e.Fields[inbox.FieldDuration] != "" {
		t.Errorf("the verified event of an answer from the record: %+v", e)
	}
	entries := lab.journal()
	if len(entries) != 1 || entries[0].Outcome != "done" || entries[0].Verdict != closeout.VerdictPassed || entries[0].DurationMS != 0 || !strings.HasPrefix(entries[0].Detail, fromTheRecord) {
		t.Fatalf("one end entry from the record: %+v", entries)
	}

	lab.release()
	if fresh := lab.ask(closeout.Request{Fresh: true}); !fresh.Started || fresh.Recorded {
		t.Errorf("fresh runs whatever the record holds: %+v", fresh)
	}
	if e := lab.awaitVerified(2); e.Fields[inbox.FieldRecorded] != "false" || e.To != closeout.VerdictPassed {
		t.Errorf("the fresh job's verified event: %+v", e)
	}
}

// ADR-0229 §2: a job an earlier flai serve left running is ended interrupted
// at the pass's first look, journaled failed, and announced.
func TestAJobLeftRunningIsEndedInterruptedWhenServeStarts(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, false)
	dir := closeout.Dir(lab.repo.CacheDir())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	left := closeout.Job{ID: lab.story + "-20261031T080000Z", Story: lab.story, By: "agent-" + lab.story, Started: time.Now().Add(-time.Hour).UTC()}
	if err := closeout.WriteJob(dir, left); err != nil {
		t.Fatal(err)
	}
	lab.start()
	e := lab.awaitVerified(1)
	if e.To != closeout.VerdictInterrupted || e.Fields[inbox.FieldJob] != left.ID || e.Fields[inbox.FieldCommit] != "head-a" {
		t.Errorf("the verified event of an interrupted job: %+v", e)
	}
	job, _, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story)
	if err != nil || job.Running() || job.Verdict != closeout.VerdictInterrupted || job.ExitCode == nil || *job.ExitCode != -1 {
		t.Errorf("the interrupted job's state: %+v %v", job, err)
	}
	entries := lab.journal()
	if len(entries) != 1 || entries[0].Outcome != "failed" || entries[0].Verdict != closeout.VerdictInterrupted || entries[0].By != left.By {
		t.Errorf("the interrupted job journaled failed: %+v", entries)
	}
}

// A job running when flai serve stops is stopped with its process group and
// left running in its state, for the next flai serve to end interrupted;
// nothing announces an end it did not reach.
func TestAJobIsStoppedWithServeAndLeftForTheNextStart(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, true)
	lab.start()
	if a := lab.ask(closeout.Request{}); !a.Started {
		t.Fatalf("the job did not start: %+v", a)
	}
	lab.stop()
	done := make(chan struct{})
	go func() {
		lab.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("the job was not stopped with flai serve")
	}
	if job, ok, err := closeout.ReadJob(lab.repo.CacheDir(), lab.story); err != nil || !ok || !job.Running() {
		t.Errorf("the stopped job's state: %+v %v %v", job, ok, err)
	}
	if n := len(lab.verified()); n != 0 {
		t.Errorf("no verified event for a job flai serve stopped: %d", n)
	}
}

// Where the worktree has no close-out script, the job runs this flai's
// verify with --record-issues, and --all and --fresh as asked.
func TestAWorktreeWithoutTheScriptRunsFlaiVerify(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, false)
	stub := filepath.Join(t.TempDir(), "flai")
	if err := os.WriteFile(stub, []byte("#!/bin/sh\necho \"flai $*\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.pass.exe = func() (string, error) { return stub, nil }
	lab.start()
	lab.ask(closeout.Request{All: true, Fresh: true, Message: "unused"})
	if e := lab.awaitVerified(1); e.To != closeout.VerdictPassed {
		t.Errorf("the verified event: %+v", e)
	}
	out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if err != nil || strings.TrimSpace(string(out)) != "flai verify "+lab.story+" --record-issues --all --fresh" {
		t.Errorf("the command run: %q %v", out, err)
	}
}

// A story with no worktree is refused, saying how to open it, and nothing
// runs or is journaled.
func TestACloseOutWithNoWorktreeIsRefused(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, false)
	if err := os.RemoveAll(lab.worktree); err != nil {
		t.Fatal(err)
	}
	lab.start()
	a := lab.ask(closeout.Request{})
	if !a.Refused() || a.Job != "" || !strings.Contains(a.Reason, "flai stream open "+lab.story) {
		t.Fatalf("the refusal: %+v", a)
	}
	if entries := lab.journal(); len(entries) != 0 {
		t.Errorf("a refusal journals no close-out: %+v", entries)
	}
}

package serve

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/harness"
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
echo "fresh=${CLOSE_OUT_FRESH:-} all=${CLOSE_OUT_ALL:-}"
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

// I-0280 (S-0671): a fresh request at a commit the record passed ran the
// worktree's script, which answered from the record and ran nothing, since
// the job dropped fresh and all for it. The job now asks the script for them
// by CLOSE_OUT_FRESH and CLOSE_OUT_ALL, never by a flag an older script would
// take for the story, and sets neither when the request does not ask.
func TestAFreshOrAllCloseOutAsksTheWorktreesScriptForThem(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		req  closeout.Request
		want string
	}{
		{"fresh and all", closeout.Request{Fresh: true, All: true, Message: "close out"}, "fresh=1 all=1"},
		{"fresh", closeout.Request{Fresh: true}, "fresh=1 all="},
		{"all", closeout.Request{All: true}, "fresh= all=1"},
		{"neither", closeout.Request{}, "fresh= all="},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			lab := newVerifyLab(t, true)
			lab.release()
			lab.start()
			if a := lab.ask(c.req); !a.Started {
				t.Fatalf("the answer: %+v", a)
			}
			lab.awaitVerified(1)
			out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			args := "close-out " + lab.story
			if c.req.Message != "" {
				args += " -m " + c.req.Message
			}
			if lines[0] != args {
				t.Errorf("the script's arguments: no flag before the story, the story, and the message: %q", out)
			}
			if !slices.Contains(lines, c.want) {
				t.Errorf("the script's environment: want %q in %q", c.want, out)
			}
		})
	}
}

// The job's environment never carries flai serve's own CLOSE_OUT_FRESH or
// CLOSE_OUT_ALL: only the request sets them.
func TestACloseOutJobDropsServesOwnFreshAndAll(t *testing.T) {
	t.Parallel()
	env := closeOutEnv([]string{"PATH=/bin", "CLOSE_OUT_FRESH=1", "CLOSE_OUT_ALL=1"}, "agent-S-0001")
	if !slices.Equal(env, []string{"PATH=/bin", "FLAI_AGENT=agent-S-0001"}) {
		t.Errorf("the job's environment: %q", env)
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

// closeOutCall is one process a close-out job ran.
type closeOutCall struct {
	dir  string
	env  []string
	argv []string
}

// fakeCloseOut is a closeOutRunner that runs nothing: it records each call
// and answers it with answer, which may write to out; nil answers 0.
type fakeCloseOut struct {
	mu     sync.Mutex
	calls  []closeOutCall
	answer func(argv []string, out io.Writer) int
}

func (f *fakeCloseOut) run(_ context.Context, dir string, env []string, out io.Writer, argv []string) (int, error) {
	f.mu.Lock()
	f.calls = append(f.calls, closeOutCall{dir: dir, env: slices.Clone(env), argv: slices.Clone(argv)})
	answer := f.answer
	f.mu.Unlock()
	if answer == nil {
		return 0, nil
	}
	return answer(argv, out), nil
}

// argvs is the commands of calls, which a failure prints without their
// environments.
func argvs(calls []closeOutCall) [][]string {
	out := make([][]string, 0, len(calls))
	for _, c := range calls {
		out = append(out, c.argv)
	}
	return out
}

func (f *fakeCloseOut) ran() []closeOutCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls)
}

// imageDocker is a docker client on PATH that has the image when has is
// true, and runs nothing else.
type imageDocker struct{ has bool }

func (d imageDocker) Run(_, _ string, args ...string) (string, error) {
	if len(args) > 1 && args[0] == "image" && args[1] == "inspect" {
		if d.has {
			return `{"Id":"sha256:abc","Config":{"Labels":null}}` + "\n", nil
		}
		return "Error: No such image: " + args[len(args)-1], errors.New("exit status 1")
	}
	return "", errors.New("imageDocker answers only docker image inspect")
}

func (d imageDocker) RunInput(dir, name, _ string, args ...string) (string, error) {
	return d.Run(dir, name, args...)
}

func (imageDocker) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

// hostFlai is the flai the host's steps of a container story's close-out run.
const hostFlai = "/opt/host/bin/flai"

// newContainerLab is a verify lab whose story's agent runs on the container
// harness, the project's default agent, with its worktree linked to its
// registration in the main checkout's .git as git worktree add leaves it, the
// stub close-out script in it, docker answering for the image when has, and
// the job's processes recorded by a fake runner rather than run.
func newContainerLab(t *testing.T, has bool) (*verifyLab, *fakeCloseOut) {
	t.Helper()
	lab := newVerifyLab(t, true)
	lab.release()
	root := lab.repo.Root
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nagent:\n  harness: container\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reg := filepath.Join(root, ".git", "worktrees", lab.story)
	if err := os.MkdirAll(reg, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reg, "commondir"), []byte("../..\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.write(".git", "gitdir: "+reg+"\n")
	fake := &fakeCloseOut{}
	lab.pass.run = fake.run
	lab.pass.exe = func() (string, error) { return hostFlai, nil }
	lab.pass.o.Docker = imageDocker{has: has}
	return lab, fake
}

// envHas says whether env holds kv.
func envHas(env []string, kv string) bool { return slices.Contains(env, kv) }

// gitEnvSaid is the part of env a failure may print: git's variables and
// FLAI_AGENT, never the host's others, which may hold secrets.
func gitEnvSaid(env []string) []string {
	return slices.DeleteFunc(slices.Clone(env), func(kv string) bool {
		return !strings.HasPrefix(kv, "GIT_") && !strings.HasPrefix(kv, "FLAI_AGENT=")
	})
}

// I-0295, ADR-0274: the close-out of a story whose agent runs in a container
// runs flai verify --record-issues in a docker container with the story's
// mounts, and runs nothing from the worktree on the host: not the worktree's
// scripts/close-out.sh, and no program under the worktree. The host runs only
// docker, git, and its own flai. On flai serve before ADR-0274 the job ran sh
// scripts/close-out.sh in the worktree on the host.
func TestAContainerStorysCloseOutRunsItsVerifyInItsContainerAndNothingFromTheWorktreeOnTheHost(t *testing.T) {
	t.Parallel()
	lab, fake := newContainerLab(t, true)
	lab.start()
	a := lab.ask(closeout.Request{Message: "close out"})
	if !a.Started {
		t.Fatalf("the job did not start: %+v", a)
	}
	if e := lab.awaitVerified(1); e.To != closeout.VerdictPassed {
		t.Errorf("the verified event: %+v", e)
	}
	calls := fake.ran()
	if len(calls) == 0 {
		t.Fatal("the job ran nothing")
	}
	for _, c := range calls {
		program := c.argv[0]
		if program == "sh" || strings.HasPrefix(program, lab.worktree) || !filepath.IsAbs(program) && strings.Contains(program, "/") || slices.Contains(c.argv, closeOutScript) {
			t.Errorf("the job ran %q on the host, from the worktree", c.argv)
		}
		if program != "docker" && program != "git" && program != hostFlai {
			t.Errorf("the host ran %q, which is neither docker, git, nor its own flai", c.argv)
		}
	}
	root := lab.repo.Root
	reg := filepath.Join(root, ".git", "worktrees", lab.story)
	verify := calls[0]
	whole := " " + strings.Join(verify.argv, " ") + " "
	for _, want := range []string{
		"docker run --rm --name flai-closeout-" + lab.story + "-",
		" -v " + lab.worktree + ":" + lab.worktree + " -w " + lab.worktree + " ",
		" -v " + filepath.Join(root, ".git") + ":" + filepath.Join(root, ".git") + ":ro ",
		" -v " + reg + ":" + reg + " ",
		" -v " + filepath.Join(root, "wip") + ":" + filepath.Join(root, "wip") + ":ro ",
		" -v " + filepath.Join(root, "system-flow.yaml") + ":" + filepath.Join(root, "system-flow.yaml") + ":ro ",
		" -e FLAI_AGENT ",
		" --entrypoint " + harness.ImageFlai + " " + harness.DefaultImages[harness.ClaudeCode] + " verify " + lab.story + " --record-issues ",
	} {
		if !strings.Contains(whole, want) {
			t.Errorf("the verify's docker run lacks %q:\n%s", want, whole)
		}
	}
	if !envHas(verify.env, "FLAI_AGENT=agent-"+lab.story) || !envHas(verify.env, "GIT_CONFIG_VALUE_0="+root) {
		t.Errorf("the docker client's environment lacks the asker or safe.directory: %q", gitEnvSaid(verify.env))
	}
	last := calls[len(calls)-1]
	if strings.Join(last.argv, " ") != "git status --porcelain" {
		t.Errorf("the last step: %q, want the clean worktree check", last.argv)
	}
	out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if err != nil || !strings.HasSuffix(string(out), "close-out: "+lab.story+" passed every step; ready to move to review\n") {
		t.Errorf("the job's log: %q %v", out, err)
	}
}

// ADR-0274 §2, §3, ADR-0156: the host commits what is outstanding with git and
// the request's message, runs its own flai verify --sync-only, and, when main
// moved, flai stream sync and the next round in a container of its own name.
// Every git it runs names its git dir, common dir, and work tree from the main
// checkout and runs the main checkout's hooks with no fsmonitor; its own flai
// gets the hooks and fsmonitor alone.
func TestAContainerStorysCloseOutCommitsAndSyncsOnTheHost(t *testing.T) {
	t.Parallel()
	lab, fake := newContainerLab(t, true)
	statuses, syncs := []string{" M a.go\n", "", ""}, 0
	fake.answer = func(argv []string, out io.Writer) int {
		switch strings.Join(argv[1:], " ") {
		case "status --porcelain":
			s := statuses[0]
			statuses = statuses[1:]
			_, _ = io.WriteString(out, s)
		case "verify " + lab.story + " --sync-only":
			syncs++
			if syncs == 1 {
				_, _ = io.WriteString(out, "verify: "+lab.story+" stopped at sync (main has 1 commit the branch lacks)\n")
				return 1
			}
		}
		return 0
	}
	lab.start()
	lab.ask(closeout.Request{Message: "close out"})
	if e := lab.awaitVerified(1); e.To != closeout.VerdictPassed {
		t.Errorf("the verified event: %+v", e)
	}
	root := lab.repo.Root
	var got []string
	for _, c := range fake.ran() {
		switch c.argv[0] {
		case "docker":
			got = append(got, "docker "+c.argv[4])
		default:
			got = append(got, strings.Join(c.argv, " "))
		}
		hooks := envHas(c.env, "GIT_CONFIG_KEY_0=core.hooksPath") && envHas(c.env, "GIT_CONFIG_VALUE_0="+filepath.Join(root, ".git", "hooks")) &&
			envHas(c.env, "GIT_CONFIG_KEY_1=core.fsmonitor") && envHas(c.env, "GIT_CONFIG_VALUE_1=false")
		gitDir := envHas(c.env, "GIT_DIR="+filepath.Join(root, ".git", "worktrees", lab.story)) && envHas(c.env, "GIT_COMMON_DIR="+filepath.Join(root, ".git")) &&
			envHas(c.env, "GIT_WORK_TREE="+lab.worktree)
		switch c.argv[0] {
		case "git":
			if !hooks || !gitDir {
				t.Errorf("the host's %q runs without the main checkout's git dir, hooks, and no fsmonitor: %q", c.argv, gitEnvSaid(c.env))
			}
		case hostFlai:
			if !hooks || gitDir {
				t.Errorf("the host's %q should have the main checkout's hooks and no fsmonitor, and no git dir: %q", c.argv, gitEnvSaid(c.env))
			}
		}
	}
	id := strings.TrimPrefix(lab.verified()[0].Fields[inbox.FieldJob], lab.story+"-")
	want := []string{
		"docker flai-closeout-" + lab.story + "-" + id,
		"git status --porcelain", "git add -A", "git commit -m close out",
		hostFlai + " verify " + lab.story + " --sync-only",
		hostFlai + " stream sync " + lab.story,
		"docker flai-closeout-" + lab.story + "-" + id + "-r2",
		"git status --porcelain",
		hostFlai + " verify " + lab.story + " --sync-only",
		"git status --porcelain",
	}
	if !slices.Equal(got, want) {
		t.Errorf("the steps\n%q\nwant\n%q", got, want)
	}
	out, err := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if err != nil || !strings.HasSuffix(string(out), "close-out: "+lab.story+" passed every step in round 2 of 3; ready to move to review\n") {
		t.Errorf("the job's log: %q %v", out, err)
	}
}

// As the script does, a container story's close-out with something
// outstanding and no message stops at the commit, and commits nothing.
func TestAContainerStorysCloseOutWithChangesAndNoMessageStopsAtTheCommit(t *testing.T) {
	t.Parallel()
	lab, fake := newContainerLab(t, true)
	fake.answer = func(argv []string, out io.Writer) int {
		if strings.Join(argv, " ") == "git status --porcelain" {
			_, _ = io.WriteString(out, "?? new.go\n")
		}
		return 0
	}
	lab.start()
	lab.ask(closeout.Request{})
	if e := lab.awaitVerified(1); e.To != closeout.VerdictFailed {
		t.Errorf("the verified event: %+v", e)
	}
	for _, c := range fake.ran() {
		if c.argv[0] == "git" && c.argv[1] != "status" {
			t.Errorf("committed with no message: %q", c.argv)
		}
	}
	out, _ := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if !strings.Contains(string(out), "uncommitted changes and no message") || !strings.HasSuffix(string(out), "stopped at commit (exit 1)\n") {
		t.Errorf("the job's log: %q", out)
	}
}

// ADR-0274 §5: a container story whose image is not on this host is not
// closed out on the host instead: the job runs nothing, its verdict is
// failed, and its log and journal say how to build the image.
func TestAContainerStorysCloseOutWithNoImageFailsAndRunsNothing(t *testing.T) {
	t.Parallel()
	lab, fake := newContainerLab(t, false)
	lab.start()
	lab.ask(closeout.Request{Message: "close out"})
	if e := lab.awaitVerified(1); e.To != closeout.VerdictFailed {
		t.Errorf("the verified event: %+v", e)
	}
	if calls := fake.ran(); len(calls) != 0 {
		t.Errorf("a story with no image ran %q", argvs(calls))
	}
	out, _ := os.ReadFile(closeout.LogPath(lab.repo.CacheDir(), lab.story))
	if !strings.Contains(string(out), "nothing runs on the host") || !strings.Contains(string(out), "make agent-images") {
		t.Errorf("the job's log: %q", out)
	}
	entries := lab.journal()
	if end := entries[len(entries)-1]; end.Verdict != closeout.VerdictFailed || !strings.Contains(end.Detail, "make agent-images") {
		t.Errorf("the end entry: %+v", end)
	}
}

// ADR-0274 §6: a story whose agent runs on any other harness keeps the job it
// had: sh scripts/close-out.sh with the message, in the worktree.
func TestAHostHarnessStorysCloseOutStillRunsTheScript(t *testing.T) {
	t.Parallel()
	lab := newVerifyLab(t, true)
	lab.release()
	fake := &fakeCloseOut{}
	lab.pass.run = fake.run
	lab.pass.o.Docker = imageDocker{has: true}
	lab.start()
	lab.ask(closeout.Request{Message: "close out"})
	lab.awaitVerified(1)
	calls := fake.ran()
	if len(calls) != 1 || calls[0].dir != lab.worktree || !slices.Equal(calls[0].argv, []string{"sh", closeOutScript, lab.story, "-m", "close out"}) {
		t.Errorf("a host-harness story's close-out: %+v", argvs(calls))
	}
}

// ADR-0274 §3: the container can write the worktree's .git file and the
// registration's commondir; a host step refuses to run git through either
// once it no longer names the main checkout's, rather than read a repository
// the container chose, and the request is refused before the record is
// asked.
func TestAContainerStorysCloseOutRefusesGitFilesTheContainerRedirected(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name  string
		write func(lab *verifyLab, evil string)
	}{
		{"the worktree's .git file", func(lab *verifyLab, evil string) { lab.write(".git", "gitdir: "+evil+"\n") }},
		{"the registration's commondir", func(lab *verifyLab, evil string) {
			if err := os.WriteFile(filepath.Join(lab.repo.Root, ".git", "worktrees", lab.story, "commondir"), []byte(evil+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			lab, fake := newContainerLab(t, true)
			evil := filepath.Join(lab.worktree, "evil.git")
			if err := os.MkdirAll(evil, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := gitPointers(lab.repo.Root, lab.story, lab.worktree); err != nil {
				t.Fatalf("the files as git worktree add leaves them: %v", err)
			}
			c.write(lab, evil)
			err := gitPointers(lab.repo.Root, lab.story, lab.worktree)
			if err == nil || !strings.Contains(err.Error(), "git worktree repair "+lab.worktree) {
				t.Fatalf("the redirected files: %v", err)
			}
			lab.start()
			if a := lab.ask(closeout.Request{Message: "close out"}); !a.Refused() || !strings.Contains(a.Reason, "git worktree repair") {
				t.Errorf("the request: %+v", a)
			}
			if calls := fake.ran(); len(calls) != 0 {
				t.Errorf("a redirected story ran %q", argvs(calls))
			}
		})
	}
}

// ADR-0274 §3, with real git: a core.hooksPath and a core.fsmonitor the
// container wrote into the registration's config.worktree, which git reads
// once the main checkout turns extensions.worktreeConfig on, run nothing when
// the host commits the story's work; the main checkout's own hook runs, and
// the commit lands on the story's branch in the main checkout's repository.
func TestTheHostsCommitRunsNoHookTheContainerWrote(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not on PATH")
	}
	root, story := t.TempDir(), "S-0001"
	worktree := filepath.Join(root, ".flai-cache", "worktrees", story)
	env := append(withoutGitEnv(os.Environ()), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@example.com", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@example.com")
	git := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir, cmd.Env = dir, env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	writeFile := func(path, content string, mode os.FileMode) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), mode); err != nil {
			t.Fatal(err)
		}
	}
	git(root, "init", "-q", "-b", "main")
	writeFile(filepath.Join(root, "system-flow.yaml"), "version: 1\n", 0o644)
	git(root, "add", "-A")
	git(root, "commit", "-q", "-m", "init")
	git(root, "worktree", "add", "-q", "-b", "story/"+story, worktree)
	git(root, "config", "extensions.worktreeConfig", "true")

	marks := t.TempDir()
	evilHooks := filepath.Join(worktree, "evil-hooks")
	writeFile(filepath.Join(evilHooks, "pre-commit"), "#!/bin/sh\ntouch "+filepath.Join(marks, "container-hook")+"\n", 0o755)
	writeFile(filepath.Join(evilHooks, "fsmonitor"), "#!/bin/sh\ntouch "+filepath.Join(marks, "container-fsmonitor")+"\n", 0o755)
	writeFile(filepath.Join(root, ".git", "worktrees", story, "config.worktree"),
		"[core]\n\thooksPath = "+evilHooks+"\n\tfsmonitor = "+filepath.Join(evilHooks, "fsmonitor")+"\n", 0o644)
	writeFile(filepath.Join(root, ".git", "hooks", "pre-commit"), "#!/bin/sh\ntouch "+filepath.Join(marks, "main-hook")+"\n", 0o755)
	if got := git(worktree, "config", "core.hooksPath"); got != evilHooks {
		t.Fatalf("the container's hooksPath should be live for a git that trusts the worktree's files: %q", got)
	}
	writeFile(filepath.Join(worktree, "work.go"), "package work\n", 0o644)

	p := &verifyPass{run: runCloseOut, e: Entry{Root: root}}
	c := &containerCloseOut{p: p, story: story, worktree: worktree, req: closeout.Request{Story: story, Message: "the story's work"}}
	gitEnv := append(append(slices.Clone(env), hostGitConfig(root)...), storyGitDir(root, story, worktree)...)
	var out strings.Builder
	if code, err := c.commit(context.Background(), &out, gitEnv); code != 0 || err != nil {
		t.Fatalf("the host's commit: %d %v\n%s", code, err, out.String())
	}
	for _, mark := range []string{"container-hook", "container-fsmonitor"} {
		if _, err := os.Stat(filepath.Join(marks, mark)); err == nil {
			t.Errorf("the host's commit ran the %s the container wrote", mark)
		}
	}
	if _, err := os.Stat(filepath.Join(marks, "main-hook")); err != nil {
		t.Errorf("the host's commit should run the main checkout's hook: %v", err)
	}
	if got := git(root, "log", "-1", "--format=%s", "story/"+story); got != "the story's work" {
		t.Errorf("the story's branch's last commit: %q", got)
	}
}

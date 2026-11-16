package closeout

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// quick is a wait short enough for a test, with flai serve's answer allowed
// a few seconds once the request is taken.
var quick = Wait{Take: 50 * time.Millisecond, Answer: 5 * time.Second, Poll: 5 * time.Millisecond}

// take stands in for flai serve: it waits for story's request in cacheDir,
// takes it by renaming it, and answers with what answer makes of it. The
// request it read arrives on the channel.
func take(t *testing.T, cacheDir, story string, answer func(Request) Answer) <-chan Request {
	t.Helper()
	got := make(chan Request, 1)
	dir := Dir(cacheDir)
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			taken := File(dir, story, TakenSuffix)
			if os.Rename(File(dir, story, RequestSuffix), taken) == nil {
				data, err := os.ReadFile(taken)
				var req Request
				if err == nil {
					err = json.Unmarshal(data, &req)
				}
				if err != nil {
					t.Errorf("the request taken: %v", err)
				}
				got <- req
				if err := WriteAnswer(dir, answer(req)); err != nil {
					t.Errorf("the answer: %v", err)
				}
				return
			}
			time.Sleep(2 * time.Millisecond)
		}
		close(got)
	}()
	return got
}

// ADR-0229 §1: a request flai serve takes is answered with the job's id, the
// request stamped with when it was asked, and the files left behind are
// cleared.
func TestAskAnswersTheJobFlaiServeStarts(t *testing.T) {
	t.Parallel()
	cache := t.TempDir()
	got := take(t, cache, "S-0001", func(r Request) Answer {
		return Answer{Job: r.Story + "-20261031T090000Z", Story: r.Story, Started: true, Log: LogPath(cache, r.Story)}
	})
	a, err := Ask(context.Background(), cache, Request{Story: "S-0001", By: "agent-S-0001", Message: "close out", Fresh: true}, quick)
	if err != nil {
		t.Fatal(err)
	}
	if a.Job != "S-0001-20261031T090000Z" || !a.Started || a.Running() || a.Refused() || a.Log != LogPath(cache, "S-0001") {
		t.Errorf("the answer: %+v", a)
	}
	req := <-got
	if req.Story != "S-0001" || req.By != "agent-S-0001" || req.Message != "close out" || !req.Fresh || req.At.IsZero() {
		t.Errorf("the request flai serve read: %+v", req)
	}
	for _, suffix := range []string{RequestSuffix, TakenSuffix, AnswerSuffix} {
		if _, err := os.Stat(File(Dir(cache), "S-0001", suffix)); !os.IsNotExist(err) {
			t.Errorf("the %s file is left behind: %v", suffix, err)
		}
	}
}

// An answer naming a job flai serve did not start says it was running; one
// with a reason is a refusal.
func TestAnAnswerSaysWhetherItsJobWasRunningOrItWasRefused(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		a                Answer
		running, refused bool
	}{
		{Answer{Job: "S-0001-x", Started: true}, false, false},
		{Answer{Job: "S-0001-x"}, true, false},
		{Answer{Recorded: true, Verdict: VerdictPassed}, false, false},
		{Answer{Reason: "no worktree"}, false, true},
	} {
		if c.a.Running() != c.running || c.a.Refused() != c.refused {
			t.Errorf("%+v: running %v refused %v, want %v %v", c.a, c.a.Running(), c.a.Refused(), c.running, c.refused)
		}
	}
}

// ADR-0229 §1: a request no flai serve takes in time is withdrawn, and the
// error is ErrServeNotTaking, saying to run the close-out yourself.
func TestAskNotTakenIsWithdrawn(t *testing.T) {
	t.Parallel()
	cache := t.TempDir()
	_, err := Ask(context.Background(), cache, Request{Story: "S-0001", By: "agent-S-0001"}, quick)
	if !errors.Is(err, ErrServeNotTaking) || !strings.Contains(err.Error(), "run the close-out script") {
		t.Fatalf("a request no one took: %v", err)
	}
	if _, err := os.Stat(File(Dir(cache), "S-0001", RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// A context that ends before flai serve takes the request withdraws it.
func TestAskStoppedBeforeItIsTakenWithdrawsTheRequest(t *testing.T) {
	t.Parallel()
	cache := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Ask(ctx, cache, Request{Story: "S-0001"}, Wait{Take: time.Minute, Poll: 5 * time.Millisecond})
	if !errors.Is(err, context.Canceled) || !strings.Contains(err.Error(), "withdrew it") {
		t.Fatalf("a request stopped: %v", err)
	}
	if _, err := os.Stat(File(Dir(cache), "S-0001", RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// A job's state reads back as written; a story with none reads false, and
// one whose state is not JSON is an error saying what to do.
func TestReadJobReadsWhatWriteJobWrote(t *testing.T) {
	t.Parallel()
	cache := t.TempDir()
	if _, ok, err := ReadJob(cache, "S-0001"); ok || err != nil {
		t.Fatalf("no job yet: %v %v", ok, err)
	}
	dir := Dir(cache)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	started := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	if err := WriteJob(dir, Job{ID: "S-0001-20261031T090000Z", Story: "S-0001", By: "agent-S-0001", Started: started}); err != nil {
		t.Fatal(err)
	}
	j, ok, err := ReadJob(cache, "S-0001")
	if err != nil || !ok || !j.Running() || j.ID != "S-0001-20261031T090000Z" || !j.Started.Equal(started) {
		t.Fatalf("the running job: %+v %v %v", j, ok, err)
	}
	code := 0
	j.Ended, j.ExitCode, j.Verdict, j.Commit = started.Add(time.Minute), &code, VerdictPassed, "abc"
	if err := WriteJob(dir, j); err != nil {
		t.Fatal(err)
	}
	if j, _, err := ReadJob(cache, "S-0001"); err != nil || j.Running() || j.Verdict != VerdictPassed || j.ExitCode == nil || *j.ExitCode != 0 {
		t.Errorf("the ended job: %+v %v", j, err)
	}
	if err := os.WriteFile(File(dir, "S-0002", JobSuffix), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := ReadJob(cache, "S-0002"); ok || err == nil || !strings.Contains(err.Error(), "delete it and ask for the close-out again") {
		t.Errorf("a job that does not read: %v %v", ok, err)
	}
}

// closeOutLog is the log of S-0001's close-out script so far: a round whose
// flai verify passed, printing its steps at its end, a commit, and a sync
// check that stopped at sync, then a second round whose verify has begun.
var closeOutLog = []string{
	"close-out: flai verify S-0001 --record-issues",
	"passed rebase (1s)",
	"passed sync (2s)",
	"passed narrative (0s)",
	"passed check (3s)",
	"passed tier go-test (40s)",
	"verify: S-0001 passed every step",
	"close-out: commit",
	" flai/x.go | 2 +-",
	"close-out: flai verify S-0001 --sync-only",
	"passed rebase (1s)",
	"failed sync (1s)",
	"verify: S-0001 stopped at sync (failed)",
	"close-out: main moved during the run; flai stream sync S-0001 (round 2 of 3)",
	"close-out: flai verify S-0001 --record-issues",
}

// I-0384: a running job answers its id, start, the head it runs at, and the
// step its log says it is on; an ended job answers false; a job written
// before S-0799, with no head, still reads.
func TestRunningSaysTheHeadAndTheStepOfARunningJob(t *testing.T) {
	t.Parallel()
	cache := t.TempDir()
	dir := Dir(cache)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := Running(cache, "S-0001"); ok || err != nil {
		t.Fatalf("no job yet: %v %v", ok, err)
	}
	started := time.Date(2026, 11, 16, 2, 0, 0, 0, time.UTC)
	job := Job{ID: "S-0001-20261116T020000Z", Story: "S-0001", By: "agent-S-0001", Started: started, Head: "head-b"}
	if err := WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	p, ok, err := Running(cache, "S-0001")
	if err != nil || !ok || p != (Progress{ID: job.ID, Started: started, Head: "head-b"}) {
		t.Fatalf("a running job with no log yet: %+v %v %v", p, ok, err)
	}
	data, err := json.Marshal(p)
	if err != nil || string(data) != `{"id":"S-0001-20261116T020000Z","started":"2026-11-16T02:00:00Z","head":"head-b"}` {
		t.Errorf("the progress as JSON: %s %v", data, err)
	}
	for _, c := range []struct {
		lines int
		want  string
	}{
		{0, ""},
		{1, StepVerify},
		{6, StepVerify}, // verify's own steps are no step of the close-out's
		{7, StepCommit}, // verify passed: the commit is next
		{9, StepCommit}, // what the commit printed
		{10, StepSyncCheck},
		{13, StepSyncCheck},
		{14, StepSync},
		{15, StepVerify}, // the next round
	} {
		log := strings.Join(closeOutLog[:c.lines], "\n")
		if c.lines > 0 {
			log += "\n"
		}
		if err := os.WriteFile(LogPath(cache, "S-0001"), []byte(log), 0o644); err != nil {
			t.Fatal(err)
		}
		if p, ok, err := Running(cache, "S-0001"); err != nil || !ok || p.Step != c.want || p.Head != "head-b" {
			t.Errorf("after %d lines of the log: %+v %v %v, want the step %q", c.lines, p, ok, err, c.want)
		}
	}
	// a verify that stopped leaves the close-out at it until its last line
	stopped := "close-out: flai verify S-0001 --record-issues\nfailed rebase (0s)\nverify: S-0001 stopped at rebase (failed)\n"
	if err := os.WriteFile(LogPath(cache, "S-0001"), []byte(stopped), 0o644); err != nil {
		t.Fatal(err)
	}
	if p, _, err := Running(cache, "S-0001"); err != nil || p.Step != StepVerify {
		t.Errorf("a verify that stopped: %+v %v", p, err)
	}
	// another story's lines are not this one's
	if err := os.WriteFile(LogPath(cache, "S-0001"), []byte("close-out: flai verify S-0002 --sync-only\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if p, _, err := Running(cache, "S-0001"); err != nil || p.Step != "" {
		t.Errorf("another story's line: %+v %v", p, err)
	}

	before := `{"id":"S-0001-20261116T020000Z","story":"S-0001","by":"agent-S-0001","started":"2026-11-16T02:00:00Z"}`
	if err := os.WriteFile(File(dir, "S-0001", JobSuffix), []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(LogPath(cache, "S-0001"), []byte(closeOutLog[0]+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if p, ok, err := Running(cache, "S-0001"); err != nil || !ok || p != (Progress{ID: job.ID, Started: started, Step: StepVerify}) {
		t.Errorf("a job written before S-0799: %+v %v %v", p, ok, err)
	}

	code := 1
	job.Ended, job.ExitCode, job.Verdict, job.Commit = started.Add(time.Minute), &code, VerdictFailed, "head-c"
	if err := WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	if p, ok, err := Running(cache, "S-0001"); ok || err != nil || p != (Progress{}) {
		t.Errorf("an ended job: %+v %v %v", p, ok, err)
	}
	if j, _, err := ReadJob(cache, "S-0001"); err != nil || j.Head != "head-b" || j.Commit != "head-c" {
		t.Errorf("the ended job keeps the head it ran at beside the one it ended at: %+v %v", j, err)
	}
}

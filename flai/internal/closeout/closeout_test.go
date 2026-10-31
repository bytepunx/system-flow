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

package kluster

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

// quick is a wait short enough for a test, with flai serve's answer allowed
// a few seconds once the request is taken.
var quick = Wait{Take: 50 * time.Millisecond, Answer: 5 * time.Second, Poll: 5 * time.Millisecond}

// serve stands in for flai serve: it looks for pending requests in root's
// requests folder, takes the first with Take, and answers with what answer
// makes of it. The request it took arrives on the channel.
func serve(t *testing.T, root string, answer func(Request) Answer) <-chan Request {
	t.Helper()
	got := make(chan Request, 1)
	go func() {
		defer close(got)
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			pending, err := Pending(root)
			if err != nil {
				t.Errorf("the pending requests: %v", err)
				return
			}
			for _, p := range pending {
				req, ok, err := Take(root, p.Story, p.Kind)
				if err != nil {
					t.Errorf("take %s's %s: %v", p.Story, p.Kind, err)
					return
				}
				if !ok {
					continue
				}
				got <- req
				if err := WriteAnswer(root, answer(req)); err != nil {
					t.Errorf("the answer: %v", err)
				}
				return
			}
			time.Sleep(2 * time.Millisecond)
		}
	}()
	return got
}

// ADR-0292 §1: a request flai serve takes is answered with the job's id and
// the jobs ahead of it, the request stamped with when it was asked, and the
// files left behind are cleared.
func TestAskAnswersTheJobFlaiServeQueues(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	queued := time.Date(2026, 11, 13, 7, 33, 53, 0, time.UTC)
	got := serve(t, root, func(r Request) Answer {
		return Answer{Job: JobID(r.Story, r.Kind, queued), Kind: r.Kind, Story: r.Story, Ahead: 2, Started: true}
	})
	a, err := Ask(context.Background(), root, Request{Story: "s-0660", Kind: KindDeploy, By: "agent-S-0660"}, quick)
	if err != nil {
		t.Fatal(err)
	}
	if a.Job != "S-0660-deploy-20261113T073353Z" || a.Ahead != 2 || !a.Started || a.Found() || a.Refused() || a.Kind != KindDeploy {
		t.Errorf("the answer: %+v", a)
	}
	req := <-got
	if req.Story != "S-0660" || req.Kind != KindDeploy || req.By != "agent-S-0660" || req.At.IsZero() {
		t.Errorf("the request flai serve took: %+v", req)
	}
	for _, suffix := range []string{RequestSuffix, TakenSuffix, AnswerSuffix} {
		if _, err := os.Stat(RequestFile(root, "S-0660", KindDeploy, suffix)); !os.IsNotExist(err) {
			t.Errorf("the %s file is left behind: %v", suffix, err)
		}
	}
}

// ADR-0292 §1: a request while a job of the story waits or runs is answered
// with that job, which it did not queue.
func TestAskAnswersAJobAlreadyQueued(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	serve(t, root, func(r Request) Answer {
		return Answer{Job: "S-0660-deploy-20261113T070000Z", Kind: r.Kind, Story: r.Story, Ahead: 1}
	})
	a, err := Ask(context.Background(), root, Request{Story: "S-0660", Kind: KindTeardown, By: "agent-S-0660"}, quick)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Found() || a.Started || a.Refused() || a.Job != "S-0660-deploy-20261113T070000Z" || a.Ahead != 1 {
		t.Errorf("an answer naming the job waiting: %+v", a)
	}
}

// ADR-0292 §1: a refusal queues nothing and says why and what to do.
func TestAskAnswersARefusal(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	reason := "the kluster host action is off; the operator turns it on with flai serve enable kluster"
	serve(t, root, func(r Request) Answer { return Answer{Kind: r.Kind, Story: r.Story, Reason: reason} })
	a, err := Ask(context.Background(), root, Request{Story: "S-0660", Kind: KindDeploy, By: "agent-S-0660"}, quick)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Refused() || a.Found() || a.Job != "" || a.Reason != reason {
		t.Errorf("a refusal: %+v", a)
	}
}

// ADR-0292 §1: a request no flai serve takes in time is withdrawn, and the
// error is ErrServeNotTaking, saying the operator runs the job in a shell.
func TestAskNotTakenIsWithdrawn(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	_, err := Ask(context.Background(), root, Request{Story: "S-0660", Kind: KindDeploy, By: "agent-S-0660"}, quick)
	if !errors.Is(err, ErrServeNotTaking) || !strings.Contains(err.Error(), "the operator runs flai kluster deploy") {
		t.Fatalf("a request no one took: %v", err)
	}
	if _, err := os.Stat(RequestFile(root, "S-0660", KindDeploy, RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
	if pending, err := Pending(root); err != nil || len(pending) != 0 {
		t.Errorf("nothing is pending once it is withdrawn: %v %v", pending, err)
	}
}

// A request for no story, or for a kind flai serve does not run, writes
// nothing.
func TestAskRefusesARequestItCannotName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, req := range []Request{{Story: "T-0001", Kind: KindDeploy}, {Story: "S-0660", Kind: "promote"}, {Story: "S-0660/../x", Kind: KindDeploy}} {
		if _, err := Ask(context.Background(), root, req, quick); err == nil || errors.Is(err, ErrServeNotTaking) {
			t.Errorf("%+v: %v", req, err)
		}
	}
	if _, err := os.Stat(RequestsDir(root)); !os.IsNotExist(err) {
		t.Errorf("nothing is written: %v", err)
	}
}

// Pending lists the requests oldest first; Take renames one to its taken
// file, and a request withdrawn first is not taken.
func TestPendingAndTake(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if pending, err := Pending(root); err != nil || len(pending) != 0 {
		t.Fatalf("no folder, nothing pending: %v %v", pending, err)
	}
	if err := os.MkdirAll(RequestsDir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 11, 13, 7, 0, 0, 0, time.UTC)
	for i, story := range []string{"S-0002", "S-0001"} {
		data := `{"story":"` + story + `","kind":"deploy","by":"agent-` + story + `","at":"` + at.Add(time.Duration(i)*time.Second).Format(time.RFC3339) + `"}`
		if err := os.WriteFile(RequestFile(root, story, KindDeploy, RequestSuffix), []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	pending, err := Pending(root)
	if err != nil || len(pending) != 2 || pending[0].Story != "S-0002" || pending[1].Story != "S-0001" || pending[0].By != "agent-S-0002" {
		t.Fatalf("pending, oldest first: %+v %v", pending, err)
	}
	req, ok, err := Take(root, "S-0002", KindDeploy)
	if err != nil || !ok || req.Story != "S-0002" || !req.At.Equal(at) {
		t.Fatalf("take: %+v %v %v", req, ok, err)
	}
	if _, err := os.Stat(RequestFile(root, "S-0002", KindDeploy, TakenSuffix)); err != nil {
		t.Errorf("the taken file: %v", err)
	}
	if _, ok, err := Take(root, "S-0002", KindDeploy); ok || err != nil {
		t.Errorf("a request taken already is not taken again: %v %v", ok, err)
	}
}

// A job's state reads back as written, waiting, running, then ended; a story
// with none reads false, and one whose state is not JSON is an error saying
// what to do.
func TestReadAskedJobReadsWhatWriteAskedJobWrote(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, ok, err := ReadAskedJob(root, "S-0660"); ok || err != nil {
		t.Fatalf("no job yet: %v %v", ok, err)
	}
	queued := time.Date(2026, 11, 13, 7, 0, 0, 0, time.UTC)
	j := AskedJob{ID: JobID("S-0660", KindDeploy, queued), Kind: KindDeploy, Story: "S-0660", By: "agent-S-0660", QueuedAt: queued}
	if err := WriteAskedJob(root, j); err != nil {
		t.Fatal(err)
	}
	got, ok, err := ReadAskedJob(root, "S-0660")
	if err != nil || !ok || !got.Waiting() || got.Running() || got.Ended() || got.By != "agent-S-0660" || !got.QueuedAt.Equal(queued) {
		t.Fatalf("the job waiting: %+v %v %v", got, ok, err)
	}
	j.StartedAt = queued.Add(time.Minute)
	if err := WriteAskedJob(root, j); err != nil {
		t.Fatal(err)
	}
	if got, _, _ := ReadAskedJob(root, "S-0660"); got.Waiting() || !got.Running() || got.Ended() {
		t.Errorf("the job running: %+v", got)
	}
	code := 3
	j.EndedAt, j.Commit, j.Outcome, j.FailedStep, j.ExitCode, j.Log = queued.Add(2*time.Minute), commitA, OutcomeFailed, StepBuild, &code, LogPath(root, "s-0660", StepBuild)
	j.Steps = []Step{{Name: StepBuild, Outcome: OutcomeFailed, ExitCode: 3, Log: LogPath(root, "s-0660", StepBuild)}, {Name: StepLoad, Outcome: OutcomeSkipped}, {Name: StepDeploy, Outcome: OutcomeSkipped}}
	if err := WriteAskedJob(root, j); err != nil {
		t.Fatal(err)
	}
	got, _, err = ReadAskedJob(root, "S-0660")
	if err != nil || !got.Ended() || got.Running() || got.Outcome != OutcomeFailed || got.FailedStep != StepBuild || got.ExitCode == nil || *got.ExitCode != 3 || len(got.Steps) != 3 || got.Steps[0].Log != j.Log {
		t.Errorf("the job ended: %+v %v", got, err)
	}
	if err := os.WriteFile(JobFile(root, "S-0661"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := ReadAskedJob(root, "S-0661"); ok || err == nil || !strings.Contains(err.Error(), "delete it and ask for the job again") {
		t.Errorf("a job that does not read: %v %v", ok, err)
	}
}

// ADR-0292 §5: Await answers the job once it ends within the wait; one that
// outlasts it is answered as it stands, not ended; a wait of 0 answers at
// once; and a job another took the place of is an error naming it.
func TestAwaitWaitsForTheJobToEnd(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	queued := time.Now().UTC()
	j := AskedJob{ID: JobID("S-0660", KindDeploy, queued), Kind: KindDeploy, Story: "S-0660", By: "agent-S-0660", QueuedAt: queued, StartedAt: queued}
	if err := WriteAskedJob(root, j); err != nil {
		t.Fatal(err)
	}
	got, ended, err := Await(context.Background(), root, "S-0660", j.ID, 0, 0)
	if err != nil || ended || !got.Running() {
		t.Errorf("a wait of 0 answers the job running at once: %+v %v %v", got, ended, err)
	}
	got, ended, err = Await(context.Background(), root, "S-0660", j.ID, 30*time.Millisecond, 5*time.Millisecond)
	if err != nil || ended || !got.Running() {
		t.Errorf("a job that outlasts the wait: %+v %v %v", got, ended, err)
	}
	go func() {
		time.Sleep(20 * time.Millisecond)
		done := j
		done.EndedAt, done.Outcome, done.Commit = time.Now().UTC(), OutcomePassed, commitA
		if err := WriteAskedJob(root, done); err != nil {
			t.Errorf("end the job: %v", err)
		}
	}()
	got, ended, err = Await(context.Background(), root, "S-0660", j.ID, 5*time.Second, 5*time.Millisecond)
	if err != nil || !ended || got.Outcome != OutcomePassed || got.Commit != commitA {
		t.Errorf("a job that ends within the wait: %+v %v %v", got, ended, err)
	}
	if _, _, err := Await(context.Background(), root, "S-0660", "S-0660-deploy-older", time.Second, 5*time.Millisecond); err == nil || !strings.Contains(err.Error(), "await "+j.ID) {
		t.Errorf("a job another took the place of: %v", err)
	}
	if _, _, err := Await(context.Background(), root, "S-0661", "S-0661-deploy-x", time.Second, 5*time.Millisecond); err == nil {
		t.Error("a story with no job is an error")
	}
}

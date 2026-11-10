package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/closeout"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// closeOutFixture is setup with the story's worktree there.
func closeOutFixture(t *testing.T) *fixture {
	t.Helper()
	f := setup(t)
	if err := os.MkdirAll(f.repo.WorktreePath(f.story.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	return f
}

// verifyFixtureWith is testFixture, with the server's options changed by
// with when it is not nil, with the story's worktree on its own branch from
// main and its narrative's Current state and Next steps written; change
// commits a file on the story's branch, and git runs git.
func verifyFixtureWith(t *testing.T, with func(*Options)) (f *fixture, change func(rel string), git func(dir string, args ...string)) {
	t.Helper()
	f, git = testFixtureWith(t, with)
	wt := f.repo.WorktreePath(f.story.ID)
	git(f.repo.Root, "worktree", "add", "-q", "-b", "story/"+f.story.ID, wt, "main")
	path := f.repo.NarrativePath(f.story.ID)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	narrative := string(data)
	for _, h := range []string{"## Current state\n", "## Next steps\n"} {
		if !strings.Contains(narrative, h) {
			t.Fatalf("the narrative has no %q: %s", h, narrative)
		}
		narrative = strings.Replace(narrative, h, h+"\nWritten.\n", 1)
	}
	if err := os.WriteFile(path, []byte(narrative), 0o644); err != nil {
		t.Fatal(err)
	}
	return f, func(rel string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(wt, filepath.FromSlash(rel)), []byte("changed\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(wt, "add", "-A")
		git(wt, "commit", "-q", "-m", "change "+rel)
	}, git
}

// takeCloseOut stands in for flai serve: it waits for story's close-out
// request in the project's cache dir, takes it, and answers with what answer
// makes of it. The request it read arrives on the channel, which is closed
// when none came within five seconds.
func takeCloseOut(t *testing.T, cacheDir, story string, answer func(closeout.Request) closeout.Answer) <-chan closeout.Request {
	t.Helper()
	got := make(chan closeout.Request, 1)
	dir := closeout.Dir(cacheDir)
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			taken := closeout.File(dir, story, closeout.TakenSuffix)
			if os.Rename(closeout.File(dir, story, closeout.RequestSuffix), taken) == nil {
				var req closeout.Request
				data, err := os.ReadFile(taken)
				if err == nil {
					err = json.Unmarshal(data, &req)
				}
				if err != nil {
					t.Errorf("the request taken: %v", err)
				}
				got <- req
				if err := closeout.WriteAnswer(dir, answer(req)); err != nil {
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

// S-0546 criterion 1, ADR-0229 §1: the tool is listed, saying that it asks
// flai serve for the close-out as a job and what to wait for; it asks as the
// session's agent, with the message, all, and fresh, and answers the job's
// id at once, or the id of the story's job already running.
func TestVerifyAsksFlaiServeForTheCloseOutJob(t *testing.T) {
	t.Parallel()
	f := closeOutFixture(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	listed := false
	for _, tool := range res.Tools {
		if tool.Name != "verify" {
			continue
		}
		listed = true
		for _, want := range []string{"Ask flai serve for a story's close-out", "as a job", "scripts/close-out.sh S-nnnn -m <message>", "answers at once",
			"a commit the verify record holds a passing verdict for", "hold wait_for_events until it reports verified", "flai verify S-nnnn --last",
			"flai serve is not running", "run the close-out script, or flai verify where the project has none, yourself", "acceptance criteria and the conventions"} {
			if !strings.Contains(tool.Description, want) {
				t.Errorf("the description does not say %q: %s", want, tool.Description)
			}
		}
	}
	if !listed {
		t.Fatal("the tool verify is not listed")
	}

	cache := f.repo.CacheDir()
	id := f.story.ID + "-20261031T090000Z"
	got := takeCloseOut(t, cache, f.story.ID, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Job: id, Story: r.Story, Started: true, Log: closeout.LogPath(cache, r.Story)}
	})
	out, failed := f.call(t, "verify", map[string]any{"story": f.story.ID, "message": "feat: close out", "all": true, "fresh": true})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["job"] != id || out["story"] != f.story.ID || out["started"] != true || out["running"] != nil || out["recorded"] != nil || out["log"] != closeout.LogPath(cache, f.story.ID) {
		t.Errorf("the job started: %v", out)
	}
	if req := <-got; req.Story != f.story.ID || req.By != "claude" || req.Message != "feat: close out" || !req.All || !req.Fresh {
		t.Errorf("the request flai serve read: %+v", req)
	}

	takeCloseOut(t, cache, f.story.ID, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Job: id, Story: r.Story, Log: closeout.LogPath(cache, r.Story)}
	})
	out, failed = f.call(t, "verify", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["job"] != id || out["started"] != false || out["running"] != true {
		t.Errorf("the job already running: %v", out)
	}
}

// S-0546 criterion 2, ADR-0229 §5: a commit the verify record holds a
// verdict for is answered with it, saying it came from the record, with its
// commit and the time of its run, and no job.
func TestVerifyAnswersTheVerdictTheRecordHolds(t *testing.T) {
	t.Parallel()
	f := closeOutFixture(t)
	ranAt := time.Date(2026, 10, 31, 9, 0, 0, 0, time.UTC)
	takeCloseOut(t, f.repo.CacheDir(), f.story.ID, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Story: r.Story, Recorded: true, Verdict: closeout.VerdictPassed, Commit: "0123abc4567", RanAt: ranAt}
	})
	out, failed := f.call(t, "verify", map[string]any{"story": f.story.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	if out["recorded"] != true || out["verdict"] != closeout.VerdictPassed || out["commit"] != "0123abc4567" || out["ran_at"] != "2026-10-31T09:00:00Z" ||
		out["job"] != nil || out["started"] != false {
		t.Errorf("the answer from the record: %v", out)
	}
}

// A request flai serve refuses is an error with its reason.
func TestVerifyFailsWithTheReasonFlaiServeRefused(t *testing.T) {
	t.Parallel()
	f := closeOutFixture(t)
	takeCloseOut(t, f.repo.CacheDir(), f.story.ID, func(r closeout.Request) closeout.Answer {
		return closeout.Answer{Story: r.Story, Reason: "flai serve could not start " + r.Story + "'s close-out: the disk is full"}
	})
	if _, failed := f.call(t, "verify", map[string]any{"story": f.story.ID}); !strings.Contains(failed, "the disk is full") {
		t.Errorf("the refusal: %q", failed)
	}
}

// ADR-0229 §1: with no flai serve to take the request, it is withdrawn and
// the call fails, saying to run the close-out yourself.
func TestVerifySaysFlaiServeIsNotRunning(t *testing.T) {
	// it shortens the package's wait, so it runs before the parallel tests
	was := closeOutWait
	closeOutWait = closeout.Wait{Take: 50 * time.Millisecond, Poll: 5 * time.Millisecond}
	t.Cleanup(func() { closeOutWait = was })
	f := closeOutFixture(t)
	_, failed := f.call(t, "verify", map[string]any{"story": f.story.ID})
	if !strings.Contains(failed, "flai serve is not running") || !strings.Contains(failed, "run the close-out script") {
		t.Errorf("no flai serve: %q", failed)
	}
	if _, err := os.Stat(closeout.File(closeout.Dir(f.repo.CacheDir()), f.story.ID, closeout.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// S-0270: a story that is not there, an item that is not a story, and a
// story with no worktree are refused, saying what to do, and nothing is
// asked of flai serve.
func TestVerifyRefusesWhatItCannotVerify(t *testing.T) {
	f := closeOutFixture(t)
	other, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Other", Owner: "alex", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args map[string]any
		want string
	}{
		{map[string]any{"story": " "}, "name the story to verify"},
		{map[string]any{"story": "S-0099"}, "S-0099 is not a work item here"},
		{map[string]any{"story": f.task.ID}, f.task.ID + " is a task; name a story"},
		{map[string]any{"story": other.ID}, other.ID + " has no worktree at .flai-cache/worktrees/" + other.ID + "; open it with flai stream open " + other.ID},
	} {
		if _, failed := f.call(t, "verify", c.args); !strings.Contains(failed, c.want) {
			t.Errorf("%v: %q, want %q", c.args, failed, c.want)
		}
	}
	if _, err := os.Stat(closeout.Dir(f.repo.CacheDir())); !os.IsNotExist(err) {
		t.Errorf("a refusal asked flai serve: %v", err)
	}
}

// ADR-0229 §4: a story's agent with its question open to the designer is not
// told to end while a close-out job of its story runs, nor when the job's
// state cannot be read, since it waits for the job's verified event; once
// the job ended, it is.
func TestWaitForEventsHoldsAStorysAgentWhileItsCloseOutJobRuns(t *testing.T) {
	t.Parallel()
	f, d := storysAgent(t)
	f.ask(t, f.story.ID)
	dir := closeout.Dir(f.repo.CacheDir())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	job := closeout.Job{ID: f.story.ID + "-20261031T090000Z", Story: f.story.ID, By: "claude", Started: t0}
	if err := closeout.WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)

	path := closeout.File(dir, f.story.ID, closeout.JobSuffix)
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)

	code := 0
	job.Ended, job.ExitCode, job.Verdict = t0.Add(time.Minute), &code, closeout.VerdictPassed
	if err := closeout.WriteJob(dir, job); err != nil {
		t.Fatal(err)
	}
	f.endsAtOnce(t, d)
}

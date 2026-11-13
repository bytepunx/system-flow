package mcpserver

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
)

// klusterWaitForTests is the kluster tools' wait in these tests.
var klusterWaitForTests = kluster.Wait{Take: 3 * time.Second, Answer: 3 * time.Second, Poll: 5 * time.Millisecond}

// klusterFixture is setup for the agent of the fixture's story, S-0001, with
// the kluster tools' waits shortened and the options changed by with when it
// is not nil.
func klusterFixture(t *testing.T, with func(*Options)) *fixture {
	t.Helper()
	return setupWith(t, func(o *Options) {
		o.Story, o.KlusterWait = "S-0001", klusterWaitForTests
		if with != nil {
			with(o)
		}
	})
}

// takeKluster stands in for flai serve: it takes story's request of kind in
// the project at root, writes jobs[0] as the story's job state, answers with
// what answer makes of the request, and writes each later job state 30ms
// after the one before. The request arrives on the channel, which is closed
// when none came within five seconds.
func takeKluster(t *testing.T, root, story, kind string, answer func(kluster.Request) kluster.Answer, jobs ...kluster.AskedJob) <-chan kluster.Request {
	t.Helper()
	got := make(chan kluster.Request, 1)
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			req, ok, err := kluster.Take(root, story, kind)
			if err != nil {
				t.Errorf("take the request: %v", err)
				break
			}
			if !ok {
				time.Sleep(2 * time.Millisecond)
				continue
			}
			got <- req
			if len(jobs) > 0 {
				if err := kluster.WriteAskedJob(root, jobs[0]); err != nil {
					t.Errorf("the job state: %v", err)
				}
			}
			if err := kluster.WriteAnswer(root, answer(req)); err != nil {
				t.Errorf("the answer: %v", err)
			}
			for _, j := range jobs[min(1, len(jobs)):] {
				time.Sleep(30 * time.Millisecond)
				if err := kluster.WriteAskedJob(root, j); err != nil {
					t.Errorf("the job state: %v", err)
				}
			}
			return
		}
		close(got)
	}()
	return got
}

// klusterQueuedAt is when the jobs of these tests were queued.
var klusterQueuedAt = time.Date(2026, 11, 13, 7, 33, 53, 0, time.UTC)

// askedKlusterJob is S-0001's job of kind, running once started, and ended with
// outcome once outcome is not empty.
func askedKlusterJob(kind string, started bool, outcome string) kluster.AskedJob {
	j := kluster.AskedJob{ID: kluster.JobID("S-0001", kind, klusterQueuedAt), Kind: kind, Story: "S-0001", By: "claude", QueuedAt: klusterQueuedAt}
	if started {
		j.StartedAt = klusterQueuedAt.Add(time.Second)
	}
	if outcome != "" {
		j.EndedAt, j.Outcome, j.Commit = klusterQueuedAt.Add(time.Minute), outcome, "0123456789abcdef0123456789abcdef01234567"
	}
	return j
}

// S-0660 criterion 1, ADR-0292 §1, §5: the tools are listed, saying what
// they answer and what to wait for; they ask flai serve as the session's
// agent for its own story, by default or named in any padding, and answer
// the job's id and the jobs ahead, or the job of the story found running.
func TestKlusterToolsAskFlaiServeForTheOwnStory(t *testing.T) {
	t.Parallel()
	f := klusterFixture(t, nil)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	listed := 0
	for _, tool := range res.Tools {
		if tool.Name != "kluster_deploy" && tool.Name != "kluster_teardown" {
			continue
		}
		listed++
		for _, want := range []string{"Ask flai serve to", "as a job (ADR-0292)", "for its own story alone", "It answers at once", "ahead, the jobs of the queue before it",
			"wait_seconds, 0 to 45", "more than 45 is refused", "hold wait_for_events for the kluster event", "the failed step, its exit code, and the log",
			"for another story", "while the kluster host action is off", "with no kluster.deploy set", "over kluster.agent_namespaces", "uncommitted changes", "When flai serve is not running"} {
			if !strings.Contains(tool.Description, want) {
				t.Errorf("%s's description does not say %q: %s", tool.Name, want, tool.Description)
			}
		}
	}
	if listed != 2 {
		t.Fatalf("the tools kluster_deploy and kluster_teardown are not both listed: %d", listed)
	}

	root := projectRoot(f.repo)
	deploy := askedKlusterJob(kluster.KindDeploy, false, "")
	got := takeKluster(t, root, f.story.ID, kluster.KindDeploy, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: deploy.ID, Kind: r.Kind, Story: r.Story, Ahead: 2, Started: true, Log: "logs/s-0001"}
	}, deploy)
	out, failed := f.call(t, "kluster_deploy", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	ans, _ := out["answer"].(map[string]any)
	if ans["job"] != deploy.ID || ans["kind"] != kluster.KindDeploy || ans["story"] != f.story.ID || ans["ahead"] != float64(2) || ans["started"] != true || ans["log"] != "logs/s-0001" || out["job"] != nil || out["ended"] != nil {
		t.Errorf("the job queued: %v", out)
	}
	if req := <-got; req.Story != f.story.ID || req.Kind != kluster.KindDeploy || req.By != "claude" {
		t.Errorf("the request flai serve read: %+v", req)
	}

	teardown := askedKlusterJob(kluster.KindTeardown, true, "")
	takeKluster(t, root, f.story.ID, kluster.KindTeardown, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: teardown.ID, Kind: r.Kind, Story: r.Story}
	}, teardown)
	out, failed = f.call(t, "kluster_teardown", map[string]any{"story": "s-1"})
	if failed != "" {
		t.Fatal(failed)
	}
	if ans, _ := out["answer"].(map[string]any); ans["job"] != teardown.ID || ans["started"] != false {
		t.Errorf("the job found running: %v", out)
	}
}

// S-0660 criterion 3, ADR-0292 §1: another story than the session's own,
// FLAI_STORY or the story in an agent name of the form agent-S-nnnn, is
// refused before any request is written, as is a wait over 45 seconds; a
// session with no story of its own names one.
func TestKlusterToolsRefuseAnotherStory(t *testing.T) {
	t.Parallel()
	f := klusterFixture(t, nil)
	for _, c := range []struct {
		tool string
		args map[string]any
		want string
	}{
		{"kluster_deploy", map[string]any{"story": "S-0002"}, "this session works S-0001"},
		{"kluster_teardown", map[string]any{"story": "S-2"}, "nothing was asked for S-0002"},
		{"kluster_deploy", map[string]any{"wait_seconds": 46}, "wait_seconds 46 is not 0 to 45"},
		{"kluster_teardown", map[string]any{"wait_seconds": -1}, "wait_seconds -1 is not 0 to 45"},
	} {
		if _, failed := f.call(t, c.tool, c.args); !strings.Contains(failed, c.want) {
			t.Errorf("%s %v: %q, want %q", c.tool, c.args, failed, c.want)
		}
	}
	if _, err := os.Stat(kluster.RequestsDir(projectRoot(f.repo))); !os.IsNotExist(err) {
		t.Errorf("a refusal wrote a request: %v", err)
	}

	named := klusterFixture(t, func(o *Options) { o.Story, o.Agent = "", "agent-S-0001" })
	if _, failed := named.call(t, "kluster_deploy", map[string]any{"story": "S-0002"}); !strings.Contains(failed, "this session works S-0001") {
		t.Errorf("the story in the agent's name: %q", failed)
	}
	none := klusterFixture(t, func(o *Options) { o.Story = "" })
	if _, failed := none.call(t, "kluster_deploy", map[string]any{}); !strings.Contains(failed, "name the story whose kluster deploy to ask for") {
		t.Errorf("no story at all: %q", failed)
	}
}

// ADR-0292 §1: a refusal from flai serve is an error with its reason; a
// request no flai serve takes is withdrawn, and the call says so.
func TestKlusterToolsRefusedOrNotTaken(t *testing.T) {
	t.Parallel()
	f := klusterFixture(t, nil)
	reason := "S-0001 has uncommitted changes in its worktree: a.go; commit first"
	takeKluster(t, projectRoot(f.repo), f.story.ID, kluster.KindDeploy, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Kind: r.Kind, Story: r.Story, Reason: reason}
	})
	if _, failed := f.call(t, "kluster_deploy", map[string]any{}); failed != reason {
		t.Errorf("the refusal: %q", failed)
	}

	idle := klusterFixture(t, func(o *Options) {
		o.KlusterWait = kluster.Wait{Take: 50 * time.Millisecond, Poll: 5 * time.Millisecond}
	})
	if _, failed := idle.call(t, "kluster_teardown", map[string]any{}); !strings.Contains(failed, "flai serve is not running") {
		t.Errorf("no flai serve: %q", failed)
	}
	if _, err := os.Stat(kluster.RequestFile(projectRoot(idle.repo), idle.story.ID, kluster.KindTeardown, kluster.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// S-0660 criterion 2, ADR-0292 §5: with wait_seconds the tool waits for the
// job's end and answers its state, outcome, commit, step, exit code, and log,
// with ended true; a job that outlasts the wait is answered running, with
// ended false.
func TestKlusterToolsWaitForTheEnd(t *testing.T) {
	t.Parallel()
	f := klusterFixture(t, nil)
	root := projectRoot(f.repo)
	running := askedKlusterJob(kluster.KindDeploy, true, "")
	failedJob := askedKlusterJob(kluster.KindDeploy, true, kluster.OutcomeFailed)
	code := 4
	failedJob.FailedStep, failedJob.ExitCode, failedJob.Log = "deploy", &code, "logs/s-0001/deploy.log"
	failedJob.Steps = []kluster.Step{{Name: "deploy", Outcome: kluster.OutcomeFailed, ExitCode: 4, Log: "logs/s-0001/deploy.log"}}
	answer := func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: running.ID, Kind: r.Kind, Story: r.Story, Started: true}
	}
	takeKluster(t, root, f.story.ID, kluster.KindDeploy, answer, running, failedJob)
	out, failed := f.call(t, "kluster_deploy", map[string]any{"wait_seconds": 5})
	if failed != "" {
		t.Fatal(failed)
	}
	job, _ := out["job"].(map[string]any)
	steps, _ := job["steps"].([]any)
	if out["ended"] != true || job["id"] != running.ID || job["outcome"] != kluster.OutcomeFailed || job["failed_step"] != "deploy" || job["exit_code"] != float64(4) ||
		job["log"] != "logs/s-0001/deploy.log" || job["commit"] != failedJob.Commit || len(steps) != 1 {
		t.Errorf("the job waited for to its end: %v", out)
	}

	takeKluster(t, root, f.story.ID, kluster.KindDeploy, answer, running)
	out, failed = f.call(t, "kluster_deploy", map[string]any{"wait_seconds": 1})
	if failed != "" {
		t.Fatal(failed)
	}
	if job, _ := out["job"].(map[string]any); out["ended"] != false || job["id"] != running.ID || job["outcome"] != nil || job["started_at"] == nil {
		t.Errorf("a job that outlasts the wait: %v", out)
	}
}

// ADR-0292 §5: a story's agent with its question open to the designer is not
// told to end while a kluster job of its story waits or runs, nor when the
// job's state cannot be read, since it waits for the job's kluster event;
// once the job ended, it is.
func TestWaitForEventsHoldsAStorysAgentWhileItsKlusterJobRuns(t *testing.T) {
	t.Parallel()
	f, d := storysAgent(t)
	f.ask(t, f.story.ID)
	root := projectRoot(f.repo)
	job := askedKlusterJob(kluster.KindDeploy, false, "")
	if err := kluster.WriteAskedJob(root, job); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)

	job.StartedAt = klusterQueuedAt.Add(time.Second)
	if err := kluster.WriteAskedJob(root, job); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)

	if err := os.WriteFile(kluster.JobFile(root, f.story.ID), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.holds(t, d)

	if err := kluster.WriteAskedJob(root, askedKlusterJob(kluster.KindDeploy, true, kluster.OutcomePassed)); err != nil {
		t.Fatal(err)
	}
	f.endsAtOnce(t, d)
}

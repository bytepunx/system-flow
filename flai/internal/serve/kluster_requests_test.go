package serve

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// klAgent is the story's agent that asks in these tests.
const klAgent = "agent-S-0002"

// inProgress makes a story with a worktree folder, its branch at head, a
// criterion and a task, and moves it to in-progress.
func (lab *klusterLab) inProgress(title, head string) string {
	lab.t.Helper()
	epic, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic " + title, Owner: "alex", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	st, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: title, Parent: epic.ID, Owner: "alex", Now: time.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] deploys\n", 1)), 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.move(st.ID, workitem.Ready)
	lab.move(st.ID, workitem.InProgress)
	if _, err := lab.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: st.ID, Owner: "builder", Now: time.Now()}); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.MkdirAll(lab.repo.WorktreePath(st.ID), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	lab.git.set(st.ID, head)
	return st.ID
}

// tick moves the lab's clock a second on, so that two jobs asked for one
// after the other have ids of their own.
func (lab *klusterLab) tick() {
	lab.mu.Lock()
	defer lab.mu.Unlock()
	lab.clock = lab.clock.Add(time.Second)
}

// ask writes by's request for story's job of kind, has the pass take it as
// its every second does, and answers its answer.
func (lab *klusterLab) ask(story, kind, by string) kluster.Answer {
	lab.t.Helper()
	lab.tick()
	data, err := json.Marshal(kluster.Request{Story: story, Kind: kind, By: by, At: lab.opts.Now()})
	if err != nil {
		lab.t.Fatal(err)
	}
	if err := os.MkdirAll(kluster.RequestsDir(lab.root), 0o755); err != nil {
		lab.t.Fatal(err)
	}
	if err := os.WriteFile(kluster.RequestFile(lab.root, story, kind, kluster.RequestSuffix), data, 0o644); err != nil {
		lab.t.Fatal(err)
	}
	lab.pass.takeRequests(context.Background())
	path := kluster.RequestFile(lab.root, story, kind, kluster.AnswerSuffix)
	data, err = os.ReadFile(path)
	if err != nil {
		lab.t.Fatalf("no answer to %s's %s request: %v", story, kind, err)
	}
	var a kluster.Answer
	if err := json.Unmarshal(data, &a); err != nil {
		lab.t.Fatal(err)
	}
	_ = os.Remove(path)
	_ = os.Remove(kluster.RequestFile(lab.root, story, kind, kluster.TakenSuffix))
	return a
}

// drain runs the jobs queued, in this goroutine.
func (lab *klusterLab) drain() { lab.pass.drain(context.Background()) }

// asked is story's job state, failing the test when it has none.
func (lab *klusterLab) asked(story string) kluster.AskedJob {
	lab.t.Helper()
	j, ok, err := kluster.ReadAskedJob(lab.root, story)
	if err != nil || !ok {
		lab.t.Fatalf("%s has no kluster job state: %v", story, err)
	}
	return j
}

// klusterEvents are the kluster events in the project's event log.
func (lab *klusterLab) klusterEvents() []inbox.LoggedEvent {
	var out []inbox.LoggedEvent
	for _, e := range inbox.LoggedEvents(lab.repo) {
		if e.Kind == inbox.KindKluster {
			out = append(out, e)
		}
	}
	return out
}

// queued is the number of jobs waiting in the pass's queue.
func (lab *klusterLab) queued() int {
	lab.pass.mu.Lock()
	defer lab.pass.mu.Unlock()
	return len(lab.pass.queue)
}

// ADR-0292 §1, §3, §6: a story in progress is deployed on its agent's
// request at its branch's head, the cluster brought up first, its record and
// the journal naming who asked; its job state is kept at queue and at end, a
// kluster event tells its end, and its namespace stays while it is in
// progress.
func TestAStoryInProgressIsDeployedOnRequestAtItsHead(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	lab.kluster.absent = true
	story := lab.inProgress("Asked", klHeadA)

	a := lab.ask(story, kluster.KindDeploy, klAgent)
	if a.Refused() || !a.Started || a.Job == "" || a.Ahead != 0 || a.Story != story || a.Kind != kluster.KindDeploy {
		t.Fatalf("answer = %+v, want a job started with none ahead", a)
	}
	if j := lab.asked(story); j.ID != a.Job || !j.Waiting() || j.By != klAgent || j.Kind != kluster.KindDeploy {
		t.Fatalf("the job's state at queue = %+v, want job %s waiting, asked by %s", j, a.Job, klAgent)
	}

	lab.drain()
	want := "absent;up;deploy " + story + " " + klHeadA
	if got := strings.Join(lab.kluster.all(), ";"); got != want {
		t.Fatalf("calls = %s, want %s", got, want)
	}
	ns, _ := kluster.StoryNamespace(story)
	if r, ok, _ := kluster.ReadRecord(lab.root, ns); !ok || r.AskedBy != klAgent || r.Commit != klHeadA {
		t.Fatalf("the namespace's record = %+v, want it deployed at %s asked by %s", r, short(klHeadA), klAgent)
	}
	j := lab.asked(story)
	if !j.Ended() || j.StartedAt.IsZero() || j.Outcome != kluster.OutcomePassed || j.Commit != klHeadA || j.ID != a.Job {
		t.Fatalf("the job's state at its end = %+v, want it passed at %s", j, short(klHeadA))
	}
	evs := lab.klusterEvents()
	if len(evs) != 1 || evs[0].ID != story || evs[0].Fields[inbox.FieldJob] != a.Job || evs[0].Fields[inbox.FieldOutcome] != kluster.OutcomePassed ||
		evs[0].Fields[inbox.FieldKind] != kluster.KindDeploy || evs[0].Fields[inbox.FieldCommit] != klHeadA {
		t.Fatalf("kluster events = %+v, want one telling the deploy passed", evs)
	}
	for _, e := range lab.entries() {
		if e.By != klAgent || !strings.Contains(e.Detail, "asked by "+klAgent) {
			t.Errorf("journal entry %+v does not name %s as who asked", e, klAgent)
		}
	}

	n := len(lab.kluster.all())
	lab.round()
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("looks at the story in progress ran %v, want its namespace kept and no deploy", got)
	}
}

// ADR-0292 §1, §4: each request flai serve cannot run is refused at once,
// saying why and what to do, and queues nothing.
func TestAKlusterRequestIsRefusedSayingWhy(t *testing.T) {
	for _, c := range []struct {
		name, block, kind string
		set               func(lab *klusterLab, story string) string // the story to ask for
		want              []string
	}{
		{"action off", klusterBlock, kluster.KindDeploy, func(lab *klusterLab, s string) string { lab.setOn(false); return s },
			[]string{"host action is off", "flai serve enable kluster", "operator's"}},
		{"no block", "", kluster.KindDeploy, nil, []string{"no kluster: block", "flai serve enable kluster", "operator's"}},
		{"no deploy", klusterNoDeploy, kluster.KindDeploy, nil, []string{"not set up yet"}},
		{"deploys before review off", klusterBlock + "  agent_namespaces: 0\n", kluster.KindDeploy, nil, []string{"deploys before review are off", "agent_namespaces is 0"}},
		{"not in progress", klusterBlock, kluster.KindDeploy, func(lab *klusterLab, s string) string { lab.move(s, workitem.Ready); return s },
			[]string{"is ready", "only while the story is in progress"}},
		{"not found", klusterBlock, kluster.KindDeploy, func(*klusterLab, string) string { return "S-0999" }, []string{"S-0999 is not a story"}},
		{"no worktree", klusterBlock, kluster.KindDeploy, func(lab *klusterLab, s string) string {
			_ = os.RemoveAll(lab.repo.WorktreePath(s))
			return s
		}, []string{"has no worktree", "flai stream open"}},
		{"uncommitted changes", klusterBlock, kluster.KindDeploy, func(lab *klusterLab, s string) string {
			lab.git.setDirty(lab.repo.WorktreePath(s), " M flai/cmd/x.go\nM  docs/y.md\n")
			return s
		}, []string{"uncommitted changes", "flai/cmd/x.go", "docs/y.md", "Commit them"}},
		{"teardown with no record", klusterBlock, kluster.KindTeardown, nil, []string{"has no record", "nothing to tear down"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			lab := newKlusterLab(t, c.block)
			story := lab.inProgress("Refused", klHeadA)
			if c.set != nil {
				story = c.set(lab, story)
			}
			a := lab.ask(story, c.kind, klAgent)
			if !a.Refused() || a.Job != "" || a.Started {
				t.Fatalf("answer = %+v, want a refusal", a)
			}
			for _, w := range c.want {
				if !strings.Contains(a.Reason, w) {
					t.Errorf("the refusal does not say %q:\n%s", w, a.Reason)
				}
			}
			if n := lab.queued(); n != 0 {
				t.Errorf("a refused request queued %d jobs", n)
			}
			if _, ok, _ := kluster.ReadAskedJob(lab.root, story); ok {
				t.Errorf("a refused request wrote a job state")
			}
			lab.drain()
			if got := lab.kluster.all(); len(got) != 0 {
				t.Errorf("a refused request ran %v", got)
			}
		})
	}
}

// ADR-0292 §1: one job per story: a request while its job waits is answered
// with that job, and a deploy at a newer head takes the waiting one's place;
// the answer counts the jobs ahead.
func TestARequestWhileAJobOfTheStoryWaitsIsAnsweredWithIt(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	first := lab.inProgress("First", klHeadA)
	story := lab.inProgress("Second", klHeadB)

	if a := lab.ask(first, kluster.KindDeploy, "agent-"+first); !a.Started || a.Ahead != 0 {
		t.Fatalf("first answer = %+v", a)
	}
	a := lab.ask(story, kluster.KindDeploy, klAgent)
	if !a.Started || a.Ahead != 1 {
		t.Fatalf("answer behind another story's job = %+v, want started with 1 ahead", a)
	}
	again := lab.ask(story, kluster.KindDeploy, klAgent)
	if again.Started || again.Job != a.Job || again.Ahead != 1 || again.Refused() {
		t.Fatalf("answer while the job waits = %+v, want job %s found, not started", again, a.Job)
	}
	if td := lab.ask(story, kluster.KindTeardown, klAgent); td.Started || td.Job != a.Job {
		t.Fatalf("a teardown request while the deploy waits = %+v, want the deploy's job %s", td, a.Job)
	}

	lab.git.set(story, klHeadC)
	newer := lab.ask(story, kluster.KindDeploy, klAgent)
	if !newer.Started || newer.Job == a.Job || newer.Ahead != 1 {
		t.Fatalf("answer at a newer head = %+v, want a new job in the waiting one's place", newer)
	}
	if n := lab.queued(); n != 2 {
		t.Fatalf("queue holds %d jobs, want one for each story", n)
	}
	lab.drain()
	want := "deploy " + first + " " + klHeadA + ";deploy " + story + " " + klHeadC
	if got := strings.Join(lab.kluster.deploys(), ";"); got != want {
		t.Fatalf("deploys = %s, want %s", got, want)
	}
	if j := lab.asked(story); j.ID != newer.Job || j.Commit != klHeadC {
		t.Fatalf("the story's job state = %+v, want job %s at %s", j, newer.Job, short(klHeadC))
	}
}

// ADR-0292 §4: kluster.agent_namespaces stories in progress hold a namespace
// at most: a third story's request is refused naming the limit and each
// holder, never a request of a story holding one, and stories in review do
// not count.
func TestTheLimitRefusesAThirdStoryAndNeverOneHoldingANamespace(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	one := lab.inProgress("One", klHeadA)
	two := lab.inProgress("Two", klHeadB)
	three := lab.inProgress("Three", klHeadC)
	lab.ask(one, kluster.KindDeploy, "agent-"+one)
	lab.ask(two, kluster.KindDeploy, "agent-"+two)
	lab.drain()

	a := lab.ask(three, kluster.KindDeploy, "agent-"+three)
	if !a.Refused() || a.Job != "" {
		t.Fatalf("a third story's answer = %+v, want a refusal", a)
	}
	for _, w := range []string{"kluster.agent_namespaces allows (2)", one, two, short(klHeadA), short(klHeadB), "ask again later"} {
		if !strings.Contains(a.Reason, w) {
			t.Errorf("the refusal does not name %q:\n%s", w, a.Reason)
		}
	}

	lab.git.set(one, klHeadC)
	if again := lab.ask(one, kluster.KindDeploy, "agent-"+one); !again.Started {
		t.Fatalf("a story holding a namespace was refused for the limit: %+v", again)
	}
	lab.drain()

	lab.move(two, workitem.Review)
	if b := lab.ask(three, kluster.KindDeploy, "agent-"+three); !b.Started {
		t.Fatalf("a story in review counted against the limit: %+v", b)
	}
}

// ADR-0292 §3: a namespace stays while its story is in progress, and is torn
// down on a move to ready, on cancellation, and on a teardown its agent asks
// for.
func TestANamespaceIsKeptWhileInProgressAndTornDownWhenItLeaves(t *testing.T) {
	for _, how := range []string{"ready", "cancelled", "asked"} {
		t.Run(how, func(t *testing.T) {
			lab := newKlusterLab(t, klusterBlock)
			story := lab.inProgress("Kept", klHeadA)
			lab.ask(story, kluster.KindDeploy, klAgent)
			lab.drain()
			n := len(lab.kluster.all())
			lab.round()
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("looks at the story in progress ran %v, want its namespace kept", got)
			}
			switch how {
			case "ready":
				lab.move(story, workitem.Ready)
				lab.round()
			case "cancelled":
				lab.move(story, workitem.Cancelled)
				lab.round()
			case "asked":
				a := lab.ask(story, kluster.KindTeardown, klAgent)
				if !a.Started || a.Kind != kluster.KindTeardown {
					t.Fatalf("teardown answer = %+v, want a job started", a)
				}
				lab.drain()
				if j := lab.asked(story); j.Kind != kluster.KindTeardown || j.Outcome != kluster.OutcomePassed || j.ID != a.Job {
					t.Fatalf("the teardown's state = %+v, want it passed", j)
				}
				if evs := lab.klusterEvents(); len(evs) != 2 || evs[1].Fields[inbox.FieldKind] != kluster.KindTeardown {
					t.Fatalf("kluster events = %+v, want the deploy's and the teardown's", evs)
				}
				if lab.kluster.by[len(lab.kluster.by)-1] != klAgent {
					t.Errorf("the teardown was for %v, want %s", lab.kluster.by, klAgent)
				}
			}
			if got := lab.kluster.since(n); strings.Join(got, ";") != "teardown "+story {
				t.Fatalf("calls once %s = %v, want one teardown", how, got)
			}
			ns, _ := kluster.StoryNamespace(story)
			if _, ok, _ := kluster.ReadRecord(lab.root, ns); ok {
				t.Fatalf("the namespace's record is kept after its teardown")
			}
			n = len(lab.kluster.all())
			lab.round()
			if got := lab.kluster.since(n); len(got) != 0 {
				t.Fatalf("a look after the teardown ran %v", got)
			}
		})
	}
}

// ADR-0292 §1, §3: a story deployed on request that enters review at that
// head is not deployed again, its request answered from the record; at a new
// head it is deployed in place; accepted, it is torn down and main promoted.
func TestAStoryInReviewAtItsDeployedHeadIsNotDeployedAgain(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inProgress("Review", klHeadA)
	lab.ask(story, kluster.KindDeploy, klAgent)
	lab.drain()
	lab.move(story, workitem.Review)
	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a look at the story in review at its deployed head ran %v, want nothing", got)
	}
	a := lab.ask(story, kluster.KindDeploy, klAgent)
	if !a.Recorded || a.Commit != klHeadA || a.Outcome != kluster.OutcomePassed || a.Job != "" || a.Refused() {
		t.Fatalf("answer in review at the deployed head = %+v, want the record's", a)
	}
	if n := lab.queued(); n != 0 {
		t.Fatalf("a request in review queued %d jobs", n)
	}

	lab.git.set(story, klHeadB)
	lab.round()
	if got := strings.Join(lab.kluster.since(n), ";"); got != "absent;deploy "+story+" "+klHeadB {
		t.Fatalf("calls at a new head in review = %s, want one deploy in place", got)
	}
	ns, _ := kluster.StoryNamespace(story)
	if r, _, _ := kluster.ReadRecord(lab.root, ns); r.AskedBy != forgeAuthor {
		t.Errorf("the review's deploy is recorded as asked by %q, want %q", r.AskedBy, forgeAuthor)
	}

	n = len(lab.kluster.all())
	lab.accept(story)
	lab.round()
	if got := strings.Join(lab.kluster.since(n), ";"); got != "teardown "+story+";promote main" {
		t.Fatalf("calls after the acceptance = %s, want its teardown and a promotion", got)
	}
}

// ADR-0292 §1: a request for a story in review is answered with the review's
// job while it waits, and otherwise says the review deploys it.
func TestARequestForAStoryInReviewIsAnsweredWithTheReviewsJob(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inReview("Waiting", klHeadA)
	lab.pass.look(context.Background())
	a := lab.ask(story, kluster.KindDeploy, klAgent)
	if a.Job == "" || a.Started || a.Refused() {
		t.Fatalf("answer while the review's deploy waits = %+v, want its job", a)
	}
	if n := lab.queued(); n != 1 {
		t.Fatalf("the request queued another job: %d waiting", n)
	}

	other := lab.inReview("Other", klHeadB) // not looked at yet
	b := lab.ask(other, kluster.KindDeploy, klAgent)
	if !b.Refused() || !strings.Contains(b.Reason, "deploys it itself") {
		t.Fatalf("answer for a story in review with nothing queued = %+v, want it told the review deploys it", b)
	}
	if td := lab.ask(other, kluster.KindTeardown, klAgent); !td.Refused() || !strings.Contains(td.Reason, "in review") {
		t.Fatalf("a teardown for a story in review = %+v, want a refusal", td)
	}
}

// ADR-0292 §5: an asked-for deploy that fails opens no thread, moves nothing,
// and tells its step, exit code, and log in its state and its kluster event.
func TestAFailedAskedDeployOpensNoThreadAndTellsItsEvent(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	story := lab.inProgress("Broken", klHeadA)
	ns, _ := kluster.StoryNamespace(story)
	log := filepath.Join(kluster.LogDir(lab.root, ns), "deploy.log")
	lab.kluster.fail = func(string, string) error {
		return &kluster.StepError{Namespace: ns, Step: kluster.StepDeploy, ExitCode: 3, Log: log, Err: errors.New("exit status 3")}
	}
	a := lab.ask(story, kluster.KindDeploy, klAgent)
	lab.drain()
	if ths := lab.threadsOn(story); len(ths) != 0 {
		t.Fatalf("a failed asked-for deploy opened %d threads", len(ths))
	}
	if got := lab.status(story); got != workitem.InProgress {
		t.Fatalf("the story is %s after its asked-for deploy failed, want it left in progress", got)
	}
	j := lab.asked(story)
	if j.Outcome != kluster.OutcomeFailed || j.FailedStep != kluster.StepDeploy || j.ExitCode == nil || *j.ExitCode != 3 || j.Log != log || j.Commit != klHeadA {
		t.Fatalf("the job's state = %+v, want it failed at the deploy step with exit code 3 and its log", j)
	}
	evs := lab.klusterEvents()
	if len(evs) != 1 {
		t.Fatalf("kluster events = %+v, want one", evs)
	}
	f := evs[0].Fields
	if f[inbox.FieldJob] != a.Job || f[inbox.FieldOutcome] != kluster.OutcomeFailed || f[inbox.FieldFailedStep] != kluster.StepDeploy || f[inbox.FieldExitCode] != "3" || f[inbox.FieldLog] != log {
		t.Fatalf("the event's fields = %v, want the failure's step, exit code, and log", f)
	}
	n := len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a look after the failed asked-for deploy ran %v, want nothing: the agent asks again", got)
	}
}

// ADR-0292 §5, §3: with checks on, an asked-for deploy runs its checks; a
// run that fails ends the job failed at checks with no thread, and the story
// entering review at that head is deployed and checked again.
func TestAnAskedDeploysFailingChecksOpenNoThreadAndReviewDeploysAgain(t *testing.T) {
	lab := newKlusterLab(t, klusterBlock)
	failing, _ := nsCheck(t, 2)
	lab.withChecks(failing)
	story := lab.inProgress("Checked", klHeadA)
	lab.ask(story, kluster.KindDeploy, klAgent)
	lab.drain()
	if ths := lab.titled(story, checksFailureTitle(story)); len(ths) != 0 {
		t.Fatalf("failing checks of an asked-for deploy opened %d threads", len(ths))
	}
	if j := lab.asked(story); j.Outcome != kluster.OutcomeFailed || j.FailedStep != "checks" || j.ExitCode == nil || *j.ExitCode != 2 {
		t.Fatalf("the job's state = %+v, want it failed at checks with exit code 2", j)
	}
	if got := lab.status(story); got != workitem.InProgress {
		t.Fatalf("the story is %s, want it left in progress", got)
	}
	passing, _ := nsCheck(t, 0)
	lab.withChecks(passing)
	lab.move(story, workitem.Review)
	n := len(lab.kluster.all())
	lab.round()
	if got := strings.Join(lab.kluster.since(n), ";"); got != "absent;deploy "+story+" "+klHeadA {
		t.Fatalf("calls once in review at a head whose checks failed = %s, want one deploy", got)
	}
	n = len(lab.kluster.all())
	lab.round()
	if got := lab.kluster.since(n); len(got) != 0 {
		t.Fatalf("a look once deploy and checks passed at the head ran %v", got)
	}
}

// ADR-0292 §5: a job an earlier flai serve left waiting or running ends
// interrupted when the pass starts again, with its event.
func TestAnAskedJobLeftByAnEarlierServeEndsInterrupted(t *testing.T) {
	for _, started := range []bool{false, true} {
		name := "waiting"
		if started {
			name = "running"
		}
		t.Run(name, func(t *testing.T) {
			lab := newKlusterLab(t, klusterBlock)
			story := lab.inProgress("Interrupted", klHeadA)
			a := lab.ask(story, kluster.KindDeploy, klAgent)
			if started {
				j := lab.asked(story)
				j.StartedAt = lab.opts.Now().UTC()
				if err := kluster.WriteAskedJob(lab.root, j); err != nil {
					t.Fatal(err)
				}
			}
			lab.restart()
			lab.pass.takeRequests(context.Background())
			j := lab.asked(story)
			if !j.Ended() || j.Outcome != kluster.OutcomeInterrupted || j.ID != a.Job {
				t.Fatalf("the job's state after a restart = %+v, want it ended interrupted", j)
			}
			evs := lab.klusterEvents()
			if len(evs) != 1 || evs[0].Fields[inbox.FieldOutcome] != kluster.OutcomeInterrupted || evs[0].Fields[inbox.FieldJob] != a.Job {
				t.Fatalf("kluster events = %+v, want one telling the job was interrupted", evs)
			}
			lab.pass.takeRequests(context.Background())
			if evs := lab.klusterEvents(); len(evs) != 1 {
				t.Fatalf("a second look told the interruption again: %+v", evs)
			}
		})
	}
}

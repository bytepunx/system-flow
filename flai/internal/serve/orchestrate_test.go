//go:build !windows

package serve

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// orchestrateLab is an agentLab whose command is an orchestrator stub, with
// the project's orchestrator on the lab's launcher and a clock of its own,
// later than the launcher's by shift, or pinned to at when at is set. The stub prints stream-orchestrator as
// its output when there is one, then writes what it was given to
// orchestrator-<pid>.txt, so that once given has read it, the run's usage is
// in its log and a stop cannot come before it (I-0090). It then waits while
// held until released, marks its exit, as it would have ended, when asked to
// stop, and exits with the code in exit-orchestrator when there is one. The
// orchestrate action is off.
type orchestrateLab struct {
	*agentLab
	orch  *orchestrator
	shift time.Duration
	at    time.Time
}

func newOrchestrateLab(t *testing.T) *orchestrateLab {
	t.Helper()
	lab := &orchestrateLab{agentLab: newAgentLab(t)}
	d := lab.outDir
	stub := filepath.Join(d, "orch-agent")
	script := "#!/bin/sh\n" +
		"trap 'touch \"" + d + "/ended-$$\"; exit 143' TERM\n" +
		"[ -f \"" + d + "/stream-orchestrator\" ] && cat \"" + d + "/stream-orchestrator\"\n" +
		"{ echo \"args: $*\"; echo \"dir: $(pwd)\"; echo \"agent: $FLAI_AGENT\"; echo \"story: ${FLAI_STORY-unset}\"; echo \"role: $FLAI_ROLE\"; echo \"item: ${FLAI_ITEM-unset}\"; echo \"session: $FLAI_SESSION\"; echo \"by: $FLAI_STARTED_BY\"; } > \"" + d + "/orchestrator-$$.txt\"\n" +
		"while [ -f \"" + d + "/hold\" ] && [ ! -f \"" + d + "/release-orchestrator\" ]; do sleep 0.05; done\n" +
		"touch \"" + d + "/ended-$$\"\n" +
		"[ -f \"" + d + "/exit-orchestrator\" ] && exit \"$(cat \"" + d + "/exit-orchestrator\")\"\n" +
		"exit 0\n"
	if err := os.WriteFile(stub, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	lab.cfg.Command = []string{stub, "story={story}"}
	o := lab.o
	o.Now = func() time.Time {
		if !lab.at.IsZero() {
			return lab.at
		}
		return time.Now().Add(lab.shift)
	}
	lab.orch = newOrchestrator(o, Entry{Key: "t", Name: "t", Root: lab.root}, lab.l)
	return lab
}

func (lab *orchestrateLab) look() { lab.orch.look(context.Background()) }

func (lab *orchestrateLab) orchestrator() *AgentRun { return lab.state().Orchestrator }

// given is what the stub was given in the run with pid, once it has written
// it.
func (lab *orchestrateLab) given(pid int) string {
	lab.t.Helper()
	out := filepath.Join(lab.outDir, "orchestrator-"+strconv.Itoa(pid)+".txt")
	waitFor(lab.t, "the orchestrator stub ran", func() bool {
		data, _ := os.ReadFile(out)
		return strings.Contains(string(data), "by: ")
	})
	data, _ := os.ReadFile(out)
	return string(data)
}

// ended waits until the orchestrator run with session has ended, and
// returns it.
func (lab *orchestrateLab) ended(session string) *AgentRun {
	lab.t.Helper()
	waitFor(lab.t, "the orchestrator run ended", func() bool {
		r := lab.orchestrator()
		return r != nil && r.Session == session && r.Ended != ""
	})
	return lab.orchestrator()
}

// orchestrateEntries are the journal's entries for the orchestrate action.
func (lab *orchestrateLab) orchestrateEntries() []hostapi.Entry {
	var out []hostapi.Entry
	for _, e := range lab.entries() {
		if e.Action == hostapi.ActionOrchestrate {
			out = append(out, e)
		}
	}
	return out
}

// S-0218: with orchestrate off nothing starts; turned on, one orchestrator
// run starts for the project, in its main checkout, as the orchestrator and
// for no story or item, logged as the orchestrator's and recorded under
// orchestrator, apart from the stories' runs; while it runs, no other is
// started, and it holds no story's agent back.
func TestTheOrchestratorRunsWhileTheActionIsOn(t *testing.T) {
	t.Setenv("FLAI_STORY", "S-0999") // the session flai serve was started from
	lab := newOrchestrateLab(t)
	lab.hold()
	lab.look()
	if st := lab.state(); st.Orchestrator != nil || len(lab.orchestrateEntries()) != 0 {
		t.Fatalf("started while orchestrate is off: %+v", st.Orchestrator)
	}
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	if run == nil || !run.live() || run.Story != "" || run.Item != "" || run.Agent != "orchestrator" || run.Harness != "command" || run.Command != "orch-agent" || run.Session == "" {
		t.Fatalf("run = %+v", run)
	}
	logs, err := lab.o.Dir.ActivityLogs("t", workitem.ActivityOrchestrator)
	if err != nil || len(logs) != 1 || logs[0] != run.Log || !strings.HasPrefix(filepath.Base(run.Log), "t-orchestrator-") {
		t.Errorf("orchestrator logs = %v (%v), want the run's %s", logs, err, run.Log)
	}
	got := lab.given(run.PID)
	real, _ := filepath.EvalSymlinks(lab.root)
	for _, want := range []string{"args: story=\n", "dir: " + real + "\n", "agent: orchestrator\n", "story: unset\n", "role: orchestrate\n", "item: unset\n", "session: 2", "by: flai-serve\n"} {
		if !strings.Contains(got, want) {
			t.Errorf("the orchestrator was given\n%s\nwithout %q", got, want)
		}
	}
	if data, err := os.ReadFile(lab.o.Dir.agents()); err != nil || !strings.Contains(string(data), `"orchestrator": {`) {
		t.Errorf("serve/agents.json does not record the run under orchestrator (%v):\n%s", err, data)
	}
	if st := lab.state(); st.Running != nil || st.Last != nil || len(st.Stories) != 0 || len(st.Plans) != 0 {
		t.Errorf("state = %+v, want the run under orchestrator alone", st)
	}
	j := lab.orchestrateEntries()
	if len(j) != 1 || j[0].Method != "serve.orchestrate" || j[0].Outcome != "done" || !strings.Contains(j[0].Detail, "started orch-agent (command) as orchestrator (pid ") {
		t.Errorf("journal: %+v", j)
	}
	lab.look()
	if again := lab.orchestrator(); !again.same(run) || len(lab.orchestrateEntries()) != 1 {
		t.Errorf("a second run was started beside the first: %+v", again)
	}
	// the in-progress limit does not count it
	lab.limit(1)
	story := lab.ready("Beside")
	lab.l.look(context.Background(), false)
	waitFor(t, "the story's agent started beside the orchestrator", func() bool { return lab.run(story).live() })
	lab.release(story)
	lab.release("orchestrator")
}

// S-0218: a run that ends while the action is on is started again at the
// next look, and, when it failed, no sooner than a minute after it ended;
// each end is journalled.
func TestTheOrchestratorIsStartedAgainWhenItEnds(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	first := lab.orchestrator()
	lab.release("orchestrator")
	if r := lab.ended(first.Session); r.Outcome != OutcomeWorked || r.Exit == nil || *r.Exit != 0 {
		t.Fatalf("first = %+v, want worked with exit 0", r)
	}
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-orchestrator"), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.look()
	second := lab.orchestrator()
	if second.Session == first.Session || second.Started == "" {
		t.Fatalf("not started again after it worked: %+v", second)
	}
	failed := lab.ended(second.Session)
	if failed.Outcome != OutcomeFailed || failed.Exit == nil || *failed.Exit != 3 || failed.Why != "ended (exit 3)" {
		t.Fatalf("second = %+v, want failed with exit 3", failed)
	}
	// the looks around the retry take their time from the recorded end, which
	// is to the second, not from the host's clock, so neither depends on how
	// long the host took to get there (I-0106)
	ended, err := time.Parse(time.RFC3339, failed.Ended)
	if err != nil {
		t.Fatalf("the failed run's end %q: %v", failed.Ended, err)
	}
	lab.at = ended
	lab.look()
	lab.at = ended.Add(orchestrateRetry - time.Second)
	lab.look()
	if r := lab.orchestrator(); r.Session != second.Session {
		t.Fatalf("started again within a minute of a failure: %+v", r)
	}
	lab.at = ended.Add(orchestrateRetry)
	lab.look()
	third := lab.orchestrator()
	if third.Session == second.Session {
		t.Fatalf("not started again a minute after a failure: %+v", third)
	}
	lab.ended(third.Session)
	var ends []hostapi.Entry
	for _, e := range lab.orchestrateEntries() {
		if strings.Contains(e.Detail, ") ended, ") {
			ends = append(ends, e)
		}
	}
	if len(ends) != 3 || ends[0].Outcome != "done" || !strings.HasSuffix(ends[0].Detail, "ended, worked") ||
		ends[1].Outcome != "failed" || !strings.HasSuffix(ends[1].Detail, "ended, failed: ended (exit 3)") {
		t.Errorf("ends journalled: %+v", ends)
	}
}

// S-0218: turning the action off stops the run going, as the operator's stop
// of a story's agent does, whoever waits for it, journals the stop, and logs
// the run's activity in orchestrator.md as stopped for the action.
func TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff(t *testing.T) {
	for name, handOver := range map[string]bool{"a run this flai serve waits for": false, "a run an earlier flai serve started": true} {
		t.Run(name, func(t *testing.T) {
			lab := newOrchestrateLab(t)
			lab.hold()
			stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999)}, "\n") + "\n"
			if err := os.WriteFile(filepath.Join(lab.outDir, "stream-orchestrator"), []byte(stream), 0o644); err != nil {
				t.Fatal(err)
			}
			lab.l.handOver = handOver
			lab.cfg.Orchestrate = true
			lab.look()
			run := lab.orchestrator()
			if handOver {
				// the process an earlier flai serve left is reaped by the
				// system once it ends, not by this one
				go func() {
					if p, err := os.FindProcess(run.PID); err == nil {
						_, _ = p.Wait()
					}
				}()
			}
			lab.given(run.PID) // it runs, has set its trap, and has printed its usage
			lab.cfg.Orchestrate = false
			lab.look()
			if Alive(run.PID) {
				t.Errorf("pid %d still runs", run.PID)
			}
			stopped := lab.ended(run.Session)
			if stopped.Outcome != OutcomeStopped || stopped.Stopped == "" || stopped.live() {
				t.Errorf("not recorded as stopped: %+v", stopped)
			}
			// S-0532: and its stop is one stopped agent event
			lab.orchestratorEventsAre([2]string{inbox.AgentStarted, "started while orchestrate is on"}, [2]string{inbox.AgentStopped, stopped.Why})
			waitFor(t, "its activity logged", func() bool {
				doc, err := lab.repo.Activity(workitem.ActivityOrchestrator)
				return err == nil && len(doc.Entries) == 1
			})
			doc, _ := lab.repo.Activity(workitem.ActivityOrchestrator)
			if e := doc.Entries[0]; e.Summary != orchestrateOff || len(e.Items) != 0 || e.Seconds != 61 {
				t.Errorf("activity = %+v, want %q over the run's 61 s and no items", e, orchestrateOff)
			}
			j := lab.orchestrateEntries()
			if last := j[len(j)-1]; last.Outcome != "done" || !strings.Contains(last.Detail, "stopped the orchestrator orchestrator, since orchestrate was turned off: ended its process group (pid "+strconv.Itoa(run.PID)+")") {
				t.Errorf("journal: %+v", j)
			}
			// off, it is not started again
			lab.look()
			if r := lab.orchestrator(); r.Session != run.Session {
				t.Errorf("started while off: %+v", r)
			}
		})
	}
}

// S-0218: the orchestrator is refused, saying why, when its agent names no
// harness and no command is set, and when the harness refuses it, as
// claude-code does without orchestrator.md; each refusal is recorded and
// journalled once, not at every look.
func TestTheOrchestratorIsRefusedOnce(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.cfg.Command, lab.cfg.Orchestrate = nil, true
	lab.look()
	run := lab.orchestrator()
	if run == nil || run.Outcome != OutcomeFailed || !strings.Contains(run.Error, "the orchestrator's agent names no harness (orchestration.agent or agent in system-flow.yaml), and no command is set on the host") {
		t.Fatalf("no harness: %+v", run)
	}
	for _, shift := range []time.Duration{0, 2 * orchestrateRetry} {
		lab.shift = shift
		lab.look()
	}
	if j := lab.orchestrateEntries(); len(j) != 1 || j[0].Outcome != "failed" || !strings.Contains(j[0].Detail, "the orchestrator could not be started: the orchestrator's agent names no harness") {
		t.Errorf("journal: %+v", j)
	}
	if r := lab.orchestrator(); r.Started != run.Started {
		t.Errorf("the refusal was recorded again: %+v", r)
	}
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  agent:\n    harness: claude-code\n"
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.shift = 4 * orchestrateRetry
	lab.look()
	if r := lab.orchestrator(); r.Harness != "claude-code" || !strings.Contains(r.Error, "the orchestrator's session runs as the agent orchestrator, and its definition could not be read") {
		t.Errorf("no orchestrator.md: %+v", r)
	}
	lab.shift = 6 * orchestrateRetry
	lab.look()
	if j := lab.orchestrateEntries(); len(j) != 2 {
		t.Errorf("journal: %+v", j)
	}
	lab.logs.mu.Lock()
	logged := lab.logs.buf.String()
	lab.logs.mu.Unlock()
	if n := strings.Count(logged, `msg="orchestrator could not be started"`); n != 2 {
		t.Errorf("the refusals were logged %d times, want once each:\n%s", n, logged)
	}
}

// The orchestrator is told the merge queue runs only while the merge-queue
// host action is on for its project, accept_reviews is on, and acceptance
// is local, so that its prompt says to enqueue then and to accept otherwise
// (S-0403).
func TestTheOrchestratorIsToldTheQueueRunsOnlyWhileItBatches(t *testing.T) {
	root := t.TempDir()
	e := Entry{Root: root}
	for _, tc := range []struct {
		name       string
		on, accept bool
		mode       string
		other      bool
		want       bool
	}{
		{name: "action, accept_reviews, local", on: true, accept: true, want: true},
		{name: "action off", accept: true},
		{name: "accept_reviews off", on: true},
		{name: "pull-request acceptance", on: true, accept: true, mode: manifest.AcceptPullRequest},
		{name: "action on for another project", on: true, accept: true, other: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			at := root
			if tc.other {
				at = filepath.Join(root, "other")
			}
			o := Options{Host: hostapi.Host{Enabled: func(action, r string) bool {
				return tc.on && action == hostapi.ActionMergeQueue && r == at
			}}}
			var m manifest.Manifest
			m.Orchestration.Permissions.AcceptReviews = tc.accept
			m.Acceptance.Mode = tc.mode
			if got := queueBatches(o, e, m); got != tc.want {
				t.Errorf("queueBatches = %v, want %v", got, tc.want)
			}
		})
	}
	if queueBatches(Options{}, e, manifest.Manifest{}) {
		t.Error("queueBatches with no host = true, want false")
	}
}

// orchestratorEventsAre says that the agent events on the orchestrator's
// role are one per change, to and why as want, in turn (S-0532).
func (lab *orchestrateLab) orchestratorEventsAre(want ...[2]string) {
	lab.t.Helper()
	lab.agentEventsAre(workitem.ActivityOrchestrator, inbox.RoleType, workitem.ActivityOrchestrator, workitem.ActivityOrchestrator, conventions.RoleOrchestrate, want...)
}

// S-0532: each orchestrator run appends one agent event on the orchestrator's
// role at its start, saying why it was started, and one at its end, ended for
// a run that worked and failed for one that did not, so that a run that fails
// and is retried every minute writes one failed and one started per attempt.
func TestAnOrchestratorRunAppendsAnAgentEventAtEachStartAndEnd(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	first := lab.orchestrator()
	lab.release("orchestrator")
	lab.ended(first.Session)
	started := [2]string{inbox.AgentStarted, "started while orchestrate is on"}
	worked := [2]string{inbox.AgentEnded, "ended with its work done"}
	lab.orchestratorEventsAre(started, worked)
	if err := os.WriteFile(filepath.Join(lab.outDir, "exit-orchestrator"), []byte("3"), 0o644); err != nil {
		t.Fatal(err)
	}
	want := [][2]string{started, worked, {inbox.AgentStarted, "started again while orchestrate is on: its last run ended, worked"}}
	failed := [2]string{inbox.AgentFailed, "ended (exit 3)"}
	retried := [2]string{inbox.AgentStarted, "started again after its last run failed: ended (exit 3)"}
	lab.look()
	for attempt := range 3 {
		run := lab.ended(lab.orchestrator().Session)
		want = append(want, failed)
		lab.orchestratorEventsAre(want...)
		if attempt == 2 {
			break
		}
		ended, err := time.Parse(time.RFC3339, run.Ended)
		if err != nil {
			t.Fatalf("the failed run's end %q: %v", run.Ended, err)
		}
		lab.at = ended.Add(orchestrateRetry - time.Second)
		lab.look()
		lab.orchestratorEventsAre(want...)
		lab.at = ended.Add(orchestrateRetry)
		lab.look()
		if r := lab.orchestrator(); r.Session == run.Session {
			t.Fatalf("attempt %d: not started again a minute after a failure: %+v", attempt+1, r)
		}
		want = append(want, retried)
	}
}

// S-0532: an orchestrator start that could not be made is one failed agent
// event, for its error, as it is journalled once and not at every look.
func TestAnOrchestratorStartThatCouldNotBeMadeIsOneFailedAgentEvent(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.cfg.Command, lab.cfg.Orchestrate = nil, true
	lab.look()
	run := lab.orchestrator()
	if run == nil || !strings.HasPrefix(run.Why, "could not be started: ") {
		t.Fatalf("not recorded as failed to start: %+v", run)
	}
	for _, shift := range []time.Duration{0, 2 * orchestrateRetry, 4 * orchestrateRetry} {
		lab.shift = shift
		lab.look()
	}
	lab.orchestratorEventsAre([2]string{inbox.AgentFailed, run.Why})
}

// S-0534: an orchestrator run whose log ends on the recorded usage-limit
// refusal is failed on the limit and sets the budget state until the reset
// the API reported; orchestrateRetry passing starts no other run while the
// state holds, and the first look once its hold has passed does.
func TestAnOrchestratorRunOnTheUsageLimitWaitsForTheBudgetStateNotOrchestrateRetry(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.at = limitReset.Add(-time.Hour)
	lab.budgetClock(&lab.at)
	lab.cfg.Orchestrate = true
	one := 1
	run := &AgentRun{Agent: "orchestrator", Harness: harness.ClaudeCode, Session: "limited", Started: lab.at.Format(time.RFC3339), Log: limitResult(t, true)}
	lab.l.orchestrateEnded(run, &one)
	lab.heldUntil(limitReset, true, "5m0s")
	if r := lab.orchestrator(); r.Outcome != OutcomeFailed || !r.UsageLimited || r.ResetsAt != limitResetsAt || r.Why != limitWhyText {
		t.Errorf("orchestrator run = %+v; want failed on the limit", r)
	}
	ended := false
	for _, e := range lab.entries() {
		ended = ended || (e.Action == hostapi.ActionOrchestrate && e.Outcome == "failed" && strings.HasSuffix(e.Detail, "ended, failed: "+limitWhyText))
	}
	if !ended {
		t.Errorf("journal: %+v; want the end failed on the limit", lab.entries())
	}

	lab.at = lab.at.Add(2 * orchestrateRetry)
	lab.look()
	lab.look()
	if r := lab.orchestrator(); r.Session != "limited" || r.live() || lab.startedCount() != 0 {
		t.Errorf("started while the budget state holds: %+v", r)
	}

	lab.at = limitReset
	lab.look()
	waitFor(t, "the orchestrator is started once the hold has passed", func() bool { return lab.orchestrator().Session != "limited" })
	lab.settle()
}

// S-0534: while the budget state of the orchestrator's provider holds, as
// another role's run on the limit set it, no look starts the orchestrator,
// nothing is recorded or journalled, and the hold is logged once; a run
// that failed otherwise before it is not started again as orchestrateRetry
// passes, and the first look once the state clears starts it.
func TestTheOrchestratorWaitsWhileTheBudgetStateHoldsItsProvider(t *testing.T) {
	lab := newOrchestrateLab(t)
	lab.cfg.Orchestrate = true
	lab.budgetHeldOn(harness.Command, anHour())
	lab.look()
	lab.look()
	if r := lab.orchestrator(); r != nil || len(lab.orchestrateEntries()) != 0 {
		t.Fatalf("started while held: %+v, journalled %+v", r, lab.orchestrateEntries())
	}
	lab.o.Dir.updateAgent(lab.root, func(s *AgentState) {
		s.put(&AgentRun{Agent: workitem.ActivityOrchestrator, Session: "failed", Started: "2026-10-28T08:00:00Z", Ended: "2026-10-28T08:00:05Z", Outcome: OutcomeFailed, Why: "ended (exit 3)"})
	})
	lab.at = time.Date(2026, 10, 28, 8, 0, 5, 0, time.UTC).Add(2 * orchestrateRetry)
	lab.look()
	lab.at = lab.at.Add(orchestrateRetry)
	lab.look()
	if r := lab.orchestrator(); r.Session != "failed" || len(lab.orchestrateEntries()) != 0 {
		t.Fatalf("started again after orchestrateRetry while held: %+v", r)
	}
	if n := strings.Count(lab.logText(), `msg="orchestrator start held"`); n != 1 {
		t.Errorf("the hold was logged %d times, want once:\n%s", n, lab.logText())
	}
	lab.budgetCleared()
	lab.look()
	if r := lab.orchestrator(); r.Session == "failed" || r.Started == "" {
		t.Fatalf("not started at the first look once the state cleared: %+v", r)
	}
	lab.settle()
}

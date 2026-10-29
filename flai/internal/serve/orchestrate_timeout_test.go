//go:build !windows

package serve

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// timeoutLab is an orchestrateLab whose project sets orchestration.run_timeout
// to 2s and coalesce to 1s, with its orchestrator turned on and its first run
// started, held, and running, its usage printed so that its activity is
// logged when it ends; it returns the run and when it started.
func timeoutLab(t *testing.T) (*orchestrateLab, *AgentRun, time.Time) {
	t.Helper()
	lab := newOrchestrateLab(t)
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\norchestration:\n  run_timeout: 2s\n  coalesce: 1s\n"
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	stream := strings.Join([]string{streamCall("s", "m1", runStart, 999), streamCall("s", "m2", runStart.Add(time.Minute), 999)}, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(lab.outDir, "stream-orchestrator"), []byte(stream), 0o644); err != nil {
		t.Fatal(err)
	}
	lab.hold()
	lab.cfg.Orchestrate = true
	lab.look()
	run := lab.orchestrator()
	if run == nil || !run.live() {
		t.Fatalf("the orchestrator did not start: %+v", run)
	}
	lab.given(run.PID) // it runs and has set its trap
	started, err := time.Parse(time.RFC3339, run.Started)
	if err != nil {
		t.Fatalf("the run's start %q: %v", run.Started, err)
	}
	return lab, run, started
}

// stopsJournalled are the journal's entries for the orchestrator's stops.
func (lab *orchestrateLab) stopsJournalled() []string {
	var out []string
	for _, e := range lab.orchestrateEntries() {
		if strings.HasPrefix(e.Detail, "stopped the orchestrator ") {
			out = append(out, e.Detail)
		}
	}
	return out
}

// S-0536 criterion 3, ADR-0209 §4: a run past orchestration.run_timeout is
// stopped at the next look, recorded stopped and not failed, for a why that
// names the timeout, which its stopped agent event and its activity carry;
// the stop is journalled and logged with the pid and the timeout; and a
// reason queued while it ran starts the next run at the first look once its
// window has passed, with no orchestrateRetry to wait.
func TestAnOrchestratorRunPastRunTimeoutIsStopped(t *testing.T) {
	lab, run, started := timeoutLab(t)
	lab.queueReason(reviewReason, started.Add(time.Second))
	lab.at = started.Add(3 * time.Second)
	lab.look()
	if Alive(run.PID) {
		t.Errorf("pid %d still runs past run_timeout", run.PID)
	}
	why := "stopped: ran past orchestration.run_timeout (2s)"
	stopped := lab.ended(run.Session)
	if stopped.Outcome != OutcomeStopped || stopped.Why != why {
		t.Errorf("recorded %+v, want stopped for %q", stopped, why)
	}
	lab.orchestratorEventsAre([2]string{inbox.AgentStarted, turnedOn}, [2]string{inbox.AgentStopped, why})
	waitFor(t, "its activity logged", func() bool {
		doc, err := lab.repo.Activity(workitem.ActivityOrchestrator)
		return err == nil && len(doc.Entries) == 1
	})
	if doc, _ := lab.repo.Activity(workitem.ActivityOrchestrator); doc.Entries[0].Summary != why {
		t.Errorf("activity = %+v, want %q", doc.Entries[0], why)
	}
	wantStop := "stopped the orchestrator orchestrator, since it ran past orchestration.run_timeout (2s): ended its process group (pid " + strconv.Itoa(run.PID) + ")"
	if stops := lab.stopsJournalled(); len(stops) != 1 || stops[0] != wantStop {
		t.Errorf("stops journalled: %q, want %q", stops, wantStop)
	}
	for _, e := range lab.orchestrateEntries() {
		if e.Outcome == "failed" {
			t.Errorf("journalled as a failure: %+v", e)
		}
	}
	logged := lab.logText()
	if n := strings.Count(logged, `msg="orchestrator stopped past its run timeout"`); n != 1 ||
		!strings.Contains(logged, "pid="+strconv.Itoa(run.PID)+" run_timeout=2s") {
		t.Errorf("the stop was logged %d times, want once with its pid and run_timeout:\n%s", n, logged)
	}

	lab.at = started.Add(4 * time.Second)
	lab.look()
	next := lab.orchestrator()
	if next.Session == run.Session || !next.live() {
		t.Fatalf("the reason queued while it ran did not start the next run: %+v", next)
	}
	if ended, err := time.Parse(time.RFC3339, stopped.Ended); err != nil || !lab.at.Before(ended.Add(orchestrateRetry)) {
		t.Fatalf("the next run started %s, not within orchestrateRetry of the stop at %q (%v)", lab.at, stopped.Ended, err)
	}
	if got := lab.reasonsGiven(next.PID); len(got) != 1 || got[0] != "moved S-0001 to review: Reviewed" {
		t.Errorf("FLAI_REASONS = %q, want the move queued while the stopped run went", got)
	}
	lab.orchestratorEventsAre([2]string{inbox.AgentStarted, turnedOn}, [2]string{inbox.AgentStopped, why},
		[2]string{inbox.AgentStarted, "started for 1 reason: moved S-0001 to review: Reviewed"})
	lab.release("orchestrator")
	lab.ended(next.Session)
}

// ADR-0209 §4: a run within orchestration.run_timeout, up to the moment it
// is reached, is left alone.
func TestAnOrchestratorRunWithinRunTimeoutIsLeftAlone(t *testing.T) {
	lab, run, started := timeoutLab(t)
	for _, at := range []time.Time{started.Add(time.Second), started.Add(2 * time.Second)} {
		lab.at = at
		lab.look()
		if r := lab.orchestrator(); !r.same(run) || !r.live() || !Alive(run.PID) {
			t.Fatalf("at %s, within run_timeout, the run was not left going: %+v", at, r)
		}
	}
	if stops := lab.stopsJournalled(); len(stops) != 0 {
		t.Errorf("stops journalled: %q, want none", stops)
	}
	lab.release("orchestrator")
	if r := lab.ended(run.Session); r.Outcome != OutcomeWorked {
		t.Errorf("ended %+v, want worked", r)
	}
}

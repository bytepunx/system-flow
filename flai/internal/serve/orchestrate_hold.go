package serve

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator stopped and held stopped, and started again, on the
// operator's word (S-0228): flai serve orchestrate stop and start, which the
// dashboard's orchestrate.stop and orchestrate.start run. flai serve starts
// the orchestrator again whenever its run ends while the orchestrate action
// is on (S-0218), so a stop alone would be undone at the next look: the stop
// holds it stopped, on its run in serve/agents.json, and the start lifts the
// hold. Turning the action off lifts it too (orchestrator.look).

// orchestrateSettle is how long OrchestratorStop waits, once the run's
// process is gone, for the serving flai that waits for it to record its end,
// before it records the end itself.
var orchestrateSettle = 2 * time.Second

// OrchestratorStop stops the orchestrator of project e and holds it stopped,
// so that flai serve does not start it again while the orchestrate action
// stays on. A run going is ended as Stop ends a story's agent: marked
// stopped, its process group asked to end while it is still the one flai
// started, and made to once o.StopGrace has passed. The serving flai that
// waits for it records its end, and logs its activity in orchestrator.md as
// orchestrateHeld; when none has within orchestrateSettle, it is recorded
// here. A run that is not going, one that failed and waits to be started
// again, is held as it is. The stop is journalled, and the run returned as
// recorded.
//
// It is refused, as a *Refused, while the orchestrate action is off for the
// project, when flai serve has started no orchestrator for it, and when the
// orchestrator is already held stopped.
func OrchestratorStop(o Options, e Entry) (*AgentRun, error) {
	if err := orchestrateOn(o, e, "stop"); err != nil {
		return nil, err
	}
	now := o.Now
	if now == nil {
		now = time.Now
	}
	at := now().UTC().Format(time.RFC3339)
	// Held, and marked stopped, before the signal, so that whoever records
	// the end records the operator's stop, and no look starts another.
	var run *AgentRun
	var no error
	o.Dir.updateAgent(e.Root, func(s *AgentState) {
		cur := s.Orchestrator
		switch {
		case cur == nil:
			no = refused("flai serve has started no orchestrator for this project, so there is none to stop")
			return
		case cur.Held:
			no = refused("the orchestrator is already held stopped: flai serve orchestrate start starts it again")
			return
		}
		marked := *cur
		marked.Held = true
		if marked.live() {
			marked.Stopped = at
		}
		s.put(&marked)
		run = &marked
	})
	if no != nil {
		return nil, no
	}
	l := newLauncher(o, e)
	l.handOver = true
	signalled, killed := false, false
	if run.running() {
		signalled = true
		_ = TerminateGroup(run.PID)
		if !gone(run.PID, o.stopGrace()) {
			killed = true
			_ = KillGroup(run.PID)
			gone(run.PID, 2*time.Second)
		}
	}
	if run.live() && !settled(o, e, run) {
		if cur := o.Dir.AgentStates()[e.Root].Orchestrator; cur.same(run) && cur.live() {
			l.orchestrateEnded(cur, nil)
		}
	}
	what := fmt.Sprintf("it was not running: it ended %s, %s", run.Ended, outcomeOf(run))
	switch {
	case killed:
		what = fmt.Sprintf("killed its process group (pid %d), which had not ended %s after it was asked to", run.PID, o.stopGrace())
	case signalled:
		what = fmt.Sprintf("ended its process group (pid %d)", run.PID)
	case run.live():
		what = fmt.Sprintf("recorded it as stopped: pid %d is no longer the orchestrator flai started", run.PID)
	}
	if o.Host.Record != nil {
		o.Host.Record(hostapi.Entry{At: at, Action: hostapi.ActionOrchestrate, Method: "serve.orchestrate", Project: e.Key, Root: e.Root, By: "flai serve",
			Outcome: "done", Detail: fmt.Sprintf("stopped the orchestrator %s on the operator's word, and held it stopped while orchestrate stays on: %s", run.Agent, what)})
	}
	l.log("orchestrator stopped and held", "pid", run.PID, "signalled", signalled, "killed", killed)
	out := o.Dir.AgentStates()[e.Root].Orchestrator
	if out == nil {
		return nil, fmt.Errorf("the orchestrator's run is no longer recorded in %s; it was held stopped, but flai serve may start it again", o.Dir.agents())
	}
	return out, nil
}

// settled waits up to orchestrateSettle for run, the orchestrator's, to be
// recorded as ended by whoever waits for it, and says whether it was.
func settled(o Options, e Entry, run *AgentRun) bool {
	deadline := time.Now().Add(orchestrateSettle)
	for {
		if cur := o.Dir.AgentStates()[e.Root].Orchestrator; !cur.same(run) || !cur.live() {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// OrchestratorStart lifts the operator's hold on the orchestrator of project
// e and starts it, as flai serve does when the orchestrate action is turned
// on: in the project's main checkout, with its orchestration agent, the run
// handed over to the serving flai, which settles it once its process is gone
// and starts it again whenever it ends. Recording the new run lifts the hold;
// so does a start that fails, which flai serve tries again a minute later, as
// it does after any failure. It returns the run as recorded, and an error
// saying why when the orchestrator could not be started.
//
// It is refused, as a *Refused, while the orchestrate action is off for the
// project, when the orchestrator is not held stopped, which flai serve runs
// by itself, and while the run it was held on is still being stopped.
func OrchestratorStart(ctx context.Context, o Options, e Entry) (*AgentRun, error) {
	if err := orchestrateOn(o, e, "start"); err != nil {
		return nil, err
	}
	run := o.Dir.AgentStates()[e.Root].Orchestrator
	switch {
	case run == nil || !run.Held:
		return nil, refused("the orchestrator is not held stopped: flai serve runs it for as long as orchestrate is on, and starts it again when it ends")
	case run.running():
		return nil, refused("the orchestrator is still being stopped (pid %d); start it once it has ended", run.PID)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return nil, err
	}
	l := newLauncher(o, e)
	l.handOver = true
	l.orchestrate(ctx, o.Agent(e.Root), repo.Manifest.OrchestrationAgent(), "")
	out := o.Dir.AgentStates()[e.Root].Orchestrator
	switch {
	case out == nil || out.Session == run.Session:
		return nil, fmt.Errorf("the orchestrator's new run was not recorded in %s; it is still held stopped", o.Dir.agents())
	case out.Error != "":
		return out, errors.New(out.Why)
	}
	return out, nil
}

// orchestrateOn refuses verb, stop or start, while the orchestrate host
// action is off for project e: flai serve then runs no orchestrator to hold.
func orchestrateOn(o Options, e Entry, verb string) error {
	if o.Agent != nil && o.Agent(e.Root).Orchestrate {
		return nil
	}
	return refused("the orchestrate host action is off for this project, so flai serve runs no orchestrator to %s: flai serve enable orchestrate turns it on, and starts it", verb)
}

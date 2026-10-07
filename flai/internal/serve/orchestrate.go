package serve

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator, run for a project while the orchestrate host action is
// on (S-0218).
//
// An orchestrator per served project looks whenever the project's launcher
// does: when its work items or threads change, when an agent ends, and every
// minute. With the action on and no run going, it starts the orchestrator in
// the project's main checkout with the project's orchestration agent
// (orchestration.agent over agent) and the harness's orchestrator request,
// as orchestrator, with FLAI_ROLE=orchestrate. The orchestrator waits on
// events and does not end by itself; when a run ends all the same, the next
// look starts another, and after a failure no sooner than orchestrateRetry
// later, so that a run that fails at once does not spin. With the action off
// and a run going, it stops the run as Stop stops a story's agent.
//
// Like the planner's, the run is recorded in serve/agents.json, under
// orchestrator, the newest alone; its log is the orchestrator's
// (Dir.ActivityLogs); and once it has ended, its activity since its last
// activity_log is logged in wip/agents/orchestrator.md. It is no story's
// agent: the in-progress limit does not count it, and it holds no story back.
//
// The operator may stop it and hold it stopped, and start it again, from its
// dashboard page or in a shell (S-0228, OrchestratorStop and
// OrchestratorStart). The hold is kept on its run in serve/agents.json:
// while it is held no look starts it, and turning the action off lifts it,
// so that turned on again the orchestrator starts.

// orchestrateRetry is how long after a failed run the orchestrator is
// started again.
const orchestrateRetry = time.Minute

// orchestrateOff is the activity of a run stopped because the orchestrate
// host action was turned off.
const orchestrateOff = "stopped: orchestrate turned off"

// orchestrateHeld is the activity of a run the operator stopped and held
// stopped (S-0228).
const orchestrateHeld = "stopped: stopped from the dashboard"

// orchestrator runs the orchestrator for one project while the orchestrate
// host action is on (S-0218).
type orchestrator struct {
	o       Options
	e       Entry
	starter *launcher

	mu sync.Mutex
	// said is the refusal to start last recorded, so that a look every
	// minute does not record it again; "" when the last start was not
	// refused.
	said string
	// seen is the run as the last look left it, its session and whether it
	// was held, so that a run started or held by a process of its own (flai
	// serve orchestrate start or stop) is told to the dashboard at the next
	// look (S-0228); "" before the first look.
	seen string
}

// newOrchestrator is the orchestrator for project e, starting its runs with
// starter, the project's launcher.
func newOrchestrator(o Options, e Entry, starter *launcher) *orchestrator {
	return &orchestrator{o: o, e: e, starter: starter}
}

// look starts the orchestrator while the action is on and no run is going,
// unless the operator holds it stopped (S-0228) or the last run failed less
// than orchestrateRetry ago, and stops the run going while the action is
// off, which also lifts the hold.
func (r *orchestrator) look(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg AgentConfig
	if r.o.Agent != nil {
		cfg = r.o.Agent(r.e.Root)
	}
	run := r.o.Dir.AgentStates()[r.e.Root].Orchestrator
	r.notice(run)
	switch {
	case !cfg.Orchestrate:
		r.said = ""
		if run.live() {
			r.stop(run)
		}
		r.unhold()
		return
	case run.live():
		return
	case run != nil && run.Held:
		return
	case run != nil && run.Outcome == OutcomeFailed && !r.due(run):
		return
	}
	repo, err := workitem.Open(r.e.Root)
	if err != nil {
		return
	}
	r.starter.mu.Lock()
	defer r.starter.mu.Unlock()
	r.said = r.starter.orchestrate(ctx, cfg, repo.Manifest.OrchestrationAgent(), r.said)
	r.seen = seenAs(r.o.Dir.AgentStates()[r.e.Root].Orchestrator) // told as it was recorded
}

// notice tells the dashboard of run when the last look left another, or the
// same one held otherwise: a process of its own, which tells no dashboard,
// started or held it since. The first look tells nothing.
func (r *orchestrator) notice(run *AgentRun) {
	seen := seenAs(run)
	if seen == r.seen {
		return
	}
	if r.seen != "" && run != nil && r.starter.changed != nil {
		r.starter.changed(run)
	}
	r.seen = seen
}

// seenAs is what notice compares of the orchestrator's run.
func seenAs(run *AgentRun) string {
	if run == nil {
		return "none"
	}
	return fmt.Sprintf("%s held=%t", run.Session, run.Held)
}

// unhold lifts the operator's hold on the orchestrator, if there is one,
// because the action was turned off (S-0228): turned on again, it starts the
// orchestrator. The dashboard is told, as of a run that ended.
func (r *orchestrator) unhold() {
	if run := r.o.Dir.AgentStates()[r.e.Root].Orchestrator; run == nil || !run.Held {
		return
	}
	var lifted *AgentRun
	r.o.Dir.updateAgent(r.e.Root, func(s *AgentState) {
		if cur := s.Orchestrator; cur != nil && cur.Held {
			c := *cur
			c.Held = false
			s.put(&c)
			lifted = &c
		}
	})
	if lifted == nil {
		return
	}
	r.seen = seenAs(lifted)
	r.starter.log("orchestrator no longer held stopped, since orchestrate was turned off", "pid", lifted.PID)
	if r.starter.changed != nil {
		r.starter.changed(lifted)
	}
}

// due says whether orchestrateRetry has passed since run ended; a run whose
// end cannot be read is due.
func (r *orchestrator) due(run *AgentRun) bool {
	ended, err := time.Parse(time.RFC3339, run.Ended)
	return err != nil || !r.o.Now().Before(ended.Add(orchestrateRetry))
}

// stop ends run, the orchestrator's, because the action was turned off, as
// Stop ends a story's agent: it is marked stopped first, its process group
// is asked to end while it is still the one flai started, and made to once
// o.StopGrace has passed. The launcher waiting for it settles its end; when
// none does, it is settled here. The stop is journalled.
func (r *orchestrator) stop(run *AgentRun) {
	l := r.starter
	at := r.o.Now().UTC().Format(time.RFC3339)
	// Marked before the signal, so that whoever settles the run records the
	// stop and not a failure.
	r.o.Dir.updateAgent(r.e.Root, func(s *AgentState) {
		if cur := s.Orchestrator; cur.same(run) && cur.live() {
			marked := *cur
			marked.Stopped = at
			s.put(&marked)
		}
	})
	signalled, killed := false, false
	if run.running() {
		signalled = true
		_ = TerminateGroup(run.PID)
		if !gone(run.PID, r.o.stopGrace()) {
			killed = true
			_ = KillGroup(run.PID)
			gone(run.PID, 2*time.Second)
		}
	}
	l.mu.Lock()
	if cur := r.o.Dir.AgentStates()[r.e.Root].Orchestrator; cur.same(run) && cur.live() && !l.waiting[run.PID] {
		l.orchestrateEnded(cur, nil)
	}
	l.mu.Unlock()
	what := fmt.Sprintf("recorded it as stopped: pid %d is no longer the orchestrator flai started", run.PID)
	switch {
	case killed:
		what = fmt.Sprintf("killed its process group (pid %d), which had not ended %s after it was asked to", run.PID, r.o.stopGrace())
	case signalled:
		what = fmt.Sprintf("ended its process group (pid %d)", run.PID)
	}
	if l.record != nil {
		l.record(hostapi.Entry{At: at, Action: hostapi.ActionOrchestrate, Method: "serve.orchestrate", Project: r.e.Key, Root: r.e.Root, By: "flai serve",
			Outcome: "done", Detail: fmt.Sprintf("stopped the orchestrator %s, since orchestrate was turned off: %s", run.Agent, what)})
	}
	l.log("orchestrator stopped", "pid", run.PID, "signalled", signalled, "killed", killed)
}

// orchestrate starts the orchestrator with agent, records the run, journals
// it, and returns why it was refused before it could be started, "" when it
// was not. It is refused when the agent names no harness and no command is
// set on the host, and when the harness refuses the request, as claude-code
// does without the project's orchestrator definition. A refusal that is
// said, the one returned last time, is not recorded, journalled, or logged
// again.
func (l *launcher) orchestrate(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, said string) string {
	now := l.now().UTC()
	run := &AgentRun{Agent: workitem.ActivityOrchestrator, Started: now.Format(time.RFC3339), Session: newSession()}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionOrchestrate, Method: "serve.orchestrate", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", "the orchestrator could not be started: "+run.Error
		if l.record != nil {
			l.record(entry)
		}
		l.warn("orchestrator could not be started", "harness", run.Harness, "err", run.Error)
	}
	refuse := func(err error) string {
		if err.Error() != said {
			fail(err)
		}
		return err.Error()
	}
	commandSet := cfg.host(harness.Command).Program != ""
	if (agent == nil || agent.Harness == "") && !commandSet {
		return refuse(errors.New("the orchestrator's agent names no harness (orchestration.agent or agent in system-flow.yaml), and no command is set on the host (flai serve agent set -- <program> [args...])"))
	}
	name, adapter, err := harness.For(agent, commandSet)
	run.Harness = name
	if err != nil {
		return refuse(err)
	}
	spec, err := adapter.Start(harness.Request{Role: conventions.RoleOrchestrate, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session}, cfg.host(name))
	if err != nil {
		return refuse(err)
	}
	cmd, out, err := l.spawn(ctx, run, spec.Argv, workitem.ActivityOrchestrator, now, spec.Env)
	if err != nil {
		fail(err)
		return ""
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) as %s (pid %d); log %s", run.Command, run.Harness, run.Agent, run.PID, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("orchestrator started", "harness", run.Harness, "command", run.Command, "pid", run.PID)
	l.await(cmd, out, run.PID, func(code int) { l.orchestrateEnded(run, &code) })
	return ""
}

// orchestrateEnded records an orchestrator run as ended, with its exit code
// when it was seen: failed on a failed exit, worked otherwise, an exit
// nobody saw included, and stopped when it was stopped. It journals the end,
// save a stop's, which the stop journals, and logs the run's activity since
// the last entry in the orchestrator's activity document, with its final
// reply as the summary, or, when it was stopped, orchestrateHeld when the
// operator held it stopped and orchestrateOff otherwise. An activity that
// cannot be logged is warned of, and the run stays recorded as it ended.
func (l *launcher) orchestrateEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	ended.Outcome, ended.Why = OutcomeWorked, ""
	if exit != nil && *exit != 0 {
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (exit %d)", *exit)
	}
	l.ended(&ended)
	args := []any{"pid", run.PID, "outcome", ended.Outcome}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	l.log("orchestrator ended", args...)
	summary := ""
	if ended.Outcome == OutcomeStopped {
		summary = orchestrateOff
		if ended.Held {
			summary = orchestrateHeld
		}
	} else if l.record != nil {
		e := hostapi.Entry{At: ended.Ended, Action: hostapi.ActionOrchestrate, Method: "serve.orchestrate", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve",
			Outcome: "done", Detail: fmt.Sprintf("the orchestrator %s (pid %d) ended, %s", run.Agent, run.PID, ended.Outcome)}
		if ended.Outcome == OutcomeFailed {
			e.Outcome, e.Detail = "failed", e.Detail+": "+ended.Why
		}
		l.record(e)
	}
	if _, err := logRunEnd(l.dir, l.entry.Root, l.entry.Key, workitem.ActivityOrchestrator, "", summary, nil); err != nil {
		l.warn("orchestrator activity not logged", "err", err)
	}
}

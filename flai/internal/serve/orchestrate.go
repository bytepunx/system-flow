package serve

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The orchestrator, run for a project while the orchestrate host action is
// on (S-0218), for reasons (ADR-0209).
//
// An orchestrator per served project looks whenever the project's launcher
// does: when its work items or threads change, when an agent ends, and every
// minute, and when the reasons it holds come due between those. At each look
// it reads the inbox's catch-up as the orchestrator would, under a cursor of
// flai serve's own (serveCursor), and queues as a reason each event a
// trigger in the orchestrator's role file matches (orchestratorTriggers,
// reasonsOf); orchestration.schedule coming round is one too. The queue is
// kept per project in serve/orchestrate-reasons.json, so that a flai serve
// started again loses none. With the action on, no run going, and the
// orchestrator not held, it starts one run once orchestration.coalesce has
// passed since the first reason was queued, or at once for the operator's
// asked, with every reason queued, in the project's main checkout with the
// project's orchestration agent (orchestration.agent over agent) and the
// harness's orchestrator request, as orchestrator, with FLAI_ROLE=orchestrate
// and the reasons in FLAI_REASONS. A start clears the reasons it took; one
// refused keeps them, and after a failed run the next is started no sooner
// than orchestrateRetry later, so that a run that fails at once does not
// spin. Reasons that come while a run goes wait for its end, and start the
// next run at the first look after it. With nothing queued, nothing starts:
// a run that ends is not started again by itself.
//
// The action turned on, or found on at flai serve's first look with no
// queue kept, queues one asked reason, and none of the catch-up's history.
// While it is off the cursor and the schedule keep up and nothing is
// queued, so that turning it on does not act on the past; and a run going
// is stopped as Stop stops a story's agent.
//
// Like the planner's, the run is recorded in serve/agents.json, under
// orchestrator, the newest alone; its log is the orchestrator's
// (Dir.ActivityLogs); and once it has ended, its activity since its last
// activity_log is logged in wip/agents/orchestrator.md. It is no story's
// agent: the in-progress limit does not count it, and it holds no story back.
//
// The operator may stop it and hold it stopped, and start it, from its
// dashboard page or in a shell (S-0228, OrchestratorStop and
// OrchestratorStart). The hold is kept on its run in serve/agents.json:
// while it is held no look starts it, and turning the action off lifts it.

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
	// said is the refusal to start last recorded, or why the budget state
	// held the last start, as logged, so that a look every minute does not
	// record or log it again; "" when the last start was neither.
	said string
	// seen is the run as the last look left it, its session and whether it
	// was held, so that a run started or held by a process of its own (flai
	// serve orchestrate start or stop) is told to the dashboard at the next
	// look (S-0228); "" before the first look.
	seen string
	// cadence follows orchestration.schedule.
	cadence cadence
	// warned are the failures to read the project warned of, so that a look
	// every minute does not warn of them again.
	warned map[string]bool
	// wake is sent on when the reasons queued come due between looks, and
	// timer sends it.
	wake  chan struct{}
	timer *time.Timer
}

// newOrchestrator is the orchestrator for project e, starting its runs with
// starter, the project's launcher. orchestration.schedule's times before now
// do not come round.
func newOrchestrator(o Options, e Entry, starter *launcher) *orchestrator {
	now := o.Now().UTC()
	r := &orchestrator{o: o, e: e, starter: starter, cadence: cadence{mark: now}, warned: map[string]bool{}, wake: make(chan struct{}, 1)}
	if repo, err := workitem.Open(e.Root); err == nil {
		r.cadence.see(repo.Manifest.Orchestration.Schedule, now)
	}
	return r
}

// queueBatches reports whether the merge queue batches for project e with
// manifest m, as flai serve's queue pass judges it: the merge-queue host
// action on, accept_reviews on, and local acceptance (ADR-0173, S-0403).
func queueBatches(o Options, e Entry, m manifest.Manifest) bool {
	return mergequeue.Batches(o.Host.Enabled != nil && o.Host.Enabled(hostapi.ActionMergeQueue, e.Root), m)
}

// look queues the reasons the events since the last look give while the
// action is on, and starts the orchestrator for those queued once they are
// due (reasonQueue.due) and no run is going, unless the operator holds it
// stopped (S-0228), the last run ended on the harness's usage limit and its
// provider's budget state holds (S-0534), or the last run failed otherwise
// less than orchestrateRetry ago. A start clears the reasons it took; while
// the budget state holds, they stay queued (ADR-0209 §5). A last run on the
// usage limit waits for the budget state alone, not orchestrateRetry: the
// first look once the state no longer holds starts the next. A start the
// budget state holds otherwise, as when another role's run set it, records
// nothing, so it is not a failure orchestrateRetry counts, and the next look
// tries again. While the action is off it keeps up with the events and the
// schedule and queues nothing, and stops the run going, which also lifts the
// hold.
func (r *orchestrator) look(ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg AgentConfig
	if r.o.Agent != nil {
		cfg = r.o.Agent(r.e.Root)
	}
	run := r.o.Dir.AgentStates()[r.e.Root].Orchestrator
	r.notice(run)
	now := r.o.Now().UTC()
	repo, err := workitem.Open(r.e.Root)
	if err != nil {
		r.warnOnce("orchestrator could not read the project", err)
	}
	if !cfg.Orchestrate {
		r.said = ""
		if run.live() {
			r.stop(run)
		}
		r.unhold()
		if repo != nil {
			r.keepUp(repo, now)
		}
		return
	}
	if repo == nil {
		return
	}
	q := r.gather(repo, now)
	switch {
	case run.live():
		return
	case run != nil && run.Held:
		return
	case run != nil && run.UsageLimited:
		if _, held := r.starter.budgetHolds(budgetKey(run.Provider, run.Harness)); held {
			return
		}
	case run != nil && run.Outcome == OutcomeFailed && !r.due(run):
		return
	}
	coalesce, err := repo.Manifest.Orchestration.CoalesceDuration()
	if err != nil {
		r.warnOnce("orchestration coalesce not valid", err)
		coalesce = manifest.DefaultCoalesce
	}
	ready, at := q.due(now, coalesce)
	if !ready {
		if !at.IsZero() {
			r.wakeIn(at.Sub(now))
		}
		return
	}
	r.starter.mu.Lock()
	defer r.starter.mu.Unlock()
	var started bool
	r.said, _, started = r.starter.orchestrate(ctx, cfg, repo.Manifest.OrchestrationAgent(), r.said, queueBatches(r.o, r.e, repo.Manifest), q.Reasons)
	if started {
		r.took(q.Reasons)
	}
	r.seen = seenAs(r.o.Dir.AgentStates()[r.e.Root].Orchestrator) // told as it was recorded
}

// gather reads the events since the last look as the orchestrator, under
// serveCursor, and queues those its triggers match, and
// orchestration.schedule when it has come round, and returns the queue. The
// action turned on, or on at the first look with no queue kept, queues an
// asked reason in their place: what came before is history.
func (r *orchestrator) gather(repo *workitem.Repo, now time.Time) reasonQueue {
	items, events := r.catchUp(repo, now)
	triggers := r.triggers(repo)
	sched := r.schedule(repo.Manifest.Orchestration, now)
	due := r.cadence.due(sched, now) && scheduled(triggers)
	queued := 0
	q, err := r.o.Dir.updateReasonQueue(r.e.Root, func(q *reasonQueue) {
		if !q.On {
			q.On = true
			queued = q.add(now, askedReason(askedTurnedOn, now))
			return
		}
		queued = q.add(now, reasonsOf(events, triggers, items)...)
		if due {
			queued += q.add(now, scheduleReason(sched, now))
		}
	})
	if err != nil {
		r.starter.warn("orchestrator reasons not kept", "err", err)
	}
	if queued > 0 {
		r.starter.log("orchestrator reasons queued", "queued", queued, "waiting", len(q.Reasons))
	}
	return q
}

// keepUp reads the events since the last look, while the action is off, so
// that turned on the orchestrator does not act on them, marks the schedule
// as come round, and empties the queue, the action off.
func (r *orchestrator) keepUp(repo *workitem.Repo, now time.Time) {
	r.catchUp(repo, now)
	r.cadence.see(repo.Manifest.Orchestration.Schedule, now)
	r.cadence.mark = now
	if _, err := r.o.Dir.updateReasonQueue(r.e.Root, func(q *reasonQueue) { q.On = false; q.clear() }); err != nil {
		r.starter.warn("orchestrator reasons not kept", "err", err)
	}
}

// catchUp is the project's items, archive included, and the events since
// serveCursor's last look as the orchestrator reads them: its own changes
// left out, and threads that await it or the operator. A failure is warned
// of once, and is no events.
func (r *orchestrator) catchUp(repo *workitem.Repo, now time.Time) ([]*workitem.Item, []inbox.Event) {
	items, err := repo.List(true)
	if err != nil {
		r.warnOnce("orchestrator could not read the project", err)
		return nil, nil
	}
	events, _, err := inbox.CatchUp(repo, inbox.Reader{Agent: workitem.ActivityOrchestrator, Cursor: serveCursor}, "", true, now, items)
	if err != nil {
		r.warnOnce("orchestrator could not read the project's events", err)
		return items, nil
	}
	return items, events
}

// triggers are the orchestrator's triggers, from its role file in repo's
// design folder or else builtinTriggers.
func (r *orchestrator) triggers(repo *workitem.Repo) []roles.Trigger {
	all, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		r.warnOnce("orchestrator role file not read", err)
	}
	return orchestratorTriggers(all)
}

// schedule is orchestration.schedule parsed, nil when it is unset or cannot
// be parsed, which is warned of once. A schedule set or changed comes round
// first at its next time from now.
func (r *orchestrator) schedule(o manifest.Orchestration, now time.Time) *cron.Schedule {
	r.cadence.see(o.Schedule, now)
	sched, err := o.OrchestrationSchedule()
	if err != nil {
		r.warnOnce("orchestration schedule not valid", err)
		return nil
	}
	return sched
}

// took takes the reasons a run was started with out of the queue.
func (r *orchestrator) took(reasons []orchestrateReason) {
	if _, err := r.o.Dir.updateReasonQueue(r.e.Root, func(q *reasonQueue) { q.remove(reasons) }); err != nil {
		r.starter.warn("orchestrator reasons not kept", "err", err)
	}
}

// wakeIn has the serving loop look again in d, when the reasons queued come
// due, in place of any earlier wake.
func (r *orchestrator) wakeIn(d time.Duration) {
	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(max(d, 0), func() {
		select {
		case r.wake <- struct{}{}:
		default:
		}
	})
}

// warnOnce warns of err with msg the first time it is seen.
func (r *orchestrator) warnOnce(msg string, err error) {
	if key := msg + "\x00" + err.Error(); !r.warned[key] {
		r.warned[key] = true
		r.starter.warn(msg, "err", err)
	}
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

// orchestrate starts the orchestrator with agent for reasons, records the
// run, journals it with its reasons, and returns why it was refused before it
// could be started, "" when it was not, and whether it was started. It is
// refused when the agent names no harness and no command is
// set on the host, when the harness refuses the request, as claude-code
// does without the project's orchestrator definition, and when the agent's
// provider cannot be had or its key is not set (ADR-0129). A refusal that is
// said, the one returned last time, is not recorded, journalled, or logged
// again. queue says the merge queue batches for the project, which the
// orchestrator's prompt then says, so that it enqueues the stories it has
// verified instead of accepting them (S-0403). The run is given the reasons'
// lines in FLAI_REASONS (harness.Request.Reasons). Its start, each retry's
// included, is an agent event on the orchestrator's role, saying why it was
// started and for which reasons, and so is a start it could not make, once
// as it is journalled (S-0532). A start the budget state of its provider
// holds is not made and nothing is recorded (spawnGated, S-0534): held is
// that state, nil when it was not held, and why it is held is returned as
// said, logged when it is not the one said last, so that a hold is logged
// once per state. started says whether a run was started, so that the
// reasons it took leave the queue.
func (l *launcher) orchestrate(ctx context.Context, cfg AgentConfig, agent *manifest.Agent, said string, queue bool, reasons []orchestrateReason) (string, *Budget, bool) {
	now := l.now().UTC()
	last := l.dir.AgentStates()[l.entry.Root].Orchestrator
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
		appendAgentEvent(l.entry.Root, run, inbox.AgentFailed, run.Why, l.warn)
	}
	refuse := func(err error) (string, *Budget, bool) {
		if err.Error() != said {
			fail(err)
		}
		return err.Error(), nil, false
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
	spec, keys, err := l.startOf(cfg, run, adapter, name, harness.Request{Role: conventions.RoleOrchestrate, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session, Queue: queue, Reasons: reasonLines(reasons)})
	if err != nil {
		return refuse(err)
	}
	cmd, out, held, err := l.spawnGated(ctx, run, spec.Argv, workitem.ActivityOrchestrator, now, append(spec.Env, keys...))
	if held != nil {
		why := budgetHeldWhy(*held, now)
		if why != said {
			l.log("orchestrator start held", "provider", held.Provider, "until", held.Until.UTC().Format(time.RFC3339), "reason", why)
		}
		return why, held, false
	}
	if err != nil {
		fail(err)
		return "", nil, false
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) as %s (pid %d) for %s; log %s", run.Command, run.Harness, run.Agent, run.PID, reasonsSaid(reasons), run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("orchestrator started", "harness", run.Harness, "command", run.Command, "pid", run.PID, "reasons", len(reasons))
	appendAgentEvent(l.entry.Root, run, inbox.AgentStarted, orchestrateWhy(last, reasons), l.warn)
	l.await(cmd, out, run.PID, func(code int) { l.orchestrateEnded(run, &code) })
	return "", nil, true
}

// orchestrateEnded records an orchestrator run as ended, with its exit code
// when it was seen: failed on a failed exit, worked otherwise, an exit
// nobody saw included, and stopped when it was stopped. It journals the end,
// save a stop's, which the stop journals, and logs the run's activity since
// the last entry in the orchestrator's activity document, with its final
// reply as the summary, or, when it was stopped, orchestrateHeld when the
// operator held it stopped and orchestrateOff otherwise. An activity that
// cannot be logged is warned of, and the run stays recorded as it ended. The
// end, a stop's included, is an agent event on the orchestrator's role,
// appended once, by the end that recorded the run as ended first (S-0532):
// a stop settled here and by the launcher that waits for it is told once.
// A run whose log ends on the harness's usage limit is failed on it, with
// why naming the limit, whether its exit was seen or not, and marked so
// (limitEnded), so that look holds the next start on its provider's budget
// state rather than orchestrateRetry; every end feeds that state once
// (budgetEnd, S-0534).
func (l *launcher) orchestrateEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	ended.Outcome, ended.Why = OutcomeWorked, ""
	lim, _ := runLimit(run, l.warn)
	ended.limitEnded(lim)
	switch {
	case lim != nil:
		ended.Outcome, ended.Why = OutcomeFailed, limitWhy(lim)
	case exit != nil && *exit != 0:
		ended.Outcome, ended.Why = OutcomeFailed, fmt.Sprintf("ended (exit %d)", *exit)
	}
	first := l.ended(&ended)
	if first {
		appendAgentEnd(l.entry.Root, &ended, l.warn)
	}
	l.budgetEnd(&ended, lim, first)
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

// orchestrateWhy is why the orchestrator was started, for its started agent
// event (S-0532): for reasons, named as reasonsSaid names them, and, told by
// last, the run recorded before it, the operator's start of one they held
// stopped, or a start after a run that failed, saying how it failed.
func orchestrateWhy(last *AgentRun, reasons []orchestrateReason) string {
	why := "started"
	switch {
	case last == nil:
	case last.Held:
		why = "the operator started it again, held stopped until now,"
	case last.Outcome == OutcomeFailed && last.Why != "":
		why = "started again after its last run failed (" + last.Why + "),"
	}
	return why + " for " + reasonsSaid(reasons)
}

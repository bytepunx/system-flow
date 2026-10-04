package serve

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The planner, started for an epic or a story on the operator's word
// (S-0208, ADR-0075, ADR-0079).
//
// flai plan, or the Plan button through plan.run, which runs it, starts the
// planner in the project's main checkout with the project's planning agent
// (planning.agent over agent) and the harness's planner request: it plans
// one item, writes work items and threads through flai, and ends. Like flai
// serve agent start it runs in a process of its own and records the run in
// serve/agents.json, by item, and its log beside the story agents' as the
// planner's (Dir.ActivityLogs); the serving flai settles the run once its
// process is gone, judges how it went, and logs its activity in
// wip/agents/planner.md. It is no story's agent: it is not counted against
// the in-progress limit, and holds no story back.

// planned matches the ID of an item the planner plans: an epic's or a
// story's.
var planned = regexp.MustCompile(`^[ES]-\d{3,}$`)

// Plan starts the planner for item, an epic or a story, in project e, and
// returns its run as recorded. It is refused while the plan host action is
// off for the project, for an ID that is not an epic's or a story's, for an
// item that is archived, done, or cancelled, while a planner runs for the
// item, and when nothing can start it. It does not wait for the planner:
// the serving flai settles the run once the process is gone.
func Plan(ctx context.Context, o Options, e Entry, item string) (*AgentRun, error) {
	cfg := o.Agent(e.Root)
	if !cfg.Plan {
		return nil, refused("the plan host action is off for this project: flai serve enable plan")
	}
	if strings.HasPrefix(item, "T-") {
		return nil, refused("%s is a task; the planner plans an epic or a story", item)
	}
	if !planned.MatchString(item) {
		return nil, refused("%q is not an epic's or a story's ID; the planner plans an E-nnnn or an S-nnnn", item)
	}
	repo, err := workitem.Open(e.Root)
	if err != nil {
		return nil, err
	}
	it, err := repo.Get(item)
	if err != nil {
		return nil, err
	}
	switch {
	case it.Archived:
		return nil, refused("%s is archived; the planner plans an epic or a story still on the board", it.ID)
	case it.Status == workitem.Done || it.Status == workitem.Cancelled:
		return nil, refused("%s is %s; there is nothing left to plan", it.ID, it.Status)
	}
	if run := o.Dir.AgentStates()[e.Root].Plans[it.ID]; run.running() {
		return nil, refused("the planner is already running for %s (pid %d, started %s); one item has one planner at a time", it.ID, run.PID, run.Started)
	}
	agent := repo.Manifest.PlanningAgent()
	if (agent == nil || agent.Harness == "") && cfg.host(harness.Command).Program == "" {
		return nil, refused("the planner's agent names no harness (planning.agent or agent in system-flow.yaml), and no command is set on the host (flai serve agent set -- <program> [args...])")
	}
	l := newLauncher(o, e)
	l.handOver = true
	l.plan(ctx, cfg, it.ID, agent)
	run := o.Dir.AgentStates()[e.Root].Plans[it.ID]
	if run == nil {
		return nil, fmt.Errorf("the planner run for %s was not recorded in %s", it.ID, o.Dir.agents())
	}
	if run.Error != "" {
		return run, errors.New(run.Why)
	}
	return run, nil
}

// plan starts the planner for item with agent, records the run, journals
// it, and says whether it started.
func (l *launcher) plan(ctx context.Context, cfg AgentConfig, item string, agent *manifest.Agent) bool {
	now := l.now().UTC()
	run := &AgentRun{Item: item, Agent: workitem.ActivityPlanner + "-" + item, Started: now.Format(time.RFC3339), Session: newSession()}
	if agent != nil {
		run.Model = agent.Model
	}
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionPlan, Method: "serve.plan", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	fail := func(err error) bool {
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("the planner for %s could not be started: %s", item, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("planner could not be started", "item", item, "harness", run.Harness, "err", run.Error)
		return false
	}
	name, adapter, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	spec, err := adapter.Start(harness.Request{Role: conventions.RolePlan, Item: item, Root: l.entry.Root, Project: l.entry.Key, Agent: agent,
		Name: run.Agent, Flai: cfg.Flai, Session: run.Session}, cfg.host(name))
	if err != nil {
		return fail(err)
	}
	cmd, out, err := l.spawn(ctx, run, spec.Argv, workitem.ActivityPlanner, now, spec.Env)
	if err != nil {
		return fail(err)
	}
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) to plan %s as %s (pid %d); log %s", run.Command, run.Harness, item, run.Agent, run.PID, run.Log)
	if l.record != nil {
		l.record(entry)
	}
	l.log("planner started", "item", item, "harness", run.Harness, "command", run.Command, "pid", run.PID)
	l.await(cmd, out, run.PID, func(code int) { l.planEnded(run, &code) })
	return true
}

// planEnded records a planner run as ended, with its exit code when it was
// seen, as judgePlan judges it, and logs its activity in the planner's
// activity document (ADR-0079), naming the items it planned. An activity
// that cannot be logged is warned of, and the run stays recorded as it
// ended; items that cannot be read are warned of, and the activity names
// the planned item alone.
func (l *launcher) planEnded(run *AgentRun, exit *int) {
	ended := *run
	ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), exit
	ended.Outcome, ended.Why, ended.Thread = judgePlan(l.entry.Root, run.Item, run.Agent, exit)
	l.ended(&ended)
	args := []any{"item", run.Item, "pid", run.PID, "outcome", ended.Outcome}
	if exit != nil {
		args = append(args, "exit", *exit)
	}
	l.log("planner ended", args...)
	items := []string{run.Item}
	if since, err := time.Parse(time.RFC3339, run.Started); err != nil {
		l.warn("planner's items not named", "item", run.Item, "err", err)
	} else if items, err = plannedItems(l.entry.Root, run.Item, since); err != nil {
		items = []string{run.Item}
		l.warn("planner's items not named", "item", run.Item, "err", err)
	}
	if _, err := LogRunEnd(l.dir, l.entry.Root, l.entry.Key, workitem.ActivityPlanner, items); err != nil {
		l.warn("planner activity not logged", "item", run.Item, "err", err)
	}
}

// plannedItems are the items a planner run on item, started at since, is
// said to have planned (S-0209): item first, then the items under it created
// since, then those under it created before and changed since, each in ID
// order. Under an epic are its stories and their tasks; under a story, its
// tasks. Archived items are not under it. The second since started counts
// as since.
func plannedItems(root, item string, since time.Time) ([]string, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	all, err := repo.List(false)
	if err != nil {
		return nil, err
	}
	since = since.UTC().Truncate(time.Second)
	after := func(at string) bool {
		t, err := time.Parse(workitem.TimeFormat, at)
		return err == nil && !t.Before(since)
	}
	under := map[string]bool{item: true}
	var created, changed []string
	// List puts stories before tasks, so a story is known to be under the
	// epic before its tasks are seen.
	for _, it := range all {
		if it.ID == item || !under[it.Parent] {
			continue
		}
		under[it.ID] = true
		switch {
		case after(it.Created):
			created = append(created, it.ID)
		case after(it.Updated):
			changed = append(changed, it.ID)
		}
	}
	return append(append([]string{item}, created...), changed...), nil
}

// judgePlan is how a planner run that has ended went: asked when a question
// of its own on its item is open and unanswered, failed on a failed exit,
// and worked otherwise, an exit nobody saw included, since a planner leaves
// its item in no state that tells.
func judgePlan(root, item, agent string, exit *int) (outcome, why, thread string) {
	code := "an exit code nobody saw"
	if exit != nil {
		code = fmt.Sprintf("exit %d", *exit)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the project could not be read: %v", code, err), ""
	}
	on, err := threads.For(repo, item)
	if err != nil {
		return OutcomeFailed, fmt.Sprintf("ended (%s); the threads on %s could not be read: %v", code, item, err), ""
	}
	for _, th := range on {
		if e := th.Entries(); th.Open() && len(e) > 0 && e[len(e)-1].Author == agent {
			return OutcomeAsked, "waiting for an answer to " + th.ID + ": " + th.Title, th.ID
		}
	}
	if exit != nil && *exit != 0 {
		return OutcomeFailed, fmt.Sprintf("ended (%s) planning %s", code, item), ""
	}
	return OutcomeWorked, "", ""
}

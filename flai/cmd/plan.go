package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai plan: the planner for an epic or a story, on the operator's word
// (S-0208).
func newPlanCmd(a *app) *cobra.Command {
	var candidates bool
	c := &cobra.Command{
		Use:   "plan <epic-or-story-id> | plan --candidates",
		Short: "Start the planner for an epic or a story: it drafts and enriches the item's stories or tasks through flai",
		Long: `Starts the planner for an epic or a story, now, on this host and as you
(S-0208, ADR-0075). It runs in the project's main checkout with the
project's planning agent: planning.agent in system-flow.yaml over the
project's agent, started with the harnesses and the command you set with
flai serve agent. For an epic with no stories it drafts the stories that
deliver its outcome, enriches each one, and drafts each one's tasks
(S-0300); for a story it adds touches, a forecast, and a cost of delay, and
drafts its tasks or revisits those it has (S-0255); for an epic with
stories it revisits each one not done or cancelled, drafts what the outcome
still lacks, and drafts the tasks of each story it adds and of each draft
it revisits that has none. It writes work items and threads through flai
alone, moves nothing past backlog, finalizes no draft, and
asks you on a thread on the item when an input of yours is missing.

The run is recorded where flai serve tracks agents, by item, and its output
goes to a log beside flai serve's state. Once it ends, the serving flai
records how (worked, asked, or failed) and logs what it did and cost in
wip/agents/planner.md. It is no story's agent: the in-progress limit does
not count it, and it holds no story back.

It is a host action, off until you enable it (flai serve enable plan). It
refuses, and says why, while the action is off for the project, for a task
or an ID that is neither an epic's nor a story's, for an item that is
archived, done, or cancelled, while a planner runs for the item, and when
nothing can start it. The Plan button on an epic's or a story's page runs
this.

--candidates lists, in ID order and with the reason for each, the epics the
planner should plan (S-0219): an epic in the backlog with no story, and an
epic not done or cancelled whose stories, archived ones included, are all
done or cancelled with at least one done. After them it lists the stories
(S-0328): a story in the backlog, draft or not, that has no touches, no
forecast duration, no cost of delay value, or no task that is not
cancelled, each with what it lacks. It leaves out, and lists apart with
why, an item a planner runs for now and one whose newest planner run ended
asking a question still awaiting the operator, as flai serve's record of
planner runs on this host says; and a story whose epic a planner runs for
now, or whose newest planner run ended and which has not changed since. It
starts nothing and writes nothing.

The orchestrator (FLAI_ROLE=orchestrate) asks for the planner on an epic
that --candidates lists alone, while orchestration.permissions gives it
plan_backlog_epics, and on a story that --candidates lists alone, while it
gives it plan_backlog_stories (S-0328, ADR-0119); its run records
orchestrator, in place of asked, as what started it (S-0219).`,
		Example: `  flai serve enable plan
  flai plan E-0016
  flai plan S-0208 --json
  flai plan --candidates --json`,
		Args: func(cmd *cobra.Command, args []string) error {
			if candidates {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(_ *cobra.Command, args []string) error {
			if candidates {
				return a.planCandidates()
			}
			return a.planNow(workitem.CanonicalID(args[0]))
		},
	}
	c.Flags().BoolVar(&candidates, "candidates", false, "list the epics and stories the planner should plan, each with why, and those left out while a planner runs or awaits the operator; writes nothing")
	return c
}

// planCandidates prints the epics and stories the planner should plan in
// the project, leaving out those flai serve's planner runs on this host keep
// back.
func (a *app) planCandidates() error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	got, err := a.planCandidatesOf(repo)
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(got)
	}
	if len(got.Candidates) == 0 {
		fmt.Fprintln(a.out, "nothing to plan")
	} else {
		fmt.Fprintln(a.out, "to plan:")
	}
	for _, c := range got.Candidates {
		fmt.Fprintf(a.out, "  %s  %s\n    - %s\n", c.ID, c.Title, c.Reason)
	}
	if len(got.LeftOut) > 0 {
		fmt.Fprintln(a.out, "left out:")
	}
	for _, c := range got.LeftOut {
		fmt.Fprintf(a.out, "  %s  %s\n    - %s\n", c.ID, c.Title, c.Reason)
	}
	return nil
}

// planCandidatesOf is what flai plan --candidates lists in repo's project,
// judged with flai serve's planner runs on this host.
func (a *app) planCandidatesOf(repo *workitem.Repo) (workitem.PlanCandidates, error) {
	items, err := repo.List(true)
	if err != nil {
		return workitem.PlanCandidates{}, err
	}
	return workitem.PlanCandidatesOf(items, planHeld(repo, items, a.serveDir().AgentStates()[mainRootOf(repo)])), nil
}

// planHeld is why each epic and story st's planner runs keep from being a
// candidate, judged among items: a planner runs for it now, or its newest
// planner run ended asking on a thread whose last word is still the
// planner's. A story is kept back too while a planner runs for its epic, and
// when its newest planner run has ended and the story has not been updated
// since, so that a planner that could not finish is not started again until
// something changes (S-0328, ADR-0119). A run that has not ended but whose
// process is gone is neither running nor ended: flai serve settles it at its
// next look.
func planHeld(repo *workitem.Repo, items []*workitem.Item, st serve.AgentState) map[string]string {
	out := map[string]string{}
	byID := map[string]*workitem.Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	running := map[string]string{}
	for id, run := range st.Plans {
		typ := workitem.TypeOfID(id)
		if run == nil || (typ != workitem.Epic && typ != workitem.Story) {
			continue
		}
		if run.Ended == "" && run.Error == "" && run.PID > 0 && serve.Owns(run.PID, run.Start) {
			running[id] = fmt.Sprintf("(pid %d, started %s)", run.PID, run.Started)
			out[id] = "a planner runs for it now " + running[id]
			continue
		}
		if asked := plannerAsked(repo, run); asked != "" {
			out[id] = asked
			continue
		}
		if story := byID[id]; typ == workitem.Story && story != nil && run.Ended != "" && !updatedAfter(story, run.Ended) {
			out[id] = fmt.Sprintf("its planner ran until %s and it has not changed since: a planner is not started on it again until it does", run.Ended)
		}
	}
	for _, it := range items {
		if _, held := out[it.ID]; held || it.Type != workitem.Story {
			continue
		}
		if pid, ok := running[it.Parent]; ok {
			out[it.ID] = fmt.Sprintf("a planner runs for its epic %s now %s", it.Parent, pid)
		}
	}
	return out
}

// plannerAsked says that run ended asking on a thread whose last word is
// still the planner's, or nothing when it did not.
func plannerAsked(repo *workitem.Repo, run *serve.AgentRun) string {
	if run.Outcome != serve.OutcomeAsked || run.Thread == "" {
		return ""
	}
	th, err := threads.Get(repo, run.Thread)
	if err != nil {
		return ""
	}
	if e := th.Entries(); th.Open() && len(e) > 0 && e[len(e)-1].Author == run.Agent {
		return fmt.Sprintf("its planner asked on %s, which awaits the operator: %s", th.ID, th.Title)
	}
	return ""
}

// updatedAfter says whether it was updated after ended, a run's end. An
// item or an end whose time does not read counts as updated, so that it is
// not held for want of a time.
func updatedAfter(it *workitem.Item, ended string) bool {
	end, err := time.Parse(time.RFC3339, ended)
	if err != nil {
		return true
	}
	updated, err := time.Parse(time.RFC3339, it.Updated)
	if err != nil {
		return true
	}
	return updated.After(end)
}

// planNow starts the planner for item and prints the run, and says on
// standard error when no flai serve is running to settle it. The
// orchestrator's shell is held as its MCP tool plan is (orchestratorPlans),
// and its run records the orchestrator as what started it.
func (a *app) planNow(item string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	start := serve.Plan
	if os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate {
		if err := a.orchestratorPlans(repo, item); err != nil {
			return fmt.Errorf("rule: %w", err)
		}
		start = serve.PlanForOrchestrator
	}
	o, e := a.planOn(repo)
	run, err := start(context.Background(), o, e, item)
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if _, running := o.Dir.ReadStatus(a.now()); !running {
		fmt.Fprintln(a.errOut, "flai serve is not running (flai serve start, or flai host start): the planner runs, and its end is recorded and its activity logged once flai serve is")
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"item": run.Item, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "session": run.Session, "started": run.Started})
	}
	fmt.Fprintf(a.out, "started %s (%s) to plan %s as %s (pid %d); log %s\n", run.Command, run.Harness, run.Item, run.Agent, run.PID, run.Log)
	return nil
}

// planOn is what the planner starts with in repo's project, for flai plan
// and flai mcp's plan alike: flai serve's options on this host, and the
// project in its main checkout.
func (a *app) planOn(repo *workitem.Repo) (serve.Options, serve.Entry) {
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	return o, serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
}

// mcpPlan starts the planner for item for flai mcp's tool plan (S-0208), as
// flai plan does, under the same plan host action, and journals who asked
// and what came of it. Of the agents flai serve starts, the orchestrator
// alone may ask, for an epic flai plan --candidates lists, while the project
// gives it plan_backlog_epics (S-0218, S-0219), and for a story it lists,
// while the project gives it plan_backlog_stories (S-0328), and its run
// records the orchestrator as what started it.
func (a *app) mcpPlan(ctx context.Context, root, item, by string) (mcpserver.PlanStarted, error) {
	orchestrator := os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" && !orchestrator {
		return mcpserver.PlanStarted{}, fmt.Errorf("an agent flai serve started does not start the planner: planning %s is the operator's to ask for, from its page or with flai plan %s on the host", item, item)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.PlanStarted{}, err
	}
	start := serve.Plan
	if orchestrator {
		if err := a.orchestratorPlans(repo, item); err != nil {
			return mcpserver.PlanStarted{}, err
		}
		by = orchestratorBy(by)
		start = serve.PlanForOrchestrator
	}
	o, e := a.planOn(repo)
	run, err := start(ctx, o, e, item)
	entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: hostapi.ActionPlan, Method: "mcp.plan", Project: e.Key, Root: e.Root, By: by, Outcome: "done"}
	var no *serve.Refused
	switch {
	case errors.As(err, &no) && !o.Agent(e.Root).Plan:
		entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s asked to plan %s: %s", by, item, no.Why)
	case err != nil:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s asked to plan %s: %s", by, item, err)
	default:
		entry.Detail = fmt.Sprintf("%s asked to plan %s: started %s as %s (pid %d)", by, item, run.Command, run.Agent, run.PID)
	}
	o.Host.Record(entry)
	if err != nil {
		return mcpserver.PlanStarted{}, err
	}
	return mcpserver.PlanStarted{Item: run.Item, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID, Log: run.Log, Session: run.Session, Started: run.Started}, nil
}

// orchestratorPlans says why the orchestrator may not ask for the planner on
// item in repo, or nil when it may: item is an epic flai plan --candidates
// lists, judged with flai serve's planner runs on this host, and the project
// gives the orchestrator plan_backlog_epics (S-0219); or it is a story the
// candidates list, and the project gives it plan_backlog_stories (S-0328,
// ADR-0119).
func (a *app) orchestratorPlans(repo *workitem.Repo, item string) error {
	var what, alone string
	switch workitem.TypeOfID(item) {
	case workitem.Epic:
		if !repo.Manifest.Orchestration.Permissions.Allows(manifest.PermitPlanBacklogEpics) {
			return fmt.Errorf("the orchestrator asks for the planner only with orchestration.permissions.plan_backlog_epics, which is off: ask the operator with thread_open on %s", item)
		}
		what = "is neither in the backlog with no stories nor open with every story done or cancelled and one done"
		alone = "the orchestrator asks for the planner on an epic flai plan --candidates lists alone, and "
	case workitem.Story:
		if !repo.Manifest.Orchestration.Permissions.Allows(manifest.PermitPlanBacklogStories) {
			return fmt.Errorf("the orchestrator asks for the planner on a story only with orchestration.permissions.plan_backlog_stories, which is off: ask the operator with thread_open on %s", item)
		}
		what = "has a plan: touches, a forecast duration, a cost of delay value, and a task"
		alone = "the orchestrator asks for the planner on a story flai plan --candidates lists alone, and "
	default:
		return fmt.Errorf("the orchestrator asks for the planner on an epic or a story flai plan --candidates lists alone, and %s is neither", item)
	}
	it, err := repo.Get(item)
	if err != nil {
		return err
	}
	got, err := a.planCandidatesOf(repo)
	if err != nil {
		return err
	}
	for _, c := range got.Candidates {
		if c.ID == it.ID {
			return nil
		}
	}
	for _, c := range got.LeftOut {
		if c.ID == it.ID {
			return fmt.Errorf("%s%s is left out: %s", alone, it.ID, c.Reason)
		}
	}
	switch {
	case it.Archived:
		what = "is archived"
	case it.Closed():
		what = "is " + it.Status
	case it.Type == workitem.Story && it.Status != workitem.Backlog:
		what = "is " + it.Status + ", not in the backlog"
	}
	return fmt.Errorf("%s%s is not one: it %s", alone, it.ID, what)
}

// orchestratorBy names the orchestrator as who asks, with the name flai mcp
// knows it by when that does not already.
func orchestratorBy(by string) string {
	switch {
	case by == "":
		return workitem.ActivityOrchestrator
	case strings.HasPrefix(by, workitem.ActivityOrchestrator):
		return by
	}
	return fmt.Sprintf("%s (%s)", workitem.ActivityOrchestrator, by)
}

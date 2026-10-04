package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai plan: the planner for an epic or a story, on the operator's word
// (S-0208).
func newPlanCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "plan <epic-or-story-id>",
		Short: "Start the planner for an epic or a story: it drafts and enriches the item's stories through flai",
		Long: `Starts the planner for an epic or a story, now, on this host and as you
(S-0208, ADR-0075). It runs in the project's main checkout with the
project's planning agent: planning.agent in system-flow.yaml over the
project's agent, started with the harnesses and the command you set with
flai serve agent. For an epic with no stories it drafts the stories that
deliver its outcome; for a story it adds touches, a forecast, and a cost of
delay; for an epic with stories it revisits each one not done or cancelled
and drafts what the outcome still lacks. It writes work items and threads
through flai alone, moves nothing past backlog, finalizes no draft, and
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
this.`,
		Example: `  flai serve enable plan
  flai plan E-0016
  flai plan S-0208 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.planNow(workitem.CanonicalID(args[0]))
		},
	}
}

// planNow starts the planner for item and prints the run, and says on
// standard error when no flai serve is running to settle it.
func (a *app) planNow(item string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	run, err := serve.Plan(context.Background(), o, e, item)
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

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai curate: the curator for the project, on the operator's word (S-0431,
// ADR-0196).
func newCurateCmd(a *app) *cobra.Command {
	var since string
	c := &cobra.Command{
		Use:   "curate",
		Short: "Start the curator for the project: it takes in what the archive and the threads gained as digest entries, through flai digest alone",
		Long: `Starts the curator for the project, now, on this host and as you
(S-0431, ADR-0196). It runs in the project's main checkout with the
project's digest agent: digest.agent in system-flow.yaml over the project's
agent, started with the harnesses and the command you set with flai serve
agent. It reads flai digest harvest from --since, a YYYY-MM-DD day in UTC or
an RFC 3339 time, or without one from the curator's last activity in
wip/agents/curator.md, or everything on its first run. For each candidate
it writes a digest entry with flai digest new, or merges it into the entry
that holds the fact with flai digest edit. It edits no file and authors no
work item; flai guard holds it to that.

The run is recorded where flai serve tracks agents, as the curator's, the
newest alone, and its output goes to a log beside flai serve's state. Once
it ends, the serving flai records how (worked or failed) and logs the
entries it wrote or changed and what the run cost in wip/agents/curator.md.
It is no story's agent: the in-progress limit does not count it, and it
holds no story back. One curator runs for a project at a time.

A run whose log ends on the harness's usage limit is recorded failed, with
the limit and when it resets, and logs no activity (ADR-0214). Until the
reset flai serve starts no curator on its own and holds the acceptances
that come meanwhile; at its first look after the reset it starts that run
again, once, with them. A retry that ends on the limit again, or a limit
with no reset time, is raised on a thread on system-flow.yaml instead.
This command is never held.

It is a host action, off until you enable it (flai serve enable curate).
It refuses, and says why, while the action is off for the project, for a
--since in neither form, while a curator runs for the project, naming it,
and when nothing can start it. The host API's curate.run runs this.`,
		Example: `  flai serve enable curate
  flai curate
  flai curate --since 2026-10-01 --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			return a.curateNow(since)
		},
	}
	c.Flags().StringVar(&since, "since", "", "take in what was gained since this YYYY-MM-DD day in UTC or RFC 3339 time; since the curator's last activity, or everything on its first run, when not given")
	return c
}

// curateNow starts the curator for the project, taking in what was gained
// since since, and prints the run, and says on standard error when no flai
// serve is running to settle it, as flai analyze does.
func (a *app) curateNow(since string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o, e := a.planOn(repo)
	run, err := serve.Curate(context.Background(), o, e, since, serve.TriggerAsked)
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if _, running := o.Dir.ReadStatus(a.now()); !running {
		fmt.Fprintln(a.errOut, "flai serve is not running (flai serve start, or flai host start): the curator runs, and its end is recorded and its entries and cost logged once flai serve is")
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"since": run.Since, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "session": run.Session, "started": run.Started, "trigger": run.Trigger})
	}
	fmt.Fprintf(a.out, "started %s (%s) to curate, %s, as %s (pid %d); log %s\n", run.Command, run.Harness, curatedSince(run.Since), run.Agent, run.PID, run.Log)
	return nil
}

// curatedSince says what a curator run takes in, as its journal entry says
// it: since its since, or everything on its first run.
func curatedSince(since string) string {
	if since == "" {
		return "everything, its first run"
	}
	return "since " + since
}

// mcpCurate starts the curator for flai mcp's tool curate (S-0431), as flai
// curate does, under the same curate host action, and journals who asked
// and what came of it. No agent flai serve started may ask: the curator is
// the operator's to start.
func (a *app) mcpCurate(ctx context.Context, root, since, by string) (mcpserver.CurateStarted, error) {
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" {
		return mcpserver.CurateStarted{}, errors.New("an agent flai serve started does not start the curator: it is the operator's to ask for, from the dashboard or with flai curate on the host")
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.CurateStarted{}, err
	}
	o, e := a.planOn(repo)
	run, err := serve.Curate(ctx, o, e, since, serve.TriggerAsked)
	asked := curatedSince(since)
	if since == "" {
		asked = "since its last run"
	}
	entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: hostapi.ActionCurate, Method: "mcp.curate", Project: e.Key, Root: e.Root, By: by, Outcome: "done"}
	var no *serve.Refused
	switch {
	case errors.As(err, &no) && !o.Agent(e.Root).Curate:
		entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s asked to curate, %s: %s", by, asked, no.Why)
	case err != nil:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s asked to curate, %s: %s", by, asked, err)
	default:
		entry.Detail = fmt.Sprintf("%s asked to curate, %s: started %s as %s (pid %d)", by, curatedSince(run.Since), run.Command, run.Agent, run.PID)
	}
	o.Host.Record(entry)
	if err != nil {
		return mcpserver.CurateStarted{}, err
	}
	return mcpserver.CurateStarted{Since: run.Since, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID, Log: run.Log, Session: run.Session, Started: run.Started, Trigger: run.Trigger}, nil
}

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai analyze: the analyzer for the project, on the operator's word
// (S-0223).
func newAnalyzeCmd(a *app) *cobra.Command {
	var focus string
	c := &cobra.Command{
		Use:   "analyze",
		Short: "Start the analyzer for the project: it looks for bottlenecks, intent, or risk and writes one report under design/analysis",
		Long: `Starts the analyzer for the project, now, on this host and as you
(S-0223). It runs in the project's main checkout with the project's
analysis agent: analysis.agent in system-flow.yaml over the project's
agent, started with the harnesses and the command you set with flai serve
agent. It reads the metrics, the design, the code, and the issues, looks
for what --focus names, or for all of it without one: bottlenecks in the
flow, gaps between design/system and the code (intent), and technical and
security risks (risk). It writes one report under design/analysis, edits
nothing else, and authors no stories.

The run is recorded where flai serve tracks agents, as the analyzer's, the
newest alone, and its output goes to a log beside flai serve's state. Once
it ends, the serving flai records how (worked or failed) and logs the
report it wrote and what the run cost in wip/agents/analyzer.md. It is no
story's agent: the in-progress limit does not count it, and it holds no
story back. One analyzer runs for a project at a time.

It is a host action, off until you enable it (flai serve enable analyze).
It refuses, and says why, while the action is off for the project, for a
focus other than bottlenecks, intent, or risk, while an analyzer runs for
the project, naming it, and when nothing can start it. The host API's
analyze.run runs this, and analysis.schedule in system-flow.yaml starts it
on a schedule while flai serve runs.`,
		Example: `  flai serve enable analyze
  flai analyze
  flai analyze --focus risk --json`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			return a.analyzeNow(focus)
		},
	}
	c.Flags().StringVar(&focus, "focus", "", "what to look for: bottlenecks, intent, or risk; all three when not given")
	return c
}

// analyzeNow starts the analyzer for the project, looking for focus, and
// prints the run, and says on standard error when no flai serve is running
// to settle it, as flai plan does.
func (a *app) analyzeNow(focus string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o, e := a.planOn(repo)
	run, err := serve.Analyze(context.Background(), o, e, focus, serve.TriggerAsked)
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if _, running := o.Dir.ReadStatus(a.now()); !running {
		fmt.Fprintln(a.errOut, "flai serve is not running (flai serve start, or flai host start): the analyzer runs, and its end is recorded and its report and cost logged once flai serve is")
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"focus": run.Focus, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "session": run.Session, "started": run.Started, "trigger": run.Trigger})
	}
	fmt.Fprintf(a.out, "started %s (%s) to analyze, focus %s, as %s (pid %d); log %s\n", run.Command, run.Harness, run.Focus, run.Agent, run.PID, run.Log)
	return nil
}

// mcpAnalyze starts the analyzer for flai mcp's tool analyze (S-0223), as
// flai analyze does, under the same analyze host action, and journals who
// asked and what came of it. No agent flai serve started may ask: the
// analyzer is the operator's to start, or the schedule's.
func (a *app) mcpAnalyze(ctx context.Context, root, focus, by string) (mcpserver.AnalyzeStarted, error) {
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" {
		return mcpserver.AnalyzeStarted{}, errors.New("an agent flai serve started does not start the analyzer: it is the operator's to ask for, from the dashboard or with flai analyze on the host, or analysis.schedule's")
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.AnalyzeStarted{}, err
	}
	o, e := a.planOn(repo)
	run, err := serve.Analyze(ctx, o, e, focus, serve.TriggerAsked)
	looking := focus
	if looking == "" {
		looking = harness.AllFocus
	}
	entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: hostapi.ActionAnalyze, Method: "mcp.analyze", Project: e.Key, Root: e.Root, By: by, Outcome: "done"}
	var no *serve.Refused
	switch {
	case errors.As(err, &no) && !o.Agent(e.Root).Analyze:
		entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s asked to analyze, focus %s: %s", by, looking, no.Why)
	case err != nil:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s asked to analyze, focus %s: %s", by, looking, err)
	default:
		entry.Detail = fmt.Sprintf("%s asked to analyze, focus %s: started %s as %s (pid %d)", by, run.Focus, run.Command, run.Agent, run.PID)
	}
	o.Host.Record(entry)
	if err != nil {
		return mcpserver.AnalyzeStarted{}, err
	}
	return mcpserver.AnalyzeStarted{Focus: run.Focus, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID, Log: run.Log, Session: run.Session, Started: run.Started, Trigger: run.Trigger}, nil
}

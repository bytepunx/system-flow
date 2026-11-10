package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newServeAgentRestartCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "restart <story-id>",
		Short: "Start a new agent for a story in ready or in progress whose agent dropped or failed, or that has none here",
		Long: `Starts a new agent, in a new session, for a story in ready or
in-progress whose last agent flai serve started has ended or dropped, the
way flai serve starts one when a story enters ready: the story's harness,
model, and options, or the host's command. The agent is told how its last
one ended and to go on from the story's narrative. The run is recorded
where the serving flai tracks it: the dot on the card, the outcome, the
restart on an answer (S-0116, ADR-0043).

A story that this host's flai serve has started no agent for gets one too
(S-0177, ADR-0064). One in progress, begun on another host, gets an agent
told who began it, when, and on which host, which threads have been written
to since, and to reconcile with what is committed rather than start over.

For a story in ready while the in-progress limit is full, review is full,
or a claim holds it, it queues the new agent instead: flai serve starts it as soon as the
story can start, as it starts a story that enters ready, and until then the
story's agent reads as waiting, queued (S-0118).

It refuses, and says why, while the agent action is off for the project,
when the story is in another state, while its agent runs or waits for an
answer, when one is already queued, and when nothing can start it. The story
page's Retry button runs this, and its Start agent button for a story in
progress that has no agent here.`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Restart)
		},
	}
}

func newServeAgentStartCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "start <story-id>",
		Short: "Start a ready story's agent now, whatever flai serve's own rules say about when",
		Long: `Starts the agent of a story in ready now, the way flai serve starts one
when a story enters ready: the story's harness, model, and options, the
project's default under them, the host's command when neither names a
harness, and flai's built-in default, claude-code, when no command is set
either (ADR-0261). It starts whether or not the story was ready before flai
serve began, whether or not it has had an agent since it entered ready, and
whoever is attending the project. The run is recorded where the serving
flai tracks it: the dot on the card, the outcome when it ends, the start
again on an answer (S-0115).

It refuses, and says why, while the agent action is off for the project,
when the story is not in ready, while its agent runs or waits for an
answer, and when its harness cannot start: a harness named command with no
command set for the project, or one flai has no adapter for, each refusal
naming the dashboard's settings page that fixes it. A full
in-progress limit or a full review does not stop it: it starts past them,
with a warning, as a move does. A story in progress whose agent dropped or failed,
or that has had no agent here, gets one from flai serve agent restart. The story page's Start agent
button runs this.`,
		Example: `  flai serve agent start S-0115`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Start)
		},
	}
}

func newServeAgentCommitCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "commit <story-id>",
		Short: "Start an agent to commit what a story in review left uncommitted in its worktree",
		Long: `Starts the story's agent, in a new session, to commit what its worktree
under .flai-cache/worktrees holds and nothing else: a story in review whose
worktree has uncommitted changes cannot be accepted, because acceptance
merges the branch as it is committed and removes the worktree. The agent is
told to commit the changes on the story branch the way the project's
conventions say, to discard only what plainly does not belong in the
repository, to log what it did in the narrative, and to leave the story in
review. The run is recorded where the serving flai tracks it, like a start
(S-0140).

It refuses, and says why, while the agent action is off for the project,
when the story is not in review, while its agent runs or waits for an
answer, when the story has no worktree or the worktree has nothing
uncommitted, and when nothing can start it. The Have an agent commit them
button in the dashboard's acceptance confirmation runs this.`,
		Example: `  flai serve agent commit S-0140`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentNow(args[0], serve.Commit)
		},
	}
}

func newServeAgentStopCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <story-id>",
		Short: "Stop a story's agent: end its process and everything it started",
		Long: `Stops the agent flai serve started for a story (S-0170). An agent that
runs is sent SIGTERM with every process it started (its process group), and
SIGKILL if it has not ended ten seconds later. An agent that ended waiting
for an answer is not started again when the answer comes. The run is
recorded as stopped by the operator: the story stays where it is, its
worktree holds what the agent left, committed or not, and it gets no new
agent until it is retried (flai serve agent restart) or moved back to ready.

A process is signalled only while it is still the agent flai started: the
leader of its own session, which on Linux started at the clock tick since
boot recorded with the run. A PID the system has given to another process
since, after a reboot, is left alone, and the run is only recorded as
stopped; flai serve settles such a run by itself at its next look.

It refuses, and says why, when flai serve has started no agent for the story
and when its agent is neither running nor waiting for an answer. It works
whether or not the agent action is on: the dashboard's Stop, on the
activity page, asks for it only while it is.`,
		Example: `  flai serve agent stop S-0170`,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentStop(args[0])
		},
	}
}

// flai serve orchestrate: the orchestrator stopped and held stopped, and
// started again, on the operator's word (S-0228).
func newServeOrchestrateCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "orchestrate",
		Short: "Stop the project's orchestrator and hold it stopped, or start it again",
		Long: `While the orchestrate host action is on, flai serve runs the project's
orchestrator, and starts it again whenever its run ends (S-0218). stop ends
the run and holds the orchestrator stopped, so that flai serve does not start
it again while the action stays on; start lifts the hold and starts it, as
turning the action on does. Turning the action off and on lifts the hold too.
The dashboard's orchestrator page runs these; both are refused while the
action is off.`,
		Example: `  flai serve orchestrate stop
  flai serve orchestrate start`,
		Args: cobra.NoArgs,
	}
	c.AddCommand(
		&cobra.Command{
			Use:   "stop",
			Short: "End the orchestrator's run and hold it stopped while orchestrate stays on",
			Long: `Ends the orchestrator's run as flai serve agent stop ends a story's agent:
its process group is sent SIGTERM, and SIGKILL if it has not ended ten
seconds later, while it is still the orchestrator flai started. The run is
recorded as stopped by the operator, its activity is logged in
wip/agents/orchestrator.md as stopped from the dashboard, and the
orchestrator is held stopped: flai serve does not start it again until flai
serve orchestrate start, or until orchestrate is turned off and on. A run
that failed and waits to be started again is held as it is.

It refuses, and says why, while the orchestrate action is off for the
project, when flai serve has started no orchestrator for it, and when the
orchestrator is already held stopped.`,
			Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.orchestrateNow(false) },
		},
		&cobra.Command{
			Use:   "start",
			Short: "Start the orchestrator now, lifting a hold on it",
			Long: `Starts an orchestrator run at once, for the reason asked, with any reasons
flai serve has queued for it (ADR-0209): in the project's main checkout, with
its orchestration agent and the harnesses and the command set with flai serve
agent. It lifts the hold flai serve orchestrate stop put on the orchestrator,
if there is one. The run acts on its reasons and ends, and flai serve starts
the next when another reason comes. A start that fails lifts the hold as
well, keeps the reasons, and flai serve tries again a minute later.

It refuses, and says why, while the orchestrate action is off for the
project, and while an orchestrator run is going or still being stopped.`,
			Args: cobra.NoArgs,
			RunE: func(*cobra.Command, []string) error { return a.orchestrateNow(true) },
		},
	)
	return c
}

// orchestrateNow stops and holds, or starts, the orchestrator of the project
// in the working directory, and prints the run as recorded.
func (a *app) orchestrateNow(start bool) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	var run *serve.AgentRun
	if start {
		run, err = serve.OrchestratorStart(context.Background(), o, e)
	} else {
		run, err = serve.OrchestratorStop(o, e)
	}
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"held": run.Held, "orchestrator": run})
	}
	if start {
		fmt.Fprintf(a.out, "started %s (%s) as %s (pid %d); log %s\n", run.Command, run.Harness, run.Agent, run.PID, run.Log)
		return nil
	}
	fmt.Fprintf(a.out, "stopped the orchestrator %s (pid %d), and held it stopped; flai serve orchestrate start starts it again\n", run.Agent, run.PID)
	return nil
}

// agentStop stops story's agent and prints the run as it ended.
func (a *app) agentStop(story string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	run, err := serve.Stop(o, e, workitem.CanonicalID(story))
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "pid": run.PID, "stopped": run.Stopped, "ended": run.Ended, "why": run.Why})
	}
	fmt.Fprintf(a.out, "stopped %s's agent %s (pid %d)\n", run.Story, run.Agent, run.PID)
	return nil
}

func newServeAgentStreamCmd(a *app) *cobra.Command {
	var from int64
	var follow bool
	c := &cobra.Command{
		Use:   "stream <story-id>",
		Short: "Print what a story's agent said and did, from the log flai serve gave it",
		Long: `Prints the stream of the newest agent flai serve started for a story:
what it said, the tools it called and what they answered, the background
tasks it ran, and how its session ended, one entry to a line, each cut to a
few hundred characters. It is read from the log flai serve writes the
agent's output to; Claude Code's stream-json is turned into entries, and any
other output is printed as it is. The dashboard's activity page shows the
same stream (S-0142).

Without --from it starts at the last 256 KiB of the log. --follow goes on
printing what the agent writes until it ends. --json prints each read as
one JSON object on a line: the entries, whether the agent runs, and next,
the offset to read from again with --from.`,
		Example: `  flai serve agent stream S-0142
  flai serve agent stream S-0142 --follow
  flai serve agent stream S-0142 --from 0 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentStream(workitem.CanonicalID(args[0]), from, follow)
		},
	}
	c.Flags().Int64Var(&from, "from", -1, "the byte offset of the log to read from (default: its last 256 KiB)")
	c.Flags().BoolVarP(&follow, "follow", "f", false, "go on printing until the agent ends")
	return c
}

// agentStream prints a story's agent's stream, once or until it ends.
func (a *app) agentStream(story string, from int64, follow bool) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	root := mainRootOf(repo)
	started := ""
	for {
		got, err := serve.Stream(a.serveDir().AgentStates()[root], story, from)
		if errors.Is(err, hostapi.ErrNoAgent) {
			return fmt.Errorf("rule: flai serve has started no agent for %s", story)
		}
		if err != nil {
			return err
		}
		if started != "" && got.Started != started {
			// another agent was started for the story: its log is new
			from = 0
			started = got.Started
			continue
		}
		if a.jsonOut {
			line, err := json.Marshal(got)
			if err != nil {
				return err
			}
			fmt.Fprintln(a.out, string(line))
		} else {
			if started == "" {
				state := "ended " + got.Ended
				if got.Running {
					state = "running"
				}
				fmt.Fprintf(a.out, "%s: %s, started %s, %s\n", story, got.Agent, got.Started, state)
			}
			printStream(a, got)
		}
		started, from = got.Started, got.Next
		if !follow || (!got.Running && !got.More) {
			return nil
		}
		if !got.More {
			a.sleep(time.Second)
		}
	}
}

// printStream prints the entries of one read, one to a line, the lines of
// an entry's text after its first indented under it.
func printStream(a *app, got *serve.StreamRead) {
	if got.Skipped > 0 {
		fmt.Fprintf(a.out, "(%d earlier %s left out)\n", got.Skipped, oneOrMany(got.Skipped, "entry", "entries"))
	}
	for _, e := range got.Entries {
		text := e.Text
		if e.Tool != "" {
			text = e.Tool + ": " + text
		}
		kind := e.Kind
		if e.Error {
			kind += "!"
		}
		fmt.Fprintf(a.out, "%-9s %s\n", kind, strings.ReplaceAll(text, "\n", "\n          "))
	}
}

// mcpAgents starts or restarts a story's agent for flai mcp's tools
// agent_start and agent_restart (S-0177, ADR-0064), as flai serve agent
// start and restart do, under the same agent host action, and journals who
// asked and what came of it. They are for the operator's own agent: one that
// flai serve started, whose flai mcp has FLAI_STARTED_BY, is refused, since
// a start goes past holds, the in-progress limit, and a full review on the
// operator's word.
func (a *app) mcpAgents(ctx context.Context, verb, root, story, by string) (mcpserver.AgentStarted, error) {
	if os.Getenv("FLAI_STARTED_BY") == "flai-serve" {
		return mcpserver.AgentStarted{}, fmt.Errorf("an agent flai serve started does not start agents: %s's agent is the operator's to %s, from the story's page or with flai serve agent %s %s on the host", story, verb, verb, story)
	}
	how := serve.Start
	if verb == "restart" {
		how = serve.Restart
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.AgentStarted{}, err
	}
	host := a.host()
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: host}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: root}
	run, err := how(ctx, o, e, story)
	entry := hostapi.Entry{At: a.now().UTC().Format(time.RFC3339), Action: hostapi.ActionAgent, Method: "mcp.agent_" + verb, Project: e.Key, Root: root, By: by, Outcome: "done"}
	var no *serve.Refused
	switch {
	case errors.As(err, &no) && !o.Agent(root).Enabled:
		entry.Outcome, entry.Detail = "disabled", fmt.Sprintf("%s asked to %s %s's agent: %s", by, verb, story, no.Why)
	case err != nil:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s asked to %s %s's agent: %s", by, verb, story, err)
	case run.Queued != "":
		entry.Detail = fmt.Sprintf("%s asked to %s %s's agent: queued, and flai serve starts it when it can", by, verb, story)
	default:
		entry.Detail = fmt.Sprintf("%s asked to %s %s's agent: started %s as %s (pid %d)", by, verb, story, run.Command, run.Agent, run.PID)
	}
	host.Record(entry)
	if err != nil {
		return mcpserver.AgentStarted{}, err
	}
	return mcpserver.AgentStarted{Story: run.Story, Agent: run.Agent, Harness: run.Harness, Command: run.Command, PID: run.PID, Log: run.Log, Started: run.Started, Queued: run.Queued}, nil
}

// agentNow starts story's agent on the operator's word, with start or
// restart deciding whether it may, and prints the run, or that it was
// queued.
func (a *app) agentNow(story string, how func(context.Context, serve.Options, serve.Entry, string) (*serve.AgentRun, error)) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	o := serve.Options{Dir: a.serveDir(), Logger: a.logger(), Now: a.now, Agent: a.agentConfig, Host: a.host()}
	e := serve.Entry{Key: repo.Manifest.Key, Name: repo.Manifest.Name, Root: mainRootOf(repo)}
	run, err := how(context.Background(), o, e, workitem.CanonicalID(story))
	var no *serve.Refused
	if errors.As(err, &no) {
		return fmt.Errorf("rule: %s", no.Why)
	}
	if err != nil {
		return err
	}
	if run.Queued != "" {
		if a.jsonOut {
			return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "queued": run.Queued})
		}
		fmt.Fprintf(a.out, "queued another agent for %s: the in-progress limit, a full review, or a hold leaves no room, and flai serve starts it when there is\n", run.Story)
		return nil
	}
	if a.jsonOut {
		return a.printJSON(map[string]any{"story": run.Story, "agent": run.Agent, "harness": run.Harness, "command": run.Command, "pid": run.PID, "log": run.Log, "started": run.Started})
	}
	fmt.Fprintf(a.out, "started %s (%s) for %s as %s (pid %d); log %s\n", run.Command, run.Harness, run.Story, run.Agent, run.PID, run.Log)
	return nil
}

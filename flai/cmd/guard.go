package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// exitGuardRefused is the exit status with which a Claude Code hook refuses
// a tool call; the hook's standard error goes back to the caller.
const exitGuardRefused = 2

func newGuardCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "guard",
		Short: "Refuse a sub-agent's writes, hold the planner to planning and the orchestrator to its permissions, and keep a story's agent from waiting on its sub-agents with wait_for_events, as a Claude Code hook",
		Long: `Reads a Claude Code PreToolUse hook's input on standard input and refuses
the call when a sub-agent makes it (the input carries an agent_id) and it
would change a work item, a
thread, a narrative, or the repository's history (ADR-0059, ADR-0060): any
of flai's MCP tools but board, doc_get, doc_search, item_get,
order_by_policy, prime, promote_candidates, release_evaluate, shared_paths,
thread_get, and who_touches; a flai command other than one that reads
(board, check, cod, doc search and show, forecast, help, issue list, order
--by without --apply and without a story to place, plan --candidates,
prime, promote --candidates or --drafts, release --evaluate, shared list
and check, show, stats, stream diff,
thread list and show, touches suggest, version, or any with --help); and a
git
command other than one that reads (blame, cat-file, describe, diff, grep, log, ls-files,
ls-tree, merge-base, rev-list, rev-parse, shortlog, show, status). A
refusal prints why on standard error and exits 2, which Claude Code hands
back to the sub-agent. Every word of a command line is looked at, so a
command run through env, sudo, timeout, xargs, find -exec, or a shell's -c
is found too. The story's agent's own calls carry no agent_id and pass, save
the wait below, as does anything it cannot read: the guard fails open. An
input whose hook_event_name is neither SubagentStart nor SubagentStop is a
PreToolUse's, named or not. It is not a shell, and
a command hidden on purpose (a backslash in its name, a variable holding
it) gets past it.

No sub-agent, and no session flai serve starts (FLAI_STARTED_BY=flai-serve,
FLAI_ROLE, or FLAI_STORY set), changes the manifest's shared paths: the MCP
tool shared_paths_edit and flai shared add and remove are refused, because
claims.shared decides which overlaps hold a story (ADR-0096). Only the
operator's own session changes them.

In a planner session, one flai serve starts with FLAI_ROLE=plan, the
session's own calls are held to planning too (strategic-agents.md): besides
what a sub-agent may do, the MCP tools inbox, item_new, item_edit,
thread_open, thread_reply, activity_log, wait_for_events, and item_move to
backlog; the commands story new, epic new, task new, edit (but not
--no-draft), touches, thread new and reply, issue new and bump, and move to
backlog. A story the planner creates is a draft for the operator to
finalize: item_new of a story needs draft true, and story new needs
--draft. It refuses the planner every other flai tool and command, git's
writes, and the Edit, Write, and NotebookEdit tools. The planner's
sub-agents are held as any sub-agent is.

In an orchestrator session, one flai serve starts with
FLAI_ROLE=orchestrate, the session's own calls are held to
orchestration.permissions in the system-flow.yaml of the project the hook
runs in (S-0218), each off when unset or when the manifest is unreadable.
Whatever its permissions it may do what a sub-agent may, call inbox,
activity_log, wait_for_events, and thread_open, and run thread new and
issue new and bump. Each permission allows more: plan_backlog_epics the
MCP tool plan and flai plan on an epic; finalize_drafts item_edit with
draft false and flai edit --no-draft, each with nothing else to change;
promote_to_ready item_move and flai move of a story to ready; order_ready
flai order --by <policy> --apply; answer_threads a reply on a thread
another opened, with thread_reply or flai thread reply, as a recommendation
(recommendation true, --recommend) while it is recommend, and also as an
answer citing a source (source, --source) while it is autonomous, so that an
answer it cannot source goes to the operator as a recommendation (S-0220);
accept_reviews flai accept and flai move <story> done, each only with --by
orchestrator, so that it accepts as itself and as nobody else (S-0221,
ADR-0093); publish the MCP tool release_publish alone, while flai
release other than --evaluate, flai push, git push, and git tag are refused
whatever its permissions (S-0222, ADR-0094).
On a thread it opened it follows up and resolves, but never recommends or
answers; it never resolves another's thread, never confirms a
recommendation, and never names another author with --by. A call a permission
would allow is refused while it is off, naming it (it needs
orchestration.permissions.<name>); anything else that writes is refused
as what the orchestrator never does: plan for a story, flai order placing
a story by hand, item_move to done, whose refusal names flai accept --by
orchestrator, other flai tools and commands, git's writes, and the
Edit, Write, and NotebookEdit tools. Whether a story it promotes is held,
a draft, or over the ready limit is flai's to check as it moves it, not
the guard's: the guard reads the call, not the board. Each
refusal is logged under ## Refusals in wip/agents/orchestrator.md, with
its time, the call, and the permission it needs; a refusal that cannot be
logged is warned of and refused all the same. The orchestrator's
sub-agents are held as any sub-agent is.

Given a SubagentStart or SubagentStop hook's input, as hook_event_name
says, it records the sub-agent (its agent_id and agent_type, and when it
started) as running in the hook's session, or no longer, in
.flai-cache/guard/<session_id>.json in the project's main checkout, under a
lock, since a layer's sub-agents start at once. It refuses neither and says
nothing; a record it cannot write is warned of. In a story's agent's
session, one flai serve starts with FLAI_STORY and no FLAI_ROLE, it refuses
the agent's own wait_for_events while the session's record lists a
sub-agent running and no unresolved thread is on the story or one of its
tasks (S-0285): wait_for_events reports work items and threads, not
sub-agents, so it would run to its timeout. The refusal names the running
sub-agents and says to wait for one by launching it with the Agent tool's
run_in_background set to false, which returns its result as the tool's
result. With a thread on the story open, the agent waits on the designer,
and the call passes and returns on the answer. A project, record, or
threads it cannot read let the call through.

The template's .claude/settings.json runs it before Bash and flai's MCP
tools, on SubagentStart and SubagentStop, and, in a planner's or an
orchestrator's session alone, before Edit, Write, and NotebookEdit as well.`,
		Example: `  flai guard < hook-input.json`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				a.logger().Warn("hook input unreadable", "component", "guard", "err", err.Error())
				return nil
			}
			var e guard.Event
			if err := json.Unmarshal(data, &e); err != nil {
				a.logger().Warn("hook input is not a tool call", "component", "guard", "err", err.Error())
				return nil
			}
			if e.HookEventName == guard.HookSubagentStart || e.HookEventName == guard.HookSubagentStop {
				a.recordSubAgent(e)
				return nil
			}
			g := guard.Guard{Role: os.Getenv("FLAI_ROLE"), Story: os.Getenv("FLAI_STORY"), Served: os.Getenv(guard.StartedByEnv) == guard.StartedByServe}
			for _, c := range cmd.Root().Commands() {
				g.Commands = append(g.Commands, c.Name())
			}
			var repo *workitem.Repo
			if g.Role == guard.RoleOrchestrate && e.AgentID == "" {
				if repo, err = a.guardProject(); err != nil {
					a.logger().Warn("project unreadable, the orchestrator's permissions are taken as off", "component", "guard", "err", err.Error())
				} else {
					g.Permissions = repo.Manifest.Orchestration.Permissions
					g.Opener = guard.ThreadOpener(repo)
				}
			}
			if g.Role == "" && g.Story != "" && e.AgentID == "" {
				a.holdWaits(&g, e.SessionID)
			}
			r := g.Decide(e)
			if r.Why == "" {
				return nil
			}
			fmt.Fprintln(a.errOut, r.Why)
			if repo != nil && r.Call != "" {
				if _, err := repo.AppendRefusal(workitem.ActivityOrchestrator, workitem.ActivityRefusal{At: a.now(), Call: r.Call, Needs: r.Needs}); err != nil {
					a.logger().Warn("refusal not logged in the orchestrator's activity document", "component", "guard", "err", err.Error())
				}
			}
			return &exitError{code: exitGuardRefused}
		},
	}
}

// guardProject is the project the hook runs in, from the working directory,
// for its manifest and the orchestrator's activity document.
func (a *app) guardProject() (*workitem.Repo, error) {
	start := a.cwd
	if start == "" {
		var err error
		if start, err = os.Getwd(); err != nil {
			return nil, err
		}
	}
	return workitem.Open(start)
}

// recordSubAgent notes in the project's .flai-cache/guard that e's
// sub-agent started or stopped in its session. It says nothing when it
// does, and warns when it cannot: the guard fails open.
func (a *app) recordSubAgent(e guard.Event) {
	repo, err := a.guardProject()
	if err == nil {
		err = guard.Record(filepath.Join(repo.CacheDir(), guard.RecordDir), e, a.now())
	}
	if err != nil {
		a.logger().Warn("sub-agent not recorded", "component", "guard", "event", e.HookEventName, "agent_id", e.AgentID, "err", err.Error())
	}
}

// holdWaits gives g what it needs to hold a story's agent's wait_for_events:
// the sub-agents session's record lists as running, and whether a thread on
// the story is open, each read from the project only when the call is
// decided on. A project, a record, or threads it cannot read are warned of
// and let the call through.
func (a *app) holdWaits(g *guard.Guard, session string) {
	project := sync.OnceValues(a.guardProject)
	g.Running = func() ([]guard.SubAgent, error) {
		repo, err := project()
		if err != nil {
			a.logger().Warn("project unreadable, the story's agent's wait is let through", "component", "guard", "err", err.Error())
			return nil, err
		}
		running, err := guard.Running(filepath.Join(repo.CacheDir(), guard.RecordDir), session)
		if err != nil {
			a.logger().Warn("sub-agents unreadable, the story's agent's wait is let through", "component", "guard", "err", err.Error())
		}
		return running, err
	}
	g.ThreadOpen = func() (bool, error) {
		repo, err := project()
		if err == nil {
			var open bool
			if open, err = guard.StoryThreadOpen(repo, g.Story)(); err == nil {
				return open, nil
			}
		}
		a.logger().Warn("threads unreadable, the story's agent's wait is let through", "component", "guard", "story", g.Story, "err", err.Error())
		return false, err
	}
}

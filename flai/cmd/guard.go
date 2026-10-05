package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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
		Short: "Refuse a sub-agent's writes, hold the planner to planning and the orchestrator to its permissions, as a Claude Code PreToolUse hook",
		Long: `Reads a Claude Code PreToolUse hook's input on standard input and refuses
the call when a sub-agent makes it (the input carries an agent_id) and it
would change a work item, a
thread, a narrative, or the repository's history (ADR-0059, ADR-0060): any
of flai's MCP tools but board, doc_get, doc_search, item_get,
order_by_policy, prime, promote_candidates, release_evaluate, thread_get,
and who_touches; a flai command other than one that reads (board, check,
cod, doc search and show, forecast, help, issue list, order --by without
--apply, prime, promote --candidates, release --evaluate, show, stats,
stream diff, thread list and show, touches suggest, version, or any with
--help); and a git
command other than one that reads (blame, cat-file, describe, diff, grep, log, ls-files,
ls-tree, merge-base, rev-list, rev-parse, shortlog, show, status). A
refusal prints why on standard error and exits 2, which Claude Code hands
back to the sub-agent. Every word of a command line is looked at, so a
command run through env, sudo, timeout, xargs, find -exec, or a shell's -c
is found too. The story's agent's own calls carry no agent_id and pass, as
does anything it cannot read: the guard fails open. It is not a shell, and
a command hidden on purpose (a backslash in its name, a variable holding
it) gets past it.

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
MCP tool plan and flai plan on an epic; finalize_drafts flai edit
--no-draft with nothing else to change; promote_to_ready item_move and
flai move to ready; order_ready flai order that writes; answer_threads,
recommend or autonomous, thread_reply and flai thread reply;
accept_reviews flai accept; publish flai release --pending and flai push.
A call a permission would allow is refused while it is off, naming it
(it needs orchestration.permissions.<name>); anything else that writes is
refused as what the orchestrator never does: other flai tools and
commands, git's writes, and the Edit, Write, and NotebookEdit tools. Each
refusal is logged under ## Refusals in wip/agents/orchestrator.md, with
its time, the call, and the permission it needs; a refusal that cannot be
logged is warned of and refused all the same. The orchestrator's
sub-agents are held as any sub-agent is.

The template's .claude/settings.json runs it before Bash and flai's MCP
tools, and, in a planner's or an orchestrator's session alone, before Edit,
Write, and NotebookEdit as well.`,
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
			g := guard.Guard{Role: os.Getenv("FLAI_ROLE")}
			for _, c := range cmd.Root().Commands() {
				g.Commands = append(g.Commands, c.Name())
			}
			var repo *workitem.Repo
			if g.Role == guard.RoleOrchestrate && e.AgentID == "" {
				if repo, err = a.guardProject(); err != nil {
					a.logger().Warn("project unreadable, the orchestrator's permissions are taken as off", "component", "guard", "err", err.Error())
				} else {
					g.Permissions = repo.Manifest.Orchestration.Permissions
				}
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

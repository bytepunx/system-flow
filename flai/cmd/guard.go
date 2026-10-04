package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// exitGuardRefused is the exit status with which a Claude Code hook refuses
// a tool call; the hook's standard error goes back to the caller.
const exitGuardRefused = 2

func newGuardCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "guard",
		Short: "Refuse a sub-agent's writes and hold the planner to planning, as a Claude Code PreToolUse hook",
		Long: `Reads a Claude Code PreToolUse hook's input on standard input and refuses
the call when a sub-agent makes it (the input carries an agent_id) and it
would change a work item, a
thread, a narrative, or the repository's history (ADR-0059, ADR-0060): any
of flai's MCP tools but board, doc_get, doc_search, item_get, prime,
thread_get, and who_touches; a flai command other than one that reads
(board, check, cod, doc search and show, forecast, help, issue list, prime,
show, stats, stream diff, thread list and show, touches suggest, version, or
any with --help); and a git
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

The template's .claude/settings.json runs it before Bash and flai's MCP
tools, and, in a planner session alone, before Edit, Write, and NotebookEdit
as well.`,
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
			why := g.Check(e)
			if why == "" {
				return nil
			}
			fmt.Fprintln(a.errOut, why)
			return &exitError{code: exitGuardRefused}
		},
	}
}

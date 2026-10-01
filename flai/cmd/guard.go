package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// exitGuardRefused is the exit status with which a Claude Code hook refuses
// a tool call; the hook's standard error goes back to the caller.
const exitGuardRefused = 2

func newGuardCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "guard",
		Short: "Refuse a sub-agent's writes, as a Claude Code PreToolUse hook",
		Long: `Reads a Claude Code PreToolUse hook's input on standard input and refuses
the call when a sub-agent makes it (the input carries an agent_id) and it
would change a work item, a
thread, a narrative, or the repository's history (ADR-0059, ADR-0060): any
of flai's MCP tools but board, doc_get, doc_search, item_get, prime,
thread_get, and who_touches; a flai command other than one that reads
(board, check, doc search and show, help, issue list, prime, show, stats,
stream diff, thread list and show, version, or any with --help); and a git
command other than one that reads (blame, cat-file, describe, diff, grep, log, ls-files,
ls-tree, merge-base, rev-list, rev-parse, shortlog, show, status). A
refusal prints why on standard error and exits 2, which Claude Code hands
back to the sub-agent. Every word of a command line is looked at, so a
command run through env, sudo, timeout, xargs, find -exec, or a shell's -c
is found too. The story's agent's own calls carry no agent_id and pass, as
does anything it cannot read: the guard fails open. It is not a shell, and
a command hidden on purpose (a backslash in its name, a variable holding
it) gets past it.

The template's .claude/settings.json runs it before Bash and flai's MCP
tools.`,
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
			g := guard.Guard{}
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

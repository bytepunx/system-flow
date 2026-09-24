package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newMoveCmd(a *app) *cobra.Command {
	var reason, by string
	var dryRun bool
	var accept acceptOptions
	c := &cobra.Command{
		Use:   "move <id> <state>",
		Short: "Transition a work item, enforcing the workflow rules",
		Long: `Move an item to a new state. States: ` + strings.Join(workitem.States, ", ") + `.

Rules from design/system/workflow.md are enforced: a story needs acceptance
criteria before ready, and at least one task and no open question in its
narrative before review; children must be closed before done; cancelling or
sending review back needs --reason. WIP limit breaches warn.

Cancelling an epic cancels every open story under it and their open tasks;
cancelling a story cancels its open tasks (ADR-0028). The items are listed
first, a terminal is asked unless --yes is given, and --dry-run changes
nothing. Branches, worktrees, and narratives are left as they are.

Moving a story from review to done is acceptance: it runs the same flow as
flai accept (merge the story branch, archive, commit), with the same flags.
There is no other way for a story to become done. Acceptance computes no
release; see flai release --pending.`,
		Example: `  flai move S-004 in-progress
  flai move T-021 done
  flai move S-004 in-progress --reason "tests missing"   # from review
  flai move S-009 cancelled --reason "superseded by S-012"
  flai move E-003 cancelled --reason "a different route" --dry-run`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if it.Type == workitem.Story && args[1] == workitem.Done && it.Status == workitem.Review {
				accept.by = by
				accept.dryRun = accept.dryRun || dryRun
				res, err := a.acceptItem(repo, it, accept)
				if err != nil {
					return err
				}
				return a.printAccept(res)
			}
			if args[1] == workitem.Cancelled {
				return a.cancelItem(repo, it, a.movedBy(by), reason, dryRun)
			}
			if dryRun {
				return fmt.Errorf("--dry-run previews a cancellation or an acceptance; other moves have nothing to preview")
			}
			warnings, err := repo.Transition(it, args[1], a.movedBy(by), reason, a.now())
			if err != nil {
				return err
			}
			for _, w := range warnings {
				a.logger().Warn("workflow policy warning", "component", "workitem", "item", it.ID, "detail", w)
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "status": it.Status, "warnings": warnings})
			}
			fmt.Fprintf(a.out, "%s → %s\n", it.ID, it.Status)
			if it.Status == workitem.Done && it.Type == workitem.Epic {
				a.logger().Info("epic closed without an archive or a commit; flai accept does move, archive, and commit in one step", "component", "workitem", "item", it.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&reason, "reason", "", "why (required for cancelled and review → in-progress)")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "for a move to cancelled: list what would be cancelled with it and change nothing")
	c.Flags().StringVar(&by, "by", "", "who made the change (default: FLAI_AGENT when set, else the config author)")
	addAcceptFlags(c, &accept)
	return c
}

// movedBy is who a transition is recorded as made by: --by, else FLAI_AGENT
// when an agent's session sets it, else the config author. An agent's moves
// must carry its name so that they are not reported back to it as the
// designer's (S-0058). Acceptance keeps the config author: it is the
// operator's act whoever types the command.
func (a *app) movedBy(by string) string {
	if by != "" {
		return by
	}
	if agent := os.Getenv("FLAI_AGENT"); agent != "" {
		return agent
	}
	return a.author()
}

func newBlockCmd(a *app) *cobra.Command {
	var reason string
	c := &cobra.Command{
		Use:   "block <id>",
		Short: "Open a blocked interval on an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if err := workitem.BlockItem(it, reason, a.now()); err != nil {
				return err
			}
			if err := repo.Save(it); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "blocked": true, "reason": reason})
			}
			fmt.Fprintf(a.out, "%s blocked: %s\n", it.ID, reason)
			return nil
		},
	}
	c.Flags().StringVar(&reason, "reason", "", "why the item is blocked")
	_ = c.MarkFlagRequired("reason")
	return c
}

func newUnblockCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "unblock <id>",
		Short: "Close the open blocked interval on an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			if err := workitem.UnblockItem(it, a.now()); err != nil {
				return err
			}
			if err := repo.Save(it); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": it.ID, "blocked": false})
			}
			fmt.Fprintf(a.out, "%s unblocked\n", it.ID)
			return nil
		},
	}
}

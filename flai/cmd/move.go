package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newMoveCmd(a *app) *cobra.Command {
	var reason, by string
	var accept acceptOptions
	c := &cobra.Command{
		Use:   "move <id> <state>",
		Short: "Transition a work item, enforcing the workflow rules",
		Long: `Move an item to a new state. States: ` + strings.Join(workitem.States, ", ") + `.

Rules from design/system/workflow.md are enforced: a story needs tasks and
acceptance criteria before ready, children must be closed before done, and
cancelling or sending review back needs --reason. WIP limit breaches warn.

Moving a story from review to done is acceptance: it runs the same flow as
flai accept (merge the story branch, archive, commit, tag, push), with the
same flags. There is no other way for a story to become done.`,
		Example: `  flai move S-004 in-progress
  flai move T-021 done
  flai move S-004 in-progress --reason "tests missing"   # from review
  flai move S-009 cancelled --reason "superseded by S-012"`,
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
				res, err := a.acceptItem(repo, it, accept)
				if err != nil {
					return err
				}
				return a.printAccept(res)
			}
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			warnings, err := repo.Move(it, args[1], workitem.MoveOptions{
				By: orDefault(by, a.author()), Reason: reason, Now: a.now(), Items: items, Board: board,
			})
			if err != nil {
				return err
			}
			if err := repo.Save(it); err != nil {
				return err
			}
			if it.Type == workitem.Story {
				if err := board.Save(a.now().Format("2006-01-02")); err != nil {
					return err
				}
			}
			if err := a.refreshIndex(repo); err != nil {
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
				a.logger().Info("epic closed without a release; flai accept does move, archive, commit, release, and push in one step", "component", "workitem", "item", it.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&reason, "reason", "", "why (required for cancelled and review → in-progress)")
	c.Flags().StringVar(&by, "by", "", "who made the change (default: config author)")
	addAcceptFlags(c, &accept)
	return c
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

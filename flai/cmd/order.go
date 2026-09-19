package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newOrderCmd(a *app) *cobra.Command {
	var p workitem.Placement
	c := &cobra.Command{
		Use:   "order <story>",
		Short: "Place a ready or backlog story in the pull order",
		Long: `Put a story at a position in its column's pull order, the order list in
wip/kanban/board.md that agents pull from.

Only the ready and backlog columns have an order, and only stories are in it:
ready stories come first, then backlog stories in the order they should be
refined. A story the list does not name comes after the ones it does, by ID.
The position is relative to another story of the same column, or the top or
bottom of that column. To change a story's column use flai move.`,
		Example: `  flai order S-0061 --top
  flai order S-0059 --before S-0061
  flai order S-0047 --after S-0053
  flai order S-0056 --bottom`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			id := args[0]
			if err := board.Place(items, id, p); err != nil {
				return err
			}
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return err
			}
			it, err := repo.Get(id)
			if err != nil {
				return err
			}
			sequence := workitem.PullSequence(board.Order, items, it.Status)
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": id, "status": it.Status, "sequence": sequence, "order": board.Order})
			}
			fmt.Fprintf(a.out, "%s: %s\n", it.Status, strings.Join(sequence, ", "))
			return nil
		},
	}
	c.Flags().StringVar(&p.Before, "before", "", "place it just before this story of the same column")
	c.Flags().StringVar(&p.After, "after", "", "place it just after this story of the same column")
	c.Flags().BoolVar(&p.Top, "top", false, "place it first in its column")
	c.Flags().BoolVar(&p.Bottom, "bottom", false, "place it last in its column")
	return c
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newOrderCmd(a *app) *cobra.Command {
	var p workitem.Placement
	var by, placedBy string
	var apply bool
	var keep time.Duration
	c := &cobra.Command{
		Use:   "order <story>",
		Short: "Place a ready or backlog story in the pull order",
		Long: `Put a story at a position in its column's pull order, the order list in
wip/kanban/board.md that agents pull from.

Only the ready and backlog columns have an order, and only stories are in it:
ready stories come first, then backlog stories in the order they should be
refined. A story the list does not name comes after the ones it does, by ID.
The position is relative to another story of the same column, or the top or
bottom of that column. To change a story's column use flai move.

Each placement is recorded in board.md's placed map with who made it and
when: --placed-by, else FLAI_AGENT, else the config author, and always
orchestrator when flai serve runs the orchestrator (FLAI_ROLE=orchestrate).
The dashboard's drag runs this command, so it is recorded too. flai move
drops a story's record when it takes the story out of its column.

With --by, it computes the ready column's order by a policy instead and
prints it with the figure each story was ordered by: cod, cost of delay
value, highest first; wsjf, that value over the forecast duration in hours,
highest first; throughput, forecast duration, shortest first; fifo, created,
oldest first. A story without the figure goes after those with it, in its
current order, and ties keep the current order. A ready story placed by
anyone but the orchestrator within --keep-placed, a day unless given, keeps
its place, and the rest are ordered around it; the order marks each kept
story with who placed it and when. --keep-placed 0 keeps none. It writes
nothing unless --apply is given, which writes the computed order to
board.md.`,
		Example: `  flai order S-0061 --top
  flai order S-0059 --before S-0061
  flai order S-0047 --after S-0053
  flai order S-0056 --bottom --placed-by alex
  flai order --by wsjf
  flai order --by cod --apply
  flai order --by cod --apply --keep-placed 2h`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if by != "" {
				if p != (workitem.Placement{}) {
					return fmt.Errorf("--by orders the whole ready column, and --before, --after, --top, and --bottom place one story; give one or the other")
				}
				if len(args) > 0 {
					return fmt.Errorf("--by orders the whole ready column, so it takes no story; drop %s", args[0])
				}
				if placedBy != "" {
					return fmt.Errorf("--placed-by names who places one story, and --by orders the whole ready column; drop --placed-by")
				}
				if keep < 0 {
					return fmt.Errorf("--keep-placed is how long a story placed by hand keeps its place, 0 for none; got %s", keep)
				}
				return orderBy(a, by, apply, keep)
			}
			if apply {
				return fmt.Errorf("--apply writes the order --by computes; give --by too")
			}
			if cmd.Flags().Changed("keep-placed") {
				return fmt.Errorf("--keep-placed keeps placements when --by orders the ready column; give --by too")
			}
			if len(args) != 1 {
				return fmt.Errorf("name the story to place, or give --by to order the ready column")
			}
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
			board.RecordPlacement(id, a.placedBy(placedBy), a.now())
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return err
			}
			it, err := repo.Get(id)
			if err != nil {
				return err
			}
			sequence := workitem.PullSequence(board.Order, items, it.Status)
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": id, "status": it.Status, "sequence": sequence, "order": board.Order, "placed": board.Placed[id]})
			}
			fmt.Fprintf(a.out, "%s: %s\n", it.Status, strings.Join(sequence, ", "))
			return nil
		},
	}
	c.Flags().StringVar(&p.Before, "before", "", "place it just before this story of the same column")
	c.Flags().StringVar(&p.After, "after", "", "place it just after this story of the same column")
	c.Flags().BoolVar(&p.Top, "top", false, "place it first in its column")
	c.Flags().BoolVar(&p.Bottom, "bottom", false, "place it last in its column")
	c.Flags().StringVar(&placedBy, "placed-by", "", "who placed it (default: FLAI_AGENT, then the config author)")
	c.Flags().StringVar(&by, "by", "", "compute the ready column's order by a policy: "+strings.Join(workitem.OrderPolicies, ", "))
	c.Flags().BoolVar(&apply, "apply", false, "write the order --by computes to board.md")
	c.Flags().DurationVar(&keep, "keep-placed", workitem.DefaultKeepPlaced, "with --by, a ready story placed by hand this recently keeps its place; 0 for none")
	return c
}

// placedBy is who a placement is recorded as made by: the orchestrator when
// flai serve runs it, whatever FLAI_AGENT a harness's settings give its
// shell and whatever it says with --placed-by, so that its placements never
// pass for one by hand; else --placed-by, FLAI_AGENT, or the config author,
// as movedBy.
func (a *app) placedBy(by string) string {
	if os.Getenv("FLAI_ROLE") == conventions.RoleOrchestrate {
		return workitem.ActivityOrchestrator
	}
	return a.movedBy(by)
}

// orderBy computes the ready column's order by a policy, keeping the places
// of stories placed by hand within keep, prints it, and with apply writes it
// to board.md as placing each story would.
func orderBy(a *app, policy string, apply bool, keep time.Duration) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	got, board, items, err := repo.ReadyOrderKeeping(policy, keep, a.now())
	if err != nil {
		return err
	}
	if apply {
		if err := got.Apply(board, items); err != nil {
			return err
		}
		if err := board.Save(a.now().Format("2006-01-02")); err != nil {
			return err
		}
	}
	if a.jsonOut {
		return a.printJSON(got)
	}
	ranked := got.Stories
	if len(ranked) == 0 {
		fmt.Fprintln(a.out, "no ready stories to order")
		return nil
	}
	fmt.Fprintf(a.out, "ready by %s:\n", policy)
	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
	for _, r := range ranked {
		figure := r.Text
		if r.Missing != "" {
			figure = "no " + r.Missing
		}
		title := r.Title
		if r.Kept != nil {
			title += fmt.Sprintf(" (kept: placed by %s at %s)", r.Kept.By, r.Kept.At)
		}
		fmt.Fprintf(w, "  %d\t%s\t%s\t%s\n", r.Position, r.ID, figure, title)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if apply {
		fmt.Fprintln(a.out, "board.md's ready order is now this one")
	}
	return nil
}

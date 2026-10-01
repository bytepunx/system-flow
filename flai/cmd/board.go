package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/pending"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newBoardCmd(a *app) *cobra.Command {
	var all bool
	c := &cobra.Command{
		Use:   "board",
		Short: "Print the kanban board",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			view := workitem.NewBoardView(items, board, a.now(), all, release.PendingIDs(a.runner, repo.Root, repo.Manifest, repo), repo.Manifest.Projects)
			if u := pending.Detect(a.runner, mainRootOf(repo)); u.Pending() {
				view.Unpushed = u
			}
			columns, counts, breaches := view.Columns, view.Counts, view.Breaches
			if a.jsonOut {
				return a.printJSON(view)
			}
			for _, st := range workitem.States {
				cards := columns[st]
				limit := ""
				if l, ok := board.WIPLimits[st]; ok && l > 0 {
					limit = fmt.Sprintf(" (%d/%d)", counts[st], l)
				}
				fmt.Fprintf(a.out, "%s%s\n", st, limit)
				if len(cards) == 0 {
					fmt.Fprintln(a.out, "  -")
					continue
				}
				for _, cd := range cards {
					flag := ""
					if cd.Blocked {
						flag = " BLOCKED"
					}
					if cd.Held != nil {
						flag += " HELD"
					}
					fmt.Fprintf(a.out, "  %-6s %-46s %-12s %6s%s\n", cd.ID, truncate(cd.Title, 46), cd.Nature, cd.Age, flag)
					if len(cd.Touches) > 0 {
						fmt.Fprintf(a.out, "         touches %s\n", strings.Join(cd.Touches, ", "))
					}
					if cd.Held != nil {
						fmt.Fprintf(a.out, "         %s\n", cd.Held.Reason)
					}
					// an open story's plan; a closed one's is history (S-0176)
					if t := cd.Tasks; t != nil && st != workitem.Done && st != workitem.Cancelled {
						fmt.Fprintf(a.out, "         tasks %d ready, %d waiting, %d in progress, %d done; %s\n", t.Ready, t.Waiting, t.InProgress, t.Done, plural(t.Layers, "layer"))
					}
				}
			}
			if len(board.Order) > 0 {
				fmt.Fprintf(a.out, "\npull order: %v\n", board.Order)
			}
			if u := view.Unpushed; u != nil {
				fmt.Fprintf(a.out, "\naccepted, not pushed: %s (%d commit(s) ahead of %s%s); run: %s\n", strings.Join(u.Acceptances, ", "), u.Commits, u.Upstream, tagsNote(u.Tags), u.Command)
			}
			for _, b := range breaches {
				a.logger().Warn("wip limit exceeded", "component", "workitem", "detail", b)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&all, "all", false, "include epics and tasks")
	c.AddCommand(newBoardLimitCmd(a))
	return c
}

func newBoardLimitCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "limit <column> <n>",
		Short: "Set a column's WIP limit",
		Long: `Set the WIP limit of ready, in-progress, or review in wip/kanban/board.md,
the one place flai, flai serve, flai check, and the dashboard read it from.
0 removes the limit. Limits count stories; a move past one warns.`,
		Example: `  flai board limit in-progress 3
  flai board limit review 0`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("a WIP limit is a whole number, 0 for none; got %q", args[1])
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			if err := board.SetLimit(args[0], n); err != nil {
				return err
			}
			if err := board.Save(a.now().Format("2006-01-02")); err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"column": args[0], "limit": n, "wip_limits": board.WIPLimits})
			}
			if n == 0 {
				fmt.Fprintf(a.out, "%s: no WIP limit\n", args[0])
			} else {
				fmt.Fprintf(a.out, "%s: WIP limit %d\n", args[0], n)
			}
			return nil
		},
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

// mainRootOf is the main checkout, where acceptances are committed.
func mainRootOf(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

func tagsNote(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return ", tags " + strings.Join(tags, ", ")
}

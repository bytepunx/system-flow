package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/pending"
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
			items, err := repo.List(false)
			if err != nil {
				return err
			}
			board, err := repo.LoadBoard()
			if err != nil {
				return err
			}
			view := workitem.NewBoardView(items, board, a.now(), all)
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
					fmt.Fprintf(a.out, "  %-6s %-46s %-12s %6s%s\n", cd.ID, truncate(cd.Title, 46), cd.Nature, cd.Age, flag)
					if len(cd.Touches) > 0 {
						fmt.Fprintf(a.out, "         touches %s\n", strings.Join(cd.Touches, ", "))
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
	return c
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

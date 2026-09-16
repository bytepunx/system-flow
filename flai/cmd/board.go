package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

type boardCard struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Title   string `json:"title"`
	Nature  string `json:"nature"`
	Parent  string `json:"parent,omitempty"`
	Blocked bool   `json:"blocked"`
	Age     string `json:"age_in_column"`
	AgeSecs int64  `json:"age_in_column_seconds"`
}

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
			now := a.now()
			columns := map[string][]boardCard{}
			counts := map[string]int{}
			for _, it := range items {
				if !all && it.Type != workitem.Story {
					continue
				}
				age := now.Sub(it.EnteredAt())
				columns[it.Status] = append(columns[it.Status], boardCard{
					ID: it.ID, Type: it.Type, Title: it.Title, Nature: it.Nature, Parent: it.Parent,
					Blocked: it.IsBlocked(), Age: humanDuration(age), AgeSecs: int64(age.Seconds()),
				})
				if it.Type == workitem.Story {
					counts[it.Status]++
				}
			}
			var breaches []string
			for _, st := range workitem.States {
				if limit, ok := board.WIPLimits[st]; ok && limit > 0 && counts[st] > limit {
					breaches = append(breaches, fmt.Sprintf("%s has %d stories, limit %d", st, counts[st], limit))
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"columns": columns, "wip_limits": board.WIPLimits, "order": board.Order, "breaches": breaches})
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
				}
			}
			if len(board.Order) > 0 {
				fmt.Fprintf(a.out, "\npull order: %v\n", board.Order)
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

func humanDuration(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

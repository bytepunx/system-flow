package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newStatsCmd(a *app) *cobra.Command {
	var since, typ, by string
	c := &cobra.Command{
		Use:   "stats",
		Short: "Print flow metrics: throughput, cycle time, WIP, flow efficiency, time in state",
		Long: `The reference implementation of design/system/metrics.md. Aggregates cover
items completed in the window (default 30d); WIP and aging are as of now.
--json includes per-item values, weekly throughput, burn-up and cumulative
flow series, and aging items, for dashboards and scripts.`,
		Example: `  flai stats
  flai stats --since 90d --by nature
  flai stats --type task --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			window, err := parseWindow(since)
			if err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			rep := metrics.Compute(items, metrics.Options{Now: a.now(), Since: window, Type: typ, By: by})
			if a.jsonOut {
				return a.printJSON(rep)
			}
			printSummary(a, rep)
			return nil
		},
	}
	c.Flags().StringVar(&since, "since", "30d", "window for completed items: 7d, 90d, 12w, 720h")
	c.Flags().StringVar(&typ, "type", workitem.Story, "item type: epic, story, task")
	c.Flags().StringVar(&by, "by", "", "group by nature, type, or parent")
	return c
}

func parseWindow(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "d") || strings.HasSuffix(s, "w") {
		var n int
		unit := s[len(s)-1]
		if _, err := fmt.Sscanf(s[:len(s)-1], "%d", &n); err != nil || n <= 0 {
			return 0, fmt.Errorf("--since expects a window like 30d, 12w, or 720h, got %q", s)
		}
		if unit == 'w' {
			n *= 7
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("--since expects a window like 30d, 12w, or 720h, got %q", s)
	}
	return d, nil
}

func printSummary(a *app, rep *metrics.Report) {
	fmt.Fprintf(a.out, "%s completed in the last %gd (since %s)\n", workitem.Folder(rep.Type), rep.WindowDays, rep.WindowStart[:10])
	printGroup(a, rep.Summary)
	for _, g := range rep.Groups {
		fmt.Fprintf(a.out, "\n[%s]\n", g.Group)
		printGroup(a, g)
	}
	if len(rep.Aging) > 0 {
		fmt.Fprintln(a.out, "\naging WIP (age since started, vs cycle time p85):")
		for _, ag := range rep.Aging {
			flag := ""
			if ag.OverP85 {
				flag = "  OVER P85"
			}
			if ag.Blocked {
				flag += "  BLOCKED"
			}
			fmt.Fprintf(a.out, "  %-6s %-12s %-6s %s%s\n", ag.ID, ag.Status, ag.AgeHuman, ag.Title, flag)
		}
	}
	if len(rep.Throughput) > 0 {
		fmt.Fprintln(a.out, "\nthroughput by week:")
		for _, w := range rep.Throughput {
			fmt.Fprintf(a.out, "  %s (%s)  %d\n", w.Week, w.Start, w.Done)
		}
	}
}

func printGroup(a *app, s metrics.Summary) {
	fmt.Fprintf(a.out, "  completed %d · cancelled %d (%.0f%%) · throughput %.1f/week · WIP now %d\n",
		s.Completed, s.Cancelled, s.CancellationRate*100, s.ThroughputWeek, s.WIP)
	dist := func(name string, d metrics.Distribution) {
		if d.Count == 0 {
			fmt.Fprintf(a.out, "  %-12s no data\n", name)
			return
		}
		fmt.Fprintf(a.out, "  %-12s p50 %s · p85 %s · max %s · mean %s (n=%d)\n", name,
			metrics.Human(d.P50), metrics.Human(d.P85), metrics.Human(d.Max), metrics.Human(d.Mean), d.Count)
	}
	dist("cycle time", s.CycleTime)
	dist("lead time", s.LeadTime)
	dist("queue time", s.QueueTime)
	if s.CycleTime.Count > 0 {
		fmt.Fprintf(a.out, "  flow efficiency %.0f%%\n", s.FlowEfficiency*100)
	}
	if len(s.InStateShare) > 0 {
		var parts []string
		states := make([]string, 0, len(s.InStateShare))
		for st := range s.InStateShare {
			states = append(states, st)
		}
		sort.Slice(states, func(i, j int) bool { return stateRank(states[i]) < stateRank(states[j]) })
		for _, st := range states {
			parts = append(parts, fmt.Sprintf("%s %.0f%%", st, s.InStateShare[st]*100))
		}
		fmt.Fprintf(a.out, "  time in state: %s\n", strings.Join(parts, " · "))
	}
}

func stateRank(s string) int {
	for i, st := range workitem.States {
		if st == s {
			return i
		}
	}
	return len(workitem.States)
}

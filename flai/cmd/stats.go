package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newStatsCmd(a *app) *cobra.Command {
	var since, typ, by, bucket string
	c := &cobra.Command{
		Use:   "stats",
		Short: "Print flow metrics: throughput, cycle time, WIP, flow efficiency, time in state, tokens and cost",
		Long: `The reference implementation of design/system/metrics.md. Aggregates cover
items completed in the window (default 30d); WIP and aging are as of now.
Items that carry usage add what agents spent on those done in the window:
tokens, cost, and agent time, what that comes to per item, per minute of
agent work, and per dollar, and the same per model. --json includes per-item
values of every item of the type, throughput for each week of the window,
burn-up and cumulative flow for each day of it, aging items, and usage: items
done against time and cost, and under usage.spend what was
spent on epics, on stories, and on tasks over time, one point per --bucket
(hour, day, or week) from the first with spend to now, for dashboards and
scripts. A bucket of an hour needs a window of 31 days or less.`,
		Example: `  flai stats
  flai stats --since 90d --by nature
  flai stats --type task --json
  flai stats --since 7d --bucket hour --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			window, err := metrics.ParseWindow(since)
			if err != nil {
				return err
			}
			if err := metrics.CheckBucket(bucket, window); err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			rep := metrics.Compute(items, metrics.Options{Now: a.now(), Since: window, Type: typ, By: by, Bucket: bucket})
			if a.jsonOut {
				return a.printJSON(rep)
			}
			printSummary(a, rep)
			return nil
		},
	}
	c.Flags().StringVar(&since, "since", metrics.DefaultWindow, "window for completed items: 7d, 90d, 12w, 720h")
	c.Flags().StringVar(&typ, "type", workitem.Story, "item type: epic, story, task")
	c.Flags().StringVar(&by, "by", "", "group by nature, type, or parent")
	c.Flags().StringVar(&bucket, "bucket", metrics.DefaultBucket, "what --json lays spend over time out in: hour, day, or week")
	return c
}

func printSummary(a *app, rep *metrics.Report) {
	fmt.Fprintf(a.out, "%s completed in the last %gd (since %s)\n", workitem.Folder(rep.Type), rep.WindowDays, rep.WindowStart[:10])
	printGroup(a, rep.Summary)
	printUsage(a, rep.Type, rep.Usage)
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

// printUsage prints what agents spent on the items done in the window: in
// total, per item, per minute of agent work, and per dollar, and per model
// (S-0143, S-0163).
func printUsage(a *app, typ string, u metrics.UsageReport) {
	if u.Items == 0 {
		return
	}
	estimated := ""
	if u.Estimated {
		estimated = " (estimated in part)"
	}
	fmt.Fprintf(a.out, "  usage: %s tokens · $%.2f%s · %s of agent work, over %d done\n",
		usage.Count(u.Tokens), u.Cost, estimated, (time.Duration(u.Seconds) * time.Second).String(), u.Items)
	if s := u.Spend[typ]; s != nil && s.Items > 0 {
		fmt.Fprintf(a.out, "    per %s %s\n", typ, strings.Join(averages(s.Spend), " · "))
	}
	for _, m := range u.Models {
		rate := ""
		if m.TokensPerMinute != nil {
			rate = fmt.Sprintf(" · %s tokens/min", usage.Count(int64(*m.TokensPerMinute)))
		}
		fmt.Fprintf(a.out, "    %s  %s tokens · $%.2f%s (%d)\n", m.Model, usage.Count(m.Tokens), m.Cost, rate, m.Items)
	}
}

// averages says what a spend comes to per item, then per minute of agent
// work and per dollar, leaving out what has no divisor.
func averages(s metrics.Spend) []string {
	var parts []string
	if s.TokensPerItem != nil && s.CostPerItem != nil {
		parts = append(parts, fmt.Sprintf("%s tokens, $%.2f", usage.Count(int64(*s.TokensPerItem)), *s.CostPerItem))
	}
	if s.TokensPerMinute != nil {
		parts = append(parts, fmt.Sprintf("per agent minute %s tokens", usage.Count(int64(*s.TokensPerMinute))))
	}
	if s.TokensPerDollar != nil {
		parts = append(parts, fmt.Sprintf("per dollar %s tokens", usage.Count(int64(*s.TokensPerDollar))))
	}
	return parts
}

func stateRank(s string) int {
	for i, st := range workitem.States {
		if st == s {
			return i
		}
	}
	return len(workitem.States)
}

package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/statsread"
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
scripts. A bucket of an hour needs a window of 31 days or less. It also
carries forecasts (forecast, delivery, and estimate error), cost_of_delay
(outstanding per column and incurred, by the day and week), waiting (on
threads and in review, by the week), claims (in progress by the day against
the limit, and each story's touches against the files its commits changed),
and strategic_days (the strategic agents' cost and time beside delivery).
What strategic agents spent on items is reported apart from what agents did,
under usage.strategic and each item's usage.strategic, beside the project's
usage.cost_per_agent_hour and each item's expected_cost (ADR-0083).
Touches drift needs git; without it flai stats warns and leaves it out.`,
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
			items, opt, err := statsread.Read(a.runner, repo, a.logger())
			if err != nil {
				return err
			}
			opt.Now, opt.Since, opt.Type, opt.By, opt.Bucket = a.now(), window, typ, by, bucket
			rep := metrics.Compute(items, opt)
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
	printStrategic(a, rep.Strategic)
	printStrategicDays(a, rep)
	printForecasts(a, rep.Forecasts)
	printCostOfDelay(a, rep)
	printWaiting(a, rep.Waiting)
	printClaims(a, rep)
}

// printForecasts prints how far forecasts and estimates were from what
// happened, for each kind of error some item has (S-0205).
func printForecasts(a *app, f metrics.Forecasts) {
	header := false
	for _, e := range []struct {
		name string
		s    metrics.ErrorStats
	}{{"forecast", f.Forecast}, {"delivery", f.Delivery}, {"estimate", f.Estimate}} {
		if e.s.Count == 0 || e.s.P50 == nil || e.s.P85 == nil {
			continue
		}
		if !header {
			fmt.Fprintln(a.out, "\nforecast and estimate error (absolute):")
			header = true
		}
		fmt.Fprintf(a.out, "  %-12s p50 %s · p85 %s (n=%d)\n", e.name, metrics.Human(*e.s.P50), metrics.Human(*e.s.P85), e.s.Count)
	}
}

// printCostOfDelay prints, when any item has a cost of delay, what is
// outstanding in each column now and what waiting cost over the window
// (S-0205).
func printCostOfDelay(a *app, rep *metrics.Report) {
	valued := false
	for _, m := range rep.Items {
		valued = valued || m.CostOfDelay != nil
	}
	days := rep.CostOfDelay.Days
	if !valued || len(days) == 0 {
		return
	}
	last := days[len(days)-1].Outstanding
	cols := make([]string, 0, len(last))
	for c := range last {
		cols = append(cols, c)
	}
	sort.Slice(cols, func(i, j int) bool { return stateRank(cols[i]) < stateRank(cols[j]) })
	parts := make([]string, 0, len(cols))
	for _, c := range cols {
		parts = append(parts, fmt.Sprintf("%s %.2f", c, last[c]))
	}
	incurred := 0.0
	for _, d := range days {
		incurred += d.Incurred
	}
	fmt.Fprintln(a.out, "\ncost of delay (per week of waiting):")
	fmt.Fprintf(a.out, "  outstanding now  %s\n", strings.Join(parts, " · "))
	fmt.Fprintf(a.out, "  incurred in the window %.2f\n", incurred)
}

// printWaiting prints how long the agents of the items completed in the
// window waited on threads and in review, in all and per item (S-0205).
func printWaiting(a *app, w metrics.Waiting) {
	var items int
	var threads, review float64
	for _, wk := range w.Weeks {
		items += wk.Items
		threads += wk.Threads.Total
		review += wk.Review.Total
	}
	if items == 0 {
		return
	}
	fmt.Fprintf(a.out, "\nwaiting, over %d completed:\n", items)
	fmt.Fprintf(a.out, "  %-12s total %s · mean %s\n", "on threads", metrics.Human(threads), metrics.Human(threads/float64(items)))
	fmt.Fprintf(a.out, "  %-12s total %s · mean %s\n", "in review", metrics.Human(review), metrics.Human(review/float64(items)))
}

// printClaims prints the time the stories completed in the window were held
// in ready, the items in progress now against the limit, and the touches
// drift when git could be read (S-0205).
func printClaims(a *app, rep *metrics.Report) {
	var lines []string
	start, _ := time.Parse(workitem.TimeFormat, rep.WindowStart)
	var held float64
	var stories int
	for _, m := range rep.Items {
		done, _ := time.Parse(workitem.TimeFormat, m.Completed)
		if m.HeldSeconds == nil || m.Status != workitem.Done || done.Before(start) {
			continue
		}
		held += *m.HeldSeconds
		stories++
	}
	if stories > 0 {
		lines = append(lines, fmt.Sprintf("  held in ready  total %s · mean %s, over %d completed", metrics.Human(held), metrics.Human(held/float64(stories)), stories))
	}
	c := rep.Claims
	busy := false
	for _, d := range c.Days {
		busy = busy || d.InProgress > 0
	}
	if n := len(c.Days); busy {
		limit := ""
		if c.Limit > 0 {
			limit = fmt.Sprintf(" of a limit of %d", c.Limit)
		}
		lines = append(lines, fmt.Sprintf("  in progress now %d%s", c.Days[n-1].InProgress, limit))
	}
	if c.Drift != nil && len(*c.Drift) > 0 {
		var outside, unchanged int
		for _, d := range *c.Drift {
			outside += d.OutsideCount
			unchanged += d.UnchangedCount
		}
		lines = append(lines, fmt.Sprintf("  touches drift  %d stories with commits · %d files outside their touches · %d touches unchanged", len(*c.Drift), outside, unchanged))
	}
	if len(lines) == 0 {
		return
	}
	fmt.Fprintln(a.out, "\nclaims:")
	fmt.Fprintln(a.out, strings.Join(lines, "\n"))
}

// printStrategicDays prints what the strategic agents spent over the window
// beside the items completed in it (S-0205).
func printStrategicDays(a *app, rep *metrics.Report) {
	var cost float64
	var seconds int64
	var completed int
	for _, d := range rep.StrategicDays {
		cost += d.Cost
		seconds += d.Seconds
		completed += d.Completed
	}
	if cost == 0 && seconds == 0 {
		return
	}
	per := ""
	if u := rep.Usage; u.Items > 0 {
		per = fmt.Sprintf(" (usage $%.2f per item)", u.Cost/float64(u.Items))
	}
	fmt.Fprintln(a.out, "\nstrategic agents in the window:")
	fmt.Fprintf(a.out, "  $%.2f · %s, beside %d completed%s\n", cost, metrics.Human(float64(seconds)), completed, per)
}

// printStrategic prints each strategic agent's totals as its activity
// document holds them, all time (ADR-0079), with what of them the items carry
// and the project strategic total (ADR-0095).
func printStrategic(a *app, agents []metrics.StrategicAgent) {
	if len(agents) == 0 {
		return
	}
	fmt.Fprintln(a.out, "\nstrategic agents (all time):")
	for _, s := range agents {
		noun := "activities"
		if s.Activities == 1 {
			noun = "activity"
		}
		last := ""
		if s.LastRun != "" {
			last = ", last " + s.LastRun
		}
		fmt.Fprintf(a.out, "  %s: %d %s, %.4f USD, %d s (on items %.4f USD, %d s; project %.4f USD, %d s)%s\n",
			s.Kind, s.Activities, noun, s.Cost, s.Seconds, s.Items.Cost, s.Items.Seconds, s.Project.Cost, s.Project.Seconds, last)
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
// (S-0143, S-0163); then apart what strategic agents spent on them, and the
// project's cost per agent hour (ADR-0083).
func printUsage(a *app, typ string, u metrics.UsageReport) {
	if u.Items > 0 {
		printAgentUsage(a, typ, u)
	}
	if s := u.Strategic; s.Items > 0 {
		kinds := make([]string, 0, len(s.Kinds))
		for _, k := range s.Kinds {
			kinds = append(kinds, fmt.Sprintf("%s $%.2f", k.Kind, k.Cost))
		}
		fmt.Fprintf(a.out, "  strategic usage, apart: %s tokens · $%.2f (estimated) · %s of strategic agent work, over %d done (%s)\n",
			usage.Count(s.Tokens), s.Cost, (time.Duration(s.Seconds) * time.Second).String(), s.Items, strings.Join(kinds, " · "))
	}
	if r := u.CostPerAgentHour; r != nil {
		fmt.Fprintf(a.out, "  cost per agent hour $%.2f, over every story measured from its logs\n", *r)
	}
}

// printAgentUsage prints what agents spent on the items done in the window.
func printAgentUsage(a *app, typ string, u metrics.UsageReport) {
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

// averages says what a spend comes to per item, with the agent time an item
// took, then per minute of agent work and per dollar, leaving out what has
// no divisor.
func averages(s metrics.Spend) []string {
	var parts []string
	if s.TokensPerItem != nil && s.CostPerItem != nil {
		per := fmt.Sprintf("%s tokens, $%.2f", usage.Count(int64(*s.TokensPerItem)), *s.CostPerItem)
		if s.MinutesPerItem != nil {
			per += fmt.Sprintf(", %s of agent work", metrics.Human(*s.MinutesPerItem*60))
		}
		parts = append(parts, per)
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

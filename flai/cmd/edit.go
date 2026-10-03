package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
)

// flai edit changes a work item after it was made (S-0085).
func newEditCmd(a *app) *cobra.Command {
	var title, nature, parent, hash, message, byFlag, harness, model string
	var agentConfig []string
	var rf roleFlags
	var tags, touches, topics, after, trailers []string
	var clearTags, clearTouches, clearTopics, clearAfter, clearAgent, bodyStdin, autocommit, show bool
	var draft, noDraft, clearCost, clearForecast bool
	var revenue, penalty, timeLost, costValue, duration, delivery, basis string
	c := &cobra.Command{
		Use:   "edit <id>",
		Short: "Change an item's title, nature, tags, topics, touches, after, parent, planning data, or body, checked and in one step",
		Long: `Change what an item says about itself. Any of the fields and the body can
change together. What is the item's state stays flai's and is changed by its
own commands: status by flai move, blocking by flai block, never here.

A retitle keeps everything that carries the title in step: the front matter,
the heading, the file's name, the line in the parent's list, the story's
narrative, and links to the old file name in design, docs, and wip. A new
parent must be an open item of the right type; the item leaves the old
parent's list and joins the new one. An archived or closed item is refused.

--topics names what a story or epic is about beyond the components its tags
and touches reach, such as logging or release (ADR-0047); flai show prints a
story's topics with where each came from.

--after names the stories a story waits for: while any of them is not done,
the story is held in ready, and flai serve and wait_for_work pass it over
(ADR-0046). On a task it names the tasks of the same story the task waits for
(S-0176). flai check refuses an entry that does not exist, a task of another
story, and a cycle.

Planning data (S-0199). --draft makes a backlog story a draft and
--no-draft finalizes one, which may then go to ready. A cost of delay, on a
story or an epic, has inputs (--revenue-per-week, --penalty-per-week, as
amounts in planning.currency, and --time-lost-per-cycle, a Go duration) and
a value per week (--cost-of-delay-value). A story's forecast has a duration
(--forecast-duration, a Go duration), a delivery (--forecast-delivery, a UTC
timestamp like 2026-10-09T17:00:00Z), and a basis (--forecast-basis, one
sentence). Each flag changes its key only: an empty value removes it, and
removing the last input or value, or the last of duration and delivery,
removes the block. --clear-cost-of-delay and --clear-forecast remove a
block. A block that changes records who changed it (--by) and when.

--body-stdin reads what lies below the heading; the heading is the ID and the
title, and flai writes it. With --hash, the hash flai edit --show printed, a
change someone made meanwhile is a conflict (exit 3) and nothing is written.
flai check runs with the change in place: if it reports anything the change
introduces, every file is put back and the findings are printed (exit 4).
--autocommit commits every file the edit touched in one commit, unless the
project sets dashboard.autocommit: false. Nothing is pushed.

Agents connected over MCP are told of an edit someone else made, and what of
the item changed.`,
		Example: `  flai edit S-0085 --show
  flai edit S-0085 --title "Items are editable from the dashboard" --autocommit
  flai edit S-0085 --tag dashboard --tag cli --touches flaiover/src
  flai edit S-0085 --parent E-0004
  flai edit S-0130 --after S-0128,S-0129
  flai edit T-0042 --after T-0040,T-0041
  flai edit S-0135 --topics logging,release
  flai edit S-0199 --no-draft
  flai edit S-0199 --revenue-per-week 1200 --time-lost-per-cycle 4h
  flai edit S-0199 --penalty-per-week ""
  flai edit S-0199 --forecast-duration 6h --forecast-basis "three tasks like S-0185's"
  flai edit S-0085 --body-stdin --hash 3f0c... < body.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if show {
				v, err := itemedit.Show(repo, args[0])
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(v)
				}
				fmt.Fprintf(a.out, "%s %s\n  %s, %s, %s\n  tags: %s\n", v.ID, v.Title, v.Type, v.Nature, v.Status, strings.Join(v.Tags, ", "))
				if v.Type != "task" {
					fmt.Fprintf(a.out, "  topics: %s\n", strings.Join(v.Topics, ", "))
				}
				if v.Type != "epic" {
					fmt.Fprintf(a.out, "  touches: %s\n  parent: %s\n", strings.Join(v.Touches, ", "), v.Parent)
				}
				if v.Type != "epic" {
					fmt.Fprintf(a.out, "  after: %s\n", strings.Join(v.After, ", "))
				}
				if v.Type == "story" {
					fmt.Fprintf(a.out, "  agent: %s (project default: %s)\n", v.Agent, v.DefaultAgent)
				}
				for _, l := range planningLines(v.Draft, v.CostOfDelay, v.Forecast, v.Currency) {
					fmt.Fprintf(a.out, "  %s\n", l)
				}
				fmt.Fprintf(a.out, "  file: %s\n  hash: %s\n", v.Path, v.Hash)
				if !v.Editable {
					fmt.Fprintf(a.out, "  not editable: %s\n", v.Reason)
				}
				fmt.Fprintf(a.out, "\n%s", v.Body)
				return nil
			}
			var ch itemedit.Change
			f := cmd.Flags()
			if f.Changed("title") {
				ch.Title = &title
			}
			if f.Changed("nature") {
				ch.Nature = &nature
			}
			if f.Changed("parent") {
				ch.Parent = &parent
			}
			switch {
			case clearTags && f.Changed("tag"):
				return fmt.Errorf("--tag and --clear-tags contradict each other")
			case clearTags:
				ch.Tags = &[]string{}
			case f.Changed("tag"):
				ch.Tags = &tags
			}
			switch {
			case clearTouches && f.Changed("touches"):
				return fmt.Errorf("--touches and --clear-touches contradict each other")
			case clearTouches:
				ch.Touches = &[]string{}
			case f.Changed("touches"):
				ch.Touches = &touches
			}
			switch {
			case clearTopics && f.Changed("topics"):
				return fmt.Errorf("--topics and --clear-topics contradict each other")
			case clearTopics:
				ch.Topics = &[]string{}
			case f.Changed("topics"):
				ch.Topics = &topics
			}
			switch {
			case clearAfter && f.Changed("after"):
				return fmt.Errorf("--after and --clear-after contradict each other")
			case clearAfter:
				ch.After = &[]string{}
			case f.Changed("after"):
				ch.After = &after
			}
			if clearAgent || f.Changed("harness") || f.Changed("model") || f.Changed("agent-config") || rf.given() {
				next, err := editedAgent(repo, args[0], clearAgent, harness, model, agentConfig, rf)
				if err != nil {
					return err
				}
				if next == nil {
					ch.ClearAgent = true
				} else {
					ch.Agent = next
				}
			}
			if err := planningChange(f, &ch, draft, noDraft, clearCost, clearForecast, map[string]*string{
				"revenue-per-week": &revenue, "penalty-per-week": &penalty, "time-lost-per-cycle": &timeLost, "cost-of-delay-value": &costValue,
				"forecast-duration": &duration, "forecast-delivery": &delivery, "forecast-basis": &basis,
			}); err != nil {
				return err
			}
			if bodyStdin {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				body := string(data)
				ch.Body = &body
			}
			if ch == (itemedit.Change{}) {
				return fmt.Errorf("nothing to change: give --title, --nature, --tag, --topics, --clear-topics, --touches, --after, --clear-after, --parent, --harness, --model, --agent-config, a --role- flag, --unset-role, --clear-agent, --draft, --no-draft, a cost of delay or forecast flag, or --body-stdin (flai edit %s --show prints what is there)", args[0])
			}
			by, _ := agentIdentity()
			if cfg, _, err := a.loadConfig(); err == nil && by == "agent" && cfg.Author != "" {
				by = cfg.Author
			}
			if byFlag != "" {
				by = byFlag
			}
			res, err := itemedit.Apply(repo, a.runner, args[0], ch, itemedit.Options{Hash: hash, By: by, Message: message, Trailers: trailers, NoCommit: !autocommit, Now: a.now()})
			if c, ok := docedit.IsConflict(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"conflict": c})
				}
				return &exitError{code: exitDocConflict, msg: c.Error()}
			}
			if r, ok := docedit.IsRefused(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": r})
				} else {
					for _, fd := range r.Findings {
						fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", fd.Path, fd.Line, fd.Level, fd.Rule, fd.Message)
					}
				}
				return &exitError{code: exitDocRefused, msg: r.Error()}
			}
			var inv *itemedit.InvalidError
			if errors.As(err, &inv) {
				// a rule, as flai move says its refusals: the caller's mistake, which a
				// dashboard answers with 400 and not with 500
				return fmt.Errorf("rule: %s", inv.Error())
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			if res.Unchanged {
				fmt.Fprintf(a.out, "%s is unchanged\n", res.ID)
				return nil
			}
			fmt.Fprintf(a.out, "%s: changed %s\n  %s\n", res.ID, strings.Join(res.Changed, ", "), res.Path)
			if res.Renamed != "" {
				fmt.Fprintf(a.out, "  renamed from %s\n", res.Renamed)
			}
			switch {
			case res.Committed:
				fmt.Fprintf(a.out, "  committed %s (%d file(s))\n", res.Commit, len(res.Files))
			case res.CommitError != "":
				fmt.Fprintf(a.out, "  changed, not committed: %s\n", res.CommitError)
			}
			return nil
		},
	}
	f := c.Flags()
	f.BoolVar(&show, "show", false, "print the fields, the body, the hash, and whether the item may be edited")
	f.StringVar(&title, "title", "", "the new title")
	f.StringVar(&nature, "nature", "", "one of feature, improvement, remediation, research, experiment")
	f.StringSliceVar(&tags, "tag", nil, "the tags, replacing the ones there (repeatable or comma separated)")
	f.BoolVar(&clearTags, "clear-tags", false, "remove every tag")
	f.StringSliceVar(&touches, "touches", nil, "the paths or components the work changes, replacing the ones there")
	f.BoolVar(&clearTouches, "clear-touches", false, "remove the list")
	f.StringSliceVar(&topics, "topics", nil, "what a story or epic is about, such as logging or release, replacing the ones there")
	f.BoolVar(&clearTopics, "clear-topics", false, "remove a story's or epic's topics")
	f.StringSliceVar(&after, "after", nil, "what it waits for until they are done, replacing the ones there: a story's stories, a task's tasks of the same story")
	f.BoolVar(&clearAfter, "clear-after", false, "remove what a story or a task waits for")
	f.StringVar(&parent, "parent", "", "the new parent: an epic for a story, a story for a task")
	f.StringVar(&harness, "harness", "", "a story's agent: the harness that runs it")
	f.StringVar(&model, "model", "", "a story's agent: the model it runs")
	f.StringArrayVar(&agentConfig, "agent-config", nil, "a story's agent: an option, key=value, and key= to remove one (repeatable)")
	f.BoolVar(&clearAgent, "clear-agent", false, "remove the story's agent; with --harness, --model, or --agent-config, replace it with exactly those")
	rf.register(c, true)
	f.BoolVar(&draft, "draft", false, "make a story in the backlog a draft, which must be finalized before it is ready")
	f.BoolVar(&noDraft, "no-draft", false, "finalize a draft story, so that it may go to ready")
	f.StringVar(&revenue, "revenue-per-week", "", "cost of delay input: revenue each week it is done brings, in planning.currency; empty removes it")
	f.StringVar(&penalty, "penalty-per-week", "", "cost of delay input: what each week it is not done costs beyond revenue; empty removes it")
	f.StringVar(&timeLost, "time-lost-per-cycle", "", "cost of delay input: work lost each cycle it is not done, a Go duration; empty removes it")
	f.StringVar(&costValue, "cost-of-delay-value", "", "the cost of delay per week, in planning.currency; empty removes it")
	f.BoolVar(&clearCost, "clear-cost-of-delay", false, "remove the cost of delay")
	f.StringVar(&duration, "forecast-duration", "", "a story's forecast: the work it is expected to take, a Go duration; empty removes it")
	f.StringVar(&delivery, "forecast-delivery", "", "a story's forecast: when it is expected done, a UTC timestamp; empty removes it")
	f.StringVar(&basis, "forecast-basis", "", "a story's forecast: what it rests on, in one sentence; empty removes it")
	f.BoolVar(&clearForecast, "clear-forecast", false, "remove the story's forecast")
	f.BoolVar(&bodyStdin, "body-stdin", false, "read the body below the heading from standard input")
	f.StringVar(&hash, "hash", "", "the hash flai edit --show printed; a change made meanwhile is then a conflict")
	f.StringVar(&byFlag, "by", "", "who edits, as agents are told (default: FLAI_AGENT, then the config author)")
	f.StringVar(&message, "message", "", "commit subject after the prefix (default names what changed)")
	f.BoolVar(&autocommit, "autocommit", false, "commit every file the edit touched, unless dashboard.autocommit is false")
	f.StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	return c
}

// planningChange puts the draft, cost of delay, and forecast flags given on
// the command line into ch (S-0199); flags names each string flag's value.
func planningChange(f *pflag.FlagSet, ch *itemedit.Change, draft, noDraft, clearCost, clearForecast bool, flags map[string]*string) error {
	if draft && noDraft {
		return fmt.Errorf("--draft and --no-draft contradict each other")
	}
	if draft || noDraft {
		ch.Draft = &draft
	}
	given := func(name string) *string {
		if f.Changed(name) {
			return flags[name]
		}
		return nil
	}
	cost := itemedit.CostOfDelayEdit{RevenuePerWeek: given("revenue-per-week"), PenaltyPerWeek: given("penalty-per-week"),
		TimeLostPerCycle: given("time-lost-per-cycle"), Value: given("cost-of-delay-value")}
	if cost != (itemedit.CostOfDelayEdit{}) {
		if clearCost {
			return fmt.Errorf("--clear-cost-of-delay and a cost of delay flag contradict each other; give one or the other")
		}
		ch.CostOfDelay = &cost
	}
	ch.ClearCostOfDelay = clearCost
	fc := itemedit.ForecastEdit{Duration: given("forecast-duration"), Delivery: given("forecast-delivery"), Basis: given("forecast-basis")}
	if fc != (itemedit.ForecastEdit{}) {
		if clearForecast {
			return fmt.Errorf("--clear-forecast and a --forecast- flag contradict each other; give one or the other")
		}
		ch.Forecast = &fc
	}
	ch.ClearForecast = clearForecast
	return nil
}

package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newItemCmd(a *app, typ string) *cobra.Command {
	c := &cobra.Command{
		Use:   typ,
		Short: fmt.Sprintf("Create %s (flai show prints one, flai move transitions it)", pluralType(typ)),
	}
	c.AddCommand(newItemNewCmd(a, typ))
	if typ == workitem.Task {
		// a task is closed in one call (ADR-0107); a story or an epic is not
		c.AddCommand(newTaskDoneCmd(a))
	}
	return c
}

func newItemNewCmd(a *app, typ string) *cobra.Command {
	var nature, owner, parent, harness, model, revenue, penalty, timeLost string
	var tags, touches, topics, after, trailers, agentConfig []string
	var rf roleFlags
	var bodyStdin, autocommit, printBody, draft bool
	parentFlag := map[string]string{workitem.Story: "epic", workitem.Task: "story"}[typ]
	c := &cobra.Command{
		Use:   "new \"<title>\"",
		Short: fmt.Sprintf("Create %s from the item template", withArticle(typ)),
		Long: fmt.Sprintf(`Create %s from the project's item template with the next free ID,
linked into its parent.

With --body-stdin the body below the item's heading is read from standard
input instead of the template's empty sections, and the creation is one step
that happens or does not: flai check runs with the new item in place, and if
it reports anything the item introduces, the item is removed, its parent is
restored, and the findings are printed (exit 4). --autocommit commits the new
item and its parent on their own, unless the project sets
dashboard.autocommit: false. Nothing is pushed. --print-body prints the body
the template gives, for a form or a script to start from, and creates nothing.%s%s%s`, withArticle(typ), afterHelp(typ), draftHelp(typ), costHelp(typ)),
		Args: func(cmd *cobra.Command, args []string) error {
			if printBody {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if printBody {
				body, err := repo.TemplateBody(typ)
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(map[string]string{"type": typ, "body": body})
				}
				fmt.Fprint(a.out, body)
				return nil
			}
			agent, err := agentFlags(harness, model, agentConfig, rf)
			if err != nil {
				return err
			}
			opt := workitem.NewOptions{
				Type: typ, Title: args[0], Nature: nature, Parent: parent,
				Owner: orDefault(owner, a.author()), Tags: tags, Touches: touches, Topics: topics, After: after, Agent: agent, Draft: draft, Now: a.now(),
			}
			// the inputs are set by the owner, whom the dashboard names as the operator
			if opt.CostOfDelay, err = newCostOfDelay(revenue, penalty, timeLost, repo.Manifest.Planning.CurrencyCode(), opt.Owner, opt.Now); err != nil {
				return err
			}
			// a task's touches may grow its story's claim into another's (I-0059)
			var watch *itemedit.ClaimWatch
			if typ == workitem.Task {
				watch = itemedit.WatchClaim(repo, parent)
			}
			// an after: entry that names nothing, or forms a cycle, is the
			// check's to find, so a creation that sets one is checked
			if bodyStdin || autocommit || len(after) > 0 {
				if bodyStdin {
					data, err := io.ReadAll(cmd.InOrStdin())
					if err != nil {
						return err
					}
					if strings.TrimSpace(string(data)) == "" {
						return fmt.Errorf("--body-stdin was given and standard input is empty; write the body, or leave the flag out for the template's empty sections")
					}
					opt.Body = string(data)
				}
				res, err := itemnew.Create(repo, a.runner, itemnew.Options{New: opt, Autocommit: autocommit, Trailers: trailers})
				if r, ok := docedit.IsRefused(err); ok {
					if a.jsonOut {
						_ = a.printJSON(map[string]any{"refused": r})
					} else {
						for _, f := range r.Findings {
							fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
						}
					}
					return &exitError{code: exitDocRefused, msg: r.Error()}
				}
				if err != nil {
					return err
				}
				told := a.grown(watch, res.Item.ID, a.writer())
				if a.jsonOut {
					return a.printJSON(struct {
						*itemnew.Result
						Overlaps []itemedit.Overlapping `json:"overlaps,omitempty"`
					}{res, told})
				}
				fmt.Fprintf(a.out, "%s %s\n  %s\n", res.Item.ID, res.Item.Title, res.Path)
				switch {
				case res.Committed:
					fmt.Fprintf(a.out, "  committed %s\n", res.Commit)
				case res.CommitError != "":
					fmt.Fprintf(a.out, "  NOT committed (%s)\n", firstLine(res.CommitError))
				}
				printOverlapping(a.out, told)
				return nil
			}
			it, err := repo.Create(opt)
			if err != nil {
				return err
			}
			told := a.grown(watch, it.ID, a.writer())
			if a.jsonOut {
				return a.printJSON(struct {
					*workitem.Item
					Overlaps []itemedit.Overlapping `json:"overlaps,omitempty"`
				}{it, told})
			}
			fmt.Fprintf(a.out, "%s %s\n  %s\n", it.ID, it.Title, relPath(repo.Root, it.Path))
			printOverlapping(a.out, told)
			return nil
		},
	}
	c.Flags().StringVar(&nature, "nature", "feature", "one of "+strings.Join(workitem.Natures, ", "))
	c.Flags().StringVar(&owner, "owner", "", "owner (default: config author)")
	c.Flags().StringSliceVar(&tags, "tag", nil, "tag (repeatable or comma separated)")
	c.Flags().StringSliceVar(&touches, "touches", nil, "paths or components this work changes (repeatable or comma separated)")
	if typ != workitem.Task {
		// what it is about beyond its components (S-0135, ADR-0047)
		c.Flags().StringSliceVar(&topics, "topics", nil, "topics the "+typ+" is about beyond the components it reaches, such as logging or release (repeatable or comma separated)")
		// what waiting for it costs (S-0204), the inputs the edit command also sets
		c.Flags().StringVar(&revenue, "revenue-per-week", "", "cost of delay input: revenue each week it is done brings, in planning.currency")
		c.Flags().StringVar(&penalty, "penalty-per-week", "", "cost of delay input: what each week it is not done costs beyond revenue, in planning.currency")
		c.Flags().StringVar(&timeLost, "time-lost-per-cycle", "", "cost of delay input: work lost each cycle it is not done, a Go duration")
	}
	switch typ {
	case workitem.Story:
		c.Flags().StringSliceVar(&after, "after", nil, "the stories this story waits for until they are done (comma separated); checked before it is kept")
	case workitem.Task:
		c.Flags().StringSliceVar(&after, "after", nil, "the tasks of the same story this task waits for until they are done (comma separated); checked before it is kept")
	}
	c.Flags().BoolVar(&bodyStdin, "body-stdin", false, "read the body below the heading from standard input; checked before it is kept")
	c.Flags().BoolVar(&autocommit, "autocommit", false, "commit the new item and its parent on their own, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	c.Flags().BoolVar(&printBody, "print-body", false, "print the body the template gives this type and create nothing")
	if typ == workitem.Story {
		// who works it (S-0103), over the project's default (flai agent)
		c.Flags().StringVar(&harness, "harness", "", "the harness that runs the agent for this story, over the project's default")
		c.Flags().StringVar(&model, "model", "", "the model it runs, over the project's default")
		c.Flags().StringArrayVar(&agentConfig, "agent-config", nil, "an option for the harness, key=value, over the project's default (repeatable)")
		rf.register(c, false)
		c.Flags().BoolVar(&draft, "draft", false, "make the story a draft, which must be finalized before it is ready")
	}
	if parentFlag != "" {
		help := "parent " + parentFlag + " ID"
		if typ == workitem.Story {
			help += " (optional: a story need not belong to one, S-0092)"
		}
		c.Flags().StringVar(&parent, parentFlag, "", help)
		if typ != workitem.Story {
			c.PreRunE = func(cmd *cobra.Command, args []string) error {
				if !printBody && parent == "" {
					return fmt.Errorf("required flag \"%s\" not set", parentFlag)
				}
				return nil
			}
		}
	}
	return c
}

func newShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Print one work item with its children and history",
		Long: `Print one work item with its children and history.

A story with tasks also gets its task plan (S-0176): each task's state,
ready to start, waiting with the tasks of its after it waits for, in
progress, done, or cancelled; and the layers, the open and done tasks
grouped by the longest chain of after steps before each, so that the tasks
of a layer can run at once when those of the layers before it are done. A
task on a cycle of after, or waiting on one, is in no layer. --json gives
the plan as plan, beside item and children.

An item with a forecast duration, or else an estimate, gets its expected
cost: that duration priced at the project's cost per agent hour, an
estimate, absent while no story has been measured from its logs
(ADR-0083). --json gives it as expected_cost. What strategic agents such as
the planner spent on the item is printed per kind, apart from what its
agents spent.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			items, err := repo.List(true)
			if err != nil {
				return err
			}
			children := workitem.Children(items, it.ID)
			// a story's topics and where each came from (S-0135, ADR-0047)
			var storyTopics []topics.StoryTopic
			if it.Type == workitem.Story {
				if storyTopics, err = topics.ForStory(repo, it.ID); err != nil {
					return err
				}
			}
			// a story's task plan (S-0176)
			var plan *workitem.Plan
			if it.Type == workitem.Story {
				plan = workitem.PlanOf(children, it.ID)
			}
			// what the item is expected to cost, priced only when it has a
			// duration to price (ADR-0083)
			var rate *float64
			var expected *metrics.ExpectedCost
			if (it.Forecast != nil && it.Forecast.Duration != "") || it.Estimate != "" {
				rate = metrics.CostPerAgentHour(items)
				expected = metrics.ExpectedCostOf(it, rate)
			}
			if a.jsonOut {
				out := map[string]any{"item": it, "children": children}
				if expected != nil {
					out["expected_cost"] = expected
				}
				if storyTopics != nil {
					out["topics"] = storyTopics
				}
				if plan != nil {
					out["plan"] = plan
				}
				return a.printJSON(out)
			}
			fmt.Fprintf(a.out, "%s %s\n", it.ID, it.Title)
			fmt.Fprintf(a.out, "  %s · %s · %s", it.Type, it.Nature, it.Status)
			if it.IsBlocked() {
				fmt.Fprint(a.out, " · BLOCKED")
			}
			if it.Archived {
				fmt.Fprint(a.out, " · archived")
			}
			fmt.Fprintln(a.out)
			if it.Parent != "" {
				fmt.Fprintf(a.out, "  parent: %s\n", it.Parent)
			}
			fmt.Fprintf(a.out, "  owner: %s · created %s · updated %s\n", it.Owner, it.Created, it.Updated)
			if len(it.Tags) > 0 {
				fmt.Fprintf(a.out, "  tags: %s\n", strings.Join(it.Tags, ", "))
			}
			if len(it.After) > 0 {
				fmt.Fprintf(a.out, "  after: %s\n", strings.Join(it.After, ", "))
			}
			switch {
			case storyTopics != nil:
				fmt.Fprintln(a.out, "  topics:")
				width := 0
				for _, t := range storyTopics {
					width = max(width, len(t.Topic))
				}
				for _, t := range storyTopics {
					fmt.Fprintf(a.out, "    %-*s  %s\n", width, t.Topic, t.Summary())
				}
			case len(it.Topics) > 0:
				fmt.Fprintf(a.out, "  topics: %s\n", strings.Join(it.Topics, ", "))
			}
			if !it.Agent.IsZero() {
				fmt.Fprintf(a.out, "  agent: %s\n", it.Agent)
			}
			for _, l := range planningLines(it.Draft, it.Finalized, it.CostOfDelay, it.Forecast, repo.Manifest.Planning.CurrencyCode()) {
				fmt.Fprintf(a.out, "  %s\n", l)
			}
			if expected != nil {
				fmt.Fprintf(a.out, "  %s\n", expectedCostLine(it, expected, *rate))
			}
			if !it.Usage.Empty() {
				fmt.Fprintf(a.out, "  usage: %s\n", it.Usage.Summary())
				for _, m := range it.Usage.Models {
					fmt.Fprintf(a.out, "    %s\n", m)
				}
			}
			// strategic usage stands apart, and alone when no agent has
			// worked the item yet (ADR-0083)
			if it.Usage != nil {
				for _, s := range it.Usage.Strategic {
					fmt.Fprintf(a.out, "  %s\n", strategicLine(s))
					for _, m := range s.Models {
						fmt.Fprintf(a.out, "    %s\n", m)
					}
				}
			}
			fmt.Fprintf(a.out, "  file: %s\n", relPath(repo.Root, it.Path))
			if len(it.Transitions) > 0 {
				fmt.Fprintln(a.out, "  history:")
				for _, tr := range it.Transitions {
					fmt.Fprintf(a.out, "    %s  %-12s by %s\n", tr.At, tr.To, tr.By)
				}
			}
			for _, b := range it.Blocked {
				until := orDefault(b.Until, "open")
				fmt.Fprintf(a.out, "  blocked: %s to %s: %s\n", b.From, until, b.Reason)
			}
			if len(children) > 0 {
				fmt.Fprintln(a.out, "  children:")
				for _, c := range children {
					fmt.Fprintf(a.out, "    %s  %-12s %s\n", c.ID, c.Status, c.Title)
				}
			}
			if plan != nil {
				a.printTaskPlan(plan)
			}
			return nil
		},
	}
}

// expectedCostLine is flai show's line for an item's expected cost: the
// cost, that it is an estimate, and the duration and rate it comes from.
func expectedCostLine(it *workitem.Item, e *metrics.ExpectedCost, rate float64) string {
	d := it.Estimate
	if e.From == "forecast" {
		d = it.Forecast.Duration
	}
	return fmt.Sprintf("expected cost: $%.2f (estimated, from the %s of %s at $%.2f per agent hour)", e.Cost, e.From, d, rate)
}

// strategicLine is flai show's line for what one kind of strategic agent
// spent on an item, kept apart from what its agents spent.
func strategicLine(s usage.Strategic) string {
	cost := fmt.Sprintf("$%.2f", s.Cost())
	if s.Estimated {
		cost += " (estimated)"
	}
	return fmt.Sprintf("%s, strategic: %s tokens · %s · %s, apart from the agents' usage", s.Kind, usage.Count(s.Tokens()), cost, (time.Duration(s.Seconds) * time.Second).String())
}

// printTaskPlan prints a story's task plan as flai show gives it: each task's
// state and what it waits for, then the layers, numbered from 1.
func (a *app) printTaskPlan(p *workitem.Plan) {
	fmt.Fprintln(a.out, "  plan:")
	inLayer := map[string]bool{}
	for _, l := range p.Layers {
		for _, id := range l {
			inLayer[id] = true
		}
	}
	var outside []string
	for _, t := range p.Tasks {
		note := ""
		switch {
		case len(t.WaitingFor) > 0:
			note = "  for " + strings.Join(t.WaitingFor, ", ")
		case len(t.After) > 0:
			note = "  after " + strings.Join(t.After, ", ")
		}
		fmt.Fprintln(a.out, strings.TrimRight(fmt.Sprintf("    %s  %-11s%s", t.ID, t.State, note), " "))
		if t.State != workitem.PlanCancelled && !inLayer[t.ID] {
			outside = append(outside, t.ID)
		}
	}
	fmt.Fprintln(a.out, "  layers:")
	for i, l := range p.Layers {
		fmt.Fprintf(a.out, "    %d  %s\n", i+1, strings.Join(l, ", "))
	}
	if len(outside) > 0 {
		fmt.Fprintf(a.out, "    none  %s: after forms a cycle, which flai check names\n", strings.Join(outside, ", "))
	}
}

// afterHelp is what flai story new and flai task new say of --after.
func afterHelp(typ string) string {
	switch typ {
	case workitem.Story:
		return `

--after names the stories this story waits for: while any of them is not
done, it is held in ready (ADR-0046). It is checked as --body-stdin is: a
story that does not exist, or a cycle, refuses the creation.`
	case workitem.Task:
		return `

--after names the tasks of the same story this task waits for: the story's
agent starts it when they are done, and runs together the tasks that wait
for nothing undone (S-0176). It is checked as --body-stdin is: a task that
does not exist, a task of another story, or a cycle refuses the creation.`
	}
	return ""
}

// draftHelp is what flai story new says of --draft.
func draftHelp(typ string) string {
	if typ != workitem.Story {
		return ""
	}
	return `

--draft makes the story a draft, as the stories an agent writes are (S-0199):
moving it to ready is refused until it is finalized, by flai edit with
--no-draft, or by flai move to ready with --yes.`
}

// costHelp is what flai story new and flai epic new say of the cost of
// delay inputs.
func costHelp(typ string) string {
	if typ == workitem.Task {
		return ""
	}
	return `

--revenue-per-week and --penalty-per-week, amounts in planning.currency, and
--time-lost-per-cycle, a Go duration, are the inputs of the ` + typ + `'s cost
of delay (S-0204), recorded with its owner and the time it is made, for the
planner to turn into a value. Without them the ` + typ + ` has no cost of
delay; flai edit sets and changes them later.`
}

// newCostOfDelay is a new item's cost of delay from its input flags, the
// inputs set by by at now (ADR-0080), or nil when none is given (S-0204). An
// amount that is not a number is refused as flai edit refuses it; Create
// checks the rest.
func newCostOfDelay(revenue, penalty, timeLost, currency, by string, now time.Time) (*workitem.CostOfDelay, error) {
	var in workitem.CostInputs
	for _, a := range []struct {
		key, value string
		to         **float64
	}{
		{"cost_of_delay.inputs.revenue_per_week", revenue, &in.RevenuePerWeek},
		{"cost_of_delay.inputs.penalty_per_week", penalty, &in.PenaltyPerWeek},
	} {
		v := strings.TrimSpace(a.value)
		if v == "" {
			continue
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("%s %q is not a number: write an amount in %s such as 1200 or 99.5, or leave it out", a.key, a.value, currency)
		}
		*a.to = &f
	}
	in.TimeLostPerCycle = strings.TrimSpace(timeLost)
	if in.IsZero() {
		return nil, nil
	}
	in.By, in.At = by, now.UTC().Format(workitem.TimeFormat)
	return &workitem.CostOfDelay{Inputs: &in}, nil
}

// pluralType is an item type's plural: epics, stories, tasks.
func pluralType(typ string) string {
	if typ == workitem.Story {
		return "stories"
	}
	return typ + "s"
}

// withArticle is an item type with its indefinite article: an epic, a story.
func withArticle(typ string) string {
	if strings.IndexByte("aeiou", typ[0]) >= 0 {
		return "an " + typ
	}
	return "a " + typ
}

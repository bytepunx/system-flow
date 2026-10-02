package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newServeAgentUsageCmd(a *app) *cobra.Command {
	var all, write bool
	c := &cobra.Command{
		Use:   "usage [story-id...]",
		Short: "Measure the tokens and cost a story's agents spent, from the logs flai serve kept",
		Long: `Measures what the agents flai serve started for a story spent, from the
logs it keeps of them in the serve folder beside flai's config file: the
story whole, and each of its tasks (S-0143, ADR-0051). Only Claude Code's
stream-json is read. Each session's last result gives its tokens and cost
per model; what came after it, a run still going or one that died, is
counted from its calls and priced at the rate the logs report for the
model, and marked estimated. A task gets its story's totals in the share
of the input and cache tokens of its own calls: those of the sub-agents
started for it, known by the task ID their Agent call's description, or
else its prompt, names, and an even share, among the tasks in progress at
the time, of every other call made while it was in progress (ADR-0071). It
is marked estimated.

It prints what it finds. --write records it in the items' usage, as flai
serve does when an agent ends, and sums the story's epic again. --all
measures every story of the project flai serve has kept a log for, to fill
in stories worked before flai measured them. Use the flai your flai serve
runs, or give its configuration with --config, so that the logs are the
ones it keeps.`,
		Example: `  flai serve agent usage S-0142
  flai serve agent usage S-0142 --write
  flai serve agent usage --all --write --json`,
		RunE: func(_ *cobra.Command, args []string) error {
			return a.agentUsage(args, all, write)
		},
	}
	c.Flags().BoolVar(&all, "all", false, "every story flai serve has kept a log for")
	c.Flags().BoolVar(&write, "write", false, "record what the logs say in the items' usage")
	return c
}

// agentUsage measures stories from the logs flai serve kept of their agents.
func (a *app) agentUsage(args []string, all, write bool) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	dir := a.serveDir()
	stories := make([]string, 0, len(args))
	for _, id := range args {
		stories = append(stories, workitem.CanonicalID(id))
	}
	if all {
		if len(args) > 0 {
			return fmt.Errorf("rule: name stories or give --all, not both")
		}
		if stories, err = dir.LoggedStories(repo.Manifest.Key); err != nil {
			return err
		}
	}
	if len(stories) == 0 {
		return fmt.Errorf("rule: name a story, or give --all")
	}
	root := mainRootOf(repo)
	var out []*serve.Measured
	for _, id := range stories {
		m, err := serve.Measure(dir, root, repo.Manifest.Key, id, write)
		if err != nil {
			return fmt.Errorf("measuring %s: %w", id, err)
		}
		out = append(out, m)
	}
	if a.jsonOut {
		return a.printJSON(out)
	}
	for _, m := range out {
		switch {
		case len(m.Logs) == 0:
			fmt.Fprintf(a.out, "%s: flai serve kept no log of an agent for it\n", m.Story)
			continue
		case m.Usage.Empty():
			fmt.Fprintf(a.out, "%s: %d %s, nothing spent in them\n", m.Story, len(m.Logs), oneOrMany(len(m.Logs), "log", "logs"))
		default:
			fmt.Fprintf(a.out, "%s: %s (%d %s)\n", m.Story, m.Usage.Summary(), len(m.Logs), oneOrMany(len(m.Logs), "log", "logs"))
			for _, mu := range m.Usage.Models {
				fmt.Fprintf(a.out, "  %s\n", mu)
			}
		}
		for _, id := range slices.Sorted(maps.Keys(m.Tasks)) {
			if u := m.Tasks[id]; !u.Empty() {
				fmt.Fprintf(a.out, "  %s: %s\n", id, u.Summary())
			}
		}
		if write && len(m.Changed) > 0 {
			fmt.Fprintf(a.out, "  written: %s\n", strings.Join(m.Changed, ", "))
		}
	}
	return nil
}

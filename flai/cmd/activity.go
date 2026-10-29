package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// flai activity: a strategic agent's activities, in its activity document
// under wip/agents (S-0206, ADR-0079).
func newActivityCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "activity",
		Short: "Log a strategic agent's activity in its activity document under wip/agents",
		Long: `A strategic agent, the planner, the orchestrator, the analyzer, the curator,
or a role the project defines under its design folder's workflows/roles
(ADR-0207), keeps an activity document, wip/agents/<kind>.md in the main
checkout, with one entry per activity it finishes and their totals in its
front matter (S-0206, ADR-0079). flai serve logs a run's activity by itself
when the run ends.`,
	}
	c.AddCommand(newActivityLogCmd(a))
	return c
}

func newActivityLogCmd(a *app) *cobra.Command {
	var items []string
	c := &cobra.Command{
		Use:   "log <kind> <summary>",
		Short: "Log an activity of a strategic agent, measured from its run's log and charged as the MCP tool activity_log charges it",
		Long: `Logs an activity of the strategic agent kind, the planner, the orchestrator,
the analyzer, the curator, or a role the project defines (its name), in its
activity document, wip/agents/<kind>.md in the main checkout, as the MCP
tool activity_log does (S-0417): flai measures its wall-clock seconds and
cost from the kind's run log, from the later of the run's start and the
last entry to now, appends the entry with the summary and each --item, and
adds it to the document's totals. The cost is charged under
usage.strategic as activity_log charges it; what no item takes is left in
the kind's project strategic total. An activity outside a run flai serve
logged is logged with no seconds and no cost. A kind that is none of these
is refused, naming those there are. Nothing is committed.`,
		Example: `  flai activity log reviewer "reviewed S-0042 against its criteria" --item S-0042
  flai activity log planner "planned E-0016" --item E-0016 --json`,
		Args: cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			summary := strings.Join(strings.Fields(args[1]), " ")
			if summary == "" {
				return errors.New("an activity needs a summary: say in one line what you did")
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			ids := make([]string, 0, len(items))
			for _, it := range items {
				ids = append(ids, workitem.CanonicalID(strings.TrimSpace(it)))
			}
			logged, err := a.mcpActivity(context.Background(), mainRootOf(repo), strings.TrimSpace(args[0]), summary, ids, a.writer())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(logged)
			}
			e, doc := logged.Entry, logged.Activity
			fmt.Fprintf(a.out, "logged in wip/agents/%s.md at %s: %s (%ds, $%.2f)\n", doc.Kind, e.At, e.Summary, e.Seconds, e.Cost)
			if logged.Charge != "" {
				fmt.Fprintf(a.out, "  %s\n", logged.Charge)
			}
			fmt.Fprintf(a.out, "  %d activities, %ds, $%.2f in all\n", doc.TasksCompleted, doc.AccruedSeconds, doc.AccruedCost)
			return nil
		},
	}
	c.Flags().StringArrayVar(&items, "item", nil, "the ID of an item the activity touched, such as S-0042; give it once for each")
	return c
}

// mcpActivity logs a strategic agent's activity for flai mcp's tool
// activity_log (S-0206, ADR-0079), and for flai activity log (S-0417):
// measured from the kind's run logs in the
// serve folder, appended to wip/agents/<kind>.md in the project's main
// checkout, charged to the items or issues it concerned (ADR-0083,
// ADR-0095, S-0227), and returned with the document's totals rather than
// every entry and with what it was charged to.
func (a *app) mcpActivity(_ context.Context, root, kind, summary string, items []string, by string) (mcpserver.ActivityLogged, error) {
	repo, err := workitem.Open(root)
	if err != nil {
		return mcpserver.ActivityLogged{}, err
	}
	logged, err := serve.LogActivity(a.serveDir(), mainRootOf(repo), repo.Manifest.Key, kind, summary, items, a.now())
	if err != nil && logged == nil {
		return mcpserver.ActivityLogged{}, err
	}
	if err != nil {
		// the activity is logged; only its charge to its items failed (ADR-0083, ADR-0095, S-0227)
		a.logger().Warn("activity not charged to its items", "component", "mcp", "agent", by, "project", repo.Manifest.Key, "kind", kind, "err", err)
	}
	e, doc := logged.Entry, logged.Activity
	to, charge := chargedTo(logged, err)
	a.logger().Info("activity logged", "component", "mcp", "agent", by, "project", repo.Manifest.Key, "kind", kind, "seconds", e.Seconds, "cost", e.Cost, "charged", strings.Join(logged.Charged, ","))
	return mcpserver.ActivityLogged{
		ChargedTo: to, Charge: charge,
		Entry: mcpserver.ActivityEntry{
			At: e.At.UTC().Format(workitem.TimeFormat), Summary: e.Summary, Items: e.Items,
			Seconds: e.Seconds, Cost: e.Cost, Estimated: e.Estimated,
		},
		Activity: mcpserver.ActivityTotals{
			Kind: doc.Kind, AccruedCost: doc.AccruedCost, AccruedSeconds: doc.AccruedSeconds,
			TasksCompleted: doc.TasksCompleted, LastRun: doc.LastRun,
		},
	}, nil
}

// chargedTo are the items a logged activity's cost was charged to, and a
// line saying so: the item a planner's run planned, the work items an
// orchestrator's activity named, the issues an analyzer's named (S-0227),
// or, when none took any of a cost, the kind's project strategic total
// (ADR-0095). The line is empty when the activity cost nothing; failed, the
// charge's error, when it failed.
func chargedTo(logged *serve.Logged, failed error) ([]string, string) {
	kind := logged.Activity.Kind
	var to []string
	var line string
	switch {
	case logged.Planned != "":
		to = []string{logged.Planned}
		line = fmt.Sprintf("charged to %s, the item its run planned, and the items above it", logged.Planned)
	case kind == workitem.ActivityAnalyzer && len(logged.Shared) == 1:
		to = logged.Shared
		line = fmt.Sprintf("charged to %s, the issue it named", logged.Shared[0])
	case kind == workitem.ActivityAnalyzer && len(logged.Shared) > 1:
		to = logged.Shared
		line = fmt.Sprintf("charged evenly to %s, the issues it named", strings.Join(logged.Shared, ", "))
	case len(logged.Shared) == 1:
		to = logged.Shared
		line = fmt.Sprintf("charged to %s and the items above it", logged.Shared[0])
	case len(logged.Shared) > 1:
		to = logged.Shared
		line = fmt.Sprintf("charged evenly to %s and the items above them", strings.Join(logged.Shared, ", "))
	case logged.Entry.Cost > 0:
		line = fmt.Sprintf("charged to no item: left in the %s's project strategic total", kind)
	}
	if failed != nil {
		line = strings.TrimPrefix(line+"; ", "; ") + failed.Error() + "; what was not charged is left in the project strategic total"
	}
	return to, line
}

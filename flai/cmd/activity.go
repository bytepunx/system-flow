package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// mcpActivity logs a strategic agent's activity for flai mcp's tool
// activity_log (S-0206, ADR-0079): measured from the kind's run logs in the
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

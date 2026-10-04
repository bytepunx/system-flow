package cmd

import (
	"context"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// mcpActivity logs a strategic agent's activity for flai mcp's tool
// activity_log (S-0206, ADR-0079): measured from the kind's run logs in the
// serve folder, appended to wip/agents/<kind>.md in the project's main
// checkout, and returned with the document's totals rather than every entry.
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
		// the activity is logged; only its charge to the item it planned failed (ADR-0083)
		a.logger().Warn("activity not charged to its item", "component", "mcp", "agent", by, "project", repo.Manifest.Key, "kind", kind, "err", err)
	}
	e, doc := logged.Entry, logged.Activity
	a.logger().Info("activity logged", "component", "mcp", "agent", by, "project", repo.Manifest.Key, "kind", kind, "seconds", e.Seconds, "cost", e.Cost)
	return mcpserver.ActivityLogged{
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

---
id: TH-0147
title: "S-0212 plan: five tasks in four layers, charts built in flaiover from flai stats as it is"
anchor:
  path: wip/kanban/stories/S-0212-charts-compare-forecasts-and-estimates-with-what-happened.md
  item: S-0212
status: resolved
participants: [planner-S-0212, alex]
created: 2026-10-05T05:47:13Z
updated: 2026-10-05T05:50:50Z
---

# TH-0147 S-0212 plan: five tasks in four layers, charts built in flaiover from flai stats as it is

On wip/kanban/stories/S-0212-charts-compare-forecasts-and-estimates-with-what-happened.md.

## Entries

### 2026-10-05T05:47:13Z planner-S-0212
Plan for S-0212. I have drafted five tasks, all in the backlog.

Layers:

- Layer 1: T-0950 adds the chart data types (per-item forecast errors and `forecasts`), `PLANNING_KINDS`, the nature and model filters, and the forecast-accuracy chart, in `flaiover/src/lib/viz/charts.ts` and its test.
- Layer 2: T-0951 adds the delivery-accuracy chart, with delivery error in days and the share delivered on time per ISO week. It waits for T-0950, because it uses the same types and changes the same files.
- Layer 3: T-0956 adds the forecast-by-model chart (p50 absolute error per bucket per model) and the rows for the table view. It waits for T-0951, for the same files.
- Layer 4: these two run together, since their paths are disjoint:
  - T-0958 adds the Planning group on the chart page, the nature and model selects, `ForecastTable.svelte`, and the summary strip. It waits for T-0956.
  - T-0962 updates `flaiover-dashboard.md`, the Charts table in `metrics.md`, and `docs/users/flaiover.md`. It waits for T-0956.

Assumptions:

1. No flai change is needed. `flai stats --json` already sends `forecast_error_seconds`, `delivery_error_seconds`, and `estimate_error_seconds` per item, and `forecasts` with p50 and p85 overall, by nature, and by model, limited to the window (S-0205, `flai/internal/metrics/forecast.go`).
   - The dashboard takes the percentile lines from `forecasts` when one filter is set, so they match flai to the second.
   - It works out three figures from the per-item values: the lines when both filters are set, the weekly on-time share, and the per-bucket p50 per model. None of these is in `flai stats`.
   - If you want those figures in flai, it becomes a metrics change with an ADR, and I would add a task for it.
2. "The Charts menu lists them under a Planning group" means the chart page's header groups (flow and usage today, from `FLOW_KINDS` and `USAGE_KINDS`), not the site menu. So `flaiover/src/lib/sitemenu.ts` and `flaiover/src/lib/charts` (which does not exist) are declared touches that will likely stay unchanged. I kept them because they are declared; dropping them would spare a touches-drift finding.
3. `design/system/metrics.md` gets rows in its Charts table. That is a design touch I added; it needs no ADR because the metrics are unchanged.

Figures:

- Forecast: 31m, delivery 2026-10-05T13:46Z, from flai's 14 matched stories. It falls within the 12–33 agent minutes that past chart stories took.
- Cost of delay: 33.79 USD a week. This is S-0212's share, by duration, of E-0016's 1500 USD a week, because the story has no inputs of its own. It replaces planner-E-0016's 60.98, which was worked out on a 1h15m forecast.

Proposals: none to split, merge, or drop.

### 2026-10-05T05:50:50Z alex
Resolved.

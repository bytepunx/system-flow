---
id: S-0212
type: story
nature: feature
title: Charts compare forecasts and estimates with what happened
status: review
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-07T07:31:41Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:50Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:38Z
    by: alex
  - to: ready
    at: 2026-10-07T03:29:31Z
    by: alex
  - to: in-progress
    at: 2026-10-07T06:46:47Z
    by: agent-S-0212
  - to: review
    at: 2026-10-07T07:31:41Z
    by: agent-S-0212
tags: [dashboard]
topics: [planning]
touches: [flaiover/src/routes/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, design/system/metrics.md, flaiover/src/lib/components/ForecastTable.svelte, flaiover/src/lib/components/ForecastTable.svelte.test.ts, design/adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md, design/adrs/README.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md, flai/internal/metrics/forecast.go, flai/internal/metrics/forecast_test.go, flai/internal/metrics/metrics.go, design/issues/I-0105-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md, design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2719
  models:
    - model: claude-opus-5-5
      input: 392
      output: 214625
      cache_read: 23096485
      cache_write: 754164
      cost: 13.3775
    - model: claude-sonnet-5-5
      input: 46
      output: 9960
      cache_read: 715336
      cache_write: 89485
      cost: 0.4665
cost_of_delay:
  value: 33.79
  by: planner-S-0212
  at: 2026-10-05T05:45:27Z
forecast:
  duration: 31m
  delivery: 2026-10-07T04:48:00Z
  basis: "Its own forecast of 31m; 4th in the pull order with an in-progress limit of 3, behind S-0270, S-0271, S-0274 and S-0275."
  by: flai
  at: 2026-10-07T03:30:07Z
---
# S-0212 Charts compare forecasts and estimates with what happened

## Goal

The operator should see whether the planner's forecasts can be trusted before letting the orchestrator act on them.

## Acceptance criteria
- [x] `/charts/forecast-accuracy`: one point per story completed in the window, x completed, y forecast error (actual minus forecast duration), with p50 and p85 lines of absolute error, filterable by nature and by model; estimate error as a second series when a human estimate exists
- [x] `/charts/delivery-accuracy`: delivery error in days per story over time, with the share delivered on or before the forecast date per week
- [x] `/charts/forecast-by-model`: p50 absolute forecast error per bucket per model
- [x] Every chart spans the window as ADR-0054 requires, reads `/api/stats`, and matches `flai stats --json` to the second; the Charts menu lists them under a Planning group
- [x] `design/system/flaiover-dashboard.md` and the user guide describe them; tests cover each chart's data mapping

## Tasks
- T-0950 The chart data types carry the forecast errors, and the forecast-accuracy chart maps them with nature and model filters
- T-0951 The delivery-accuracy chart maps each story's delivery error in days and the weekly share delivered on time
- T-0956 The forecast-by-model chart maps the p50 absolute forecast error per bucket per model
- T-0958 The chart page lists the planning charts under a Planning group, with nature and model filters and a table view
- T-0962 The dashboard design, the metrics chart table, and the user guide describe the planning charts
- T-1156 flai stats gives each item the model its forecasts are grouped under

## Notes

### Planning

Touches, by where each came from:

- Declared, kept: `flaiover/src/routes/charts` (the chart page `[kind]/+page.svelte`, its Planning group, filters, and `charts.svelte.test.ts`), `flaiover/src/lib/viz` (`charts.ts` builders, `Report` and `ItemMetrics` types, `controls`, and `charts.test.ts`), `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flaiover/src/lib/charts`, and `flaiover/src/lib/sitemenu.ts`. The last two are likely to stay unchanged: `flaiover/src/lib/charts` does not exist, the chart mapping lives in `flaiover/src/lib/viz/charts.ts`, and the "Charts menu" groups (flow, usage) are the chart page's header, built from `FLOW_KINDS` and `USAGE_KINDS`, not the site menu. The operator may drop them to spare touches drift.
- Design: `design/system/metrics.md`, whose Charts table lists every chart the dashboard draws; the three planning charts get rows there, mapping metrics S-0205 already defines, so no ADR.
- Layout: `flaiover/src/lib/components/ForecastTable.svelte` and its test, the table view under the planning charts, after `SpendTable.svelte`'s precedent.
- Co-change, left out: `design/system/flai-cli.md` (44%), `docs/users/flai.md` (40%), and `docs/operators/index.md` (34%) change with the dashboard docs often, but this story changes neither flai nor operator settings: `flai stats --json` already carries `forecasts` and the per-item errors, windowed (`flai/internal/metrics/forecast.go`).

Figures:

- Forecast 31m, delivery 2026-10-05T13:46Z: flai's figure with the predicted touches (132 s a unit over 14 done feature stories on claude-opus-5-5 in the large band, times size 14). Before the touches were predicted flai said 14m on only 3 medium stories, which was low against the 12 to 33 agent minutes comparable chart stories took (S-0163, S-0166, S-0168, S-0169); 31m sits in that range, so it stands. The delivery assumes the operator promotes S-0212 when its place in the pull order comes; it is in the backlog today.
- Cost of delay 33.79 USD a week: S-0212 has no inputs of its own, so `flai cod` gives its share of E-0016's 1500 USD a week (the operator's 10h time lost per cycle) by forecast duration, 31m of 22h56m over the epic's 17 open stories without inputs. It stands: the story brings no revenue alone, its value is in gating the orchestrator's trust in forecasts, which the epic's figure already prices. It replaces planner-E-0016's 60.98, worked out on a 1h15m forecast.

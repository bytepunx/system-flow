---
id: S-0213
type: story
nature: feature
title: Charts show cost of delay outstanding, incurred, and what the pull order costs
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-06T11:14:34Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:52Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:46Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs, flai/internal/statsread, flai/internal/planning/forecast.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai.md]
after: [S-0205, S-0217]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 91.22
  by: planner-S-0213
  at: 2026-10-05T05:45:55Z
forecast:
  duration: 1h30m
  delivery: 2026-10-06T17:48:00Z
  basis: "Its own forecast of 1h30m; 10th in the pull order with an in-progress limit of 3, behind S-0221, S-0222, S-0224, S-0226, S-0223, S-0227, S-0229, S-0284 and S-0212."
  by: flai
  at: 2026-10-06T11:14:34Z
---
# S-0213 Charts show cost of delay outstanding, incurred, and what the pull order costs

## Goal

Cost of delay only matters if the operator can see it accumulate and what ordering by it would save.

## Acceptance criteria
- [ ] `/charts/cod-outstanding`: stacked area per day of the window, cost of delay value per column (backlog, ready, in-progress, review)
- [ ] `/charts/cod-incurred`: bar per week of cost of delay incurred (value × time waited), with the running mean
- [ ] `/charts/cod-order`: for the current ready column, the incurred cost projected under the current pull order against ordering by cost of delay (and by WSJF, value over forecast duration), as two lines over the forecast horizon, with the difference stated
- [ ] Items without a value are shown as a count so their absence is visible
- [ ] Charts span the window, read `/api/stats`, match `flai stats --json`; the Charts menu lists them under Planning; the design and user guide describe them; tests cover the data mapping

## Tasks
- T-0953 metrics.md defines the pull order's projected cost of delay and the count of items without a value, and an ADR records them
- T-0955 flai stats counts the open items without a cost of delay value per column per day
- T-0959 flai stats projects the ready column's cost of delay under the pull order, by cost of delay, and by WSJF
- T-0965 flai stats prints the pull order's projected saving and the items without a cost of delay value
- T-0967 charts.ts maps cost_of_delay to the cod-outstanding, cod-incurred, and cod-order charts in a planning group
- T-0969 The charts page lists the cost of delay charts under Planning and states the saving and the items without a value
- T-0970 flaiover-dashboard.md and the user guide describe the cost of delay charts

## Notes

### Planning

Touches, by where each came from:

- Declared, kept: `flaiover/src/routes/charts`, `flaiover/src/lib/charts`, `flaiover/src/lib/viz`, `flaiover/src/lib/sitemenu.ts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flai/internal/metrics`, `design/system/metrics.md`, `design/adrs`. `flaiover/src/lib/charts` does not exist (the chart code is `src/lib/viz/charts.ts`) and the Charts menu's groups are in `src/routes/charts/[kind]/+page.svelte`, not `sitemenu.ts`; both are kept as declared.
- Layout: `flai/internal/statsread` (it builds `metrics.Options` for `flai stats` and the dashboard's `stats.get`, and must pass the board's order and the planning settings for the projection); `flai/internal/planning/forecast.go` (its `playOut` already lays the pull order onto the in-progress limit, which the projection replays).
- Co-change and design: `flai/cmd/stats.go` and `flai/cmd/check_stats_test.go` (it prints the cost of delay aggregates, so the projection's difference and the count without a value join them); `docs/users/flai.md` (its `flai stats` section names the cost of delay output). `design/system/flai-cli.md` and `docs/operators/index.md` lead `flai touches suggest` (41%, 27%) but describe nothing this story changes, so they are left out.

Figures:

- Forecast 1h30m, kept over `flai forecast`'s 31m (132 s × size 14): the size counts five criteria and nine touches, but the work is an ADR, a projection in Go, three chart builders, the page, and two guides, six tasks across Go and Svelte, which is S-0205's shape (seven tasks, 1h09m in progress). Delivery 2026-10-05T15:19Z is flai's for 1h30m at 11th in the pull order; it waits for S-0217.
- Cost of delay 91.22 USD a week, from `flai cod`: the story has no inputs of its own, so it is its share of E-0016's 1500 USD a week by forecast, 1h30m of 24h40m over 17 open stories. It replaces the 73.17 planner-E-0016 set, which was the share when the epic had other forecasts.

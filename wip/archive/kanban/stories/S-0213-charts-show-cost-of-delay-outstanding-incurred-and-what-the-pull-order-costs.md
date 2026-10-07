---
id: S-0213
type: story
nature: feature
title: Charts show cost of delay outstanding, incurred, and what the pull order costs
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-07T14:57:53Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:52Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:46Z
    by: alex
  - to: ready
    at: 2026-10-07T03:29:29Z
    by: alex
  - to: in-progress
    at: 2026-10-07T07:32:18Z
    by: agent-S-0213
  - to: review
    at: 2026-10-07T08:17:48Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T08:18:03Z
    by: agent-S-0213
  - to: review
    at: 2026-10-07T08:42:19Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T14:57:53Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs, flai/internal/statsread, flai/internal/planning/forecast.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/summary.md]
after: [S-0205, S-0217]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5358
  models:
    - model: claude-opus-5-5
      input: 730
      output: 253266
      cache_read: 46810702
      cache_write: 1077606
      cost: 21.1967
    - model: claude-sonnet-5-5
      input: 48
      output: 14212
      cache_read: 785691
      cache_write: 152282
      cost: 0.6801
  strategic:
    - kind: orchestrator
      seconds: 965
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 75
          output: 992
          cache_read: 5575774
          cache_write: 32505
          cost: 1.4633
        - model: claude-sonnet-5-5
          input: 10
          output: 43
          cache_read: 136303
          cache_write: 59643
          cost: 0.1314
cost_of_delay:
  value: 91.22
  by: planner-S-0213
  at: 2026-10-05T05:45:55Z
forecast:
  duration: 1h30m
  delivery: 2026-10-07T09:09:00Z
  basis: "Its own forecast of 1h30m; 3rd in the pull order with an in-progress limit of 3, behind S-0212, S-0271, S-0274 and S-0275."
  by: flai
  at: 2026-10-07T06:48:44Z
---
# S-0213 Charts show cost of delay outstanding, incurred, and what the pull order costs

## Goal

Cost of delay only matters if the operator can see it accumulate and what ordering by it would save.

## Acceptance criteria
- [x] `/charts/cod-outstanding`: stacked area per day of the window, cost of delay value per column (backlog, ready, in-progress, review)
- [x] `/charts/cod-incurred`: bar per week of cost of delay incurred (value × time waited), with the running mean
- [x] `/charts/cod-order`: for the current ready column, the incurred cost projected under the current pull order against ordering by cost of delay (and by WSJF, value over forecast duration), as two lines over the forecast horizon, with the difference stated
- [x] Items without a value are shown as a count so their absence is visible
- [x] Charts span the window, read `/api/stats`, match `flai stats --json`; the Charts menu lists them under Planning; the design and user guide describe them; tests cover the data mapping

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

### Delivered

- Criterion 3 draws three lines, not two: the pull order, by cost of delay, and by WSJF, with the saving stated against the cheaper of the two policies ([ADR-0112](../../../design/adrs/0112-flai-stats-counts-the-open-items-without-a-cost-of-delay-value-per-column-and.md)).
- CoD Outstanding and CoD Incurred offer the item type, as they sum the report's type; CoD by Order projects the ready stories whatever the type.
- 2026-10-07T08:18:03Z: moved to in-progress: S-0212 was accepted as S-0213 moved to review; S-0213 takes the eight conflicting paths at its sync (TH-0250)

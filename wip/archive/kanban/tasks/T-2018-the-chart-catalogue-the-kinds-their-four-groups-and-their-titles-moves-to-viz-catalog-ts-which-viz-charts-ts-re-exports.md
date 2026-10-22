---
id: T-2018
type: task
nature: improvement
title: The chart catalogue, the kinds, their four groups, and their titles, moves to viz/catalog.ts, which viz/charts.ts re-exports
status: done
parent: S-0484
owner: alex
created: 2026-10-21T20:57:08Z
updated: 2026-10-22T20:23:10Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:29Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:19:34Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:23:10Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: [flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts]
usage:
  source: log
  seconds: 216
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 5649
      cache_read: 884176
      cache_write: 31700
      cost: 0.4829
---
# T-2018 The chart catalogue, the kinds, their four groups, and their titles, moves to viz/catalog.ts, which viz/charts.ts re-exports

## Work

Make `flaiover/src/lib/viz/catalog.ts`, a module with no import of ECharts or of the builders, and move into it from `viz/charts.ts` the kind lists (`FLOW_KINDS`, `USAGE_KINDS`, `FORECAST_KINDS`, `CLAIMS_KINDS`, `PLANNING_KINDS`, `STRATEGIC_KINDS`, `COD_KINDS`, `KINDS`), the `Kind` type, and `TITLES`. Add `CHART_GROUPS`, the four groups the charts page builds today in its local `groups` const (`flow`, `usage`, `planning` as `PLANNING_KINDS` then `COD_KINDS`, `strategic`), each a key, a label, and its kinds in order. `viz/charts.ts` re-exports everything it moved, so the page, `sitemenu.ts` later, and every test import as before.

First layer: waits for nothing. Runs beside the nav rail state task, with which it shares no path.

## Done when

`catalog.test.ts` asserts that `CHART_GROUPS` covers every kind of `KINDS` once and in `KINDS`' order, that every kind has a title, and that `catalog.ts` imports nothing from `echarts` or `viz/charts.ts`. `viz/charts.test.ts` passes unchanged, and `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. The rail is in every page's layout, so the catalogue must not drag the builders and ECharts types into that bundle.

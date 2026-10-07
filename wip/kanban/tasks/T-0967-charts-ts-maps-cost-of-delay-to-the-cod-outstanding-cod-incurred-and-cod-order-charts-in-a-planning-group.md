---
id: T-0967
type: task
nature: feature
title: charts.ts maps cost_of_delay to the cod-outstanding, cod-incurred, and cod-order charts in a planning group
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:47:14Z
updated: 2026-10-07T07:59:44Z
transitions:
  - to: ready
    at: 2026-10-07T07:51:10Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T07:51:10Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T07:59:44Z
    by: agent-S-0213
stream: S-0213
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/palette.ts]
after: [T-0959]
usage:
  source: log
  seconds: 514
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 70
      output: 393
      cache_read: 3583009
      cache_write: 131417
      cost: 1.6663
---
# T-0967 charts.ts maps cost_of_delay to the cod-outstanding, cod-incurred, and cod-order charts in a planning group

## Work

In `flaiover/src/lib/viz/charts.ts`:

- Give `Report` the `cost_of_delay` section as `flai stats --json` sends it: `days[]` with `date`, `outstanding`, `incurred`, and `without_value`; `weeks[]`; and `order`. Absent when the host's flai is older.
- Add the kinds `cod-outstanding`, `cod-incurred`, and `cod-order` to a planning list beside `FLOW_KINDS` and `USAGE_KINDS`, joining S-0212's list if it has landed, and to `KINDS`, with their titles and the controls they use: window, no type, no epic.
- Build each chart:
  - `cod-outstanding`: a stacked area per day, one series per column in board order, over the window as `span` lays it out (ADR-0054).
  - `cod-incurred`: a bar per week of `weeks[].incurred`, with the running mean per week as a dashed line, as the spend charts draw theirs.
  - `cod-order`: one line per order (`current`, `cod`, `wsjf`) of cumulative projected cost from now to the horizon, with the saving and its order returned for the page to state. Its axis runs from now to the horizon, not over the window: it is a projection.
- Return beside each chart the count of items without a value: today's per column for the first two, and `left_out` for the third.

Colours come from `palette.ts`: a column's colour as the cumulative flow diagram uses it.

It waits for T-0959, so that its fixtures are what `flai stats --json` prints. It runs with the `flai stats` printing task, whose paths it does not share.

## Done when

- `charts.test.ts` pins each chart's series from a fixture of `flai stats --json`: values per day and per week, the running mean, the three order lines, the saving, and the counts without a value
- a report with no `cost_of_delay` gives empty charts, not an error
- `npm run test` and `npm run check` in `flaiover` pass for these files

## Notes

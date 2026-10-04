---
id: S-0213
type: story
nature: feature
title: Charts show cost of delay outstanding, incurred, and what the pull order costs
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-04T21:41:59Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:52Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:46Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs]
after: [S-0205, S-0217]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 73.17
  by: planner-E-0016
  at: 2026-10-04T04:48:14Z
forecast:
  duration: 1h30m
  delivery: 2026-10-05T00:00:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 3, behind S-0212."
  by: flai
  at: 2026-10-04T21:41:59Z
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

## Notes

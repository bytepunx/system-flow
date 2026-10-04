---
id: S-0213
type: story
nature: feature
title: Charts show cost of delay outstanding, incurred, and what the pull order costs
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-03T20:33:52Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:52Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

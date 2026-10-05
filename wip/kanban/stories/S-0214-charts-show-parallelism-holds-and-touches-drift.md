---
id: S-0214
type: story
nature: feature
title: Charts show parallelism, holds, and touches drift
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-05T05:23:10Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:54Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:39Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs, docs/users/flai.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 85.37
  by: planner-E-0016
  at: 2026-10-04T04:48:15Z
forecast:
  duration: 1h45m
  delivery: 2026-10-05T16:06:00Z
  basis: "Its own forecast of 1h45m; 12th in the pull order with an in-progress limit of 3, behind S-0258, S-0276, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0226, S-0212 and S-0213."
  by: flai
  at: 2026-10-05T05:23:10Z
---
# S-0214 Charts show parallelism, holds, and touches drift

## Goal

The planner's touches exist to raise parallelism safely. The operator should see how many stories run at once, how long claims hold work, and how far declared touches are from what was changed.

## Acceptance criteria
- [ ] `/charts/parallelism`: stories in progress per day against the in-progress limit, with held stories as a second series
- [ ] `/charts/hold-time`: per week the hours stories spent held, by reason (overlap, after, empty claim)
- [ ] `/charts/touches-drift`: per story completed, files changed outside its touches and touches never changed, as stacked bars, with the share of stories whose touches were exact per week
- [ ] Charts span the window, read `/api/stats`, match `flai stats --json`; listed under Planning; design and user guide describe them; tests cover the data mapping

## Tasks

## Notes

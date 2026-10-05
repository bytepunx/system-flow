---
id: S-0212
type: story
nature: feature
title: Charts compare forecasts and estimates with what happened
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-05T00:03:14Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:50Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:38Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:48:14Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T07:01:00Z
  basis: "Its own forecast of 1h15m; 8th in the pull order with an in-progress limit of 3, behind S-0252, S-0253, S-0257, S-0244, S-0262, S-0258, S-0260 and S-0259."
  by: flai
  at: 2026-10-05T00:03:14Z
---
# S-0212 Charts compare forecasts and estimates with what happened

## Goal

The operator should see whether the planner's forecasts can be trusted before letting the orchestrator act on them.

## Acceptance criteria
- [ ] `/charts/forecast-accuracy`: one point per story completed in the window, x completed, y forecast error (actual minus forecast duration), with p50 and p85 lines of absolute error, filterable by nature and by model; estimate error as a second series when a human estimate exists
- [ ] `/charts/delivery-accuracy`: delivery error in days per story over time, with the share delivered on or before the forecast date per week
- [ ] `/charts/forecast-by-model`: p50 absolute forecast error per bucket per model
- [ ] Every chart spans the window as ADR-0054 requires, reads `/api/stats`, and matches `flai stats --json` to the second; the Charts menu lists them under a Planning group
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe them; tests cover each chart's data mapping

## Tasks

## Notes

---
id: S-0212
type: story
nature: feature
title: Charts compare forecasts and estimates with what happened
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:14Z
updated: 2026-10-02T11:54:41Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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

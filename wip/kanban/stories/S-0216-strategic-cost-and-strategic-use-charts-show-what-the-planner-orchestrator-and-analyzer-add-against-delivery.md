---
id: S-0216
type: story
nature: feature
title: Strategic Cost and Strategic Use charts show what the planner, orchestrator, and analyzer add against delivery
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-05T00:14:35Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:57Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:16Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts]
after: [S-0205, S-0225, S-0226, S-0227]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 48.78
  by: planner-E-0016
  at: 2026-10-04T04:48:16Z
forecast:
  duration: 1h
  delivery: 2026-10-05T07:43:00Z
  basis: "Its own forecast of 1h; 15th in the pull order with an in-progress limit of 3, behind S-0252, S-0259, S-0249, S-0266, S-0253, S-0257, S-0244, S-0262, S-0267, S-0268, S-0258, S-0260, S-0212, S-0213, S-0214 and S-0215."
  by: flai
  at: 2026-10-05T00:14:35Z
---
# S-0216 Strategic Cost and Strategic Use charts show what the planner, orchestrator, and analyzer add against delivery

## Goal

The operator should be able to tell whether the strategic agents pay for themselves. Two charts contrast what they cost and the time they take with the average cost and delivery time per story over the same period.

## Acceptance criteria
- [ ] `/charts/strategic-cost`: per bucket, the cost of the planner, orchestrator, and analyzer (stacked, from their activity logs) as bars, with the mean cost per story completed in the bucket as a line on a second axis, and the ratio stated
- [ ] `/charts/strategic-use`: per bucket, the agent seconds of the three as bars, with the mean cycle time and the mean agent waiting time per story completed as lines, so a fall in waiting or cycle time can be read against the agents' time
- [ ] Both span the window, read `/api/stats`, match `flai stats --json`; listed under a Strategic group; a note explains how to read them; design and user guide describe them; tests cover the mapping

## Tasks

## Notes

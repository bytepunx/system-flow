---
id: S-0215
type: story
nature: feature
title: A chart shows how long agents spend waiting on threads and on review
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-04T23:13:46Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:56Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:41Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 48.78
  by: planner-E-0016
  at: 2026-10-04T04:48:15Z
forecast:
  duration: 1h
  delivery: 2026-10-05T05:59:00Z
  basis: "Its own forecast of 1h; 9th in the pull order with an in-progress limit of 3, behind S-0248, S-0252, S-0253, S-0257, S-0244, S-0212, S-0213 and S-0214."
  by: flai
  at: 2026-10-04T23:13:46Z
---
# S-0215 A chart shows how long agents spend waiting on threads and on review

## Goal

Agents wait for the operator: on threads they open and on stories in review. That time is the operator's bottleneck and the orchestrator's reason to exist. The designer asked for a chart of it on 2026-10-02.

## Acceptance criteria
- [ ] `/charts/agent-waiting`: stacked bar per week of the window of hours agents waited, split into thread waits (opened to answered, from thread timestamps while the story was in progress) and review waits (review to done, from transitions), with the mean wait per story as a line
- [ ] A table under it lists the longest waits in the window with the story, the thread, and who was awaited
- [ ] Spans the window, reads `/api/stats`, matches `flai stats --json`; listed under Flow; design and user guide describe it; tests cover the mapping

## Tasks

## Notes

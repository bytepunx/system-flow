---
id: S-0205
type: story
nature: feature
title: flai stats and metrics.md gain the planning, waiting, and strategic-agent metrics
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-02T11:54:39Z
transitions: []
tags: [flai]
touches: [flai/internal/metrics, flai/cmd/stats.go, design/system/metrics.md, flai/internal/usage]
after: [S-0199, S-0206]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0205 flai stats and metrics.md gain the planning, waiting, and strategic-agent metrics

## Goal

The charts E-0016 asks for need numbers `flai stats` does not compute: forecast against actual, cost of delay outstanding and incurred, time agents spend waiting, claim and touches drift, and what the planner, orchestrator, and analyzer cost against delivery. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, so this is an ADR.

## Acceptance criteria
- [ ] Per story: `forecast_error` (actual cycle time minus `forecast.duration`), `delivery_error` (completed minus `forecast.delivery`), `estimate_error` (against the human `estimate`), each absent when the input is; aggregates p50 and p85 of absolute error, by nature and by the story's agent model
- [ ] Cost of delay: per item its `value`; per day of the window the value outstanding per column (summing open items' values) and the value incurred (value × days waited in backlog and ready); per week the total incurred
- [ ] Waiting: per story the time its agent waited, from each thread's opened-to-answered interval while the story was in progress, plus the review interval (review to done); aggregates mean and total per week, split into `threads` and `review`
- [ ] Claims and touches: per story the time it was held (from the launcher's hold reasons where recorded), stories in progress per day against the WIP limit, and touches drift: files committed outside its declared touches and declared touches never changed, as counts and paths
- [ ] Strategic agents: per day of the window the cost and agent seconds of the planner, orchestrator, and analyzer (from their activity logs' front matter), beside the mean cost and cycle time per story completed that day, so the dashboard can chart Strategic Cost and Strategic Use
- [ ] `flai stats --json` carries all of it; `flai stats` prints the aggregates; precision rules are written for each in `metrics.md`, and an ADR records the additions
- [ ] Tests pin each metric on fixtures, to the second

## Tasks

## Notes

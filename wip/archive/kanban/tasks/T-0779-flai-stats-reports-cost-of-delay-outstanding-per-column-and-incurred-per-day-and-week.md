---
id: T-0779
type: task
nature: feature
title: flai stats reports cost of delay outstanding per column and incurred per day and week
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:03Z
updated: 2026-10-03T20:52:51Z
transitions:
  - to: ready
    at: 2026-10-03T20:47:25Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T20:47:25Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T20:52:51Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/internal/metrics]
after: [T-0778]
usage:
  source: log
  seconds: 326
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 10478
      cache_read: 1512938
      cache_write: 40956
      cost: 0.7576
---
# T-0779 flai stats reports cost of delay outstanding per column and incurred per day and week

## Work

In `flai/internal/metrics`, per item `cost_of_delay` and `cost_of_delay_incurred`, and `cost_of_delay.days` (outstanding per column, incurred) and `cost_of_delay.weeks` (incurred), as `metrics.md` § Cost of delay defines them. Waits for the forecast task: both change `metrics.go`.

## Done when

- [ ] `Compute` reports each value
- [ ] Tests pin them on a fixture of their own, to the cent and the second

## Notes

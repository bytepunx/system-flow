---
id: T-026
type: task
nature: feature
title: "Stats: per-item metrics and aggregates"
status: done
parent: S-008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:34Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:34Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:34Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:34Z
    by: agent
stream: S-008
tags: [cli, stats]
---

# T-026 Stats: per-item metrics and aggregates

## Work
internal/metrics: per item derived timestamps (committed, started, completed), state intervals, blocked total, lead, cycle, queue, time in state, estimate error, per metrics.md. Aggregates over a window: throughput per week, cycle/lead/queue p50, p85, max, mean, WIP now, flow efficiency, cancellation rate, time-in-state share, aging WIP against p85, grouping by nature, type, or parent.

## Done when
Every definition in metrics.md has a function and a unit test.

## Notes

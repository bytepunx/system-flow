---
id: T-0781
type: task
nature: feature
title: flai stats reports held time, stories in progress against the limit, and touches drift
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:05Z
updated: 2026-10-03T21:04:17Z
transitions:
  - to: ready
    at: 2026-10-03T20:58:33Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T20:58:33Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T21:04:17Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/internal/metrics]
after: [T-0780]
usage:
  source: log
  seconds: 344
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 15949
      cache_read: 2302883
      cache_write: 62340
      cost: 1.1532
---
# T-0781 flai stats reports held time, stories in progress against the limit, and touches drift

## Work

In `flai/internal/metrics`, per item `held_seconds` replayed from the hold rules, `claims.days` (stories in progress against the limit), and `claims.drift` from the committed files `Options` carries, as `metrics.md` § Claims and touches defines them. Waits for the waiting task: both change `metrics.go`.

## Done when

- [ ] `Compute` reports each value
- [ ] Tests pin them on a fixture of their own, to the second

## Notes

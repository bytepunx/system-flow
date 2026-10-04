---
id: T-0782
type: task
nature: feature
title: flai stats reports the strategic agents' cost and seconds per day beside delivery
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:05Z
updated: 2026-10-03T21:10:02Z
transitions:
  - to: ready
    at: 2026-10-03T21:04:17Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T21:04:17Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T21:10:02Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/internal/metrics]
after: [T-0781]
usage:
  source: log
  seconds: 345
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 40
      output: 16826
      cache_read: 2429573
      cache_write: 65769
      cost: 1.2166
---
# T-0782 flai stats reports the strategic agents' cost and seconds per day beside delivery

## Work

In `flai/internal/metrics`, `strategic_days`: per day of the window each strategic agent's cost and seconds from its activity log entries, beside the cost and cycle time per item completed that day, as `metrics.md` § Strategic use per day defines them. Waits for the claims task: both change `metrics.go`.

## Done when

- [ ] `Compute` reports `strategic_days`
- [ ] Tests pin it on a fixture of their own, to the second

## Notes

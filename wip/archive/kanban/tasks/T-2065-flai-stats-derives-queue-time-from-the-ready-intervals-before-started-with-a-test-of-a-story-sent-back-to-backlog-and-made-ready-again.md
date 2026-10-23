---
id: T-2065
type: task
nature: feature
title: flai stats derives queue time from the ready intervals before started, with a test of a story sent back to backlog and made ready again
status: done
parent: S-0469
owner: alex
created: 2026-10-23T06:05:41Z
updated: 2026-10-23T06:08:05Z
transitions:
  - to: ready
    at: 2026-10-23T06:06:54Z
    by: agent-S-0469
  - to: in-progress
    at: 2026-10-23T06:06:55Z
    by: agent-S-0469
  - to: done
    at: 2026-10-23T06:08:05Z
    by: agent-S-0469
stream: S-0469
tags: []
touches: [flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go]
after: [T-2064]
usage:
  source: log
  seconds: 70
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 4654
      cache_read: 1305287
      cache_write: 35658
      cost: 0.6395
---
# T-2065 flai stats derives queue time from the ready intervals before started, with a test of a story sent back to backlog and made ready again

## Work

Change `Derive` in `flai/internal/metrics/metrics.go` so `queue_time_seconds` sums the item's `ready` intervals that end at or before `started`, as T-2064's ADR defines it. It waits for T-2064 so the code follows the recorded definition.

Add a test in `flai/internal/metrics/metrics_test.go` that reproduces I-0153: a story made ready, sent back to backlog for three days, made ready again, then started. Its queue time is the time in ready only, and fails against `started - committed`.

## Done when

- The new test passes, and fails against the old formula.
- `flai test` on the two files passes.

## Notes

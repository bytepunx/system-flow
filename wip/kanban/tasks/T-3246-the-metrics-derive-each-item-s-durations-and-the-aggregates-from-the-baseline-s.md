---
id: T-3246
type: task
nature: feature
title: The metrics derive each item's durations and the aggregates from the baseline's timers, with flai stats byte for byte as before
status: done
parent: S-0750
owner: alex
created: 2026-11-12T16:09:20Z
updated: 2026-11-15T21:11:30Z
transitions:
  - to: ready
    at: 2026-11-15T20:56:37Z
    by: agent-S-0750
  - to: in-progress
    at: 2026-11-15T20:56:37Z
    by: agent-S-0750
  - to: done
    at: 2026-11-15T21:11:30Z
    by: agent-S-0750
stream: S-0750
tags: [cli]
touches: [flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, flai/internal/metrics/waiting.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/concurrency.go, flai/internal/metrics/stats_golden_test.go, flai/internal/metrics/testdata/good.stats.golden.json, flai/internal/metrics/reopen_test.go]
after: [T-3241, T-3237]
usage:
  source: log
  seconds: 893
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 17623
      cache_read: 4570932
      cache_write: 145287
      cost: 2.1174
---
# T-3246 The metrics derive each item's durations and the aggregates from the baseline's timers, with flai stats byte for byte as before

## Work

Waits for T-3241: it calls the timer evaluation that task adds. Make `Derive` and `Compute` in `metrics.go` take time in state, queue, lead, cycle, and review time from the baseline's timers instead of their own code, and name no state. Add a test that runs `flai stats --json` on `testdata/good` and compares it byte for byte with the output before the change, and keep `windows.golden.json` as it is.

## Done when

`metrics.go` names no state, and `flai stats --json` gives the same bytes on both fixtures.

## Notes

---
id: T-1664
type: task
nature: feature
title: flai stats reports the queue per week from its log, and each item's queue wait
status: done
parent: S-0404
owner: alex
created: 2026-10-16T05:02:28Z
updated: 2026-10-20T00:42:44Z
transitions:
  - to: ready
    at: 2026-10-19T23:52:40Z
    by: agent-S-0404
  - to: in-progress
    at: 2026-10-19T23:52:41Z
    by: agent-S-0404
  - to: done
    at: 2026-10-20T00:17:06Z
    by: agent-S-0404
stream: S-0404
tags: [cli]
touches: [flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats_test.go, flai/cmd/accept_queue_test.go, docs/users/flai-reference.md]
after: [T-1651, T-1652]
usage:
  source: log
  seconds: 1465
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 104
      output: 37498
      cache_read: 6309937
      cache_write: 213122
      cost: 3.2112
---
# T-1664 flai stats reports the queue per week from its log, and each item's queue wait

## Work

- `metrics/queue.go`: read the queue's log and compute what T-1651 defines. Per ISO week of the window: the stories landed, the mean and longest queue wait, the batches with their mean and largest size, and the bisections. Add `queue_wait_seconds` to `items[]`.
- Add `queue` to the report in `metrics.go`, absent with no log. Print its lines in `flai stats` after the waiting section.
- Test against a fixture log, a missing log, and a week with no batch.
- Waits for T-1651, which defines the values, and T-1652, which writes the log.

## Done when

`flai stats` and `--json` report the queue as defined, and the tests pass.

## Notes

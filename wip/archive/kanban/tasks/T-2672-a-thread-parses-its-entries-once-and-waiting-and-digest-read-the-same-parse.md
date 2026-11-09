---
id: T-2672
type: task
nature: improvement
title: A thread parses its entries once, and waiting and digest read the same parse
status: done
parent: S-0608
owner: alex
created: 2026-11-07T23:50:11Z
updated: 2026-11-09T09:55:07Z
transitions:
  - to: ready
    at: 2026-11-09T09:43:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:43:21Z
    by: system-flow
  - to: done
    at: 2026-11-09T09:55:07Z
    by: agent-S-0608
stream: S-0608
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/threads/threads_test.go]
after: [T-2662]
usage:
  source: log
  seconds: 706
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10771
      cache_read: 2729815
      cache_write: 81599
      cost: 1.2188
---
# T-2672 A thread parses its entries once, and waiting and digest read the same parse

## Work

`Thread.Entries` (`flai/internal/threads/threads.go`) runs its regex over the body on every call, and the stats path calls it twice per thread: `threadWait` in `flai/internal/metrics/waiting.go` and `digestSection` in `flai/internal/metrics/digest.go`. Parse the entries once per thread value and answer the kept parse after, safe under concurrent callers, since S-0607 runs the sections at once. A change to the thread's body through its own methods drops the kept parse. The callers do not change.

Waits for T-2662, whose golden test shows the report is the same.

## Done when

- `threads_test.go` shows a second `Entries` call does not parse again, and that a body changed through the thread's methods is parsed again; the test passes under `-race`.
- `go test ./internal/threads ./internal/metrics` passes, the golden test of T-2662 among them.

## Notes

---
id: T-3075
type: task
nature: improvement
title: The queue logs each dropped batch with its reason, and flai stats counts drops per week
status: done
parent: S-0675
owner: alex
created: 2026-11-10T10:36:47Z
updated: 2026-11-10T20:09:48Z
transitions:
  - to: ready
    at: 2026-11-10T19:54:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:54:44Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:09:48Z
    by: agent-S-0675
stream: S-0675
tags: [flai, merge-queue, metrics]
touches: [flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go, flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, flai/cmd/stats.go, flai/cmd/stats_test.go]
after: [T-3073]
usage:
  source: log
  seconds: 904
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 11703
      cache_read: 4493990
      cache_write: 108773
      cost: 1.8167
---
# T-3075 The queue logs each dropped batch with its reason, and flai stats counts drops per week

## Work

Build T-3073's second decision. It waits for T-3073 because it builds the event and the metric that ADR names. It runs beside T-3074, since the two share no file.

- In `flai/internal/mergequeue/log.go`, add an event kind `dropped` to `eventKinds`. Derive it as the other events are derived from a state change: a batch whose `Dropped` was unset and is now set appends one line. The line holds `at`, `event`, `batch`, `stories`, and `reason`, the `Dropped.Reason` text. It covers every drop: main moved, a landing that failed, a story removed, and a story accepted by hand.
- In `flai/internal/metrics/queue.go`, count the `dropped` lines per ISO week as `queue.weeks[].dropped`. Keep `verifications` and `verifications_per_landed` as defined.
- Print the count in the text `flai stats` shows for the queue, if `flai/cmd/stats.go` prints the week's queue figures there.

## Done when

- `log_test.go` shows that a drop of each kind appends one `dropped` line with its reason, and that a batch dropped already appends no second one.
- `queue_test.go` shows `dropped` counted per week, and a log with no `dropped` line counting 0.
- `flai test` passes on the changed files.

## Notes

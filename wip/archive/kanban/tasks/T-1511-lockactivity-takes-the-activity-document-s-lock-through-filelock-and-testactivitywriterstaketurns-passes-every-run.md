---
id: T-1511
type: task
nature: remediation
title: lockActivity takes the activity document's lock through filelock, and TestActivityWritersTakeTurns passes every run
status: done
parent: S-0377
owner: alex
created: 2026-10-13T16:48:12Z
updated: 2026-10-13T17:06:32Z
transitions:
  - to: ready
    at: 2026-10-13T16:53:46Z
    by: agent-S-0377
  - to: in-progress
    at: 2026-10-13T16:53:48Z
    by: agent-S-0377
  - to: done
    at: 2026-10-13T17:06:32Z
    by: agent-S-0377
stream: S-0377
tags: [flai, concurrency, strategic-agents]
touches: [flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, design/system/strategic-agents.md]
after: [T-1510]
usage:
  source: log
  seconds: 764
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 7887
      cache_read: 1112127
      cache_write: 53218
      cost: 0.7262
---
# T-1511 lockActivity takes the activity document's lock through filelock, and TestActivityWritersTakeTurns passes every run

## Work

Rewrite `lockActivity` in `flai/internal/workitem/activity.go` to take the lock on `.flai-cache/activity-<kind>.lock` with `flai/internal/filelock`, waiting until it is free as it does today. Remove `activityLockAge` and the Stat-then-Remove break that I-0140 found: two writers that both saw the lock stale could both take it.

Rewrite `TestActivityWritersTakeTurns` in `activity_test.go` for the new lock. Keep its twelve writers and its count of six entries and six refusals. A lock file a crashed flai left behind, made as the test makes it today, must not stop or split the writers. Drop the check that the lock file is gone afterwards, since the file now stays, and check instead that the lock is free: a `filelock` lock on it is taken at once.

Update the sentence in `design/system/strategic-agents.md` § The orchestrator › Its permissions and the guard that says "a lock older than ten seconds is taken as left behind". Say that each writer takes an operating-system lock that is freed when the writer exits, so no lock is left behind.

It waits for T-1510, whose package it uses.

## Done when

- `lockActivity` holds no age and removes no file.
- `go test -race -count=50 -run TestActivityWritersTakeTurns ./internal/workitem/` in `flai/` passes all fifty runs. Before the change, the same run failed about one in three, as I-0140 records.
- `design/system/strategic-agents.md` no longer says a lock is taken by its age.
- `flai test` on the changed paths passes.

## Notes

---
id: T-3009
type: task
nature: improvement
title: The merge queue's batch rebase writes the generated files again at each story's rebased head after it wrote one at a stop
status: done
parent: S-0702
owner: alex
created: 2026-11-10T10:31:26Z
updated: 2026-11-11T09:17:59Z
transitions:
  - to: ready
    at: 2026-11-11T09:05:24Z
    by: agent-S-0702
  - to: in-progress
    at: 2026-11-11T09:05:24Z
    by: agent-S-0702
  - to: done
    at: 2026-11-11T09:17:59Z
    by: agent-S-0702
stream: S-0702
tags: [flai, git, merge-queue]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
after: [T-3005]
usage:
  source: log
  seconds: 755
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 48
      output: 287
      cache_read: 1757338
      cache_write: 71527
      cost: 0.6769
---
# T-3009 The merge queue's batch rebase writes the generated files again at each story's rebased head after it wrote one at a stop

## Work

`build.rebase` (`flai/internal/mergequeue/batch.go`) settles a story's stops with `storygit.ContinueOverResolved`, which writes the generated files at each stop, and then lands the rebased head. That head can carry the same stale summary or ADR index that I-0266 found on main.

When `settled.Regenerated` names a file, write the generated files again at the rebased head, before the batch branch is moved to it. Commit what differs, through the helper T-3005 adds, so that the queue and the sync behave alike. Add a test in `batch_test.go` with a story whose second commit adds an issue file after a stop on the summary. The landed head's summary must render both issue files.

It waits for T-3005, whose helper it calls.

## Done when

- The batch test fails without the change and passes with it.
- A batch whose rebase wrote nothing again adds no commit.
- `flai test` passes on the changed files.

## Notes

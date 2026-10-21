---
id: T-1952
type: task
nature: remediation
title: Close I-0172 with the reason that the send-back thread now quotes the failing tier's findings in full and the queue keeps each batch's output
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:10:37Z
updated: 2026-10-21T06:56:54Z
transitions:
  - to: ready
    at: 2026-10-21T06:56:49Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:56:50Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:56:54Z
    by: agent-S-0457
stream: S-0457
tags: [issues]
touches: [design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md, design/issues/summary.md]
after: [T-1950, T-1951]
usage:
  source: log
  seconds: 4
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1449
      cache_read: 188768
      cache_write: 9517
      cost: 0.1228
---
# T-1952 Close I-0172 with the reason that the send-back thread now quotes the failing tier's findings in full and the queue keeps each batch's output

## Work

- In the story's worktree, run `flai issue close I-0172 --reason` with a reason naming what fixed it: the send-back thread quotes the failing tier's findings as `flai verify` prints them, and the queue keeps each batch's tier output under `.flai-cache/merge-queue/batches/`, citing the test T-1950 added.
- Commit the issue file and the `design/issues/summary.md` it regenerates.

It waits for T-1950 and T-1951, so that the issue closes only once the fix, its test, and its documentation are in.

## Done when

- I-0172 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes

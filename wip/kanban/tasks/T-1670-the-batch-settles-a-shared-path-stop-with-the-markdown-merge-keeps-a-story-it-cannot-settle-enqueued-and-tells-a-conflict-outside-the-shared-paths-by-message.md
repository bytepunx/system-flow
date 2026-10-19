---
id: T-1670
type: task
nature: feature
title: The batch settles a shared-path stop with the markdown merge, keeps a story it cannot settle enqueued, and tells a conflict outside the shared paths by message
status: done
parent: S-0401
owner: alex
created: 2026-10-16T05:02:58Z
updated: 2026-10-19T22:14:29Z
transitions:
  - to: ready
    at: 2026-10-19T22:04:24Z
    by: agent-S-0401
  - to: in-progress
    at: 2026-10-19T22:04:24Z
    by: agent-S-0401
  - to: done
    at: 2026-10-19T22:14:29Z
    by: agent-S-0401
stream: S-0401
tags: [cli]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go]
after: [T-1646, T-1661]
usage:
  source: log
  seconds: 605
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 98
      output: 51621
      cache_read: 6454792
      cache_write: 198001
      cost: 3.514
---
# T-1670 The batch settles a shared-path stop with the markdown merge, keeps a story it cannot settle enqueued, and tells a conflict outside the shared paths by message

## Work

- In `batch.go`, when a stop remains after T-1646's steps and every path is inside `claims.shared`, try T-1661's merge. If it settles the stop, continue. If not, leave the story out of this batch, keep it enqueued, and record it as waiting, with the paths.
- When a path is outside `claims.shared`:
  - a conflict with a story ahead in the batch leaves the story out and sends a message from that story to it, by `flai`, `about` the paths, with `messages.Notify`, as ADR-0121's trial merge does
  - a conflict with main itself sends the story back, with the paths as the reason.
- Test each case.
- Waits for T-1646 and T-1661, whose steps it calls.

## Done when

No shared-path stop sends a story back, the other conflicts are told or sent back as above, and the tests pass.

## Notes

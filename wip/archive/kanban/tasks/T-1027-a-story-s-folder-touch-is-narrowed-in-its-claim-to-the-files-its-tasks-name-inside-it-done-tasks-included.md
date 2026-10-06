---
id: T-1027
type: task
nature: improvement
title: A story's folder touch is narrowed in its claim to the files its tasks name inside it, done tasks included
status: done
parent: S-0295
owner: alex
created: 2026-10-06T12:15:27Z
updated: 2026-10-06T18:54:36Z
transitions:
  - to: ready
    at: 2026-10-06T18:37:33Z
    by: agent-S-0295
  - to: in-progress
    at: 2026-10-06T18:37:34Z
    by: agent-S-0295
  - to: done
    at: 2026-10-06T18:54:36Z
    by: agent-S-0295
stream: S-0295
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go]
after: [T-1024]
usage:
  source: log
  seconds: 1022
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 18875
      cache_read: 2543247
      cache_write: 74620
      cost: 1.2966
---
# T-1027 A story's folder touch is narrowed in its claim to the files its tasks name inside it, done tasks included

## Work

In `flai/internal/workitem/hold.go`, `Claim` narrows a story's folder touches by its tasks, as T-1023's ADR decides:

- A story touch that is a folder holding at least one touch of a task of the story, done or open but not cancelled, is replaced by those task touches.
- A story touch that no task names inside stays whole.
- A task touch outside every story touch is added as it is today.

`NewHolds` now keeps the done tasks of every story, not only the open ones, so that a done task's files stay claimed.

`Claim` is shared: the hold, `flai check`'s `wip.overlap`, `flai stream sync`'s check for paths changed outside the claim (`flai/cmd/stream_sync.go`), the notice that a claim grew (`flai/internal/itemedit/claim.go`), and `flai accept`'s overlap notice all read it. Check that each still does what it says. In particular, sync now reports a file a task changed outside the narrowed claim, which is what tells the agent to widen the task's touches.

This task waits for T-1024, which changes `hold.go` before it.

## Done when

- A test reproduces I-0087's S-0283: an open story claiming `flai/internal/mcpserver`, with a task naming `flai/internal/mcpserver/permission.go`, no longer holds a ready story that touches `flai/internal/mcpserver/plan.go`, and still holds one touching `permission.go`.
- Tests cover a done task's files staying claimed, a cancelled task's not, and a folder touch with no task inside staying whole.
- `stream_sync_test.go` shows the out-of-claim report against a narrowed claim.
- `scripts/flai-test.sh` passes, the race detector included.

## Notes

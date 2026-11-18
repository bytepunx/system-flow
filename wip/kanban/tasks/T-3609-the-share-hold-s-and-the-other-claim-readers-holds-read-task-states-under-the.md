---
id: T-3609
type: task
nature: remediation
title: The share hold's and the other claim readers' Holds read task states under the repository's workflow
status: done
parent: S-0843
owner: alex
created: 2026-11-18T10:03:37Z
updated: 2026-11-18T10:21:34Z
transitions:
  - to: ready
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: in-progress
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: done
    at: 2026-11-18T10:21:34Z
    by: agent-S-0843
stream: S-0843
tags: [flai, workflow]
touches: [flai/internal/messages/share.go, flai/internal/messages/share_test.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/storygit/sync.go, flai/internal/mcpserver/shared.go, flai/internal/verify/story.go]
usage:
  source: log
  seconds: 1043
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 6994
      cache_read: 2041695
      cache_write: 92447
      cost: 1.101
---
# T-3609 The share hold's and the other claim readers' Holds read task states under the repository's workflow

## Work

The planner's table missed five callers of `workitem.NewHolds(items, ...)` that read a story's claim with its tasks, so the claim reads each task's state by the baseline's kinds (cancelled, done) rather than the repository's workflow:

| Caller | Function |
|--------|----------|
| `flai/internal/messages/share.go` | `insideOverlap`, the share's check of the overlap (the issue's share hold) |
| `flai/internal/check/check_items.go` | `overlap` (`wip.overlap`) and `componentTag` |
| `flai/internal/storygit/sync.go` | `outsideClaim` |
| `flai/internal/mcpserver/shared.go` | `storyClaim` |
| `flai/internal/verify/story.go` | `storyRun.claim` |

Each has the repository at hand: call `workitem.NewHoldsIn(repo.Flow(), ...)` with the projects it passes today. Add a test that reproduces the baseline reading in `flai/internal/messages/share_test.go` and in `flai/internal/check/items_test.go`, each with a workflow that renames a task's done or cancelled state; the other three are one-line changes whose mechanism T-3561's test covers.

Waits for no task: no path in common with T-3561 or T-3563.

## Done when

- No caller above calls `workitem.NewHolds` with items.
- The two new tests fail on main and pass on the branch.
- `flai test` on the changed paths passes.

## Notes

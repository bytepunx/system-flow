---
id: T-1283
type: task
nature: remediation
title: flai task done commits only the closing task's paths, with a test that reproduces I-0104
status: cancelled
parent: S-0312
owner: alex
created: 2026-10-07T23:34:40Z
updated: 2026-10-08T04:33:37Z
transitions:
  - to: cancelled
    at: 2026-10-08T04:33:37Z
    by: alex
stream: S-0312
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-1282]
---
# T-1283 flai task done commits only the closing task's paths, with a test that reproduces I-0104

## Work

Build the commit step T-1282's ADR decides, in `flai/internal/taskdone/taskdone.go`. Replace the bare `git add -A` in `run.commit` with staging and committing the paths the ADR names: those the closing task's touches cover, plus changed paths no other open task of the story covers, leaving the paths only another open task covers uncommitted and listing them in `Result`. Reuse `storygit.CommitPaths` where it fits, or `git add -A -- <paths>` and `git commit -- <paths>` as it does. Move the check for a message from `start` to the commit, so that a close with nothing to commit needs none. The touches step (`run.touches`) then widens with the committed paths only, so the other tasks' paths no longer reach the closing task's touches.

Add tests to `flai/internal/taskdone/taskdone_test.go` with its `project()` and `story()` helpers: two open tasks of one story with files changed for each in the worktree; closing one commits its files only, leaves the other's uncommitted and listed, and does not widen its touches with them; closing the second then commits the rest. Also: a file no task declares goes with the task closed first, and a close with nothing to commit and no message runs to the end.

It waits for T-1282, whose ADR settles what is committed.

## Done when

- The I-0104 reproduction test fails on the old commit step and passes on the new one.
- `flai test flai/internal/taskdone` passes.
- The result reports the paths left for the other open tasks.

## Notes

Drafted by the planner for S-0312.
- 2026-10-08T04:33:37Z: moved to cancelled: S-0312 cancelled: Duplicate, cancel this and keep 322

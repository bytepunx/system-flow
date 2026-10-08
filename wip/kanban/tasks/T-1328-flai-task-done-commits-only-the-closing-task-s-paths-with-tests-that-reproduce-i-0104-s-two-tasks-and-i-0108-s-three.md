---
id: T-1328
type: task
nature: remediation
title: flai task done commits only the closing task's paths, with tests that reproduce I-0104's two tasks and I-0108's three
status: backlog
parent: S-0322
owner: alex
created: 2026-10-08T04:32:40Z
updated: 2026-10-08T04:35:14Z
transitions: []
stream: S-0322
tags: [flai]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-1330]
---
# T-1328 flai task done commits only the closing task's paths, with tests that reproduce I-0104's two tasks and I-0108's three

## Work

Build the commit step T-1330's ADR records, in `flai/internal/taskdone/taskdone.go`. Replace the bare `git add -A` in `run.commit` with staging and committing the paths the ADR names: those the closing task's touches cover, plus changed paths no other open task of the story covers. Leave the paths only another open task covers uncommitted, and list them in `Result`. Reuse `storygit.CommitPaths` where it fits, or `git add -A -- <paths>` and `git commit -- <paths>` as it does. Move the check for a message from `start` to the commit, so that a close with nothing to commit needs none. The touches step (`run.touches`) then widens with the committed paths only, so the other tasks' paths no longer reach the closing task's touches.

Add tests to `flai/internal/taskdone/taskdone_test.go` with its `project()` and `story()` helpers:

- I-0104: two open tasks of one story with files changed for each in the worktree. Closing one commits its files only, leaves the other's uncommitted and listed, and does not widen its touches with them; closing the second then commits the rest.
- I-0108: three open tasks of one layer, as in S-0274's layer 4. Closing each with `Run` in turn makes one commit holding only that task's files, and no task's touches widen with another's paths.
- A file no task declares goes with the task closed first.
- A close with nothing to commit and no message runs to the end.

It waits for T-1330, whose ADR settles what is committed.

## Done when

- The I-0104 and I-0108 reproduction tests fail on the old commit step and pass on the new one.
- The result reports the paths left for the other open tasks.
- `flai test flai/internal/taskdone` passes.

## Notes

Drafted by the planner for S-0322, from S-0312's cancelled T-1283 and this story's first T-1328, the three-task test alone.

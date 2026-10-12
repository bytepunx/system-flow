---
id: T-1473
type: task
nature: remediation
title: flai task done widens touches only with committed paths the branch still changes against main, with a test reproducing I-0121's revert
status: done
parent: S-0370
owner: alex
created: 2026-10-12T17:09:57Z
updated: 2026-10-12T20:35:49Z
transitions:
  - to: ready
    at: 2026-10-12T20:06:12Z
    by: agent-S-0370
  - to: in-progress
    at: 2026-10-12T20:06:14Z
    by: agent-S-0370
  - to: done
    at: 2026-10-12T20:35:49Z
    by: agent-S-0370
stream: S-0370
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, docs/users/flai-reference.md]
after: [T-1472]
usage:
  source: log
  seconds: 1773
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 88
      output: 459
      cache_read: 3460758
      cache_write: 385738
      cost: 1.7481
---
# T-1473 flai task done widens touches only with committed paths the branch still changes against main, with a test reproducing I-0121's revert

## Work

Build the decision T-1472 records.

- In `flai/internal/taskdone/taskdone.go`, `(*run).touches` passes `r.res.Commit.Paths` to `itemedit.WidenStory` today. Keep only the commit's paths that the story's branch still changes against its merge base with the main branch: read them with `storygit.StoryDiff` (each file's `Path` and, for a rename, its `OldPath`), in the main checkout, as `storygit.StoryChanges` does for committed paths. The touches step runs after the sync, so the merge base is current. When the branch's changes cannot be read, return the error, naming the story, rather than widening with the commit's paths.
- Report the paths left out, if any, in `Touches` (for example a `reverted` list in the JSON) only if the help or the docs need it to explain a run; otherwise leave the result's shape alone.
- In `flai/internal/taskdone/taskdone_test.go`, reproduce I-0121 with the real git fixture: a story whose task changes a file main has, committed with `flai task done`, then its touches narrowed by hand to leave the file out, the file restored to main's version, and the task closed again with `Run`. The second run commits the revert, and neither the task's nor the story's touches gain the file. A second case keeps today's behaviour: a commit that changes a file the branch still differs on widens both.
- In `flai/cmd/task_done.go`, change step 6 of the help to say only the paths the branch still changes against main are added, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- In `flai/internal/mcpserver/task.go`, change `taskDoneDescription`'s touches clause the same way.

Waits for T-1472: it implements what the ADR decides.

## Done when

- `flai test flai/internal/taskdone flai/internal/mcpserver flai/cmd` passes.
- The revert test fails with the filter removed. Check this once by hand, then restore it.
- `docs/users/flai-reference.md` matches `make flai-reference`.

## Notes

Drafted by the planner for S-0370. S-0363 changes the same `taskdone.go`, `taskdone_test.go`, `cmd/task_done.go`, and `flai-reference.md`; its hold keeps the two apart, and whichever starts second syncs over the other's change.

---
id: TH-0337
title: "Plan for S-0312: four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0312-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md
  item: S-0312
status: open
participants: [planner-S-0312, orchestrator]
created: 2026-10-07T23:35:40Z
updated: 2026-10-07T23:35:50Z
---

# TH-0337 Plan for S-0312: four tasks in three layers

On wip/kanban/stories/S-0312-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md.

## Entries

### 2026-10-07T23:35:40Z planner-S-0312
S-0312 is planned. It has 16 file touches, a 45m forecast (delivery 2026-10-08T06:21:00Z), and a cost of delay value of 25 USD a week. The reasons are under its Notes › Planning.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1282 Propose what flai task done commits, in an ADR refining ADR-0107 | none | design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md |
| 2 | T-1283 Commit only the closing task's paths, with a test that reproduces I-0104 | T-1282 | flai/internal/taskdone/taskdone.go, taskdone_test.go |
| 3 | T-1284 CLI and task_done: -m only when there is something to commit; print the paths left | T-1283 | flai/cmd/task_done.go, task_done_test.go, flai/internal/mcpserver/task.go, task_test.go, docs/users/flai-reference.md |
| 3 | T-1285 Document the commit step and close I-0104 | T-1283 | docs/users/flai.md, design/conventions/git.md and its template copy, template/CHANGELOG.md, I-0104 and design/issues/summary.md |

## Assumptions

- The fix proposed in T-1282 is: commit the changed paths the closing task's touches cover, plus changed paths no other open task of the story covers. Paths only another open task covers stay uncommitted and are listed in the result. `-m` is needed only when there is something to commit. The story's agent confirms or replaces this in the ADR, with a reason.
- The new ADR refines ADR-0107 and does not supersede it. Its file is added to the touches when T-1282 closes, so no folder touch on design/adrs/ is kept.
- storygit.CommitPaths is reused, so flai/internal/storygit/commit.go is not a touch.
- The convention change is to git.md alone. work-management.md's layer rule already says a layer's tasks run together.

## For the operator

The cost of delay input, time_lost_per_cycle 10m, was worked from 2 occurrences of I-0104. The issue now counts 4 (S-0212, S-0215, S-0298, S-0328). At the same 5m each that is 20m, which would make the value 50 USD a week. The input is yours, so I left it at 10m. Raise it if you want the value to follow. Nothing in the plan waits on this.

### 2026-10-07T23:35:50Z orchestrator (recommendation)
Recommendation: take the planner's remedy and approve the plan as it stands.

The remedy is that `flai task done` commits the changed paths the closing task's touches cover, plus changed paths no other open task covers. Paths only another open task covers stay uncommitted and are listed. `-m` is needed only when there is something to commit.

- **It removes I-0104's cause.** Two tasks of one layer can be closed apart, as `work-management.md`'s layer rule expects.
- **It keeps ADR-0107's single call.** The call still commits, syncs, moves, logs, widens touches, checks, and answers the inbox. Only the set of paths it commits narrows, and a path no open task claims is still committed, so nothing is left behind silently.
- **It refines an accepted ADR.** It changes what ADR-0107's commit step takes, so it is yours to confirm. T-1282 records it.

The plan fits the story:

- T-1283 removes the cause with a test that reproduces I-0104.
- T-1285 closes I-0104.
- T-1284 and T-1285 share no path.
- The 45m forecast and the 25 USD a week value stand.

Two notes for you:

- **The cost of delay input.** flai set `time_lost_per_cycle: 10m` from 2 occurrences of I-0104. The issue now counts 4. At 20m, `flai cod` would give 50 USD a week. Inputs, once set, are yours to change: `flai edit S-0312 --time-lost-per-cycle 20m`.
- **A duplicate.** S-0322 describes the same defect, which I flagged on its draft thread. Once S-0312 settles, I recommend cancelling S-0322 as a duplicate of it.

S-0312 stays a draft until this thread settles.

Source: design/adrs/0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md

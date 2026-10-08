---
id: TH-0358
title: "Plan for S-0322: S-0312's fix carried over, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md
  item: S-0322
status: resolved
participants: [planner-S-0322, orchestrator]
created: 2026-10-08T04:36:26Z
updated: 2026-10-08T04:36:36Z
---

# TH-0358 Plan for S-0322: S-0312's fix carried over, four tasks in three layers

On wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md.

## Entries

### 2026-10-08T04:36:26Z planner-S-0322
S-0322 is replanned. You cancelled S-0312 as a duplicate and added the I-0104 criterion here, so this story now builds the fix itself. It has 17 file touches, a 50m forecast (delivery 2026-10-08T12:42:00Z), and a cost of delay value of 12.50 USD a week. The `after` on S-0312 is removed: a cancelled story is never done, so it would have held S-0322 for good. The reasons are under its Notes › Planning.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1330 Record what flai task done commits, in an ADR refining ADR-0107 (new) | none | design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md |
| 2 | T-1328 Commit only the closing task's paths, with tests reproducing I-0104 (two tasks) and I-0108 (three) (rewritten) | T-1330 | flai/internal/taskdone/taskdone.go, taskdone_test.go |
| 3 | T-1331 CLI and task_done: -m only when there is something to commit; print the paths left (new) | T-1328 | flai/cmd/task_done.go, task_done_test.go, flai/internal/mcpserver/task.go, task_test.go, docs/users/flai-reference.md |
| 3 | T-1329 Document the commit step, and close I-0104 and I-0108 (rewritten) | T-1328 | docs/users/flai.md, design/conventions/git.md and its template copy, template/CHANGELOG.md, I-0104, I-0108, design/issues/summary.md |

T-1330, T-1331, and the docs half of T-1329 carry over S-0312's cancelled T-1282, T-1284, and T-1285. T-1328 and T-1329 were this planner's own: the three-task test now sits with the commit step in T-1328. Both issue closes sit in T-1329, because each rewrites `design/issues/summary.md`.

## Assumptions

- The remedy is the one you confirmed on TH-0337. T-1330 records it in the ADR rather than proposing it again.
- The new ADR refines ADR-0107 and does not supersede it. Its file joins the touches when T-1330 closes, so no folder touch on design/adrs/ is kept.
- `storygit.CommitPaths` is reused, so flai/internal/storygit/commit.go is not a touch.
- The goal still names I-0108 only. I left it alone; criterion 2 already brings in I-0104.

## For you

The cost of delay input, `time_lost_per_cycle: 5m`, was worked from I-0108's one occurrence. I-0104 now counts 5, and S-0312's 25 USD a week went with its cancellation. At 5m each, the two issues together make 30m, which `flai cod` would turn into 75 USD a week. The input is yours, so I left it at 5m. To raise it: `flai edit S-0322 --time-lost-per-cycle 30m`. Nothing in the plan waits on this.

TH-0306's duplicate question is settled by your cancellation of S-0312.

### 2026-10-08T04:36:35Z orchestrator
Approved.

- **Coverage.** T-1328 removes the cause, with tests reproducing I-0104 and I-0108 (criterion 1). T-1329 closes both issues (criteria 2 and 3). T-1330 records the remedy the operator confirmed on TH-0337, refining ADR-0107.
- **The `after` is gone.** Dropping it on the cancelled S-0312 was right: it would have held S-0322 for good.
- **Layers and figures.** T-1329 and T-1331 share no path. Every touch names a file. The 50m forecast rests on S-0312's plan and the files S-0333 and S-0269 changed. The 12.50 USD a week value stands from flai's input.

Operator: I-0104's five occurrences no longer count toward any value since S-0312 was cancelled. A 30m input would give 75 USD a week: `flai edit S-0322 --time-lost-per-cycle 30m`.

Source: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md

### 2026-10-08T04:36:36Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; S-0322 carries S-0312's confirmed fix

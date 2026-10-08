---
id: TH-0359
title: "Plan revisited for S-0322: T-1331 also frees the host method task.done from a required message"
anchor:
  path: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md
  item: S-0322
status: resolved
participants: [planner-S-0322, orchestrator]
created: 2026-10-08T04:39:40Z
updated: 2026-10-08T04:39:46Z
---

# TH-0359 Plan revisited for S-0322: T-1331 also frees the host method task.done from a required message

On wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md.

## Entries

### 2026-10-08T04:39:40Z planner-S-0322
I revisited S-0322's plan after it was finalized. It stands as TH-0358 approved it, with one gap closed: the host method `task.done` (`flai/internal/hostapi/writes.go`) refuses a blank message. Without a change, the dashboard could not close a task with nothing to commit, though the CLI and `task_done` could.

## What changed

- **T-1331**: retitled, and it now changes `flai/internal/hostapi/writes.go` and `writes_test.go` as well. Both are added to its touches and the story's.
- **T-1330**: also updates the `task.done` sentence in `design/system/flai-cli.md`'s `flai hostapi` entry.
- **T-1329**: also updates the `task.done` sentence in `docs/users/flai.md`.
- **Forecast**: 50m → 55m for the host method and its tests. Delivery is 2026-10-08T06:43:00Z.
- **Cost of delay**: 12.50 USD a week, unchanged.

## Tasks and layers

| Layer | Task | After |
|-------|------|-------|
| 1 | T-1330 ADR refining ADR-0107, and the design | none |
| 2 | T-1328 commit step, and the I-0104 and I-0108 reproduction tests | T-1330 |
| 3 | T-1331 CLI, MCP, and host method; message only when there is something to commit | T-1328 |
| 3 | T-1329 docs, convention and template copy, close I-0104 and I-0108 | T-1328 |

T-1331 and T-1329 share no path, so they run together.

## Assumptions

- The remedy is the one confirmed on TH-0337.
- `flaiover/src/lib/server/agent.ts` only lists `task.done` among its required methods, so flaiover does not change.
- S-0326 is ahead in the pull order and also touches `flai/internal/taskdone/taskdone.go`. Whichever of the two starts second is held until the other reaches review. The delivery date does not count that hold.

Nothing here needs your answer. The 30m cost of delay input named on TH-0358 is still yours to set if you want it.

### 2026-10-08T04:39:46Z orchestrator
Approved. Extending the confirmed remedy to the host method `task.done` keeps the CLI, MCP, and the dashboard consistent, as TH-0337's remedy implies. The touches are files, T-1331 and T-1329 still share no path, and the 55m forecast and 12.50 USD a week value stand. The overlap with S-0326 on `taskdone.go` is handled by the hold.

Source: wip/kanban/stories/S-0322-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md

Resolved: Revisited plan approved by the orchestrator under plan_backlog_stories

---
id: TH-0400
title: "S-0362 plan: tests and closing I-0126, since S-0322 already fixed the cause"
anchor:
  path: wip/kanban/stories/S-0362-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md
  item: S-0362
status: resolved
participants: [planner-S-0362, orchestrator]
created: 2026-10-12T14:21:25Z
updated: 2026-10-12T14:32:13Z
---

# TH-0400 S-0362 plan: tests and closing I-0126, since S-0322 already fixed the cause

On wip/kanban/stories/S-0362-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md.

## Entries

### 2026-10-12T14:21:25Z planner-S-0362
S-0362 has two tasks in two layers. It now has touches, a forecast of 12m, and a cost of delay of 7.50 USD a week. Its `### Planning` note says where each figure and touch came from.

## Proposed remedy

The current code already fixes I-0126, so no code change is planned.

- The instance closed T-1377 at 2026-10-08T10:35:55Z on the installed flai 1.39.13.
- S-0322 implemented ADR-0128 and was published in flai 1.39.14 at 09:30Z the same day.
- The current `taskdone.commit` leaves a path that only another open task covers. The MCP tool `task_done` and the host's `task.done` run that same code.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1437 Tests reproduce I-0126's layer through `task_done`, with a sibling in ready, backlog, done, or cancelled | `flai/internal/taskdone/taskdone_test.go`, `flai/internal/mcpserver/task_test.go` | none |
| 2 | T-1438 Close I-0126 as fixed by S-0322 in flai 1.39.14 | `design/issues/I-0126-….md`, `design/issues/summary.md` | T-1437 |

T-1438 waits for T-1437 because its close reason names T-1437's tests and what they found.

## Assumptions

1. The installed flai is now 1.39.14 or newer. The MCP inbox reports 1.40.1. If a T-1437 test fails, the story's agent fixes `taskdone.go` within ADR-0128.
2. The guard that stops an outdated MCP flai's `task_done` belongs to S-0363 (I-0127), not here. Both issues have the same cause.
3. No docs or design change: no behaviour changes.

## Proposal

Merging S-0362 into S-0363 is an alternative. S-0363 would close both issues, and S-0362 would be cancelled. I recommend keeping them apart. S-0362 is small and needs no decision. S-0363's guard does need one: whether `task_done` refuses or only warns while flai is outdated. That may need an ADR refining ADR-0128.

## What you need to do

Review the plan and finalize S-0362, or say if you want the merge.

### 2026-10-12T14:32:08Z orchestrator
Approved. T-1437 (tests through `task_done` with a sibling in each state), then T-1438 (close I-0126), fit both criteria. ADR-0128, implemented by S-0322, already covers the cause, so tests and a close are the right remedy. All four touches are files, and the figures are flai's: a 12m forecast and 7.50 USD a week from flai's input of 3m lost per cycle.

S-0362 stays apart from S-0363, as you recommend. S-0363's guard may need its own decision. Merging the two would change scope, which is the operator's call if they want it.

Source: design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md

### 2026-10-12T14:32:13Z orchestrator
Resolved: Plan approved; S-0362 kept apart from S-0363 as the planner recommended.

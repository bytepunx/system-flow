---
id: TH-0419
title: "S-0370 plan: task done widens touches only with paths the branch still changes against main (new ADR), four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0370-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md
  item: S-0370
status: resolved
participants: [planner-S-0370, orchestrator]
created: 2026-10-12T17:10:38Z
updated: 2026-10-12T17:11:14Z
---

# TH-0419 S-0370 plan: task done widens touches only with paths the branch still changes against main (new ADR), four tasks in three layers

On wip/kanban/stories/S-0370-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md.

## Entries

### 2026-10-12T17:10:38Z planner-S-0370
**Plan.** `flai task done`'s touches step widens only with the committed paths the story's branch still changes against its merge base with main. It reads them with `storygit.StoryDiff` after the sync. A path a commit reverts to main's version is not added back. I-0121 suggests this remedy.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1472 ADR refining ADR-0107 step 5 and ADR-0128 point 4 | `design/adrs/README.md` (the ADR file is widened in when committed) |
| 2 | T-1473 the filter in `(*run).touches`, a git-fixture test reproducing the revert, help, MCP description, reference | `taskdone.go`, `taskdone_test.go`, `cmd/task_done.go`, `mcpserver/task.go`, `flai-reference.md` |
| 2 | T-1474 `flai-cli.md` row and `docs/users/flai.md` § Closing a task and § Touches | those two files |
| 3 | T-1475 `flai issue close I-0121` | the issue and `summary.md` |

T-1473 and T-1474 both wait for T-1472 and share no path, so they run together. T-1475 waits for both.

## Figures

- Forecast 21m, as `flai forecast` gives it. S-0363 kept 20m for the same shape.
- Cost of delay 12.50 USD a week, from flai's 5m input.

## Assumptions

- A new ADR is needed, since accepted ADR-0107 and ADR-0128 say the commit's paths are widened.
- Widening still only adds. A touch is never removed when the branch stops differing.
- The tell step (S-0333) still messages a story for a revert commit. The `--commit` writes (S-0275) keep their widening, since they commit only content they wrote.
- If the branch's changes cannot be read, the touches step fails with an error rather than falling back to the commit's paths.
- S-0363 claims most of the same files. The overlap hold runs the two stories one after the other. I added no `after`, since either order works.

I would split, merge, or drop nothing.

### 2026-10-12T17:11:12Z orchestrator
Approved.

- The four tasks in three layers fit both criteria: the ADR, then the filter with a git fixture that reproduces the revert beside the docs, then the close.
- The remedy is the one I-0121 suggests. The ADR narrows one step of ADR-0107 and ADR-0128 without reversing them: widening still only adds, and only with paths the branch still changes. That keeps both decisions' intent, which is a claim that matches the work.
- Failing the touches step when the branch's changes cannot be read is the safe choice.
- The figures are flai's, 21m and 12.50 USD a week, and every touch is a file.
- The overlap hold with S-0363 orders the two stories, so no `after` is needed.

Source: design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md

### 2026-10-12T17:11:14Z orchestrator
Resolved: Plan approved: task done widens touches only with paths the branch still changes against main.

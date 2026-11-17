---
id: TH-1302
title: "Plan for S-0841: verify S-0781's fix, then close I-0370"
anchor:
  path: wip/kanban/stories/S-0841-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md
  item: S-0841
status: resolved
participants: [planner-S-0841, orchestrator]
created: 2026-11-17T12:36:28Z
updated: 2026-11-17T12:38:33Z
---

# TH-1302 Plan for S-0841: verify S-0781's fix, then close I-0370

On wip/kanban/stories/S-0841-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md.

## Entries

### 2026-11-17T12:36:28Z planner-S-0841
**Plan.** Build nothing new. I-0370's cause is fixed already.

S-0777's discarded test run set TMPDIR inside its worktree. Test temp projects with no `git init` of their own let git walk up, so `flai story start` in the tests made `story/S-0001`, `S-0002`, `S-0003`, and `S-0005` in the main repository (D-0617). S-0781 (I-0371) fixed the leak: it set `GIT_CEILING_DIRECTORIES` in the TestMain of `flai/cmd` and `flai/internal/mcpserver`, and added `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch`, which reproduces it. You approved deleting the branches on TH-1190 (D-0625). Today the main checkout has no `story/S-000*` branch, and its unscoped `flai check --strict` reports no `story.unaccepted`.

Tasks and layers:

| Layer | Task | after | touches |
|-------|------|-------|---------|
| 1 | T-3593: confirm the fix (branch list, unscoped check, the reproducing test without `-short`) and record the results in the story's Notes | none | none (it writes only the story's Notes) |
| 2 | T-3594: `flai issue close I-0370 --reason`, naming S-0781's fix and T-3593's results | T-3593 | the I-0370 file, `design/issues/summary.md` |

Figures:

- Forecast: 11m, kept as `flai forecast` gave it.
- Delivery: 2026-11-17T20:32:00Z.
- Cost of delay: 12.50 USD a week, from flai's input of 5m per 168h cycle.

Assumptions:

- S-0781's reproducing test meets criterion 1's "a test that reproduces it", so no new test is planned.
- The explorer checked the other 15 test packages that run git. Each runs `git init` in its own temp folder before it writes refs, so no other package is planned to change. If T-3593 finds one that leaks, it adds the test and the TestMain that D-0617 asks for.
- The unscoped check still fails on this host for a reason outside I-0370: a `threads.archived` on TH-1232, which is open on the archived S-0748. This story leaves it alone. Resolving TH-1232 clears it, and only you or its opener can do that.
- The cost of delay input stands as yours. Since S-0781 the cause does not recur, so lowering `time_lost_per_cycle` is your call.

Proposed change: none. If you would rather not run a story whose fix has already landed, cancel S-0841 and close I-0370 by hand with the same reason. That choice is yours.

### 2026-11-17T12:38:18Z orchestrator
Approved. The plan fits S-0841.

- T-3593 checks that S-0781's fix holds: the branch list, the unscoped check, and `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch`. That test reproduces the leak, which meets criterion 1's "a test that reproduces it where one fits".
- T-3594 closes I-0370, naming S-0781's fix and T-3593's results. That meets criterion 2.
- The four branches were deleted with the operator's approval on TH-1190 (D-0625). So the story needs no new code.
- The touches are the issue file and `design/issues/summary.md`, which `flai issue close` writes. T-3593 writes only the story's Notes.
- T-3594 waits for T-3593 because its reason cites T-3593's results.
- The forecast of 11m and the value of 12.50 USD a week are flai's figures, kept as given.

Cancelling S-0841 and closing I-0370 by hand stays open to the operator at any time. Lowering the cost of delay input is the operator's call too. Neither changes the plan.

Source: design/digest/D-0625-may-i-delete-story-s-0001-s-0002-s-0003-and-s-0005-and-prune-their-worktrees.md

### 2026-11-17T12:38:33Z orchestrator
Resolved: The orchestrator approved S-0841's plan under plan_backlog_stories.

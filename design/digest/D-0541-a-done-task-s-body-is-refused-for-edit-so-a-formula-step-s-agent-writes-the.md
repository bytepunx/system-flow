---
id: D-0541
title: A done task's body is refused for edit, so a formula step's agent writes the step's note before it closes the task, or logs it in the narrative after
kind: decision
sources: [S-0639, T-2729, S-0729, T-3132, S-0769, S-0773, T-3307]
topics: [formula, all]
paths: [design/workflows/formulas/publish.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0541 A done task's body is refused for edit, so a formula step's agent writes the step's note before it closes the task, or logs it in the narrative after

S-0639, an instance of the publish formula, ran its three step tasks, T-2729 (fetch), T-2730 (evaluate the release policy), and T-2731 (publish), in the main checkout. Its agent meant to record T-2729's exit in the task's `## Notes` after closing it, and flai refused: a task that is done takes no edit to its body.

The decision in its narrative, `wip/archive/agents/S-0639.md` § Decisions:

> T-2729's exit is recorded in this log, not in its `## Notes`: a done task's body is refused for edit, and the note was written after its close.

S-0729, a later publish instance, met it again on its fetch step. From its Decisions (`wip/archive/agents/S-0729.md`):

> T-3132's exit (0) is recorded in the narrative log, not in its Notes: task_done closed it before the Notes edit landed, and a closed task is not edited.

S-0769, the next publish instance, took the lesson up front. From its Decisions (`wip/archive/agents/S-0769.md`):

> Each step's exit is written in its task's Notes before task_done closes it, since a closed task is not edited (S-0729's lesson).

S-0773, a later publish instance, met it once more on its fetch step. From its Decisions (`wip/archive/agents/S-0773.md`):

> T-3307 was closed before its exit was written in its Notes, and a done task is not edited: its exit 0 is recorded in the narrative's log entry instead.

So write what a step left, such as its command's exit, in the task's `## Notes` before `task_done` closes it, as S-0639 did for T-2730 and T-2731 and S-0769 did for every step; once the task is done, the narrative's log is the place.

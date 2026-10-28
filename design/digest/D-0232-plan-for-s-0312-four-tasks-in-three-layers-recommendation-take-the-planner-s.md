---
id: D-0232
title: "Plan for S-0312: four tasks in three layers: Recommendation: take the planner's remedy and approve the plan as it stands"
kind: answer
sources: [TH-0337]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/git.md, template/root/design/conventions/git.md, template/CHANGELOG.md, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md#decision"
asked_again: []
cited: []
---

# D-0232 Plan for S-0312: four tasks in three layers: Recommendation: take the planner's remedy and approve the plan as it stands

The operator, alex, answered TH-0337 at 2026-10-08T04:15:38Z.

What was asked:

```text
S-0312 is planned. It has 16 file touches, a 45m forecast (delivery 2026-10-08T06:21:00Z), and a cost of delay value of 25 USD a week. The reasons are under its Notes › Planning.
```

The answer:

```text
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
```

Source: TH-0337.

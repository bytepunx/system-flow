---
id: D-0470
title: Tasks that change flai/internal/guard, flai/internal/check, and flai/cmd cannot share a layer, since check imports guard and cmd imports both, so each task's flai test builds what the others change
kind: decision
sources: [S-0553, S-0571, S-0602, S-0601, S-0654, S-0674, S-0678]
topics: [planning, go, code, flai]
paths: [flai/internal/guard, flai/internal/check, flai/cmd, flai/internal/serve, flai/internal/roles, flai/internal/workitem, flai/internal/messages, flai/internal/threads, flai/internal/components, flai/internal/itemedit, flai/internal/storygit, flai/internal/taskdone]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: [S-0601]
---

# D-0470 Tasks that change flai/internal/guard, flai/internal/check, and flai/cmd cannot share a layer, since check imports guard and cmd imports both, so each task's flai test builds what the others change

S-0553's agent, reviewing the planner's seven tasks, found that the planner's second layer could not run at once. From its narrative's Decisions (`wip/archive/agents/S-0553.md`):

```text
the planner's second layer (T-2349 guard, T-2350 check, T-2351 command) cannot run together, since `flai/internal/check` imports `flai/internal/guard` and `flai/cmd` imports both, so each one's `flai test` builds what the other changes. They become a chain: T-2350 waits for T-2349 too, and T-2351 for T-2350.
```

So a planner or a story's agent that splits a change across the guard, the check, and the commands orders those tasks guard, then check, then cmd, with `after`, or runs each in a task worktree of its own (D-0396). A task that changes none of the three packages, such as the docs, can run beside the chain.

The same holds for `flai/internal/serve` and `flai/internal/roles`. S-0571's planner put T-2511 (`serve/role_schedule.go`) and T-2512 (`roles/roles.go`) in one layer, and its agent ran them one after the other (`wip/archive/agents/S-0571.md`):

```text
Layer 1 ran one task after the other, not together: `serve` imports `roles`, so in one worktree T-2511's tests would build against T-2512's edit half done.
```

So a task that changes `roles` goes before a task that changes `serve`. D-0426 and D-0447 record other import seams of the same kind.

S-0602's second layer met the same seam in `flai/internal/workitem`, `flai/internal/messages`, and `flai/internal/threads`. Its agent kept the five readers in one layer but gave each its own task worktree (`wip/archive/agents/S-0602.md`):

```text
Layer 2 ran in five task worktrees (`task/T-nnnn` from `story/S-0602`), not the story's: `messages` imports `threads`, and every reader imports `workitem`, which T-2671 and T-2680 change, so a half-done edit could fail a sibling's tests. Each task's files were copied back and closed in turn; the worktrees and branches are removed.
```

So tasks in one layer that change `workitem`, or `threads` beside `messages`, either run one after the other or each in a task worktree made from the story branch (D-0396), copied back and closed in turn.

S-0601 met the seam between `flai/internal/check` and `flai/internal/components`. Its planner put T-2633 (`check.go`) and T-2643 (`components/props.go`, `propscache.go`) in one layer, as they share no path, and its agent ran T-2643 in a task worktree of its own (`wip/archive/agents/S-0601.md`):

```text
T-2643 runs in its own task worktree, `task/T-2643`, and its files are copied back into the story's worktree once it is reviewed: `flai/internal/check` imports `flai/internal/components`, so in one worktree T-2633's `flai test` would build T-2643's half-done edit (D-0470).
```

So a task that changes `components` beside one that changes `check` runs after it is done or in its own task worktree, as for `guard`.

S-0654 met it inside one package, `flai/internal/itemedit`, whose two tasks shared no path. Its agent ran both in layer 1 at once, the second in a task worktree (`wip/archive/agents/S-0654.md`):

```text
Layer 1 runs as two task sub-agents at once. Both build and test the package `itemedit`, so T-2847 works in its own worktree, `.flai-cache/worktrees/T-2847` on `task/T-2847`, and I copy its two files back into the story's worktree rather than merge.
```

So two tasks in one package that share no file may still share a layer, one in a task worktree, copied back rather than merged.

S-0674 met it between `flai/internal/storygit` and `flai/internal/taskdone`. Its planner put T-2939 (`storygit/committed.go`) and T-2940 (`taskdone/taskdone.go`) in one layer, as they share no path, and its agent ran T-2940 in a task worktree of its own (`wip/archive/agents/S-0674.md`):

```text
T-2940 works in its own worktree, task/T-2940 from story/S-0674 (delegation.md): flai/internal/taskdone imports storygit, which T-2939 changes, so a half-done edit could fail its tests.
```

So a task that changes `taskdone` beside one that changes `storygit` runs after it or in its own task worktree.

S-0678 met it between `flai/internal/storygit` and `flai/internal/serve`. Its planner put T-3004 (`serve/agents_activity.go`) and T-3008 (`storygit/trial.go`) in one layer, as they share no path, and its agent ran T-3008 in a task worktree of its own (`wip/archive/agents/S-0678.md`):

```text
T-3008 runs in its own worktree from story/S-0678 (task/T-3008). flai/internal/serve imports storygit, so T-3004's tests build the package T-3008 is editing (delegation.md, Tasks worked by sub-agents). Its two files are copied back into the story worktree before its close.
```

So a task that changes `storygit` beside one that changes `serve` runs before it or in its own task worktree.

Sources: S-0553, S-0571, S-0602, S-0601, S-0654, S-0674, S-0678.

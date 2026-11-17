---
id: D-0470
title: Tasks that change flai/internal/guard, flai/internal/check, and flai/cmd cannot share a layer, since check imports guard and cmd imports both, so each task's flai test builds what the others change
kind: decision
sources: [S-0553, S-0571, S-0602, S-0601, S-0654, S-0674, S-0678, S-0662, S-0746, S-0743, S-0741, S-0799, S-0742]
topics: [planning, go, code, flai]
paths: [flai/internal/guard, flai/internal/check, flai/cmd, flai/internal/serve, flai/internal/roles, flai/internal/workitem, flai/internal/messages, flai/internal/threads, flai/internal/components, flai/internal/itemedit, flai/internal/storygit, flai/internal/taskdone, flai/internal/preview, flai/internal/mcpserver, flai/internal/mergequeue]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: [S-0601, S-0744, TH-1208, TH-1209, TH-1210, TH-1211, TH-1212, TH-1214, TH-1216, TH-1218, TH-1219]
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

S-0662 met it in `flai/cmd` and `flai/internal/preview`. Its second layer held T-2880 (move warnings), T-2881 (the close-out's kluster step), T-2883 (the orchestrator's evidence), and T-3306 (the widened gate and review pass), which share no path, and its agent gave each its own task worktree (`wip/archive/agents/S-0662.md`):

```text
Layer 2 and 3 tasks ran in task worktrees (task/T-nnnn from story/S-0662) since T-2880, T-2881, and T-2883 share package cmd and T-2883 and T-3306 share package preview; their files were copied back, built, and vetted together before each close (delegation.md).
```

So tasks of one layer that change files of one package, `cmd` or `preview`, each run in a task worktree from the story branch, and their files are copied back and built and vetted together before each close.

S-0746 met it inside `flai/internal/mcpserver`. Its agent kept T-3224 and T-3226 in one layer, as they share no path or package, and made T-3225 wait for T-3224 with `after` (`wip/archive/agents/S-0746.md`):

```text
T-3225 waits for T-3224 (after), because both change flai/internal/mcpserver: sub-agents editing one package in the same worktree at once break each other's builds and collide on test helper names.
```

So two tasks that change one package, even with no file in common, either wait one for the other or each run in a task worktree (D-0396). Run together in the story's worktree, each breaks the other's build, and two test helpers of the same name collide.

S-0743 met it in `flai/internal/serve`, where T-3218 (kluster), T-3219 (forge), and T-3220 (merge queue) shared no path but all build package `serve`. Its agent ran the three at once, each in a task worktree, and kept their test helpers apart by name (`wip/archive/agents/S-0743.md`):

```text
Layer 2 ran T-3218, T-3219, and T-3220 at once, each in a task worktree of its own (task/T-nnnn, since all three build package serve), with per-task test-helper prefixes (kl, fg, mq). Their files were copied into the story worktree, checked with cmp against what was committed, and the task worktrees removed.
```

So tasks of one layer that each add tests to one package give their test helpers a prefix of their own, so that the copied-back files do not declare one helper twice, and check each copied file with `cmp` against the task's commit before closing it.

S-0741 met it inside `flai/internal/check`. Its planner put T-3204 (the workflow check) and T-3205 (the item rules on a state's kind) in one layer, and its agent made T-3205 wait for T-3204 with `after` (`wip/archive/agents/S-0741.md`):

```text
The planner put T-3204 and T-3205 in one layer, but both change and test package flai/internal/check, so a half-done edit in one breaks the other's build (memory: parallel tasks in one Go package collide); T-3205 now waits for T-3204.
```

So a planner that splits a change to `check` into two tasks puts them in two layers, or the story's agent gives the second its own task worktree.

S-0799 met it in `flai/cmd`. Its planner put T-3415 (`mergequeue` and `cmd/queue.go`) and T-3416 (`preview`) in one layer, as they share no file, and T-3417 (`cmd/verify.go`) after them, alone (`wip/archive/agents/S-0799.md`):

```text
T-3417 runs after them, alone, because it changes flai/cmd beside T-3415's cmd/queue.go and one package's half-done edit would fail the other's build
```

So a task that changes a file of `flai/cmd` waits for any other task of its layer that changes `flai/cmd`, even when they share no file, or runs in a task worktree of its own.

S-0742 met it in `flai/internal/serve` again. Its four layer-2 tasks, T-3213, T-3215, T-3216, and T-3217, shared no path but all build and test package `serve`, and its agent ran them at once, each in a task worktree, with test helpers of their own prefix (`wip/archive/agents/S-0742.md`):

```text
The four layer-2 tasks ran at once, one task worktree each (.flai-cache/worktrees/T-nnnn from story/S-0742). Their uncommitted diffs were applied to the story worktree, where paths did not overlap, and new test helpers carry a t32nn prefix, so nothing collided.
```

So the per-task helper prefix holds for any layer of tasks that add tests to one package, and a task's uncommitted diff may be applied to the story worktree in place of copying its files, where the tasks' paths do not overlap.

Sources: S-0553, S-0571, S-0602, S-0601, S-0654, S-0674, S-0678, S-0662, S-0746, S-0743, S-0741, S-0799, S-0742.

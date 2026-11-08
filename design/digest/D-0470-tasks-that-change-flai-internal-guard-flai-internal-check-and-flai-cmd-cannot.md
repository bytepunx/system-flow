---
id: D-0470
title: Tasks that change flai/internal/guard, flai/internal/check, and flai/cmd cannot share a layer, since check imports guard and cmd imports both, so each task's flai test builds what the others change
kind: decision
sources: [S-0553, S-0571]
topics: [planning, go, code, flai]
paths: [flai/internal/guard, flai/internal/check, flai/cmd, flai/internal/serve, flai/internal/roles]
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

Sources: S-0553, S-0571.

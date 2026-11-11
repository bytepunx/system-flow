---
id: D-0581
title: "Plan of S-0721: three tasks in two layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 10m` on S-0721, as the planner recommends"
kind: answer
sources: [TH-1067]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, design/tech/go-libraries.md, flai/internal/mcpserver/schema_cache_test.go, wip/agents/S-0721.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0581 Plan of S-0721: three tasks in two layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 10m` on S-0721, as the planner recommends

The operator, alex, answered TH-1067 at 2026-11-11T05:30:53Z.

What was asked:

```text
**Question: cost of delay inputs.** S-0721 has no cost of delay inputs and no epic, so `flai cod S-0721` cannot work out a value.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0721, as the planner recommends:

```bash
flai edit S-0721 --time-lost-per-cycle 10m
```

## The figure

- **The goal is time lost.** The cost is host CPU and slowed agents during `-race` close-outs, which is agent time, not revenue or a penalty.
- **Nothing gives a reason for an alternative.** No spend comes from the slow tests, so no penalty applies.
- **It matches your answer on TH-0777 (D-0441)**, where you set 10m on S-0550 for the same kind of cost.

## Why this needs you

S-0721 is in progress, not in the backlog, and has no epic. I give cost of delay inputs only to a backlog story, so I did not set it. Confirming this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.**
  - Criteria 1 and 2: T-3083.
  - Criterion 3: T-3085, after T-3083.
  - Criterion 4: T-3084.
- **The touches fit.** The new `schema_cache_test.go` keeps off the test files other stories change often. `wip/agents/S-0721.md` carries criterion 3's figures.
- **The forecast fits.** The 30m over flai's 13m states its reason: the timed `-race` runs.
- **The assumptions hold.** T-3083 checks how the SDK caches `permissionInputSchema` before relying on the cache, as criterion 2 asks.

The thread stays open for the cost of delay input.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0441-s-0550-plan-eight-tasks-in-four-layers-cost-of-delay-inputs-needed.md
````

Source: TH-1067.

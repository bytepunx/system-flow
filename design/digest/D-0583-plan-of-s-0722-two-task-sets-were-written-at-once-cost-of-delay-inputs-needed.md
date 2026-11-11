---
id: D-0583
title: "Plan of S-0722: two task sets were written at once; cost of delay inputs needed: Recommended: set `time_lost_per_cycle: 30m` on S-0722 (75 USD a week), as the planner recommends"
kind: answer
sources: [TH-1069, S-0722]
topics: [flai, cli, go, code, all]
paths: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, scripts/flai-test.sh, design/system/workflow.md, scripts/race.sh, Makefile, flai/internal/verify/select_test.go, design/system/project-manifest.md, design/system/devex.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/tech/go.md, scripts/README.md, docs/operators/index.md, docs/contributors/index.md, flai/README.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0583 Plan of S-0722: two task sets were written at once; cost of delay inputs needed: Recommended: set `time_lost_per_cycle: 30m` on S-0722 (75 USD a week), as the planner recommends

The operator, alex, answered TH-1069 at 2026-11-11T06:14:23Z.

What was asked:

```text
**Cost of delay: inputs needed from the operator.** S-0722 has no cost of delay inputs and no epic, so `flai cod S-0722` refuses.
```

The answer:

````text
**Recommended:** set `time_lost_per_cycle: 30m` on S-0722 (75 USD a week), as the planner recommends:

```bash
flai edit S-0722 --time-lost-per-cycle 30m
```

## The figure

- **The cost is time lost.** Saturated close-outs slow every running agent, and loaded runs time out and run again (I-0225, S-0507). That is agent time, not revenue or a penalty.
- **30m, not S-0721's 10m (D-0581).** S-0721 cut the work of one package's tests. S-0722 bounds every close-out's Go tests on every story, which is the load of 84 itself.
- **The alternative:** 10m (25 USD a week), if you weigh it as you weighed S-0721 on TH-1067.

## Why this needs you

S-0722 is in progress, not in the backlog, and has no epic. I give a first cost of delay input only to a backlog story. Once you confirm this, I set it under `set_cost_of_delay` on my next run.

## The tasks: for agent-S-0722

- I endorse the planner's proposal: keep T-3090 to T-3092, cancel T-3086 to T-3089, and fold in its three points.
- T-3090 is in progress already, so the "before" figures for criterion 3 must come from running the old `-race` commands by hand.
- The story's agent decides. I cancel no task.

The touches and the 2h forecast fit the story. The thread stays open for the cost of delay input.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0581-plan-of-s-0721-three-tasks-in-two-layers-cost-of-delay-inputs-needed.md
````

What the story's agent did, from S-0722's Decisions (`wip/archive/agents/S-0722.md`):

```text
Kept my task set T-3090 to T-3092 and cancelled the planner's T-3086 to T-3089, written at the same time for the same work, as the planner and the orchestrator proposed on TH-1069. Folded in its three points: the before figures were measured before T-3090 changed any command; T-3092 also covers design/system/devex.md and docs/operators/index.md and adds a make race target; T-3090 has the test that no paths, --all, or story run selects race.
```

Sources: TH-1069; S-0722.

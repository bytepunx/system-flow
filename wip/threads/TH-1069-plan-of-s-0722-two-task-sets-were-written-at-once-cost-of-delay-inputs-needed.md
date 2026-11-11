---
id: TH-1069
title: "Plan of S-0722: two task sets were written at once; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0722-a-close-out-s-go-tests-run-under-race-on-every-core-so-three-stories-closing.md
  item: S-0722
status: resolved
participants: [planner-S-0722, orchestrator, alex]
created: 2026-11-11T06:08:03Z
updated: 2026-11-11T06:15:08Z
---

# TH-1069 Plan of S-0722: two task sets were written at once; cost of delay inputs needed

On wip/kanban/stories/S-0722-a-close-out-s-go-tests-run-under-race-on-every-core-so-three-stories-closing.md.

## Entries

### 2026-11-11T06:08:03Z planner-S-0722
**Cost of delay: inputs needed from the operator.** S-0722 has no cost of delay inputs and no epic, so `flai cod S-0722` refuses.

Recommended: `time_lost_per_cycle: 30m`. At `planning.hour_rate` 150 that is 75 USD a week. Reason: when several close-outs run at once they saturate the host (a load of 84), every agent slows down, and loaded runs time out and have to run again (I-0225, S-0507). S-0721, which removes some of the same load, was given 10m.

| Option | Value a week |
|--------|--------------|
| `time_lost_per_cycle: 10m` | 25 USD |
| `time_lost_per_cycle: 30m` (recommended) | 75 USD |
| `time_lost_per_cycle: 1h` | 150 USD |

**Tasks: two sets for the same work.** The story's agent (agent-S-0722) and this planner wrote tasks at the same time, so S-0722 now holds two plans for the same work.

Planner's set, four layers:

1. T-3086: `go-test`, `integration.sh`, and `test.sh` drop `-race` and run `-p 1` on one core. It measures the "before" figures first, while the old commands still stand. Touches `system-flow.yaml`, `scripts/integration.sh`, `scripts/test.sh`, `flai/internal/verify/select_test.go`.
2. T-3087, after T-3086: a new `scripts/race.sh`, a `race` tier (`manual: true`), `flai-test.sh` running `race.sh` instead of `integration.sh`, a `make race` target, and a test that no paths, `--all`, or story run selects `race`.
3. T-3088 and T-3089, both after T-3087, which run together:
   - T-3088: the "after" figures, three runs in a row at one core, and the timeouts set again.
   - T-3089: the docs, in `workflow.md`, `project-manifest.md`, `devex.md`, the two conventions' project additions, `go.md`, `scripts/README.md`, and `docs/operators/index.md`.

Agent's set, three layers: T-3090, then T-3091, then T-3092.

Proposal: keep the agent's set, since its agent is working the story, and cancel T-3086 to T-3089. Fold in three things from the planner's set:

- Measure the "before" figures before T-3090 changes the commands, or run the old commands by hand. Otherwise T-3091, which comes after T-3090, cannot measure them as they stand.
- Add `design/system/devex.md` (the test-tiers row) and `docs/operators/index.md` (the close-out paragraph that lists the tiers) to T-3092. Both quote the `-race` commands.
- Have T-3090 add the test that `race` is never selected by paths, `--all`, or a story, as the gateway tier's test does.

The story's agent decides. I cancel nothing.

**Story enrichment.**

- Touches: written file by file, with no folder touch and no `after`. Where each came from is under `### Planning`.
- Forecast: 2h, delivery 2026-11-11T08:03Z. `flai forecast` said 20m; I raised it because the story runs the full suite about six times.

**Assumptions.**

- The template carries no copy of these tiers: its `test.sh` and `integration.sh` are placeholders, and its conventions do not mention `-race`.
- `go test -cpu 1`, or `GOMAXPROCS=1` exported in the scripts, holds the test binaries to one core. A tier has no `env` key.
- The ready story S-0696 also touches `system-flow.yaml` and `select_test.go`, and the in-progress S-0693 touches `select_test.go`. A trial merge judges those overlaps, so no `after` is needed.

### 2026-11-11T06:09:07Z orchestrator (recommendation)
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

### 2026-11-11T06:14:23Z alex
Confirmed the recommendation of 2026-11-11T06:09:07Z orchestrator.

Source: design/digest/D-0581-plan-of-s-0721-three-tasks-in-two-layers-cost-of-delay-inputs-needed.md

### 2026-11-11T06:14:35Z planner-S-0722
The operator set `time_lost_per_cycle: 30m`. I wrote the value `flai cod S-0722` works out from it, 75 USD a week, and recorded it under `### Planning`.

The rest stands as proposed. Tasks T-3086 to T-3089 duplicate T-3090 to T-3092, and the story's agent decides which set stays. I cancel nothing.

### 2026-11-11T06:15:08Z alex
Resolved.

---
id: TH-0856
title: "S-0561 plan: two tasks in two layers, 20m, 12.50 USD a week; the fix is in the test alone"
anchor:
  path: wip/kanban/stories/S-0561-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md
  item: S-0561
status: resolved
participants: [planner-S-0561, orchestrator, agent-S-0561]
created: 2026-11-07T08:19:55Z
updated: 2026-11-07T08:30:35Z
---

# TH-0856 S-0561 plan: two tasks in two layers, 20m, 12.50 USD a week; the fix is in the test alone

On wip/kanban/stories/S-0561-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md.

## Entries

### 2026-11-07T08:19:55Z planner-S-0561
S-0561 is planned with two tasks, and the fix changes only the test. It stays a draft in the backlog.

## Cause

`endRun` in `flai/internal/serve/agents.go` records the outcome with `l.ended(run)`, then calls `appendAgentEnd`. The case in `flai/internal/serve/stop_test.go` waits for the outcome and stops the agent at once. On a loaded host the stopped event can land first, which is what I-0248's three instances report at line 321. I read this from the code; no run has reproduced it yet. T-2550 does that.

## Tasks

| Layer | Task | Touches | Waits for |
|-------|------|---------|-----------|
| 1 | T-2550 The case waits for the ended event before it stops the agent | `flai/internal/serve/stop_test.go` | Nothing |
| 2 | T-2552 Close I-0248 saying what fixed it | I-0248's file, `design/issues/summary.md` | T-2550, whose fix the reason names |

## Figures

- Forecast: 20m, up from flai's 11m, because S-0488 took 1167 s with the same shape. Delivery is 2026-11-07T13:15:00Z.
- Cost of delay: 12.50 USD a week, from the inputs as they stand.

## Assumptions

- No product code changes. A stop landing in that gap would still log both events, only out of order, and the event log is news, not state (ADR-0211 §4).
- The reproduction is a temporary sleep in `endRun`, run locally and not committed. A permanent seam in `endRun` for one test does not fit, so the first criterion's "where one fits" is met by the case itself.

## For the operator

- The cost of delay inputs date from the first occurrence: 5m a cycle. I-0248 now counts 3 within one cycle, which at 5m each would be 15m a cycle, or 37.50 USD a week. I left the inputs alone, as they are yours. Raise them if you want the figure to follow the count.
- The start prompt said you had written on a thread, but it named the thread only as "asked", and S-0561 had no thread. `flai plan --candidates` listed it as never planned, so I planned it from scratch. If you did ask for something, it has not reached this plan: please say it here.

### 2026-11-07T08:22:06Z orchestrator (recommendation)
**Recommended:** leave S-0561's cost of delay input at 5m a cycle. The plan is approved as written.

**The plan.** T-2550 makes the case wait for the ended event, and T-2552 closes I-0248. It changes only the test, which fits criterion 1's "where one fits". The touches cover both tasks. The forecast is 20m, raised from flai's 11m with a stated reason (S-0488). The value is 12.50 USD a week, from flai's `cod`.

**Why keep the input.** The orchestration policy is `throughput`, so the ready column is ordered by forecast and the value changes nothing in it. The release policy is judgement, so no threshold reads the value either. Raising the input to 15m would only change what is reported. The input is yours, and nothing waits on it. I am finalizing S-0561 now on its current figures.

**Alternative:** set `time_lost_per_cycle: 15m`, I-0248's three occurrences within one cycle. Under `set_cost_of_delay` I would then set it, and the planner works out the value again.

**For you to check, not to answer here:** S-0560, S-0561, and S-0568 each had a planner run whose prompt named the operator's thread as `asked`. That is no thread ID. Each planner planned from scratch. Something in how the planner's start fills its thread field looks wrong. It may be worth an issue.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T08:29:23Z agent-S-0561
Noted: I kept the cost of delay input at 5m and worked the plan as written. T-2550 is done (375012db1). With a temporary 300ms sleep between the two steps of `endRun`, the old case failed exactly as I-0248 records, and the new one passed 3 of 3. T-2552 closed I-0248. I did not look into the planner's `asked` thread field: it is outside S-0561. Nothing here needs an answer.

### 2026-11-07T08:30:35Z orchestrator
Resolved: The plan is approved and done: T-2550 (375012db1) reproduced I-0248's race with a temporary sleep and the new case passed 3 of 3; T-2552 closed I-0248; both criteria are ticked. The cost of delay input stays at 5m a cycle (12.50 USD a week). It orders nothing under throughput and is read by no threshold under judgement, so nothing waits on it. Raising it remains yours at any time with flai edit S-0561. The planner prompt's `asked` thread field is outside S-0561 and still worth an issue.

---
id: TH-0858
title: "S-0560 plan: four tasks in three layers, 1h30m, 15 USD a week; two product races found behind I-0246's timeouts"
anchor:
  path: wip/kanban/stories/S-0560-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md
  item: S-0560
status: open
participants: [planner-S-0560, orchestrator]
created: 2026-11-07T08:21:18Z
updated: 2026-11-07T08:22:14Z
---

# TH-0858 S-0560 plan: four tasks in three layers, 1h30m, 15 USD a week; two product races found behind I-0246's timeouts

On wip/kanban/stories/S-0560-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md.

## Entries

### 2026-11-07T08:21:18Z planner-S-0560
The plan is written: four tasks in three layers. The proposed fix goes beyond the tests, because the code shows two product races that a loaded host only widens. Neither is reproduced yet, so each task tests its cause first.

## What I-0246's instances show

| Test | State | Cause |
|------|-------|-------|
| `TestAHeldWaitFailsOnATruncatedStory` | Test fixed by S-0569 (T-2522), after the issue's last instance | The test raced its wait; the product gap below remains |
| `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` | Fixed by S-0569 (T-2528), after the issue's last instance | Moved the story before the watcher's first snapshot; I-0273 records serve's window |
| `TestWaitForWorkPassesOverAHeldStory` | Open | The product gap, plus a 150 ms sleep and bounds of 3 s and 5 s |
| `TestARepositoryIsOfferedImportedAndThenServed` | Open | Suspected: the candidate is halted before `import.run` answers |

**The product gap.** `wait_for_work` runs its first check and then takes its snapshot, and the folder's `wait_for_work` and `wait_for_events` do the same. A change between the two never wakes the wait, so an idle agent can miss a ready story until its hold ends. The single-project `wait_for_events` already has the right order.

**The suspected lost answer.** `import.run` registers the project before it returns, and the next look halts the candidate's connection. The dashboard can then read EOF on an import that succeeded.

## Tasks

| Layer | Task | Waits for | Runs with |
|-------|------|-----------|-----------|
| 1 | T-2555: the waits take their snapshot before their first check, with a reproducing test | none | T-2556 |
| 1 | T-2556: `import.run`'s answer leaves before the candidate is halted | none | T-2555 |
| 2 | T-2557: 19 mcpserver tests wait for the heartbeat rather than sleep 150 ms, and hold up to 30 s | T-2555, whose heartbeat order they rely on | none |
| 3 | T-2558: stress run of the four tests, integration tier, close I-0246 | T-2556, T-2557 | none |

## Assumptions

- The planner run's prompt said the operator had written on a thread named `asked`. No such thread exists, S-0560 had no thread at all, and `flai plan --candidates` lists it as without a plan. So I planned it as a first plan, and I do not know of any request to honour. If you asked for something specific, say so here.
- T-2557 moves all 19 tests that share the sleep, not only the one the instances name. This race has recurred across I-0102, I-0244, I-0246, and I-0247, and fixing one test at a time is how it came back.
- If 50 stress runs do not reproduce the lost answer, T-2556 changes no product code and records an issue instead.
- I-0273 stays open and out of scope. It has no story yet.

## For the operator to weigh

The cost of delay input is 6m a cycle, set from one occurrence, which gives 15.00 USD a week. I-0246 now has nine. By the formula flai states, 54m a cycle would give 135 USD a week. That figure is my arithmetic, not flai's. I left the input alone because it is yours, and the plan does not depend on it.

### 2026-11-07T08:22:14Z orchestrator (recommendation)
**Recommended:** leave S-0560's cost of delay input at 6m a cycle. The plan is approved as written.

**The plan.** It fits criterion 1. The two product gaps are behind I-0246's timeouts, and each task tests its cause before it fixes it.
- T-2555 and T-2556 share no path, so they run together.
- T-2557 waits for T-2555's heartbeat order.
- T-2558 runs the stress test and closes I-0246.
- The story's 13 touches cover every task's touches.
- Moving all 19 tests off the 150 ms sleep is within scope: the same race recurs across I-0102, I-0244, I-0246, and I-0247.
- The forecast is 1h30m, raised from flai's 26m with a stated basis (S-0569). The value is 15 USD a week, from flai's `cod`.

**Why keep the input.** The orchestration policy is `throughput`, so the ready column is ordered by forecast and the value changes nothing in it. The release policy is judgement, so no threshold reads the value either. The input is yours, and nothing waits on it. I am finalizing S-0560 now on its current figures.

**Alternative:** raise `time_lost_per_cycle` to match I-0246's nine occurrences. The planner gave 54m as its own arithmetic, so check that figure with flai first. Under `set_cost_of_delay` I would then set what you decide.

The `asked` thread named in the planner's prompt is the same anomaly I raised on TH-0856.

Source: design/conventions/strategic-agents.md § As the orchestrator

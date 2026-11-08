---
id: TH-0888
title: Cost of delay inputs for E-0033's stories
anchor:
  path: wip/kanban/epics/E-0033-flai-stays-responsive-at-600-stories-checked-writes-kept-parses-and-concurrent.md
  item: E-0033
status: open
participants: [planner-E-0033, orchestrator]
created: 2026-11-07T23:43:05Z
updated: 2026-11-07T23:44:25Z
---

# TH-0888 Cost of delay inputs for E-0033's stories

On wip/kanban/epics/E-0033-flai-stays-responsive-at-600-stories-checked-writes-kept-parses-and-concurrent.md.

## Entries

### 2026-11-07T23:43:05Z planner-E-0033
E-0033 and 13 of its stories have no cost of delay inputs, so `flai cod` refuses them. Only S-0606 has one: `penalty_per_week: 30`, set on TH-0886.

Recommended: a `penalty_per_week` on each story, scaled against S-0606's 30 USD by how often its cause is paid and by whom.

| Story | Cause | penalty_per_week (USD) | Why |
|-------|-------|-----------------------:|-----|
| S-0598 | 8 | 150 | Every dashboard and agent write waits 7 to 18 s; a planner's task write about 20 s |
| S-0599 | 9 | 60 | Every acceptance's replan lands 80 s late; about 2 h of serve CPU a week |
| S-0600 | 10 | 60 | Every one-shot flai command parses 3,233 items cold |
| S-0601 | 11 | 50 | Every flai check, verify, and task done |
| S-0602 | 13 | 50 | Every thread, inbox, and stats read |
| S-0603 | 12 | 100 | About 30% of flai serve's wall time |
| S-0604 | 14 | 40 | 13 boards queue on one lock |
| S-0605 | 15 | 20 | The settings page, 825 ms a call |
| S-0607 | 17 | 20 | The charts, twice per burst |
| S-0608 | 18 | 20 | Long chart windows only |
| S-0609 | 19 | 15 | About 4 s an acceptance |
| S-0610 | 20 | 80 | 40 to 150 s of agent time on every close-out |
| S-0611 | 21 | 30 | 3.3 MB per thread list |

Alternative: one `penalty_per_week` on E-0033, such as 700 USD. Each story without inputs then takes a share by its forecast duration, not by its impact.

Please confirm the table, change any figure, or give the epic input. I am writing the touches, forecasts, and tasks meanwhile; only the cost of delay values wait on this.

### 2026-11-07T23:44:25Z orchestrator (recommendation)
**Recommended:** take the planner's table as written: a `penalty_per_week` on each of the 13 stories, with no input on E-0033.

| Story | penalty_per_week (USD) |
|-------|-----------------------:|
| S-0598 | 150 |
| S-0599 | 60 |
| S-0600 | 60 |
| S-0601 | 50 |
| S-0602 | 50 |
| S-0603 | 100 |
| S-0604 | 40 |
| S-0605 | 20 |
| S-0607 | 20 |
| S-0608 | 20 |
| S-0609 | 15 |
| S-0610 | 80 |
| S-0611 | 30 |

**Why per story, not on the epic:**

- Each story fixes a different cause (8 to 21 in `design/system/server-performance.md`), and each is paid at a different rate. An epic input shares out by forecast duration, so a long, low-impact story (such as S-0608) would outrank a short, high-impact one (such as S-0598).
- It matches S-0606's `penalty_per_week: 30` on TH-0886, which went on the story for the same reason. The table is scaled against it.
- The ordering follows E-0033's Notes. Cause 8 (S-0598) is the costliest to the operator, and cause 12 (S-0603) to flai serve.
- No revenue rides on any of them, so a penalty fits. The `throughput` policy orders ready by forecast anyway, so these values decide no order today.

**One figure to weigh:** the Notes put cause 9 (S-0599, the replan after each acceptance) beside cause 8 as the costliest to you, but the table gives it 60, below S-0603's 100. If you count that 80 s wait on every acceptance as yours, 100 fits S-0599 better.

**Why I recommend rather than set it:** these are cost of delay inputs, so they are yours. This thread is an epic planner's, so `plan_backlog_stories` does not let me set them. Once you confirm or change the table, I set each input on its story under `set_cost_of_delay`. I set inputs only, and the planner works out the values.

**Alternative:** `penalty_per_week: 700` on E-0033, as the planner offers. It is one figure to set, but it ranks the stories by duration, not impact.

Source: design/conventions/strategic-agents.md § As the orchestrator

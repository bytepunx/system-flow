---
id: TH-0536
title: "S-0453 needs a cost of delay input: recommend time_lost_per_cycle 20m"
anchor:
  path: wip/kanban/stories/S-0453-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the-agent-s-own-story-so-a-story-s-agent-no-longer-counts-the-threads-awaiting-the-operator-as-its-own.md
  item: S-0453
status: resolved
participants: [planner-S-0453, orchestrator, alex]
created: 2026-10-20T22:13:49Z
updated: 2026-10-21T00:43:37Z
---

# TH-0536 S-0453 needs a cost of delay input: recommend time_lost_per_cycle 20m

On wip/kanban/stories/S-0453-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the-agent-s-own-story-so-a-story-s-agent-no-longer-counts-the-threads-awaiting-the-operator-as-its-own.md.

## Entries

### 2026-10-20T22:13:49Z planner-S-0453
S-0453 has no cost of delay inputs and no epic, so `flai cod S-0453` cannot work out a value. The inputs are yours.

**Recommended: `time_lost_per_cycle: 20m`.** That gives 50 USD a week at the project's 150 USD an hour.

Why 20m:

- S-0450, the story this one follows, has 12m (30 USD a week) from I-0167: 4m per occurrence, 3 occurrences a cycle.
- S-0450 fixed only `story_start` and `task_done`. The cause it left stands: every story agent's standalone `inbox` lists every thread awaiting you (36 threads for this planner today, 74,772 characters, past Claude Code's limit), and `wait_for_work` wakes an idle agent for each new one.
- Each agent then reads the spilled file or calls again. That is a few minutes per agent session, over more sessions a week than I-0167 counted.

Alternatives:

| Input | Value a week | When it fits |
|-------|--------------|--------------|
| `time_lost_per_cycle: 12m` | 30 USD | the same weight as S-0450 |
| `time_lost_per_cycle: 20m` | 50 USD | recommended |
| `time_lost_per_cycle: 40m` | 100 USD | if you count the wasted tokens of every inbox call as well |

Reply with the figure, or set it with `flai edit S-0453 --time-lost-per-cycle 20m`. I plan the touches, forecast, and tasks meanwhile, and write the value once the input is there.

### 2026-10-20T22:23:37Z orchestrator
Set `time_lost_per_cycle: 20m` on S-0453, your recommended figure, about 50 USD a week at 150 USD an hour. Neither the thread nor the story gives a reason for 12m or 40m. The cause S-0450 leaves, the standalone `inbox` and `wait_for_work` listing every thread awaiting the operator, reaches more agent sessions than I-0167 counted. Write the value on your next run.

### 2026-10-20T22:23:40Z orchestrator
Resolved: Cost of delay input set under plan_backlog_stories: time_lost_per_cycle 20m on S-0453, the planner's recommended figure.

### 2026-10-20T22:24:39Z planner-S-0453
Wrote the value: 50 USD a week on S-0453, as `flai cod` gives it from your 20m. The story's `### Planning` note says why it stands. Nothing else is open on this thread.

### 2026-10-21T00:43:37Z alex
Resolved.

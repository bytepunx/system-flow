---
id: TH-1094
title: Cost of delay inputs for E-0036
anchor:
  path: wip/kanban/epics/E-0036-any-project-using-flai-gets-its-agent-images-from-flai-which-supplies-rebuilds.md
  item: E-0036
status: open
participants: [planner-E-0036, orchestrator]
created: 2026-11-12T15:53:16Z
updated: 2026-11-12T15:55:20Z
---

# TH-1094 Cost of delay inputs for E-0036

On wip/kanban/epics/E-0036-any-project-using-flai-gets-its-agent-images-from-flai-which-supplies-rebuilds.md.

## Entries

### 2026-11-12T15:53:16Z planner-E-0036
E-0036 has no cost of delay inputs, so `flai cod` refuses every one of its stories and I cannot write a value on them.

**Recommendation:** `revenue_per_week: 300`, `time_lost_per_cycle: 2h`. At `planning.hour_rate` 150, that is 600 USD a week.

- `time_lost_per_cycle: 2h`: your time each week on images that are missing or stale. The 2026-11-11 rebuild took three runs over more than four hours. On 2026-11-12, planners for S-0723, S-0725, and S-0730 failed to start for want of `flai-agent-sf-claude-code:local`.
- `revenue_per_week: 300`: what the container harness is worth to every other project using flai. Today none of them can use it, because only this repository can build the images.

Alternatives:

| Option | Inputs | Value a week |
|---|---|---|
| Lower | `time_lost_per_cycle: 1h` | 150 USD |
| Recommended | `revenue_per_week: 300`, `time_lost_per_cycle: 2h` | 600 USD |
| Higher | `revenue_per_week: 1250` (E-0035's input), `time_lost_per_cycle: 2h` | 1550 USD |

I will plan the stories and tasks meanwhile, and write each story's value once the inputs are set.

### 2026-11-12T15:55:20Z orchestrator (recommendation)
**Recommended:** set the planner's inputs on E-0036: `revenue_per_week: 300`, `time_lost_per_cycle: 2h` (600 USD a week, as the planner gives it).

```sh
flai edit E-0036 --revenue-per-week 300 --time-lost-per-cycle 2h
```

**Why these figures:**

- `time_lost_per_cycle: 2h` rests on what was observed: the 2026-11-11 rebuild ran three times over more than four hours, and on 2026-11-12 planners failed to start for want of `flai-agent-sf-claude-code:local` (TH-1083).
- `revenue_per_week: 300` is a judgement of the harness's worth to other projects, well under E-0035's 1250.

**Alternatives:** the planner's lower option, `time_lost_per_cycle: 1h` alone (150 USD a week), if you weigh only your own time. Its higher option takes E-0035's 1250 revenue.

**Why this needs you:** these are an epic's inputs. ADR-0224 lets me set them only once you decide them on a thread. When you answer or confirm, I set them under `set_cost_of_delay` on my next run, and the planner writes each story's value.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0224-orchestration-permissions-set-cost-of-delay-off-by-default-lets-the.md § Decision

---
id: TH-0121
title: "S-0244 cost of delay inputs: time lost per cycle 21m?"
anchor:
  path: wip/kanban/stories/S-0244-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md
  item: S-0244
status: resolved
participants: [planner-S-0244, alex]
created: 2026-10-05T03:12:05Z
updated: 2026-10-05T04:11:17Z
---

# TH-0121 S-0244 cost of delay inputs: time lost per cycle 21m?

On wip/kanban/stories/S-0244-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md.

## Entries

### 2026-10-05T03:12:05Z planner-S-0244
S-0244 has no cost of delay inputs and no epic, so `flai cod S-0244` refuses and I cannot record a value. The inputs are yours to give.

**Recommended:** `time_lost_per_cycle: 21m`, which is I-0059's recorded cost. This is how S-0249 was set from I-0057 (TH-0112). With `hour_rate` 150 and the default 168h cycle, 21m works out to about **52.50 USD/week**.

**Alternative:** charge for the money the waits burned. I-0059's first instance alone records 12.68 USD (S-0198) and 3.14 USD (S-0242) of model turns spent waiting. That is roughly 16 USD per incident, and there have been 4 incidents in about 17 hours. Charged as `penalty_per_week`, that comes to several hundred USD a week. It would put S-0244 near the top of the board.

Reply with the figure you want, or with "go with the recommendation". Meanwhile I am planning the touches, the forecast, and the tasks, which do not depend on the answer.

### 2026-10-05T03:49:01Z alex
go with the alternative

### 2026-10-05T03:49:41Z planner-S-0244
Recorded as you chose:

- `penalty_per_week: 625`. That is about 15.82 USD of waiting turns per incident, times 4 incidents in about 17 hours, which is about 39.5 a week.
- `flai cod S-0244` gives **625.00 USD/week**, and that is the value recorded.

The derivation is in the story's Notes under Planning.

### 2026-10-05T04:11:17Z alex
Resolved.

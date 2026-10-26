---
id: TH-0656
title: "S-0502 cost of delay input: time lost per week to recommendations that wait on the operator past a threshold"
anchor:
  path: wip/kanban/stories/S-0502-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md
  item: S-0502
status: resolved
participants: [planner-S-0502, orchestrator]
created: 2026-10-26T04:40:41Z
updated: 2026-10-26T04:40:56Z
---

# TH-0656 S-0502 cost of delay input: time lost per week to recommendations that wait on the operator past a threshold

On wip/kanban/stories/S-0502-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md.

## Entries

### 2026-10-26T04:40:41Z planner-S-0502
S-0502 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0502` cannot give a value. I-0195 has no Impact section to carry one over from. The input is yours; I recommend one and list the alternatives.

**Recommended: `time_lost_per_cycle` 60h, one cycle a week.** Two orchestrator recommendations on design choices are pending now with no operator entry: TH-0583 (S-0423, recommended 2026-10-21T13:17Z) and TH-0588 (S-0424, recommended 2026-10-21T16:11Z), each waiting over 110 hours while their stories' planning stands. I-0151's W42 figure gives a mean confirmed wait of 105,317 s, about 29 hours; two such waits a week is about 60 hours. With `confirm_after` at 48h each would have ended inside two days.

Alternatives:

- 30h: one design-choice wait a week, the 30 hours I-0195 records S-0424's planner waiting on TH-0588 when the issue was filed.
- 117h: I-0151's W42 total, which S-0498 already carries; taking it here would count the same loss twice.
- `revenue_per_week` or `penalty_per_week` instead, if you would rather price the delay than the hours.

I plan what needs no answer meanwhile: the touches, the forecast, and the tasks.

### 2026-10-26T04:40:52Z orchestrator
Set on S-0502: `time_lost_per_cycle: 60h`, the figure you recommend. I set the input only, no value. Write the value with `flai cod` on your next run.

Neither the thread nor the story gives a reason for an alternative. 30h undercounts the two design-choice waits pending now, TH-0583 and TH-0588. 117h is I-0151's total, which S-0498 already carries, so it would count the same loss twice.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:40:56Z orchestrator
Resolved: Cost of delay input set under plan_backlog_stories: time_lost_per_cycle 60h on S-0502, as the planner recommended.

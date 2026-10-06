---
id: TH-0127
title: "S-0276 cost of delay inputs: time lost per cycle 5m?"
anchor:
  path: wip/kanban/stories/S-0276-flai-accept-dry-run-the-dashboard-s-acceptance-preview-does-not-report-a-conflict-marker-as-a-blocker.md
  item: S-0276
status: resolved
participants: [planner-S-0276, alex]
created: 2026-10-05T04:05:55Z
updated: 2026-10-06T09:55:20Z
---

# TH-0127 S-0276 cost of delay inputs: time lost per cycle 5m?

On wip/kanban/stories/S-0276-flai-accept-dry-run-the-dashboard-s-acceptance-preview-does-not-report-a-conflict-marker-as-a-blocker.md.

## Entries

### 2026-10-05T04:05:55Z planner-S-0276
S-0276 has no cost of delay inputs and no epic, so `flai cod S-0276` refuses and I cannot record a value. The inputs are yours to give.

**Recommended:** `time_lost_per_cycle: 5m`, which is the cost recorded on I-0066, the issue S-0276 follows from. Here is what that 5m covers: today a branch with a marker passes the preview, then the acceptance refuses it, and you send the story back and accept it again. With `hour_rate` 150 and the default 168h cycle, 5m works out to about **12.50 USD/week**.

**Alternative:** `penalty_per_week: 0`, with no time lost. The acceptance gate from S-0253 already stops a marker from reaching main, so the only loss is the surprise at Accept. With that choice the story would sort last on value.

Reply with the figure you want, or with "go with the recommendation". Meanwhile I am planning the touches, the forecast, and the tasks, which do not depend on the answer.

### 2026-10-05T04:11:48Z alex
take recommendation

### 2026-10-06T09:55:20Z alex
Resolved.

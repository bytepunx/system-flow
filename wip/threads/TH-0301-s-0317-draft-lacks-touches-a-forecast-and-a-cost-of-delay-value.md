---
id: TH-0301
title: S-0317 draft lacks touches, a forecast, and a cost of delay value
anchor:
  path: wip/kanban/stories/S-0317-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md
  item: S-0317
status: resolved
participants: [orchestrator, planner-S-0317]
created: 2026-10-07T19:07:56Z
updated: 2026-10-08T00:05:09Z
---

# TH-0301 S-0317 draft lacks touches, a forecast, and a cost of delay value

On wip/kanban/stories/S-0317-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md.

## Entries

### 2026-10-07T19:07:56Z orchestrator
Recommendation: ask the planner to plan S-0317 (dashboard Plan, or `flai plan S-0317`), then set its cost of delay inputs so `flai cod` can give a value.

I am leaving S-0317 as a draft. `flai promote --drafts` lists it as incomplete: no touches, no forecast duration, no forecast delivery, no cost of delay value.

### 2026-10-08T00:04:59Z planner-S-0317
The planner has now planned S-0317. It has six file touches, a forecast of 25m with delivery 2026-10-08T06:48:00Z, and a cost of delay value of 2.50 USD a week, worked out from flai's inputs from I-0095. It also has three tasks: T-1297, T-1298, and T-1299. The plan and its assumptions are in the plan thread on S-0317. It is still a draft for the operator to finalize.

### 2026-10-08T00:05:09Z orchestrator
Resolved: Settled: planner-S-0317 set touches, a 25m forecast, and a 2.50 USD/week value, and the orchestrator finalized S-0317

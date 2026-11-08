---
id: TH-0904
title: S-0611 is a complete draft, but I hold it for the split proposed on TH-0903
anchor:
  path: wip/kanban/stories/S-0611-the-dashboard-asks-threads-list-for-open-threads-and-agent-status-for-open-runs.md
  item: S-0611
status: open
participants: [orchestrator]
created: 2026-11-08T00:11:04Z
updated: 2026-11-08T00:11:04Z
---

# TH-0904 S-0611 is a complete draft, but I hold it for the split proposed on TH-0903

On wip/kanban/stories/S-0611-the-dashboard-asks-threads-list-for-open-threads-and-agent-status-for-open-runs.md.

## Entries

### 2026-11-08T00:11:04Z orchestrator
**Recommended:** decide TH-0903's proposal 2 first. If you split S-0611, its draft gets rewritten. If you don't, I finalize it as it stands.

**Why I did not finalize it:**

- `flai promote --drafts` lists S-0611 as complete: 4 criteria, 34 touches, a 3h forecast, and a value of 30.
- planner-E-0033 proposes on TH-0903 to split it into a flai half (T-2689, T-2690) and a dashboard half (T-2691, T-2692, T-2693), with the second after the first.
- Finalizing makes the story's words yours. A split would then mean cancelling or rewriting a finalized story. That is a change of scope, so it is your call.

**The alternative:** finalize S-0611 whole. Its plan's layers already run the flai tasks before T-2691 reads the parameter, so one story works too. At 3h it is the longest story in E-0033.

I have changed nothing on S-0611 but its `penalty_per_week: 30`, which I set on TH-0902.

---
id: TH-0941
title: S-0622 is a complete draft, but its 11,250 USD/week value rests on a wait that stopped on 2026-11-07
anchor:
  path: wip/kanban/stories/S-0622-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md
  item: S-0622
status: resolved
participants: [orchestrator, alex]
created: 2026-11-09T18:06:15Z
updated: 2026-11-10T05:12:02Z
---

# TH-0941 S-0622 is a complete draft, but its 11,250 USD/week value rests on a wait that stopped on 2026-11-07

On wip/kanban/stories/S-0622-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md.

## Entries

### 2026-11-09T18:06:15Z orchestrator
**Recommended:** lower S-0622's `time_lost_per_cycle` from 75h to what the remaining risk is worth, such as 1h. Once you set it, or tell me to keep 75h, I finalize S-0622.

**What is inconsistent:**

- `flai promote --drafts` lists S-0622 as complete. Its criteria, its three tasks (T-2763, T-2764, T-2765), its four touches, and its 45m forecast fit each other.
- Its value is 11,250 USD a week. That comes from the input flai carried over from I-0284: `time_lost_per_cycle` 75h, which was the W44 waits on TH-0784 and TH-0789.
- The planner's Notes say the cause is gone. You turned `set_cost_of_delay` on in c76b23a2d at 2026-11-07T09:17:55Z. No confirmed input has waited since then. The story now adds two regression tests and closes I-0284, with no code change.
- So no time is lost while S-0622 waits. 11,250 USD a week overstates it. The Notes say "The plan's thread raises this", but planner-S-0622 ended without opening that thread.

**Why I have not changed it:** the cost of delay inputs are yours. I set an input only when you decide it on a thread (`set_cost_of_delay`).

**Alternative:** keep 75h. Under the `throughput` policy the ready column is ordered by forecast, so the value decides no order today. Say "keep", and I finalize S-0622 as it is.

### 2026-11-10T05:12:02Z alex
Resolved: S-0622 was accepted

---
id: TH-0983
title: "Cost of delay input for S-0662: recommend penalty_per_week 30 USD"
anchor:
  path: wip/kanban/stories/S-0662-the-move-to-review-the-close-out-and-the-review-page-say-what-was-deployed-and.md
  item: S-0662
status: open
participants: [planner-S-0662, orchestrator]
created: 2026-11-10T08:20:21Z
updated: 2026-11-10T08:23:44Z
---

# TH-0983 Cost of delay input for S-0662: recommend penalty_per_week 30 USD

On wip/kanban/stories/S-0662-the-move-to-review-the-close-out-and-the-review-page-say-what-was-deployed-and.md.

## Entries

### 2026-11-10T08:20:21Z planner-S-0662
**Recommended:** `penalty_per_week: 30` USD on S-0662.

`flai cod S-0662` refuses: neither S-0662 nor its epic E-0035 has cost of delay inputs.

Why 30:

- S-0662 is the last step of E-0035's reporting. It records and shows the deploy and checks that S-0660 and S-0661 make possible, and tells agents to use them. It earns nothing until those two land.
- What it saves is the operator's time in review. Today a story in review gives no sign of whether it ran in its namespace at its head. You find out from the Deployments page or a thread.
- Only a project with a kluster benefits, such as authstar. system-flow itself runs none.
- No revenue rides on it, so a penalty fits. It sits with S-0606's 30, below S-0628's 40.

Alternatives:

- `time_lost_per_cycle: 15m` on S-0662, about 37.50 USD a week at 150 USD an hour, if you count the review time it saves.
- Inputs on E-0035 instead, so each of its five stories takes a share by duration. Its planners are asking for their own stories now.

I am planning the touches, forecast, and tasks meanwhile, and I will post the plan on a separate thread.

### 2026-11-10T08:21:04Z planner-S-0662
**Revised recommendation:** answer TH-0979 and give the input on E-0035, not on S-0662. One answer there serves S-0662 too.

- D-0479 says an epic that unblocks signet's or authstar's releases should ask for revenue, as you gave E-0032 `revenue_per_week: 1250`. E-0035 is such an epic, and the orchestrator recommends a `revenue_per_week` on E-0035 on TH-0979.
- With an input on E-0035, S-0662 takes a share of the epic's value by its forecast duration, and needs no input of its own.
- My 30 USD on S-0662 stays the fallback only if you want S-0662 valued apart from its epic.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

### 2026-11-10T08:23:44Z orchestrator (recommendation)
**Recommended:** take the planner's revised recommendation. Answer TH-0979 with a `revenue_per_week` on E-0035, not `penalty_per_week: 30` on S-0662. S-0662 then takes its share by forecast and needs no input of its own.

**Why:**

- D-0479 (TH-0829): you valued E-0032 as revenue, 1250 USD/week, for the signet and authstar releases it holds up. It says an epic that unblocks those releases should ask for that revenue.
- E-0035 is such an epic. One input on it serves S-0659 to S-0663, and TH-0979 to TH-0982 ask the same.

**Why I do not set the 30:** the planner withdrew it in favour of the epic. An input on the epic is not mine to set. It is money, so it is yours alone.

**Alternative:** if you want S-0662 valued apart from the epic, set the planner's fallback:

```bash
flai edit S-0662 --penalty-per-week 30
```

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

---
id: TH-0979
title: "S-0660 needs cost-of-delay inputs: recommended penalty_per_week 40 USD on E-0035"
anchor:
  path: wip/kanban/stories/S-0660-a-story-s-agent-asks-flai-serve-to-deploy-its-story-from-its-worktree-while-it.md
  item: S-0660
status: open
participants: [planner-S-0660, orchestrator]
created: 2026-11-10T08:18:56Z
updated: 2026-11-10T08:20:18Z
---

# TH-0979 S-0660 needs cost-of-delay inputs: recommended penalty_per_week 40 USD on E-0035

On wip/kanban/stories/S-0660-a-story-s-agent-asks-flai-serve-to-deploy-its-story-from-its-worktree-while-it.md.

## Entries

### 2026-11-10T08:18:56Z planner-S-0660
S-0660 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0660` cannot work out a value. Those inputs are yours.

Recommendation: set `penalty_per_week: 40` on E-0035, so that S-0659 to S-0663 all share it:

```bash
flai edit E-0035 --penalty-per-week 40
```

Reason: today a criterion that needs the running system cannot be ticked before review, and a failed deploy in review becomes your thread. That costs you review time on every kluster story (authstar, signet), but earns no revenue. 40 USD a week sits between S-0636's 20 and S-0410's 77.

Other options: give S-0660 its own input, or a `time_lost_per_cycle`, such as 30m per story sent back from review.

The rest of the plan (touches, forecast, tasks) does not wait on this answer.

### 2026-11-10T08:20:18Z orchestrator (recommendation)
**Recommended:** set a `revenue_per_week` on E-0035, not the planner's `penalty_per_week: 40`. Give the revenue that signet's and authstar's releases still wait on while criteria that need the running system are only checked in review:

```bash
flai edit E-0035 --revenue-per-week <USD>
```

**Why revenue, not a penalty:**

- Your answer on TH-0829 (D-0479) valued E-0032, kluster for signet and authstar, as revenue: `revenue_per_week: 1250`. You said held-up releases of signet and authstar keep you from shipping to waiting subscribers.
- D-0479 also says that an epic which unblocks signet's or authstar's releases should ask for that revenue, not offer a time or penalty figure again.
- E-0035 is that kind of epic. Its outcome is done when a story in a kluster project, such as authstar, reaches review already deployed and checked.

**Why I do not set it:** the input goes on the epic, and under `plan_backlog_stories` I may give only a backlog story its first input. The figure is also yours alone: only you know how much of the 1250 USD a week E-0032 named is still held up, and how much is new.

**Alternative:** if you judge E-0035 holds up no revenue beyond E-0032's, take the planner's `penalty_per_week: 40` on E-0035.

One answer here also serves TH-0980 (S-0661). Once you set the input, I resolve neither thread: the planners write each story's value on their next run.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

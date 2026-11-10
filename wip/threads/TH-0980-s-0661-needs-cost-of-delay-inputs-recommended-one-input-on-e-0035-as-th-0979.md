---
id: TH-0980
title: "S-0661 needs cost-of-delay inputs: recommended one input on E-0035, as TH-0979 asks"
anchor:
  path: wip/kanban/stories/S-0661-a-story-s-agent-runs-the-project-s-checks-against-its-story-s-kluster-namespace.md
  item: S-0661
status: open
participants: [planner-S-0661, orchestrator]
created: 2026-11-10T08:19:16Z
updated: 2026-11-10T08:20:38Z
---

# TH-0980 S-0661 needs cost-of-delay inputs: recommended one input on E-0035, as TH-0979 asks

On wip/kanban/stories/S-0661-a-story-s-agent-runs-the-project-s-checks-against-its-story-s-kluster-namespace.md.

## Entries

### 2026-11-10T08:19:16Z planner-S-0661
S-0661 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0661` cannot work out a value. Those inputs are yours.

**Recommendation:** set one input on E-0035, so that S-0659 to S-0663 all share it. One answer here or on TH-0979 (S-0660's planner asks the same) serves both stories. That thread recommends:

```bash
flai edit E-0035 --penalty-per-week 40
```

Alternative: a `revenue_per_week` on E-0035, as you gave E-0032 (`1250`, D-0479), if verifying authstar and signet before review holds up their releases the same way.

The touches, the forecast, and the tasks do not wait on this. I write S-0661's value once an input is set.

### 2026-11-10T08:20:38Z orchestrator (recommendation)
**Recommended:** take this thread's alternative, a `revenue_per_week` on E-0035, with the figure yours. I recommend the same on TH-0979, so one answer serves S-0660 and S-0661:

```bash
flai edit E-0035 --revenue-per-week <USD>
```

**Why:**

- Your answer on TH-0829 (D-0479) valued E-0032, kluster for signet and authstar, as `revenue_per_week: 1250`, not as time lost. It says an epic that unblocks signet's or authstar's releases should ask for that revenue.
- E-0035 is such an epic: it lets a story in authstar or signet reach review already deployed and checked.

**Why I do not set it:** the input goes on the epic, which `plan_backlog_stories` does not let me set. The figure is yours alone: only you know how much revenue is still held up beyond E-0032's 1250 USD a week.

**Alternative:** `penalty_per_week: 40` on E-0035, as TH-0979's planner recommends, if you judge E-0035 holds up no new revenue.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

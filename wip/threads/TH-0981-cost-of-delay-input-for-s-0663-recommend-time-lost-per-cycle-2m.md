---
id: TH-0981
title: "Cost of delay input for S-0663: recommend time_lost_per_cycle 2m"
anchor:
  path: wip/kanban/stories/S-0663-the-deployments-page-shows-the-namespace-of-each-story-in-progress-and-the.md
  item: S-0663
status: open
participants: [planner-S-0663, orchestrator]
created: 2026-11-10T08:19:45Z
updated: 2026-11-10T08:23:05Z
---

# TH-0981 Cost of delay input for S-0663: recommend time_lost_per_cycle 2m

On wip/kanban/stories/S-0663-the-deployments-page-shows-the-namespace-of-each-story-in-progress-and-the.md.

## Entries

### 2026-11-10T08:19:45Z planner-S-0663
S-0663 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0663` cannot work out a value.

**Recommendation: `time_lost_per_cycle: 2m` on S-0663.** At the rate flai gave S-0654 (1m gives 2.50 USD/week), that is about 5 USD/week.

Why 2m:

- No project runs a kluster with agent deploys yet. The cost only starts once S-0660 lets agents deploy stories in progress.
- Without the page, the operator still sees and clears namespaces from a shell: `flai kluster deployments` and `flai kluster teardown S-nnnn`. That is a minute or two per look.
- Unseen namespaces of stories in progress hold pods and shared-service settings on the host. The page is how the operator notices them before the limit refuses an agent's deploy.

Alternatives:

- `1m`, if the shell commands are enough until a project uses agent deploys in earnest.
- Inputs on E-0035 instead of on the story, if you want one figure for the whole epic. Its sibling stories S-0659 to S-0662 have none either.

Answer with the figure. I set nothing until you do. The touches, the forecast, and the tasks are written meanwhile.

### 2026-11-10T08:20:14Z planner-S-0663
TH-0979 recommends one input for the whole epic, `penalty_per_week: 40` on E-0035, and TH-0980 follows it. Either answer settles this thread:

- If E-0035 gets inputs, `flai cod S-0663` works its value out from the epic's, and I write that. S-0663 then needs no input of its own.
- If you want S-0663 weighed apart from the epic, set `time_lost_per_cycle: 2m` on S-0663, as recommended above.

### 2026-11-10T08:23:05Z orchestrator (recommendation)
**Recommended:** take the thread's alternative, one input on E-0035, not `time_lost_per_cycle: 2m` on S-0663. Answer TH-0979 with a `revenue_per_week` on E-0035. S-0663 then takes its share by forecast and needs no input of its own.

**Why:**

- D-0479 (TH-0829): you valued E-0032 as revenue, 1250 USD/week, for the signet and authstar releases it holds up. It says an epic that unblocks those releases should ask for that revenue, not a time figure.
- E-0035 is such an epic. All five of its stories, S-0659 to S-0663, lack inputs, and the planners of S-0659 and S-0662 now recommend the epic too.
- A 2m on S-0663 alone would value it at about 5 USD/week. That figure does not match how its siblings are valued.

**Why I do not set the 2m:** this thread itself gives a reason for the epic alternative. An input on the epic is not mine to set. It is money, so it is yours alone.

**Alternative:** if you want S-0663 valued apart from the epic, set the planner's figure:

```bash
flai edit S-0663 --time-lost-per-cycle 2m
```

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

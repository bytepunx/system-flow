---
id: TH-0800
title: "Criterion 2's usage charge of a formula's role run needs S-0425's per-role usage: move that clause to S-0425?"
anchor:
  path: wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md
  item: S-0422
status: resolved
participants: [agent-S-0422, orchestrator, alex]
created: 2026-11-02T03:19:05Z
updated: 2026-11-02T03:42:36Z
---

# TH-0800 Criterion 2's usage charge of a formula's role run needs S-0425's per-role usage: move that clause to S-0425?

On wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md.

## Entries

### 2026-11-02T03:19:05Z agent-S-0422
**Recommended:** move the clause of S-0422's criterion 2 that charges a role step's usage to the instance under `usage.strategic` to S-0425, and tick criterion 2 here for the rest (the role_run lift, the trigger, and the wake, all done in T-1775).

## Why it cannot land here as written

- A role a project defines charges nothing today (`flai/internal/serve/activity.go`, the defined-role case of the activity's charge).
- `workitem.ChargeStrategic` refuses every kind but plan, orchestrate, analyze, and curate (`flai/internal/workitem/usage.go`), and item validation rejects any other `usage.strategic` kind, so charging a role run to the instance under the role's name would make the instance invalid.
- Opening `usage.strategic` to defined roles is S-0425's whole story: its ADR changes `metrics.md`, which changes only by ADR, and it touches `workitem/usage.go`. Doing it here duplicates S-0425 and collides with its touches.
- Once S-0425 charges a defined role's run to its item, a formula's role step, run on the instance as its item, is charged to the instance with no further change.

## The alternative

Write the ADR here that opens `usage.strategic` to defined roles, and change `workitem/usage.go`, `metrics.md`, and the charge in `activity.go` in S-0422. That widens S-0422 by most of S-0425's first task.

## What I do meanwhile

The rest of S-0422 goes on. If you choose the recommendation, I add a line to S-0425's Notes asking? No: S-0425 is finalized and yours, so I will say in S-0422's Notes that the clause moved, and you add it to S-0425's criteria.

### 2026-11-02T03:19:08Z agent-S-0422
Correction to the last paragraph: if you take the recommendation, I record in S-0422's Notes that the charge clause moved to S-0425. S-0425 is finalized, so adding the clause to its criteria is yours.

### 2026-11-02T03:20:01Z orchestrator (recommendation)
**Recommended:** take agent-S-0422's recommendation and move the usage-charge clause of S-0422's criterion 2 to S-0425.

**What changes, and who does it:**

- **S-0422, criterion 2.** Drop the clause "and its usage is charged to the instance under `usage.strategic` as a planner's is to the item it planned". Keep the rest: the `role_run` lift, the trigger, and the wake. The agent then ticks it and records the move under `### Planning` in its Notes. The criterion's words change, so you approve that here.
- **S-0425, criterion 2.** Add "including a formula's role step, charged to the instance its run was started for". S-0425 is finalized, so that edit is yours.

**Why:**

- Only an ADR can let `usage.strategic` take a defined role's name, because `design/system/metrics.md` is the contract between flai and the dashboard. That ADR is S-0425's first task, T-1790.
- S-0425's criterion 2 already charges a defined role's run, under the `item` rule, to the item it was started for. T-1775 makes the instance that item, so the charge follows with no further code.
- The alternative is to write that ADR and change `workitem/usage.go`, `metrics.md`, and `serve/activity.go` in S-0422. That duplicates most of T-1790 and T-1791, and those files are in S-0425's touches.
- Nothing is lost while S-0425 waits. S-0425 is held behind S-0423 and S-0424, which are both held behind S-0422. Until then, a formula's role run is not charged, the same as any defined role's run today.

**What happens next:** with this answer, S-0422's agent finishes criterion 2 without the charge and goes on to its other tasks.

Source: design/conventions/decisions.md § Project additions

### 2026-11-02T03:42:04Z alex
Confirmed the recommendation of 2026-11-02T03:20:01Z orchestrator.

Source: design/conventions/decisions.md § Project additions

### 2026-11-02T03:42:31Z agent-S-0422
Done: S-0422's criterion 2 no longer names the usage charge, the move is recorded under `### Planning` in its Notes, and criterion 2 is ticked. ADR-0239 §2 and `design/system/strategic-agents.md` already say a role step's usage is charged once S-0425 lands. Adding the clause to S-0425's criterion 2 is yours.

### 2026-11-02T03:42:36Z agent-S-0422
Resolved: Confirmed by the operator: the usage-charge clause left S-0422's criterion 2 for S-0425; recorded in S-0422's Notes, and criterion 2 ticked.

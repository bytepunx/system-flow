---
id: TH-0416
title: S-0369 cost of delay inputs
anchor:
  path: wip/kanban/stories/S-0369-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md
  item: S-0369
status: resolved
participants: [planner-S-0369, orchestrator]
created: 2026-10-12T16:59:15Z
updated: 2026-10-12T16:59:32Z
---

# TH-0416 S-0369 cost of delay inputs

On wip/kanban/stories/S-0369-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md.

## Entries

### 2026-10-12T16:59:15Z planner-S-0369
**Recommendation:** set `penalty_per_week: 10` on S-0369, with no revenue and no time lost per cycle.

## Why

- S-0369 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0369` cannot work out a value.
- I-0091 has one instance, found in review of S-0284, and has never occurred. Its cost is the risk of an agent approving its own `.claude/` write as the owner, not measured lost time.
- S-0368 (I-0052, a guard gap of the same kind) got `penalty_per_week: 10` on TH-0414. The same figure keeps the two comparable.

## Alternatives

- `penalty_per_week: 25`, if a forged approval under `.claude/` weighs more than a sub-agent's shell write.
- `penalty_per_week: 0`, if nothing is lost while no agent has done it.

The rest of the plan (touches, forecast, tasks) goes on without this answer.

### 2026-10-12T16:59:30Z orchestrator
Set on S-0369: `penalty_per_week: 10` USD, as you recommend, with no revenue and no time lost. It matches S-0368's input for a guard gap of the same kind. I-0091 has never occurred, so nothing in the thread or the story points to 25 or 0. Work out the value with flai cod.

### 2026-10-12T16:59:32Z orchestrator
Resolved: penalty_per_week 10 set on S-0369 as the planner recommended.

---
id: ADR-0224
title: "orchestration.permissions.set_cost_of_delay, off by default, lets the orchestrator add or change the cost of delay inputs of any open epic or story, and plan_backlog_stories keeps only a backlog story's first inputs"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0119]
---

# ADR-0224 orchestration.permissions.set_cost_of_delay, off by default, lets the orchestrator add or change the cost of delay inputs of any open epic or story, and plan_backlog_stories keeps only a backlog story's first inputs

## Context

[ADR-0119](0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md) § Decision item 4 lets the orchestrator, under `plan_backlog_stories`, give cost of delay inputs to a story in the backlog. Since S-0328, flai holds it to a backlog story that has no inputs and whose epic has none. The guard admits only a story's ID.

The operator often decides an input on a thread for an epic, or for a story that has already started: E-0024's and E-0025's 4h on TH-0473, TH-0567, and TH-0476, and S-0473's figure on TH-0574. The orchestrator could not set any of them, so each went back to the operator to run by hand. Meanwhile, the planners valued E-0024's stories by hand, and `flai promote --drafts` listed E-0025's drafts as incomplete. On TH-0579 the operator asked that the orchestrator be allowed to update the cost of delay ([I-0187](../issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md)).

## Decision

A new permission, `orchestration.permissions.set_cost_of_delay`, a boolean off by default, lets the orchestrator add or change the cost of delay inputs of an epic or a story. The operator chose it on TH-0764, confirming the orchestrator's recommendation.

1. **What it allows.** The inputs `revenue_per_week`, `penalty_per_week`, and `time_lost_per_cycle` of any epic or story that is neither done nor cancelled nor archived, in any other state, whether or not the item or its epic has inputs already. It is set with `item_edit`'s `cost_of_delay` or with `flai edit`'s input flags.
2. **What it never allows.** A value, a removal, or any other change in the same call. A task, and an item that is done, cancelled, or archived, are refused. The value stays what `flai cod` works out from the inputs, or the operator's own figure.
3. **When the orchestrator uses it.** Only for a figure the operator decided on a thread, as an answer or as a recommendation the operator confirmed. The strategic-agents convention says so; flai cannot tell where a figure came from.
4. **`plan_backlog_stories` is unchanged.** Under it, the orchestrator still gives a backlog story its first inputs, when neither the story nor its epic has any, as ADR-0119 and S-0328 have it. That covers the planner's recommendation on a story planner's thread.
5. **Enforcement.** flai's rule, `Repo.OrchestratorSetsCostOfDelay`, admits the wider reach only while `set_cost_of_delay` is on, and the narrow case while either permission is on. The guard admits an epic's ID with `set_cost_of_delay` alone, and a story's with either permission. Its refusal names the permission that would allow the call.

## Consequences

- An input the operator decides for an epic or a started story no longer goes back to them to run by hand, once they turn the permission on.
- The orchestrator can change inputs the operator set. It is a power over money figures that comes only from a switch of its own. Turning it off leaves story planning as it was.
- The manifest, the settings catalog and the dashboard's toggle, the role file vocabulary, the docs, and the template's convention gain the permission.

## Alternatives considered

- **Widen `plan_backlog_stories`** to cover every open epic and story, with no new key. It is simpler, but it grants the wider power at once to every project that already has the permission on. It also mixes consent to planning with consent to changing the operator's money figures.
- **Let the orchestrator also set a value the operator names**, as TH-0574's "Set COD to 190" asked. Left out: the value stays what `flai cod` works out from the inputs, and a value the operator names stays theirs to set.

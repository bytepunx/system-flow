---
id: T-2241
type: task
nature: improvement
title: An ADR refining ADR-0119 lets the orchestrator add or change cost of delay inputs on an epic or an open story under plan_backlog_stories, inputs only
status: backlog
parent: S-0523
owner: alex
created: 2026-10-28T14:31:26Z
updated: 2026-10-28T14:31:26Z
transitions: []
stream: S-0523
tags: [adr, orchestrator, cost-of-delay]
touches: [design/adrs/README.md]
---
# T-2241 An ADR refining ADR-0119 lets the orchestrator add or change cost of delay inputs on an epic or an open story under plan_backlog_stories, inputs only

## Work

First layer: waits for nothing. Write the ADR with `flai adr new` (it allots the number and adds the index row in `design/adrs/README.md`), refining [ADR-0119](../../../design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md), with topics `orchestration` and `planning`.

Context: I-0187's instances. On TH-0579 the operator decided `time_lost_per_cycle: 4h` on E-0024 and E-0025 and asked that the orchestrator be allowed to update the cost of delay; `flai guard` refused it `item_edit` with `cost_of_delay` on an epic, since the guard admits a story ID alone. On TH-0574 the operator asked for a cost of delay on S-0473, in progress, and flai refused: `Repo.OrchestratorSetsCostOfDelay` gives inputs to a backlog story without any alone. Until the operator ran the edits by hand, E-0024's stories carried values the planners worked out by hand and E-0025's ten stories had none, so `flai promote --drafts` listed them incomplete.

Decision: under `orchestration.permissions.plan_backlog_stories`, the orchestrator may add or change the cost of delay inputs (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`) of an epic, and of a story in any state but `done` or `cancelled`, archived items left out, with `item_edit`'s `cost_of_delay` or `flai edit`'s input flags, and nothing else in the call: no value (the planner's to work out on its next run), no removal, no other field. It does so only for a figure the operator decided on a thread (an answer, or a confirmed recommendation) or, as today, the planner's recommendation on a story planner's thread for a backlog story that has none; that is the convention's rule, not flai's. flai holds the rule as well as the guard, so a call the guard does not see is held all the same. Record the alternatives considered: a new permission of its own (`set_cost_of_delay`), rejected as a second toggle for one decision the operator has already handed over; letting it set a value the operator names, left out so that the value stays what `flai cod` works out from the inputs; keeping the backlog-only rule, rejected because it sends every epic input and every input on a started story back to the operator by hand.

Consequence to state: the permission's risk grows from backlog stories' inputs to any open epic's or story's, which the settings text says.

## Done when

- `design/adrs/<nnnn>-<slug>.md` exists with status accepted, `refines: [ADR-0119]`, topics `orchestration` and `planning`, and the context, decision, consequences, and alternatives above.
- `design/adrs/README.md` lists it, written by `flai adr new`.
- `flai check --strict` passes on the ADR and the index.

## Notes

Drafted by planner-S-0523. The ADR's number is allotted when it is written, so the story keeps the folder touch `design/adrs` for it.

---
id: T-0878
type: task
nature: improvement
title: An ADR extends ADR-0083 to the orchestrator's items and a project strategic total
status: backlog
parent: S-0226
owner: alex
created: 2026-10-05T04:44:31Z
updated: 2026-10-05T04:44:31Z
transitions: []
stream: S-0226
tags: [flai, docs]
touches: [design/adrs, design/system/strategic-agents.md, design/system/work-hierarchy.md]
---
# T-0878 An ADR extends ADR-0083 to the orchestrator's items and a project strategic total

## Work

Write the ADR that refines ADR-0083 (and ADR-0079) for the orchestrator, so that the code tasks build one decision. It waits for nothing: it is the first layer.

- The charge: an orchestrator activity, logged through `activity_log` or as its run ends, is apportioned to its span as ADR-0083 apportions a planner's, split evenly between the items the entry names, and each share goes under the `orchestrator` entry of that item's `usage.strategic` and is summed up to its epic when it is charged, with `ChargeStrategic`.
- Edge cases to decide and write down: an item named twice counts once; a named ID that does not exist is left out of the split; an activity that names a story and its epic gives the epic both shares; how the split rounds so that the shares add up to the whole.
- The project total: where an activity that names no item, or none that exists, is charged, and how `flai stats` reads it. Recommended: derive it at read time from the activity documents, the entries that charged no item, so nothing new is written to the repository; say why against storing it.
- Whether the charge is written for every kind but the planner (so S-0227's analyzer reuses it) or for the orchestrator alone.
- Add the ADR to `design/adrs/README.md`, set `refines`, and update `design/system/strategic-agents.md` and `design/system/work-hierarchy.md` to say how the orchestrator's charge works.

## Done when

- The ADR is accepted in `design/adrs` with its context, decision, consequences, and alternatives, and is listed in the ADR index.
- `design/system/strategic-agents.md` and `design/system/work-hierarchy.md` describe the orchestrator's charge and the project total as the ADR decides.
- `flai check --strict` and the markdown lint pass.

## Notes

- Drafted by planner-S-0226.

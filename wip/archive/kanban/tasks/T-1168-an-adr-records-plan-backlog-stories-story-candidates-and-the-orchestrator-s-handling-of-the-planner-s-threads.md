---
id: T-1168
type: task
nature: feature
title: An ADR records plan_backlog_stories, story candidates, and the orchestrator's handling of the planner's threads
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:39:20Z
updated: 2026-10-07T19:42:14Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:14Z
    by: agent-S-0328
stream: S-0328
tags: [adr, orchestrator]
touches: [design/adrs, design/adrs/README.md]
---
# T-1168 An ADR records plan_backlog_stories, story candidates, and the orchestrator's handling of the planner's threads

## Work

Write an ADR refining ADR-0087, ADR-0082, and S-0220's rules for answering threads, and add it to `design/adrs/README.md`. It decides:

- `orchestration.permissions.plan_backlog_stories`, off by default, lets the orchestrator start the planner, with the MCP tool `plan` or `flai plan`, for a story `flai plan --candidates` lists, one at a time.
- A story is a candidate when it is in the backlog, not archived, and lacks any part of a plan: touches, a forecast duration, a cost of delay value, or a task. It is left out while a planner runs for it, or while a thread its planner opened awaits the operator, as an epic is.
- What the orchestrator does with the planner's threads: it answers a plan thread with an approval, and a cost of delay question by choosing the inputs or taking the planner's recommendation, writing them itself with `item_edit` giving only `cost_of_delay` inputs, on a story with none set on it or its epic. The permissions that allow it, and how this narrows the rule that money questions escalate to the operator.

The plan thread on S-0328 records the recommended answer for the last point. Write what the operator answers there, or that recommendation if the thread is still open, and say so in the narrative's `## Decisions`.

It waits for nothing: every later task builds what it decides.

## Done when

- The ADR is under `design/adrs/` with status `accepted` and the `topics` `planning` and `orchestration`, and `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

The ADR's file name follows its title, so the task claims the `design/adrs` folder until it is written.
- 2026-10-07T19:42:14Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178

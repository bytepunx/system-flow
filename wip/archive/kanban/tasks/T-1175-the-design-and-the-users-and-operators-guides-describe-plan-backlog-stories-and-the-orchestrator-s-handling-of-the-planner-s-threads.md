---
id: T-1175
type: task
nature: feature
title: The design and the users' and operators' guides describe plan_backlog_stories and the orchestrator's handling of the planner's threads
status: cancelled
parent: S-0328
owner: alex
created: 2026-10-07T19:40:13Z
updated: 2026-10-07T19:42:17Z
transitions:
  - to: cancelled
    at: 2026-10-07T19:42:17Z
    by: agent-S-0328
stream: S-0328
tags: [docs, design, orchestrator]
touches: [design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/index.md]
after: [T-1171, T-1172]
---
# T-1175 The design and the users' and operators' guides describe plan_backlog_stories and the orchestrator's handling of the planner's threads

## Work

Bring the design and the guides up to what T-1169 to T-1172 built:

- `design/system/strategic-agents.md`: the planner's Starting it (the orchestrator's exception covers stories), the orchestrator's What it does with each permission (a `plan_backlog_stories` row), Answering threads (the planner's threads and the cost of delay rule), and Its permissions and the guard (the new rows), linking T-1168's ADR.
- `design/system/project-manifest.md`: the key, its default, and its meaning.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai plan --candidates` lists stories, and what the orchestrator may plan and edit. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- `docs/users/flaiover.md`: the Orchestrator page's settings panel offers the permission.
- `docs/operators/index.md`: the permission's risk beside `plan_backlog_epics`.

It waits for T-1171 and T-1172, whose behaviour it documents. It shares no path with T-1173 and runs beside it.

## Done when

- Each document says what the code does, with its `updated` date bumped.
- `docs/users/flai-reference.md` matches `make flai-reference`.
- `flai check --strict` and the markdown lint pass on the paths changed.

## Notes
- 2026-10-07T19:42:17Z: moved to cancelled: duplicates the story agent's plan (T-1174, T-1176 to T-1179), written at start while the planner ran; its refinements are folded into T-1176 to T-1178

---
id: T-3626
type: task
nature: improvement
title: The design, the user and operator guides, and the template changelog describe the orchestrator's keep of a removal
status: backlog
parent: S-0852
owner: alex
created: 2026-11-27T23:58:14Z
updated: 2026-11-27T23:58:21Z
transitions: []
stream: S-0852
tags: [flai, template]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-3621, T-3622]
---
# T-3626 The design, the user and operator guides, and the template changelog describe the orchestrator's keep of a removal

## Work

Describe what T-3620 to T-3622 built, linking the ADR of T-3618 rather than restating it:

- `design/system/strategic-agents.md`, under The orchestrator › Accepting a story (S-0221) and Its permissions and the guard: the permission, the narrative condition, and `keep_removals` on `accept` and `queue_add`.
- `design/system/flai-cli.md`, at `flai accept` and `flai queue add`: `--keep-removal` with `--by orchestrator`.
- `design/system/workflow.md`, where it says a removal of main's lines waits for the operator: the exception.
- `docs/users/flai.md`, at `flai accept` and `flai queue add`.
- `docs/operators/index.md`: the permission, with the orchestration permissions.
- `template/CHANGELOG.md`: an entry for the baseline role and convention change.

It waits for T-3621 and T-3622, whose behaviour it describes. It shares no file with T-3624 and runs beside it.

## Done when

- Each of the six documents says it, with `updated` bumped where it has front matter.
- `flai test` passes on the changed paths, the markdown lint and the `design-routes` tier among them.

## Notes

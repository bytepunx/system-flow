---
id: T-2344
type: task
nature: feature
title: orchestration.permissions.set_cost_of_delay gates the wider cost of delay inputs in flai's rule and the guard, and plan_backlog_stories keeps a backlog story's first inputs
status: in-progress
parent: S-0523
owner: alex
created: 2026-10-30T16:49:04Z
updated: 2026-10-30T16:50:28Z
transitions:
  - to: ready
    at: 2026-10-30T16:50:27Z
    by: agent-S-0523
  - to: in-progress
    at: 2026-10-30T16:50:28Z
    by: agent-S-0523
stream: S-0523
tags: []
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/workitem/promotable.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/guard_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
after: [T-2241]
---

# T-2344 orchestration.permissions.set_cost_of_delay gates the wider cost of delay inputs in flai's rule and the guard, and plan_backlog_stories keeps a backlog story's first inputs

## Work

Second layer: waits for T-2241, whose ADR-0224 states the rule. It was added after the operator chose option 1 on TH-0764. Layer 1's T-2242 and T-2243 built the wider reach under `plan_backlog_stories`. This task moves that reach to a permission of its own.

- `flai/internal/manifest`: add `orchestration.permissions.set_cost_of_delay`, a boolean off by default. Add it to `Permissions`, the `Permit*` constants, `PermissionNames` after `plan_backlog_stories`, and `Allows`. Give it its settings meaning and risk text in `settings.go`. Narrow `plan_backlog_stories`' text back to a backlog story's first inputs.
- `flai/internal/roles/conditions.go`: give the role-file vocabulary the new permission, and let the `cost-inputs` condition sit under either permission. Update `design/workflows/roles/orchestrator.md` and its template copy to match (front matter only).
- `Repo.OrchestratorSetsCostOfDelay`:
  - With `set_cost_of_delay` on, the wide rule: any open epic or story.
  - With only `plan_backlog_stories` on, the narrow rule S-0328 had: a backlog story with no inputs whose epic has none.
  - With neither on, refuse, naming both.
  - With only `plan_backlog_stories` on and an item past the narrow rule, refuse, naming `set_cost_of_delay`.
- The guard: an epic's ID needs `set_cost_of_delay`. A story's ID passes with either permission, and its refusal names `plan_backlog_stories` while both are off.

## Done when

- `flai test` on the changed paths passes.
- With `set_cost_of_delay` on alone, inputs are set on an epic, a ready story, an in-progress story, and a story with inputs, by MCP and CLI.
- With `plan_backlog_stories` on alone, a backlog story without inputs takes them, and the wider cases are refused naming `set_cost_of_delay`.
- With both off, the guard refuses each.

## Notes

Written by agent-S-0523 after TH-0764.

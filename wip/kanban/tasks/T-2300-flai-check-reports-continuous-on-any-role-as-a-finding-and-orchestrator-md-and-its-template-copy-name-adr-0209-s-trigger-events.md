---
id: T-2300
type: task
nature: feature
title: flai check reports continuous on any role as a finding, and orchestrator.md and its template copy name ADR-0209's trigger events
status: backlog
parent: S-0536
owner: alex
created: 2026-10-28T21:44:58Z
updated: 2026-10-28T21:45:32Z
transitions: []
stream: S-0536
tags: [flai, roles, orchestrator]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/roles_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
after: [T-2301]
---
# T-2300 flai check reports continuous on any role as a finding, and orchestrator.md and its template copy name ADR-0209's trigger events

## Work

ADR-0209 §1, criterion 1. `roles.Validate` (`checker.triggers` in `flai/internal/roles/roles.go`) accepts `continuous` today as a role's only trigger. Make it a `roles.file` finding on any role, naming ADR-0209 and the vocabulary to use instead: `asked`, an `{event}`, or a `{schedule}`. Keep `TriggerContinuous` only for the finding's text.

In the same task, so that `flai check` stays clean on the branch, replace `trigger: [continuous]` in `design/workflows/roles/orchestrator.md` and `template/root/design/workflows/roles/orchestrator.md` with ADR-0209 §1's list: `asked`; `{event: thread, when: awaits-orchestrator}`; `{event: moved, when: to-review}`; `{event: moved, when: to-backlog}`; `{event: edited, when: draft-changed}`; `{event: planned}`; `{event: accepted, when: release-due}`; `{event: ready-empty}`; `{event: reordered}`; and `{schedule: orchestration.schedule}`. The events and conditions are S-0535's, on main by then; use the names it settled. Leave the role file's body, the prompt, to the prompt task.

Waits for T-2301, which makes `orchestration.schedule` a manifest key the `{schedule}` trigger can name.

## Done when

- A test in `flai/internal/roles/roles_test.go` gives a role `trigger: [continuous]` and gets a `roles.file` finding naming ADR-0209, and one in `flai/internal/check/roles_test.go` sees it through `flai check`.
- Both copies of `orchestrator.md` carry the trigger list above, and `flai check --strict` reports nothing on them.
- `flai test` passes on the changed paths.

## Notes

---
id: T-2300
type: task
nature: feature
title: flai check reports continuous on any role as a finding, and orchestrator.md and its template copy name ADR-0209's trigger events
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:44:58Z
updated: 2026-10-29T07:43:11Z
transitions:
  - to: ready
    at: 2026-10-29T07:25:39Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:25:40Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T07:43:11Z
    by: agent-S-0536
stream: S-0536
tags: [flai, roles, orchestrator]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/check/roles_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md]
after: [T-2301]
usage:
  source: log
  seconds: 1051
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 35
      output: 13283
      cache_read: 2572452
      cache_write: 81760
      cost: 1.2423
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

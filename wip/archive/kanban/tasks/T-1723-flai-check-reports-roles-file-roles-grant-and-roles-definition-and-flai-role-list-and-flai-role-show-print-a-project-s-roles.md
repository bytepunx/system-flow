---
id: T-1723
type: task
nature: feature
title: flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles
status: done
parent: S-0414
owner: alex
created: 2026-10-16T05:43:11Z
updated: 2026-10-28T13:45:27Z
transitions:
  - to: ready
    at: 2026-10-28T13:32:34Z
    by: agent-S-0414
  - to: in-progress
    at: 2026-10-28T13:32:35Z
    by: agent-S-0414
  - to: done
    at: 2026-10-28T13:45:27Z
    by: agent-S-0414
stream: S-0414
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/root.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1721, T-1722]
usage:
  source: log
  seconds: 772
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 73
      output: 35100
      cache_read: 4134176
      cache_write: 174118
      cost: 2.5168
---
# T-1723 flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles

## Work

- Call `roles.Validate` from `flai check`. Report `roles.file`, `roles.grant`, and `roles.definition`, each on its field's line. `roles.definition` reads the role file's own definition fields, the ones T-1720's ADR says flai generates the harness definition from. It reports them missing, or listing a tool the role's rules refuse.
- Add `flai role list [--json]` and `flai role show <name> [--json]`.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Waits for T-1721, the loader, and T-1722, so the check runs clean on the shipped baseline files.

## Done when

The check tests cover each finding. The command tests cover both commands. `flai check --strict` in this repository reports nothing on `design/workflows/roles/`.

## Notes

- Register `c.roles()` in `Run` in `check.go`, after `c.conventions()`. Register `newRoleCmd` in the second `AddCommand` batch in `root.go`.
- `scripts/flai-reference.sh`, behind `make flai-reference`, also rewrites the flag index in `docs/operators/settings.md`. `TestReferenceIsCurrent` fails while either file is stale.

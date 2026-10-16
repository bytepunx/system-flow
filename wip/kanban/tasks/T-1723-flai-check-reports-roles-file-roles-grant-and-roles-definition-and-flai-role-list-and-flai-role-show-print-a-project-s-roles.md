---
id: T-1723
type: task
nature: feature
title: flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles
status: backlog
parent: S-0414
owner: alex
created: 2026-10-16T05:43:11Z
updated: 2026-10-16T05:43:11Z
transitions: []
stream: S-0414
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/root.go, docs/users/flai-reference.md]
after: [T-1721, T-1722]
---
# T-1723 flai check reports roles.file, roles.grant, and roles.definition, and flai role list and flai role show print a project's roles

## Work

- Call `roles.Validate` from `flai check`. Report `roles.file`, `roles.grant`, and `roles.definition`, each on its field's line. `roles.definition` reads the `.claude/agents/<name>.md` the file names and compares its `tools:` with the role's rules.
- Add `flai role list [--json]` and `flai role show <name> [--json]`.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Waits for T-1721, the loader, and T-1722, so the check runs clean on the shipped baseline files.

## Done when

The check tests cover each finding. The command tests cover both commands. `flai check --strict` in this repository reports nothing on `design/roles/`.

## Notes

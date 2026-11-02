---
id: T-2432
type: task
nature: feature
title: "roles.Validate knows a flai formula run <name> entry as the command flai formula run, so the shipped orchestrator role passes TestShippedRolesNameFlaisToolsAndCommands"
status: in-progress
parent: S-0423
owner: alex
created: 2026-11-02T05:01:36Z
updated: 2026-11-02T05:01:46Z
transitions:
  - to: ready
    at: 2026-11-02T05:01:45Z
    by: agent-S-0423
  - to: in-progress
    at: 2026-11-02T05:01:46Z
    by: agent-S-0423
stream: S-0423
tags: []
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go]
---
# T-2432 roles.Validate knows a flai formula run <name> entry as the command flai formula run, so the shipped orchestrator role passes TestShippedRolesNameFlaisToolsAndCommands

## Work

The close-out's integration tier failed `TestShippedRolesNameFlaisToolsAndCommands`: `design/workflows/roles/orchestrator.md:79: flai formula run publish is not a flai command`. `knownCommand` in `flai/internal/roles/roles.go` matches an entry's words against the command paths exactly. The guard already takes `flai formula run <name>` as a run of that formula alone (S-0422, `formulaEntry` in `flai/internal/guard/roles.go`). Make `knownCommand` take an entry `flai formula run <name>` (name not `*`) as known when `formula run` is a command, and add a case to `roles_test.go`.

## Done when

`flai test` on the two paths passes, and the close-out's integration tier passes `TestShippedRolesNameFlaisToolsAndCommands`.

## Notes

Found by the close-out at 8ac61fd. The test skips under `-short`, so the task sub-agents' `flai test` runs did not reach it.

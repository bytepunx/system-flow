---
id: T-1735
type: task
nature: feature
title: The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them
status: backlog
parent: S-0417
owner: alex
created: 2026-10-16T05:44:41Z
updated: 2026-10-16T05:44:56Z
transitions: []
stream: S-0417
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go]
---
# T-1735 The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them

## Work

- Add `roles` to the manifest: a map from a role's name to its `agent`, merged over the project's `agent` as `Manifest.PlanningAgent` merges, and its `permissions`, each a boolean or the value its role file declares.
- `Manifest.RoleAgent(name)` returns the merged agent. `planning.agent`, `orchestration.agent`, and `analysis.agent` win for their roles, so nothing changes for the three.
- `flai check` reports `manifest.roles` for an agent that is not valid, or a permission the role's file does not declare.
- Add the keys to the settings table so `flai manifest set` writes them.

First layer: the serve path (T-1738) reads it.

## Done when

`go test ./internal/manifest/...` passes. The existing settings and manifest tests are unchanged.

## Notes

The `roles.<name>.agent` key assumes TH-0471's answer to question 4. If the operator says no, this task keeps only `permissions`.

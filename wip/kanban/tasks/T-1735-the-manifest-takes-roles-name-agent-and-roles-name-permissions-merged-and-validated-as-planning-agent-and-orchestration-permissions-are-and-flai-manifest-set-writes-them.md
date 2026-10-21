---
id: T-1735
type: task
nature: feature
title: The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them
status: backlog
parent: S-0417
owner: alex
created: 2026-10-16T05:44:41Z
updated: 2026-10-21T11:58:02Z
transitions: []
stream: S-0417
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, docs/operators/settings.md]
---
# T-1735 The manifest takes roles.name.agent and roles.name.permissions, merged and validated as planning.agent and orchestration.permissions are, and flai manifest set writes them

## Work

- Add `roles` to the manifest: a map from a role's name to its `agent`, merged over the project's `agent` as `Manifest.PlanningAgent` merges, and its `permissions`, each a boolean or the value its role file declares.
- `Manifest.RoleAgent(name)` returns the merged agent. `planning.agent`, `orchestration.agent`, and `analysis.agent` win for their roles, so nothing changes for the three.
- `flai check` reports `manifest.roles` for an agent that is not valid, or a permission the role's file does not declare.
- Add the keys to the settings table so `flai manifest set` writes them.
- Add rows for `roles.<name>.agent` and `roles.<name>.permissions` to `docs/operators/settings.md`, since `TestSettingsIndexIsComplete` fails on a manifest key with no row and T-1737 runs `./cmd/...` in the same layer. T-1740 writes the prose.

First layer: the serve path (T-1738) and the guard task read it.

## Done when

`go test ./internal/manifest/...` and `go test ./cmd/ -run TestSettingsIndexIsComplete` pass. The existing settings and manifest tests are unchanged.

## Notes

- TH-0471's answer 4 confirmed `roles.<name>.agent`. The settings rows were added by planner-S-0417 on TH-0570's approval.

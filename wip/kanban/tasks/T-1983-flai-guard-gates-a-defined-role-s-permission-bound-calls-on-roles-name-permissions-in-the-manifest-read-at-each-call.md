---
id: T-1983
type: task
nature: feature
title: flai guard gates a defined role's permission-bound calls on roles.name.permissions in the manifest, read at each call
status: cancelled
parent: S-0417
owner: alex
created: 2026-10-21T11:58:10Z
updated: 2026-10-29T03:06:50Z
transitions:
  - to: cancelled
    at: 2026-10-29T03:06:50Z
    by: agent-S-0417
stream: S-0417
tags: [cli]
touches: [flai/cmd/guard.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go]
after: [T-1735]
---
# T-1983 flai guard gates a defined role's permission-bound calls on roles.name.permissions in the manifest, read at each call

## Work

S-0416's role path gates a defined role's calls on the operator's permission in the manifest, but `flai/cmd/guard.go` loads only `orchestration.permissions`. Wire the key T-1735 adds into the guard.

- With `FLAI_ROLE` naming a defined role, `flai/cmd/guard.go` reads `roles.<name>.permissions` from the manifest of the project the hook runs in, at each call, as it reads the orchestrator's.
- The role path in `flai/internal/guard/roles.go` lets a call its role file binds to a permission pass only while that permission is on, and refuses it otherwise, naming `roles.<name>.permissions.<permission>`.
- A manifest it cannot read leaves every permission off, with a warning, so the guard fails closed.
- Test a fixture role's permission-bound call with the permission on, off, unset, and turned on between two calls of one session. Keep the orchestrator's permission tests unchanged.

Waits for T-1735, which adds `roles.<name>.permissions` and its reader. Runs beside T-1738, with which it shares no path.

## Done when

`go test ./internal/guard/... ./cmd/ -run Guard` passes, and the existing guard tests pass unchanged.

## Notes

- Added by planner-S-0417 on TH-0570's approval: without it, the key T-1735 adds does nothing.
- 2026-10-29T03:06:50Z: moved to cancelled: S-0416 delivered it: cmd/guard.go guardRole reads repo.Manifest.Roles[name].Permissions at each call and internal/guard/roles_test.go covers on, off, and fail-closed

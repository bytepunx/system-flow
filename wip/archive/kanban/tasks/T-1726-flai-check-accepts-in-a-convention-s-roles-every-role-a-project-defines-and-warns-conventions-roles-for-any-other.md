---
id: T-1726
type: task
nature: feature
title: flai check accepts in a convention's roles every role a project defines and warns conventions.roles for any other
status: done
parent: S-0415
owner: alex
created: 2026-10-16T05:43:42Z
updated: 2026-10-28T20:17:09Z
transitions:
  - to: ready
    at: 2026-10-28T15:11:46Z
    by: agent-S-0415
  - to: in-progress
    at: 2026-10-28T15:11:47Z
    by: agent-S-0415
  - to: done
    at: 2026-10-28T20:17:09Z
    by: agent-S-0415
stream: S-0415
tags: [cli]
touches: [flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go]
after: [T-1725]
usage:
  source: log
  seconds: 103
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 5348
      cache_read: 1068726
      cache_write: 38542
      cost: 0.6099
---
# T-1726 flai check accepts in a convention's roles every role a project defines and warns conventions.roles for any other

## Work

`conventions.Roles` and `StrategicRoles` are fixed lists today. Make the `conventions.roles` rule take the roles the project defines, from `roles.Load`, beside `story`, `explore`, and `verify`. Keep the warning's text for an unknown role and list the known ones in it. Test a convention that lists a defined role, and one that lists a role no file defines.

First layer: shares no path with T-1725.

## Done when

`go test ./internal/conventions/...` passes. `flai check --strict` here reports nothing new.

## Notes

- Planner, S-0415 (2026-10-21): `flai/internal/check/check.go` is added to the touches. `Set.Validate` has no repository; its one caller, `checker.conventions` in `check.go`, loads the conventions and would pass the defined roles in. S-0414's T-1723 changes `check.go` too, but S-0415 waits for S-0414, so the two never run together.
- Mind the import direction: `flai/internal/roles` (S-0414) must not import `conventions` if `conventions` or `check` imports `roles`.

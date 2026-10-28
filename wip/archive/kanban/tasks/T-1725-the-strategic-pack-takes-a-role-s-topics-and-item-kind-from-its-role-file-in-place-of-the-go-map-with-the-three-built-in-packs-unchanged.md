---
id: T-1725
type: task
nature: feature
title: The strategic pack takes a role's topics and item kind from its role file in place of the Go map, with the three built-in packs unchanged
status: done
parent: S-0415
owner: alex
created: 2026-10-16T05:43:39Z
updated: 2026-10-28T15:11:24Z
transitions:
  - to: ready
    at: 2026-10-28T14:54:15Z
    by: agent-S-0415
  - to: in-progress
    at: 2026-10-28T14:54:17Z
    by: agent-S-0415
  - to: done
    at: 2026-10-28T15:11:24Z
    by: agent-S-0415
stream: S-0415
tags: [cli]
touches: [flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/story.go, flai/internal/context/context.go]
usage:
  source: log
  seconds: 1027
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 80
      output: 26713
      cache_read: 4278206
      cache_write: 245571
      cost: 2.98
---
# T-1725 The strategic pack takes a role's topics and item kind from its role file in place of the Go map, with the three built-in packs unchanged

## Work

Replace the role-to-topic map and the item checks in `flai/internal/context/strategic.go` (`ForStrategic`, `strategicItem`, `withRole`) with what `roles.Load` gives: the role's topics, and whether it takes an epic, a story, or neither. Add a test with a fixture role, `reviewer`, whose file names its own topics and takes a story. Keep every existing test in `strategic_test.go` as it is.

First layer: runs beside T-1726, which changes the conventions check.

## Done when

`go test ./internal/context/...` passes, the existing strategic tests unchanged. The fixture role's pack holds its topics' conventions and briefs.

## Notes

- Planner, S-0415 (2026-10-21): `story.go` is added to the touches. `ForRole` in `flai/internal/context/story.go` refuses an unknown role naming `conventions.StrategicRoles`; with defined roles, it names the roles `roles.Load` gives. `roleAgent` (the planner, the orchestrator, the analyzer) and `Pack.strategic()`, which `context.go` reads, also hold the three roles in Go; the role file's name or description, and a flag `ForStrategic` sets, take their place. Proposed for `## Work` on the plan thread.

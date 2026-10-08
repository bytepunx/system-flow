---
id: T-1388
type: task
nature: improvement
title: Each adapter states its capabilities, claude-code all of them and command none
status: backlog
parent: S-0352
owner: alex
created: 2026-10-08T08:45:54Z
updated: 2026-10-08T08:45:54Z
transitions: []
stream: S-0352
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
---
# T-1388 Each adapter states its capabilities, claude-code all of them and command none

## Work

- `harness.Capabilities` with `Resume`, `Guard`, `Ask`, `RoleModels`, `StrategicAgents`, and `Cost`, each documented by what `flai serve` may count on when it is true.
- `Adapter` gains `Capabilities() Capabilities`; `ClaudeCode` returns all true, `Command` all false, since flai can neither resume, guard, ask, nor read it.
- `harness_test.go` pins both, and a test fails when an adapter is added to `Adapters` without a capabilities case.

First layer: it runs together with T-1387, which touches other files.

## Done when

- `flai test flai/internal/harness/` passes.

## Notes

Layer 1 of S-0352. Nothing reads the capabilities yet: the start's refusals are the harness-without-a-guard story's.

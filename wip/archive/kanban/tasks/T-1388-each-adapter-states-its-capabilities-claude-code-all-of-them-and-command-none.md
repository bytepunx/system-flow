---
id: T-1388
type: task
nature: improvement
title: Each adapter states its capabilities, claude-code all of them and command none
status: done
parent: S-0352
owner: alex
created: 2026-10-08T08:45:54Z
updated: 2026-10-12T20:48:58Z
transitions:
  - to: ready
    at: 2026-10-12T20:40:52Z
    by: agent-S-0352
  - to: in-progress
    at: 2026-10-12T20:40:53Z
    by: agent-S-0352
  - to: done
    at: 2026-10-12T20:48:58Z
    by: agent-S-0352
stream: S-0352
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 485
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 13
      output: 5152
      cache_read: 832950
      cache_write: 23519
      cost: 0.4123
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

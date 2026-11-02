---
id: T-1769
type: task
nature: feature
title: flai/internal/formulas reads and validates a project's formulas into steps, inputs, and layers
status: done
parent: S-0421
owner: alex
created: 2026-10-16T05:47:14Z
updated: 2026-11-02T02:20:14Z
transitions:
  - to: ready
    at: 2026-11-02T02:09:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T02:09:49Z
    by: system-flow
  - to: done
    at: 2026-11-02T02:20:14Z
    by: agent-S-0421
stream: S-0421
tags: [cli]
touches: [flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go]
usage:
  source: log
  seconds: 625
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 12877
      cache_read: 1935669
      cache_write: 60629
      cost: 1.0065
---
# T-1769 flai/internal/formulas reads and validates a project's formulas into steps, inputs, and layers

## Work

Add the package `flai/internal/formulas`.

- `Load(project)` reads every formula under the folder.
- `Validate` returns each problem the check reports: a field missing, a step whose `run` is none of the three kinds, a script not under `scripts/`, a role `roles.Load` does not define, an `after` naming no step, a cycle, and an input a step uses that the formula does not declare.
- `Layers` orders the steps as the plan's layers order tasks.
- Tests cover a valid formula and each problem.

Waits for the ADR.

## Done when

`go test ./internal/formulas/...` passes.

## Notes

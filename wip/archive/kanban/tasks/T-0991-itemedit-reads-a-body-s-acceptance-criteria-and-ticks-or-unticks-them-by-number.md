---
id: T-0991
type: task
nature: remediation
title: itemedit reads a body's acceptance criteria and ticks or unticks them by number
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:48Z
updated: 2026-10-06T03:54:08Z
transitions:
  - to: ready
    at: 2026-10-06T03:48:15Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T03:48:15Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T03:54:08Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [flai/internal/itemedit/criteria.go, flai/internal/itemedit/criteria_test.go]
usage:
  source: log
  seconds: 353
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 4848
      cache_read: 776493
      cache_write: 22854
      cost: 0.3899
---
# T-0991 itemedit reads a body's acceptance criteria and ticks or unticks them by number

## Work

Add `flai/internal/itemedit/criteria.go`: a behaviour module, no I/O. `Criteria(body)` returns the checkboxes of the `## Acceptance criteria` section in order, numbered from 1, each with its text and whether it is ticked. `Tick(body, tick, untick []int)` returns the body with those boxes set to `[x]` or `[ ]` and nothing else changed, and refuses (InvalidError) a number out of range, a number in both lists, an empty change, and a body with no criteria.

Waits for nothing: it is the module the command, the MCP tool, and the host action call.

## Done when

- [ ] `Criteria` and `Tick` exist with one-line doc comments and table tests in `criteria_test.go`, covering nested and `[X]` boxes, text outside the section left alone, and each refusal.
- [ ] `go test -race ./internal/itemedit/` passes.

## Notes

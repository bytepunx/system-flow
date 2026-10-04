---
id: T-0798
type: task
nature: feature
title: The guard lets the planner run flai task new
status: done
parent: S-0255
owner: alex
created: 2026-10-04T04:03:35Z
updated: 2026-10-04T04:07:52Z
transitions:
  - to: ready
    at: 2026-10-04T04:04:26Z
    by: agent-S-0255
  - to: in-progress
    at: 2026-10-04T04:04:26Z
    by: agent-S-0255
  - to: done
    at: 2026-10-04T04:07:52Z
    by: agent-S-0255
stream: S-0255
tags: []
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go]
usage:
  source: log
  seconds: 206
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2972
      cache_read: 1286538
      cache_write: 20392
      cost: 0.4595
---

# T-0798 The guard lets the planner run flai task new

## Work

The guard holds the planner to its planning commands (`cliPlans` in `flai/internal/guard/guard.go`); `task new` is not among them, so the planner can write a task only through the MCP tool `item_new`. Add `task new`, update the refusal message that lists the planner's commands, and `flai guard`'s help in `flai/cmd/guard.go`. Waits for nothing; shares no path with T-0797.

## Done when

- a planner session running `flai task new --story S-nnnn ...` passes the guard; `flai task` with any other sub-command stays refused
- a test in `guard_test.go` fails without it
- `go test ./internal/guard/ ./cmd/ -run Guard` passes

## Notes

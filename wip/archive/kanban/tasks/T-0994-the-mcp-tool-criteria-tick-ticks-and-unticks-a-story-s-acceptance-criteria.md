---
id: T-0994
type: task
nature: remediation
title: The MCP tool criteria_tick ticks and unticks a story's acceptance criteria
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:57Z
updated: 2026-10-06T04:06:49Z
transitions:
  - to: ready
    at: 2026-10-06T03:54:10Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T03:54:10Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T04:06:49Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [flai/internal/mcpserver/criteria.go, flai/internal/mcpserver/criteria_test.go, flai/internal/mcpserver/folder.go]
after: [T-0991]
usage:
  source: log
  seconds: 759
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 10580
      cache_read: 1694550
      cache_write: 49875
      cost: 0.8509
---
# T-0994 The MCP tool criteria_tick ticks and unticks a story's acceptance criteria

## Work

Add the MCP tool `criteria_tick` in `flai/internal/mcpserver`: input `id`, `tick` and `untick` (1-based numbers), and `hash`; it ticks through `itemedit.Tick` and `itemedit.Apply` as the agent (nothing committed, as `item_edit`), and returns the criteria numbered with their state, the new hash, and what changed. A sub-agent is refused by the guard already, since the tool is not in `MCPReads`.

Waits for the itemedit task: it calls `Criteria` and `Tick`.

## Done when

- [ ] `criteria_tick` is registered and tested in `flai/internal/mcpserver/criteria_test.go`: tick, untick, a bad number, a stale hash.
- [ ] `go test -race ./internal/mcpserver/` passes.

## Notes

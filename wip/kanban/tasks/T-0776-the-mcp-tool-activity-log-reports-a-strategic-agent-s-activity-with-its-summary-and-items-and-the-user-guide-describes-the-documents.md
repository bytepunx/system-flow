---
id: T-0776
type: task
nature: feature
title: The MCP tool activity_log reports a strategic agent's activity with its summary and items, and the user guide describes the documents
status: done
parent: S-0206
owner: alex
created: 2026-10-03T18:37:51Z
updated: 2026-10-03T19:12:18Z
transitions:
  - to: ready
    at: 2026-10-03T18:38:47Z
    by: agent-S-0206
  - to: in-progress
    at: 2026-10-03T18:59:43Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T19:12:18Z
    by: agent-S-0206
stream: S-0206
tags: []
touches: [flai/internal/mcpserver, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/internal/guard, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0773]
usage:
  source: log
  seconds: 755
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 108
      output: 34353
      cache_read: 7321072
      cache_write: 130240
      cost: 2.9496
---
# T-0776 The MCP tool activity_log reports a strategic agent's activity with its summary and items, and the user guide describes the documents

## Work

`flai mcp` gains `activity_log` with `kind` (planner, orchestrator, or analyzer), `summary` (one line), and `items` (the IDs it touched). Like `agent_start`, the server calls back into the host through an option the command wires (`cmd/activity.go`), which logs the activity with T-0773's serve function against the project's main checkout and returns the entry. The guard refuses it to a sub-agent as it refuses other writes. Describe the tool in `design/system/flai-cli.md`, and the activity documents and the tool in `docs/users/flai.md` and the reference.

Waits for T-0773: it calls the serve function T-0773 writes.

## Done when

- Tests cover a call writing an entry through the tool, a bad kind and an empty summary refused, and the guard refusing it to a sub-agent
- The user guide describes the documents and the tool

## Notes

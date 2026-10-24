---
id: T-1838
type: task
nature: feature
title: The MCP thread_reply writes the operator's answer as a digest entry and names it in its result
status: done
parent: S-0433
owner: alex
created: 2026-10-16T06:07:41Z
updated: 2026-10-24T20:58:04Z
transitions:
  - to: ready
    at: 2026-10-24T20:35:50Z
    by: agent-S-0433
  - to: in-progress
    at: 2026-10-24T20:35:51Z
    by: agent-S-0433
  - to: done
    at: 2026-10-24T20:58:04Z
    by: agent-S-0433
stream: S-0433
tags: [cli]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go]
after: [T-1837]
usage:
  source: log
  seconds: 1333
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11580
      cache_read: 1789741
      cache_write: 98881
      cost: 1.1658
---
# T-1838 The MCP thread_reply writes the operator's answer as a digest entry and names it in its result

## Work

In `flai/internal/mcpserver/server.go`, after `threadReply` saves the entry, call `digest.OnAnswer` and add `digest` to the result with the entry's ID and whether it was written or updated, or its warning.

Waits for T-1837. Runs beside the CLI and host API tasks.

## Done when

`server_test.go` covers the operator's reply through the tool, which writes an entry, and an agent's reply, which writes none; `flai test flai/internal/mcpserver` passes.

## Notes

The operator's own agent replies as the operator only when its `FLAI_AGENT` is the owner; an agent `flai serve` started never is.

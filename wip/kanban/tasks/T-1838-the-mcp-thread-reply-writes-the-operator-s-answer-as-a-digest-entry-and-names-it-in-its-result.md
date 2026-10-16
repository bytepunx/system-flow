---
id: T-1838
type: task
nature: feature
title: The MCP thread_reply writes the operator's answer as a digest entry and names it in its result
status: backlog
parent: S-0433
owner: alex
created: 2026-10-16T06:07:41Z
updated: 2026-10-16T06:07:41Z
transitions: []
stream: S-0433
tags: [cli]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [T-1837]
---
# T-1838 The MCP thread_reply writes the operator's answer as a digest entry and names it in its result

## Work

In `flai/internal/mcpserver/server.go`, after `threadReply` saves the entry, call `digest.OnAnswer` and add `digest` to the result with the entry's ID and whether it was written or updated, or its warning.

Waits for T-1837. Runs beside the CLI and host API tasks.

## Done when

`server_test.go` covers the operator's reply through the tool, which writes an entry, and an agent's reply, which writes none; `flai test flai/internal/mcpserver` passes.

## Notes

The operator's own agent replies as the operator only when its `FLAI_AGENT` is the owner; an agent `flai serve` started never is.

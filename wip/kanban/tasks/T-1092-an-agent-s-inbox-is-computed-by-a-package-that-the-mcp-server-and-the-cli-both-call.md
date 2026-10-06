---
id: T-1092
type: task
nature: improvement
title: An agent's inbox is computed by a package that the MCP server and the CLI both call
status: backlog
parent: S-0269
owner: alex
created: 2026-10-06T22:52:57Z
updated: 2026-10-06T22:52:57Z
transitions: []
stream: S-0269
tags: [mcp]
touches: [flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/folder.go]
---
# T-1092 An agent's inbox is computed by a package that the MCP server and the CLI both call

## Work

The agent's inbox is built only inside the MCP server. That covers `(*server).inbox` in `flai/internal/mcpserver/server.go`, the per-agent cursor and `catchUp` in `cursor.go`, and the folder's inbox in `folder.go`. `flai task done` on the CLI has to answer the same inbox, so move it into a new package, `flai/internal/inbox`.

- `inbox.go` builds the answer for an agent name and an optional story:
  - the threads awaiting the agent;
  - the ready stories in pull order, with their holds, `can_pull`, and `pull_hold`;
  - the changes since the agent's cursor, at most 50, with `changes_omitted`;
  - what is unpublished.
- It reads and advances the same cursor file the MCP server uses, so a change reported by one is not reported again by the other.
- The MCP tool `inbox`, and the folder's inbox, call the package. Their answers, field for field, and their descriptions stay as they are.
- `inbox_test.go` covers threads awaiting the agent, the pull order with a held story, and changes reported once across two calls.

This task waits for nothing and touches no file the other layer 1 tasks touch, so it runs with them.

## Done when

- [ ] `flai/internal/inbox` builds an agent's inbox and advances its cursor, with tests.
- [ ] The MCP `inbox` tools call it, and `flai/internal/mcpserver`'s tests pass unchanged.
- [ ] `scripts/flai-test.sh` passes for `flai/internal/inbox` and `flai/internal/mcpserver`.

## Notes

Drafted by the planner. There is no `flai inbox` command today; this task does not add one.

---
id: T-1099
type: task
nature: improvement
title: The inbox is an exported mcpserver function over the agent's on-disk cursor, so the CLI can answer it as MCP does
status: done
parent: S-0274
owner: alex
created: 2026-10-06T22:53:10Z
updated: 2026-10-07T07:34:43Z
transitions:
  - to: ready
    at: 2026-10-07T07:28:17Z
    by: agent-S-0274
  - to: in-progress
    at: 2026-10-07T07:28:17Z
    by: agent-S-0274
  - to: done
    at: 2026-10-07T07:34:43Z
    by: agent-S-0274
stream: S-0274
tags: [mcp, go]
touches: [flai/internal/mcpserver/inbox.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/cursor.go]
usage:
  source: log
  seconds: 386
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 8004
      cache_read: 1155246
      cache_write: 36916
      cost: 0.5963
---
# T-1099 The inbox is an exported mcpserver function over the agent's on-disk cursor, so the CLI can answer it as MCP does

## Work

Today the inbox is computed only inside the MCP tool handler `(*server).inbox` in `flai/internal/mcpserver/server.go`. It works from the open threads, the board view's ready stories, `can_pull` and `pull_hold`, `flai_outdated`, and the changes since the agent's cursor (`catchUpWith` and `loadCursor` in `cursor.go`, a file per agent name). `flai story start` must answer the same inbox from the CLI.

Move the computation into an exported function in `flai/internal/mcpserver/inbox.go`. It takes the repo, the agent name, the runner, flai's version, and the clock, and returns `InboxOut` while advancing the same cursor file that the MCP tool advances. That way a change reported by `flai story start` is not reported again by the agent's next `inbox` call. `(*server).inbox` then calls it. `cursor.go` gives the function the cursor's path and the load and save it needs, without a `*server`.

This task waits for nothing. It runs alongside the storygit extraction, which touches no file this one does.

## Done when

- `server.go`'s `inbox` tool calls the exported function, and its answer is unchanged: the existing `server_test.go` inbox tests pass without changes.
- `flai/internal/mcpserver/inbox_test.go` shows that a change returned by the function under an agent name is not returned again by the MCP `inbox` under the same name, and that a first look covers the last 24 hours of stories and epics.
- `scripts/flai-test.sh` passes.

## Notes

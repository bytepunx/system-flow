---
id: T-1829
type: task
nature: feature
title: flai curate, the MCP tool curate, the host API's curate.run, and the curate host action start the curator
status: backlog
parent: S-0431
owner: alex
created: 2026-10-16T06:06:18Z
updated: 2026-10-23T17:03:55Z
transitions: []
stream: S-0431
tags: [cli]
touches: [flai/cmd/curate.go, flai/cmd/curate_test.go, flai/cmd/root.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/curate.go, flai/internal/mcpserver/curate_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1827]
---
# T-1829 flai curate, the MCP tool curate, the host API's curate.run, and the curate host action start the curator

## Work

- `flai/cmd/curate.go`, registered in `root.go`: `flai curate [--since <date>] [--json]`, through `serve.Curate`, starting the run itself as `flai analyze` does when no `flai serve` runs; journalled.
- `flai/internal/mcpserver/curate.go`, registered in `folder.go`: the tool `curate`, refused to an agent `flai serve` started; journalled as `mcp.curate`.
- `flai/internal/hostapi/writes.go`: `ActionCurate` and `curate.run` (`{since?}`), which runs `flai curate --json`, with its description.
- `flai/cmd/serve_actions.go`: the host action `curate` in the status and enable lists, off by default.
- Regenerate `docs/users/flai-reference.md`.

Waits for T-1827, whose `serve.Curate` every entry point calls.

## Done when

Tests cover each entry point's start and refusal, the action's enable and status, and the journal entries; `flai test flai/cmd flai/internal/mcpserver flai/internal/hostapi` passes.

## Notes

The dashboard's Settings page lists the host actions flai reports, so `curate` shows there with no flaiover change; a page of its own for the curator is proposed in the plan thread, not drafted.

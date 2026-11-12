---
id: T-3148
type: task
nature: improvement
title: flai host starts each project's MCP server on every address mcp.bind gives, resolving auto at each start, on a port free on all of them
status: backlog
parent: S-0732
owner: alex
created: 2026-11-12T15:54:59Z
updated: 2026-11-12T15:54:59Z
transitions: []
stream: S-0732
tags: [cli]
touches: [flai/cmd/host.go, flai/cmd/host_test.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go]
after: [T-3146, T-3147]
---
# T-3148 flai host starts each project's MCP server on every address mcp.bind gives, resolving auto at each start, on a port free on all of them

## Work

Have `hostLauncher.mcpAddr` in `flai/cmd/host.go` and `mcpAddrFor` in `flai/cmd/serve_mcp.go` give each project one port and the list of addresses. Read `mcp.bind` at every start, as now, and resolve `auto` there through `mcpBind`. Take a port only when it is free on every address, so the S-0697 move to a free port still holds. Compare ports without their host, as now, so a change of addresses moves no project off its port. Have `hostLauncher.mcp` pass one `--addr` per address to `flai mcp http`. A server already running keeps the addresses it started on until its next start or `flai host restart mcp`, as a change of `mcp.bind` does today; the plan's thread settles whether a change of engine restarts it.

It waits for T-3146, whose `mcpBind` answers the list and resolves `auto`, and which changes `serve_mcp.go` before this one does. It waits for T-3147 too, whose `flai mcp http` takes several `--addr`.

## Done when

- `host_test.go` or `serve_mcp_test.go` shows a server started on 127.0.0.1 and a second address with one port.
- `auto` is resolved at each start, the fake lookup's gateway added on a native engine and left out on Docker Desktop.
- A port taken on any one address moves the project to a free one.
- `flai test` on the changed files passes.

## Notes

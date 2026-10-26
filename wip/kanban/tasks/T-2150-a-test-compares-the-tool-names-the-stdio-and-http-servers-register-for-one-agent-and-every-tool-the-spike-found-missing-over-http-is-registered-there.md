---
id: T-2150
type: task
nature: feature
title: A test compares the tool names the stdio and HTTP servers register for one agent, and every tool the spike found missing over HTTP is registered there
status: backlog
parent: S-0506
owner: alex
created: 2026-10-26T04:47:12Z
updated: 2026-10-26T04:47:12Z
transitions: []
stream: S-0506
tags: [cli, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/cmd/mcp_http_test.go, flai/internal/mcphttp/mcphttp_test.go]
---
# T-2150 A test compares the tool names the stdio and HTTP servers register for one agent, and every tool the spike found missing over HTTP is registered there

## Work

Read the issues S-0505 filed on the HTTP MCP server and the results' list of tools that answered and did not. Write a test in `flai/internal/mcpserver/server_test.go` that builds the server as `flai mcp` on stdio does (`cmd/mcp.go`) and as `flai mcp http` does (`cmd/mcp_http.go`, through `mcphttp`), for the same agent name, and fails on any tool name one registers and the other does not. Register over HTTP, in `server.go` and `folder.go`, each tool the test or the spike's issues find missing, so that `task_done`, `stream_sync`, `permission_prompt`, and `wait_for_events` answer over HTTP as on stdio; a tool whose behaviour differs by transport for a reason records that reason in the test. This task waits for nothing.

## Done when

The parity test passes, each tool the spike's issues named as missing over HTTP answers in `flai/cmd/mcp_http_test.go` or `flai/internal/mcphttp/mcphttp_test.go`, and `flai test flai/internal/mcpserver flai/internal/mcphttp flai/cmd` passes.

## Notes

`cmd/mcp_http.go` omits `AutoApprove` where `cmd/mcp.go` passes it; that difference is the next task's, not this one's.

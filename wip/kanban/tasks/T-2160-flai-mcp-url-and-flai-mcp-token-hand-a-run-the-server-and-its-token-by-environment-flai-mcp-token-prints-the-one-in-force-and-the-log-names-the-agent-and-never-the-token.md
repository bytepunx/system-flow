---
id: T-2160
type: task
nature: feature
title: FLAI_MCP_URL and FLAI_MCP_TOKEN hand a run the server and its token by environment, flai mcp token prints the one in force, and the log names the agent and never the token
status: backlog
parent: S-0506
owner: alex
created: 2026-10-26T04:49:04Z
updated: 2026-10-26T04:49:04Z
transitions: []
stream: S-0506
tags: [cli, mcp]
touches: [flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcphttp/mcphttp.go, flai/internal/mcphttp/mcphttp_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2151]
---
# T-2160 FLAI_MCP_URL and FLAI_MCP_TOKEN hand a run the server and its token by environment, flai mcp token prints the one in force, and the log names the agent and never the token

## Work

Make `flai serve` able to hand a run it starts the project's HTTP MCP server and its token by environment: `agents.go`'s `spawn` sets `FLAI_MCP_URL` to the address `flai host` gave the project's server and `FLAI_MCP_TOKEN` to the token `flai mcp token` reads from `.flai-cache/mcp.token`, for every harness whose adapter asks for them (a capability or a field on `Start`, left unset by `claude-code` on the host so nothing changes there), and `withoutSpendKeys` and the log redaction keep the token out of every log line. `flai mcp token` keeps printing the token in force. In `mcphttp.go`, make the request log name the agent from `X-Flai-Agent` and never the bearer. Test in `mcphttp_test.go` that a log line of a request carries the agent and not the token, and in `agents_test.go` that a start for an adapter asking for them gets both variables and one for `claude-code` does not. Waits for T-2151, which changes `cmd/mcp_http.go` and `mcphttp.go` too.

## Done when

An adapter that asks for them gets `FLAI_MCP_URL` and `FLAI_MCP_TOKEN` in its environment, no log line written by the server or `flai serve` carries the token, and `flai test flai/cmd flai/internal/mcphttp flai/internal/serve` passes.

## Notes

The container adapter (S-0508) is the first to ask; this task gives it the two variables without naming it.

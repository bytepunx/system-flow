---
id: T-2542
type: task
nature: remediation
title: The MCP file flai writes for its HTTP server, and the one flai mcp status prints, name the per-server timeout the ADR sets
status: done
parent: S-0583
owner: alex
created: 2026-11-07T07:06:29Z
updated: 2026-11-07T07:31:39Z
transitions:
  - to: ready
    at: 2026-11-07T07:26:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T07:26:43Z
    by: system-flow
  - to: done
    at: 2026-11-07T07:31:39Z
    by: agent-S-0583
stream: S-0583
tags: [cli]
touches: [flai/internal/mcphttp/mcphttp.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/pilot/run_test.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go]
after: [T-2541]
usage:
  source: log
  seconds: 296
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7030
      cache_read: 3629155
      cache_write: 53886
      cost: 1.2773
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 18267
          cache_write: 781
          cost: 0.0047
---
# T-2542 The MCP file flai writes for its HTTP server, and the one flai mcp status prints, name the per-server timeout the ADR sets

## Work

This task waits for T-2541, which measures the figure and records it in the ADR that this code cites.

Add the figure as one exported constant in `flai/internal/mcphttp/mcphttp.go`, beside `AgentHeader`, with a comment naming the ADR.

- `httpMCPConfig` in `flai/internal/harness/adapters.go` writes it as the `timeout` of the server named `flai`. That one file serves a story's agent and a strategic agent in a container alike, through `claudeCodeMCP`.
- `mcpClientConfig` in `flai/cmd/mcp_http.go`, the sample `flai mcp status` prints, names the same `timeout`, so that a hand-written client copies it.

Leave `mcpserver.HTTPLongestWait` and `HTTPPermissionHold` as they are, and leave the stdio configuration in `mcpConfig` alone. Set no `MCP_TOOL_TIMEOUT`: S-0519 measured it cut by a second idle limit at 300 s.

Test it:

- `TestClaudeCodeInAContainerReachesTheHostsHTTPServer` in `flai/internal/harness/harness_test.go` asserts the `timeout` in the file, and that the stdio configuration carries none.
- `flai/internal/pilot/run_test.go` asserts the `mcp.json` pilot writes under the cache carries it.
- `flai/cmd/mcp_http_test.go` asserts the printed sample carries it.

## Done when

- `flai/internal/harness`, `flai/internal/pilot`, and the tests above pass with `flai test` on the changed paths.
- The `timeout` comes from the one constant in every place that writes it.
- The holds of `permission_prompt`, `wait_for_events`, and `wait_for_work` are unchanged, and their tests in `flai/internal/mcpserver` pass untouched.

## Notes

---
id: T-0971
type: task
nature: feature
title: "flai analyze [--focus] and the MCP tool analyze start an analyzer run through flai serve"
status: done
parent: S-0223
owner: alex
created: 2026-10-05T05:47:31Z
updated: 2026-10-06T20:35:36Z
transitions:
  - to: ready
    at: 2026-10-06T20:28:06Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:28:07Z
    by: agent-S-0223
  - to: review
    at: 2026-10-06T20:35:36Z
    by: agent-S-0223
  - to: done
    at: 2026-10-06T20:35:36Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/cmd/analyze.go, flai/cmd/analyze_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/analyze.go, flai/internal/mcpserver/analyze_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go]
after: [T-0968]
usage:
  source: log
  seconds: 449
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 18287
      cache_read: 2923177
      cache_write: 71308
      cost: 1.3444
---
# T-0971 flai analyze [--focus] and the MCP tool analyze start an analyzer run through flai serve

## Work

Add `flai analyze [--focus bottlenecks|intent|risk]` (`flai/cmd/analyze.go`, as `flai plan` is in `plan.go`): it asks the running `flai serve` to start the analyzer for the project through `analyze.run`, prints the run it started or the refusal (the action off, a run already going, a bad focus), and refuses without a running `flai serve`, saying how to start it.

Add the MCP tool `analyze`, with an optional `focus` (`flai/internal/mcpserver/analyze.go`, registered in `folder.go` beside `plan`), and an `Analyses` start function in the server's options (`server.go`), wired in `flai/cmd/mcp.go` and `mcp_http.go` as `Plans` is, logging a host entry as `mcpPlan` does. The guard refuses it to sub-agents and strategic agents, since it is not among their tools.

It waits for T-0968, whose `analyze.run` and serve start it calls.

## Done when

- a test runs `flai analyze` and `flai analyze --focus risk` against a serve fixture and sees the run start, and sees each refusal named
- a test calls the MCP tool `analyze` and sees the run start, and a bad focus refused
- `go test ./cmd/ ./internal/mcpserver/` passes

## Notes

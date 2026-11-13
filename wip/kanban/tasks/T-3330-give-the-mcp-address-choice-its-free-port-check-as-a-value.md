---
id: T-3330
type: task
nature: remediation
title: Give the MCP address choice its free-port check as a value
status: done
parent: S-0779
owner: alex
created: 2026-11-13T20:33:33Z
updated: 2026-11-13T20:43:55Z
transitions:
  - to: ready
    at: 2026-11-13T20:39:52Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T20:39:52Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T20:43:55Z
    by: agent-S-0779
stream: S-0779
tags: [flai]
touches: [flai/cmd/serve_mcp.go, flai/cmd/host.go, flai/cmd/serve_mcp_test.go]
usage:
  source: log
  seconds: 243
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 124
      cache_read: 1273044
      cache_write: 6283
      cost: 0.4802
---
# T-3330 Give the MCP address choice its free-port check as a value

## Work

`mcpAddrsFor` (`flai/cmd/serve_mcp.go`) and `hostLauncher.mcpAddrs` (`flai/cmd/host.go`) ask the operating system through `freeOnAll` and `freePort` whether a port is free, so a test cannot say which ports are free or held without racing for a real one.

Give that check as a value: `freeOnAll` and `mcpAddrsFor` take it, such as `free func(addr string) bool`, and `hostLauncher` carries one, used as `freePort` when nil, so that the test literals that build a `hostLauncher` keep working. Production callers pass `freePort`. Update the calls in `flai/cmd/serve_mcp_test.go` to pass `freePort` too, and change nothing else in the tests: that is the next task's. flai's behaviour does not change.

The story's Notes leave the way to the story's agent. Should it choose another, such as holding the listener and handing it over, or retrying with a new port, record the choice in the narrative's `## Decisions` and rewrite this task's Work and Done when to match before working it.

Waits for nothing.

## Done when

- `mcpAddrsFor` and `hostLauncher` take the free-port check as a value, `freePort` by default, and no production path changes behaviour.
- `flai test` on `flai/cmd/serve_mcp.go`, `flai/cmd/host.go`, and `flai/cmd/serve_mcp_test.go` passes.

## Notes

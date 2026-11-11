---
id: T-2974
type: task
nature: remediation
title: flai host moves a project's MCP server off a port another process holds, and a test reproduces I-0271
status: done
parent: S-0697
owner: alex
created: 2026-11-10T10:28:58Z
updated: 2026-11-11T08:07:50Z
transitions:
  - to: ready
    at: 2026-11-11T08:03:37Z
    by: agent-S-0697
  - to: in-progress
    at: 2026-11-11T08:03:38Z
    by: agent-S-0697
  - to: done
    at: 2026-11-11T08:07:50Z
    by: agent-S-0697
stream: S-0697
tags: [flai, host, mcp]
touches: [flai/cmd/host.go, flai/cmd/serve_mcp.go, flai/cmd/serve_mcp_test.go]
usage:
  source: log
  seconds: 252
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 3336
      cache_read: 1561740
      cache_write: 46986
      cost: 0.755
---
# T-2974 flai host moves a project's MCP server off a port another process holds, and a test reproduces I-0271

## Work

Waits for nothing: this is the fix the other tasks describe and close.

In `flai/cmd/host.go`, `hostLauncher.mcpAddr` answers the port in `l.taken[repo.MainRoot]` without looking at it again. In `flai/cmd/serve_mcp.go`, `mcpAddrFor` answers the remembered port without `freePort`. So when another process takes the port, the child exits with `cannot listen on ...: address already in use` and the host starts it again on the same port, forever.

Make both look at the port at each start:

- In `mcpAddr`, when the port in `l.taken` cannot be listened on (`freePort` false), drop it from `l.taken` and pick one as for a new project. `mcp` has already answered a live server of the project's own as external, so a held port here is another process's.
- In `mcpAddrFor`, answer the remembered port only when `freePort` says it is free, else fall through to the range scan.
- Log `mcp address moved` with the root, `from`, and `to`, and a reason saying another process holds the port, as S-0110's move does.

In `flai/cmd/serve_mcp_test.go`, add a test that reproduces I-0271. Listen on a port, remember it for the project, and call `l.mcpAddr`: it answers another port in range and logs the move. Then call it again with the project already in `l.taken` and the port held: it moves too. Close the listener at the end. Keep the existing `mcpAddr` and `mcpAddrFor` tests passing, so that a free remembered port is still kept.

## Done when

- `mcpAddr` never answers a port that `freePort` says is held, whether it came from `l.taken` or was remembered.
- The new test fails on the code before the change and passes after it.
- `flai test flai/cmd/host.go flai/cmd/serve_mcp.go flai/cmd/serve_mcp_test.go` passes.

## Notes

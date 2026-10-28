---
id: T-1727
type: task
nature: feature
title: flai prime --role and the MCP tool prime take any role a project defines, with its item flag, and refuse a role or item kind its file does not allow
status: done
parent: S-0415
owner: alex
created: 2026-10-16T05:43:51Z
updated: 2026-10-28T20:27:52Z
transitions:
  - to: ready
    at: 2026-10-28T20:17:30Z
    by: agent-S-0415
  - to: in-progress
    at: 2026-10-28T20:17:31Z
    by: agent-S-0415
  - to: done
    at: 2026-10-28T20:27:52Z
    by: agent-S-0415
stream: S-0415
tags: [cli]
touches: [flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, docs/users/flai-reference.md]
after: [T-1725]
usage:
  source: log
  seconds: 621
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 7488
      cache_read: 1598963
      cache_write: 48845
      cost: 0.8604
---
# T-1727 flai prime --role and the MCP tool prime take any role a project defines, with its item flag, and refuse a role or item kind its file does not allow

## Work

`flai/cmd/prime.go` checks `--role` against fixed values (`primeStrategic`, and the checks before it). Make it accept any defined role and take `--epic` or `--story` only as the role's file allows. A role no file defines is refused naming the defined roles; a file that does not parse is refused naming its finding. Do the same for the MCP tool `prime` in `flai/internal/mcpserver/server.go`, its description included. Add tests for a fixture role and for each refusal. Regenerate `docs/users/flai-reference.md`.

Waits for T-1725, whose pack it returns.

## Done when

`go test ./cmd/... ./internal/mcpserver/...` passes, the existing prime tests unchanged.

## Notes

- Planner, S-0415 (2026-10-21): `flai/internal/mcpserver/folder.go` is added to the touches. The `prime` tool's description, which names `plan`, `orchestrate`, and `analyze`, is registered there (`mcp.AddTool` with `Name: "prime"`); `server.go` holds `PrimeIn`'s `role` schema and `pack`, which choose by `conventions.StrategicRoles`. `prime.go`'s help (`--role plan, orchestrate, or analyze`) and its `--role` flag text change too, so the reference is regenerated.

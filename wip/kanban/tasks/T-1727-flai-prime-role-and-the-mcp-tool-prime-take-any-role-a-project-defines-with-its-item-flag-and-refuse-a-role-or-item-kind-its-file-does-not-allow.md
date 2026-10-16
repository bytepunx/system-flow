---
id: T-1727
type: task
nature: feature
title: flai prime --role and the MCP tool prime take any role a project defines, with its item flag, and refuse a role or item kind its file does not allow
status: backlog
parent: S-0415
owner: alex
created: 2026-10-16T05:43:51Z
updated: 2026-10-16T05:43:51Z
transitions: []
stream: S-0415
tags: [cli]
touches: [flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/prime_test.go, docs/users/flai-reference.md]
after: [T-1725]
---
# T-1727 flai prime --role and the MCP tool prime take any role a project defines, with its item flag, and refuse a role or item kind its file does not allow

## Work

`flai/cmd/prime.go` checks `--role` against fixed values (`primeStrategic`, and the checks before it). Make it accept any defined role and take `--epic` or `--story` only as the role's file allows. A role no file defines is refused naming the defined roles; a file that does not parse is refused naming its finding. Do the same for the MCP tool `prime` in `flai/internal/mcpserver/server.go`, its description included. Add tests for a fixture role and for each refusal. Regenerate `docs/users/flai-reference.md`.

Waits for T-1725, whose pack it returns.

## Done when

`go test ./cmd/... ./internal/mcpserver/...` passes, the existing prime tests unchanged.

## Notes

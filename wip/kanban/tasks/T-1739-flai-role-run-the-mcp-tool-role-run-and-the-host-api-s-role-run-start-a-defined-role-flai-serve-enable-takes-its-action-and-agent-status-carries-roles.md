---
id: T-1739
type: task
nature: feature
title: flai role run, the MCP tool role_run, and the host API's role.run start a defined role, flai serve enable takes its action, and agent.status carries roles
status: backlog
parent: S-0417
owner: alex
created: 2026-10-16T05:45:09Z
updated: 2026-10-16T05:45:09Z
transitions: []
stream: S-0417
tags: [cli]
touches: [flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/users/flai-reference.md]
after: [T-1738]
---
# T-1739 flai role run, the MCP tool role_run, and the host API's role.run start a defined role, flai serve enable takes its action, and agent.status carries roles

## Work

- Add `flai role run <name> [<item>] [--input key=value]... [--json]` beside `list` and `show` in `flai/cmd/role.go`, running `serve.RunRole` as `flai plan` runs `serve.Plan`.
- Add the MCP tool `role_run` `{name, item?}`, registered in `folder.go`, refused to an agent `flai serve` started, as `plan` is, and journalled as `mcp.role_run`.
- Add `role.run` to the host API.
- Make `flai serve enable` and the Settings page's action list take a defined role's action.
- `agent.status` carries `roles`.
- Regenerate `docs/users/flai-reference.md`.

Waits for T-1738, the serve path all three call.

## Done when

`go test ./cmd/... ./internal/mcpserver/... ./internal/hostapi/...` passes, with a refusal test for each entry point.

## Notes

---
id: T-1762
type: task
nature: feature
title: agent.status carries every role by name under roles with its action, run, past runs, and hold, and agent.stream reads any role's stream by name
status: backlog
parent: S-0420
owner: alex
created: 2026-10-16T05:46:34Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0420
tags: [cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/contract_test.go, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
---
# T-1762 agent.status carries every role by name under roles with its action, run, past runs, and hold, and agent.stream reads any role's stream by name

## Work

- Fill `roles` in `agent.status` for every role the project defines, the three built-in ones included. Each entry has its action's state, its run under way, up to 20 past runs, and `held`.
- Keep `plans`, `orchestrator`, `analyzer`, `past_orchestrators`, and `past_analyzers` for an older dashboard.
- `agent.stream` takes `role: <name>`, and `RoleStream` reads any role's log. Keep `role: orchestrate` and `analyze` as they are.
- Update the contract test.

First layer: the dashboard reads what this gives.

## Done when

`go test ./internal/hostapi/... ./internal/serve/...` passes.

## Notes

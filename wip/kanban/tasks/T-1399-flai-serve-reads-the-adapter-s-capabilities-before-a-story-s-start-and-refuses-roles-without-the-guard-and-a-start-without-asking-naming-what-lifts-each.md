---
id: T-1399
type: task
nature: feature
title: flai serve reads the adapter's capabilities before a story's start and refuses roles without the guard and a start without asking, naming what lifts each
status: backlog
parent: S-0355
owner: alex
created: 2026-10-08T08:49:09Z
updated: 2026-10-08T08:53:02Z
transitions: []
stream: S-0355
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/serve_actions.go]
after: [T-1398]
---
# T-1399 flai serve reads the adapter's capabilities before a story's start and refuses roles without the guard and a start without asking, naming what lifts each

## Work

- `serve.AgentConfig` gains each harness's `Guard` and `DenyProtected` and whether the project's `auto-approve` host action is on, filled by `agentConfig` in `flai/cmd/serve_actions.go` at every look.
- Before `adapter.Start` for a story's agent in `agents.go`, read `adapter.Capabilities()` and those settings:
  - without `Guard`, a story whose agent has `roles` is refused with `harness <name> cannot run flai guard, so its roles cannot run; flai serve agent harness <name> --guard none runs them unguarded`; with `guard: none` it starts and the start line says the roles run unguarded;
  - without `Ask`, the start is refused unless auto-approve is on or `deny_protected` is set, naming both.
- The refusal is recorded as a refused start is today, so that the story's agent status carries it to the dashboard.
- Tests in `agents_test.go` for the `command` harness: roles refused, roles with `guard: none`, no ask refused, auto-approve on, `deny_protected` set; and one showing a `claude-code` start is unchanged.

It waits for T-1398, whose settings it reads, and changes `serve_actions.go` after it.

## Done when

- The tests above pass.
- `flai test flai/internal/serve/ flai/cmd/` passes.

## Notes

Layer 2 of S-0355.

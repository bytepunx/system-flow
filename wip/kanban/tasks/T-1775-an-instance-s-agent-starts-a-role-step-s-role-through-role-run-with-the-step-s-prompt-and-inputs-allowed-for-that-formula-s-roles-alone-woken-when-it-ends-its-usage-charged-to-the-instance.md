---
id: T-1775
type: task
nature: feature
title: An instance's agent starts a role step's role through role_run with the step's prompt and inputs, allowed for that formula's roles alone, woken when it ends, its usage charged to the instance
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:45Z
updated: 2026-10-21T12:59:09Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go]
---
# T-1775 An instance's agent starts a role step's role through role_run with the step's prompt and inputs, allowed for that formula's roles alone, woken when it ends, its usage charged to the instance

## Work

- S-0417's `role_run` and `flai role run` refuse every agent `flai serve` started, as `plan` does. Lift that in `cmd/role.go`'s hook for the agent of an instance: the caller's `FLAI_STORY` names a story with `formula:`, the formula's steps name the role asked for, and the item given is that story. Any other role, item, or served agent is refused as before, naming why.
- The run takes the step's prompt and inputs as its input, the instance as its item, and the trigger `formula <name> <S-nnnn>`, recorded in `serve/agents.json` and the role's activity entry (`serve/role.go`).
- When the run settles, the agent of the instance is woken: its end is a change `wait_for_events` reports on the story, with the run's outcome (`worked`, `failed`, or `asked`), so that the agent holds `wait_for_events` as it does for a thread's answer and then closes the task or takes the failure path (T-1776).
- The run's usage is charged to the instance under `usage.strategic`, as `chargePlanned` in `serve/activity.go` charges a planner's run to the item it planned: a role run with an item is charged to that item.
- Tests: an instance's agent passes for its formula's role and is refused another; a served agent of an ordinary story is still refused; the trigger and the charge, with a fake harness.

First layer. It shares no path with the other tasks of the layer.

## Done when

`go test ./cmd/... ./internal/mcpserver/... ./internal/serve/...` passes.

## Notes

Rewritten by the planner on 2026-10-21, on TH-0471's answer 5 and TH-0582's approval. It was the serve-side role step before.

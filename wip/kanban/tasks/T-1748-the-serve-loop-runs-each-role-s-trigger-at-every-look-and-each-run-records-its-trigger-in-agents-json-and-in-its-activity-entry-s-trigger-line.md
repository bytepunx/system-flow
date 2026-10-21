---
id: T-1748
type: task
nature: feature
title: The serve loop runs each role's trigger at every look, and each run records its trigger in agents.json and in its activity entry's Trigger line
status: backlog
parent: S-0418
owner: alex
created: 2026-10-16T05:45:51Z
updated: 2026-10-21T12:01:54Z
transitions: []
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go]
after: [T-1743, T-1744, T-1746]
---
# T-1748 The serve loop runs each role's trigger at every look, and each run records its trigger in agents.json and in its activity entry's Trigger line

## Work

- Wire the event queue, the schedulers, and the continuous launchers into the serve loop beside the replanner, the orchestrator, and the analysis scheduler. Each project gets one of each per defined role whose trigger needs it, rebuilt when its role files change.
- `RunRole` takes the trigger.
- The run's record and its activity entry's `- Trigger:` line carry `asked`, `<event> <ID>`, `schedule <spec>`, or `continuous`, joined by semicolons when coalesced.

Waits for the three trigger tasks, which it wires together.

## Done when

`go test ./internal/serve/...` passes, with a loop test that starts a fixture role once per trigger kind.

## Notes

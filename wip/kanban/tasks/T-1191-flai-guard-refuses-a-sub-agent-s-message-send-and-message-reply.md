---
id: T-1191
type: task
nature: feature
title: flai guard refuses a sub-agent's message_send and message_reply
status: backlog
parent: S-0331
owner: alex
created: 2026-10-07T20:14:37Z
updated: 2026-10-07T20:14:37Z
transitions: []
stream: S-0331
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
---
# T-1191 flai guard refuses a sub-agent's message_send and message_reply

## Work

Hold sub-agents to reading, as ADR-0060 does for thread writes. It waits for nothing: it needs only the tool names, which S-0331's criteria fix. It shares no path with T-1190 and runs beside it.

- Add `message_send`, `message_reply`, and the `flai message send` and `reply` commands to what `flai guard` refuses a sub-agent; `message_get` and `flai message list` and `show` pass.
- The planner, the orchestrator, and the analyzer have no story and are refused by the tools themselves; the guard needs no rule for them.

## Done when

- Tests in `guard_test.go` cover a sub-agent refused each write and allowed each read, and the story's agent allowed both.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes

---
id: T-1389
type: task
nature: improvement
title: Each adapter gives its reader, and usage.Read, the task shares, empty wakes, and turns are measured from events
status: backlog
parent: S-0352
owner: alex
created: 2026-10-08T08:46:02Z
updated: 2026-10-08T08:46:02Z
transitions: []
stream: S-0352
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, flai/internal/usage/turns.go, flai/internal/usage/turns_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go]
after: [T-1387, T-1388]
---
# T-1389 Each adapter gives its reader, and usage.Read, the task shares, empty wakes, and turns are measured from events

## Work

- `Adapter` gains `Reader() usage.Reader`: `ClaudeCode` returns T-1387's reader; `Command` a reader that yields nothing.
- `usage.Read` takes a reader, or a reader per path, and builds its `Record` from events: calls by message ID, attribution by `ParentCall`, starts by the `Agent`/`Task` tool call's description and prompt, results by tool ID, empty wakes from a `wait_for_events` result, turns from the agent's own calls, and the session's totals from `End`. Nothing in `log.go` or `turns.go` reads stream-json fields any more.
- `serve/usage.go` picks the reader from the run's harness, `harness.For` on the agent it was started with.
- Every existing test's expectations stay as they are; only how the input is fed changes.

It waits for T-1387, whose events it reads, and T-1388, which changes the same harness files.

## Done when

- `flai test flai/internal/usage/ flai/internal/serve/ flai/internal/harness/` passes with no expected figure changed.

## Notes

Layer 2 of S-0352.

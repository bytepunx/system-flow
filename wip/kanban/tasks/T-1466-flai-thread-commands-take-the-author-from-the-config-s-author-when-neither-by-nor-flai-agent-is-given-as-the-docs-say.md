---
id: T-1466
type: task
nature: remediation
title: flai thread commands take the author from the config's author when neither --by nor FLAI_AGENT is given, as the docs say
status: backlog
parent: S-0369
owner: alex
created: 2026-10-12T17:05:39Z
updated: 2026-10-12T17:05:39Z
transitions: []
stream: S-0369
tags: [flai, threads]
touches: [flai/cmd/thread.go, flai/cmd/thread_test.go]
---
# T-1466 flai thread commands take the author from the config's author when neither --by nor FLAI_AGENT is given, as the docs say

## Work

`threadAuthor` in `flai/cmd/thread.go` returns `--by`, then `agentIdentity()`, then the config's author. `agentIdentity()` in `flai/cmd/project.go` returns `agent` when `FLAI_AGENT` is unset, so the config author is never reached. The operator's own shell, where `FLAI_AGENT` is unset, then writes as `agent`. `permission_prompt` does not take an answer by `agent` as the owner's. `docs/users/flai.md` § Threads and the flag's help say the config author is used.

Make `threadAuthor` fall back to the config's author when `FLAI_AGENT` is unset, as `writer` in `flai/cmd/touches.go` and `movedBy` in `flai/cmd/move.go` already do. The fix covers every caller: `thread new`, `reply`, `confirm`, `resolve`, the `message` commands, and `stream answer`. Change only `thread.go`, so `project.go` and the other callers of `agentIdentity` keep their behaviour.

It waits for no task: it shares no path with the ADR task or the guard task.

## Done when

- A test in `flai/cmd/thread_test.go` replies with `FLAI_AGENT` unset and no `--by`, and finds the config's author on the entry. With `FLAI_AGENT` set, the entry carries that name.
- `flai test flai/cmd/thread.go flai/cmd/thread_test.go` passes.

## Notes

Found while planning S-0369: once the guard refuses an agent `--by`, the operator's session is the one place a name is typed, and it should get the config's author without one.

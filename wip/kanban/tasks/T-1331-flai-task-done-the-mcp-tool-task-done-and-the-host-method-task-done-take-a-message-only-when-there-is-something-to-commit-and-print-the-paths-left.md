---
id: T-1331
type: task
nature: improvement
title: flai task done, the MCP tool task_done, and the host method task.done take a message only when there is something to commit, and print the paths left
status: backlog
parent: S-0322
owner: alex
created: 2026-10-08T04:35:19Z
updated: 2026-10-08T04:38:15Z
transitions: []
stream: S-0322
tags: [flai]
touches: [flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md]
after: [T-1328]
---
# T-1331 flai task done, the MCP tool task_done, and the host method task.done take a message only when there is something to commit, and print the paths left

## Work

Carry T-1328's commit step to the three front ends ADR-0107 says answer the same result.

- In `flai/cmd/task_done.go`, drop `MarkFlagRequired("message")`, rewrite step 1 of the Long help ("git add -A and git commit -m") to say which paths are committed and which are left, and print the paths left for the other open tasks in `printTaskDone`.
- In `flai/internal/mcpserver/task.go`, rewrite `taskDoneDescription` ("commit everything in the worktree") and the `message` field's schema text to match.
- In `flai/internal/hostapi/writes.go`, the `task.done` method refuses a blank `message` ("message is required"). Pass `--message=` only when one is given, so the dashboard can close a task with nothing to commit as the CLI and MCP can. In `flai/internal/hostapi/writes_test.go`, move the blank-message cases out of `task.done`'s refusals and add a valid call without one.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add tests to `flai/cmd/task_done_test.go` and `flai/internal/mcpserver/task_test.go`: with two tasks' files in the worktree, closing one leaves the other's and names them; a close with nothing to commit and no `-m` succeeds.

It waits for T-1328, whose `Result` fields and message rule it prints. It shares no path with T-1329 and runs beside it in the third layer.

## Done when

- `flai task done T-nnnn` with nothing to commit runs without `-m`, and `task_done` and `task.done` accept a call without a message.
- The text and JSON output, and `task_done`'s answer, name the paths left uncommitted.
- `flai test flai/cmd flai/internal/mcpserver flai/internal/hostapi` passes and `docs/users/flai-reference.md` is current.

## Notes

Drafted by the planner for S-0322, from S-0312's cancelled T-1284. The host method `task.done` was added on the planner's second run: it still required a message.

---
id: T-1284
type: task
nature: improvement
title: flai task done and the MCP tool task_done take -m only when there is something to commit, and print the paths left
status: backlog
parent: S-0312
owner: alex
created: 2026-10-07T23:34:48Z
updated: 2026-10-07T23:34:48Z
transitions: []
stream: S-0312
tags: [flai]
touches: [flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, docs/users/flai-reference.md]
after: [T-1283]
---
# T-1284 flai task done and the MCP tool task_done take -m only when there is something to commit, and print the paths left

## Work

Carry T-1283's commit step to the two front ends. In `flai/cmd/task_done.go`, drop `MarkFlagRequired("message")`, rewrite step 1 of the Long help to say which paths are committed and which are left, and print the paths left for the other open tasks in `printTaskDone`. In `flai/internal/mcpserver/task.go`, rewrite `taskDoneDescription` ("commit everything in the worktree") and the `message` field's schema text to match. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add tests to `flai/cmd/task_done_test.go` and `flai/internal/mcpserver/task_test.go`: two tasks' files in the worktree, closing one leaves the other's and names them; a close with nothing to commit and no `-m` succeeds.

It waits for T-1283, whose `Result` fields and message rule it prints.

## Done when

- `flai task done T-nnnn` with nothing to commit runs without `-m`.
- The text and JSON output, and `task_done`'s answer, name the paths left uncommitted.
- `flai test flai/cmd flai/internal/mcpserver` passes and `docs/users/flai-reference.md` is current.

## Notes

Drafted by the planner for S-0312.

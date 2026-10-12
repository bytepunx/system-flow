---
id: T-1440
type: task
nature: remediation
title: taskdone.Run refuses, changing nothing, while the running flai is older than the project's history, tested through Run and the MCP tool
status: backlog
parent: S-0363
owner: alex
created: 2026-10-12T14:36:30Z
updated: 2026-10-12T14:36:30Z
transitions: []
stream: S-0363
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md]
after: [T-1439]
---
# T-1440 taskdone.Run refuses, changing nothing, while the running flai is older than the project's history, tested through Run and the MCP tool

## Work

Build the decision T-1439 records.

- In `flai/internal/taskdone/taskdone.go`, have `start` (or a check before the first step) call `release.FlaiOutdated(o.Runner, <main checkout root>, o.Version)`. When it reports, return an error with its `Message`, so `Run` changes nothing: no commit, no move, no log. Use the main checkout's root, not the worktree's, since `Options.Repo` may be opened from either. An empty or dev `Version` is never refused.
- In `flai/internal/taskdone/taskdone_test.go`, reproduce I-0127's shape: a story with two open tasks of disjoint touches, both changed in the worktree, and a main checkout tagged `flai/v1.10.0`. With `Version` 1.9.3, `Run` returns the outdated error, and the worktree, the task's state, its touches, and the narrative are unchanged. With 1.10.0 and with `dev`, it closes as today.
- In `flai/internal/mcpserver/task_test.go`, add a test that the MCP tool `task_done` on a server whose version is below the project's newest tag answers a tool error naming `flai host upgrade`, and leaves the task open.
- In `flai/cmd/task_done.go`, add a sentence to the help saying it refuses while the running flai is older than the project's history, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

The host's `task.done` runs this CLI (`flai/internal/hostapi/writes.go`), so it needs no change of its own.

Waits for T-1439: it implements what the ADR decides.

## Done when

- The new tests pass under `flai test flai/internal/taskdone flai/internal/mcpserver flai/cmd`.
- The `taskdone` test fails when the check is removed. Check this once by hand, then restore it.
- `docs/users/flai-reference.md` matches `make flai-reference`.

## Notes

Drafted by the planner for S-0363. S-0362 (in progress) adds tests to the same two test files; sync after its acceptance before you start.

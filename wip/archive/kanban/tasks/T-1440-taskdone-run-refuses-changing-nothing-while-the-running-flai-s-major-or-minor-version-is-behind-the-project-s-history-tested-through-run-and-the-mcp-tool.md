---
id: T-1440
type: task
nature: remediation
title: taskdone.Run refuses, changing nothing, while the running flai's major or minor version is behind the project's history, tested through Run and the MCP tool
status: done
parent: S-0363
owner: alex
created: 2026-10-12T14:36:30Z
updated: 2026-10-12T18:22:43Z
transitions:
  - to: ready
    at: 2026-10-12T18:15:27Z
    by: agent-S-0363
  - to: in-progress
    at: 2026-10-12T18:15:28Z
    by: agent-S-0363
  - to: done
    at: 2026-10-12T18:22:43Z
    by: agent-S-0363
stream: S-0363
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md]
after: [T-1439]
usage:
  source: log
  seconds: 435
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 48
      output: 12141
      cache_read: 2612950
      cache_write: 101254
      cost: 1.4497
---
# T-1440 taskdone.Run refuses, changing nothing, while the running flai's major or minor version is behind the project's history, tested through Run and the MCP tool

## Work

Build the decision T-1439 records.

- In `flai/internal/taskdone/taskdone.go`, have `start` call `release.FlaiOutdated(o.Runner, <main checkout root>, o.Version)`. When it reports and the running version's major or minor is below the newest's (compare with `buildinfo.Semver`), return an error naming the running and newest versions and the upgrade command, so `Run` changes nothing: no commit, no move, no log. One patch behind is not refused; the inbox's `flai_outdated` warns as today. Use the main checkout's root, not the worktree's, since `Options.Repo` may be opened from either. An empty or dev `Version` is never refused.
- In `flai/internal/taskdone/taskdone_test.go`, reproduce I-0127's shape: a story with two open tasks of disjoint touches, both changed in the worktree, and a main checkout tagged `flai/v1.10.0`. With `Version` 1.9.3, `Run` returns the outdated error, and the worktree, the task's state, its touches, and the narrative are unchanged. With 1.10.0, with a patch behind such as 1.10.0 against a tag `flai/v1.10.2`, and with `dev`, it closes as today.
- In `flai/internal/mcpserver/task_test.go`, add a test that the MCP tool `task_done` on a server whose version's minor is below the project's newest tag answers an error naming the upgrade command, and leaves the task open.
- In `flai/cmd/task_done.go`, add a sentence to the help saying it refuses while the running flai's major or minor version is behind the project's history, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

The host's `task.done` runs this CLI (`flai/internal/hostapi/writes.go`), so it needs no change of its own.

Waits for T-1439: it implements what the ADR decides.

## Done when

- The new tests pass under `flai test flai/internal/taskdone flai/internal/mcpserver flai/cmd`.
- The `taskdone` test fails when the check is removed. Check this once by hand, then restore it.
- `docs/users/flai-reference.md` matches `make flai-reference`.

## Notes

Drafted by the planner for S-0363; rethresholded by S-0363's agent after TH-0401. S-0370 shares `taskdone.go`, `taskdone_test.go`, and `cmd/task_done.go` on MS-0081: S-0363 adds the check in `start`, its tests, and the help sentence, and does not change the touches step.

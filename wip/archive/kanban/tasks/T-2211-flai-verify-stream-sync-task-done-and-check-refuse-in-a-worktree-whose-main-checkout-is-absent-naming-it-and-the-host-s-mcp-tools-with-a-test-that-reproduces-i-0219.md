---
id: T-2211
type: task
nature: remediation
title: flai verify, stream sync, task done, and check refuse in a worktree whose main checkout is absent, naming it and the host's MCP tools, with a test that reproduces I-0219
status: done
parent: S-0521
owner: alex
created: 2026-10-28T13:03:13Z
updated: 2026-10-28T13:31:09Z
transitions:
  - to: ready
    at: 2026-10-28T13:12:02Z
    by: agent-S-0521
  - to: in-progress
    at: 2026-10-28T13:12:03Z
    by: agent-S-0521
  - to: done
    at: 2026-10-28T13:31:09Z
    by: agent-S-0521
stream: S-0521
tags: [cli, docker]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/stream.go]
after: [T-2207]
usage:
  source: log
  seconds: 1146
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 57
      output: 14375
      cache_read: 2634821
      cache_write: 151368
      cost: 1.748
---
# T-2211 flai verify, stream sync, task done, and check refuse in a worktree whose main checkout is absent, naming it and the host's MCP tools, with a test that reproduces I-0219

## Work

With `RequireMain` from T-2207, which this task waits for, the commands that need the main checkout refuse early when it is absent, before any of them acts on the wrong root. `flai verify` in `flai/cmd/verify.go`: `storyWorktree`, or the step before it, calls `RequireMain("verify")` so the refusal names the main checkout and the MCP tool `verify`, in place of `it has no worktree at <worktree>/.flai-cache/worktrees/S-nnnn`. `flai stream sync` in `flai/cmd/stream_sync.go`, since `storygit` takes its base branch from `repo.MainRoot`. `flai task done` in `flai/internal/taskdone/taskdone.go`, where `mainRoot` is read, since it writes `wip/` on the main checkout. `flai check` in `flai/cmd/check.go`, since S-0505 measured it reading the worktree's tracked `wip/` for the board's (2,725 items against the host's 2,738). Each refusal is one line with the command, the story worktree, the main checkout's path, and what to do: run it on the host, or call the MCP tool that runs there (`verify`, `stream_sync`, `task_done`; `flai check` has none and runs on the host). The reproducing test goes in `flai/cmd/verify_test.go` from `verifyFixture`, which builds `root/.flai-cache/worktrees/S-004` with a `.git` file naming `root/.git/worktrees/S-004`: remove the root's `system-flow.yaml` (and `wip/`) before running `flai verify S-004` from the worktree, and check the error names the root and the MCP tool and does not contain `.flai-cache/worktrees/S-004/.flai-cache`. Add one test per other command in its test file, from the same shape of fixture. The commands' behaviour from a normal checkout and from a worktree with its main checkout present stays as it is, which the existing tests show.

## Done when

- `flai verify`, `flai stream sync`, `flai task done`, and `flai check` run from a linked worktree whose main checkout is absent each refuse with one line naming the command, the main checkout's path, and the host or the MCP tool, and change nothing
- A test in `verify_test.go` reproduces I-0219's instance (the worktree opened with the main checkout's files absent) and passes with the fix, and each other command has one test of its refusal
- From a normal checkout and from a worktree whose main checkout is present, the four commands behave as before, and `go test ./cmd/... ./internal/taskdone/...`, `gofmt`, `vet`, and `golangci-lint` pass

## Notes

The refusal is the fix, not a new way to run these commands in the container: the design keeps `wip/` and the main checkout on the host, reached through flai mcp (agent-adapters.md § A container harness of flai's own). Running the close-out inside the container with what it needs is S-0510's.

---
id: T-1091
type: task
nature: improvement
title: Opening a story's branch and worktree is a storygit function that the CLI and the MCP server can both call
status: done
parent: S-0274
owner: alex
created: 2026-10-06T22:52:57Z
updated: 2026-10-07T07:33:07Z
transitions:
  - to: ready
    at: 2026-10-07T07:28:16Z
    by: agent-S-0274
  - to: in-progress
    at: 2026-10-07T07:28:17Z
    by: agent-S-0274
  - to: done
    at: 2026-10-07T07:33:07Z
    by: agent-S-0274
stream: S-0274
tags: [cli, go]
touches: [flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/cmd/branch.go, flai/cmd/stream.go]
usage:
  source: log
  seconds: 290
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 27
      output: 9623
      cache_read: 1388906
      cache_write: 44383
      cost: 0.7169
---
# T-1091 Opening a story's branch and worktree is a storygit function that the CLI and the MCP server can both call

## Work

`flai stream open` opens the story's branch and worktree through `(*app).openStoryBranch` in `flai/cmd/branch.go`. The MCP server cannot reach that code, because it lives in package `cmd`. Move the branch and worktree opening into `flai/internal/storygit/open.go` as a function that takes an `execx.Runner`, the repo, the story ID, and whether to fetch from the remote, and returns the branch, the worktree path, and where the branch came from.

Keep in that function everything the method does now: the check that the root is a git work tree, an existing worktree reused, `git worktree prune` before fetching from the remote, the fetch of `story/S-nnnn`, relative worktree paths when `worktrees.relative_paths` is set and git supports them, and the log line. `newStreamOpenCmd` and `openStoryBranch` in `flai/cmd` then call it, and the behaviour of `flai stream open` does not change.

This task waits for nothing. It runs alongside the inbox extraction, which touches no file this one does.

## Done when

- `storygit` exports the function, and `flai/cmd/branch.go` and `flai/cmd/stream.go` call it instead of their own copy.
- `flai/internal/storygit/open_test.go` covers, against a git fixture: a new branch from the main branch, an existing local branch, an existing worktree reused, and a branch fetched from the remote.
- `scripts/flai-test.sh` passes, and so do the existing `flai stream open` tests (`stream_reopen_test.go` and others) without changes.

## Notes

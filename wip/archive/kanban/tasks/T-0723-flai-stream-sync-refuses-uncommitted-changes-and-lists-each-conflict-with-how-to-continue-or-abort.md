---
id: T-0723
type: task
nature: improvement
title: flai stream sync refuses uncommitted changes and lists each conflict with how to continue or abort
status: done
parent: S-0197
owner: arobson
created: 2026-10-03T00:40:00Z
updated: 2026-10-03T00:46:17Z
transitions:
  - to: ready
    at: 2026-10-03T00:40:34Z
    by: agent-S-0197
  - to: in-progress
    at: 2026-10-03T00:40:34Z
    by: agent-S-0197
  - to: done
    at: 2026-10-03T00:46:17Z
    by: agent-S-0197
stream: S-0197
tags: []
touches: [flai/internal/storygit, flai/cmd/branch.go, flai/cmd/stream.go, flai/cmd/stream_sync_test.go]
usage:
  source: log
  seconds: 343
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 48
      output: 18474
      cache_read: 1618033
      cache_write: 64260
      cost: 1.0892
---
# T-0723 flai stream sync refuses uncommitted changes and lists each conflict with how to continue or abort

## Work

Replace `git rebase --autostash` in `syncStoryBranch` (`flai/cmd/branch.go`) with a refusal when the worktree has uncommitted changes, listing them and saying to commit them (or stash them) and sync again. When a rebase is already in progress in the worktree, refuse too, listing its conflicting paths and how to continue or abort. On a conflict, list each conflicting path on a line of its own and the steps to continue (`git add`, `git rebase --continue`, sync again) or abort (`git rebase --abort`). Put the git reads (uncommitted paths, rebase in progress, conflicting paths) in `flai/internal/storygit`. Update `stream sync`'s help in `flai/cmd/stream.go`. Waits for nothing.

## Done when

- [ ] A sync over uncommitted changes is refused, names them, and leaves the branch and worktree as they were
- [ ] A sync that stops on a conflict lists each conflicting path and how to continue or abort, in text and in `--json`
- [ ] Tests in `flai/cmd/stream_sync_test.go` cover both

## Notes

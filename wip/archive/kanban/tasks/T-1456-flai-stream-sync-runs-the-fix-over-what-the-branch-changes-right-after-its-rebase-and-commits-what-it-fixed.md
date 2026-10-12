---
id: T-1456
type: task
nature: improvement
title: flai stream sync runs the fix over what the branch changes right after its rebase, and commits what it fixed
status: done
parent: S-0367
owner: alex
created: 2026-10-12T16:52:53Z
updated: 2026-10-12T17:06:42Z
transitions:
  - to: ready
    at: 2026-10-12T16:57:57Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T16:57:57Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:06:42Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/fix.go, flai/internal/storygit/fix_test.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/branch.go, flai/cmd/stream_sync_test.go]
after: [T-1455]
usage:
  source: log
  seconds: 525
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 16447
      cache_read: 2599193
      cache_write: 71938
      cost: 1.2955
---
# T-1456 flai stream sync runs the fix over what the branch changes right after its rebase, and commits what it fixed

## Work

- `storygit.SyncOptions` gets `Fix func(dir string, paths []string) ([]string, error)`, injected by its callers, since `verify` imports `storygit`.
- In `Rebase`, after the clean rebase and the fold, run it over every file the branch changes against the main branch, and commit what it changed on the branch as `chore: [S-nnnn] ...`; record the paths in `SyncResult.Fixed`. A fix that fails is logged and the sync stands.
- `flai stream sync` and the rebase `flai accept` runs (`cmd/branch.go`) pass `verify.Fix` with the worktree's tiers; `flai stream sync` prints what it fixed.
- Waits for T-1455, which makes `verify.Fix`.

## Done when

- [ ] A sync of a branch whose markdown a fix command rewrites commits the fix, and the worktree is clean after it, under a test.
- [ ] A sync with no fix command, or a failing one, syncs as before.

## Notes

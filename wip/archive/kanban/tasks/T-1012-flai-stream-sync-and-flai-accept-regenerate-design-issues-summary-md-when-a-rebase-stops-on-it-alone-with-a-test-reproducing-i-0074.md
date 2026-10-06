---
id: T-1012
type: task
nature: improvement
title: flai stream sync and flai accept regenerate design/issues/summary.md when a rebase stops on it alone, with a test reproducing I-0074
status: done
parent: S-0278
owner: alex
created: 2026-10-06T11:32:15Z
updated: 2026-10-06T19:56:21Z
transitions:
  - to: ready
    at: 2026-10-06T19:47:35Z
    by: agent-S-0278
  - to: in-progress
    at: 2026-10-06T19:47:35Z
    by: agent-S-0278
  - to: done
    at: 2026-10-06T19:56:21Z
    by: agent-S-0278
stream: S-0278
tags: [flai]
touches: [flai/cmd/branch.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/cmd/stream_sync_test.go]
after: [T-1011]
usage:
  source: log
  seconds: 526
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 100
      output: 545
      cache_read: 3799779
      cache_write: 121353
      cost: 1.7042
---
# T-1012 flai stream sync and flai accept regenerate design/issues/summary.md when a rebase stops on it alone, with a test reproducing I-0074

## Work

Build what T-1011's ADR decides. The steps below assume the planner's recommendation. If the ADR chose otherwise, rewrite them to match it before starting, and say why in the narrative's `## Decisions`.

In `syncStoryBranch` (`flai/cmd/branch.go`), when `git rebase` stops and every path `storygit.Conflicts` lists is a generated file, the rebase goes on without the agent:

- The generated files are `design/issues/summary.md`, the path `issues.Dir` and `issues.SummaryFile` give, relative to the worktree. Keep them in one list, which T-1013 reuses.
- Regenerate the file from the issue files as the worktree has them at that stop. `issues.WriteSummary` takes a repo rooted at the worktree, so it needs no new code beyond finding that root. `git add` it.
- Continue with `git -c core.editor=true rebase --continue`. A later commit may stop on it again, so this repeats.
- A stop where any other path also conflicts is left as it is today, with every conflicting path listed, `summary.md` included.

`mergeStoryBranch` syncs through `syncStoryBranch`, so `flai accept` gets the fix too. A helper that continues a rebase belongs in `flai/internal/storygit/sync.go`, with its test in `sync_test.go`.

Reproduce I-0074 in `flai/cmd/stream_sync_test.go`, with real git as the file's other tests use:

- Story branches A and B each record or close a different issue with flai's issue commands, so each rewrites `summary.md`. Merge A into main, as an acceptance does.
- `flai stream sync` on B then ends clean. B's `summary.md` lists the open issues of both sides and carries no conflict markers.
- `flai accept` of B merges.
- When B also conflicts in another file, the sync still stops and names that file.

This task waits for T-1011, whose ADR decides what it builds.

## Done when

- A sync or acceptance whose only conflict is `design/issues/summary.md` completes, with the file regenerated.
- A conflict in any other path still stops the sync as before.
- The reproduction test fails without the change and passes with it, and `scripts/flai-test.sh` passes.

## Notes

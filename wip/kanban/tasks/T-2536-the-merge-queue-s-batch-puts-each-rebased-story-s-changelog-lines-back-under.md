---
id: T-2536
type: task
nature: remediation
title: The merge queue's batch puts each rebased story's changelog lines back under the top Unreleased, so a batch over a release lints clean
status: ready
parent: S-0570
owner: alex
created: 2026-11-07T03:48:37Z
updated: 2026-11-07T04:16:36Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:36Z
    by: system-flow
stream: S-0570
tags: [cli, git]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
after: [T-2535]
---
# T-2536 The merge queue's batch puts each rebased story's changelog lines back under the top Unreleased, so a batch over a release lints clean

## Work

Waits for T-2535: it calls the changelog step T-2535 adds to `storygit`.

`(*build).rebase` in `flai/internal/mergequeue/batch.go` runs its own `git rebase --onto` for each story of a batch and settles stops with `storygit.ContinueOverResolved`. It never reaches `storygit.Rebase`, so T-2535's step does not run there. Merge queue batch 20261102-035614 failed at `markdown` with `MD024` on `template/CHANGELOG.md` for this reason (TH-0802). After each story's rebase settles, run the same changelog step on the batch's worktree, with the merge base `rebase` already reads as the old base and `e.Head` as the old head. Log what it rewrote, as the batch logs what it merged and wrote again.

In `flai/internal/mergequeue/batch_test.go`, add a batch of one story whose changelog line sits under an `## Unreleased` that main has since released. The batch's tree must hold one `## Unreleased`, with the story's line under it.

## Done when

- The batch test fails before the change and passes after it.
- `flai test flai/internal/mergequeue` passes.

## Notes

- Drafted by the planner. If T-2524's ADR leaves the batch to the story's own sync, cancel this task with that reason, and say so in the narrative's Decisions.

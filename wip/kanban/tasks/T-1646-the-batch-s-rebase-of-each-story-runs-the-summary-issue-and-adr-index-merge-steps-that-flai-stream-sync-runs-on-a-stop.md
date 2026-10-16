---
id: T-1646
type: task
nature: feature
title: The batch's rebase of each story runs the summary, issue, and ADR index merge steps that flai stream sync runs on a stop
status: backlog
parent: S-0401
owner: alex
created: 2026-10-16T05:01:10Z
updated: 2026-10-16T05:01:10Z
transitions: []
stream: S-0401
tags: [cli]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
---
# T-1646 The batch's rebase of each story runs the summary, issue, and ADR index merge steps that flai stream sync runs on a stop

## Work

- Expose from `storygit/sync.go` the step that continues a rebase over generated and mergeable files (`continueOverGenerated`), so that the batch's rebase can call it in the batch's worktree.
- In `batch.go`, call it when a story's rebase stops. A stop it settles goes on, and the story stays in the batch.
- Test a batch of two stories that both add an ADR, both bump one issue, and both change `summary.md`.
- First in its layer: it waits for no task.

## Done when

Those batches build without a stop, and the tests of the sync and the batch pass.

## Notes

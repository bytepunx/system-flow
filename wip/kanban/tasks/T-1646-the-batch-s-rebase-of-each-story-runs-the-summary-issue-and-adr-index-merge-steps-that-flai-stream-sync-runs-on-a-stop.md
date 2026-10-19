---
id: T-1646
type: task
nature: feature
title: The batch's rebase of each story runs the summary, issue, and ADR index merge steps that flai stream sync runs on a stop
status: done
parent: S-0401
owner: alex
created: 2026-10-16T05:01:10Z
updated: 2026-10-19T22:03:43Z
transitions:
  - to: ready
    at: 2026-10-19T21:54:29Z
    by: agent-S-0401
  - to: in-progress
    at: 2026-10-19T21:54:30Z
    by: agent-S-0401
  - to: done
    at: 2026-10-19T22:03:43Z
    by: agent-S-0401
stream: S-0401
tags: [cli]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
usage:
  source: log
  seconds: 553
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 25989
      cache_read: 3249701
      cache_write: 99685
      cost: 1.7691
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

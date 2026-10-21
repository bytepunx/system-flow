---
id: T-1949
type: task
nature: improvement
title: Keep each merge-queue batch's and bisection step's tier output under .flai-cache/merge-queue/batches and record it on the batch's result
status: done
parent: S-0457
owner: alex
created: 2026-10-21T06:09:54Z
updated: 2026-10-21T06:48:22Z
transitions:
  - to: ready
    at: 2026-10-21T06:40:01Z
    by: agent-S-0457
  - to: in-progress
    at: 2026-10-21T06:40:02Z
    by: agent-S-0457
  - to: done
    at: 2026-10-21T06:48:22Z
    by: agent-S-0457
stream: S-0457
tags: [flai]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/land_test.go, flai/internal/serve/mergequeue_test.go]
after: [T-1948]
usage:
  source: log
  seconds: 500
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 90
      output: 30177
      cache_read: 3932119
      cache_write: 198247
      cost: 2.5575
---
# T-1949 Keep each merge-queue batch's and bisection step's tier output under .flai-cache/merge-queue/batches and record it on the batch's result

## Work

- Give the `Verifier` in `flai/internal/mergequeue/batch.go` the folder its tiers log to, and make `CheckoutVerifier` pass it to `verify.Tiers` as the `LogDir` T-1948 adds. The batch's run logs to `.flai-cache/merge-queue/batches/<batch id>/`, and each bisection step in `flai/internal/mergequeue/bisect.go` to a folder of its own below it, such as `step-<n>/`, so a step does not overwrite the batch's output.
- Record the folder on `Result` in `flai/internal/mergequeue/state.go`, such as `Log string`, the path relative to the project, so `state.json`, the event log's readers, and the send-back thread can name it.
- The folder outlives the batch worktree, which `clean` removes. Prune the oldest batch folders when a new batch starts, keeping the newest 20, so the cache does not grow without end.
- Update the stub verifiers in `flai/internal/mergequeue/batch_test.go` and `flai/internal/mergequeue/bisect_test.go` to the new signature. Test that a failing batch's tier output is in its folder after the batch worktree is cleaned, that a bisection step's is in its own, and that pruning keeps the newest 20.

It waits for T-1948, whose `LogDir` and step log path it passes and reads.

## Done when

- After a red batch and its bisection, each run's tier output is under `.flai-cache/merge-queue/batches/<batch id>/`, and the batch's and the step's `Result` name the folder.
- At most 20 batch folders are kept.
- `flai test` on the changed paths passes.

## Notes

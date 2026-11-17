---
id: T-3453
type: task
nature: feature
title: The batch and each bisection step run the suites verify_at.batch lists, and the batch's record names them
status: backlog
parent: S-0809
owner: alex
created: 2026-11-16T23:49:26Z
updated: 2026-11-16T23:49:26Z
transitions: []
stream: S-0809
tags: [flai]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go]
after: [T-3451]
---
# T-3453 The batch and each bisection step run the suites verify_at.batch lists, and the batch's record names them

## Work

ADR-0312 §6, criteria 1, 2, and 3. Give the merge queue the batch's suites: `Deps` (`flai/internal/mergequeue/batch.go`) carries the `verify_at.batch` S-0806 parses from the workflow file, read from the main checkout's workflow (`Repo.Flow()`) when `flai serve` builds `Deps` (`flai/internal/serve/mergequeue.go`), not from the batch's worktree, so that no story in the batch changes what its own batch runs. Widen the `Verifier` type with the suites and have `CheckoutVerifier` pass them to `verify.Tiers` as `TierOptions`. `RunBatch` and each bisection step (`bisect.go`) pass the same suites, so a step runs no suite its batch did not. Add `Suites` to `Batch` in `state.go`, written when the batch is built and kept in `state.json`, with `all` written as `all`. A workflow with no `verify_at` is the baseline, `batch: all`. Leave the `verified` and `step_verified` lines of `events.jsonl` as they are: their fields are in `metrics.md`, which changes only with an ADR. Waits for T-3451, whose `TierOptions` suites this passes.

## Done when

- [ ] A test in `batch_test.go` shows a batch given `verify_at.batch: [behavior]` runs only the `behavior` tiers its paths select, and a test in `bisect_test.go` shows each bisection step runs the same suites.
- [ ] `TestTheBatchRunsTheTiersMarkedQueue` and the other batch and bisection tests pass unchanged under the baseline workflow, so the baseline batch runs what it ran before.
- [ ] The batch in `state.json` names its suites, and a test checks it for a list and for `all`.
- [ ] A test in `flai/internal/serve/mergequeue_test.go` shows `flai serve` gives the queue the main checkout's `verify_at.batch`.
- [ ] `flai test` on the changed paths passes.

## Notes

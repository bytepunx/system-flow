---
id: T-1974
type: task
nature: remediation
title: The merge queue's batch and bisection run their tiers with MERGE_QUEUE_STORIES naming the stories they verify
status: done
parent: S-0464
owner: alex
created: 2026-10-21T07:30:16Z
updated: 2026-10-21T07:47:38Z
transitions:
  - to: ready
    at: 2026-10-21T07:32:21Z
    by: agent-S-0464
  - to: in-progress
    at: 2026-10-21T07:32:23Z
    by: agent-S-0464
  - to: done
    at: 2026-10-21T07:47:38Z
    by: agent-S-0464
stream: S-0464
tags: [flai]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/bisect_test.go, flai/internal/verify/story.go, flai/internal/mergequeue/land_test.go, flai/internal/serve/mergequeue_test.go]
usage:
  source: log
  seconds: 915
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 11385
      cache_read: 2209517
      cache_write: 79836
      cost: 1.1711
---
# T-1974 The merge queue's batch and bisection run their tiers with MERGE_QUEUE_STORIES naming the stories they verify

## Work

This is the cause I-0174 names: `CheckoutVerifier` runs `verify.Tiers`, which "sets no CLOSE_OUT_STORY", so the tiers inside the batch (smoke's `scripts/check.sh`, and the scoped checks of integration) cannot tell which stories the batch verifies and count every finding on main.

- Give `mergequeue.Verifier` the stories it verifies: the batch passes every story it rebased (its `heads`, in order); each bisection step passes the stories of that step.
- `CheckoutVerifier` sets `MERGE_QUEUE_STORIES` to them, comma-separated, through `TierOptions.RunOptions.Env`, which `verify.Run` already adds to each tier's environment. It still sets no `CLOSE_OUT_STORY`, which means a close-out and leads the scripts to record issues.
- Name the variable with a constant beside `StoryEnv` in `flai/internal/verify/story.go`, and say in `Tiers`' comment that the caller may name a batch's stories this way.
- `flai/internal/serve/mergequeue.go` only builds `CheckoutVerifier`; change it only if the build breaks.

It waits for nothing: layer 1, beside T-1973, with no path in common. The variable's name and its comma-separated form are the plan's, so the scripts task can be written to them.

## Done when

- A test in `flai/internal/mergequeue/batch_test.go` builds a batch of two stories with a tier whose command writes its environment to a file, and finds `MERGE_QUEUE_STORIES` set to the two stories in batch order and no `CLOSE_OUT_STORY`: this reproduces I-0174's cause, the batch's tiers given no scope.
- A test in `flai/internal/mergequeue/bisect_test.go` finds each bisection step's verifier given that step's stories.
- `flai test` passes on the changed paths.

## Notes

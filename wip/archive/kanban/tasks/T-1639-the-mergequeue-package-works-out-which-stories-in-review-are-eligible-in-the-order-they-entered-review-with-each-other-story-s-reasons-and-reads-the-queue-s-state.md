---
id: T-1639
type: task
nature: feature
title: The mergequeue package works out which stories in review are eligible, in the order they entered review, with each other story's reasons, and reads the queue's state
status: done
parent: S-0398
owner: alex
created: 2026-10-16T05:00:42Z
updated: 2026-10-19T14:43:12Z
transitions:
  - to: ready
    at: 2026-10-19T14:35:44Z
    by: agent-S-0398
  - to: in-progress
    at: 2026-10-19T14:35:45Z
    by: agent-S-0398
  - to: done
    at: 2026-10-19T14:43:12Z
    by: agent-S-0398
stream: S-0398
tags: [cli]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go]
usage:
  source: log
  seconds: 447
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 75
      output: 24850
      cache_read: 3788106
      cache_write: 218989
      cost: 2.6284
---
# T-1639 The mergequeue package works out which stories in review are eligible, in the order they entered review, with each other story's reasons, and reads the queue's state

## Work

- New package `flai/internal/mergequeue`.
- `eligible.go`: for each story in review, check four things:
  - its last `flai verify` report passed at its branch's head (`verify.LastReport`)
  - every criterion is ticked
  - no thread on it or its tasks is open
  - no path is recorded as undoing main's change (`storygit.RecordedUndone`).
- Order the eligible stories by their last transition to review. Give each other story every reason it fails, with a stable `kind` per reason.
- Under pull-request acceptance (`manifest.Acceptance.PullRequest()`), report the queue as off for that mode, with the reason, and no story as eligible.
- `state.go`: the state file's types and its reader, under `.flai-cache` in the main checkout. Stories enqueued, with `by` and `at`; the batch and its step; outcomes. A missing file reads as empty. Write the writer's signature only if the reader's tests need it.
- Test each reason, the order, the pull-request mode, and a missing or malformed state file.
- First in its layer: it waits for no task.

## Done when

The package answers eligibility, including off under pull-request acceptance, and reads the state as above, and its tests pass under `flai test`.

## Notes

---
id: T-2345
type: task
nature: feature
title: TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch expects verify --all after a run that left the queue tiers to reuse the tiers that run passed, as ADR-0221 makes it
status: done
parent: S-0529
owner: alex
created: 2026-10-30T18:31:29Z
updated: 2026-10-30T18:40:18Z
transitions:
  - to: ready
    at: 2026-10-30T18:31:36Z
    by: system-flow
  - to: in-progress
    at: 2026-10-30T18:31:37Z
    by: system-flow
  - to: done
    at: 2026-10-30T18:40:18Z
    by: agent-S-0529
stream: S-0529
tags: []
touches: [flai/internal/mcpserver/verify_test.go]
usage:
  source: log
  seconds: 521
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 7536
      cache_read: 2012038
      cache_write: 83055
      cost: 1.0923
---

# T-2345 TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch expects verify --all after a run that left the queue tiers to reuse the tiers that run passed, as ADR-0221 makes it

## Work

Merge queue batch 20261030-182125 failed S-0529 at `integration` (TH-0767). `TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch` in `flai/internal/mcpserver/verify_test.go` pinned ADR-0174's old rule: a `verify --all` after a run that left the queue tiers reused nothing. T-2262 changed that on purpose, so the `all` case now reuses `ok`. The test was outside the packages S-0529 changed, so the close-out's `go-test` tier did not run it. Only `integration` did, and the close-out leaves that tier to the queue. Update the expectation and the doc comment.

It waits for nothing: T-2262 is done.

## Done when

- The test expects `ok=reused late=passed` for `all`, and its comment cites ADR-0221.
- `flai test flai/internal/mcpserver` passes.

## Notes

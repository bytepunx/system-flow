---
id: T-1666
type: task
nature: improvement
title: The batch's verification always runs the tiers marked queue that the batch's paths select
status: done
parent: S-0405
owner: alex
created: 2026-10-16T05:02:34Z
updated: 2026-10-19T21:26:55Z
transitions:
  - to: ready
    at: 2026-10-19T21:18:22Z
    by: system-flow
  - to: in-progress
    at: 2026-10-19T21:18:23Z
    by: system-flow
  - to: done
    at: 2026-10-19T21:26:55Z
    by: system-flow
stream: S-0405
tags: [cli]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
after: [T-1665]
usage:
  source: log
  seconds: 512
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 2987
      cache_read: 570956
      cache_write: 23626
      cost: 0.31
---
# T-1666 The batch's verification always runs the tiers marked queue that the batch's paths select

## Work

- In `batch.go`, run the batch's verification with every selected tier, those marked `queue` included, whatever the host action says.
- Test that a batch whose paths select `integration` runs it with the action on.
- Waits for T-1654, which adds the key.

## Done when

No batch skips a tier marked `queue`, and the tests pass.

## Notes

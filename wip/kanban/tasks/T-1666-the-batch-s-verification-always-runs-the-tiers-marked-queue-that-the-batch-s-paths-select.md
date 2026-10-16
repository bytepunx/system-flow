---
id: T-1666
type: task
nature: improvement
title: The batch's verification always runs the tiers marked queue that the batch's paths select
status: backlog
parent: S-0405
owner: alex
created: 2026-10-16T05:02:34Z
updated: 2026-10-16T05:02:34Z
transitions: []
stream: S-0405
tags: [cli]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go]
after: [T-1654]
---
# T-1666 The batch's verification always runs the tiers marked queue that the batch's paths select

## Work

- In `batch.go`, run the batch's verification with every selected tier, those marked `queue` included, whatever the host action says.
- Test that a batch whose paths select `integration` runs it with the action on.
- Waits for T-1654, which adds the key.

## Done when

No batch skips a tier marked `queue`, and the tests pass.

## Notes

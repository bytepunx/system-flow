---
id: T-1653
type: task
nature: improvement
title: An ADR records that while the queue runs, the close-out and flai verify leave the tiers marked for the queue to the batch
status: backlog
parent: S-0405
owner: alex
created: 2026-10-16T05:01:35Z
updated: 2026-10-16T05:01:35Z
transitions: []
stream: S-0405
tags: [cli]
touches: [design/adrs]
---
# T-1653 An ADR records that while the queue runs, the close-out and flai verify leave the tiers marked for the queue to the batch

## Work

- Write the ADR with `flai adr new`, refining ADR-0110 and ADR-0156. It records:
  - the tier key `queue`
  - that `flai verify` leaves those tiers out while the `merge-queue` host action is on, and says so in its step
  - that the batch always runs them
  - why the risk is bounded: no story lands without them.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, indexed, and `flai check --strict` passes.

## Notes

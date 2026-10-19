---
id: T-1653
type: task
nature: improvement
title: An ADR records that while the queue runs, the close-out and flai verify leave the tiers marked for the queue to the batch
status: done
parent: S-0405
owner: alex
created: 2026-10-16T05:01:35Z
updated: 2026-10-19T20:43:35Z
transitions:
  - to: ready
    at: 2026-10-19T20:39:49Z
    by: system-flow
  - to: in-progress
    at: 2026-10-19T20:39:50Z
    by: system-flow
  - to: done
    at: 2026-10-19T20:43:35Z
    by: system-flow
stream: S-0405
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 225
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 5344
      cache_read: 1021409
      cache_write: 42266
      cost: 0.5546
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

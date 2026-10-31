---
id: T-2374
type: task
nature: improvement
title: ADRs and the design, user, and operator docs say acceptances land only through the queue, after its coalesce window, and stats count verifications per landing
status: done
parent: S-0552
owner: alex
created: 2026-10-31T07:34:17Z
updated: 2026-10-31T08:22:58Z
transitions:
  - to: ready
    at: 2026-10-31T08:05:35Z
    by: system-flow
  - to: in-progress
    at: 2026-10-31T08:05:36Z
    by: system-flow
  - to: done
    at: 2026-10-31T08:22:58Z
    by: agent-S-0552
stream: S-0552
tags: [flai]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, design/adrs/README.md, design/system/workflow.md, design/system/flaiover-dashboard.md, design/adrs/0225-while-the-merge-queue-batches-and-is-not-held-the-operator-s-acceptance-of-a.md, design/adrs/0226-the-merge-queue-waits-merge-queue-coalesce-3-minutes-by-default-from-the.md, design/adrs/0227-the-merge-queue-logs-each-bisection-step-s-verification-as-a-step-verified.md]
after: [T-2370, T-2371, T-2372]
usage:
  source: log
  seconds: 1042
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 97
      output: 581
      cache_read: 5096535
      cache_write: 159783
      cost: 1.8279
---
# T-2374 ADRs and the design, user, and operator docs say acceptances land only through the queue, after its coalesce window, and stats count verifications per landing

## Work

Record with `flai adr new`, one decision each: an operator's acceptance enqueues the story while the queue batches and is not held, so landing goes only through the queue and no hand acceptance drops a batch (refines ADR-0173 and ADR-0221); the queue waits `merge_queue.coalesce` from the first pending enqueue before it builds a batch, unless `batch_max` stories are pending (refines ADR-0173); `flai stats` reports verifications per landed story, with bisection steps logged as verifications (refines ADR-0177, since `metrics.md` is a contract). Update `design/system/flai-cli.md`'s merge queue sections, `docs/users/flai.md`'s merge queue sections (enqueueing, accepting by hand while the queue batches, the batch), `docs/operators/index.md`'s merge-queue host action section, and regenerate `docs/users/flai-reference.md` with `make flai-reference` if a command's help changed. Waits for T-2370, T-2371, and T-2372, whose behaviour it describes.

## Done when

- [ ] The three ADRs exist with their index rows, and every doc above describes the window, the enqueue on acceptance, and the new figures.

## Notes

Layer 2.

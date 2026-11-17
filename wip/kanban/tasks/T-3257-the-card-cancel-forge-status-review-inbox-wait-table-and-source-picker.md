---
id: T-3257
type: task
nature: feature
title: The card, cancel, forge status, review, inbox, wait table, and source picker components read a state's kind, tested on a state of the workflow's own
status: done
parent: S-0758
owner: alex
created: 2026-11-12T16:10:05Z
updated: 2026-11-17T08:12:31Z
transitions:
  - to: ready
    at: 2026-11-17T08:01:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:01:29Z
    by: system-flow
  - to: done
    at: 2026-11-17T08:12:31Z
    by: agent-S-0758
stream: S-0758
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/CancelConfirm.svelte, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/explorer/SourcePicker.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/CancelConfirm.svelte.test.ts]
after: [T-3252]
usage:
  source: log
  seconds: 643
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 106
      output: 521
      cache_read: 4446531
      cache_write: 118891
      cost: 1.7387
---
# T-3257 The card, cancel, forge status, review, inbox, wait table, and source picker components read a state's kind, tested on a state of the workflow's own

## Work

Waits for T-3252: these components read the module that task adds. Replace each comparison with a state's name in these components with the module's kind, colour, or moves. Add component tests that render a card and the review page for a story in a state added between `in-progress` and `review`, and keep the components' examples under `design/components/` in step with any prop that changes.

## Done when

None of these components compares a status with a state's name, and the new and existing tests pass.

## Notes

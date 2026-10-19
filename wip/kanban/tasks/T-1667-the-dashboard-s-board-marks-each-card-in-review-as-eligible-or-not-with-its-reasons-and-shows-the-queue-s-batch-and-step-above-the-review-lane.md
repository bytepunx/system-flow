---
id: T-1667
type: task
nature: feature
title: The dashboard's board marks each card in review as eligible or not with its reasons, and shows the queue's batch and step above the review lane
status: done
parent: S-0398
owner: alex
created: 2026-10-16T05:02:45Z
updated: 2026-10-19T15:22:26Z
transitions:
  - to: ready
    at: 2026-10-19T15:16:44Z
    by: agent-S-0398
  - to: in-progress
    at: 2026-10-19T15:16:45Z
    by: agent-S-0398
  - to: done
    at: 2026-10-19T15:22:26Z
    by: agent-S-0398
stream: S-0398
tags: [dashboard]
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/BoardCard.examples.json, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flaiover/src/lib/server/board.ts]
after: [T-1656]
usage:
  source: log
  seconds: 341
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 18041
      cache_read: 2354457
      cache_write: 117622
      cost: 1.4835
---
# T-1667 The dashboard's board marks each card in review as eligible or not with its reasons, and shows the queue's batch and step above the review lane

## Work

- Read `queue` from `board.get`.
- `BoardCard.svelte`: on a card in review, a mark for its place in the queue, or for not eligible with its reasons on hover and for screen readers.
- New `MergeQueue.svelte` above the review lane: the order, the batch, and its step, or a line saying the queue is empty. Build it from the design system's tokens.
- Say so when the flai on the host sends no `queue`.
- Test the marks, the panel, and an older flai.
- Waits for T-1656, which puts `queue` on `board.get`.

## Done when

The board shows the queue as above, and the vitest tests and the flaiover lint pass.

## Notes

---
id: T-2373
type: task
nature: improvement
title: The board's merge queue panel shows the coalesce window and the batch filling
status: in-progress
parent: S-0552
owner: alex
created: 2026-10-31T07:34:12Z
updated: 2026-10-31T08:05:33Z
transitions:
  - to: ready
    at: 2026-10-31T08:05:32Z
    by: system-flow
  - to: in-progress
    at: 2026-10-31T08:05:33Z
    by: system-flow
stream: S-0552
tags: [dashboard]
touches: [flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json, docs/users/flaiover.md]
after: [T-2370]
---
# T-2373 The board's merge queue panel shows the coalesce window and the batch filling

## Work

`MergeQueue.svelte` reads the window T-2370 adds to the queue's state and shows, while no batch runs and stories are enqueued, when the window ends and how many stories fill it of `batch_max`; with the window off or ended it shows nothing new. Keep its examples in step with its props, and say so in `docs/users/flaiover.md` where the queue panel is described. Waits for T-2370, whose `Queue` field it reads.

## Done when

- [ ] A component test shows the window's end and the count, and nothing when there is no window.

## Notes

Layer 2.

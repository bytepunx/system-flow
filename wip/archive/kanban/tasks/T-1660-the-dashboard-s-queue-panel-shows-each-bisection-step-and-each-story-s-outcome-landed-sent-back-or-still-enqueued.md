---
id: T-1660
type: task
nature: feature
title: "The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued"
status: done
parent: S-0400
owner: alex
created: 2026-10-16T05:02:11Z
updated: 2026-10-19T22:24:22Z
transitions:
  - to: ready
    at: 2026-10-19T22:05:11Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T22:05:12Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T22:24:22Z
    by: agent-S-0400
stream: S-0400
tags: [dashboard]
touches: [flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts]
usage:
  source: log
  seconds: 1150
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 27221
      cache_read: 3940840
      cache_write: 112170
      cost: 1.9812
---
# T-1660 The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued

## Work

- `MergeQueue.svelte` shows the batch's bisection: each step, the stories it verified with main, the commit, and its result. It also shows the story found failing.
- It shows each batched story's outcome: landed, sent back, or still enqueued, from `batch.landing`, `batch.bisect`, the `outcomes`, and the `enqueued` list.
- The state's shape is fixed by `flai/internal/mergequeue/progress.go` (T-1909): `batch.heads`, `batch.bisect`, `batch.landing`, and the outcomes `landed` and `failing`.
- Tests in `MergeQueue.svelte.test.ts` cover a bisection under way, a bisection found, and a landing part done.
- It waits for no task: it reads only the JSON shape `progress.go` fixes, which is written already.

## Done when

The panel shows each bisection step and each story's outcome, and the tests pass.

## Notes

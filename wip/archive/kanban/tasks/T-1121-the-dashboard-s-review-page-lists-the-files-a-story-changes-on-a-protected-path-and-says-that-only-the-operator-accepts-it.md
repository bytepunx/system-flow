---
id: T-1121
type: task
nature: improvement
title: The dashboard's review page lists the files a story changes on a protected path and says that only the operator accepts it
status: done
parent: S-0286
owner: alex
created: 2026-10-06T22:54:10Z
updated: 2026-10-07T01:36:16Z
transitions:
  - to: ready
    at: 2026-10-07T01:28:31Z
    by: agent-S-0286
  - to: in-progress
    at: 2026-10-07T01:28:32Z
    by: agent-S-0286
  - to: done
    at: 2026-10-07T01:36:16Z
    by: agent-S-0286
stream: S-0286
tags: [flaiover]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-1120]
usage:
  source: log
  seconds: 464
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 9285
      cache_read: 1590094
      cache_write: 40234
      cost: 0.7479
---
# T-1121 The dashboard's review page lists the files a story changes on a protected path and says that only the operator accepts it

## Work

Add the new preview field from T-1120 to the `Preview` type in `Review.svelte`. When the field lists files, show them on the review page with a line saying that only the operator accepts the story, because it changes a path Claude Code protects. The operator's own Accept stays available.

Waits for T-1120, which names the field and its shape.

## Done when

- `Review.svelte.test.ts` shows the files and the line for a preview that lists protected files, and shows neither for one that does not.
- The flaiover tests and lint pass.

## Notes

Drafted by planner-S-0286. The acceptance API route passes the preview through. If the story's agent finds that `flaiover/src/routes/api/items/[id]/acceptance/+server.ts` or `AcceptConfirm.svelte` filters it, it widens the touches with `flai touches`.

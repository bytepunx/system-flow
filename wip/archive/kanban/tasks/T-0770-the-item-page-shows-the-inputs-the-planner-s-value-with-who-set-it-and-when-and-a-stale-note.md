---
id: T-0770
type: task
nature: feature
title: The item page shows the inputs, the planner's value with who set it and when, and a stale note
status: done
parent: S-0204
owner: alex
created: 2026-10-03T18:35:52Z
updated: 2026-10-03T19:30:50Z
transitions:
  - to: ready
    at: 2026-10-03T18:36:19Z
    by: agent-S-0204
  - to: in-progress
    at: 2026-10-03T19:25:13Z
    by: agent-S-0204
  - to: done
    at: 2026-10-03T19:30:50Z
    by: agent-S-0204
stream: S-0204
tags: []
touches: [flaiover/src/lib/planning.ts, flaiover/src/lib/planning.test.ts, flaiover/src/routes/items]
after: [T-0769]
usage:
  source: log
  seconds: 337
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 8841
      cache_read: 1889841
      cache_write: 35255
      cost: 0.7721
---
# T-0770 The item page shows the inputs, the planner's value with who set it and when, and a stale note

## Work

The item page of an epic or story shows the cost of delay inputs, the planner's `value` with its `by` and `at`, and a note when the value is stale: there is a value and the inputs were changed after it. `flaiover/src/lib/planning.ts` gives the lines and the stale test from the shape T-0769 settles, and the item page shows them. Waits for T-0769, which decides where the value's and the inputs' stamps are.

## Done when

- Tests in `flaiover/src/lib/planning.test.ts` and the item page's test cover the inputs, the value with its stamp, the stale note shown and not shown, and a value with no inputs.

## Notes

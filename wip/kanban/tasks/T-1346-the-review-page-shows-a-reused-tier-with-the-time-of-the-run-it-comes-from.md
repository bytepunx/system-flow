---
id: T-1346
type: task
nature: improvement
title: The review page shows a reused tier with the time of the run it comes from
status: backlog
parent: S-0341
owner: alex
created: 2026-10-08T08:05:24Z
updated: 2026-10-08T08:05:24Z
transitions: []
stream: S-0341
tags: [flaiover]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte]
after: [T-1343]
---
# T-1346 The review page shows a reused tier with the time of the run it comes from

## Work

Criterion 2's review page. It waits for T-1343, which fixes the record's `reused` state and `reused_from` field that this reads. It shares no path with T-1344, so the two run together.

- Add `'reused'` to `VerifyState` and `reused_from?: string` to `VerifyStep` in `flaiover/src/lib/review.ts`, with the label `reused` in `STATE_LABEL`.
- In `verifyView`, give a reused row its origin in place of a duration, such as `from 07:41`, through the page's local time, as the summary line shows `ranAt`.
- In `Review.svelte`, colour a reused row as a pass that was not run now: `text-good` with the muted weight, not `text-muted` like a step not reached.

## Done when

- `review.test.ts` covers a report with reused steps: the label, the origin in place of the duration, and no findings under them.
- `flai test flaiover/src/lib/review.ts flaiover/src/lib/review.test.ts flaiover/src/lib/components/Review.svelte` passes, the lint and the type check with it.

## Notes

Today a `reused` state would already render with its raw name (`STATE_LABEL[s.state] ?? s.state`), but grey like a step not reached and with a blank duration.

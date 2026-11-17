---
id: T-3485
type: task
nature: feature
title: The review page follows a job a state's runs started in the collapsible activity stream, with its current step, its skipped cluster steps and why, and how it ended
status: backlog
parent: S-0814
owner: alex
created: 2026-11-16T23:51:29Z
updated: 2026-11-16T23:51:29Z
transitions: []
stream: S-0814
tags: [flaiover]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/LogFeed.svelte, flaiover/src/lib/components/LogFeed.svelte.test.ts, design/components/flaiover/src/lib/components/LogFeed.examples.json, "flaiover/src/routes/api/items/[id]/resolve/+server.ts", "flaiover/src/routes/api/items/[id]/resolve/resolve.test.ts"]
after: [T-3480, T-3481]
---
# T-3485 The review page follows a job a state's runs started in the collapsible activity stream, with its current step, its skipped cluster steps and why, and how it ended

## Work

Criterion 2 on the page, and the examples of criterion 4. S-0804's T-3431 follows a Rebase and resolve job in a collapsible activity stream on the review page, built on `LogFeed.svelte`, through the resolve route. Follow a job a state's `runs:` started in the same stream, from what T-3480 adds to `resolve.status`: the state it runs for, its current step and elapsed time, its output as it is written, each skipped cluster step with its reason, and how it ended, passed and left in the state, sent back with its thread, or interrupted. Give `LogFeed.svelte` the skipped steps as a prop if it lacks one, and add an example of a job with skipped steps to `LogFeed.examples.json`.

Waits for T-3480, whose answer this reads, and for T-3481, which changes `Review.svelte` and `review.ts` too.

## Done when

- The stream follows a running `runs:` job to its end without a reload, and shows a passed job, a job sent back with its thread, an interrupted job, and a job whose cluster steps were skipped, with the reason.
- A Rebase and resolve job shows as S-0804 made it.
- `LogFeed.svelte.test.ts`, `Review.svelte.test.ts`, `review.test.ts`, and `resolve.test.ts` cover these; `flai check` reports no `components.examples` finding for `LogFeed`; and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. If S-0804 puts the stream in a component of its own rather than `Review.svelte`, change that component and its examples instead, and widen the touches.

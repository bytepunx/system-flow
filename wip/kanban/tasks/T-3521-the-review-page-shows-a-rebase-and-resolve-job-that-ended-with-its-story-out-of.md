---
id: T-3521
type: task
nature: remediation
title: The review page shows a rebase and resolve job that ended with its story out of review
status: backlog
parent: S-0822
owner: alex
created: 2026-11-17T08:19:55Z
updated: 2026-11-17T08:19:55Z
transitions: []
stream: S-0822
tags: [flaiover, dashboard]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
---
# T-3521 The review page shows a rebase and resolve job that ended with its story out of review

## Work

The review page shows the new outcome as it shows the others. In `flaiover/src/lib/review.ts`, add the kind `'left-review'` to `JobEndKind`. `endOf` maps the job's outcome `left review` to it, with a headline that says the story left review while the job ran, and the job's reason as the detail. This holds for a passed verdict and for a stopped or failed one. Today a job with no outcome falls through to the close-out wording, "no Rebase and resolve ran", which is wrong for a resolve job.

In `flaiover/src/lib/components/Review.svelte`, give the kind a tone in `END_TONE`. `data-end` on the `resolve-end` block then reads `left-review`.

Add cases to `flaiover/src/lib/review.test.ts` under the `jobView` tests, and to `flaiover/src/lib/components/Review.svelte.test.ts` under the job's tests, beside the stayed and sent-back cases.

This task waits for no other. The outcome string is `left review`, fixed by the story's plan, so it can run beside the fix in flai.

## Done when

- A resolve job with outcome `left review` shows its own headline and the job's reason on the review page, for a passed verdict and a stopped one.
- The stayed, sent-back, and interrupted cases read as before.
- vitest and svelte-check pass with `flai test` on the changed paths.

## Notes

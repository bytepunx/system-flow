---
id: T-3555
type: task
nature: improvement
title: The review page shows which round a close-out or Rebase and resolve job is on and why it started each round after the first
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:34:01Z
updated: 2026-11-17T08:34:01Z
transitions: []
stream: S-0824
tags: [flaiover, dashboard]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-3553]
---
# T-3555 The review page shows which round a close-out or Rebase and resolve job is on and why it started each round after the first

## Work

Show the round fields T-3553 adds to `resolve.status`. It waits for T-3553, whose field names the page reads.

- In `flaiover/src/lib/review.ts`, `ResolveJob` and `JobView` gain the round, the rounds, and each later round's why. `runningNote` and `endOf` say them: for example "round 2 of 3: main moved during the run, sending back integration for flai/internal/verify/story.go". Round 1 of a job that ran one round says nothing of rounds, as today.
- In `flaiover/src/lib/components/Review.svelte`, the resolve job's section (`resolve-progress`) shows the round while the job runs, and, once it ended, the round it ended on with every later round's why.

Tests:

- in `flaiover/src/lib/review.test.ts`, under the describe for a close-out or Rebase and resolve job (S-0804): a status in round 2 gives the note with its why, and a job that ended in round 1 names no round;
- in `flaiover/src/lib/components/Review.svelte.test.ts`: the section renders the round and the why.

## Done when

- The review page says which round a running job is on and why it started each round after the first, and an ended job says the round it ended on.
- A one-round job reads as before.
- `pnpm` vitest and svelte-check pass on the changed files.

## Notes

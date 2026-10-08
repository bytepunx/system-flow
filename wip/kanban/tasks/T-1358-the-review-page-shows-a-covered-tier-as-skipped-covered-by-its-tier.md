---
id: T-1358
type: task
nature: improvement
title: The review page shows a covered tier as skipped, covered by its tier
status: backlog
parent: S-0342
owner: alex
created: 2026-10-08T08:06:34Z
updated: 2026-10-08T08:06:34Z
transitions: []
stream: S-0342
tags: [dashboard]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts]
after: [T-1349]
---
# T-1358 The review page shows a covered tier as skipped, covered by its tier

## Work

The dashboard's review page reads the verify record's steps, and `VerifyState` in `flaiover/src/lib/review.ts` is `'passed' | 'failed' | 'not-reached'`. A record with a `skipped` step would fall back to the raw state and show a duration of nothing run.

- Add `'skipped'` to `VerifyState` and `covered_by?: string` to the step type.
- Give `skipped` a label in `STATE_LABEL`, `skipped: covered by <tier>` when `covered_by` is set, and no duration, as `not-reached` has none.
- `review.test.ts` shows a skipped step's row with its label and an empty duration beside a passed one.
- Run `flai test` on the two files (prettier, eslint, vitest).

It waits for T-1349, which fixes the JSON shape: `state: "skipped"` and `covered_by`.

## Done when

- A verify record with a step `{"state": "skipped", "covered_by": "integration"}` shows as `skipped: covered by integration` with no duration on the review page's rows.
- `flai test flaiover/src/lib/review.ts flaiover/src/lib/review.test.ts` passes.

## Notes

Layer 3 of S-0342's plan, beside T-1355. The story's criteria name the text and `--json` only; this keeps the dashboard's reading of the same record true.

---
id: T-3481
type: task
nature: feature
title: The review page's verification section groups the report's tiers by suite in a VerifySuites component and names the suites the next transitions need that have not passed at the head
status: in-progress
parent: S-0814
owner: alex
created: 2026-11-16T23:51:13Z
updated: 2026-11-28T01:13:48Z
transitions:
  - to: ready
    at: 2026-11-28T01:13:48Z
    by: agent-S-0814
  - to: in-progress
    at: 2026-11-28T01:13:48Z
    by: agent-S-0814
stream: S-0814
tags: [flaiover]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/VerifySuites.svelte, flaiover/src/lib/components/VerifySuites.svelte.test.ts, design/components/flaiover/src/lib/components/VerifySuites.examples.json]
after: [T-3478]
---
# T-3481 The review page's verification section groups the report's tiers by suite in a VerifySuites component and names the suites the next transitions need that have not passed at the head

## Work

Criterion 1 on the page, and the examples of criterion 4. Extend `VerifyReport` and `verifyView` in `review.ts` with what T-3478 adds to `verify.status`. Draw the verification section's tiers in a new `VerifySuites.svelte`: one group per suite, each with its tiers' rows as `verify-step` draws them today, and a line above the groups naming each suite the story's next transitions need, the transition it is for, and whether it passed at the head or why not (not run, failed, passed at another head). `Review.svelte` renders it in place of the inline tier rows, and keeps the steps that are not tiers, the notes, the kluster block, and the running close-out sentence as they are.

Write `VerifySuites.examples.json` with the props of a story whose suites all passed, one with a needed suite not run, one with a suite failed, and a report with no suites.

Waits for T-3478, whose answer this reads.

## Done when

- The section groups tiers by suite and names the needed suites and their state at the head, for each case the examples hold.
- A report with no suites draws each tier as its own suite, as before.
- `VerifySuites.svelte.test.ts`, `review.test.ts`, and `Review.svelte.test.ts` cover these; `flai check` reports no `components.examples` finding for the new file; and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. The section is `verify-section` in `Review.svelte`; `Review.svelte` has no examples file, so the grouped view is a component of its own to carry examples.

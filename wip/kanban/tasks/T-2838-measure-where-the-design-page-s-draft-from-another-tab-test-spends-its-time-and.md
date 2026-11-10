---
id: T-2838
type: task
nature: remediation
title: Measure where the design page's draft-from-another-tab test spends its time, and cut it below vitest's 5,000 ms under load
status: in-progress
parent: S-0655
owner: alex
created: 2026-11-10T07:00:48Z
updated: 2026-11-10T07:21:13Z
transitions:
  - to: ready
    at: 2026-11-10T07:21:13Z
    by: agent-S-0655
  - to: in-progress
    at: 2026-11-10T07:21:13Z
    by: agent-S-0655
stream: S-0655
tags: [flaiover, testing, flaky-test]
touches: [flaiover/src/routes/design/design.svelte.test.ts]
---
# T-2838 Measure where the design page's draft-from-another-tab test spends its time, and cut it below vitest's 5,000 ms under load

## Work

Find what makes "takes up the draft another tab writes, and drops its edits when another tab drops it" in `flaiover/src/routes/design/design.svelte.test.ts` run past 5,000 ms at a load of about 110, and remove that cost. This task waits for nothing and is the first layer.

- Time the test and the other tests of its block, "editing the design system on the design page (S-0383)", with vitest's verbose reporter on the client project. Time them at the host's load and, where you can, under a load like the instance's.
- Split the test's time between its parts: `open()`'s mount and first `settle()`, then each storage event and the `settle()` after it. Record the figures.
- Cut the part that dominates in the test, keeping every assertion. Examples: settle on the condition the next assertion needs instead of a fixed 20 turns, or mount less of the page. S-0640's fix of I-0319 in `components.svelte.test.ts` is the pattern.
- If the time turns out to be the page's own mount, which the test cannot cut, give the slow tests a timeout sized from the figures, and record why in the story's Notes. If the page does work on mount that it should not, say so in the narrative, and change `flaiover/src/routes/design/+page.svelte` only with that reason. `flai task done` widens the touches.
- Write the before and after figures, and the load they were taken at, in the story's Notes under `### Outcome`.

## Done when

- The test's own duration, measured before and after at a stated load, shows the cost I-0330 names cut, or a timeout sized from those figures with the reason recorded.
- Every assertion of the test still holds, and the client project's vitest passes for the file.
- `flai test` passes on the changed paths.

## Notes

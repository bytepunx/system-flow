---
id: T-0204
type: task
nature: improvement
title: Apply the chosen size, update the component test, and check the row in every layout
status: done
parent: S-0054
owner: alex
created: 2026-09-19T05:44:38Z
updated: 2026-09-19T05:55:04Z
transitions:
  - to: ready
    at: 2026-09-19T05:51:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:51:22Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:55:04Z
    by: system-flow
stream: S-0054
tags: []
touches: [flaiover/src/lib/components]
---

# T-0204 Apply the chosen size, update the component test, and check the row in every layout

## Work
Set the chosen size. Update `BoardCard.svelte.test.ts`: the divider is present, every item in the bottom row shares the size class, the parent is still last and pushed right, and the existing cases still pass. Check the bottom row at the narrowest card the layout produces: it stays on one line where it did before, or wraps cleanly with the parent still on the right.

## Done when
- The component tests pass; flaiover lint and svelte-check are clean
- A production build succeeds

## Notes

---
id: T-1882
type: task
nature: feature
title: The review page offers Open in Components for a story whose branch changes a component, opening the page with the story picked
status: done
parent: S-0442
owner: alex
created: 2026-10-16T21:37:38Z
updated: 2026-10-17T03:26:02Z
transitions:
  - to: ready
    at: 2026-10-17T02:56:04Z
    by: agent-S-0442
  - to: in-progress
    at: 2026-10-17T02:56:04Z
    by: agent-S-0442
  - to: done
    at: 2026-10-17T03:26:02Z
    by: agent-S-0442
stream: S-0442
tags: [dashboard]
touches: ["flaiover/src/routes/review/[id]/+page.svelte", flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
after: [T-1881]
usage:
  source: log
  seconds: 1798
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 8217
      cache_read: 1459710
      cache_write: 38398
      cost: 0.6814
---
# T-1882 The review page offers Open in Components for a story whose branch changes a component, opening the page with the story picked

## Work

- `changed.ts`: `changedComponents(diffPaths, catalogue)`, the catalogue's components whose paths the diff changes.
- `review/[id]/+page.svelte`: when the list is not empty, a link "Open in Components" to `/design/components?story=<id>`, with the count.

Waits for nothing in this story.

## Done when

Tests cover a diff that changes components, one that does not, and the link's target. `flai test flaiover/src/lib/explorer flaiover/src/routes/review` passes.

## Notes

None.

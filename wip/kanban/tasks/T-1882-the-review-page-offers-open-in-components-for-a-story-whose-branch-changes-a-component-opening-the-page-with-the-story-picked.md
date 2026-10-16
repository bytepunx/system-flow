---
id: T-1882
type: task
nature: feature
title: The review page offers Open in Components for a story whose branch changes a component, opening the page with the story picked
status: backlog
parent: S-0442
owner: alex
created: 2026-10-16T21:37:38Z
updated: 2026-10-16T21:37:38Z
transitions: []
stream: S-0442
tags: [dashboard]
touches: ["flaiover/src/routes/review/[id]/+page.svelte", flaiover/src/lib/explorer/changed.ts, flaiover/src/lib/explorer/changed.test.ts]
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

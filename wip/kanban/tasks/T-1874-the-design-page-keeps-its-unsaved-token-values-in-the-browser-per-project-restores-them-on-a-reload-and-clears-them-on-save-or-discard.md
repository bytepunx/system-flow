---
id: T-1874
type: task
nature: feature
title: The Design page keeps its unsaved token values in the browser per project, restores them on a reload, and clears them on save or discard
status: backlog
parent: S-0441
owner: alex
created: 2026-10-16T21:36:19Z
updated: 2026-10-16T21:36:19Z
transitions: []
stream: S-0441
tags: [dashboard]
touches: [flaiover/src/lib/design/draft.ts, flaiover/src/lib/design/draft.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
---
# T-1874 The Design page keeps its unsaved token values in the browser per project, restores them on a reload, and clears them on save or discard

## Work

- `draft.ts`: `loadDraft(project)`, `saveDraft(project, edits)`, `clearDraft(project)`, and `onDraft(project, fn)` over `localStorage`, under a key per project, and the `storage` event for other tabs. A value that does not parse is dropped.
- `+page.svelte`: write the edits on each change, restore them on load over the saved tokens, and clear them on save and discard. The leave guard stays.

Waits for nothing.

## Done when

Tests cover the round trip, two projects apart, a corrupt value dropped, and the page restoring and clearing. `flai test flaiover/src/lib/design flaiover/src/routes/design` passes.

## Notes

None.

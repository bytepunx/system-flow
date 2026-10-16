---
id: T-1881
type: task
nature: feature
title: The Components page's source picker chooses main or a story in progress or in review, and marks and lists first the components the story's branch changes
status: backlog
parent: S-0442
owner: alex
created: 2026-10-16T21:37:29Z
updated: 2026-10-16T21:37:29Z
transitions: []
stream: S-0442
tags: [dashboard]
touches: [flaiover/src/lib/components/explorer/SourcePicker.svelte, flaiover/src/lib/components/explorer/SourcePicker.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/routes/api/components/+server.ts, flaiover/src/routes/api/components/start/+server.ts, flaiover/src/routes/api/components/stop/+server.ts, flaiover/src/routes/api/components/components.test.ts, flaiover/src/lib/server/repo.ts]
---
# T-1881 The Components page's source picker chooses main or a story in progress or in review, and marks and lists first the components the story's branch changes

## Work

- `SourcePicker.svelte`: Main and each story in progress or in review, from the board, kept in the URL as `story`.
- The three routes and `repo.ts` pass `story`. `GET /api/components` adds the paths the story's branch changes, from `stream.diff`.
- `+page.svelte`: the picker above the list; with a story, its changed components marked and first.

Waits for nothing in this story; flai's half is T-1880, and the routes pass `story` through either way.

## Done when

Tests cover the picker, the URL, the routes passing `story`, and the changed components ordered first. `flai test flaiover/src/lib/components/explorer flaiover/src/routes/design/components flaiover/src/routes/api/components` passes.

## Notes

None.

---
id: T-1868
type: task
nature: feature
title: The Components page lists the catalogue with a filter, frames the picked component in light and dark at chosen widths, starts the dev server, and shows its errors
status: backlog
parent: S-0439
owner: alex
created: 2026-10-16T21:35:00Z
updated: 2026-10-16T21:35:00Z
transitions: []
stream: S-0439
tags: [dashboard]
touches: [flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/components/explorer/ComponentList.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte]
after: [T-1865, T-1866, T-1867]
---
# T-1868 The Components page lists the catalogue with a filter, frames the picked component in light and dark at chosen widths, starts the dev server, and shows its errors

## Work

- `ComponentList.svelte`: the catalogue by sub-project and folder, with a filter on name and path, and the picked component marked. An unsupported sub-project is listed with its kind.
- `ComponentFrame.svelte`: an `<iframe>` on the harness URL with the dashboard's origin in its query. It connects with T-1866's client, sends `select` and `theme`, sizes itself from `size`, and shows an `error` with its file and line in its place. Draw one frame per theme when both are shown, at the width chosen: 375, 768, 1280, or the panel's.
- `+page.svelte`: the empty states (host action off, server stopped with Start, a host flai too old), `components.status` every minute while the page is open, and the picked component kept in the URL.
- Check by hand on flaiover that an edit to `Gallery.svelte` shows in the frame without a reload, and record it in the story's notes.

Waits for T-1865 for the routes, T-1866 for the client, and T-1867 for the menu entry.

## Done when

Component tests cover the list and filter, each empty state, an error shown in place, and both themes. `flai test flaiover/src/routes/design/components flaiover/src/lib/components/explorer` passes, and the manual check is recorded.

## Notes

None.

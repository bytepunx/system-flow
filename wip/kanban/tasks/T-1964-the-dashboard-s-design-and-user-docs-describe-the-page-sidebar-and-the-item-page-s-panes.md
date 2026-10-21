---
id: T-1964
type: task
nature: feature
title: The dashboard's design and user docs describe the page sidebar and the item page's panes
status: backlog
parent: S-0460
owner: alex
created: 2026-10-21T07:04:00Z
updated: 2026-10-21T07:04:00Z
transitions: []
stream: S-0460
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1963]
---
# T-1964 The dashboard's design and user docs describe the page sidebar and the item page's panes

## Work

- In `design/system/flaiover-dashboard.md`, describe `PageSidebar.svelte` and `$lib/pagesidebar.svelte.ts`: pane and scroll entries, the hash, and the layout below and from `md`. Rewrite the `/items/:id` row as the header over the panes, naming each pane and what it holds.
- In `docs/users/flaiover.md` § The item page, say how to move between panes and that a link keeps the pane.
- Waits for T-1963, so that the docs name the panes as built.

## Done when

Both documents describe the sidebar and the panes, and the markdown tier passes.

## Notes

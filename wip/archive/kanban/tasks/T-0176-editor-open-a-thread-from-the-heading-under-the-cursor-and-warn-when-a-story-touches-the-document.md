---
id: T-0176
type: task
nature: feature
title: "Editor: open a thread from the heading under the cursor, and warn when a story touches the document"
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:30Z
updated: 2026-09-19T02:26:59Z
transitions:
  - to: ready
    at: 2026-09-19T02:23:20Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:23:20Z
    by: alex
  - to: done
    at: 2026-09-19T02:26:59Z
    by: alex
stream: S-0040
tags: []
touches: [flaiover/src/routes, flaiover/src/lib/components]
---

# T-0176 Editor: open a thread from the heading under the cursor, and warn when a story touches the document

## Work
In the editor, work out the heading whose section contains the caret and offer "Open a thread on <heading>", using the existing `Threads` component and API (S-0038) with that heading preselected; show existing threads for the document under the editor. Show the same "being worked on by" warning the document page shows (S-0037) at the top of the editor, from one shared helper rather than a copy, and require one acknowledgement before the first save when it applies.

## Done when
- Tests for the heading-at-caret function and the shared touches helper
- The warning and the thread action appear in the editor, checked in a browser

## Notes

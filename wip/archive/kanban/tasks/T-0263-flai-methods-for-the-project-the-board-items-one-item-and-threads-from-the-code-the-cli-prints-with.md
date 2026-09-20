---
id: T-0263
type: task
nature: feature
title: "flai: methods for the project, the board, items, one item, and threads, from the code the CLI prints with"
status: done
parent: S-0073
owner: alex
created: 2026-09-20T07:52:14Z
updated: 2026-09-20T07:55:54Z
transitions:
  - to: ready
    at: 2026-09-20T07:52:16Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:52:16Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:55:54Z
    by: system-flow
stream: S-0073
tags: []
---
# T-0263 flai: methods for the project, the board, items, one item, and threads, from the code the CLI prints with

## Work
project.info grows the owner, layout, and the dashboard settings the server needs; board.get answers the board view with the fields the board page shows (status, parent title, entered at) added to flai board --json; items.list by type, state, and archive, bodies on request, paths relative to the repository; item.get in any padding with its children; threads.list with entries and by anchor, as flai thread list --json gives it. Arguments are validated as data.

## Done when
- Tests against a scratch project for each method, including the refusals
- flai board --json and flai thread list --json are unchanged except for added fields

## Notes

---
id: T-1964
type: task
nature: feature
title: The dashboard's design and user docs describe the page sidebar and the item page's panes
status: done
parent: S-0460
owner: alex
created: 2026-10-21T07:04:00Z
updated: 2026-10-21T11:38:13Z
transitions:
  - to: ready
    at: 2026-10-21T11:36:40Z
    by: agent-S-0460
  - to: in-progress
    at: 2026-10-21T11:36:41Z
    by: agent-S-0460
  - to: done
    at: 2026-10-21T11:38:13Z
    by: agent-S-0460
stream: S-0460
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1963]
usage:
  source: log
  seconds: 92
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 5457
      cache_read: 2198869
      cache_write: 45277
      cost: 0.9113
---
# T-1964 The dashboard's design and user docs describe the page sidebar and the item page's panes

## Work

- In `design/system/flaiover-dashboard.md`, describe `PageSidebar.svelte` and `$lib/pagesidebar.svelte.ts`: pane and scroll entries, the hash, and the layout below and from `md`. Rewrite the `/items/:id` row as the header over the panes, naming each pane and what it holds.
- In `docs/users/flaiover.md` § The item page, say how to move between panes and that a link keeps the pane.
- Waits for T-1963, so that the docs name the panes as built.

## Done when

Both documents describe the sidebar and the panes, and the markdown tier passes.

## Notes

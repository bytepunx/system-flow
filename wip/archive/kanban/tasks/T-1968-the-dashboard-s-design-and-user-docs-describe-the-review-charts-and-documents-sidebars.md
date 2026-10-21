---
id: T-1968
type: task
nature: improvement
title: The dashboard's design and user docs describe the review, charts, and documents sidebars
status: done
parent: S-0461
owner: alex
created: 2026-10-21T07:04:36Z
updated: 2026-10-21T12:48:37Z
transitions:
  - to: ready
    at: 2026-10-21T12:18:35Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T12:18:36Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:48:37Z
    by: agent-S-0461
stream: S-0461
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1965, T-1966, T-1967, T-1984]
usage:
  source: log
  seconds: 1801
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 93
      output: 43360
      cache_read: 5539086
      cache_write: 192100
      cost: 3.1129
---
# T-1968 The dashboard's design and user docs describe the review, charts, and documents sidebars

## Work

- In `design/system/flaiover-dashboard.md` § Views, add each sidebar and its entries to the `/review/:id`, `/charts/:kind`, and `/docs/:path` rows.
- In `docs/users/flaiover.md`, add one or two sentences each to § Reviewing a story, § Charts, and § Docs on what the sidebar holds.
- Waits for T-1965, T-1966, and T-1967, so that it describes all three pages as built.

## Done when

Both documents describe the three sidebars, and the markdown tier passes.

## Notes

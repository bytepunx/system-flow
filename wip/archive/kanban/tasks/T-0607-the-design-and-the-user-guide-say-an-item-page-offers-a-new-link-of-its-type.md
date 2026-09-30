---
id: T-0607
type: task
nature: improvement
title: The design and the user guide say an item page offers a New link of its type
status: done
parent: S-0171
owner: alex
created: 2026-09-30T01:12:22Z
updated: 2026-09-30T01:16:44Z
transitions:
  - to: ready
    at: 2026-09-30T01:12:25Z
    by: agent-S-0171
  - to: in-progress
    at: 2026-09-30T01:16:08Z
    by: agent-S-0171
  - to: done
    at: 2026-09-30T01:16:44Z
    by: agent-S-0171
stream: S-0171
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users]
usage:
  source: log
  seconds: 36
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 4751
      cache_read: 1009877
      cache_write: 25913
      cost: 0.5044
---
# T-0607 The design and the user guide say an item page offers a New link of its type

## Work

- `design/system/flaiover-dashboard.md`: the item view names the `New story` / `New epic` link and what it starts the form with.
- The user guide for the dashboard under `docs/users/` says the same where it describes the item page.

## Done when

- Both documents describe the link, with `updated` bumped, and `flai check --strict` passes.

## Notes

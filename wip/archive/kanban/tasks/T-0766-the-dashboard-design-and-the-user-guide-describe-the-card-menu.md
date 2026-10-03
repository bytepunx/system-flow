---
id: T-0766
type: task
nature: feature
title: The dashboard design and the user guide describe the card menu
status: done
parent: S-0202
owner: alex
created: 2026-10-03T18:32:22Z
updated: 2026-10-03T18:45:47Z
transitions:
  - to: ready
    at: 2026-10-03T18:32:55Z
    by: agent-S-0202
  - to: in-progress
    at: 2026-10-03T18:44:58Z
    by: agent-S-0202
  - to: done
    at: 2026-10-03T18:45:47Z
    by: agent-S-0202
stream: S-0202
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0765]
usage:
  source: log
  seconds: 49
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 8860
      cache_read: 1389417
      cache_write: 35605
      cost: 0.6946
---

# T-0766 The dashboard design and the user guide describe the card menu

## Work

Describe the card menu in `design/system/flaiover-dashboard.md` (the `/board` row) and in `docs/users/flaiover.md` (the board, beside the lane menu), as T-0765 built it. Waits for T-0765, so that the documents say what was built.

## Done when

- Both documents describe how the menu opens, its entries and when each shows, what each runs, and how it closes, with `updated` bumped.

## Notes

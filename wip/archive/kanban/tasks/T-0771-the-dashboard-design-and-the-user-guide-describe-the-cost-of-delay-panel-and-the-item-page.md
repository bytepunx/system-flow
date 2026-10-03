---
id: T-0771
type: task
nature: feature
title: The dashboard design and the user guide describe the Cost of delay panel and the item page
status: done
parent: S-0204
owner: alex
created: 2026-10-03T18:35:53Z
updated: 2026-10-03T19:31:35Z
transitions:
  - to: ready
    at: 2026-10-03T18:36:19Z
    by: agent-S-0204
  - to: in-progress
    at: 2026-10-03T19:30:50Z
    by: agent-S-0204
  - to: done
    at: 2026-10-03T19:31:35Z
    by: agent-S-0204
stream: S-0204
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0768, T-0770]
usage:
  source: log
  seconds: 45
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 32
      output: 8553
      cache_read: 1828267
      cache_write: 34106
      cost: 0.7469
---
# T-0771 The dashboard design and the user guide describe the Cost of delay panel and the item page

## Work

`design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the Cost of delay panel on the new-item and edit forms (collapsed by default, the three inputs in the project's currency, the summary line, clearing) and what the item page shows (inputs, the planner's value with who set it and when, the stale note). Waits for T-0768 and T-0770, whose behaviour it describes.

## Done when

- Both documents describe the panel and the item page as built, with `updated` bumped.

## Notes

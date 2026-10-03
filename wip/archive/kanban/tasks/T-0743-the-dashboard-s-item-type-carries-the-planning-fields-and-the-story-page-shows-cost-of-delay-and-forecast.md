---
id: T-0743
type: task
nature: feature
title: The dashboard's item type carries the planning fields and the story page shows cost of delay and forecast
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:26Z
updated: 2026-10-03T06:20:03Z
transitions:
  - to: ready
    at: 2026-10-03T06:12:24Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T06:12:24Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T06:20:03Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [flaiover/src]
after: [T-0742]
usage:
  source: log
  seconds: 459
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 17994
      cache_read: 2909520
      cache_write: 73400
      cost: 1.3749
---
# T-0743 The dashboard's item type carries the planning fields and the story page shows cost of delay and forecast

## Work

Add `draft`, `cost_of_delay`, and `forecast` to the item types in `flaiover/src`, and show cost of delay (inputs, value, who set it and when) and the forecast (duration, delivery, basis, who and when) beside `estimate` on an item's page, read-only. The `[Draft]` indicator and Finalize are S-0201's, and the forms S-0204's. Waits for T-0742, whose hostapi reads it shows.

## Done when

The item page shows both blocks when set and nothing when not, with a test, and `npm run check` and the tests pass.

## Notes

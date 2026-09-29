---
id: T-0590
type: task
nature: feature
title: The charts plot what the window holds, on a time axis that spans it
status: done
parent: S-0166
owner: alex
created: 2026-09-29T22:47:47Z
updated: 2026-09-29T22:55:36Z
transitions:
  - to: ready
    at: 2026-09-29T22:48:15Z
    by: agent-S-0166
  - to: in-progress
    at: 2026-09-29T22:51:43Z
    by: agent-S-0166
  - to: done
    at: 2026-09-29T22:55:36Z
    by: agent-S-0166
stream: S-0166
tags: []
touches: [flaiover/src/lib/viz, flaiover/src/routes/charts]
usage:
  source: log
  seconds: 233
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 64
      output: 19026
      cache_read: 4703464
      cache_write: 74301
      cost: 1.9159
---

# T-0590 The charts plot what the window holds, on a time axis that spans it

## Work

In `flaiover/src/lib/viz/charts.ts`, the charts per item (cycle time, time in state, estimates, cost by item) and their tables plot only the items completed in the report's window. Every time axis (cycle time, burn-up, cumulative flow, completion over time, and the charts of spend over time) runs from the window's start to the report's now, not from the first to the last point. The charts page ignores an answer to a window it no longer shows.

## Done when

- Behaviour tests show that a narrower window drops the items outside it and moves each time axis's ends to the window's.
- The page test shows that changing the window asks flai for it and redraws with the new axis.
- `make flaiover-test` (lint, type check, unit tests) passes.

## Notes

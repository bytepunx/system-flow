---
id: T-0596
type: task
nature: feature
title: The charts page keeps the window chosen and loads it at once
status: done
parent: S-0168
owner: alex
created: 2026-09-29T23:58:07Z
updated: 2026-09-30T00:01:07Z
transitions:
  - to: ready
    at: 2026-09-29T23:58:27Z
    by: agent-S-0168
  - to: in-progress
    at: 2026-09-29T23:58:27Z
    by: agent-S-0168
  - to: done
    at: 2026-09-30T00:01:07Z
    by: agent-S-0168
stream: S-0168
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 160
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 10747
      cache_read: 2661842
      cache_write: 41966
      cost: 1.0832
---

# T-0596 The charts page keeps the window chosen and loads it at once

## Work

- Remember the window chosen on the charts page in the browser, so a remount (leaving Charts and coming back, a reload) opens every chart at it rather than at 30d.
- Ask flai for the stats at mount without waiting for the epic list, and draw the charts even when the epic list cannot be read.
- Page tests for both.

## Done when

- The charts page test shows a remount opening cycle time and time in state at the window chosen before, and the stats asked for even when the epic list fails.
- `scripts/flaiover-test.sh` passes.

## Notes

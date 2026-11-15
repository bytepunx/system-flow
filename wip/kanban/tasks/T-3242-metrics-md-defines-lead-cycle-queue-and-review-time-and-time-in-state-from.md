---
id: T-3242
type: task
nature: feature
title: metrics.md defines lead, cycle, queue, and review time and time in state from their timers and links the ADR
status: done
parent: S-0750
owner: alex
created: 2026-11-12T16:09:01Z
updated: 2026-11-15T21:12:35Z
transitions:
  - to: ready
    at: 2026-11-15T21:12:26Z
    by: agent-S-0750
  - to: in-progress
    at: 2026-11-15T21:12:26Z
    by: agent-S-0750
  - to: done
    at: 2026-11-15T21:12:35Z
    by: agent-S-0750
stream: S-0750
tags: [cli]
touches: [design/system/metrics.md]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 1
      output: 12
      cache_read: 155204
      cache_write: 227
      cost: 0.0588
---
# T-3242 metrics.md defines lead, cycle, queue, and review time and time in state from their timers and links the ADR

## Work

Waits for T-3236: it states what that ADR decides. In `design/system/metrics.md` § Derived timestamps per item and § Per item durations, define each figure as its baseline timer, its points or states, and link the ADR. Keep the numbers, the precision rules, and the data access as they are.

## Done when

`metrics.md` defines each of the five figures from its timer.

## Notes

---
id: T-2064
type: task
nature: feature
title: An ADR settles that a story's queue time is the seconds it spent in ready before it started, so a stretch back in backlog is not counted as waiting
status: done
parent: S-0469
owner: alex
created: 2026-10-23T06:05:41Z
updated: 2026-10-23T06:06:48Z
transitions:
  - to: ready
    at: 2026-10-23T06:05:56Z
    by: agent-S-0469
  - to: in-progress
    at: 2026-10-23T06:05:56Z
    by: agent-S-0469
  - to: done
    at: 2026-10-23T06:06:48Z
    by: agent-S-0469
stream: S-0469
tags: []
touches: [design/system/metrics.md]
usage:
  source: log
  seconds: 52
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4169
      cache_read: 1169363
      cache_write: 31945
      cost: 0.5729
---
# T-2064 An ADR settles that a story's queue time is the seconds it spent in ready before it started, so a stretch back in backlog is not counted as waiting

## Work

Record with `flai adr new --commit` that queue time is the seconds an item spent in `ready` before `started`, the sum of its `ready` intervals up to its first move to `in-progress`, instead of `started - committed`. A story the operator sends back from `ready` to `backlog` and makes ready again no longer counts the stretch in backlog as waiting, so queue time minus `held_seconds` is the unheld wait in ready, as I-0153's analysis read it. Update the Queue time row and its precision rule in `design/system/metrics.md`, linking the ADR. `committed` stays the first move to `ready`.

## Done when

- The ADR is committed on `story/S-0469` and indexed.
- `design/system/metrics.md` defines queue time as the time in `ready` before `started`, with the ADR linked.

## Notes

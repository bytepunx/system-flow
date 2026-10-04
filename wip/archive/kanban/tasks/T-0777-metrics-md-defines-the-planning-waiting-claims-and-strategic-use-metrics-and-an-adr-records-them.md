---
id: T-0777
type: task
nature: feature
title: metrics.md defines the planning, waiting, claims, and strategic-use metrics, and an ADR records them
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:39:56Z
updated: 2026-10-03T20:41:50Z
transitions:
  - to: ready
    at: 2026-10-03T20:40:14Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T20:40:15Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T20:41:50Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [design/system/metrics.md, design/adrs]
usage:
  source: log
  seconds: 95
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 8050
      cache_read: 1162410
      cache_write: 31467
      cost: 0.5821
---
# T-0777 metrics.md defines the planning, waiting, claims, and strategic-use metrics, and an ADR records them

## Work

Write the contract first, so the code tasks implement it: a section of `design/system/metrics.md` per group (forecasts and estimates, cost of delay, waiting, claims and touches, strategic use per day), each with its JSON keys and precision rules, and an ADR recording the additions and the choices a reader could have made differently. Waits for nothing: every other task implements what it writes.

## Done when

- [ ] `metrics.md` defines every value the story's criteria name, with its JSON key and precision
- [ ] An ADR made with `flai adr new` records the additions, and `metrics.md` links it

## Notes

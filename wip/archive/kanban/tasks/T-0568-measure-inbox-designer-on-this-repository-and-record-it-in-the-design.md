---
id: T-0568
type: task
nature: feature
title: Measure inbox.designer on this repository and record it in the design
status: done
parent: S-0158
owner: alex
created: 2026-09-29T19:30:30Z
updated: 2026-09-29T19:33:12Z
transitions:
  - to: ready
    at: 2026-09-29T19:30:43Z
    by: agent-S-0158
  - to: in-progress
    at: 2026-09-29T19:33:11Z
    by: agent-S-0158
  - to: done
    at: 2026-09-29T19:33:12Z
    by: agent-S-0158
stream: S-0158
tags: []
usage:
  source: log
  seconds: 1
  estimated: true
  models: []
---

# T-0568 Measure inbox.designer on this repository and record it in the design

## Work

Time `inbox.designer` on this repository with `flai hostapi --timing`, three runs, before and after the change, and record the median and its phases in `design/system/server-performance.md`.

## Done when

- `inbox.designer` on this repository takes under 20 ms, and its timing has no `check.run` phase.
- `design/system/server-performance.md` records the new number against cause 3.

## Notes

---
id: T-0591
type: task
nature: feature
title: The design, the ADR, and the user guide say each chart spans the window
status: done
parent: S-0166
owner: alex
created: 2026-09-29T22:47:47Z
updated: 2026-09-29T22:57:21Z
transitions:
  - to: ready
    at: 2026-09-29T22:48:15Z
    by: agent-S-0166
  - to: in-progress
    at: 2026-09-29T22:55:36Z
    by: agent-S-0166
  - to: done
    at: 2026-09-29T22:57:21Z
    by: agent-S-0166
stream: S-0166
tags: []
touches: [design/system/metrics.md, design/system/flaiover-dashboard.md, design/adrs, docs/users]
usage:
  source: log
  seconds: 105
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 11527
      cache_read: 2849448
      cache_write: 45013
      cost: 1.1607
---

# T-0591 The design, the ADR, and the user guide say each chart spans the window

## Work

Record with `flai adr new` that every chart spans the report's window, since `design/system/metrics.md` changes only with an ADR. Update the charts table and precision rules in `design/system/metrics.md`, the charts section of `design/system/flaiover-dashboard.md`, and the user guide where it describes the window.

## Done when

- The ADR is in `design/adrs` and linked from `design/system/metrics.md`.
- `flai check --strict` passes.

## Notes

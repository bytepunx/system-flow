---
id: S-0168
type: story
nature: remediation
title: Several Charts Don't Use Window correctly
status: done
parent: E-0013
owner: alex
created: 2026-09-29T23:43:11Z
updated: 2026-09-30T00:19:31Z
transitions:
  - to: ready
    at: 2026-09-29T23:43:21Z
    by: alex
  - to: in-progress
    at: 2026-09-29T23:51:58Z
    by: agent-S-0168
  - to: review
    at: 2026-09-30T00:18:35Z
    by: agent-S-0168
  - to: done
    at: 2026-09-30T00:19:31Z
    by: alex
tags: [dashboard]
topics: [client-side-charts]
touches: [flaiover/src, design/adrs/0056-time-in-state-is-one-stacked-bar-per-day-of-the-window-the-mean-hours-per-state.md, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/metrics.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1489
  models:
    - model: claude-opus-5-5
      input: 240
      output: 60973
      cache_read: 16926002
      cache_write: 203399
      cost: 6.2328
---
# S-0168 Several Charts Don't Use Window correctly

## Goal

All charts use the window selection when the user changes to them and all charts correctly use the window selection on change.

## Acceptance criteria
- [x] the cycle time chart loads with the correct window and changes to the right time range in the X axis
- [x] the time in state chart loads with the correct window and changes to the right time range in the X axis

## Tasks
- T-0596 The charts page keeps the window chosen and loads it at once
- T-0597 Time in state is drawn on a time axis that spans the window

## Notes

- Cycle time already followed the window once drawn (S-0166); it opened at 30d whenever the charts page remounted, and waited on the epic list. Both fixed in T-0596.
- Time in state is one stacked bar per day on a time axis (ADR-0056, TH-0039).

---
id: S-0168
type: story
nature: remediation
title: Several Charts Don't Use Window correctly
status: in-progress
parent: E-0013
owner: alex
created: 2026-09-29T23:43:11Z
updated: 2026-09-29T23:58:07Z
transitions:
  - to: ready
    at: 2026-09-29T23:43:21Z
    by: alex
  - to: in-progress
    at: 2026-09-29T23:51:58Z
    by: agent-S-0168
tags: [dashboard]
topics: [client-side-charts]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1207
  models:
    - model: claude-opus-5-5
      input: 174
      output: 42081
      cache_read: 10422889
      cache_write: 164326
      cost: 4.2415
---
# S-0168 Several Charts Don't Use Window correctly

## Goal

All charts use the window selection when the user changes to them and all charts correctly use the window selection on change.

## Acceptance criteria
- [ ] the cycle time chart loads with the correct window and changes to the right time range in the X axis
- [ ] the time in state chart loads with the correct window and changes to the right time range in the X axis

## Tasks
- T-0596 The charts page keeps the window chosen and loads it at once
- T-0597 Time in state is drawn on a time axis that spans the window

## Notes

---
id: S-0452
type: story
nature: improvement
title: Keep the top menu visible at all times
status: cancelled
owner: alex
created: 2026-10-20T20:56:45Z
updated: 2026-10-20T21:19:38Z
transitions:
  - to: ready
    at: 2026-10-20T20:56:48Z
    by: alex
  - to: cancelled
    at: 2026-10-20T21:19:38Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 925
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 1118
          cache_read: 13582013
          cache_write: 20244
          cost: 3.4101
---
# S-0452 Keep the top menu visible at all times

## Goal

By default, the top menu and submenu should stay visible at all times. A setting in host configuration should allow the operator to turn this off, but it should be on by default.

## Acceptance criteria
- [ ] When scrolling, the menu bar should be visible at all times
- [ ] When an operator uses host settings to turn this off, the menu no longer stays pinned to the top of the page

## Tasks

## Notes
- 2026-10-20T21:19:38Z: moved to cancelled: Lets look at a different screen layout instead

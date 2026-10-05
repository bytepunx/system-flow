---
id: S-0263
type: story
nature: improvement
title: "Add \"Plan\" to the context menu for epic and story cards"
status: done
owner: alex
created: 2026-10-04T23:38:12Z
updated: 2026-10-05T00:03:14Z
transitions:
  - to: ready
    at: 2026-10-04T23:38:13Z
    by: alex
  - to: in-progress
    at: 2026-10-04T23:38:29Z
    by: agent-S-0263
  - to: review
    at: 2026-10-04T23:47:53Z
    by: agent-S-0263
  - to: done
    at: 2026-10-05T00:03:14Z
    by: alex
tags: [dashboard]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 610
  models:
    - model: claude-opus-5-5
      input: 168
      output: 44585
      cache_read: 6410202
      cache_write: 270435
      cost: 3.9868
    - model: claude-sonnet-5
      input: 34
      output: 7554
      cache_read: 677839
      cache_write: 60117
      cost: 0.3615
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: alex
    at: 2026-10-04T23:38:12Z
---
# S-0263 Add "Plan" to the context menu for epic and story cards

## Goal

Add the Plan action to an epic or story's right-click context menu so that the operator can kick off a plan step for unplanned stories.

## Acceptance criteria
- [x] Plan appears as a context menu option for epic or story cards on the board
- [x] Selecting `Plan` on an eligible epic or story starts the planning agent for that card's item.

## Tasks
- T-0834 The card menu offers Plan on an open epic or story while the plan host action is on, and starts the planner
- T-0835 The design and the user guide say the card menu offers Plan

## Notes

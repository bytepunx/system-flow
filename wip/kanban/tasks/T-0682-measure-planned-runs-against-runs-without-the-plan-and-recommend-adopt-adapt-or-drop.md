---
id: T-0682
type: task
nature: experiment
title: Measure planned runs against runs without the plan and recommend adopt, adapt, or drop
status: done
parent: S-0176
owner: arobson
created: 2026-10-01T11:42:11Z
updated: 2026-10-01T13:14:46Z
transitions:
  - to: ready
    at: 2026-10-01T12:51:50Z
    by: agent-S-0176
  - to: in-progress
    at: 2026-10-01T12:51:50Z
    by: agent-S-0176
  - to: done
    at: 2026-10-01T13:14:46Z
    by: agent-S-0176
stream: S-0176
tags: []
touches: [design/system/agent-context.md, design/adrs]
usage:
  source: log
  seconds: 1376
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 90
      output: 29009
      cache_read: 6410801
      cache_write: 91478
      cost: 2.4292
    - model: claude-sonnet-5-5
      input: 53
      output: 11641
      cache_read: 1008813
      cache_write: 105171
      cost: 0.5812
---
# T-0682 Measure planned runs against runs without the plan and recommend adopt, adapt, or drop

## Work

- Measure the runs TH-0059 settles on, with the plan, against comparable runs without it, as `design/system/agent-context.md § Sub-agents` measures them (ADR-0051, sub-agents included): wall-clock time to review, cost and tokens, the number of layers, conflicts, and review defects.
- Record the table, what it shows against the hypothesis (a third less time to review, at most half again the cost, no more review defects), and a recommendation: adopt, adapt, or drop. With adopt, an ADR.

## Done when

- agent-context.md has the table, the reading, and the recommendation, and an ADR exists if the recommendation is to adopt.

## Notes

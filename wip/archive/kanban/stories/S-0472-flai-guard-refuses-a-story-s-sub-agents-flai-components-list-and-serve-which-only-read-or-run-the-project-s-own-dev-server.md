---
id: S-0472
type: story
nature: improvement
title: flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server
status: done
owner: alex
created: 2026-10-21T11:58:44Z
updated: 2026-10-23T17:53:26Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:27Z
    by: alex
  - to: in-progress
    at: 2026-10-23T06:21:24Z
    by: agent-S-0472
  - to: review
    at: 2026-10-23T06:23:18Z
    by: agent-S-0472
  - to: done
    at: 2026-10-23T17:53:26Z
    by: alex
tags: []
touches: [design/issues/I-0156-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 125
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 3
      work: 14
  models:
    - model: claude-opus-5-5
      input: 36
      output: 5261
      cache_read: 1203905
      cache_write: 69405
      cost: 0.9014
  strategic:
    - kind: orchestrator
      seconds: 1137
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 166
          output: 2773
          cache_read: 48769653
          cache_write: 39173
          cost: 10.7774
        - model: claude-sonnet-5-5
          input: 14
          output: 114
          cache_read: 125511
          cache_write: 76517
          cost: 0.1399
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T11:58:44Z
finalized:
  by: alex
  at: 2026-10-22T20:53:07Z
---
# S-0472 flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server

## Goal

This story remediates [I-0156](../../../design/issues/I-0156-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md), "flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0156 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0156 is closed with `flai issue close I-0156 --reason` saying what fixed it

## Tasks
- T-2067 I-0156 is closed as fixed by S-0477, whose guard test reproduces the refusal

## Notes

Cost of delay inputs set by flai from I-0156. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T00:41:49Z, 4.5 days before this story; under one cycle counts as one).

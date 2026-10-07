---
id: S-0312
type: story
nature: improvement
title: flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart
status: backlog
owner: alex
created: 2026-10-07T14:26:01Z
updated: 2026-10-07T14:26:01Z
transitions: []
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
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 32
          cache_read: 568691
          cache_write: 4204
          cost: 0.1492
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T14:26:01Z
---
# S-0312 flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart

## Goal

This story remediates [I-0104](../../../design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md), "flai task done commits everything in the worktree, so two tasks of one layer cannot be closed apart". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0104 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0104 is closed with `flai issue close I-0104 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0104. time_lost_per_cycle 10m: 5m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-07T07:17:53Z, 0.3 days before this story; under one cycle counts as one).

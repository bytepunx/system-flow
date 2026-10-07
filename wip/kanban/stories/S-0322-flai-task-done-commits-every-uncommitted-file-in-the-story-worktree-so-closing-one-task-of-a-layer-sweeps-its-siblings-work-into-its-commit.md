---
id: S-0322
type: story
nature: improvement
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
status: backlog
owner: alex
created: 2026-10-07T18:59:53Z
updated: 2026-10-07T18:59:53Z
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
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 64217
          cache_write: 5432
          cost: 0.0172
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T18:59:53Z
---
# S-0322 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Goal

This story remediates [I-0108](../../../design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md), "flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0108 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0108 is closed with `flai issue close I-0108 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0108. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T08:06:26Z, 0.5 days before this story; under one cycle counts as one).

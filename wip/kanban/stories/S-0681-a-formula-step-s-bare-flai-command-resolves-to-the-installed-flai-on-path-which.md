---
id: S-0681
type: story
nature: improvement
title: A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use
status: ready
owner: alex
created: 2026-11-10T08:37:05Z
updated: 2026-11-10T10:58:14Z
transitions:
  - to: ready
    at: 2026-11-10T10:58:14Z
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
      seconds: 95
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 35
          output: 260
          cache_read: 1365229
          cache_write: 153744
          cost: 0.4283
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T08:37:05Z
finalized:
  by: alex
  at: 2026-11-10T10:58:10Z
---
# S-0681 A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use

## Goal

This story remediates [I-0332](../../../design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md), "A formula step's bare flai command resolves to the installed flai on PATH, which can be older than the flags the formula's steps use". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0332 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0332 is closed with `flai issue close I-0332 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0332. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T08:28:53Z, 0 days before this story; under one cycle counts as one).

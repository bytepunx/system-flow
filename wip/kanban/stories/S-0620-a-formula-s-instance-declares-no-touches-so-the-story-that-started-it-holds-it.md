---
id: S-0620
type: story
nature: remediation
title: A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress
status: backlog
owner: alex
created: 2026-11-07T23:39:44Z
updated: 2026-11-07T23:39:44Z
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
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 69952
          cache_write: 17190
          cost: 0.0225
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-07T23:39:44Z
---
# S-0620 A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress

## Goal

This story remediates [I-0282](../../../design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md), "A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0282 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0282 is closed with `flai issue close I-0282 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0282. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:04:18Z, 0.6 days before this story; under one cycle counts as one).

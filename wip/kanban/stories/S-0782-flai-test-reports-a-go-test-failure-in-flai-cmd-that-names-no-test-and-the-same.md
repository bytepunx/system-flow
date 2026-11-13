---
id: S-0782
type: story
nature: remediation
title: flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again
status: backlog
owner: alex
created: 2026-11-13T22:46:52Z
updated: 2026-11-13T22:46:52Z
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
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 39
          cache_read: 183456
          cache_write: 22472
          cost: 0.0607
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-11-13T22:46:52Z
---
# S-0782 flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Goal

This story remediates [I-0365](../../../design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md), "flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0365 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0365 is closed with `flai issue close I-0365 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0365. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:17:49Z, 0.1 days before this story; under one cycle counts as one).

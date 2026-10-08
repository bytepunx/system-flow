---
id: S-0327
type: story
nature: improvement
title: flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named
status: cancelled
owner: alex
created: 2026-10-07T18:59:59Z
updated: 2026-10-08T04:38:46Z
transitions:
  - to: cancelled
    at: 2026-10-08T04:38:46Z
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
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 18
          cache_read: 378711
          cache_write: 5983
          cost: 0.0948
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-07T18:59:59Z
---
# S-0327 flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named

## Goal

This story remediates [I-0114](../../../design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md), "flai verify's integration tier keeps only the last lines of go test's output, so the failing test is not named". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0114 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0114 is closed with `flai issue close I-0114 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0114. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T09:20:56Z, 0.4 days before this story; under one cycle counts as one).
- 2026-10-08T04:38:46Z: moved to cancelled: Duplicate of 313

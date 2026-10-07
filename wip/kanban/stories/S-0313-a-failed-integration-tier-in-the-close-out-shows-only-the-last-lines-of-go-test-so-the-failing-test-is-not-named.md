---
id: S-0313
type: story
nature: improvement
title: A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
status: backlog
owner: alex
created: 2026-10-07T14:26:02Z
updated: 2026-10-07T14:26:02Z
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
          input: 3
          output: 28
          cache_read: 552193
          cache_write: 4062
          cost: 0.1451
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-07T14:26:02Z
---
# S-0313 A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Goal

This story remediates [I-0113](../../../design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md), "A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0113 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0113 is closed with `flai issue close I-0113 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0113. time_lost_per_cycle 4m: 4m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T09:13:42Z, 0.2 days before this story; under one cycle counts as one).

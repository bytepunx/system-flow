---
id: S-0316
type: story
nature: remediation
title: The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use
status: backlog
owner: alex
created: 2026-10-07T18:59:46Z
updated: 2026-10-07T18:59:46Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-07T18:59:46Z
---
# S-0316 The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use

## Goal

This story remediates [I-0094](../../../design/issues/I-0094-the-dashboard-upgrade-stops-the-old-container-and-cannot-start-the-new-one-because-the-name-flaiover-is-still-in-use.md), "The dashboard upgrade stops the old container and cannot start the new one, because the name flaiover is still in use". The issue recommends this solution:

Directions to weigh: after stopping the old container, wait until Docker no longer lists the name before running the new one, with a short limit; or run the new container under a temporary name and rename it; and when the start fails, start the previous image again so that the operator is not left without a dashboard. A test with a stand-in for docker that keeps the name for a moment after `stop` would reproduce it.

## Acceptance criteria
- [ ] The cause I-0094 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0094 is closed with `flai issue close I-0094 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0094. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T21:01:06Z, 0.9 days before this story; under one cycle counts as one).

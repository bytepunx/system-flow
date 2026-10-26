---
id: S-0500
type: story
nature: improvement
title: Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer
status: backlog
owner: alex
created: 2026-10-26T04:01:29Z
updated: 2026-10-26T04:01:29Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-10-26T04:01:29Z
---
# S-0500 Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer

## Goal

This story remediates [I-0184](../../../design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md), "Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0184 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0184 is closed with `flai issue close I-0184 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0184. time_lost_per_cycle 1s: 450ms per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T12:50:20Z, 4.6 days before this story; under one cycle counts as one).

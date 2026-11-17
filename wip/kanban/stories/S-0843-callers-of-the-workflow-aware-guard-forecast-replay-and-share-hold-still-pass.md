---
id: S-0843
type: story
nature: remediation
title: Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them
status: backlog
owner: alex
created: 2026-11-17T08:32:54Z
updated: 2026-11-17T08:32:54Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:54Z
---
# S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them

## Goal

This story remediates [I-0379](../../../design/issues/I-0379-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md), "Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0379 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0379 is closed with `flai issue close I-0379 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0379. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T20:04:03Z, 1.5 days before this story; under one cycle counts as one).

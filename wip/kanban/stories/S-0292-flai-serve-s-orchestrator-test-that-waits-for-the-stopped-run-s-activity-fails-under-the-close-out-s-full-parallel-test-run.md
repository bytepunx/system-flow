---
id: S-0292
type: story
nature: remediation
title: flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run
status: backlog
owner: alex
created: 2026-10-06T11:14:34Z
updated: 2026-10-06T11:14:34Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-06T11:14:34Z
---
# S-0292 flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run

## Goal

This story remediates [I-0090](../../../design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md), "flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0090 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0090 is closed with `flai issue close I-0090 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0090. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T11:05:36Z, 0 days before this story; under one cycle counts as one).

---
id: S-0320
type: story
nature: remediation
title: flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
status: backlog
owner: alex
created: 2026-10-07T18:59:51Z
updated: 2026-10-07T18:59:51Z
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
    at: 2026-10-07T18:59:51Z
---
# S-0320 flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run

## Goal

This story remediates [I-0106](../../../design/issues/I-0106-flai-serve-s-test-that-the-orchestrator-is-started-again-when-it-ends-fails-under-the-close-out-s-full-integration-run.md), "flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0106 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0106 is closed with `flai issue close I-0106 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0106. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T07:31:02Z, 0.5 days before this story; under one cycle counts as one).

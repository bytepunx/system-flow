---
id: S-0310
type: story
nature: remediation
title: TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
status: backlog
owner: alex
created: 2026-10-07T06:48:45Z
updated: 2026-10-07T06:48:45Z
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
    at: 2026-10-07T06:48:45Z
---
# S-0310 TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded

## Goal

This story remediates [I-0102](../../../design/issues/I-0102-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md), "TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0102 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0102 is closed with `flai issue close I-0102 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0102. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T03:13:07Z, 0.1 days before this story; under one cycle counts as one).

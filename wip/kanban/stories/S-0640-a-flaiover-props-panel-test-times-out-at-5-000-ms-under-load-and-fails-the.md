---
id: S-0640
type: story
nature: remediation
title: A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration
status: backlog
owner: alex
created: 2026-11-09T16:55:57Z
updated: 2026-11-09T16:55:57Z
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
    at: 2026-11-09T16:55:57Z
---
# S-0640 A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Goal

This story remediates [I-0319](../../../design/issues/I-0319-a-flaiover-props-panel-test-times-out-at-5-000-ms-under-load-and-fails-the.md), "A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0319 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0319 is closed with `flai issue close I-0319 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0319. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T10:34:28Z, 0.3 days before this story; under one cycle counts as one).

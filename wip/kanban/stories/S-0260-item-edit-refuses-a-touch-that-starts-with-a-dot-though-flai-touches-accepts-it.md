---
id: S-0260
type: story
nature: remediation
title: item_edit refuses a touch that starts with a dot, though flai touches accepts it
status: ready
owner: alex
created: 2026-10-04T06:29:48Z
updated: 2026-10-04T23:14:26Z
transitions:
  - to: ready
    at: 2026-10-04T23:14:26Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-04T06:29:48Z
finalized:
  by: alex
  at: 2026-10-04T23:14:24Z
---
# S-0260 item_edit refuses a touch that starts with a dot, though flai touches accepts it

## Goal

This story remediates [I-0071](../../../design/issues/I-0071-item-edit-refuses-a-touch-that-starts-with-a-dot-though-flai-touches-accepts-it.md), "item_edit refuses a touch that starts with a dot, though flai touches accepts it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0071 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0071 is closed with `flai issue close I-0071 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0071. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-04T04:45:59Z, 0.1 days before this story; under one cycle counts as one).

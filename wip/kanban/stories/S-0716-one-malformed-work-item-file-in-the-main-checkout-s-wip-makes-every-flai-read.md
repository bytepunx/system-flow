---
id: S-0716
type: story
nature: remediation
title: One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them
status: backlog
owner: alex
created: 2026-11-11T01:40:46Z
updated: 2026-11-11T01:40:46Z
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
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-11T01:40:46Z
---
# S-0716 One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them

## Goal

This story remediates [I-0344](../../../design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md), "One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0344 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0344 is closed with `flai issue close I-0344 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0344. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T11:05:08Z, 0.6 days before this story; under one cycle counts as one).

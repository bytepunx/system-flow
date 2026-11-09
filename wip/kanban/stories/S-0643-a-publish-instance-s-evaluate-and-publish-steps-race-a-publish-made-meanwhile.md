---
id: S-0643
type: story
nature: improvement
title: A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make
status: backlog
owner: alex
created: 2026-11-09T21:22:33Z
updated: 2026-11-09T21:22:33Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-09T21:22:33Z
---
# S-0643 A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make

## Goal

This story remediates [I-0320](../../../design/issues/I-0320-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md), "A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0320 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0320 is closed with `flai issue close I-0320 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0320. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T16:58:32Z, 0.2 days before this story; under one cycle counts as one).

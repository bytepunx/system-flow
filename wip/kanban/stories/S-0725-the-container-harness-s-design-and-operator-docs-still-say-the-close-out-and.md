---
id: S-0725
type: story
nature: remediation
title: The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container
status: backlog
owner: alex
created: 2026-11-11T18:16:51Z
updated: 2026-11-11T18:16:51Z
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
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 10
          cache_read: 62150
          cache_write: 9071
          cost: 0.0202
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-11T18:16:51Z
---
# S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container

## Goal

This story remediates [I-0350](../../../design/issues/I-0350-the-container-harness-s-design-and-operator-docs-still-say-the-close-out-and.md), "The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0350 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0350 is closed with `flai issue close I-0350 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0350. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T16:04:40Z, 0.1 days before this story; under one cycle counts as one).

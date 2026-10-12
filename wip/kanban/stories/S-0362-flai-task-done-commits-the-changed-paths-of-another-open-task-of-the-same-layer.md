---
id: S-0362
type: story
nature: remediation
title: flai task done commits the changed paths of another open task of the same layer
status: backlog
owner: alex
created: 2026-10-12T14:12:29Z
updated: 2026-10-12T14:12:29Z
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
    at: 2026-10-12T14:12:29Z
---
# S-0362 flai task done commits the changed paths of another open task of the same layer

## Goal

This story remediates [I-0126](../../../design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md), "flai task done commits the changed paths of another open task of the same layer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0126 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0126 is closed with `flai issue close I-0126 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0126. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:36:24Z, 4.2 days before this story; under one cycle counts as one).

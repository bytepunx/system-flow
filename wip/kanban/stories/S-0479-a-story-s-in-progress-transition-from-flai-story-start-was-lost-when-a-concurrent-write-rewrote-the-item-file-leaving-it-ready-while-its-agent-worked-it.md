---
id: S-0479
type: story
nature: remediation
title: A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it
status: backlog
owner: alex
created: 2026-10-21T11:59:43Z
updated: 2026-10-21T12:04:48Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-21T11:59:43Z
finalized:
  by: alex
  at: 2026-10-21T12:04:48Z
---
# S-0479 A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it

## Goal

This story remediates [I-0175](../../../design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md), "A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0175 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0175 is closed with `flai issue close I-0175 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0175. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T07:37:23Z, 0.2 days before this story; under one cycle counts as one).

---
id: S-0451
type: story
nature: improvement
title: message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new
status: backlog
owner: alex
created: 2026-10-20T18:59:50Z
updated: 2026-10-20T19:01:44Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-20T18:59:50Z
finalized:
  by: alex
  at: 2026-10-20T19:01:44Z
---
# S-0451 message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new

## Goal

This story remediates [I-0169](../../../design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md), "message_share refuses a split between two stories in progress, so an agreed split cannot lift the overlap that refuses the other story's flai task new". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0169 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0169 is closed with `flai issue close I-0169 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0169. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T00:34:27Z, 0.8 days before this story; under one cycle counts as one).

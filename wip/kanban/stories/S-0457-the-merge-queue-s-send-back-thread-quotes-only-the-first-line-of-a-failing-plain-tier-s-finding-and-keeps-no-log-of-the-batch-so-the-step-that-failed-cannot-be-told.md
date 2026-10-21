---
id: S-0457
type: story
nature: improvement
title: The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told
status: backlog
owner: alex
created: 2026-10-21T02:13:46Z
updated: 2026-10-21T02:13:46Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-21T02:13:46Z
---
# S-0457 The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told

## Goal

This story remediates [I-0172](../../../design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md), "The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0172 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0172 is closed with `flai issue close I-0172 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0172. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T01:54:16Z, 0 days before this story; under one cycle counts as one).

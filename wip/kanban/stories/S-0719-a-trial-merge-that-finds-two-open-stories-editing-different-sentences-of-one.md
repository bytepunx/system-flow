---
id: S-0719
type: story
nature: improvement
title: A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance
status: backlog
owner: alex
created: 2026-11-11T02:29:04Z
updated: 2026-11-11T08:50:48Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-11T02:29:04Z
finalized:
  by: alex
  at: 2026-11-11T08:50:48Z
---
# S-0719 A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance

## Goal

This story remediates [I-0346](../../../design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md), "A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0346 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0346 is closed with `flai issue close I-0346 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0346. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T01:43:39Z, 0 days before this story; under one cycle counts as one).

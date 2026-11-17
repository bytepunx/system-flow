---
id: S-0838
type: story
nature: remediation
title: flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there
status: backlog
owner: alex
created: 2026-11-17T08:32:49Z
updated: 2026-11-17T08:32:49Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:32:49Z
---
# S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there

## Goal

This story remediates [I-0399](../../../design/issues/I-0399-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md), "flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0399 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0399 is closed with `flai issue close I-0399 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0399. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-17T08:25:17Z, 0 days before this story; under one cycle counts as one).

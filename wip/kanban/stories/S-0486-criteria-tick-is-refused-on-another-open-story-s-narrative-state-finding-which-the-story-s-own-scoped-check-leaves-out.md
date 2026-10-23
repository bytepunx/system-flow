---
id: S-0486
type: story
nature: remediation
title: criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out
status: backlog
owner: alex
created: 2026-10-23T17:43:07Z
updated: 2026-10-23T17:43:07Z
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
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-23T17:43:07Z
---
# S-0486 criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out

## Goal

This story remediates [I-0182](../../../design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md), "criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0182 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0182 is closed with `flai issue close I-0182 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0182. time_lost_per_cycle 4m: 2m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-21T12:14:14Z, 2.2 days before this story; under one cycle counts as one).

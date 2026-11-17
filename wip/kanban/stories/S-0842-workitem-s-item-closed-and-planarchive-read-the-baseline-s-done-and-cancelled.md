---
id: S-0842
type: story
nature: remediation
title: workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it
status: backlog
owner: alex
created: 2026-11-17T08:32:53Z
updated: 2026-11-17T08:32:53Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:53Z
---
# S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it

## Goal

This story remediates [I-0377](../../../design/issues/I-0377-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md), "workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0377 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0377 is closed with `flai issue close I-0377 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0377. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T01:36:52Z, 2.3 days before this story; under one cycle counts as one).

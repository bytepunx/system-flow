---
id: S-0644
type: story
nature: improvement
title: A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route
status: backlog
owner: alex
created: 2026-11-09T21:23:04Z
updated: 2026-11-09T21:23:04Z
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
    at: 2026-11-09T21:23:04Z
---
# S-0644 A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route

## Goal

This story remediates [I-0322](../../../design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md), "A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0322 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0322 is closed with `flai issue close I-0322 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0322. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T17:35:44Z, 0.2 days before this story; under one cycle counts as one).

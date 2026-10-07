---
id: S-0314
type: story
nature: remediation
title: The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off
status: backlog
owner: alex
created: 2026-10-07T18:59:44Z
updated: 2026-10-07T18:59:44Z
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
    at: 2026-10-07T18:59:44Z
---
# S-0314 The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off

## Goal

This story remediates [I-0115](../../../design/issues/I-0115-the-flaiover-guide-says-check-for-updates-always-works-but-the-dashboard-area-hides-it-while-the-dashboard-host-action-is-off.md), "The flaiover guide says Check for updates always works, but the Dashboard area hides it while the dashboard host action is off". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0115 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0115 is closed with `flai issue close I-0115 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0115. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T15:11:31Z, 0.2 days before this story; under one cycle counts as one).

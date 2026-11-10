---
id: S-0655
type: story
nature: remediation
title: "A flaiover design-page test, \"takes up the draft another tab writes\", times out at 5,000 ms under load and stops the flaiover tier"
status: backlog
owner: alex
created: 2026-11-10T05:54:35Z
updated: 2026-11-10T05:54:35Z
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
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-11-10T05:54:35Z
---
# S-0655 A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier

## Goal

This story remediates [I-0330](../../../design/issues/I-0330-a-flaiover-design-page-test-takes-up-the-draft-another-tab-writes-times-out-at.md), "A flaiover design-page test, "takes up the draft another tab writes", times out at 5,000 ms under load and stops the flaiover tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0330 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0330 is closed with `flai issue close I-0330 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0330. time_lost_per_cycle 16m: 16m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T05:37:30Z, 0 days before this story; under one cycle counts as one).

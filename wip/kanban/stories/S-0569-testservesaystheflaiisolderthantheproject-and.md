---
id: S-0569
type: story
nature: remediation
title: TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve
status: backlog
owner: alex
created: 2026-11-02T21:41:39Z
updated: 2026-11-02T21:41:39Z
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
    time_lost_per_cycle: 21m
    by: flai
    at: 2026-11-02T21:41:39Z
---
# S-0569 TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve

## Goal

This story remediates [I-0247](../../../design/issues/I-0247-testservesaystheflaiisolderthantheproject-and.md), "TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0247 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0247 is closed with `flai issue close I-0247 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0247. time_lost_per_cycle 21m: 7m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-02T01:00:19Z, 0.9 days before this story; under one cycle counts as one).

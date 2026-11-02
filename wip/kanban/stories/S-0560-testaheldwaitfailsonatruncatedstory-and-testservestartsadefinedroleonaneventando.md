---
id: S-0560
type: story
nature: remediation
title: TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80
status: backlog
owner: alex
created: 2026-11-02T03:55:02Z
updated: 2026-11-02T03:55:02Z
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
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-11-02T03:55:02Z
---
# S-0560 TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80

## Goal

This story remediates [I-0246](../../../design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md), "TestAHeldWaitFailsOnATruncatedStory and TestServeStartsADefinedRoleOnAnEventAndOnASchedule time out under flai test's Go run when the host's load is near 80". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0246 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0246 is closed with `flai issue close I-0246 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0246. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T00:55:15Z, 0.1 days before this story; under one cycle counts as one).

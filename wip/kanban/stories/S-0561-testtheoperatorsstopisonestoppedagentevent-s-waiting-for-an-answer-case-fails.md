---
id: S-0561
type: story
nature: remediation
title: "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped"
status: backlog
owner: alex
created: 2026-11-02T03:56:56Z
updated: 2026-11-02T03:56:56Z
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
    at: 2026-11-02T03:56:56Z
---
# S-0561 TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Goal

This story remediates [I-0248](../../../design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md), "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0248 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0248 is closed with `flai issue close I-0248 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0248. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T02:50:06Z, 0 days before this story; under one cycle counts as one).

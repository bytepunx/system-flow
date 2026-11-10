---
id: S-0652
type: story
nature: remediation
title: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
status: backlog
owner: alex
created: 2026-11-10T05:44:45Z
updated: 2026-11-10T05:44:45Z
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
    time_lost_per_cycle: 11m
    by: flai
    at: 2026-11-10T05:44:45Z
---
# S-0652 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Goal

This story remediates [I-0329](../../../design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md), "TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0329 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0329 is closed with `flai issue close I-0329 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0329. time_lost_per_cycle 11m: 11m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T05:34:34Z, 0 days before this story; under one cycle counts as one).

---
id: S-0389
type: story
nature: remediation
title: "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged"
status: backlog
owner: alex
created: 2026-10-14T23:09:14Z
updated: 2026-10-14T23:09:14Z
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
    at: 2026-10-14T23:09:14Z
---
# S-0389 TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Goal

This story remediates [I-0146](../../../design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md), "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0146 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0146 is closed with `flai issue close I-0146 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0146. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-13T19:01:05Z, 1.2 days before this story; under one cycle counts as one).

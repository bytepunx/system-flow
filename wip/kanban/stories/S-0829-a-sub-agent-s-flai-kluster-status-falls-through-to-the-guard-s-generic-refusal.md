---
id: S-0829
type: story
nature: remediation
title: A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run
status: backlog
owner: alex
created: 2026-11-17T08:17:21Z
updated: 2026-11-17T08:17:21Z
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
    at: 2026-11-17T08:17:21Z
---
# S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run

## Goal

This story remediates [I-0355](../../../design/issues/I-0355-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md), "A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0355 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0355 is closed with `flai issue close I-0355 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0355. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T07:33:55Z, 4 days before this story; under one cycle counts as one).

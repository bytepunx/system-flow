---
id: S-0678
type: story
nature: improvement
title: agent.status with 13 projects at once is not under 200 ms
status: backlog
owner: alex
created: 2026-11-10T08:34:55Z
updated: 2026-11-10T08:34:55Z
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
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-11-10T08:34:55Z
---
# S-0678 agent.status with 13 projects at once is not under 200 ms

## Goal

This story remediates [I-0313](../../../design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md), "agent.status with 13 projects at once is not under 200 ms". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0313 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0313 is closed with `flai issue close I-0313 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0313. time_lost_per_cycle 1s: 1s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T02:32:22Z, 2.3 days before this story; under one cycle counts as one).

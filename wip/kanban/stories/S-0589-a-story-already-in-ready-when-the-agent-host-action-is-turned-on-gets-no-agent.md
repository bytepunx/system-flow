---
id: S-0589
type: story
nature: remediation
title: A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent
status: backlog
owner: alex
created: 2026-11-07T09:44:04Z
updated: 2026-11-07T10:02:06Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T09:44:04Z
finalized:
  by: alex
  at: 2026-11-07T10:02:06Z
---
# S-0589 A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent

## Goal

This story remediates [I-0300](../../../design/issues/I-0300-a-story-already-in-ready-when-the-agent-host-action-is-turned-on-gets-no-agent.md), "A story already in ready when the agent host action is turned on gets no agent until someone presses Start agent". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0300 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0300 is closed with `flai issue close I-0300 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0300. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:31:47Z, 0.1 days before this story; under one cycle counts as one).

---
id: S-0568
type: story
nature: improvement
title: flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them
status: backlog
owner: alex
created: 2026-11-02T20:31:49Z
updated: 2026-11-07T08:01:38Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-02T20:31:49Z
finalized:
  by: alex
  at: 2026-11-07T08:01:38Z
---
# S-0568 flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Goal

This story remediates [I-0245](../../../design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md), "flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0245 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0245 is closed with `flai issue close I-0245 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0245. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T00:53:54Z, 0.8 days before this story; under one cycle counts as one).

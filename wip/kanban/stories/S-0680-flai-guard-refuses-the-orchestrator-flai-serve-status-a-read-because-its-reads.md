---
id: S-0680
type: story
nature: improvement
title: flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status
status: backlog
owner: alex
created: 2026-11-10T08:34:57Z
updated: 2026-11-10T08:34:57Z
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
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-11-10T08:34:57Z
---
# S-0680 flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Goal

This story remediates [I-0318](../../../design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md), "flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0318 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0318 is closed with `flai issue close I-0318 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0318. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T09:06:05Z, 1 day before this story; under one cycle counts as one).

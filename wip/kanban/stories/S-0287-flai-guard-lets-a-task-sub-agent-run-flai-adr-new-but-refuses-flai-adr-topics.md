---
id: S-0287
type: story
nature: remediation
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
status: backlog
owner: alex
created: 2026-10-06T09:56:50Z
updated: 2026-10-06T09:56:50Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-06T09:56:50Z
---
# S-0287 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Goal

This story remediates [I-0062](../../../design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md), "flai guard lets a task sub-agent run flai adr new but refuses flai adr topics". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0062 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0062 is closed with `flai issue close I-0062 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0062. time_lost_per_cycle 12m: 4m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-03T07:27:30Z, 3.1 days before this story; under one cycle counts as one).

---
id: S-0296
type: story
nature: remediation
title: A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts
status: backlog
owner: alex
created: 2026-10-06T11:44:51Z
updated: 2026-10-06T11:45:52Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h48m
    by: flai
    at: 2026-10-06T11:44:51Z
finalized:
  by: alex
  at: 2026-10-06T11:45:52Z
---
# S-0296 A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts

## Goal

This story remediates [I-0088](../../../design/issues/I-0088-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md), "A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts". The issue recommends this solution:

S-0221, already in ready, lets the orchestrator accept a story in review when the operator turns that permission on. S-0286 keeps the operator's acceptance for a story that changes a path Claude Code protects. Between them, the wait would remain only for those stories and for the ones the operator chooses to review. Until S-0221 is built, an agent in the operator's session can run the acceptance on the operator's word.

## Acceptance criteria
- [ ] The cause I-0088 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0088 is closed with `flai issue close I-0088 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0088. time_lost_per_cycle 2h48m: 1h24m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 0.1 days before this story; under one cycle counts as one).

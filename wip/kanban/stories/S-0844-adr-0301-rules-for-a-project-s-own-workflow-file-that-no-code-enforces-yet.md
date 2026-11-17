---
id: S-0844
type: story
nature: remediation
title: "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds"
status: backlog
owner: alex
created: 2026-11-17T08:32:55Z
updated: 2026-11-17T08:32:55Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:55Z
---
# S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds

## Goal

This story remediates [I-0380](../../../design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md), "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0380 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0380 is closed with `flai issue close I-0380 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0380. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T23:00:06Z, 1.4 days before this story; under one cycle counts as one).

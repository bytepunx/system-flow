---
id: S-0671
type: story
nature: improvement
title: scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again
status: backlog
owner: alex
created: 2026-11-10T08:34:48Z
updated: 2026-11-10T08:34:48Z
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
    at: 2026-11-10T08:34:48Z
---
# S-0671 scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again

## Goal

This story remediates [I-0280](../../../design/issues/I-0280-scripts-close-out-sh-takes-no-fresh-so-at-a-commit-the-verify-record-holds-a.md), "scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0280 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0280 is closed with `flai issue close I-0280 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0280. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:58:50Z, 3.1 days before this story; under one cycle counts as one).

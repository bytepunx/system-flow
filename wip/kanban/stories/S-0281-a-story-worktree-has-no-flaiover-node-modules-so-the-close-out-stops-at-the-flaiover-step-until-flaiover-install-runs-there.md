---
id: S-0281
type: story
nature: improvement
title: A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
status: backlog
owner: alex
created: 2026-10-05T07:09:06Z
updated: 2026-10-06T22:49:26Z
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
    at: 2026-10-05T07:09:06Z
finalized:
  by: alex
  at: 2026-10-06T22:49:26Z
---
# S-0281 A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Goal

This story remediates [I-0080](../../../design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md), "A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0080 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0080 is closed with `flai issue close I-0080 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0080. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T07:07:27Z, 0 days before this story; under one cycle counts as one).

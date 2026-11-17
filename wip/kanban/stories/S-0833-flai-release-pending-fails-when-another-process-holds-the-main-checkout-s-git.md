---
id: S-0833
type: story
nature: remediation
title: flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1
status: backlog
owner: alex
created: 2026-11-17T08:17:27Z
updated: 2026-11-17T08:17:27Z
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
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-11-17T08:17:27Z
---
# S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1

## Goal

This story remediates [I-0359](../../../design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md), "flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0359 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0359 is closed with `flai issue close I-0359 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0359. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-13T08:58:25Z, 4 days before this story; under one cycle counts as one).

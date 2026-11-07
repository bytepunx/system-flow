---
id: S-0634
type: story
nature: improvement
title: flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier
status: backlog
owner: alex
created: 2026-11-07T23:41:12Z
updated: 2026-11-07T23:41:12Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T23:41:12Z
---
# S-0634 flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier

## Goal

This story remediates [I-0308](../../../design/issues/I-0308-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md), "flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0308 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0308 is closed with `flai issue close I-0308 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0308. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T22:56:50Z, 0 days before this story; under one cycle counts as one).

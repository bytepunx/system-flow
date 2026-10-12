---
id: S-0371
type: story
nature: remediation
title: flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
status: backlog
owner: alex
created: 2026-10-12T16:52:47Z
updated: 2026-10-12T16:52:47Z
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
    at: 2026-10-12T16:52:47Z
---
# S-0371 flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

## Goal

This story remediates [I-0124](../../../design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md), "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0124 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0124 is closed with `flai issue close I-0124 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0124. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T09:00:31Z, 4.3 days before this story; under one cycle counts as one).

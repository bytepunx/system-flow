---
id: S-0480
type: story
nature: remediation
title: A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again
status: cancelled
owner: alex
created: 2026-10-21T11:59:52Z
updated: 2026-10-21T12:46:53Z
transitions:
  - to: cancelled
    at: 2026-10-21T12:46:53Z
    by: system-flow
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T11:59:52Z
finalized:
  by: alex
  at: 2026-10-21T12:03:36Z
---
# S-0480 A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again

## Goal

This story remediates [I-0176](../../../design/issues/I-0176-a-usage-write-to-a-story-overwrites-its-move-to-in-progress-made-meanwhile-so-a-started-story-reads-as-ready-again.md), "A usage write to a story overwrites its move to in-progress made meanwhile, so a started story reads as ready again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0176 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0176 is closed with `flai issue close I-0176 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0176. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T07:48:30Z, 0.2 days before this story; under one cycle counts as one).
- 2026-10-21T12:46:53Z: moved to cancelled: duplicate of S-0479, merged on the operator's yes on TH-0577: S-0479's Repo.Save merge fixes I-0176 and closes it

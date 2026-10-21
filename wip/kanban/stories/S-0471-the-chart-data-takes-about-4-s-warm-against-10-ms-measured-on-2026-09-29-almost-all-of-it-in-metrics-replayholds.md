---
id: S-0471
type: story
nature: remediation
title: The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds
status: backlog
owner: alex
created: 2026-10-21T11:58:35Z
updated: 2026-10-21T12:03:57Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4s
    by: flai
    at: 2026-10-21T11:58:35Z
finalized:
  by: alex
  at: 2026-10-21T12:03:57Z
---
# S-0471 The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds

## Goal

This story remediates [I-0155](../../../design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md), "The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0155 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0155 is closed with `flai issue close I-0155 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0155. time_lost_per_cycle 4s: 4s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-16T06:09:48Z, 5.2 days before this story; under one cycle counts as one).

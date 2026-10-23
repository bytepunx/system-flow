---
id: S-0489
type: story
nature: remediation
title: "metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does"
status: backlog
owner: alex
created: 2026-10-23T17:53:20Z
updated: 2026-10-23T17:53:20Z
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
    at: 2026-10-23T17:53:20Z
---
# S-0489 metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does

## Goal

This story remediates [I-0193](../../../design/issues/I-0193-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md), "metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0193 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0193 is closed with `flai issue close I-0193 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0193. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-22T21:09:37Z, 0.9 days before this story; under one cycle counts as one).

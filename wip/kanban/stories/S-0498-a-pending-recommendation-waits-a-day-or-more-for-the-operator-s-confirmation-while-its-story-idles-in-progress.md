---
id: S-0498
type: story
nature: improvement
title: A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress
status: backlog
owner: alex
created: 2026-10-26T04:01:18Z
updated: 2026-10-26T04:01:18Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: analyzer
      seconds: 114
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 98
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 117h
    by: flai
    at: 2026-10-26T04:01:18Z
---
# S-0498 A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Goal

This story remediates [I-0151](../../../design/issues/I-0151-a-pending-recommendation-waits-a-day-or-more-for-the-operator-s-confirmation-while-its-story-idles-in-progress.md), "A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0151 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0151 is closed with `flai issue close I-0151 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0151: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0151. time_lost_per_cycle 117h carried over from I-0151's Impact section. Its Impact time_lost_per_cycle was taken rather than the 19h20m derived from its cost and count.

Strategic usage carried over from I-0151: analyzer 1m54s, 1.2510 USD.

---
id: S-0469
type: story
nature: improvement
title: Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit
status: backlog
owner: alex
created: 2026-10-21T11:58:17Z
updated: 2026-10-22T20:53:14Z
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
      seconds: 113
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 97
          output: 6305
          cache_read: 370754
          cache_write: 41687
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23258
          cost: 0.0083
cost_of_delay:
  inputs:
    time_lost_per_cycle: 75h
    by: flai
    at: 2026-10-21T11:58:17Z
finalized:
  by: alex
  at: 2026-10-22T20:53:14Z
---
# S-0469 Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit

## Goal

This story remediates [I-0153](../../../design/issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md), "Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0153 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0153 is closed with `flai issue close I-0153 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0153: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0153. time_lost_per_cycle 75h carried over from I-0153's Impact section. Its Impact time_lost_per_cycle was taken rather than the 75h derived from its cost and count.

Strategic usage carried over from I-0153: analyzer 1m53s, 1.2510 USD.

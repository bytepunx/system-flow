---
id: S-0847
type: story
nature: improvement
title: flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator
status: backlog
owner: alex
created: 2026-11-17T08:32:58Z
updated: 2026-11-17T08:32:58Z
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
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-17T08:32:58Z
---
# S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator

## Goal

This story remediates [I-0389](../../../design/issues/I-0389-flai-serve-opens-a-new-board-starvation-thread-at-each-flip-of-the-board-48-in.md), "flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0389 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0389 is closed with `flai issue close I-0389 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0389: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0389. time_lost_per_cycle 30m carried over from I-0389's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1m derived from its cost and count.

Strategic usage carried over from I-0389: analyzer 1m39s, 1.2588 USD.

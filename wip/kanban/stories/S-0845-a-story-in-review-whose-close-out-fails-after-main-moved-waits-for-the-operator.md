---
id: S-0845
type: story
nature: improvement
title: A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves
status: backlog
owner: alex
created: 2026-11-17T08:32:56Z
updated: 2026-11-17T08:32:56Z
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
    time_lost_per_cycle: 30h
    by: flai
    at: 2026-11-17T08:32:56Z
---
# S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves

## Goal

This story remediates [I-0386](../../../design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md), "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0386 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0386 is closed with `flai issue close I-0386 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0386: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0386. time_lost_per_cycle 30h carried over from I-0386's Impact section. Its Impact time_lost_per_cycle was taken rather than the 30h derived from its cost and count.

Strategic usage carried over from I-0386: analyzer 1m39s, 1.2588 USD.

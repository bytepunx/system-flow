---
id: S-0688
type: story
nature: improvement
title: A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host
status: backlog
owner: alex
created: 2026-11-10T08:37:15Z
updated: 2026-11-10T08:37:15Z
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
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:37:15Z
---
# S-0688 A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host

## Goal

This story remediates [I-0257](../../../design/issues/I-0257-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md), "A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0257 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0257 is closed with `flai issue close I-0257 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0257: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0257. time_lost_per_cycle 9m: 5m per occurrence × 2 occurrences ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:40Z, 7.6 days before this story).

Strategic usage carried over from I-0257: analyzer 42s, 0.8156 USD.

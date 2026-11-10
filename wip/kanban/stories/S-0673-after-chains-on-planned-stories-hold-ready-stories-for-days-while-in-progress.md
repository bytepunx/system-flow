---
id: S-0673
type: story
nature: improvement
title: after chains on planned stories hold ready stories for days while in-progress stays under two of five
status: backlog
owner: alex
created: 2026-11-10T08:34:51Z
updated: 2026-11-10T08:34:51Z
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
    time_lost_per_cycle: 56h
    by: flai
    at: 2026-11-10T08:34:51Z
---
# S-0673 after chains on planned stories hold ready stories for days while in-progress stays under two of five

## Goal

This story remediates [I-0285](../../../design/issues/I-0285-after-chains-on-planned-stories-hold-ready-stories-for-days-while-in-progress.md), "after chains on planned stories hold ready stories for days while in-progress stays under two of five", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0285 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0285 is closed with `flai issue close I-0285 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0285: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0285. time_lost_per_cycle 56h carried over from I-0285's Impact section.

Strategic usage carried over from I-0285: analyzer 42s, 0.8156 USD.

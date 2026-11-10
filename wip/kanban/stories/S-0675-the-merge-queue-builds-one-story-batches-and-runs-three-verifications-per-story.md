---
id: S-0675
type: story
nature: improvement
title: The merge queue builds one-story batches and runs three verifications per story landed
status: backlog
owner: alex
created: 2026-11-10T08:34:53Z
updated: 2026-11-10T08:34:53Z
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
          input: 32
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7h
    by: flai
    at: 2026-11-10T08:34:53Z
---
# S-0675 The merge queue builds one-story batches and runs three verifications per story landed

## Goal

This story remediates [I-0287](../../../design/issues/I-0287-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md), "The merge queue builds one-story batches and runs three verifications per story landed", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0287 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0287 is closed with `flai issue close I-0287 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0287: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0287. time_lost_per_cycle 7h carried over from I-0287's Impact section.

Strategic usage carried over from I-0287: analyzer 42s, 0.8156 USD.

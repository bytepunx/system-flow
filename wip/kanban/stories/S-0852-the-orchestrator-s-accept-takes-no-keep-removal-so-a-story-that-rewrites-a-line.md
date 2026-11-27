---
id: S-0852
type: story
nature: improvement
title: The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right
status: backlog
owner: alex
created: 2026-11-27T22:55:19Z
updated: 2026-11-27T22:55:19Z
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
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29513
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 148h
    by: flai
    at: 2026-11-27T22:55:19Z
---
# S-0852 The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right

## Goal

This story remediates [I-0415](../../../design/issues/I-0415-the-orchestrator-s-accept-takes-no-keep-removal-so-a-story-that-rewrites-a-line.md), "The orchestrator's accept takes no --keep-removal, so a story that rewrites a line main added waits in review for the operator even after the story's agent and the orchestrator agree the removal is right", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0415 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0415 is closed with `flai issue close I-0415 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0415: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0415. time_lost_per_cycle 148h carried over from I-0415's Impact section. Its Impact time_lost_per_cycle was taken rather than the 148h20m derived from its cost and count.

Strategic usage carried over from I-0415: analyzer 1m25s, 1.4522 USD.

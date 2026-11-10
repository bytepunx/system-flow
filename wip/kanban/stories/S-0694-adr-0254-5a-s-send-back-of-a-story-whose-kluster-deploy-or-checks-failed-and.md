---
id: S-0694
type: story
nature: improvement
title: "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built"
status: backlog
owner: alex
created: 2026-11-10T08:54:25Z
updated: 2026-11-10T08:54:25Z
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
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:54:25Z
---
# S-0694 ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built

## Goal

This story remediates [I-0268](../../../design/issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md), "ADR-0254 §5a's send-back of a story whose kluster deploy or checks failed, and its kluster acceptance blocker, are not built", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0268 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0268 is closed with `flai issue close I-0268 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0268: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0268. time_lost_per_cycle 20m: 10m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T03:22:56Z, 3.2 days before this story; under one cycle counts as one).

Strategic usage carried over from I-0268: analyzer 42s, 0.8156 USD.

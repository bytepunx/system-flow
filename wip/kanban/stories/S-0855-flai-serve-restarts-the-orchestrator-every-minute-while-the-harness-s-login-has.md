---
id: S-0855
type: story
nature: remediation
title: flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator
status: backlog
owner: alex
created: 2026-11-27T22:55:27Z
updated: 2026-11-27T22:55:27Z
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
          cache_write: 29514
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
    time_lost_per_cycle: 14h
    by: flai
    at: 2026-11-27T22:55:27Z
---
# S-0855 flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator

## Goal

This story remediates [I-0412](../../../design/issues/I-0412-flai-serve-restarts-the-orchestrator-every-minute-while-the-harness-s-login-has.md), "flai serve restarts the orchestrator every minute while the harness's login has expired, 666 runs in 14 hours, with no hold and no thread to the operator", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0412 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0412 is closed with `flai issue close I-0412 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0412: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0412. time_lost_per_cycle 14h carried over from I-0412's Impact section. Its Impact time_lost_per_cycle was taken rather than the 14h derived from its cost and count.

Strategic usage carried over from I-0412: analyzer 1m25s, 1.4522 USD.

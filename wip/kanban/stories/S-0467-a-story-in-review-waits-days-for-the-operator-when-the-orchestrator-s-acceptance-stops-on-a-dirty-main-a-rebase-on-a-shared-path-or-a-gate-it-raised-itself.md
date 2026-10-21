---
id: S-0467
type: story
nature: improvement
title: A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself
status: backlog
owner: alex
created: 2026-10-21T11:57:59Z
updated: 2026-10-21T11:57:59Z
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
          cache_read: 204810
          cache_write: 23259
          cost: 0.0083
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 96h
    by: flai
    at: 2026-10-21T11:57:59Z
---
# S-0467 A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself

## Goal

This story remediates [I-0150](../../../design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md), "A story in review waits days for the operator when the orchestrator's acceptance stops on a dirty main, a rebase on a shared path, or a gate it raised itself", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0150 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0150 is closed with `flai issue close I-0150 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0150: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0150. time_lost_per_cycle 96h carried over from I-0150's Impact section. Its Impact time_lost_per_cycle was taken rather than the 96h derived from its cost and count.

Strategic usage carried over from I-0150: analyzer 1m54s, 1.2510 USD.

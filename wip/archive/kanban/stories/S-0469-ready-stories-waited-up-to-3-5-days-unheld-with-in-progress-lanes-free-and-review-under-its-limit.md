---
id: S-0469
type: story
nature: improvement
title: Ready stories waited up to 3.5 days unheld with in-progress lanes free and review under its limit
status: done
owner: alex
created: 2026-10-21T11:58:17Z
updated: 2026-10-23T06:21:15Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:21Z
    by: alex
  - to: in-progress
    at: 2026-10-23T06:00:06Z
    by: agent-S-0469
  - to: review
    at: 2026-10-23T06:08:41Z
    by: agent-S-0469
  - to: done
    at: 2026-10-23T06:21:15Z
    by: orchestrator
tags: []
touches: [design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md, design/adrs/README.md, design/system/metrics.md, flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, design/issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 529
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      test_runs: 1
      hand_edits: 1
      work: 38
  models:
    - model: claude-haiku-5-5
      input: 62
      output: 46884
      cache_read: 3552734
      cache_write: 188468
      cost: 0.3233
    - model: claude-opus-5-5
      input: 84
      output: 19222
      cache_read: 5391614
      cache_write: 147288
      cost: 2.6414
  strategic:
    - kind: orchestrator
      seconds: 1785
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 89
          output: 1443
          cache_read: 28919466
          cache_write: 27605
          cost: 6.3918
        - model: claude-sonnet-5-5
          input: 19
          output: 110
          cache_read: 287700
          cache_write: 73514
          cost: 0.25
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
- [x] The cause I-0153 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0153 is closed with `flai issue close I-0153 --reason` saying what fixed it

## Tasks
- T-2064 An ADR settles that a story's queue time is the seconds it spent in ready before it started, so a stretch back in backlog is not counted as waiting
- T-2065 flai stats derives queue time from the ready intervals before started, with a test of a story sent back to backlog and made ready again

## Notes

The analyzer found I-0153: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0153. time_lost_per_cycle 75h carried over from I-0153's Impact section. Its Impact time_lost_per_cycle was taken rather than the 75h derived from its cost and count.

Strategic usage carried over from I-0153: analyzer 1m53s, 1.2510 USD.

### Accepted by the orchestrator

- Verified: 90cae2651d0d5e1b9e454e5b117038b6aee57366
- Batch: 20261023-060954
- Enqueued verified: 90cae2651d0d5e1b9e454e5b117038b6aee57366
- At: 2026-10-23T06:21:15Z

Verdict: pass at 90cae2651d0d5e1b9e454e5b117038b6aee57366
- 1: flai/internal/metrics/metrics.go, flai/internal/metrics/metrics_test.go, design/system/metrics.md, design/adrs/0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md, design/adrs/README.md
- 2: design/issues/I-0153-ready-stories-waited-up-to-3-5-days-unheld-with-in-progress-lanes-free-and-review-under-its-limit.md, design/issues/summary.md

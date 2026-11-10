---
id: S-0692
type: story
nature: improvement
title: No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days
status: backlog
owner: alex
created: 2026-11-10T08:49:03Z
updated: 2026-11-10T08:49:03Z
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
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
---
# S-0692 No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days

## Goal

This story remediates [I-0292](../../../design/issues/I-0292-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md), "No analysis.schedule or analysis.formula is set, so the analyzer ran once in thirty days", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0292 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0292 is closed with `flai issue close I-0292 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0292: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0292: analyzer 42s, 0.8156 USD.

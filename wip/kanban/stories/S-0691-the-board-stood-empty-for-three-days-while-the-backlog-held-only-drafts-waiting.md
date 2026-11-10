---
id: S-0691
type: story
nature: improvement
title: The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again
status: backlog
owner: alex
created: 2026-11-10T08:49:02Z
updated: 2026-11-10T08:49:02Z
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
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
---
# S-0691 The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again

## Goal

This story remediates [I-0288](../../../design/issues/I-0288-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md), "The board stood empty for three days while the backlog held only drafts waiting on operator decisions nobody asked for again", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0288 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0288 is closed with `flai issue close I-0288 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0288: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0288: analyzer 42s, 0.8156 USD.

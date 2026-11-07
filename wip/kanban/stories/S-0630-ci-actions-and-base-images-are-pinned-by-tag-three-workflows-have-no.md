---
id: S-0630
type: story
nature: improvement
title: CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles
status: backlog
owner: alex
created: 2026-11-07T23:40:47Z
updated: 2026-11-07T23:40:47Z
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
      seconds: 41
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
          cache_read: 759049
          cache_write: 13854
          cost: 0.278
draft: true
---
# S-0630 CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles

## Goal

This story remediates [I-0297](../../../design/issues/I-0297-ci-actions-and-base-images-are-pinned-by-tag-three-workflows-have-no.md), "CI actions and base images are pinned by tag, three workflows have no permissions block, and flaiover tracks two lockfiles", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0297 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0297 is closed with `flai issue close I-0297 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0297: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0297: analyzer 41s, 0.8156 USD.

---
id: S-0625
type: story
nature: improvement
title: "Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand"
status: backlog
owner: alex
created: 2026-11-07T23:40:15Z
updated: 2026-11-09T21:56:42Z
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
          input: 15
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
finalized:
  by: alex
  at: 2026-11-09T21:56:42Z
---
# S-0625 Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand

## Goal

This story remediates [I-0291](../../../design/issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md), "Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0291 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0291 is closed with `flai issue close I-0291 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0291: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0291: analyzer 42s, 0.8156 USD.

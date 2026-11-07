---
id: S-0624
type: story
nature: remediation
title: One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats
status: backlog
owner: alex
created: 2026-11-07T23:40:09Z
updated: 2026-11-07T23:40:09Z
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
draft: true
---
# S-0624 One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats

## Goal

This story remediates [I-0290](../../../design/issues/I-0290-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md), "One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0290 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0290 is closed with `flai issue close I-0290 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0290: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0290: analyzer 42s, 0.8156 USD.

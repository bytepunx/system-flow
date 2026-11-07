---
id: S-0629
type: story
nature: remediation
title: flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again
status: backlog
owner: alex
created: 2026-11-07T23:40:41Z
updated: 2026-11-07T23:40:41Z
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
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
draft: true
---
# S-0629 flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again

## Goal

This story remediates [I-0296](../../../design/issues/I-0296-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md), "flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0296 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0296 is closed with `flai issue close I-0296 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0296: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0296: analyzer 41s, 0.8156 USD.

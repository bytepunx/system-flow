---
id: S-0623
type: story
nature: improvement
title: flai guard refuses the orchestrator git tag --list and git worktree list, which only read
status: backlog
owner: alex
created: 2026-11-07T23:40:02Z
updated: 2026-11-07T23:40:02Z
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
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-07T23:40:02Z
---
# S-0623 flai guard refuses the orchestrator git tag --list and git worktree list, which only read

## Goal

This story remediates [I-0289](../../../design/issues/I-0289-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md), "flai guard refuses the orchestrator git tag --list and git worktree list, which only read", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0289 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0289 is closed with `flai issue close I-0289 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0289: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0289. time_lost_per_cycle 30m carried over from I-0289's Impact section.

Strategic usage carried over from I-0289: analyzer 42s, 0.8156 USD.

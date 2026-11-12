---
id: S-0726
type: story
nature: improvement
title: A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself
status: backlog
owner: alex
created: 2026-11-11T21:55:30Z
updated: 2026-11-11T21:55:30Z
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
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 10
          cache_read: 62150
          cache_write: 9071
          cost: 0.0202
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-11T21:55:30Z
---
# S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself

## Goal

This story remediates [I-0351](../../../design/issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md), "A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0351 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0351 is closed with `flai issue close I-0351 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0351. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T18:18:50Z, 0.2 days before this story; under one cycle counts as one).

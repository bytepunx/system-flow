---
id: S-0674
type: story
nature: improvement
title: Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations
status: backlog
owner: alex
created: 2026-11-10T08:34:52Z
updated: 2026-11-10T08:34:52Z
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
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: flai
    at: 2026-11-10T08:34:52Z
---
# S-0674 Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations

## Goal

This story remediates [I-0286](../../../design/issues/I-0286-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md), "Declared touches match what a story changes in 16 percent of completed stories, so claims mislead the pull hold, the trial merge, and the conversations", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0286 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0286 is closed with `flai issue close I-0286 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0286: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0286. time_lost_per_cycle 4h carried over from I-0286's Impact section.

Strategic usage carried over from I-0286: analyzer 42s, 0.8156 USD.

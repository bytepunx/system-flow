---
id: S-0470
type: story
nature: improvement
title: task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result
status: ready
owner: alex
created: 2026-10-21T11:58:27Z
updated: 2026-10-22T20:55:24Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:24Z
    by: alex
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
      seconds: 44
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 308
          cache_read: 3107007
          cache_write: 2860
          cost: 0.6867
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-21T11:58:27Z
finalized:
  by: alex
  at: 2026-10-22T20:53:05Z
---
# S-0470 task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result

## Goal

This story remediates [I-0154](../../../design/issues/I-0154-task-done-s-answer-repeats-every-open-thread-not-on-the-agent-s-story-so-it-overflows-a-tool-result.md), "task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0154 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0154 is closed with `flai issue close I-0154 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0154. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-16T06:05:59Z, 5.2 days before this story; under one cycle counts as one).

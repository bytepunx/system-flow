---
id: S-0305
type: story
nature: improvement
title: flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories
status: backlog
owner: alex
created: 2026-10-07T01:07:10Z
updated: 2026-10-07T02:19:51Z
transitions: []
tags: []
topics: [planning]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md, design/issues/summary.md]
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
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 1
          cache_read: 133364
          cache_write: 347
          cost: 0.0349
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T01:07:10Z
finalized:
  by: alex
  at: 2026-10-07T02:17:58Z
---
# S-0305 flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories

## Goal

This story remediates [I-0098](../../../design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md), "flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0098 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0098 is closed with `flai issue close I-0098 --reason` saying what fixed it

## Tasks
- T-1151 The replanner neither queues nor starts the planner for a story whose own agent is at work
- T-1152 The design and the user guide say an edit by a story's own agent, or a story whose agent is at work, starts no planner on its own
- T-1154 Close I-0098 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0098. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:28:35Z, 0.1 days before this story; under one cycle counts as one).

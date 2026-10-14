---
id: S-0305
type: story
nature: improvement
title: flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories
status: done
owner: alex
created: 2026-10-07T01:07:10Z
updated: 2026-10-14T21:15:25Z
transitions:
  - to: ready
    at: 2026-10-14T20:54:55Z
    by: alex
  - to: in-progress
    at: 2026-10-14T20:56:26Z
    by: agent-S-0305
  - to: review
    at: 2026-10-14T21:14:28Z
    by: system-flow
  - to: done
    at: 2026-10-14T21:15:25Z
    by: orchestrator
tags: []
topics: [planning]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md, design/issues/summary.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1092
  priced_by: harness
  turns:
    - day: 2026-10-14
      hand_edits: 1
      work: 33
  models:
    - model: claude-opus-5-5
      input: 152
      output: 35630
      cache_read: 6976628
      cache_write: 278383
      cost: 3.9432
  strategic:
    - kind: orchestrator
      seconds: 162
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 290
          cache_read: 2223604
          cache_write: 8153
          cost: 0.5609
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
- [x] The cause I-0098 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0098 is closed with `flai issue close I-0098 --reason` saying what fixed it

## Tasks
- T-1151 The replanner neither queues nor starts the planner for a story whose own agent is at work
- T-1152 The design and the user guide say an edit by a story's own agent, or a story whose agent is at work, starts no planner on its own
- T-1154 Close I-0098 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0098. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:28:35Z, 0.1 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: 9aa680e2dd2bd467fa49c7fd7275e0e8b2a99a3a
- At: 2026-10-14T21:15:25Z

Verdict: meets all criteria; flai verify passed every step at 9aa680e2; every changed path is within the story's touches.

- 1: flai/internal/serve/replan.go (editTrigger gives no trigger for the story's own agent; drain drops a queued planner for a story in progress, in review, or with a running agent), flai/internal/serve/replan_test.go (TestAnEditOfAPlannedStoryIsATrigger cases, TestTheQueueDropsAStoryItsAgentWorks), design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md
- 2: design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md (closed, fixed by S-0305), design/issues/summary.md

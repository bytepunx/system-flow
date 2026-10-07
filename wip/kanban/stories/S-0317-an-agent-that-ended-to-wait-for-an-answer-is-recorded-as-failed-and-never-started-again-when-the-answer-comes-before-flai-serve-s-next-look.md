---
id: S-0317
type: story
nature: remediation
title: An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look
status: backlog
owner: alex
created: 2026-10-07T18:59:47Z
updated: 2026-10-07T18:59:47Z
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
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 64217
          cache_write: 5432
          cost: 0.0172
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: flai
    at: 2026-10-07T18:59:47Z
---
# S-0317 An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look

## Goal

This story remediates [I-0095](../../../design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md), "An agent that ended to wait for an answer is recorded as failed, and never started again, when the answer comes before flai serve's next look". The issue recommends this solution:

Directions to weigh: when a run ends, read whether its story has a thread whose last answer came after the agent's last entry on it, and treat that as an answer to act on rather than as no question; or record at the moment an agent asks, through `thread_open` or `thread_reply` on its own story, that it is waiting, so that the end is `asked` whatever the thread says by the time of the look. Either way, a story in progress whose agent has ended, with no thread awaiting the operator and no block, could be restarted once by flai serve, which is the second remediation of I-0084.

## Acceptance criteria
- [ ] The cause I-0095 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0095 is closed with `flai issue close I-0095 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0095. time_lost_per_cycle 1m: 1m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T21:01:06Z, 0.9 days before this story; under one cycle counts as one).

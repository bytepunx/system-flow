---
id: S-0470
type: story
nature: improvement
title: task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result
status: done
owner: alex
created: 2026-10-21T11:58:27Z
updated: 2026-10-23T06:43:56Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:24Z
    by: alex
  - to: in-progress
    at: 2026-10-23T06:08:54Z
    by: agent-S-0470
  - to: review
    at: 2026-10-23T06:16:30Z
    by: agent-S-0470
  - to: done
    at: 2026-10-23T06:43:56Z
    by: orchestrator
tags: []
touches: [flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go, design/issues/I-0154-task-done-s-answer-repeats-every-open-thread-not-on-the-agent-s-story-so-it-overflows-a-tool-result.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 477
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 2
      hand_edits: 3
      work: 22
  models:
    - model: claude-opus-5-5
      input: 56
      output: 9438
      cache_read: 2905208
      cache_write: 132282
      cost: 1.8283
  strategic:
    - kind: orchestrator
      seconds: 1689
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 69
          output: 1152
          cache_read: 22438215
          cache_write: 21525
          cost: 4.9592
        - model: claude-sonnet-5-5
          input: 6
          output: 45
          cache_read: 43275
          cache_write: 36671
          cost: 0.0553
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
- [x] The cause I-0154 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0154 is closed with `flai issue close I-0154 --reason` saying what fixed it

## Tasks
- T-2066 The tests that reproduce I-0154 name it, and the story records that S-0450 and S-0453 removed its cause

## Notes

**Proposal and finding.** The cause is gone before this story started; no code change is needed. S-0450 ([ADR-0179](../../../design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md), flai v1.67.1) made `task_done` and `story_start` list only the story's threads, count the rest in `threads_omitted`, and fit each answer to 40,000 bytes. S-0453 ([ADR-0181](../../../design/adrs/0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md), flai v1.67.6) made a thread on another item that the agent never wrote to read `awaiting: operator`, so it is neither listed nor counted. The host's flai is 1.67.13. S-0470's own `story_start` answered with no thread listed and no `threads_omitted`. I-0154's instance, the operator's threads on other stories, epics, and a non-item path, none on the story, is reproduced by `TestTaskDoneCarriesNoneOfTheOperatorsThreads` (`flai/internal/mcpserver/task_test.go`) and `TestRunCarriesNoneOfTheOperatorsThreads` (`flai/internal/taskdone/taskdone_test.go`); T-2066 names I-0154 in them. I-0167 recorded the same cause on S-0398 to S-0404; `TestTaskDoneFitsTheAnswerWithManyThreadsElsewhere` holds the 40,000-byte fit.

Cost of delay inputs set by flai from I-0154. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-16T06:05:59Z, 5.2 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: 526a5aef1bda2fd930aa244f04f57e0dddebd04d
- Batch: 20261023-063332
- Enqueued verified: 513199b9bed4a919d12a6f7bfb97e9159b81fd75
- At: 2026-10-23T06:43:56Z

Verdict: pass at 513199b9bed4a919d12a6f7bfb97e9159b81fd75
- 1: flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go
- 2: design/issues/I-0154-task-done-s-answer-repeats-every-open-thread-not-on-the-agent-s-story-so-it-overflows-a-tool-result.md, design/issues/summary.md

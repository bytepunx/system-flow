---
id: S-0492
type: story
nature: improvement
title: task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result
status: done
owner: alex
created: 2026-10-23T23:38:52Z
updated: 2026-10-24T00:28:42Z
transitions:
  - to: ready
    at: 2026-10-23T23:40:08Z
    by: alex
  - to: in-progress
    at: 2026-10-23T23:47:00Z
    by: agent-S-0492
  - to: review
    at: 2026-10-24T00:06:00Z
    by: agent-S-0492
  - to: done
    at: 2026-10-24T00:28:42Z
    by: orchestrator
tags: []
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs, design/issues, design/adrs/0198-task-done-s-answer-lists-at-most-ten-of-the-check-s-notes-outside-the-story-and.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1253
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 3
      hand_edits: 2
      work: 47
    - day: 2026-10-24
      ceremony: 1
      work: 5
  models:
    - model: claude-opus-5-5
      input: 118
      output: 24910
      cache_read: 7573123
      cache_write: 170306
      cost: 3.3757
  strategic:
    - kind: orchestrator
      seconds: 626
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 43
          output: 806
          cache_read: 16248722
          cache_write: 16370
          cost: 3.5773
        - model: claude-sonnet-5-5
          input: 7
          output: 49
          cache_read: 89951
          cache_write: 18982
          cost: 0.0751
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-23T23:38:52Z
finalized:
  by: alex
  at: 2026-10-23T23:39:41Z
---
# S-0492 task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result

## Goal

This story remediates [I-0202](../../../design/issues/I-0202-task-done-s-answer-carries-every-check-note-outside-the-story-so-it-exceeds-the-40-000-bytes-it-promises-and-overflows-the-tool-result.md), "task_done's answer carries every check note outside the story, so it exceeds the 40,000 bytes it promises and overflows the tool result". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0202 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0202 is closed with `flai issue close I-0202 --reason` saying what fixed it

## Tasks
- T-2102 An ADR settles that task_done lists at most ten of the check's notes outside the story and counts the rest in notes_omitted
- T-2103 task_done's check lists at most ten notes outside the story and counts the rest in notes_omitted, with a test, the help, and the docs

## Notes

Cost of delay inputs set by flai from I-0202. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-23T21:50:41Z, 0.1 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: 3e3369483eafa7a773d28512d880822a2185f39c
- Batch: 20261024-001509
- Enqueued verified: 3e3369483eafa7a773d28512d880822a2185f39c
- At: 2026-10-24T00:28:42Z

Verdict: pass. task_done now lists at most ten notes outside the story and counts the rest in notes_omitted, a test reproduces I-0202, the docs and reference describe the field, and I-0202 is closed with a reason.
- 1: flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, design/adrs/0198-task-done-s-answer-lists-at-most-ten-of-the-check-s-notes-outside-the-story-and.md, docs/users/flai.md, design/system/flai-cli.md, docs/users/flai-reference.md
- 2: design/issues/I-0202-task-done-s-answer-carries-every-check-note-outside-the-story-so-it-exceeds-the-40-000-bytes-it-promises-and-overflows-the-tool-result.md, design/issues/summary.md

---
id: T-3252
type: task
nature: feature
title: One flaiover module gives a state's kind, colour, and moves from the columns and transitions board.get answers
status: done
parent: S-0758
owner: alex
created: 2026-11-12T16:09:42Z
updated: 2026-11-17T08:00:43Z
transitions:
  - to: ready
    at: 2026-11-17T07:56:08Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T07:56:08Z
    by: system-flow
  - to: done
    at: 2026-11-17T08:00:43Z
    by: agent-S-0758
stream: S-0758
tags: [dashboard]
touches: [flaiover/src/lib/workflow.ts, flaiover/src/lib/workflow.test.ts, flaiover/src/lib/lanes.fixture.ts]
usage:
  source: log
  seconds: 275
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 34
      output: 192
      cache_read: 1477494
      cache_write: 61362
      cost: 0.5861
---
# T-3252 One flaiover module gives a state's kind, colour, and moves from the columns and transitions board.get answers

## Work

Add `flaiover/src/lib/workflow.ts`, read from the columns and transitions `board.get` answers: a state's kind, its colour by kind, whether it is open, and the moves the operator may make from it with whether each needs a reason. Test it on the baseline's answer and on one with a state added between `in-progress` and `review`.

## Done when

`workflow.test.ts` passes on both answers.

## Notes

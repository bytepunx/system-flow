---
id: S-0799
type: story
nature: remediation
title: After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head
status: backlog
owner: alex
created: 2026-11-16T00:41:49Z
updated: 2026-11-16T00:41:49Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-16T00:41:49Z
---
# S-0799 After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head

## Goal

This story remediates [I-0384](../../../design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md), "After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head". The issue recommends this solution:

While a close-out job runs for a story at a head other than its last report's, `flai queue`, the acceptance preview, and the review page say so, naming the job and its step, in place of or beside the last finished report's failure: the job's state is in `.flai-cache/verify/jobs/<S-nnnn>.job.json`. A `verify_failed` whose report's commit is not the branch's head says the report is from before the head moved.

## Acceptance criteria
- [ ] The cause I-0384 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0384 is closed with `flai issue close I-0384 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0384. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-16T00:41:30Z, 0 days before this story; under one cycle counts as one).

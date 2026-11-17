---
id: S-0822
type: story
nature: remediation
title: A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it
status: backlog
owner: alex
created: 2026-11-17T07:13:18Z
updated: 2026-11-17T07:13:18Z
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
      seconds: 67
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 81
          cache_read: 783588
          cache_write: 96802
          cost: 0.2725
        - model: claude-sonnet-5-5
          input: 16
          output: 841
          cache_read: 167575
          cache_write: 67904
          cost: 0.1893
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-17T07:13:18Z
---
# S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it

## Goal

This story remediates [I-0396](../../../design/issues/I-0396-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md), "A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it". The issue recommends this solution:

When a rebase and resolve job passes, `settle` in `flai/internal/serve/resolvejob.go` reads the story again, as `sendBack` does. For a story no longer in review it records an outcome of its own, not `stayed in review`, and a reason that names the state the story is in, who moved it there and when, and the command that moves it back to review, since its close-out passed at its head. The review page shows that outcome as it shows the others. Decide in the story, with the operator, whether a move out of review while the job runs should also stop the job; record the decision in an ADR refining ADR-0313 if it changes what ADR-0313 says. A test in the resolve job's tests moves the story out of review while the job runs and reads the record.

## Acceptance criteria
- [ ] The cause I-0396 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0396 is closed with `flai issue close I-0396 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0396. time_lost_per_cycle 30m carried over from I-0396's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h30m derived from its cost and count.

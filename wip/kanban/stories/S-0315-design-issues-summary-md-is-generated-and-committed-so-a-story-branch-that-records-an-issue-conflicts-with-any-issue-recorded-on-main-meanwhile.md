---
id: S-0315
type: story
nature: improvement
title: design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile
status: backlog
owner: alex
created: 2026-10-07T18:59:45Z
updated: 2026-10-07T18:59:45Z
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
          cache_write: 5433
          cost: 0.0172
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-07T18:59:45Z
---
# S-0315 design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile

## Goal

This story remediates [I-0089](../../../design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md), "design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile". The issue recommends this solution:

Directions to weigh: `flai stream sync` and `flai accept` regenerate `summary.md` themselves when it is the only conflict, since it is derived from the issue files; or the file is not committed and is generated where it is read.

S-0278 built the first direction ([ADR-0098](../adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)) for I-0074, which has the same cause. The operator chooses at its acceptance whether that closes this issue too.

## Acceptance criteria
- [ ] The cause I-0089 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0089 is closed with `flai issue close I-0089 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0089. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 1.4 days before this story; under one cycle counts as one).

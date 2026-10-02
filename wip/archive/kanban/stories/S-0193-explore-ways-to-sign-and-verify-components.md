---
id: S-0193
type: story
nature: research
title: Explore ways to sign and verify components
status: done
parent: E-0015
owner: alex
created: 2026-10-01T11:18:38Z
updated: 2026-10-02T12:45:10Z
transitions:
  - to: ready
    at: 2026-10-01T11:18:38Z
    by: alex
  - to: in-progress
    at: 2026-10-02T12:10:05Z
    by: claude-fable-5-1
  - to: review
    at: 2026-10-02T12:43:26Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:45:10Z
    by: alex
tags: []
topics: [release, security]
touches: [flaiover/src, flai/cmd, design/system/release-signing.md, design/system/README.md, design/adrs, design/issues/I-0055-a-story-made-from-the-dashboard-reached-main-with-a-markdown-lint-error-flai-task-new-would-have-refused.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
usage:
  source: log
  seconds: 2081
  models:
    - model: claude-fable-5-1
      input: 2082
      output: 80320
      cache_read: 8413158
      cache_write: 202190
      cost: 10.1839
    - model: claude-haiku-4-5-20251001
      input: 266
      output: 9362
      cache_read: 1564698
      cache_write: 78430
      cost: 0.3016
    - model: claude-sonnet-5-5
      input: 38
      output: 11499
      cache_read: 451332
      cache_write: 87642
      cost: 0.4244
---
# S-0193 Explore ways to sign and verify components

## Goal

Research options for delivering on the parent epic's goals. Communicate options via questions/threads with the operator and document resulting approach as child stories for epic 15.

## Acceptance criteria
- [x] Reported options to the operator
- [x] Planned stories for epic 15 based on operator feedback and findings

## Tasks
- T-0691 Record how flai and flaiover are released, upgraded, and connected today
- T-0692 Set out the options for signing releases and verifying them at upgrade and image run
- T-0693 Set out what each side can verify about the other at connection time, with the threat model
- T-0694 Report the options to the operator and record the decision
- T-0695 Plan the child stories of E-0015 from the decision

## Notes

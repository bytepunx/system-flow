---
id: S-0278
type: story
nature: improvement
title: Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
status: backlog
owner: alex
created: 2026-10-05T04:40:47Z
updated: 2026-10-05T04:40:47Z
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
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-10-05T04:40:47Z
---
# S-0278 Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Goal

This story remediates [I-0074](../../../design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md), "Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0074 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0074 is closed with `flai issue close I-0074 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0074. time_lost_per_cycle 9m: 3m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-05T03:22:09Z, 0.1 days before this story; under one cycle counts as one).

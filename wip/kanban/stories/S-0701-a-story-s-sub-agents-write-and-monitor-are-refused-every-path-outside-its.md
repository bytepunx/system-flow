---
id: S-0701
type: story
nature: improvement
title: A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository
status: backlog
owner: alex
created: 2026-11-10T08:55:34Z
updated: 2026-11-10T08:55:34Z
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
    time_lost_per_cycle: 26m
    by: flai
    at: 2026-11-10T08:55:34Z
---
# S-0701 A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository

## Goal

This story remediates [I-0265](../../../design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md), "A story's sub-agents, Write, and Monitor are refused every path outside its worktree, though an experiment story may be written to work outside the repository". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0265 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0265 is closed with `flai issue close I-0265 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0265. time_lost_per_cycle 26m: 13m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-06T23:16:17Z, 3.4 days before this story; under one cycle counts as one).

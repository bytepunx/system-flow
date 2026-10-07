---
id: S-0297
type: story
nature: improvement
title: Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
status: backlog
owner: alex
created: 2026-10-06T19:46:46Z
updated: 2026-10-07T02:18:22Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-06T19:46:46Z
finalized:
  by: alex
  at: 2026-10-07T02:18:22Z
---
# S-0297 Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Goal

This story remediates [I-0092](../../../design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md), "Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0092 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0092 is closed with `flai issue close I-0092 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0092. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T19:46:45Z, 0 days before this story; under one cycle counts as one).

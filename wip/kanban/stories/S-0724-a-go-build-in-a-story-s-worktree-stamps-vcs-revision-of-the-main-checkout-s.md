---
id: S-0724
type: story
nature: remediation
title: A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit
status: backlog
owner: alex
created: 2026-11-11T08:46:42Z
updated: 2026-11-11T08:50:05Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-11T08:46:42Z
finalized:
  by: alex
  at: 2026-11-11T08:50:05Z
---
# S-0724 A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit

## Goal

This story remediates [I-0348](../../../design/issues/I-0348-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md), "A go build in a story's worktree stamps vcs.revision of the main checkout's HEAD, so a worktree's bin/flai names the wrong commit". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0348 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0348 is closed with `flai issue close I-0348 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0348. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T07:50:38Z, 0 days before this story; under one cycle counts as one).

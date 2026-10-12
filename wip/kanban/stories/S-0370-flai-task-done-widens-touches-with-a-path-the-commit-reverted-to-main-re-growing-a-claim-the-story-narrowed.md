---
id: S-0370
type: story
nature: remediation
title: flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed
status: backlog
owner: alex
created: 2026-10-12T16:52:44Z
updated: 2026-10-12T16:52:44Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T16:52:44Z
---
# S-0370 flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed

## Goal

This story remediates [I-0121](../../../design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md), "flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0121 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0121 is closed with `flai issue close I-0121 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0121. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T08:25:05Z, 4.4 days before this story; under one cycle counts as one).

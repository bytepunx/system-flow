---
id: S-0689
type: story
nature: remediation
title: flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached
status: backlog
owner: alex
created: 2026-11-10T08:49:00Z
updated: 2026-11-10T08:49:00Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T08:49:00Z
---
# S-0689 flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached

## Goal

This story remediates [I-0272](../../../design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md), "flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0272 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0272 is closed with `flai issue close I-0272 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0272. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:52Z, 3.2 days before this story; under one cycle counts as one).

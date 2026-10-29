---
id: S-0545
type: story
nature: improvement
title: flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply
status: backlog
owner: alex
created: 2026-10-29T18:11:17Z
updated: 2026-10-29T18:11:17Z
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
    time_lost_per_cycle: 8m
    by: flai
    at: 2026-10-29T18:11:17Z
---
# S-0545 flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply

## Goal

This story remediates [I-0232](../../../design/issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md), "flai serve wakes a story's agent that is paused on a thread for every task-done message from another story, each needing only a settled reply". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0232 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0232 is closed with `flai issue close I-0232 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0232. time_lost_per_cycle 8m: 4m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-29T15:42:30Z, 0.1 days before this story; under one cycle counts as one).

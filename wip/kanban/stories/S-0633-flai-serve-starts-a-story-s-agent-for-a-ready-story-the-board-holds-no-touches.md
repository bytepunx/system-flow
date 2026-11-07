---
id: S-0633
type: story
nature: improvement
title: flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it
status: backlog
owner: alex
created: 2026-11-07T23:41:06Z
updated: 2026-11-07T23:41:06Z
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
    at: 2026-11-07T23:41:06Z
---
# S-0633 flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it

## Goal

This story remediates [I-0301](../../../design/issues/I-0301-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md), "flai serve starts a story's agent for a ready story the board holds no-touches, and story_start then refuses it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0301 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0301 is closed with `flai issue close I-0301 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0301. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:44:09Z, 0.6 days before this story; under one cycle counts as one).

---
id: S-0304
type: story
nature: improvement
title: A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
status: backlog
owner: alex
created: 2026-10-07T01:07:09Z
updated: 2026-10-07T01:07:09Z
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
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-07T01:07:09Z
---
# S-0304 A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Goal

This story remediates [I-0097](../../../design/issues/I-0097-a-story-s-agent-that-runs-flai-stream-open-is-left-with-the-story-in-ready-so-permission-prompt-refuses-its-claude-write-until-it-moves-the-story-itself.md), "A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0097 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0097 is closed with `flai issue close I-0097 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0097. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T23:28:35Z, 0.1 days before this story; under one cycle counts as one).

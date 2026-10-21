---
id: S-0477
type: story
nature: improvement
title: flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent
status: backlog
owner: alex
created: 2026-10-21T11:59:27Z
updated: 2026-10-21T11:59:27Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T11:59:27Z
---
# S-0477 flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent

## Goal

This story remediates [I-0165](../../../design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md), "flai guard refuses a story's sub-agent flai components list and flai components serve, so a headless pass of the harness falls to the story's agent". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0165 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0165 is closed with `flai issue close I-0165 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0165. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-19T05:38:19Z, 2.3 days before this story; under one cycle counts as one).

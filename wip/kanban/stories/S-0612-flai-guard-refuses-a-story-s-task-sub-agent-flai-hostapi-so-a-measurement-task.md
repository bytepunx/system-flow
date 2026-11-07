---
id: S-0612
type: story
nature: improvement
title: flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing
status: backlog
owner: alex
created: 2026-11-07T23:38:52Z
updated: 2026-11-07T23:38:52Z
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
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-07T23:38:52Z
---
# S-0612 flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing

## Goal

This story remediates [I-0307](../../../design/issues/I-0307-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md), "flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0307 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0307 is closed with `flai issue close I-0307 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0307. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T22:43:39Z, 0 days before this story; under one cycle counts as one).

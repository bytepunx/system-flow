---
id: S-0834
type: story
nature: remediation
title: flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount
status: backlog
owner: alex
created: 2026-11-17T08:17:29Z
updated: 2026-11-17T08:17:29Z
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
    at: 2026-11-17T08:17:29Z
---
# S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount

## Goal

This story remediates [I-0361](../../../design/issues/I-0361-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md), "flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0361 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0361 is closed with `flai issue close I-0361 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0361. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T12:16:36Z, 3.8 days before this story; under one cycle counts as one).

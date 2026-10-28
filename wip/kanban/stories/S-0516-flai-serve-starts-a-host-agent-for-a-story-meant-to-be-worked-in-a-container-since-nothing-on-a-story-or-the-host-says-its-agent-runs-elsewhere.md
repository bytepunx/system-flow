---
id: S-0516
type: story
nature: remediation
title: flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere
status: backlog
owner: alex
created: 2026-10-28T12:52:40Z
updated: 2026-10-28T12:52:40Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 22
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 29
          cache_read: 1461236
          cache_write: 3927
          cost: 0.3216
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-10-28T12:52:40Z
---
# S-0516 flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere

## Goal

This story remediates [I-0214](../../../design/issues/I-0214-flai-serve-starts-a-host-agent-for-a-story-meant-to-be-worked-in-a-container-since-nothing-on-a-story-or-the-host-says-its-agent-runs-elsewhere.md), "flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0214 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0214 is closed with `flai issue close I-0214 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0214. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-28T12:27:00Z, 0 days before this story; under one cycle counts as one).

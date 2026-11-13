---
id: S-0730
type: story
nature: improvement
title: A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it
status: backlog
owner: alex
created: 2026-11-12T03:43:39Z
updated: 2026-11-12T03:43:39Z
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
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 20
          cache_read: 194240
          cache_write: 3302
          cost: 0.057
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: flai
    at: 2026-11-12T03:43:39Z
---
# S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it

## Goal

This story remediates [I-0352](../../../design/issues/I-0352-a-story-s-agent-can-overwrite-a-work-item-file-whole-with-write-front-matter.md), "A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0352 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0352 is closed with `flai issue close I-0352 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0352. time_lost_per_cycle 1m: 1m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T22:53:06Z, 0.2 days before this story; under one cycle counts as one).

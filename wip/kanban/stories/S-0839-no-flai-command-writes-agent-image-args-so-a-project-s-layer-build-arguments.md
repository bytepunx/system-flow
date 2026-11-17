---
id: S-0839
type: story
nature: improvement
title: No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand
status: backlog
owner: alex
created: 2026-11-17T08:32:50Z
updated: 2026-11-17T08:32:50Z
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
    at: 2026-11-17T08:32:50Z
---
# S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand

## Goal

This story remediates [I-0364](../../../design/issues/I-0364-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md), "No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0364 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0364 is closed with `flai issue close I-0364 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0364. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:00:47Z, 3.5 days before this story; under one cycle counts as one).

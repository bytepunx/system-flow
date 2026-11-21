---
id: S-0849
type: story
nature: remediation
title: flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount
status: backlog
owner: alex
created: 2026-11-21T18:57:50Z
updated: 2026-11-21T18:57:50Z
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
    at: 2026-11-21T18:57:50Z
---
# S-0849 flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount

## Goal

This story remediates [I-0405](../../../design/issues/I-0405-flai-kluster-deploy-job-teardown-job-and-flai-checks-run-job-from-a-story-s.md), "flai kluster deploy --job, teardown --job, and flai checks run --job from a story's container write their request under the main checkout's .flai-cache/kluster/requests, which the container does not mount". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0405 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0405 is closed with `flai issue close I-0405 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0405. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-18T08:00:16Z, 3.5 days before this story; under one cycle counts as one).

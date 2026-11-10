---
id: S-0667
type: story
nature: improvement
title: A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project
status: backlog
owner: alex
created: 2026-11-10T08:34:40Z
updated: 2026-11-10T08:48:52Z
transitions: []
tags: [docker, agent-container]
touches: [scripts/agent-container.sh, docs/operators/runbooks/agent-container.md, design/tech/docker.md, design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md, design/issues/summary.md]
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
    at: 2026-11-10T08:34:40Z
---
# S-0667 A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project

## Goal

This story remediates [I-0276](../../../design/issues/I-0276-a-docker-run-from-a-story-s-agent-s-shell-mounts-the-main-checkout-as-an-empty.md), "A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0276 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0276 is closed with `flai issue close I-0276 --reason` saying what fixed it

## Tasks
- T-2926 The cause of I-0276's empty mount is found on the host and a fix is proposed
- T-2927 A hand-run container sees the main checkout's files, and refuses with the cause when a mount comes up empty

## Notes

Cost of delay inputs set by flai from I-0276. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:43:57Z, 3.1 days before this story; under one cycle counts as one).

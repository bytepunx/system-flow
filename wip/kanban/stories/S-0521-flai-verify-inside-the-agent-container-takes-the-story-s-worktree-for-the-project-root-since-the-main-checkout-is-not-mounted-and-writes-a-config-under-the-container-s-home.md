---
id: S-0521
type: story
nature: remediation
title: flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home
status: backlog
owner: alex
created: 2026-10-28T12:53:27Z
updated: 2026-10-28T12:53:59Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-28T12:53:27Z
finalized:
  by: alex
  at: 2026-10-28T12:53:59Z
---
# S-0521 flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home

## Goal

This story remediates [I-0219](../../../design/issues/I-0219-flai-verify-inside-the-agent-container-takes-the-story-s-worktree-for-the-project-root-since-the-main-checkout-is-not-mounted-and-writes-a-config-under-the-container-s-home.md), "flai verify inside the agent container takes the story's worktree for the project root, since the main checkout is not mounted, and writes a config under the container's home". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0219 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0219 is closed with `flai issue close I-0219 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0219. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:31:04Z, 0 days before this story; under one cycle counts as one).

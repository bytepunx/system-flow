---
id: S-0262
type: story
nature: remediation
title: flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out
status: ready
owner: alex
created: 2026-10-04T21:41:59Z
updated: 2026-10-04T23:14:03Z
transitions:
  - to: ready
    at: 2026-10-04T23:14:03Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-04T21:41:59Z
finalized:
  by: alex
  at: 2026-10-04T23:14:01Z
---
# S-0262 flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out

## Goal

This story remediates [I-0072](../../../design/issues/I-0072-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md), "flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0072 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0072 is closed with `flai issue close I-0072 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0072. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-04T21:31:32Z, 0 days before this story; under one cycle counts as one).

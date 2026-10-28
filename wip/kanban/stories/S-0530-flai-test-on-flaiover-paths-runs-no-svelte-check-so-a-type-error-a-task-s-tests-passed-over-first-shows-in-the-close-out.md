---
id: S-0530
type: story
nature: improvement
title: flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out
status: backlog
owner: alex
created: 2026-10-28T12:54:31Z
updated: 2026-10-28T12:56:09Z
transitions: []
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
    at: 2026-10-28T12:54:31Z
finalized:
  by: alex
  at: 2026-10-28T12:56:09Z
---
# S-0530 flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out

## Goal

This story remediates [I-0206](../../../design/issues/I-0206-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md), "flai test on flaiover paths runs no svelte-check, so a type error a task's tests passed over first shows in the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0206 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0206 is closed with `flai issue close I-0206 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0206. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T22:32:29Z, 3.6 days before this story; under one cycle counts as one).

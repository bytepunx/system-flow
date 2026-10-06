---
id: S-0291
type: story
nature: improvement
title: The close-out's install smoke test failed once and passed when run alone, with no cause in its output
status: backlog
owner: alex
created: 2026-10-06T10:31:55Z
updated: 2026-10-06T10:31:55Z
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
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T10:31:55Z
---
# S-0291 The close-out's install smoke test failed once and passed when run alone, with no cause in its output

## Goal

This story remediates [I-0086](../../../design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md), "The close-out's install smoke test failed once and passed when run alone, with no cause in its output". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0086 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0086 is closed with `flai issue close I-0086 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0086. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T10:22:10Z, 0 days before this story; under one cycle counts as one).

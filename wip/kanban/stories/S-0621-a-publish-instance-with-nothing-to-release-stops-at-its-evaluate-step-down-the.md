---
id: S-0621
type: story
nature: improvement
title: A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time
status: backlog
owner: alex
created: 2026-11-07T23:39:50Z
updated: 2026-11-07T23:39:50Z
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
    time_lost_per_cycle: 24m
    by: flai
    at: 2026-11-07T23:39:50Z
---
# S-0621 A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Goal

This story remediates [I-0283](../../../design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md), "A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0283 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0283 is closed with `flai issue close I-0283 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0283. time_lost_per_cycle 24m: 8m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T08:16:08Z, 0.6 days before this story; under one cycle counts as one).

---
id: S-0347
type: story
nature: improvement
title: The close-out's last check that the branch contains main fails when flai commits wip on main during its run
status: backlog
owner: alex
created: 2026-10-08T08:08:22Z
updated: 2026-10-08T08:08:22Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-08T08:08:22Z
---
# S-0347 The close-out's last check that the branch contains main fails when flai commits wip on main during its run

## Goal

This story remediates [I-0119](../../../design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md), "The close-out's last check that the branch contains main fails when flai commits wip on main during its run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0119 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0119 is closed with `flai issue close I-0119 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0119. time_lost_per_cycle 30m: 10m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T04:29:28Z, 0.2 days before this story; under one cycle counts as one).

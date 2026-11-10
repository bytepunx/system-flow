---
id: S-0709
type: story
nature: improvement
title: flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs
status: backlog
owner: alex
created: 2026-11-10T09:55:57Z
updated: 2026-11-10T09:55:57Z
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
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-11-10T09:55:57Z
---
# S-0709 flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs

## Goal

This story remediates [I-0335](../../../design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md), "flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0335 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0335. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:03:43Z, 0 days before this story; under one cycle counts as one).

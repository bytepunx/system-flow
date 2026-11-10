---
id: S-0715
type: story
nature: remediation
title: "TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with \"directory not empty\", as I-0048's and I-0233's serve tests did"
status: backlog
owner: alex
created: 2026-11-10T10:55:22Z
updated: 2026-11-10T10:55:22Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T10:55:22Z
---
# S-0715 TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did

## Goal

This story remediates [I-0343](../../../design/issues/I-0343-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md), "TestARestartStartsWhatIsReadyAndNothingTwice fails its temp folder's cleanup with "directory not empty", as I-0048's and I-0233's serve tests did". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0343 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0343 is closed with `flai issue close I-0343 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0343. time_lost_per_cycle 12m: 12m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T10:41:09Z, 0 days before this story; under one cycle counts as one).

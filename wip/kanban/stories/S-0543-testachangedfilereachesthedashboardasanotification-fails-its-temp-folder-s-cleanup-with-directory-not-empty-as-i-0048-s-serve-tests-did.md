---
id: S-0543
type: story
nature: remediation
title: "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with \"directory not empty\", as I-0048's serve tests did"
status: backlog
owner: alex
created: 2026-10-29T17:25:41Z
updated: 2026-10-29T17:25:41Z
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
    at: 2026-10-29T17:25:41Z
---
# S-0543 TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Goal

This story remediates [I-0233](../../../design/issues/I-0233-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md), "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0233 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0233 is closed with `flai issue close I-0233 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0233. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-29T16:37:28Z, 0 days before this story; under one cycle counts as one).

---
id: S-0345
type: story
nature: remediation
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
status: backlog
owner: alex
created: 2026-10-08T08:08:19Z
updated: 2026-10-08T08:08:19Z
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
    time_lost_per_cycle: 1h12m
    by: flai
    at: 2026-10-08T08:08:19Z
---
# S-0345 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Goal

This story remediates [I-0117](../../../design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md), "A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0117 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0117 is closed with `flai issue close I-0117 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0117. time_lost_per_cycle 1h12m: 18m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T04:12:50Z, 0.2 days before this story; under one cycle counts as one).

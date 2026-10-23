---
id: S-0488
type: story
nature: remediation
title: TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve
status: backlog
owner: alex
created: 2026-10-23T17:53:08Z
updated: 2026-10-23T17:53:08Z
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
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-10-23T17:53:08Z
---
# S-0488 TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve

## Goal

This story remediates [I-0197](../../../design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md), "TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0197 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0197 is closed with `flai issue close I-0197 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0197. time_lost_per_cycle 7m: 7m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-23T05:12:51Z, 0.5 days before this story; under one cycle counts as one).

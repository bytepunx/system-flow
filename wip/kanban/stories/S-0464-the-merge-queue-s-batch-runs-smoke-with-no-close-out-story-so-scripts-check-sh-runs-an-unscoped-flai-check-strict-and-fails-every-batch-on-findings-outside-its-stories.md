---
id: S-0464
type: story
nature: remediation
title: The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories
status: backlog
owner: alex
created: 2026-10-21T07:19:58Z
updated: 2026-10-21T07:19:58Z
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
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-21T07:19:58Z
---
# S-0464 The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories

## Goal

This story remediates [I-0174](../../../design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md), "The merge queue's batch runs smoke with no CLOSE_OUT_STORY, so scripts/check.sh runs an unscoped flai check --strict and fails every batch on findings outside its stories". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0174 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0174 is closed with `flai issue close I-0174 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0174. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T06:29:55Z, 0 days before this story; under one cycle counts as one).

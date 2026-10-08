---
id: S-0326
type: story
nature: remediation
title: flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile
status: backlog
owner: alex
created: 2026-10-07T18:59:58Z
updated: 2026-10-07T18:59:58Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 37
          cache_read: 730566
          cache_write: 5839
          cost: 0.1815
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T18:59:58Z
---
# S-0326 flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Goal

This story remediates [I-0112](../../../design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md), "flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0112 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0112 is closed with `flai issue close I-0112 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0112. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T09:09:31Z, 0.4 days before this story; under one cycle counts as one).

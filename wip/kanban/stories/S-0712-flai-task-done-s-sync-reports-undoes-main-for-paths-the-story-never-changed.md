---
id: S-0712
type: story
nature: remediation
title: flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync
status: backlog
owner: alex
created: 2026-11-10T10:10:59Z
updated: 2026-11-10T10:10:59Z
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
    at: 2026-11-10T10:10:59Z
---
# S-0712 flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync

## Goal

This story remediates [I-0338](../../../design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md), "flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0338 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0338 is closed with `flai issue close I-0338 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0338. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:13:03Z, 0 days before this story; under one cycle counts as one).

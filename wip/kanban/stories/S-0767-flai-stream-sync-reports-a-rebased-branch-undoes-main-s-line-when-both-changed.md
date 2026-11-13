---
id: S-0767
type: story
nature: remediation
title: flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's
status: backlog
owner: alex
created: 2026-11-13T11:42:36Z
updated: 2026-11-13T11:42:36Z
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
    at: 2026-11-13T11:42:36Z
---
# S-0767 flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's

## Goal

This story remediates [I-0360](../../../design/issues/I-0360-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md), "flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0360 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0360 is closed with `flai issue close I-0360 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0360. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T09:21:49Z, 0.1 days before this story; under one cycle counts as one).

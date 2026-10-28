---
id: S-0529
type: story
nature: remediation
title: A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration
status: backlog
owner: alex
created: 2026-10-28T12:54:24Z
updated: 2026-10-28T12:54:24Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:54:24Z
---
# S-0529 A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Goal

This story remediates [I-0203](../../../design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md), "A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0203 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0203 is closed with `flai issue close I-0203 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0203. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T00:15:33Z, 4.5 days before this story; under one cycle counts as one).

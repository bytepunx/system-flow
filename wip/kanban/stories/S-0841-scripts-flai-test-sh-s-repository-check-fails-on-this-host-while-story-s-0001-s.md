---
id: S-0841
type: story
nature: remediation
title: scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged
status: backlog
owner: alex
created: 2026-11-17T08:32:52Z
updated: 2026-11-17T08:32:52Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-17T08:32:52Z
---
# S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged

## Goal

This story remediates [I-0370](../../../design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md), "scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0370 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0370 is closed with `flai issue close I-0370 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0370. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:39:18Z, 3.5 days before this story; under one cycle counts as one).

---
id: S-0590
type: story
nature: remediation
title: flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock
status: backlog
owner: alex
created: 2026-11-07T09:44:23Z
updated: 2026-11-07T09:44:23Z
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
    at: 2026-11-07T09:44:23Z
---
# S-0590 flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock

## Goal

This story remediates [I-0302](../../../design/issues/I-0302-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md), "flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0302 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0302 is closed with `flai issue close I-0302 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0302. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:50:28Z, 0 days before this story; under one cycle counts as one).

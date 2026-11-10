---
id: S-0690
type: story
nature: remediation
title: flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look
status: backlog
owner: alex
created: 2026-11-10T08:49:01Z
updated: 2026-11-10T08:49:01Z
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
    at: 2026-11-10T08:49:01Z
---
# S-0690 flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look

## Goal

This story remediates [I-0273](../../../design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md), "flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0273 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0273 is closed with `flai issue close I-0273 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0273. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T05:23:04Z, 3.1 days before this story; under one cycle counts as one).

---
id: S-0677
type: story
nature: improvement
title: "A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path"
status: backlog
owner: alex
created: 2026-11-10T08:34:54Z
updated: 2026-11-10T08:34:54Z
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
    time_lost_per_cycle: 2s
    by: flai
    at: 2026-11-10T08:34:54Z
---
# S-0677 A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path

## Goal

This story remediates [I-0312](../../../design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md), "A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0312 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0312 is closed with `flai issue close I-0312 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0312. time_lost_per_cycle 2s: 2s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T01:58:32Z, 2.3 days before this story; under one cycle counts as one).

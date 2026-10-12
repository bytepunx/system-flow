---
id: S-0379
type: story
nature: improvement
title: flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify
status: backlog
owner: alex
created: 2026-10-12T17:44:10Z
updated: 2026-10-12T17:44:10Z
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
    at: 2026-10-12T17:44:10Z
---
# S-0379 flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify

## Goal

This story remediates [I-0142](../../../design/issues/I-0142-flai-self-upgrade-list-and-flai-host-versions-download-two-assets-for-every-signed-release-to-mark-the-ones-they-cannot-verify.md), "flai self-upgrade --list and flai host versions download two assets for every signed release to mark the ones they cannot verify". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0142 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0142 is closed with `flai issue close I-0142 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0142. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T17:28:30Z, 0 days before this story; under one cycle counts as one).

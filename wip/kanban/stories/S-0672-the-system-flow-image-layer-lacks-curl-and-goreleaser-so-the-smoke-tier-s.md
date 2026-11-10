---
id: S-0672
type: story
nature: remediation
title: The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there
status: backlog
owner: alex
created: 2026-11-10T08:34:49Z
updated: 2026-11-10T08:34:49Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-10T08:34:49Z
---
# S-0672 The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there

## Goal

This story remediates [I-0281](../../../design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md), "The system-flow image layer lacks curl and GoReleaser, so the smoke tier's install-test fails inside the agent container and a Go story's close-out cannot pass there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0281 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0281 is closed with `flai issue close I-0281 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0281. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T07:12:34Z, 3.1 days before this story; under one cycle counts as one).

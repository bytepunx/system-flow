---
id: S-0520
type: story
nature: remediation
title: scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks
status: backlog
owner: alex
created: 2026-10-28T12:53:17Z
updated: 2026-10-28T12:53:17Z
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
    at: 2026-10-28T12:53:17Z
---
# S-0520 scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks

## Goal

This story remediates [I-0218](../../../design/issues/I-0218-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md), "scripts/close-out.sh stops at its first step inside the agent container, since scripts/flai.sh builds flai with go, which the image lacks". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0218 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0218 is closed with `flai issue close I-0218 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0218. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:31:00Z, 0 days before this story; under one cycle counts as one).

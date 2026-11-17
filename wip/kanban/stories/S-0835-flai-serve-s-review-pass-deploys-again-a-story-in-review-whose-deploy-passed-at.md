---
id: S-0835
type: story
nature: remediation
title: "flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone"
status: backlog
owner: alex
created: 2026-11-17T08:17:30Z
updated: 2026-11-17T08:17:30Z
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
    at: 2026-11-17T08:17:30Z
---
# S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone

## Goal

This story remediates [I-0362](../../../design/issues/I-0362-flai-serve-s-review-pass-deploys-again-a-story-in-review-whose-deploy-passed-at.md), "flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0362 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0362 is closed with `flai issue close I-0362 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0362. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T13:28:22Z, 3.8 days before this story; under one cycle counts as one).

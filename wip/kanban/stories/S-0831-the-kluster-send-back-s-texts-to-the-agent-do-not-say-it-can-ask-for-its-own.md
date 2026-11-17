---
id: S-0831
type: story
nature: improvement
title: The kluster send-back's texts to the agent do not say it can ask for its own deploy before review
status: backlog
owner: alex
created: 2026-11-17T08:17:25Z
updated: 2026-11-17T08:17:25Z
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
    at: 2026-11-17T08:17:25Z
---
# S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review

## Goal

This story remediates [I-0357](../../../design/issues/I-0357-the-kluster-send-back-s-texts-to-the-agent-do-not-say-it-can-ask-for-its-own.md), "The kluster send-back's texts to the agent do not say it can ask for its own deploy before review". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0357 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0357 is closed with `flai issue close I-0357 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0357. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:12Z, 4 days before this story; under one cycle counts as one).

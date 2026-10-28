---
id: S-0531
type: story
nature: improvement
title: A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer
status: backlog
owner: alex
created: 2026-10-28T12:54:38Z
updated: 2026-10-28T12:54:38Z
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
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-10-28T12:54:38Z
---
# S-0531 A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer

## Goal

This story remediates [I-0209](../../../design/issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md), "A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0209 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0209 is closed with `flai issue close I-0209 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0209. time_lost_per_cycle 1s: 650ms per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-26T04:42:46Z, 2.3 days before this story; under one cycle counts as one).

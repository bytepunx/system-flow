---
id: S-0685
type: story
nature: remediation
title: The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails
status: backlog
owner: alex
created: 2026-11-10T08:37:11Z
updated: 2026-11-10T08:37:11Z
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
    time_lost_per_cycle: 11m
    by: flai
    at: 2026-11-10T08:37:11Z
---
# S-0685 The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails

## Goal

This story remediates [I-0237](../../../design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md), "The smoke tier's OpenRouter gateway check fails on a 402 because the operator's OpenRouter key is out of credits, so every close-out and batch that runs smoke fails". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0237 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0237 is closed with `flai issue close I-0237 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0237. time_lost_per_cycle 11m: 15m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-10-31T10:10:59Z, 9.9 days before this story).

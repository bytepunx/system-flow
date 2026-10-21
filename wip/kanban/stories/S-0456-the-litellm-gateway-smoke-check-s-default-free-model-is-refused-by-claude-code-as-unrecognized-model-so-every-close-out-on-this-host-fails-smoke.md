---
id: S-0456
type: story
nature: remediation
title: The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke
status: backlog
owner: alex
created: 2026-10-21T02:08:41Z
updated: 2026-10-21T02:08:41Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T02:08:41Z
---
# S-0456 The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke

## Goal

This story remediates [I-0170](../../../design/issues/I-0170-the-litellm-gateway-smoke-check-s-default-free-model-is-refused-by-claude-code-as-unrecognized-model-so-every-close-out-on-this-host-fails-smoke.md), "The LiteLLM gateway smoke check's default free model is refused by Claude Code as unrecognized_model, so every close-out on this host fails smoke". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0170 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0170 is closed with `flai issue close I-0170 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0170. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T19:17:06Z, 0.3 days before this story; under one cycle counts as one).

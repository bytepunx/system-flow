---
id: S-0455
type: story
nature: remediation
title: The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model
status: backlog
owner: alex
created: 2026-10-21T00:29:30Z
updated: 2026-10-21T00:37:14Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T00:29:30Z
finalized:
  by: alex
  at: 2026-10-21T00:37:14Z
---
# S-0455 The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model

## Goal

This story remediates [I-0171](../../../design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md), "The LiteLLM gateway smoke check fails because Claude Code 2.1.296 refuses the default haiku model name as unrecognized_model". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0171 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0171 is closed with `flai issue close I-0171 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0171. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-20T20:07:32Z, 0.2 days before this story; under one cycle counts as one).

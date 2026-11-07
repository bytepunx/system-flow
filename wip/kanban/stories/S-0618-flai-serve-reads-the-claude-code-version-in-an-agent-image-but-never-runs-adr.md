---
id: S-0618
type: story
nature: remediation
title: flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of
status: backlog
owner: alex
created: 2026-11-07T23:39:32Z
updated: 2026-11-07T23:39:32Z
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
    at: 2026-11-07T23:39:32Z
---
# S-0618 flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of

## Goal

This story remediates [I-0253](../../../design/issues/I-0253-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md), "flai serve reads the Claude Code version in an agent image but never runs ADR-0106's permission_prompt check inside it, so a new image's version is only warned of". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0253 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0253 is closed with `flai issue close I-0253 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0253. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T07:30:18Z, 5.7 days before this story; under one cycle counts as one).

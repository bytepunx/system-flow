---
id: S-0670
type: story
nature: remediation
title: "scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools"
status: backlog
owner: alex
created: 2026-11-10T08:34:46Z
updated: 2026-11-10T08:34:46Z
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
    at: 2026-11-10T08:34:46Z
---
# S-0670 scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools

## Goal

This story remediates [I-0279](../../../design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md), "scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0279 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0279 is closed with `flai issue close I-0279 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0279. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:58:47Z, 3.1 days before this story; under one cycle counts as one).

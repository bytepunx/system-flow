---
id: S-0687
type: story
nature: remediation
title: A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools
status: backlog
owner: alex
created: 2026-11-10T08:37:14Z
updated: 2026-11-10T08:37:14Z
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
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:37:14Z
---
# S-0687 A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools

## Goal

This story remediates [I-0256](../../../design/issues/I-0256-a-strategic-agent-s-container-runs-flai-with-no-host-configuration-so-the.md), "A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0256 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0256 is closed with `flai issue close I-0256 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0256. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:36Z, 7.6 days before this story).

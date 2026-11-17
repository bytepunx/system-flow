---
id: S-0830
type: story
nature: remediation
title: The dashboard's kluster setup form does not show kluster.agent_namespaces
status: backlog
owner: alex
created: 2026-11-17T08:17:22Z
updated: 2026-11-17T08:17:22Z
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
    at: 2026-11-17T08:17:22Z
---
# S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces

## Goal

This story remediates [I-0356](../../../design/issues/I-0356-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md), "The dashboard's kluster setup form does not show kluster.agent_namespaces". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0356 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0356 is closed with `flai issue close I-0356 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0356. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:09Z, 4 days before this story; under one cycle counts as one).

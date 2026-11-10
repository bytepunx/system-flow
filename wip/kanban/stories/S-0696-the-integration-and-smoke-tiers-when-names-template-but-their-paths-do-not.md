---
id: S-0696
type: story
nature: remediation
title: "The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say"
status: backlog
owner: alex
created: 2026-11-10T08:54:31Z
updated: 2026-11-10T08:54:31Z
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
    at: 2026-11-10T08:54:31Z
---
# S-0696 The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say

## Goal

This story remediates [I-0270](../../../design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md), "The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0270 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0270 is closed with `flai issue close I-0270 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0270. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:00Z, 3.2 days before this story; under one cycle counts as one).

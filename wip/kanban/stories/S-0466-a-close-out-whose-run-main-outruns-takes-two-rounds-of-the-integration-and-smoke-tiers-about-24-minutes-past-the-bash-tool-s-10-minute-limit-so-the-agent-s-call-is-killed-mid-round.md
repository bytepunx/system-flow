---
id: S-0466
type: story
nature: improvement
title: A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round
status: backlog
owner: alex
created: 2026-10-21T11:57:48Z
updated: 2026-10-21T11:57:48Z
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
    time_lost_per_cycle: 25m
    by: flai
    at: 2026-10-21T11:57:48Z
---
# S-0466 A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round

## Goal

This story remediates [I-0148](../../../design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md), "A close-out whose run main outruns takes two rounds of the integration and smoke tiers, about 24 minutes, past the Bash tool's 10-minute limit, so the agent's call is killed mid-round". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0148 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0148 is closed with `flai issue close I-0148 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0148. time_lost_per_cycle 25m: 25m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-15T00:07:49Z, 6.5 days before this story; under one cycle counts as one).

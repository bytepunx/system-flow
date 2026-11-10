---
id: S-0708
type: story
nature: remediation
title: "TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name"
status: backlog
owner: alex
created: 2026-11-10T09:55:53Z
updated: 2026-11-10T09:55:53Z
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
    at: 2026-11-10T09:55:53Z
---
# S-0708 TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name

## Goal

This story remediates [I-0334](../../../design/issues/I-0334-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md), "TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0334 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0334 is closed with `flai issue close I-0334 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0334. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T09:01:59Z, 0 days before this story; under one cycle counts as one).

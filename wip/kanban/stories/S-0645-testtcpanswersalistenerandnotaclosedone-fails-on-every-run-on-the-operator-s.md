---
id: S-0645
type: story
nature: remediation
title: TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there
status: backlog
owner: alex
created: 2026-11-09T21:50:55Z
updated: 2026-11-09T21:50:55Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-09T21:50:55Z
---
# S-0645 TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Goal

This story remediates [I-0321](../../../design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md), "TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0321 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0321 is closed with `flai issue close I-0321 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0321. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-09T17:13:14Z, 0.2 days before this story; under one cycle counts as one).

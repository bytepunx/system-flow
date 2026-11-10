---
id: S-0649
type: story
nature: remediation
title: TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611
status: backlog
owner: alex
created: 2026-11-10T05:12:08Z
updated: 2026-11-10T05:12:08Z
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
    time_lost_per_cycle: 36m
    by: flai
    at: 2026-11-10T05:12:08Z
---
# S-0649 TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Goal

This story remediates [I-0324](../../../design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md), "TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0324 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0324 is closed with `flai issue close I-0324 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0324. time_lost_per_cycle 36m: 12m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-09T22:00:56Z, 0.3 days before this story; under one cycle counts as one).

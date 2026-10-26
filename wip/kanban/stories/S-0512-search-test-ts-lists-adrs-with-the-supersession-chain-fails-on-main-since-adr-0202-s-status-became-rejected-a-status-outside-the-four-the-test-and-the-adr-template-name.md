---
id: S-0512
type: story
nature: remediation
title: search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name
status: backlog
owner: alex
created: 2026-10-26T16:03:03Z
updated: 2026-10-26T16:03:03Z
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
    at: 2026-10-26T16:03:03Z
---
# S-0512 search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name

## Goal

This story remediates [I-0208](../../../design/issues/I-0208-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md), "search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0208 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0208 is closed with `flai issue close I-0208 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0208. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-26T04:35:14Z, 0.5 days before this story; under one cycle counts as one).

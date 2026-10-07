---
id: S-0290
type: story
nature: remediation
title: TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run
status: backlog
owner: alex
created: 2026-10-06T09:56:52Z
updated: 2026-10-07T02:19:37Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T09:56:52Z
finalized:
  by: alex
  at: 2026-10-07T02:19:37Z
---
# S-0290 TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run

## Goal

This story remediates [I-0079](../../../design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md), "TestRoundTripRepositoryItems reads the live main checkout and fails a close-out when another agent edits a story mid-run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0079 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0079 is closed with `flai issue close I-0079 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0079. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T05:55:58Z, 1.2 days before this story; under one cycle counts as one).

---
id: S-0377
type: story
nature: remediation
title: "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost"
status: backlog
owner: alex
created: 2026-10-12T16:53:00Z
updated: 2026-10-12T16:53:00Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-12T16:53:00Z
---
# S-0377 TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost

## Goal

This story remediates [I-0140](../../../design/issues/I-0140-testactivitywriterstaketurns-fails-about-one-run-in-three-two-writers-that-both-find-the-activity-lock-stale-can-both-take-it-and-an-entry-is-lost.md), "TestActivityWritersTakeTurns fails about one run in three: two writers that both find the activity lock stale can both take it, and an entry is lost". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0140 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0140 is closed with `flai issue close I-0140 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0140. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-12T16:22:16Z, 0 days before this story; under one cycle counts as one).

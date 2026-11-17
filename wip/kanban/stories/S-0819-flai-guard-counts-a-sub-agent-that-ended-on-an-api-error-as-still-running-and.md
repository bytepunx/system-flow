---
id: S-0819
type: story
nature: remediation
title: flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out
status: backlog
owner: alex
created: 2026-11-17T02:47:36Z
updated: 2026-11-17T02:47:36Z
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
    at: 2026-11-17T02:47:36Z
---
# S-0819 flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out

## Goal

This story remediates [I-0376](../../../design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md), "flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0376 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0376 is closed with `flai issue close I-0376 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0376. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T01:05:45Z, 2.1 days before this story; under one cycle counts as one).

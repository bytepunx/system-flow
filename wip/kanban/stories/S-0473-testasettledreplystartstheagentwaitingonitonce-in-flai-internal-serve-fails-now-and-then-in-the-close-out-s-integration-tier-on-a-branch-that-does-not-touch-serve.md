---
id: S-0473
type: story
nature: remediation
title: TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve
status: backlog
owner: alex
created: 2026-10-21T11:58:53Z
updated: 2026-10-21T12:04:09Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 45m
    by: flai
    at: 2026-10-21T11:58:53Z
finalized:
  by: alex
  at: 2026-10-21T12:04:09Z
---
# S-0473 TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve

## Goal

This story remediates [I-0158](../../../design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md), "TestASettledReplyStartsTheAgentWaitingOnItOnce in flai/internal/serve fails now and then in the close-out's integration tier on a branch that does not touch serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0158 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0158 is closed with `flai issue close I-0158 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0158. time_lost_per_cycle 45m: 9m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-17T04:54:11Z, 4.3 days before this story; under one cycle counts as one).

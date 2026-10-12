---
id: S-0372
type: story
nature: improvement
title: A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other
status: backlog
owner: alex
created: 2026-10-12T16:52:49Z
updated: 2026-10-12T16:52:49Z
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
    at: 2026-10-12T16:52:49Z
---
# S-0372 A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other

## Goal

This story remediates [I-0125](../../../design/issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md), "A settled conversation between two stories' agents keeps waking each in turn, because every reply awaits the other". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0125 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0125 is closed with `flai issue close I-0125 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0125. time_lost_per_cycle 36m: 9m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T10:21:38Z, 4.3 days before this story; under one cycle counts as one).

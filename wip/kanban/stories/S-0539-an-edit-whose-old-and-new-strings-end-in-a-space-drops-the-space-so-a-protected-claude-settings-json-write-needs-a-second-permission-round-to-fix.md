---
id: S-0539
type: story
nature: improvement
title: An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix
status: backlog
owner: alex
created: 2026-10-29T02:01:44Z
updated: 2026-10-29T02:01:44Z
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
    at: 2026-10-29T02:01:44Z
---
# S-0539 An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix

## Goal

This story remediates [I-0223](../../../design/issues/I-0223-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md), "An Edit whose old and new strings end in a space drops the space, so a protected .claude/settings.json write needs a second permission round to fix". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0223 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0223 is closed with `flai issue close I-0223 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0223. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T22:40:11Z, 0.1 days before this story; under one cycle counts as one).

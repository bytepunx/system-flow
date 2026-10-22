---
id: S-0472
type: story
nature: improvement
title: flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server
status: backlog
owner: alex
created: 2026-10-21T11:58:44Z
updated: 2026-10-22T20:53:07Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T11:58:44Z
finalized:
  by: alex
  at: 2026-10-22T20:53:07Z
---
# S-0472 flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server

## Goal

This story remediates [I-0156](../../../design/issues/I-0156-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md), "flai guard refuses a story's sub-agents flai components list and serve, which only read or run the project's own dev server". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0156 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0156 is closed with `flai issue close I-0156 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0156. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T00:41:49Z, 4.5 days before this story; under one cycle counts as one).

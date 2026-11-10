---
id: S-0684
type: story
nature: improvement
title: flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest
status: backlog
owner: alex
created: 2026-11-10T08:37:09Z
updated: 2026-11-10T08:37:09Z
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
    time_lost_per_cycle: 43s
    by: flai
    at: 2026-11-10T08:37:09Z
---
# S-0684 flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest

## Goal

This story remediates [I-0234](../../../design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md), "flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0234 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0234 is closed with `flai issue close I-0234 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0234. time_lost_per_cycle 43s: 1m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-10-31T09:50:04Z, 9.9 days before this story).

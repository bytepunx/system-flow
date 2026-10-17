---
id: S-0445
type: story
nature: improvement
title: The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
status: backlog
owner: alex
created: 2026-10-17T16:03:15Z
updated: 2026-10-17T16:03:15Z
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
    at: 2026-10-17T16:03:15Z
---
# S-0445 The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Goal

This story remediates [I-0157](../../../design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md), "The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0157 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0157 is closed with `flai issue close I-0157 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0157. time_lost_per_cycle 15m: 3m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-17T01:51:51Z, 0.6 days before this story; under one cycle counts as one).

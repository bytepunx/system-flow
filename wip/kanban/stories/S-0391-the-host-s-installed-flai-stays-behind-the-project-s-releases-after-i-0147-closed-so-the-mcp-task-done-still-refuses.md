---
id: S-0391
type: story
nature: improvement
title: The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
status: backlog
owner: alex
created: 2026-10-15T05:02:18Z
updated: 2026-10-15T05:02:18Z
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
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-15T05:02:18Z
---
# S-0391 The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Goal

This story remediates [I-0149](../../../design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md), "The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0149 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0149 is closed with `flai issue close I-0149 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0149. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-15T04:41:02Z, 0 days before this story; under one cycle counts as one).

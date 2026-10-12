---
id: S-0363
type: story
nature: remediation
title: An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
status: backlog
owner: alex
created: 2026-10-12T14:12:30Z
updated: 2026-10-12T14:12:30Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 27
          cache_read: 268757
          cache_write: 3892
          cost: 0.0687
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T14:12:30Z
---
# S-0363 An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart

## Goal

This story remediates [I-0127](../../../design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md), "An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0127 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0127 is closed with `flai issue close I-0127 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0127. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:41:45Z, 4.1 days before this story; under one cycle counts as one).

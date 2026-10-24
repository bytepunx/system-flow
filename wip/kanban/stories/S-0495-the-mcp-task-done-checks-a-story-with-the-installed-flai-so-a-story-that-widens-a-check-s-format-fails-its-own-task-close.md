---
id: S-0495
type: story
nature: improvement
title: The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close
status: backlog
owner: alex
created: 2026-10-24T20:28:57Z
updated: 2026-10-24T20:28:57Z
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
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-24T20:28:57Z
---
# S-0495 The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close

## Goal

This story remediates [I-0194](../../../design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md), "The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0194 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0194 is closed with `flai issue close I-0194 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0194. time_lost_per_cycle 20m: 4m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-10-22T22:06:51Z, 1.9 days before this story; under one cycle counts as one).

---
id: S-0390
type: story
nature: improvement
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
status: backlog
owner: alex
created: 2026-10-15T02:55:10Z
updated: 2026-10-15T02:55:10Z
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
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-15T02:55:10Z
---
# S-0390 The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Goal

This story remediates [I-0147](../../../design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md), "The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0147 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0147 is closed with `flai issue close I-0147 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0147. time_lost_per_cycle 12m: 3m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-13T21:48:16Z, 1.2 days before this story; under one cycle counts as one).

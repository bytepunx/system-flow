---
id: S-0261
type: story
nature: improvement
title: The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
status: backlog
owner: alex
created: 2026-10-04T20:34:55Z
updated: 2026-10-04T20:34:55Z
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
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-04T20:34:55Z
---
# S-0261 The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file

## Goal

This story remediates [I-0068](../../../design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md), "The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0068 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0068 is closed with `flai issue close I-0068 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0068. time_lost_per_cycle 6m: 3m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-03T21:21:24Z, 1 day before this story; under one cycle counts as one).

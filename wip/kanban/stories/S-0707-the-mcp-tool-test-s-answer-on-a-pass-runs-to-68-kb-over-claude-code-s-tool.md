---
id: S-0707
type: story
nature: improvement
title: The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given
status: backlog
owner: alex
created: 2026-11-10T09:11:42Z
updated: 2026-11-10T09:11:42Z
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
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T09:11:42Z
---
# S-0707 The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given

## Goal

This story remediates [I-0260](../../../design/issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md), "The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0260 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0260 is closed with `flai issue close I-0260 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0260. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:44:19Z, 7.6 days before this story).

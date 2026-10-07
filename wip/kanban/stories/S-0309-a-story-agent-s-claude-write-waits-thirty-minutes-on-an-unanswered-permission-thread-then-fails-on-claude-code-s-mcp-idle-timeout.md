---
id: S-0309
type: story
nature: improvement
title: A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout
status: backlog
owner: alex
created: 2026-10-07T06:48:44Z
updated: 2026-10-07T06:48:44Z
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
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-07T06:48:44Z
---
# S-0309 A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout

## Goal

This story remediates [I-0103](../../../design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md), "A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0103 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0103 is closed with `flai issue close I-0103 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0103. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T04:55:43Z, 0.1 days before this story; under one cycle counts as one).

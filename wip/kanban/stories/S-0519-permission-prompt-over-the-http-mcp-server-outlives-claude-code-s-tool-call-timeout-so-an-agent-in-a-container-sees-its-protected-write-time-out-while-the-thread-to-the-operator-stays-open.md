---
id: S-0519
type: story
nature: remediation
title: permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open
status: backlog
owner: alex
created: 2026-10-28T12:53:05Z
updated: 2026-10-28T12:53:05Z
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
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-28T12:53:05Z
---
# S-0519 permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open

## Goal

This story remediates [I-0217](../../../design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md), "permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0217 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0217 is closed with `flai issue close I-0217 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0217. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:30:56Z, 0 days before this story; under one cycle counts as one).

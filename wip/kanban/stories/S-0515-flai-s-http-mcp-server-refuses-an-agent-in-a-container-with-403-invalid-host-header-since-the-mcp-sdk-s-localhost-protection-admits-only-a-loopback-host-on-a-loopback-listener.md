---
id: S-0515
type: story
nature: remediation
title: flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener
status: backlog
owner: alex
created: 2026-10-28T12:52:34Z
updated: 2026-10-28T12:57:10Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-28T12:52:34Z
finalized:
  by: alex
  at: 2026-10-28T12:57:10Z
---
# S-0515 flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener

## Goal

This story remediates [I-0213](../../../design/issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md), "flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0213 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0213 is closed with `flai issue close I-0213 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0213. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:22:29Z, 0 days before this story; under one cycle counts as one).

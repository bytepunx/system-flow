---
id: S-0697
type: story
nature: remediation
title: flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go
status: backlog
owner: alex
created: 2026-11-10T08:54:33Z
updated: 2026-11-10T08:54:33Z
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
    at: 2026-11-10T08:54:33Z
---
# S-0697 flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go

## Goal

This story remediates [I-0271](../../../design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md), "flai host restarts a project's MCP server on the same port while another process holds it, so the project has no MCP server until that process lets go". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0271 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0271 is closed with `flai issue close I-0271 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0271. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:15Z, 3.2 days before this story; under one cycle counts as one).

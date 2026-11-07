---
id: S-0632
type: story
nature: remediation
title: A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it
status: backlog
owner: alex
created: 2026-11-07T23:40:59Z
updated: 2026-11-07T23:40:59Z
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
    at: 2026-11-07T23:40:59Z
---
# S-0632 A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it

## Goal

This story remediates [I-0299](../../../design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md), "A role run the stdio MCP server starts through role_run is never reaped, so flai serve counts its zombie as running and never settles it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0299 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0299 is closed with `flai issue close I-0299 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0299. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:21:58Z, 0.6 days before this story; under one cycle counts as one).

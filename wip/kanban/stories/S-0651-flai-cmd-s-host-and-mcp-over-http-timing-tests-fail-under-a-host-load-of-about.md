---
id: S-0651
type: story
nature: improvement
title: flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200
status: backlog
owner: alex
created: 2026-11-10T05:44:44Z
updated: 2026-11-10T05:44:44Z
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
    time_lost_per_cycle: 56m
    by: flai
    at: 2026-11-10T05:44:44Z
---
# S-0651 flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200

## Goal

This story remediates [I-0309](../../../design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md), "flai/cmd's host and MCP-over-HTTP timing tests fail under a host load of about 200". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0309 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0309 is closed with `flai issue close I-0309 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0309. time_lost_per_cycle 56m: 14m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T00:27:14Z, 2.2 days before this story; under one cycle counts as one).

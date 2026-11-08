---
id: S-0619
type: story
nature: improvement
title: An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal
status: backlog
owner: alex
created: 2026-11-07T23:39:38Z
updated: 2026-11-07T23:39:38Z
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
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-07T23:39:38Z
---
# S-0619 An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal

## Goal

This story remediates [I-0254](../../../design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md), "An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0254 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0254 is closed with `flai issue close I-0254 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0254. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T07:30:21Z, 5.7 days before this story; under one cycle counts as one).

---
id: S-0308
type: story
nature: improvement
title: golangci-lint fails at once when another story's agent is running it on the same host
status: backlog
owner: alex
created: 2026-10-07T01:07:14Z
updated: 2026-10-07T01:07:14Z
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
    at: 2026-10-07T01:07:14Z
---
# S-0308 golangci-lint fails at once when another story's agent is running it on the same host

## Goal

This story remediates [I-0101](../../../design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md), "golangci-lint fails at once when another story's agent is running it on the same host". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0101 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0101 is closed with `flai issue close I-0101 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0101. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T00:34:48Z, 0 days before this story; under one cycle counts as one).

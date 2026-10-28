---
id: S-0518
type: story
nature: remediation
title: flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start
status: backlog
owner: alex
created: 2026-10-28T12:52:56Z
updated: 2026-10-28T12:52:56Z
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
    at: 2026-10-28T12:52:56Z
---
# S-0518 flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start

## Goal

This story remediates [I-0216](../../../design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md), "flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0216 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0216 is closed with `flai issue close I-0216 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0216. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-28T12:30:51Z, 0 days before this story; under one cycle counts as one).

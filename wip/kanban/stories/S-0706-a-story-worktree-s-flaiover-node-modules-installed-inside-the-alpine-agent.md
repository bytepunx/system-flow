---
id: S-0706
type: story
nature: remediation
title: A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse
status: backlog
owner: alex
created: 2026-11-10T09:11:37Z
updated: 2026-11-10T09:11:37Z
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
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T09:11:37Z
---
# S-0706 A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse

## Goal

This story remediates [I-0259](../../../design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md), "A story worktree's flaiover/node_modules installed inside the Alpine agent container holds only musl native bindings, so the host's flaiover tier then fails, and the reverse". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0259 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0259 is closed with `flai issue close I-0259 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0259. time_lost_per_cycle 9m: 10m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:22:35Z, 7.6 days before this story).

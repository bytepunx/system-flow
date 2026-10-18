---
id: S-0447
type: story
nature: remediation
title: The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier
status: backlog
owner: alex
created: 2026-10-18T02:19:04Z
updated: 2026-10-18T02:19:04Z
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
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-10-18T02:19:04Z
---
# S-0447 The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier

## Goal

This story remediates [I-0162](../../../design/issues/I-0162-the-components-dev-server-test-picks-a-free-port-and-closes-it-before-the-fake-vite-binds-it-so-a-concurrent-test-run-on-the-host-can-take-the-port-and-fail-the-integration-tier.md), "The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0162 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0162 is closed with `flai issue close I-0162 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0162. time_lost_per_cycle 7m: 7m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T18:10:30Z, 0.3 days before this story; under one cycle counts as one).

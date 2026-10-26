---
id: S-0499
type: story
nature: remediation
title: "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed"
status: backlog
owner: alex
created: 2026-10-26T04:01:24Z
updated: 2026-10-26T04:01:24Z
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
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-26T04:01:24Z
---
# S-0499 events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Goal

This story remediates [I-0181](../../../design/issues/I-0181-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md), "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0181 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0181 is closed with `flai issue close I-0181 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0181. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T11:42:36Z, 4.7 days before this story; under one cycle counts as one).

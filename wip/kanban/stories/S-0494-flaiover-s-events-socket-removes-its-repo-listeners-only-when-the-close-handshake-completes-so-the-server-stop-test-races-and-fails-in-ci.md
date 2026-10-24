---
id: S-0494
type: story
nature: remediation
title: flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI
status: ready
owner: alex
created: 2026-10-24T13:51:34Z
updated: 2026-10-24T13:51:45Z
transitions:
  - to: ready
    at: 2026-10-24T13:51:45Z
    by: system-flow
tags: [dashboard]
topics: [dashboard, testing]
touches: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, design/issues/I-0205-flaiover-s-events-socket-removes-its-repo-listeners-only-when-the-close-handshake-completes-so-the-server-stop-test-races-and-fails-in-ci.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-24T13:51:34Z
finalized:
  by: system-flow
  at: 2026-10-24T13:51:41Z
---
# S-0494 flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI

## Goal

This story remediates [I-0205](../../../design/issues/I-0205-flaiover-s-events-socket-removes-its-repo-listeners-only-when-the-close-handshake-completes-so-the-server-stop-test-races-and-fails-in-ci.md), "flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0205 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0205 is closed with `flai issue close I-0205 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0205. time_lost_per_cycle 15m carried over from I-0205's Impact section. Its Impact time_lost_per_cycle was taken rather than the 15m derived from its cost and count.

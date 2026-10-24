---
id: S-0494
type: story
nature: remediation
title: flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI
status: done
owner: alex
created: 2026-10-24T13:51:34Z
updated: 2026-10-24T13:58:03Z
transitions:
  - to: ready
    at: 2026-10-24T13:51:45Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T13:51:55Z
    by: agent-S-0494
  - to: review
    at: 2026-10-24T13:56:38Z
    by: agent-S-0494
  - to: done
    at: 2026-10-24T13:58:03Z
    by: alex
tags: [dashboard]
topics: [dashboard, testing]
touches: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, design/issues/I-0205-flaiover-s-events-socket-removes-its-repo-listeners-only-when-the-close-handshake-completes-so-the-server-stop-test-races-and-fails-in-ci.md, design/issues/summary.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 295
  priced_by: harness
  turns:
    - day: 2026-10-24
      ceremony: 1
      hand_edits: 2
      work: 23
  models:
    - model: claude-opus-5-5
      input: 54
      output: 9943
      cache_read: 3385448
      cache_write: 140054
      cost: 1.9966
  strategic:
    - kind: orchestrator
      seconds: 507
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 794
          cache_read: 6811489
          cache_write: 10635
          cost: 1.4958
        - model: claude-sonnet-5-5
          input: 6
          output: 48
          cache_read: 39794
          cache_write: 37810
          cost: 0.0531
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
- [x] The cause I-0205 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0205 is closed with `flai issue close I-0205 --reason` saying what fixed it

## Tasks
- T-2107 EventsSocket.close tears each socket down at once, removing its Repo listeners before it sends 1001

## Notes

Cost of delay inputs set by flai from I-0205. time_lost_per_cycle 15m carried over from I-0205's Impact section. Its Impact time_lost_per_cycle was taken rather than the 15m derived from its cost and count.

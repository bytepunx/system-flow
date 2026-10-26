---
id: S-0499
type: story
nature: remediation
title: "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed"
status: backlog
owner: alex
created: 2026-10-26T04:01:24Z
updated: 2026-10-26T04:57:43Z
transitions: []
tags: [dashboard]
topics: [dashboard, testing]
touches: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, design/issues/I-0181-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md, design/issues/summary.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 232
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 56
          output: 971
          cache_read: 10635111
          cache_write: 22634
          cost: 2.3366
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-26T04:01:24Z
  value: 7.5
  by: planner-S-0499
  at: 2026-10-26T04:20:24Z
forecast:
  duration: 30m
  delivery: 2026-10-26T18:50:00Z
  basis: "Its own forecast of 30m; 21st in the pull order with an in-progress limit of 5, behind S-0498, S-0504, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424 and S-0425."
  by: flai
  at: 2026-10-26T04:57:43Z
finalized:
  by: orchestrator
  at: 2026-10-26T04:21:34Z
---
# S-0499 events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Goal

This story remediates [I-0181](../../../design/issues/I-0181-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md), "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0181 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0181 is closed with `flai issue close I-0181 --reason` saying what fixed it

## Tasks
- T-2134 Show that I-0181's race is the one S-0494 removed: the pre-fix module fails the server-stop test and the current one holds under repeated runs
- T-2135 Close I-0181 citing S-0494's teardown, name it in the dashboard design's stop bullet, and tick the criteria

## Notes

Cost of delay inputs set by flai from I-0181. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T11:42:36Z, 4.7 days before this story; under one cycle counts as one).

### Planning

The proposed solution: I-0181's one instance (2026-10-21, S-0460's close-out) is the same race I-0205 recorded on 2026-10-24 and S-0494 fixed the same day in commit `d207daaf3`. `EventsSocket.close()` used to send 1001 and leave the Repo listeners to the socket's own `close` event, which fires only when the close handshake ends, and neither the client's 1001 nor `server.close` waits for an upgraded socket. Since S-0494 `close()` runs each socket's teardown before its 1001, and the server-stop test asserts no listener right after `close()`, which fails against the pre-fix module. So the story builds no second fix: it shows the pre-fix module fails as I-0181 describes and the current one holds under repeated runs, names I-0181 beside I-0205 where the test and the design cite the fix, and closes I-0181 citing S-0494. If the repeated runs still fail, the story's agent fixes `events-socket.ts` and says so in the narrative.

Touches: `flaiover/src/lib/server/events-socket.test.ts` (design: the failing test, and where I-0181 is named beside I-0205), `flaiover/src/lib/server/events-socket.ts` (layout: the module under test, changed only if the repeated runs still fail), `design/issues/I-0181-...md` and `design/issues/summary.md` (criteria: `flai issue close` writes both), `design/system/flaiover-dashboard.md` (co-change, 2 of 3 commits: its SIGTERM bullet cites I-0205 and S-0494, and gains I-0181). `flai touches suggest` also listed `flaiover/server.js` (2 of 3): left out, since the stop path there calls the router's `close()` and needs no change.

Forecast: `flai forecast` gave 14m, the median 117 s per unit over 32 done remediation stories on claude-opus-5-5 in the medium band at size 7 (2 criteria, 5 touches), and delivery 2026-10-26T17:52Z at the 21st place in the pull order behind S-0498 and S-0504. Set to 30m, delivery 2026-10-26T18:08Z: the work runs the pre-fix module and the current test file repeatedly, and the close-out's flaiover tier dominates a story this small.

Cost of delay: 7.50 USD a week from the operator's input, 3m lost per 168h cycle at 150 USD an hour. Stands: one occurrence so far, and S-0494 may already have removed it.

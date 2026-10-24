---
id: I-0205
title: flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-10-24T13:51:23Z
last_reported: 2026-10-24T13:51:23Z
updated: 2026-10-24T13:54:15Z
---

# I-0205 flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI

## Description
flaiover's events socket removes its Repo listeners only when the close handshake completes, so the server-stop test races and fails in CI

## Instances

### 2026-10-24T13:51:23Z
CI's flaiover test step failed once on main at src/lib/server/events-socket.test.ts 'closes every events socket with 1001 when the server stops, so the server can close': harbour.listenerCount('change') was 1 after the client saw the 1001 close and server.close had called back. EventsUpgrade.close() (src/lib/server/events-socket.ts) only calls ws.close(1001) on each open socket; the Repo listeners come off in the socket's own 'close' event, which fires after the close handshake completes, and neither the client's close nor http server.close waits for that, since an upgraded socket is detached from the http server. Locally the test passes five of five; on the CI runner the handshake lost the race. Fix: have close() tear each socket down at once (clear its ping, remove its Repo listeners, forget it) before or as it sends 1001, so stopping leaves no listener behind whatever the handshake does; the test then holds deterministically, and the sibling test's vi.waitFor pattern is a fallback.

## Impact
One CI failure on 2026-10-24 on a main push; each costs a re-run or a bisect of the merge queue's batch

- time_lost_per_cycle: 15m

## Remediation

Story S-0494 remediates this issue, created from it at 2026-10-24T13:51:34Z.
Closed 2026-10-24T13:54:15Z: Fixed by S-0494: EventsSocket.close() runs each open socket's teardown (ping cleared, Repo listeners off, socket forgotten) before it sends 1001, instead of leaving it to the socket's close event after the handshake, so stopping leaves no listener whatever the handshake does. The server-stop test in flaiover/src/lib/server/events-socket.test.ts asserts no listener right after close(), before the client sees 1001; it fails against main's events-socket.ts.

---
id: I-0181
title: "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed"
class: defect
status: closed
count: 1
cost: 3m
first_reported: 2026-10-21T11:42:36Z
last_reported: 2026-10-21T11:42:36Z
updated: 2026-10-28T12:49:35Z
---

# I-0181 events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Description
events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Instances

### 2026-10-21T11:42:36Z
Story: S-0460.
S-0460's close-out failed the flaiover tier on flaiover/src/lib/server/events-socket.test.ts:193, expect(harbour.listenerCount('change')).toBe(0) got 1, though the branch changes no server code. The listener seems to be removed on the server's socket close, which can come after the client sees 1001 and the server's close callback. A second close-out run, unchanged, was started to confirm.

## Remediation

Story S-0499 remediates this issue, created from it at 2026-10-26T04:01:24Z.
Closed 2026-10-28T12:49:35Z: Fixed by S-0494 (commit d207daaf3), the fix for I-0205, the same race: EventsSocket.close() now runs each socket's teardown, taking its Repo listeners off, before it sends 1001, rather than leaving them to the socket's close event, which the server's close does not wait for. S-0499 showed the module before d207daaf3 fails the server-stop test in events-socket.test.ts with a change listener still attached, as this issue's instance did, and the current module passed it 12 of 12 runs.

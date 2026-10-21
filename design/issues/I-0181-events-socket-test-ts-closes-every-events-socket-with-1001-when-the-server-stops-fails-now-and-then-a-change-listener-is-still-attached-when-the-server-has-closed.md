---
id: I-0181
title: "events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed"
class: defect
status: open
count: 1
cost: 3m
first_reported: 2026-10-21T11:42:36Z
last_reported: 2026-10-21T11:42:36Z
updated: 2026-10-21T11:42:36Z
---

# I-0181 events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Description
events-socket.test.ts 'closes every events socket with 1001 when the server stops' fails now and then: a change listener is still attached when the server has closed

## Instances

### 2026-10-21T11:42:36Z
Story: S-0460.
S-0460's close-out failed the flaiover tier on flaiover/src/lib/server/events-socket.test.ts:193, expect(harbour.listenerCount('change')).toBe(0) got 1, though the branch changes no server code. The listener seems to be removed on the server's socket close, which can come after the client sees 1001 and the server's close callback. A second close-out run, unchanged, was started to confirm.

## Remediation

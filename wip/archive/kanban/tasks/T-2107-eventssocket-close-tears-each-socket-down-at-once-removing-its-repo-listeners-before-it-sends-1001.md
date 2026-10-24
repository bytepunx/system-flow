---
id: T-2107
type: task
nature: remediation
title: EventsSocket.close tears each socket down at once, removing its Repo listeners before it sends 1001
status: done
parent: S-0494
owner: alex
created: 2026-10-24T13:52:30Z
updated: 2026-10-24T13:54:00Z
transitions:
  - to: ready
    at: 2026-10-24T13:52:44Z
    by: agent-S-0494
  - to: in-progress
    at: 2026-10-24T13:52:45Z
    by: agent-S-0494
  - to: done
    at: 2026-10-24T13:54:00Z
    by: agent-S-0494
stream: S-0494
tags: []
touches: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 75
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 3196
      cache_read: 1088232
      cache_write: 45020
      cost: 0.6418
---
# T-2107 EventsSocket.close tears each socket down at once, removing its Repo listeners before it sends 1001

## Work

Move each socket's teardown (clear its ping, take its `change` and `agent` listeners off the Repo, forget it) out of its `close` handler into one idempotent function kept per open socket. `close()` runs it for every socket before it sends 1001, and the socket's own `close` event runs it too, for a socket that closes on its own. Waits for no task: it is the story's only one.

## Done when

- `EventsSocket.close()` leaves no listener on the Repo when it returns, whatever the close handshake does.
- The server-stop test in `events-socket.test.ts` asserts that right after `close()`, before the client sees 1001, and fails against main's `events-socket.ts`.
- `flai test` on the two files passes.

## Notes

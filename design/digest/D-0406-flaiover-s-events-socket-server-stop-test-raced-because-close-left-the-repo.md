---
id: D-0406
title: flaiover's events socket server-stop test raced because close() left the Repo listeners to each socket's close event, so stopping now tears every socket down before it sends 1001
kind: cause
sources: [I-0205, S-0494]
topics: [flaiover, dashboard, sveltekit, code]
paths: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0406 flaiover's events socket server-stop test raced because close() left the Repo listeners to each socket's close event, so stopping now tears every socket down before it sends 1001

I-0205 was closed on 2026-10-24.

The cause: flaiover's events socket removed its `Repo` listeners only when a socket's close handshake completed, so the test "closes every events socket with 1001 when the server stops" raced the handshake and failed now and then in CI.

The fix, by S-0494:

```text
EventsSocket.close() runs each open socket's teardown (ping cleared, Repo listeners off, socket forgotten) before it sends 1001, instead of leaving it to the socket's close event after the handshake, so stopping leaves no listener whatever the handshake does. The server-stop test in flaiover/src/lib/server/events-socket.test.ts asserts no listener right after close(), before the client sees 1001.
```

ADR-0182 § Decision still says a socket stops listening to the `Repo` when it closes; since S-0494 the server's own stop does it first. Keep the teardown in `close()` when changing the socket's lifecycle.

Sources: I-0205, S-0494.

---
id: D-0406
title: flaiover's events socket server-stop test raced because close() left the Repo listeners to each socket's close event, so stopping now tears every socket down before it sends 1001
kind: cause
sources: [I-0205, S-0494, I-0181, S-0499]
topics: [flaiover, dashboard, sveltekit, code, testing]
paths: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts]
date: 2026-10-28
status: absorbed
absorbed_by: "design/system/flaiover-dashboard.md#the-channel-to-flai-on-the-host-s-0072-adr-0029"
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

I-0181 reported the same flaky test, with a change listener still attached when the server had closed. S-0499 confirmed it is the race S-0494 removed and made no second fix. From its narrative's Decisions (`wip/archive/agents/S-0499.md`):

```text
No second fix: the module before `d207daaf3` fails the server-stop test with a change listener still attached, and the current one passed it 12 of 12 runs. S-0494's assertion right after `close()` is the test that reproduces I-0181, so criterion 1 needed no new test.
```

If this test fails again, look first for a change that moved the teardown out of `close()`; that assertion is the reproducing test.

ADR-0182 § Decision still says a socket stops listening to the `Repo` when it closes; since S-0494 the server's own stop does it first. Keep the teardown in `close()` when changing the socket's lifecycle.

Sources: I-0205, S-0494, I-0181, S-0499.

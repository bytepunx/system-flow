---
id: T-2134
type: task
nature: remediation
title: "Show that I-0181's race is the one S-0494 removed: the pre-fix module fails the server-stop test and the current one holds under repeated runs"
status: done
parent: S-0499
owner: alex
created: 2026-10-26T04:20:54Z
updated: 2026-10-28T12:49:14Z
transitions:
  - to: ready
    at: 2026-10-28T12:45:25Z
    by: agent-S-0499
  - to: in-progress
    at: 2026-10-28T12:45:26Z
    by: agent-S-0499
  - to: done
    at: 2026-10-28T12:49:14Z
    by: agent-S-0499
stream: S-0499
tags: [dashboard, tests]
touches: [flaiover/src/lib/server/events-socket.test.ts, flaiover/src/lib/server/events-socket.ts]
usage:
  source: log
  seconds: 228
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 12
      output: 61
      cache_read: 674418
      cache_write: 5272
      cost: 0.2299
---
# T-2134 Show that I-0181's race is the one S-0494 removed: the pre-fix module fails the server-stop test and the current one holds under repeated runs

## Work

I-0181's instance (S-0460's close-out, 2026-10-21) failed `flaiover/src/lib/server/events-socket.test.ts` 'closes every events socket with 1001 when the server stops' with one `change` listener left on the Repo after `server.close` had called back. That is the race I-0205 recorded on 2026-10-24 and S-0494 fixed in commit `d207daaf3`: `EventsSocket.close()` left the Repo listeners to the socket's own `close` event, after a close handshake the server does not wait for. Since then `close()` runs each socket's teardown before its 1001, and the test asserts no listener right after `close()`.

- Reproduce: put the pre-fix module, `git show c6843b966:flaiover/src/lib/server/events-socket.ts`, in place of the current one in the worktree, run `pnpm vitest run src/lib/server/events-socket.test.ts` in `flaiover`, and record that the server-stop test fails on the listener count as I-0181 describes. Restore the current module afterwards; `git status` must show it unchanged.
- Verify: run the same file at least 20 times on the current module (a shell loop, or vitest's `--repeat`), and record the pass count in the narrative's `## Decisions`.
- In the test's comment on the teardown assertion, name I-0181 beside I-0205, so the test says which issues it guards.
- If a repeated run still fails on the current module, the fix is not complete: find the remaining path that leaves a listener (an upgrade whose `watch()` resolves after `close()`, or a socket closed by the client while the server stops), fix it in `events-socket.ts`, and add the test that reproduces it. Say in the narrative which it was.

## Done when

- The narrative's `## Decisions` records the pre-fix failure and the number of passing repeated runs on the current module.
- The test's comment names I-0181 and I-0205.
- `events-socket.ts` is unchanged, or changed with a new test that fails without the change, and the reason is in the narrative.
- `flai test flaiover/src/lib/server/events-socket.test.ts` passes.

## Notes

Waits for nothing: the first layer.

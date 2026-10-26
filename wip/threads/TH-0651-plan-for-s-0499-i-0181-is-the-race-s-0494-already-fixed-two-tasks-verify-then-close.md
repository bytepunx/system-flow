---
id: TH-0651
title: "Plan for S-0499: I-0181 is the race S-0494 already fixed; two tasks, verify then close"
anchor:
  path: wip/kanban/stories/S-0499-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md
  item: S-0499
status: resolved
participants: [planner-S-0499, orchestrator]
created: 2026-10-26T04:21:24Z
updated: 2026-10-26T04:21:34Z
---

# TH-0651 Plan for S-0499: I-0181 is the race S-0494 already fixed; two tasks, verify then close

On wip/kanban/stories/S-0499-events-socket-test-ts-closes-every-events-socket-with-1001-when-the-server-stops-fails-now-and-then-a-change-listener-is-still-attached-when-the-server-has-closed.md.

## Entries

### 2026-10-26T04:21:24Z planner-S-0499
Proposed solution: I-0181's one instance (S-0460's close-out, 2026-10-21) is the same race I-0205 recorded on 2026-10-24 and S-0494 fixed the same day in commit d207daaf3. EventsSocket.close() used to leave the Repo listeners to the socket's own close event, after a close handshake the server does not wait for; since S-0494 close() tears each socket down before its 1001, and the server-stop test asserts no listener right after close(), which fails against the pre-fix module. The story builds no second fix unless repeated runs still fail.

Tasks and layers:

- Layer 1: T-2134 reproduces I-0181 on the pre-fix module (git show c6843b966:flaiover/src/lib/server/events-socket.ts), runs the test file at least 20 times on the current module, names I-0181 beside I-0205 in the test's comment, and fixes events-socket.ts with a new test only if a run still fails. Touches: events-socket.test.ts, events-socket.ts.
- Layer 2: T-2135, after T-2134, names I-0181 in design/system/flaiover-dashboard.md's SIGTERM bullet beside I-0205 and S-0494, closes I-0181 with flai issue close citing S-0494 and T-2134's evidence, and ticks the criteria. Touches: the I-0181 file, design/issues/summary.md, flaiover-dashboard.md.

Figures: touches from the issue, the S-0494 commit, and co-change (flaiover-dashboard.md kept, flaiover/server.js left out: its stop path only calls the router's close()). Forecast raised from flai's 14m to 30m for the repeated vitest runs and the pre-fix reproduction; delivery 2026-10-26T18:08Z. Cost of delay 7.50 USD a week from the operator's 3m input, unchanged.

Assumptions: the current module has no residual listener path (serve() closes at once while stopping, and an upgrade whose watch() resolves after close() never attaches listeners), so T-2134's verification is expected to pass; the criterion's "test that reproduces it" is met by S-0494's assertion, which I-0181's close reason cites rather than a new test. Nothing to split, merge, or drop.

### 2026-10-26T04:21:30Z orchestrator
Approved. The plan fits S-0499.

- **The remedy is grounded.** I-0181's single instance is the race that I-0205 recorded and S-0494 fixed in `d207daaf3`. T-2134 shows it by failing the pre-fix module and holding the current one over at least 20 runs. It fixes `events-socket.ts` only if a run still fails, so no second fix is built on a guess.
- **The tasks fit.** There are two tasks in two layers, T-2135 after T-2134. The touches are five files and no folder. `flaiover/server.js` is left out with a reason.
- **The figures fit.** The forecast is 30m, raised from flai's 14m for the repeated runs. The value is 7.50 USD/week from `flai cod` on the issue's 3m input. I set no input.
- **The criterion's test** is S-0494's assertion, which T-2134 shows fails on the pre-fix module. The close reason cites it.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:21:34Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0499.

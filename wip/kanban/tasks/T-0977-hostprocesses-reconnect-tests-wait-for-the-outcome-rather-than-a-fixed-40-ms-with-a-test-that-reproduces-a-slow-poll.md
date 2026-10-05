---
id: T-0977
type: task
nature: remediation
title: HostProcesses' reconnect tests wait for the outcome rather than a fixed 40 ms, with a test that reproduces a slow poll
status: backlog
parent: S-0264
owner: alex
created: 2026-10-05T05:49:37Z
updated: 2026-10-05T05:49:37Z
transitions: []
stream: S-0264
tags: [flaiover, testing]
touches: [flaiover/src/lib/components/HostProcesses.svelte.test.ts]
---
# T-0977 HostProcesses' reconnect tests wait for the outcome rather than a fixed 40 ms, with a test that reproduces a slow poll

## Work

The cause of I-0053: the reconnect tests in `flaiover/src/lib/components/HostProcesses.svelte.test.ts` run the component on real timers cut short by the `fast` props (`disconnectTimeoutMs: 5, reconnectPollMs: 5, reconnectGiveUpMs: 50`), then assert after a fixed wall-clock wait (`settleThrough(40)`, `settleThrough()` of 20 ms). Under machine load a 5 ms `setTimeout` can stretch well past 40 ms, so the assertion runs before the second poll has answered; and `waitForReconnect` measures its deadline with `Date.now()`, so a stretched poll can also run past the 50 ms give-up and report "did not come back on a newer flai in time" instead.

- First add a test that reproduces it: make one poll's answer arrive late, longer than the old 40 ms wait (for example a mock that resolves after a `setTimeout` of 60 ms, standing in for a loaded machine), and see it fail with the waits as they are.
- Then make every test that waits on a reconnect (the upgrade tests, the serve restart that drops the connection, the 502 restart) wait for its outcome with `vi.waitFor` (or an equivalent poll with a timeout of a second or more) instead of a fixed `settleThrough`, and give `reconnectGiveUpMs` a margin of seconds in the props those tests pass, so the deadline no longer races the load. Keep the poll interval short so the tests stay fast.
- Keep a test of the give-up message working: if one is added or exists, give it its own short `reconnectGiveUpMs` with polls that never satisfy it.
- The component `HostProcesses.svelte` already takes its timings as props and needs no change; if the reproduction shows otherwise, say why in the narrative's `## Decisions` and add its path to this task's touches.
- `HostPanel.svelte.test.ts` uses the same pattern; it is outside this task unless the operator says otherwise on the plan's thread.

This task waits for none.

## Done when

- [ ] The reproduction test is in the file, and it failed against the fixed waits before the change.
- [ ] No reconnect test in the file asserts after a fixed wall-clock wait, and none can reach the give-up deadline on a slow poll.
- [ ] `scripts/flaiover-unit.sh` (or `make flaiover-test`) passes, with the file's tests run alone three times in a row green.

## Notes

Drafted by the planner for S-0264. The HostPanel test's comment explains why real timers were chosen over `vi.useFakeTimers()`; waiting for the outcome keeps real timers and removes the race instead.

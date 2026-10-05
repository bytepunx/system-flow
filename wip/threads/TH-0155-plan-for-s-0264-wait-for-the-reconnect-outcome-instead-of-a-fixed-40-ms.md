---
id: TH-0155
title: "Plan for S-0264: wait for the reconnect outcome instead of a fixed 40 ms"
anchor:
  path: wip/kanban/stories/S-0264-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md
  item: S-0264
status: resolved
participants: [planner-S-0264, alex]
created: 2026-10-05T05:50:26Z
updated: 2026-10-05T05:52:26Z
---

# TH-0155 Plan for S-0264: wait for the reconnect outcome instead of a fixed 40 ms

On wip/kanban/stories/S-0264-flaiover-s-hostprocesses-test-that-waits-past-the-old-host-after-an-upgrade-fails-under-machine-load.md.

## Entries

### 2026-10-05T05:50:26Z planner-S-0264
Plan for S-0264 (draft, backlog).

Cause, from I-0053's one instance: the reconnect tests in HostProcesses.svelte.test.ts run on real timers cut to 5 ms polls and a 50 ms give-up, then assert after a fixed 40 ms wall-clock wait. Under load a 5 ms timer can stretch past 40 ms, so the assertion runs before the second poll answers. The Date.now() deadline can also pass first, and then the test reports a give-up. Proposed fix: wait for the outcome with vi.waitFor and give the give-up a margin of seconds, keeping real timers. The test file is the only code change.

Tasks and layers:
- Layer 1: T-0977. Add a test that reproduces a slow poll and fails with the current waits. Then make every reconnect test wait for its outcome (touches HostProcesses.svelte.test.ts).
- Layer 2: T-0979, after T-0977. Run the unit tier while the Go tiers run, at least twice. Then close I-0053 from the worktree (touches I-0053 and design/issues/summary.md).

Figures:
- Forecast: 30m, raised from flai's 12m for the runs under load. Delivery is 2026-10-05T23:45Z.
- Cost of delay: 10 USD a week, from the 4m input, unchanged.

Assumptions:
- HostProcesses.svelte needs no change, because its timings are already props.
- The co-change paths, the /api/host route and its test, are left out.

Proposal: HostPanel.svelte.test.ts has the same pattern (fast 5/5/50 props and fixed settleThrough waits, 6 uses) but has no recorded instance. I recommend leaving it out of this story and recording it as an issue if it ever flakes. If you want it fixed now, I can add a third task to layer 1, since it touches a different file.

### 2026-10-05T05:52:26Z alex
Resolved.
